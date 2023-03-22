package stock

import (
	"encoding/json"
	cConfig "financial-spider.go/config/category"
	sConfig "financial-spider.go/config/stock"
	"financial-spider.go/models"
	"financial-spider.go/models/vo"
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

	log.Printf("查询股票对应公司财报数据（总计 : %d 期）", reportCount)

	for i, reportDate := range reportDates {
		log.Printf("处理报表进度 : %d / %d", i+1, reportCount)

		processingCashFlowSheet(code, reportDate)
		processingBalanceSheet(code, reportDate)
		processingIncomeSheet(code, reportDate)
	}
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
func processingCashFlowSheet(code string, reportDate string) {

}

// 处理资产负债表
func processingBalanceSheet(code string, reportDate string) {

}

// 处理利润表
func processingIncomeSheet(code string, reportDate string) {

}
