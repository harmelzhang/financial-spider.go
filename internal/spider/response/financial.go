package response

// 报告日期
type ReportDate struct {
	// 日期
	Date string `json:"REPORT_DATE"`
	// 日期类型
	Type string `json:"REPORT_TYPE"`
}

// StockReportDateResult 报告期
type StockReportDateResult struct {
	// 总条数
	Count int `json:"count"`
	// 页码
	Pages int `json:"pages"`
	// 数据
	Data []ReportDate `json:"data"` // 只有一个
}
