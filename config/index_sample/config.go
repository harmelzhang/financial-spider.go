package index_sample

// FetchIndexUrl 查询指数样本信息地址
const FetchIndexUrl = "https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com/static/html/csindex/public/uploads/file/autofile/cons/%scons.xls"

// Type 指数类型
type Type string

const (
	HS300 Type = "000300" // 沪深300
	ZZ500 Type = "000905" // 中证500
	SZ50  Type = "000016" // 上证50
	KC50  Type = "000688" // 科创50
	HLZS  Type = "000015" // 红利指数
)

// TypeNameMap 指数类型名称映射
var TypeNameMap = map[Type]string{
	HS300: "沪深300",
	ZZ500: "中证500",
	SZ50:  "上证50",
	KC50:  "科创50",
	HLZS:  "红利指数",
}
