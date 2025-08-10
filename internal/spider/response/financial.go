package response

// 报告日期
type ReportDate struct {
	// 日期
	Date string `json:"REPORT_DATE"`
	// 日期类型
	Type string `json:"REPORT_TYPE"`
}

// 报告期
type ReportDateResult struct {
	// 总条数
	Count int `json:"count"`
	// 页码
	Pages int `json:"pages"`
	// 数据
	Data []ReportDate `json:"data"` // 只有一个
}

// -----

// 分红数据
type DividendData struct {
	// 年份
	Year string `json:"STATISTICS_YEAR"`
	// 金额
	Money any `json:"TOTAL_DIVIDEND"`
}

// 分红数据分页
type DividendPage struct {
	// 总条数
	Count int `json:"count"`
	// 页码
	Pages int `json:"pages"`
	// 数据
	Data []DividendData `json:"data"` // 只有一个
}

// 分红信息响应
type DividendResult struct {
	// 状态码
	Code int `json:"code"`
	// 是否成功
	Success bool `json:"success"`
	// 消息
	Message string `json:"msg"`
	// 数据分页
	Result DividendPage `json:"result"`
}

// -----

// 财务报表数据
type FinancialData struct {
	// 报告期：yyyy-MM-dd HH:mm:ss
	ReportDate string `json:"REPORT_DATE"`

	// 经营活动产生的现金流量净额
	Ocf any `json:"NETCASH_OPERATE"`
	// 投资活动产生的现金流量净额
	Cfi any `json:"NETCASH_INVEST"`
	// 筹资活动产生的现金流量净额
	Cff any `json:"NETCASH_FINANCE"`
	// 分配股利、利润或偿付利息支付的现金
	AssignDividendPorfit any `json:"ASSIGN_DIVIDEND_PORFIT"`
	// 购建固定资产、无形资产和其他长期资产支付的现金
	AcquisitionAssets any `json:"CONSTRUCT_LONG_ASSET"`
	// 存货减少额
	InventoryLiquidating any `json:"INVENTORY_REDUCE"`

	// 净利润
	Np any `json:"NETPROFIT"`
	// 营业收入
	Oi any `json:"OPERATE_INCOME"`
	// 营业成本
	Coe any `json:"OPERATE_COST"`
	// 营业总成本（含各种费用，销售费用、管理费用等）
	CoeTotal any `json:"TOTAL_OPERATE_COST"`
	// 每股盈余|基本每股收益
	Eps any `json:"BASIC_EPS"`

	// 货币资金
	MonetaryFund any `json:"MONETARYFUNDS"`
	// 交易性金融资产
	TradeFinassetNotfvtpl any `json:"TRADE_FINASSET_NOTFVTPL"`
	// 交易性金融资产（历史遗留）
	TradeFinasset any `json:"TRADE_FINASSET"`
	// 衍生金融资产
	DeriveFinasset any `json:"DERIVE_FINASSET"`

	// 固定资产
	FixedAsset any `json:"FIXED_ASSET"`
	// 在建工程
	Cip any `json:"CIP"`

	// 流动资产总额
	CaTotal any `json:"TOTAL_CURRENT_ASSETS"`
	// 非流动资产总额
	NcaTotal any `json:"TOTAL_NONCURRENT_ASSETS"`
	// 流动负债总额
	ClTotal any `json:"TOTAL_CURRENT_LIAB"`
	// 非流动负债产总额
	NclTotal any `json:"TOTAL_NONCURRENT_LIAB"`
	// 存货
	Inventory any `json:"INVENTORY"`
	// 应收账款
	AccountsRece any `json:"ACCOUNTS_RECE"`
	// 应付账款
	AccountsPayable any `json:"ACCOUNTS_PAYABLE"`
}

// 财务报表数据响应
type FinancialResult struct {
	// 总条数
	Count int `json:"count"`
	// 页码
	Pages int `json:"pages"`
	// 财务数据
	Data []FinancialData `json:"data"` // 长度等于传入日期数，这里是一个
	// 错误类型
	Type string `json:"$type"` // 1
	// 状态码
	Status int `json:"status"` // -1
}
