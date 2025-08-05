package response

// 股票基本信息
type StockBaseInfo struct{}

// 股票发行相关信息
type StockListingInfo struct{}

// 股票基本信息响应
type StockBaseInfoResult struct {
	// 基本信息
	BaseInfo []StockBaseInfo `json:"jbzl"`
	// 发行相关
	ListingInfo []StockListingInfo `json:"fxxg"`
}

// -----

// 公司类型
type CompanyType struct {
	// 类型
	Type string `json:"ORG_TYPE"`
	// 类型代码
	TypeCode string `json:"ORG_TYPE_CODE"`
}

// 公司类型数据
type CompanyTypeData struct {
	// 总条数
	Count int `json:"count"`
	// 页码
	Pages int `json:"pages"`
	// 数据
	Data []CompanyType `json:"data"` // 只有一个
}

// 公司类型响应
type CompanyTypeResult struct {
	// 状态码
	Code int `json:"code"`
	// 消息
	Message string `json:"message"`
	// 是否成功
	Success bool `json:"success"`
	// 公司类型数据
	Result CompanyTypeData `json:"result"`
}
