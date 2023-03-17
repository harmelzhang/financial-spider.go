package index_sample

// IndexType 指数类型
type IndexType string

const (
	HS300 IndexType = "000300" // 沪深300
	ZZ500 IndexType = "000905" // 中证500
	SZ50  IndexType = "000016" // 上证50
	KC50  IndexType = "000688" // 科创50
	HLZS  IndexType = "000015" // 红利指数
)

// IndexTypeNameMap 指数类型名称映射
var IndexTypeNameMap = map[IndexType]string{
	HS300: "沪深300",
	ZZ500: "中证500",
	SZ50:  "上证50",
	KC50:  "科创50",
	HLZS:  "红利指数",
}
