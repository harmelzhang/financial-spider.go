package category

// FetchCategoryUrl 查询行业分类接口（证券会：1、中证：2）
const FetchCategoryUrl = "https://www.csindex.com.cn/csindex-home/dataServer/queryCsiPeIndustryBytradeDate?classType=%s"

// FetchStockCodeUrl 查询行业下的股票信息接口
const FetchStockCodeUrl = "https://www.csindex.com.cn/csindex-home/dataServer/queryCsiPeSecurity?classType=%s&CicsCode=&level=&isAll=true"

// Type 行业分类类型
type Type string

const (
	ZQH Type = "1" // 证券会行业
	ZZ  Type = "2" // 中证行业
)

// TypeNameMap 指数类型名称映射
var TypeNameMap = map[Type]string{
	ZQH: "证券会",
	ZZ:  "中证",
}
