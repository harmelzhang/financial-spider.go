package response

// 行业分类
type Category struct {
	// ID
	Id string `json:"catId"`
	// 名称
	Name string `json:"industryName"`
	// 层级
	Level string `json:"level"`
	// 父ID
	ParentId string `json:"parentCid"`
	// 子行业
	Children []Category `json:"children"`
}

// 行业分类数据
type CategoryData struct {
	// 交易日期
	Date string `json:"tradeDate"`
	// 行业分类
	MapList map[string][]Category `json:"csipeindustryMapList"`
}

// 行业分类响应
type CategoryResult struct {
	// 状态码
	Code string `json:"code"`
	// 是否成功
	Success bool `json:"success"`
	// 消息
	Message string `json:"msg"`
	// 行业分类数据
	Data CategoryData `json:"data"`
}
