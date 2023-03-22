package category

// QueryCategoryUrl 查询行业分类接口（证券会：1、中证：2）
const QueryCategoryUrl = "https://www.csindex.com.cn/csindex-home/dataServer/queryCsiPeIndustryBytradeDate?classType=%s"

// QueryStockCodeUrl 查询行业下的股票信息接口
const QueryStockCodeUrl = "https://www.csindex.com.cn/csindex-home/dataServer/queryCsiPeSecurity?classType=%s&CicsCode=&level=&isAll=true"

// ShanghaiMarketPrefixs 上交所股票前缀
var ShanghaiMarketPrefixs = []string{"60", "68"}

// ShenzhenMarketPrefixs 深交所股票前缀
var ShenzhenMarketPrefixs = []string{"00", "30"}

// BeijingMarketPrefixs 北交所股票前缀
var BeijingMarketPrefixs = []string{"82", "83", "87", "88", "43"}

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
