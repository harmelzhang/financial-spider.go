package stock

import (
	cConfig "financial-spider.go/config/category"
	"financial-spider.go/utils/tools"
	"log"
)

// FetchStockBaseInfo 爬取股票基本信息
func FetchStockBaseInfo(code string) {
	log.Println("爬取股票基本信息")

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
