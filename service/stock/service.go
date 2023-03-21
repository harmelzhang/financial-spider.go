package stock

import (
	"encoding/json"
	cConfig "financial-spider.go/config/category"
	sConfig "financial-spider.go/config/stock"
	"financial-spider.go/models/vo"
	"financial-spider.go/utils/http"
	"financial-spider.go/utils/tools"
	"fmt"
	"log"
)

// FetchStockBaseInfo 爬取股票基本信息
func FetchStockBaseInfo(code string) {
	log.Println("爬取股票基本信息")

	_, marketShortName := QueryStockMarketPlace(code)

	stockBaseInfoRes := vo.StockBaseInfoResult{}

	url := fmt.Sprintf(sConfig.FetchStockBaseInfoUrl, marketShortName, code)
	err := json.Unmarshal(http.Get(url), &stockBaseInfoRes)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	baseInfo := stockBaseInfoRes.BaseInfo[0]
	listingInfo := stockBaseInfoRes.ListingInfo[0]

	stockMainBusinessResult := vo.StockMainBusinessResult{}
	url = fmt.Sprintf(sConfig.FetchStockMainBusinessUrl, code)
	err = json.Unmarshal(http.Get(url), &stockMainBusinessResult)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if stockMainBusinessResult.Code == 0 && stockMainBusinessResult.Success {
		mainBusiness := stockMainBusinessResult.Result.Data[0].Info
		fmt.Println(baseInfo, listingInfo, mainBusiness)
	} else {
		log.Fatalf("获取主营业务数据失败 > Code:%d, Msg: %s", stockMainBusinessResult.Code, stockMainBusinessResult.Message)
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
