package models

// IndexSample 指数样本信息
type IndexSample struct {
	TypeCode  string // 类型代码（中证指数，www.csindex.com.cn）
	TypeName  string // 类型名称（沪深300、中证500、上证50....）
	StockCode string // 股票代码
}
