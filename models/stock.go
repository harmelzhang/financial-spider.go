package models

import (
	"financial-spider.go/utils/db"
	"financial-spider.go/utils/tools"
)

// Stock 股票信息
type Stock struct {
	Code                string      `db:"WHERE"` // 股票代码
	StockName           interface{} // 股票名称
	StockNamePinyin     interface{} // 股票名称（拼音）
	StockBeforeName     interface{} // 股票曾用名称
	CompanyName         interface{} // 公司名称
	CompanyType         string      // 公司类型
	CompanyTypeCode     string      // 公司类型代码
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

// InitData 初始化数据库数据
func (stock *Stock) InitData() {
	sql := "SELECT COUNT(*) AS cnt FROM stock WHERE code = ?"
	args := []interface{}{stock.Code}
	result := db.ExecSQL(sql, args...)

	if result[0]["cnt"].(int64) == 0 {
		sql = "INSERT INTO stock(code) VALUES(?)"
		db.ExecSQL(sql, args...)
	}
}

// ReplaceData 插入或修改数据
func (stock *Stock) ReplaceData() {
	stock.InitData()
	sql, args := tools.MakeUpdateSqlAndArgs(stock)
	if sql == "" {
		return
	}
	db.ExecSQL(sql, args...)
}
