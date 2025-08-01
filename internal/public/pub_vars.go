package public

var (
	// 任务间隔天数
	SpiderTaskIntervalDays int64 = 0
	// 线程池执行器初始化个数
	SpiderExecutorPoolSize = 3
)

var (
	// 上交所股票前缀
	ShanghaiMarketPrefixs = []string{"60", "68"}
	// 深交所股票前缀
	ShenzhenMarketPrefixs = []string{"00", "30"}
	// 北交所股票前缀
	BeijingMarketPrefixs = []string{"82", "83", "87", "88", "43"}
)
