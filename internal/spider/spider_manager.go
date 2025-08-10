package spider

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/panjf2000/ants/v2"
	"harmel.cn/financial/internal/model"
	"harmel.cn/financial/internal/public"
	"harmel.cn/financial/internal/service"
	"harmel.cn/financial/internal/spider/response"
	"harmel.cn/financial/utils/http"
	"harmel.cn/financial/utils/slice"
	"harmel.cn/financial/utils/tools"
	"harmel.cn/financial/utils/xls"
)

// 爬虫管理器
type SpiderManager struct {
	// 进度管理器
	progressManager *ProgressManager
	// 待处理任务通道
	taskChan chan PendingTask
	// 线程池
	pool *ants.Pool
	// 是否有在执行的任务
	hasTaskRunning atomic.Bool
	// 完成通知
	finishNotify chan bool
}

func NewSpiderManager(rootDir string) *SpiderManager {
	pool, err := ants.NewPool(public.SpiderExecutorPoolSize, ants.WithPreAlloc(true))
	if err != nil {
		panic(err)
	}

	return &SpiderManager{
		progressManager: NewProgressManager(rootDir),
		taskChan:        make(chan PendingTask, PENDING_TASKS_INIT_CAPACITY),
		pool:            pool,
		finishNotify:    make(chan bool),
	}
}

// 开启爬虫
func (s *SpiderManager) Start(ctx context.Context) (err error) {
	g.Log("spider").Debugf(ctx, "spider is running")

	// 加载历史进度
	err = s.progressManager.Load(ctx)
	if err != nil {
		g.Log("spider").Errorf(ctx, "ProgressManager.Load failed, err is %v", err)
		return err
	}

	// 如果到了五月一日，清空所有任务全部重跑（年报全部出了）
	if time.Now().Format("01-02") == "05-01" {
		s.progressManager.ClearTasks()
	}

	// 如果上次成功了，判断时间是否大于等于配置天数
	if s.progressManager.Done() {
		if time.Now().Unix()-s.progressManager.LastTS() >= public.SpiderTaskIntervalDays*24*3600 {
			s.progressManager.ClearTasks()
		} else {
			g.Log("spider").Debugf(ctx, "task finish, the time since the last successful task completion is less than %d days", public.SpiderTaskIntervalDays)
			return
		}
	}

	// 基础数据
	err = s.fetchIndexSample(ctx)
	if err != nil {
		g.Log("spider").Errorf(ctx, "fetch index sample data failed, err is %v", err)
		return
	}
	err = s.fetchCategory(ctx)
	if err != nil {
		g.Log("spider").Errorf(ctx, "fetch category data failed, err is %v", err)
		return
	}

	// 启动处理任务线程
	go s.doProcTaskWorker(ctx)

	// 等待完成
OUT:
	for {
		select {
		case <-s.finishNotify:
			g.Log("spider").Debug(ctx, "all task execute finish")
			break OUT
		default:
			time.Sleep(2 * time.Second)
			// 如果没有正在执行的线程并且通道里面没有待执行的任务
			if !s.hasTaskRunning.Load() && len(s.taskChan) == 0 {
				tasks := s.progressManager.UnexecutedTasks()
				if len(tasks) == 0 {
					s.progressManager.SetDone()
					err = s.progressManager.Save(ctx)
					if err != nil {
						g.Log("spider").Errorf(ctx, "save process failed, err is %v", err)
					}
					return
				}
				for _, task := range tasks {
					s.taskChan <- task
				}
			}
		}
	}

	err = s.calcFinancialRatios()
	if err != nil {
		g.Log("spider").Errorf(ctx, "calc financial ratios failed, err is %v", err)
		return
	}

	return
}

// 最新指数样本信息
func (s *SpiderManager) fetchIndexSample(ctx context.Context) error {
	for typeCode, _ := range public.IndexSampleType {
		g.Log("spider").Debugf(ctx, "start fetch %s index sample data", typeCode)

		// 请求数据
		url := fmt.Sprintf(public.UrlIndexSample, typeCode)
		client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
		body, _, err := client.Get(nil)
		if err != nil {
			g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
			continue
		}

		// 读取Excel
		items, err := xls.ReadXls(body, 0, 1)
		if err != nil {
			g.Log("spider").Errorf(ctx, "read xls failed, err is %v", err)
			continue
		}

		// 删除旧数据 & 插入新数据
		service.IndexSampleService.DeleteByType(ctx, typeCode)
		for _, item := range items {
			stockCode := item[4]
			indexSample := &model.IndexSample{
				TypeCode:  typeCode,
				StockCode: stockCode,
			}
			err = service.IndexSampleService.Insert(ctx, indexSample)
			if err != nil {
				g.Log("spider").Errorf(ctx, "insert index sample failed, TypeCode is %s StockCode is %s err is %v", typeCode, stockCode, err)
			}
		}

		g.Log("spider").Debugf(ctx, "fetch %s index sample data success", typeCode)
	}
	return nil
}

// 最新行业分类信息（含分类下的股票）
func (s *SpiderManager) fetchCategory(ctx context.Context) error {
	for typeName, typeValue := range public.CategoryType {
		g.Log("spider").Debugf(ctx, "start fetch %s catagory data", typeName)

		// 查询行业分类
		url := fmt.Sprintf(public.UrlCategory, typeValue)
		client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
		body, _, err := client.Get(nil)
		if err != nil {
			g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
			continue
		}

		categoryRes, err := http.ParseResponse[response.CategoryResult](body)
		if err != nil {
			g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
			continue
		}
		if categoryRes.Code == "200" && categoryRes.Success {
			// 删除数据库中指定类型的分类数据（同时会级联删除行业下股票信息）
			err = service.CategoryService.DeleteByType(ctx, typeName)
			if err != nil {
				g.Log("spider").Warningf(ctx, "delete category type %s data failed, err is %v", typeName, err)
				continue
			}
			// 递归插入新数据
			s.recursionCategorys(ctx, typeName, categoryRes.Data.MapList["4"])
		} else {
			g.Log("spider").Errorf(ctx, "fetch %s category data response error, code is %s", typeName, categoryRes.Code)
			continue
		}

		// 查询行业下的所有股票代码
		url = fmt.Sprintf(public.UrlCategoryStock, typeValue)
		client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
		body, _, err = client.Get(nil)
		if err != nil {
			g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
			continue
		}

		stockCodeRes, err := http.ParseResponse[response.StockCodeResult](body)
		if err != nil {
			g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
			continue
		}
		if stockCodeRes.Code == "200" && stockCodeRes.Success {
			// 插入新数据
			for _, stock := range stockCodeRes.Data.List {
				var categoryCode string
				if stock.CicsLeve1Code != "" {
					// 中证
					if stock.CicsLeve4Code == "99999999" {
						continue
					}
					categoryCode = stock.CicsLeve4Code
				} else {
					// 证监会
					if stock.CsrcLeve2Code == "" {
						// FIX 证券会暂时没对新三板股票进行分类，后续待优化
						continue
					}
					categoryCode = stock.CsrcLeve1Code + stock.CsrcLeve2Code
				}
				categoryStockCode := &model.CategoryStockCode{
					CategoryCode: categoryCode,
					StockCode:    stock.Code,
				}
				err := service.CategoryStockCodeService.Insert(ctx, categoryStockCode)
				if err != nil {
					g.Log("spider").Warningf(ctx, "insert categroy stock code failed, err is %v", err)
				}
				//  丢入任务列表
				task := PendingTask{Id: stock.Code}
				exist := s.progressManager.PutTask(task)
				if !exist {
					s.taskChan <- task
				}
			}
		} else {
			g.Log("spider").Errorf(ctx, "fetch %s category stock code data response error, code is %s", typeName, categoryRes.Code)
			continue
		}

		g.Log("spider").Debugf(ctx, "fetch %s catagory data success", typeName)
	}

	return nil
}

// 递归查询分类
func (s *SpiderManager) recursionCategorys(ctx context.Context, typeName string, categorys []response.Category) {
	if len(categorys) == 0 {
		return
	}

	for order, category := range categorys {
		mCategory := &model.Category{
			Type:         typeName,
			Code:         category.Id,
			Name:         category.Name,
			Level:        category.Level,
			DisplayOrder: order + 1,
		}
		if category.ParentId != "" {
			mCategory.ParentCode = category.ParentId
		}
		// 插入数据库
		err := service.CategoryService.Insert(ctx, mCategory)
		if err != nil {
			g.Log("spider").Warningf(ctx, "insert category data failed, err is %v", err)
			continue
		}
		if len(category.Children) != 0 {
			s.recursionCategorys(ctx, typeName, category.Children)
		}
	}
}

// TODO 计算财务比率
func (s *SpiderManager) calcFinancialRatios() error {
	return nil
}

// 处理任务
func (s *SpiderManager) doProcTaskWorker(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			g.Log("spiser").Critical(ctx, "doProcTaskWorker panic: %v", err)
		}
	}()

	for {
		task := <-s.taskChan
		err := s.executeTask(ctx, task)
		if err != nil {
			g.Log("spider").Errorf(ctx, "execute task failed, err is %v", err)
		}
	}
}

// 根据股票代码和报告期查询索引位置
func (s *SpiderManager) findFinancialIndex(stockCode, reportDate string, financials []*model.Financial) int {
	for idx, financial := range financials {
		if financial.StockCode == stockCode && financial.ReportDate == reportDate {
			return idx
		}
	}
	return -1
}

// 执行实际任务
func (s *SpiderManager) executeTask(ctx context.Context, task PendingTask) (err error) {
	err = s.pool.Submit(func() {
		s.hasTaskRunning.Store(true)
		defer func() {
			s.hasTaskRunning.Store(false)
		}()

		// 是否处理完
		isFinished, err := s.progressManager.TaskStatus(ctx, task.Id)
		if err != nil {
			g.Log("spider").Errorf(ctx, "query task status failed, err is %v", err)
			return
		}
		if isFinished {
			return
		}

		g.Log("spider").Debugf(ctx, "start execute task %s", task.Id)

		// 基本信息
		stock, err := s.fetchStockBaseInfo(ctx, task.Id)
		if err != nil {
			g.Log("spider").Errorf(ctx, "fetch stock %s base info failed, err is %v", stock.Code, err)
			return
		}

		// 查询所有报告期
		reportDates, err := s.queryAllReportData(ctx, stock)
		if err != nil {
			g.Log("spider").Errorf(ctx, "fetch stock %s report date info failed, err is %v", stock.Code, err)
			return
		}

		// 初始化操作
		financials := make([]*model.Financial, 0, len(reportDates))
		for _, reportDate := range reportDates {
			financial := &model.Financial{
				StockCode:  stock.Code,
				ReportDate: reportDate,
			}
			financials = append(financials, financial)
		}

		// 分页查询财报
		reportDatePages, totalPages := slice.ArraySlice(reportDates, public.QueryReportPageSize)
		for i, reportDates := range reportDatePages {
			g.Log("spider").Debugf(ctx, "fetch stock %s report info page %d/%d", stock.Code, i+1, totalPages)
			queryDates := strings.Join(reportDates, ",")
			// 现金流量表
			s.fetchCashFlowSheet(ctx, stock, queryDates, financials)
			// 资产负债表
			s.fetchBalanceSheet(ctx, stock, queryDates, financials)
			// 利润表
			s.fetchIncomeSheet(ctx, stock, queryDates, financials)
		}
		// 分红数据
		s.fetchDividendData(ctx, stock, financials)
		// TODO 插入或更新数据库

		// 标记完成
		s.progressManager.MarkTask(ctx, task.Id, true)

		// 写入磁盘
		err = s.progressManager.Save(ctx)
		if err != nil {
			g.Log("spider").Errorf(ctx, "save process failed, err is %v", err)
			return
		}

		g.Log("spider").Debugf(ctx, "task %s execute success", task.Id)

		// 通知完成
		if s.progressManager.Done() {
			s.finishNotify <- true
		}
	})
	return
}

// 查询股票市场
func (s *SpiderManager) queryStockMarketPlace(stockCode string) (string, string) {
	name, shortName := "", ""
	stockCodePrefix := stockCode[0:2]
	if slice.IndexOf(public.ShanghaiMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "上海", "SH"
	} else if slice.IndexOf(public.ShenzhenMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "深圳", "SZ"
	} else if slice.IndexOf(public.BeijingMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "北京", "BJ"
	}
	return name, shortName
}

// 基本信息
func (s *SpiderManager) fetchStockBaseInfo(ctx context.Context, stockCode string) (stock *model.Stock, err error) {
	marketName, marketShortName := s.queryStockMarketPlace(stockCode)

	// 公司类型
	url := fmt.Sprintf(public.UrlStockCompanyType, stockCode)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}

	companyType, companyTypeCode := "普通", "4"
	companyTypeRes, err := http.ParseResponse[response.CompanyTypeResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
	}
	if companyTypeRes.Success && companyTypeRes.Code == 0 {
		if companyTypeRes.Result.Count != 0 {
			companyType = companyTypeRes.Result.Data[0].Type
			companyTypeCode = companyTypeRes.Result.Data[0].TypeCode
		}
	}

	// 主营业务
	url = fmt.Sprintf(public.UrlStockMainBusiness, stockCode)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
	}

	mainBusinessResult, err := http.ParseResponse[response.MainBusinessResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}

	mainBusiness := ""
	if mainBusinessResult.Code == 0 && mainBusinessResult.Success {
		mainBusiness = mainBusinessResult.Result.Data[0].Info
	} else {
		g.Log("spider").Errorf(ctx, "fetch %s main business data response error, code is %d", stockCode, mainBusinessResult.Code)
	}

	// 主要信息
	url = fmt.Sprintf(public.UrlStockBaseInfo, marketShortName, stockCode)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}

	baseInfoRes, err := http.ParseResponse[response.StockBaseInfoResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	baseInfo := baseInfoRes.BaseInfo[0]
	listingInfo := baseInfoRes.ListingInfo[0]

	stock = &model.Stock{
		Code:            stockCode,
		Name:            baseInfo.Name,
		NamePinYin:      tools.PinyinFirstWord(baseInfo.Name),
		BeforeName:      baseInfo.BeforeName,
		CompanyName:     baseInfo.CompanyName,
		CompanyType:     companyType,
		CompanyTypeCode: companyTypeCode,
		CompanyProfile:  strings.TrimSpace(baseInfo.CompanyProfile),
		Region:          baseInfo.Region,
		Address:         baseInfo.Address,
		Website:         baseInfo.Website,
		MainBusiness:    mainBusiness,
		BusinessScope:   baseInfo.BusinessScope,
		CreateDate:      listingInfo.CreateDate,
		ListingDate:     listingInfo.ListingDate,
		LawFirm:         baseInfo.LawFirm,
		AccountingFirm:  baseInfo.AccountingFirm,
		MarketPlace:     marketName,
	}
	if stock.BeforeName != nil {
		stock.BeforeName = strings.ReplaceAll(fmt.Sprint(stock.BeforeName), "→", "、")
	}
	err = service.StockService.Replace(ctx, stock)
	if err != nil {
		g.Log("spider").Errorf(ctx, "replace db stock failed, err is %v", err)
		return
	}

	return
}

// 查询所有报告期
func (s *SpiderManager) queryAllReportData(ctx context.Context, stock *model.Stock) (reportDates []string, err error) {
	_, shortMarketName := s.queryStockMarketPlace(stock.Code)

	appendReportDate := func(reportDateRes *response.ReportDateResult) {
		for _, item := range reportDateRes.Data {
			date := strings.Split(item.Date, " ")[0]
			if slice.IndexOf(reportDates, date) == -1 {
				reportDates = append(reportDates, date)
			}
		}
	}

	// 资产负债表
	url := fmt.Sprintf(public.UrlBalanceSheetReport, stock.CompanyType, shortMarketName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	reportDateRes, err := http.ParseResponse[response.ReportDateResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	appendReportDate(reportDateRes)

	// 利润表
	url = fmt.Sprintf(public.UrlIncomeSheetReport, stock.CompanyType, shortMarketName, stock.Code)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	reportDateRes, err = http.ParseResponse[response.ReportDateResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	appendReportDate(reportDateRes)

	// 现金流量表
	url = fmt.Sprintf(public.UrlCashFlowSheetReport, stock.CompanyType, shortMarketName, stock.Code)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	reportDateRes, err = http.ParseResponse[response.ReportDateResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	appendReportDate(reportDateRes)

	return
}

// 现金流量表
func (s *SpiderManager) fetchCashFlowSheet(ctx context.Context, stock *model.Stock, queryDates string, financials []*model.Financial) {
	_, marketShortName := s.queryStockMarketPlace(stock.Code)
	url := fmt.Sprintf(public.UrlCashFlowSheet, stock.CompanyType, queryDates, marketShortName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	financialRes, err := http.ParseResponse[response.FinancialResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	if financialRes.Type == "1" || financialRes.Status == 1 {
		g.Log("spider").Errorf(ctx, "fetch %s cash flow sheet data response error, type is %s status is %d", stock.Code, financialRes.Type, financialRes.Status)
		return
	}

	for _, sheet := range financialRes.Data {
		reportDate := strings.Split(sheet.ReportDate, " ")[0]
		idx := s.findFinancialIndex(stock.Code, reportDate, financials)
		if idx == -1 {
			continue
		}
		financial := financials[idx]

		financial.Ocf = sheet.Ocf
		financial.Cfi = sheet.Cfi
		financial.Cff = sheet.Cff
		financial.AssignDividendPorfit = sheet.AssignDividendPorfit
		financial.AcquisitionAssets = sheet.AcquisitionAssets
		financial.InventoryLiquidating = sheet.InventoryLiquidating
	}
}

// 资产负债表
func (s *SpiderManager) fetchBalanceSheet(ctx context.Context, stock *model.Stock, queryDates string, financials []*model.Financial) {
	_, marketShortName := s.queryStockMarketPlace(stock.Code)
	url := fmt.Sprintf(public.UrlBalanceSheet, stock.CompanyType, queryDates, marketShortName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	financialRes, err := http.ParseResponse[response.FinancialResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	if financialRes.Type == "1" || financialRes.Status == 1 {
		g.Log("spider").Errorf(ctx, "fetch %s balance sheet data response error, type is %s status is %d", stock.Code, financialRes.Type, financialRes.Status)
		return
	}

	for _, sheet := range financialRes.Data {
		reportDate := strings.Split(sheet.ReportDate, " ")[0]
		idx := s.findFinancialIndex(stock.Code, reportDate, financials)
		if idx == -1 {
			continue
		}
		financial := financials[idx]

		financial.MonetaryFund = sheet.MonetaryFund
		financial.TradeFinassetNotfvtpl = sheet.TradeFinassetNotfvtpl
		financial.TradeFinasset = sheet.TradeFinasset
		financial.DeriveFinasset = sheet.DeriveFinasset

		financial.FixedAsset = sheet.FixedAsset
		financial.Cip = sheet.Cip

		financial.CaTotal = sheet.CaTotal
		financial.NcaTotal = sheet.NcaTotal
		financial.ClTotal = sheet.ClTotal
		financial.NclTotal = sheet.NclTotal
		financial.Inventory = sheet.Inventory
		financial.AccountsRece = sheet.AccountsRece
		financial.AccountsPayable = sheet.AccountsPayable
	}
}

// 利润表
func (s *SpiderManager) fetchIncomeSheet(ctx context.Context, stock *model.Stock, queryDates string, financials []*model.Financial) {
	_, marketShortName := s.queryStockMarketPlace(stock.Code)
	url := fmt.Sprintf(public.UrlIncomeSheet, stock.CompanyType, queryDates, marketShortName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	financialRes, err := http.ParseResponse[response.FinancialResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	if financialRes.Type == "1" || financialRes.Status == 1 {
		g.Log("spider").Errorf(ctx, "fetch %s balance sheet data response error, type is %s status is %d", stock.Code, financialRes.Type, financialRes.Status)
		return
	}

	for _, sheet := range financialRes.Data {
		reportDate := strings.Split(sheet.ReportDate, " ")[0]
		idx := s.findFinancialIndex(stock.Code, reportDate, financials)
		if idx == -1 {
			continue
		}
		financial := financials[idx]

		financial.Np = sheet.Np
		financial.Oi = sheet.Oi
		financial.Coe = sheet.Coe
		financial.CoeTotal = sheet.CoeTotal
		financial.Eps = sheet.Eps
	}
}

// 分红数据
func (s *SpiderManager) fetchDividendData(ctx context.Context, stock *model.Stock, financials []*model.Financial) {
	url := fmt.Sprintf(public.UrlDividend, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	dividendRes, err := http.ParseResponse[response.DividendResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	if dividendRes.Code == 0 && dividendRes.Success {
		for _, dividend := range dividendRes.Result.Data {
			reportDate := dividend.Year + "-12-31"
			idx := s.findFinancialIndex(stock.Code, reportDate, financials)
			if idx == -1 {
				continue
			}
			financial := financials[idx]
			financial.Dividend = dividend.Money
		}
	} else {
		g.Log("spider").Errorf(ctx, "fetch %s dividend data response error, code is %s message is %s", stock.Code, dividendRes.Code, dividendRes.Message)
		return
	}
}
