package models

import "financial-spider.go/utils/db"

// Stock 股票信息
type Stock struct {
	Code                interface{} // 股票代码
	StockName           interface{} // 股票名称
	StockNamePinyin     interface{} // 股票名称（拼音）
	StockBeforeName     interface{} // 股票曾用名称
	CompanyName         interface{} // 公司名称
	CompanyProfile      interface{} // 公司简介
	Region              interface{} // 地域（省份）
	Address             interface{} // 办公地址
	Website             interface{} // 公司网站
	MainBusiness        interface{} // 主营业务
	BusinessScope       interface{} // 经营范围
	DateOfIncorporation interface{} // 成立日期
	ListingDate         interface{} // 上市日期
	LawFirm             interface{} // 律师事务所
	AccountingFirm      interface{} // 会计师事务所
	MarketPlace         interface{} // 交易市场（上海、深圳、北京）
}

// ReplaceData 插入或修改数据
func (stock *Stock) ReplaceData() {
	sql := `
		REPLACE INTO stock(
			code, stock_name, stock_name_pinyin, stock_before_name,
			company_name, company_profile, region, address, website, main_business, business_scope,
			date_of_incorporation, listing_date, law_firm, accounting_firm, market_place
		)
		VALUES(
			?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?
		)
	`
	args := []interface{}{
		stock.Code, stock.StockName, stock.StockNamePinyin, stock.StockBeforeName,
		stock.CompanyName, stock.CompanyProfile, stock.Region, stock.Address, stock.Website, stock.MainBusiness, stock.BusinessScope,
		stock.DateOfIncorporation, stock.ListingDate, stock.LawFirm, stock.AccountingFirm, stock.MarketPlace,
	}
	db.ExecSQL(sql, args...)
}
