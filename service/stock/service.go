package stock

import (
	"encoding/json"
	cConfig "financial-spider.go/config/category"
	fConfig "financial-spider.go/config/financial"
	sConfig "financial-spider.go/config/stock"
	"financial-spider.go/models"
	"financial-spider.go/models/vo"
	"financial-spider.go/utils/db"
	"financial-spider.go/utils/http"
	"financial-spider.go/utils/tools"
	"fmt"
	"log"
	"strings"
)

// QueryStockBaseInfo 查询股票基本信息
func QueryStockBaseInfo(code string) {
	log.Println("查询股票基本信息")

	marketName, marketShortName := QueryStockMarketPlace(code)

	stockBaseInfoRes := vo.StockBaseInfoResult{}

	url := fmt.Sprintf(sConfig.QueryStockBaseInfoUrl, marketShortName, code)
	err := json.Unmarshal(http.Get(url), &stockBaseInfoRes)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	baseInfo := stockBaseInfoRes.BaseInfo[0]
	listingInfo := stockBaseInfoRes.ListingInfo[0]

	stock := models.Stock{
		Code:                code,
		StockName:           baseInfo.StockName,
		StockNamePinyin:     tools.GetPinyinFirstWord(baseInfo.StockName.(string)),
		StockBeforeName:     baseInfo.StockBeforeName,
		CompanyName:         baseInfo.CompanyName,
		CompanyProfile:      strings.Trim(baseInfo.CompanyProfile.(string), " "),
		Region:              baseInfo.Region,
		Address:             baseInfo.Address,
		Website:             baseInfo.Website,
		BusinessScope:       baseInfo.BusinessScope,
		DateOfIncorporation: listingInfo.DateOfIncorporation,
		ListingDate:         listingInfo.ListingDate,
		LawFirm:             baseInfo.LawFirm,
		AccountingFirm:      baseInfo.AccountingFirm,
		MarketPlace:         marketName,
	}

	stockMainBusinessResult := vo.StockMainBusinessResult{}
	url = fmt.Sprintf(sConfig.QueryStockMainBusinessUrl, code)
	err = json.Unmarshal(http.Get(url), &stockMainBusinessResult)
	if err != nil {
		log.Printf("解析JSON出错，跳过主营业务 : %s", err)
	}

	if stockMainBusinessResult.Code == 0 && stockMainBusinessResult.Success {
		stock.MainBusiness = stockMainBusinessResult.Result.Data[0].Info
	} else {
		log.Printf("获取主营业务数据失败，跳过查询 > Code:%d, Msg: %s", stockMainBusinessResult.Code, stockMainBusinessResult.Message)
	}

	// 插入或修改数据
	stock.ReplaceData()
}

// QueryStockFinancialData 查询股票对应公司财报数据
func QueryStockFinancialData(code string) {
	allReportDates, _ := queryAllReportDate(code)

	log.Println("查询股票对应公司财报数据")

	// 初始化操作
	financials := make([]*models.Financial, 0)
	for _, reportDate := range allReportDates {
		financial := models.NewFinancial(code, reportDate)
		financial.InitData()

		financials = append(financials, financial)
	}

	reportDatePages, pageTotal := tools.ArraySlice(allReportDates, fConfig.QueryPageSize) // 分页查询，减少请求量
	for i, reportDates := range reportDatePages {
		log.Printf("处理报表进度 : %d / %d", i+1, pageTotal)

		queryDates := strings.Join(reportDates, ",")

		processingCashFlowSheet(financials, code, queryDates)
		processingIncomeSheet(financials, code, queryDates)
		processingBalanceSheet(financials, code, queryDates)
	}

	// 财报数据入库（也可以在每个处理函数中进行，但是会增加数据库的操作，好处是每次都能更新一部分数据）
	for _, financial := range financials {
		financial.UpdateData()
	}

	processingDividend(code)

	calcFinancialRatio(code)
}

// QueryStockMarketPlace 查询股票交易市场名称和简称（SH、SZ、BJ）
func QueryStockMarketPlace(code string) (string, string) {
	name, shortName := "", ""
	stockCodePrefix := code[0:2]
	if tools.IndexOf(cConfig.ShanghaiMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "上海", "SH"
	} else if tools.IndexOf(cConfig.ShenzhenMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "深圳", "SZ"
	} else if tools.IndexOf(cConfig.BeijingMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "北京", "BJ"
	}
	return name, shortName
}

// 根据股票代码和报告期查询财报
func queryFinancialByCodeAndReportDate(financials []*models.Financial, code string, reportDate string) *models.Financial {
	for _, financial := range financials {
		if financial.Code == code && financial.ReportDate == reportDate {
			return financial
		}
	}
	return nil
}

// 查询所有报告期
func queryAllReportDate(code string) ([]string, int) {
	result := make([]string, 0)

	_, marketShortName := QueryStockMarketPlace(code)

	reportDateResult := vo.StockReportDateResult{}

	insertDate := func() {
		for _, reportDate := range reportDateResult.Data {
			date := strings.Split(reportDate.Date, " ")[0]
			if tools.IndexOf(result, date) == -1 {
				result = append(result, date)
			}
		}
	}

	log.Println("查询资产负债表报告期")
	url := fmt.Sprintf(sConfig.QueryBalanceSheetReportDateUrl, marketShortName, code)
	err := json.Unmarshal(http.Get(url), &reportDateResult)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}
	insertDate()

	log.Println("查询利润表报告期")
	url = fmt.Sprintf(sConfig.QueryIncomeSheetReportDateUrl, marketShortName, code)
	err = json.Unmarshal(http.Get(url), &reportDateResult)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}
	insertDate()

	log.Println("查询现金流量表报告期")
	url = fmt.Sprintf(sConfig.QueryCashFlowSheetReportDateUrl, marketShortName, code)
	err = json.Unmarshal(http.Get(url), &reportDateResult)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}
	insertDate()

	return result, len(result)
}

// 处理现金流量表
func processingCashFlowSheet(financials []*models.Financial, code string, queryDates string) {
	_, marketShortName := QueryStockMarketPlace(code)

	log.Println("查询现金流量表数据")
	url := fmt.Sprintf(fConfig.QueryCashFlowSheetUrl, queryDates, marketShortName, code)
	cashFlowSheetResult := vo.FinancialResult{}
	err := json.Unmarshal(http.Get(url), &cashFlowSheetResult)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if cashFlowSheetResult.Type == "1" || cashFlowSheetResult.Status == 1 {
		log.Printf("跳过查询，没有该期报表数据或参数异常 [%s %s]", code, queryDates)
		return
	}

	if len(cashFlowSheetResult.Data) != 0 {
		for _, cashFlowSheetData := range cashFlowSheetResult.Data {
			reportDate := strings.Split(cashFlowSheetData.ReportDate, " ")[0]

			financial := queryFinancialByCodeAndReportDate(financials, code, reportDate)
			financial.Ocf = cashFlowSheetData.Ocf
			financial.Cfi = cashFlowSheetData.Cfi
			financial.Cff = cashFlowSheetData.Cff
		}
	}
}

// 处理分红数据
func processingDividend(code string) {
	log.Println("查询分红数据")
	url := fmt.Sprintf(fConfig.QueryDividendUrl, code)
	dividendResult := vo.DividendResult{}
	err := json.Unmarshal(http.Get(url), &dividendResult)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if dividendResult.Code == 0 && dividendResult.Success {
		for _, dividend := range dividendResult.Result.Data {

			reportDate := dividend.Year + "-12-31"
			financial := models.NewFinancial(code, reportDate)
			financial.InitData()

			sql := "UPDATE financial SET dividend = ? WHERE code = ? AND report_date = ?"
			args := []interface{}{dividend.Money, code, reportDate}
			db.ExecSQL(sql, args...)
		}
	} else {
		log.Printf("获取分红数据失败，跳过查询 > Code:%d, Msg: %s", dividendResult.Code, dividendResult.Message)
	}

}

// 处理资产负债表
func processingBalanceSheet(financials []*models.Financial, code string, queryDates string) {
	_, marketShortName := QueryStockMarketPlace(code)

	log.Println("查询资产负债表数据")
	url := fmt.Sprintf(fConfig.QueryBalanceSheetUrl, queryDates, marketShortName, code)
	balanceSheet := vo.FinancialResult{}
	err := json.Unmarshal(http.Get(url), &balanceSheet)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if balanceSheet.Type == "1" || balanceSheet.Status == 1 {
		log.Printf("跳过查询，没有该期报表数据或参数异常 [%s %s]", code, queryDates)
		return
	}

	if len(balanceSheet.Data) != 0 {
		for _, balanceSheetData := range balanceSheet.Data {
			reportDate := strings.Split(balanceSheetData.ReportDate, " ")[0]

			financial := queryFinancialByCodeAndReportDate(financials, code, reportDate)
			financial.MonetaryFund = balanceSheetData.MonetaryFund
			financial.TradeFinassetNotfvtpl = balanceSheetData.TradeFinassetNotfvtpl
			financial.TradeFinasset = balanceSheetData.TradeFinasset
			financial.DeriveFinasset = balanceSheetData.DeriveFinasset

			financial.FixedAsset = balanceSheetData.FixedAsset
			financial.Cip = balanceSheetData.Cip

			financial.CaTotal = balanceSheetData.CaTotal
			financial.NcaTotal = balanceSheetData.NcaTotal
			financial.ClTotal = balanceSheetData.ClTotal
			financial.NclTotal = balanceSheetData.NclTotal
			financial.Inventory = balanceSheetData.Inventory
			financial.AccountsRece = balanceSheetData.AccountsRece
			financial.AccountsPayable = balanceSheetData.AccountsPayable
		}
	}
}

// 处理利润表
func processingIncomeSheet(financials []*models.Financial, code string, queryDates string) {
	_, marketShortName := QueryStockMarketPlace(code)

	log.Println("查询利润表数据")
	url := fmt.Sprintf(fConfig.QueryIncomeSheetUrl, queryDates, marketShortName, code)
	incomeSheet := vo.FinancialResult{}
	err := json.Unmarshal(http.Get(url), &incomeSheet)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if incomeSheet.Type == "1" || incomeSheet.Status == 1 {
		log.Printf("跳过查询，没有该期报表数据或参数异常 [%s %s]", code, queryDates)
		return
	}

	if len(incomeSheet.Data) != 0 {
		for _, incomeSheetData := range incomeSheet.Data {
			reportDate := strings.Split(incomeSheetData.ReportDate, " ")[0]

			financial := queryFinancialByCodeAndReportDate(financials, code, reportDate)
			financial.Np = incomeSheetData.Np
			financial.Oi = incomeSheetData.Oi
			financial.Coe = incomeSheetData.Coe
			financial.CoeTotal = incomeSheetData.CoeTotal
			financial.Eps = incomeSheetData.Eps
		}
	}

}

// 计算财务比率
func calcFinancialRatio(code string) {
	sql := `
		UPDATE financial
		SET
		    np_ratio = ROUND(np / oi * 100, 2),
		    dividend_ratio = ROUND(dividend / np * 100, 2),
		    oi_ratio = ROUND((oi - coe) / oi * 100, 2),
		    operating_profit_ratio = ROUND((oi - coe_total) / oi * 100, 2),
		    operating_safety_ratio = ROUND(operating_profit_ratio / oi_ratio * 100, 2),
		    cash_equivalent_ratio = ROUND((monetary_fund + IFNULL(IFNULL(trade_finasset, trade_finasset_notfvtpl), 0) + IFNULL(derive_finasset, 0)) / (ca_total + nca_total) * 100, 2),
		    ca_ratio = ROUND(ca_total / (ca_total + nca_total) * 100, 2),
		    cl_ratio = ROUND(cl_total / (ca_total + nca_total) * 100, 2),
		    ncl_ratio = ROUND(ncl_total / (ca_total + nca_total) * 100, 2),
		    debt_ratio = ROUND((cl_total + ncl_total) / (ca_total + nca_total) * 100, 2),
		    long_term_funds_ratio = ROUND((ncl_total + (ca_total + nca_total - cl_total - ncl_total)) / (fixed_asset + cip) * 100, 2),
		    equity_ratio = ROUND(100 - debt_ratio, 2),
		    inventory_ratio = ROUND(inventory / (ca_total + nca_total) * 100, 2),
		    accounts_rece_ratio = ROUND(accounts_rece / (ca_total + nca_total) * 100, 2),
		    accounts_payable_ratio = ROUND(accounts_payable / (ca_total + nca_total) * 100, 2),
		    current_ratio = ROUND(ca_total / cl_total * 100, 2),
		    quick_ratio = ROUND((ca_total - inventory) / cl_total * 100, 2),
		    roe = ROUND(np / (ca_total + nca_total - cl_total -ncl_total) * 100, 2),
		    roa = ROUND(np / (ca_total + nca_total) * 100, 2),
		    accounts_rece_turnover_ratio = ROUND(oi / accounts_rece, 2),
		    average_cash_receipt_days = ROUND(360 / accounts_rece_turnover_ratio, 2),
		    inventory_turnover_ratio = ROUND(coe / inventory, 2),
		    average_sales_days = ROUND(360 / inventory_turnover_ratio, 2),
		    total_asset_turnover_ratio = ROUND(oi / (ca_total + nca_total), 2),
		    cash_flow_ratio = ROUND(ocf / cl_total * 100, 2)
		WHERE code = ?
	`
	args := []interface{}{code}
	db.ExecSQL(sql, args...)
}
