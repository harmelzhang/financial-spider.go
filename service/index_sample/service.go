package index_sample

import (
	isConfig "financial-spider.go/config/index_sample"
	"financial-spider.go/models"
	"financial-spider.go/utils/db"
	"financial-spider.go/utils/http"
	"financial-spider.go/utils/tools"
	"financial-spider.go/utils/xls"
	"fmt"
	"log"
)

// 待爬取的地址
var fetchUrls = map[isConfig.Type]string{}

// 指数样本信息
var indexSample = make(map[isConfig.Type][]string)

func init() {
	for indexType, _ := range isConfig.TypeNameMap {
		fetchUrls[indexType] = fmt.Sprintf(isConfig.FetchIndexUrl, indexType)
	}
}

// FetchIndexSample 爬取指数样本信息
func FetchIndexSample() {
	for isType, url := range fetchUrls {
		log.Printf("爬取%s样本信息", isConfig.TypeNameMap[isType])
		data := xls.ReadXls(http.Get(url), 0, 0)
		stockCodes := tools.FetchColData(data, 4)
		indexSample[isType] = append(indexSample[isType], stockCodes...)
	}

	// 删除旧数据
	db.ExecSQL("DELETE FROM index_sample")
	// 插入新数据
	for indexType, stockCodes := range indexSample {
		for _, stockCode := range stockCodes {
			is := models.IndexSample{
				TypeCode:  string(indexType),
				TypeName:  isConfig.TypeNameMap[indexType],
				StockCode: stockCode,
			}
			is.IntoDb()
		}
	}
}

// GetIndexStocks 根据指数类型查询样本股票代码
func GetIndexStocks(typ isConfig.Type) []string {
	return indexSample[typ]
}

// GetStockTypes 获取指定股票指数类型
func GetStockTypes(stockCode string) ([]isConfig.Type, []string) {
	types := make([]isConfig.Type, 0, len(isConfig.TypeNameMap))
	typeNames := make([]string, 0, len(isConfig.TypeNameMap))
	for typ, stockCodes := range indexSample {
		if tools.IndexOf(stockCodes, stockCode) != -1 {
			types = append(types, typ)
			typeNames = append(typeNames, isConfig.TypeNameMap[typ])
		}
	}
	return types, typeNames
}
