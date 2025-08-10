package model

type Stock struct {
	// 代码
	Code string
	// 名称
	Name string
	// 名称（拼音）
	NamePinYin string
	// 曾用名称
	BeforeName string
	// 公司名称
	CompanyName string
	// 公司类型
	CompanyType string
	// 公司类型代码
	CompanyTypeCode string
	// 公司简介
	CompanyProfile string
	// 地域（省份）
	Region string
	// 办公地址
	Address string
	// 公司网站
	Website string
	// 主营业务
	MainBusiness string
	// 经营范围
	BusinessScope string
	// 成立日期
	CreateDate string
	// 上市日期
	ListingDate string
	// 律师事务所
	LawFirm string
	// 会计师事务所
	AccountingFirm string
	// 交易市场（上海、深圳、北京）
	MarketPlace string
}

// 股票表所有列信息
type stockColumns struct {
	// 代码
	Code string
	// 名称
	Name string
	// 名称（拼音）
	NamePinYin string
	// 曾用名称
	BeforeName string
	// 公司名称
	CompanyName string
	// 公司类型
	CompanyType string
	// 公司类型代码
	CompanyTypeCode string
	// 公司简介
	CompanyProfile string
	// 地域（省份）
	Region string
	// 办公地址
	Address string
	// 公司网站
	Website string
	// 主营业务
	MainBusiness string
	// 经营范围
	BusinessScope string
	// 成立日期
	CreateDate string
	// 上市日期
	ListingDate string
	// 律师事务所
	LawFirm string
	// 会计师事务所
	AccountingFirm string
	// 交易市场（上海、深圳、北京）
	MarketPlace string
}

// 股票表信息
type stockTableInfo struct {
	// 表名
	table string
	// 所有列名
	columns stockColumns
}

var StockTableInfo = stockTableInfo{
	table: "stock",
	columns: stockColumns{
		Code:            "code",
		Name:            "name",
		NamePinYin:      "name_pinyin",
		BeforeName:      "before_name",
		CompanyName:     "company_name",
		CompanyType:     "company_type",
		CompanyTypeCode: "company_type_code",
		CompanyProfile:  "company_profile",
		Region:          "region",
		Address:         "address",
		Website:         "website",
		MainBusiness:    "main_business",
		BusinessScope:   "business_scope",
		CreateDate:      "create_date",
		ListingDate:     "listing_date",
		LawFirm:         "law_firm",
		AccountingFirm:  "accounting_firm",
		MarketPlace:     "market_place",
	},
}

// 数据表名
func (info *stockTableInfo) Table() string {
	return info.table
}

// 字段名（列名）
func (info *stockTableInfo) Columns() stockColumns {
	return info.columns
}
