package response

// 股票基本信息
type StockBaseInfo struct {
	// 名称
	Name string `json:"SECURITY_NAME_ABBR"`
	// 曾用名称
	BeforeName string `json:"FORMERNAME"`
	// 公司名称
	CompanyName string `json:"ORG_NAME"`
	// 公司简介
	CompanyProfile string `json:"ORG_PROFILE"`
	// 地域
	Region string `json:"PROVINCE"`
	// 办公地址
	Address string `json:"ADDRESS"`
	// 公司网站
	Website string `json:"ORG_WEB"`
	// 经营范围
	BusinessScope string `json:"BUSINESS_SCOPE"`
	// 会计师事务所
	AccountingFirm string `json:"ACCOUNTFIRM_NAME"`
	// 律师事务所
	LawFirm string `json:"LAW_FIRM"`
}

// 股票发行相关信息
type StockListingInfo struct {
	// 成立日期
	CreateDate string `json:"FOUND_DATE"`
	// 上市日期
	ListingDate string `json:"LISTING_DATE"`
}

// 股票基本信息响应
type StockBaseInfoResult struct {
	// 基本信息
	BaseInfo []StockBaseInfo `json:"jbzl"`
	// 发行相关
	ListingInfo []StockListingInfo `json:"fxxg"`
}

// -----

// 主营业务
type MainBusiness struct {
	// 主营业务
	Info string `json:"MAIN_BUSINESS"`
}

// 主营业务数据
type MainBusinessData struct {
	// 总条数
	Count int `json:"count"`
	// 页码
	Pages int `json:"pages"`
	// 数据
	Data []MainBusiness `json:"data"` // 只有一个
}

// 主营业务响应
type MainBusinessResult struct {
	// 状态码
	Code int `json:"code"`
	// 消息
	Message string `json:"message"`
	// 是否成功
	Success bool `json:"success"`
	// 主营业务数据
	Result MainBusinessData `json:"result"`
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
