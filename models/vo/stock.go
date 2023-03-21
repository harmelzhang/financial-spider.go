package vo

type StockBaseInfo struct {
	StockName       string `json:"STR_NAMEA"`        // 股票名称
	StockBeforeName string `json:"FORMERNAME"`       // 股票曾用名称
	CompanyName     string `json:"ORG_NAME"`         // 公司名称
	CompanyProfile  string `json:"ORG_PROFILE"`      // 公司简介
	Region          string `json:"PROVINCE"`         // 地域
	Address         string `json:"ADDRESS"`          // 办公地址
	Website         string `json:"ORG_WEB"`          // 公司网站
	BusinessScope   string `json:"BUSINESS_SCOPE"`   // 经营范围
	AccountingFirm  string `json:"ACCOUNTFIRM_NAME"` // 会计师事务所
	LawFirm         string `json:"LAW_FIRM"`         // 律师事务所
}

type StockListingInfo struct {
	DateOfIncorporation string `json:"FOUND_DATE"`   // 成立日期
	ListingDate         string `json:"LISTING_DATE"` // 上市日期
}

type StockBaseInfoResult struct {
	BaseInfo    []StockBaseInfo    `json:"jbzl"` // 基本信息
	ListingInfo []StockListingInfo `json:"fxxg"` // 发行相关
}

// -----

type MBInfo struct {
	Info string `json:"MAIN_BUSINESS"`
}

// StockMB 主营业务员
type StockMB struct {
	Count int      `json:"count"`
	Pages int      `json:"pages"`
	Data  []MBInfo `json:"data"` // 只有一个
}

// StockMainBusinessResult 股票主营业务
type StockMainBusinessResult struct {
	Code    int     `json:"code"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Result  StockMB `json:"result"` // 主营业务
}
