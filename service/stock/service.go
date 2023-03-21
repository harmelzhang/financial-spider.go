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
)

// FetchStockBaseInfo 爬取股票基本信息
func FetchStockBaseInfo(code string) {
	log.Println("爬取股票基本信息")

	marketName, marketShortName := QueryStockMarketPlace(code)

	stockBaseInfoRes := vo.StockBaseInfoResult{}

	url := fmt.Sprintf(sConfig.FetchStockBaseInfoUrl, marketShortName, code)
	err := json.Unmarshal(http.Get(url), &stockBaseInfoRes)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	baseInfo := stockBaseInfoRes.BaseInfo[0]
	listingInfo := stockBaseInfoRes.ListingInfo[0]

	stock := models.Stock{
		Code:                models.NewValue(code),
		StockName:           models.NewValue(baseInfo.StockName),
		StockNamePinyin:     models.NewValue(tools.GetPinyinFirstWord(baseInfo.StockName)),
		StockBeforeName:     models.NewValue(baseInfo.StockBeforeName),
		CompanyName:         models.NewValue(baseInfo.CompanyName),
		CompanyProfile:      models.NewValue(baseInfo.CompanyProfile),
		Region:              models.NewValue(baseInfo.Region),
		Address:             models.NewValue(baseInfo.Address),
		Website:             models.NewValue(baseInfo.Website),
		BusinessScope:       models.NewValue(baseInfo.BusinessScope),
		DateOfIncorporation: models.NewValue(listingInfo.DateOfIncorporation),
		ListingDate:         models.NewValue(listingInfo.ListingDate),
		LawFirm:             models.NewValue(baseInfo.LawFirm),
		AccountingFirm:      models.NewValue(baseInfo.AccountingFirm),
		MarketPlace:         models.NewValue(marketName),
	}

	stockMainBusinessResult := vo.StockMainBusinessResult{}
	url = fmt.Sprintf(sConfig.FetchStockMainBusinessUrl, code)
	err = json.Unmarshal(http.Get(url), &stockMainBusinessResult)
	if err != nil {
		log.Printf("解析JSON出错，跳过主营业务 : %s", err)
	}

	if stockMainBusinessResult.Code == 0 && stockMainBusinessResult.Success {
		mainBusiness := stockMainBusinessResult.Result.Data[0].Info
		stock.MainBusiness = models.NewValue(mainBusiness)
	} else {
		log.Printf("获取主营业务数据失败，跳过爬取 > Code:%d, Msg: %s", stockMainBusinessResult.Code, stockMainBusinessResult.Message)
	}
}

// FetchStockFinancialData 爬取股票对应公司财报数据
func FetchStockFinancialData(code string) {
	log.Println("爬取股票对应公司财报数据")
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
