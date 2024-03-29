package vo

type StockBaseInfo struct {
	StockName       interface{} `json:"SECURITY_NAME_ABBR"` // 股票名称
	StockBeforeName interface{} `json:"FORMERNAME"`         // 股票曾用名称
	CompanyName     interface{} `json:"ORG_NAME"`           // 公司名称
	CompanyProfile  interface{} `json:"ORG_PROFILE"`        // 公司简介
	Region          interface{} `json:"PROVINCE"`           // 地域
	Address         interface{} `json:"ADDRESS"`            // 办公地址
	Website         interface{} `json:"ORG_WEB"`            // 公司网站
	BusinessScope   interface{} `json:"BUSINESS_SCOPE"`     // 经营范围
	AccountingFirm  interface{} `json:"ACCOUNTFIRM_NAME"`   // 会计师事务所
	LawFirm         interface{} `json:"LAW_FIRM"`           // 律师事务所
}

type StockListingInfo struct {
	DateOfIncorporation interface{} `json:"FOUND_DATE"`   // 成立日期
	ListingDate         interface{} `json:"LISTING_DATE"` // 上市日期
}

// StockBaseInfoResult 股票基本信息
type StockBaseInfoResult struct {
	BaseInfo    []StockBaseInfo    `json:"jbzl"` // 基本信息
	ListingInfo []StockListingInfo `json:"fxxg"` // 发行相关
}

// -----

type MBInfo struct {
	Info interface{} `json:"MAIN_BUSINESS"`
}

// StockMB 主营业务员
type StockMB struct {
	Count int      `json:"count"`
	Pages int      `json:"pages"`
	Data  []MBInfo `json:"data"` // 只有一个
}

// StockMainBusinessResult 股票主营业务
type StockMainBusinessResult struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
	Result  StockMB     `json:"result"` // 主营业务
}

// -----

type StockReportDate struct {
	Date string `json:"REPORT_DATE"`
	Type string `json:"REPORT_TYPE"`
}

// StockReportDateResult 报告期
type StockReportDateResult struct {
	Count int               `json:"count"`
	Pages int               `json:"pages"`
	Data  []StockReportDate `json:"data"` // 只有一个
}

// -----

type CompanyType struct {
	Type     string `json:"ORG_TYPE"`
	TypeCode string `json:"ORG_TYPE_CODE"`
}

type CompanyTypeData struct {
	Count int           `json:"count"`
	Pages int           `json:"pages"`
	Data  []CompanyType `json:"data"` // 只有一个
}

// CompanyTypeResult 公司类型
type CompanyTypeResult struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Success bool            `json:"success"`
	Result  CompanyTypeData `json:"result"`
}
