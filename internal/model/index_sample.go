package model

// 指数样本（来源：中证指数）
type IndexSample struct {
	// 类型代码（中证指数，www.csindex.com.cn）
	TypeCode string
	// 类型名称（沪深300、中证500、上证50....）
	TypeName string
	// 股票代码
	StockCode string
}
