package public

var (
	// 任务间隔天数
	SpiderTaskIntervalDays int64 = 0
	// 线程池执行器初始化个数
	SpiderExecutorPoolSize = 3
	// 请求超时时长（秒）
	SpiderTimtout = 3

	// 上交所股票前缀
	ShanghaiMarketPrefixs []string
	// 深交所股票前缀
	ShenzhenMarketPrefixs []string
	// 北交所股票前缀
	BeijingMarketPrefixs []string
)

var (
	// 行业分类类型（CSRC：证监会、CICS：中证）
	CategoryType = map[string]string{"CSRC": "1", "CICS": "2"}
)
