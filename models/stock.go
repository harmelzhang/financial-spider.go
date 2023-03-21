package models

// Stock 股票信息
type Stock struct {
	Code                *Value // 股票代码
	StockName           *Value // 股票名称
	StockNamePinyin     *Value // 股票名称（拼音）
	StockBeforeName     *Value // 股票曾用名称
	CompanyName         *Value // 公司名称
	CompanyProfile      *Value // 公司简介
	Region              *Value // 地域（省份）
	Address             *Value // 办公地址
	Website             *Value // 公司网站
	MainBusiness        *Value // 主营业务
	BusinessScope       *Value // 经营范围
	DateOfIncorporation *Value // 成立日期
	ListingDate         *Value // 上市日期
	LawFirm             *Value // 律师事务所
	AccountingFirm      *Value // 会计师事务所
	MarketPlace         *Value // 交易市场（上海、深圳、北京）
}
