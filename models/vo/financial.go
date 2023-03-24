package vo

type FinancialData struct {
	Ocf interface{} `json:"NETCASH_OPERATE"` // 经营活动产生的现金流量净额
	Cfi interface{} `json:"NETCASH_INVEST"`  // 投资活动产生的现金流量净额
	Cff interface{} `json:"NETCASH_FINANCE"` // 筹资活动产生的现金流量净额

	Np interface{} `json:"NETPROFIT"` // 净利润
}

// FinancialResult 现金流量表
type FinancialResult struct {
	Count int             `json:"count"`
	Pages int             `json:"pages"`
	Data  []FinancialData `json:"data"` // 长度等于传入日期数，这里是一个
	// 错误信息
	Type   string `json:"$type"`  // 1
	Status int    `json:"status"` // -1
}

// -----

type DividendData struct {
	Year  string      `json:"STATISTICS_YEAR"`
	Money interface{} `json:"TOTAL_DIVIDEND"`
}

type DividendPage struct {
	Count int            `json:"count"`
	Pages int            `json:"pages"`
	Data  []DividendData `json:"data"` // 只有一个
}

// DividendResult 分红信息
type DividendResult struct {
	Code    int          `json:"code"`
	Success bool         `json:"success"`
	Message string       `json:"msg"`
	Result  DividendPage `json:"result"`
}
