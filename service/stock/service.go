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
	reportDates, reportCount := queryAllReportDate(code)

	log.Println("查询股票对应公司财报数据")

	for i, reportDate := range reportDates {
		log.Printf("处理报表进度 : %d / %d", i+1, reportCount)

		financial := models.NewFinancial(code, reportDate)
		financial.InitData()

		processingCashFlowSheet(financial)
		processingBalanceSheet(financial)
		processingIncomeSheet(financial)

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
func processingCashFlowSheet(financial *models.Financial) {
	_, marketShortName := QueryStockMarketPlace(financial.Code)

	log.Println("查询现金流量表数据")
	url := fmt.Sprintf(fConfig.QueryCashFlowSheetUrl, financial.ReportDate, marketShortName, financial.Code)
	cashFlowSheetResult := vo.FinancialResult{}
	err := json.Unmarshal(http.Get(url), &cashFlowSheetResult)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if cashFlowSheetResult.Type == "1" || cashFlowSheetResult.Status == 1 {
		log.Printf("跳过查询，没有该期报表数据或参数异常 [%s %s]", financial.Code, financial.ReportDate)
		return
	}

	if len(cashFlowSheetResult.Data) != 0 {
		cashFlowSheetData := cashFlowSheetResult.Data[0]
		financial.Ocf = cashFlowSheetData.Ocf
		financial.Cfi = cashFlowSheetData.Cfi
		financial.Cff = cashFlowSheetData.Cff
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
func processingBalanceSheet(financial *models.Financial) {
	_, marketShortName := QueryStockMarketPlace(financial.Code)

	log.Println("查询资产负债表数据")
	url := fmt.Sprintf(fConfig.QueryBalanceSheetUrl, financial.ReportDate, marketShortName, financial.Code)
	fmt.Println(url)
}

// 处理利润表
func processingIncomeSheet(financial *models.Financial) {
	_, marketShortName := QueryStockMarketPlace(financial.Code)

	log.Println("查询利润表数据")
	url := fmt.Sprintf(fConfig.QueryIncomeSheetUrl, financial.ReportDate, marketShortName, financial.Code)
	incomeSheet := vo.FinancialResult{}
	err := json.Unmarshal(http.Get(url), &incomeSheet)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if incomeSheet.Type == "1" || incomeSheet.Status == 1 {
		log.Printf("跳过查询，没有该期报表数据或参数异常 [%s %s]", financial.Code, financial.ReportDate)
		return
	}

	if len(incomeSheet.Data) != 0 {
		incomeSheetData := incomeSheet.Data[0]
		financial.Np = incomeSheetData.Np
	}

}

// 计算财务比率
func calcFinancialRatio(code string) {
	sql := `
		UPDATE financial
		SET
		    dividend_ratio = ROUND(dividend / np * 100, 2)
		WHERE code = ?
	`
	args := []interface{}{code}
	db.ExecSQL(sql, args...)
}
