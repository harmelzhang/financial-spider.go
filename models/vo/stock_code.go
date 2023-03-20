package vo

type StockCode struct {
	Code string `json:"securityCode"`
	// 中证
	CicsLeve1Code string `json:"cics1stCode"`
	CicsLeve2Code string `json:"cics2ndCode"`
	CicsLeve3Code string `json:"cics3rdCode"`
	CicsLeve4Code string `json:"cics4thCode"`
	// 证券会
	CsrcLeve1Code string `json:"csrc1stCode"`
	CsrcLeve2Code string `json:"csrc2ndCode"`
}

type StockCodeData struct {
	Date string      `json:"tradeDate"`
	List []StockCode `json:"list"`
}

type StockCodeResult struct {
	Code    string        `json:"code"`
	Success bool          `json:"success"`
	Message string        `json:"msg"`
	Data    StockCodeData `json:"data"`
}
