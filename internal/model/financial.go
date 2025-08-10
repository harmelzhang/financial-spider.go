package model

// 财务报表
type Financial struct {
	// 股票代码
	StockCode string
	// 年份
	Year string
	// 财报季期
	ReportDate string
	// 季期类型（Q1~Q4，分别代表：一季报、半年报、三季报、年报；O，代表：其他）
	ReportType string

	// 年度分红金额
	Dividend any

	// 营业活动现金流量
	Ocf any
	// 投资活动现金流量
	Cfi any
	// 筹资活动现金流量
	Cff any
	// 分配股利、利润或偿付利息支付的现金
	AssignDividendPorfit any
	// 购建固定资产、无形资产和其他长期资产支付的现金
	AcquisitionAssets any
	// 存货减少额
	InventoryLiquidating any

	// 净利润
	Np any
	// 营业收入
	Oi any
	// 营业成本
	Coe any
	// 营业总成本（含各种费用，销售费用、管理费用等）
	CoeTotal any
	// 每股盈余|基本每股收益
	Eps any

	// 货币资金
	MonetaryFund any
	// 交易性金融资产
	TradeFinassetNotfvtpl any
	// 交易性金融资产（历史遗留）
	TradeFinasset any
	// 衍生金融资产
	DeriveFinasset any

	// 固定资产
	FixedAsset any
	// 在建工程
	Cip any

	// 流动资产总额
	CaTotal any
	// 非流动资产总额
	NcaTotal any
	// 流动负债总额
	ClTotal any
	// 非流动负债产总额
	NclTotal any
	// 存货
	Inventory any
	// 应收账款
	AccountsRece any
	// 应付账款
	AccountsPayable any
}
