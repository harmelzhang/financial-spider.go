package vo

type CashFlowSheetData struct {
	Ocf interface{} `json:"NETCASH_OPERATE"` // 经营活动产生的现金流量净额
	Cfi interface{} `json:"NETCASH_INVEST"`  // 投资活动产生的现金流量净额
	Cff interface{} `json:"NETCASH_FINANCE"` // 筹资活动产生的现金流量净额
}

type CashFlowSheetResult struct {
	Count int                 `json:"count"`
	Pages int                 `json:"pages"`
	Data  []CashFlowSheetData `json:"data"` // 只有一个
	// 错误信息
	Type   string `json:"$type"`  // 1
	Status int    `json:"status"` // -1
}
