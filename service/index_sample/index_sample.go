package index_sample

import (
	"financial-spider.go/config"
	index "financial-spider.go/config/index_sample"
	"financial-spider.go/models"
	"financial-spider.go/utils/db"
	"financial-spider.go/utils/http"
	"financial-spider.go/utils/tools"
	"financial-spider.go/utils/xls"
	"fmt"
	"log"
)

// 指数样本信息
var indexSample = make(map[index.IndexType][]string)

// Init 初始化操作，抓取网络资源并做入库操作
func Init() {
	// 指数样本信息
	var fetchUrls = map[index.IndexType]string{}
	for indexType, _ := range index.IndexTypeNameMap {
		fetchUrls[indexType] = fmt.Sprintf(config.FetchIndexUrl, indexType)
	}

	log.Println("初始化主要指数样本信息")
	for indexType, url := range fetchUrls {
		data := xls.ReadXls(http.Get(url), 0, 0)
		stockCodes := tools.FetchColData(data, 4)
		indexSample[indexType] = append(indexSample[indexType], stockCodes...)
	}

	// 入库
	db.ExecSQL("DELETE FROM index_sample")
	// 插入数据
	for indexType, stockCodes := range indexSample {
		for _, stockCode := range stockCodes {
			is := models.IndexSample{
				TypeCode:  string(indexType),
				TypeName:  index.IndexTypeNameMap[indexType],
				StockCode: stockCode,
			}
			is.IntoDb()
		}
	}
}

// GetIndexStocks 根据指数类型查询样本股票代码
func GetIndexStocks(typ index.IndexType) []string {
	return indexSample[typ]
}

// GetStockTypes 获取指定股票指数类型
func GetStockTypes(stockCode string) ([]index.IndexType, []string) {
	types := make([]index.IndexType, 0, len(index.IndexTypeNameMap))
	typeNames := make([]string, 0, len(index.IndexTypeNameMap))
	for typ, stockCodes := range indexSample {
		if tools.IndexOf(stockCodes, stockCode) != -1 {
			types = append(types, typ)
			typeNames = append(typeNames, index.IndexTypeNameMap[typ])
		}
	}
	return types, typeNames
}
