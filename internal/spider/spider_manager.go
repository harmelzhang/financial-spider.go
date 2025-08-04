package spider

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/panjf2000/ants/v2"
	"harmel.cn/financial/internal/model"
	"harmel.cn/financial/internal/public"
	"harmel.cn/financial/internal/service"
	"harmel.cn/financial/internal/spider/response"
	"harmel.cn/financial/utils/http"
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
			// 删除数据库中指定类型的分类数据
			err := service.CategoryService.DeleteByType(ctx, typeName)
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
			// 删除数据库中指定类型的分类关系数据
			err := service.CategoryStockCodeService.DeleteByType(ctx, typeName)
			if err != nil {
				g.Log("spider").Warningf(ctx, "delete all categroy stock code failed, err is %v", err)
			}
			// 插入新数据
			for _, stock := range stockCodeRes.Data.List {
				var categoryType, categoryCode string
				if stock.CicsLeve1Code != "" {
					// 中证
					categoryType = "CICS"
					if stock.CicsLeve4Code == "99999999" {
						continue
					}
					categoryCode = stock.CicsLeve4Code
				} else {
					// 证监会
					categoryType = "CSRC"
					if stock.CsrcLeve2Code == "" {
						// FIX 证券会暂时没对新三板股票进行分类，后续待优化
						continue
					}
					categoryCode = stock.CsrcLeve1Code + stock.CsrcLeve2Code
				}
				categoryStockCode := &model.CategoryStockCode{
					Type:         categoryType,
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
			ParentCode:   category.ParentId,
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

// 计算财务比率
func (s *SpiderManager) calcFinancialRatios() error {
	// TODO
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
