package models

import (
	"financial-spider.go/utils/db"
	"financial-spider.go/utils/tools"
	"strings"
)

// Financial 财务报表
type Financial struct {
	Code       string `db:"WHERE"` // 股票代码
	year       string // 年份
	ReportDate string `db:"WHERE"` // 财报季期
	reportType string // 季期类型（Q1~Q4，分别代表：一季报、半年报、三季报、年报；O，代表：其他）

	Dividend interface{} // 年度分红金额

	Ocf interface{} // 营业活动现金流量
	Cfi interface{} // 投资活动现金流量
	Cff interface{} // 筹资活动现金流量

	Np interface{} // 净利润
}

// NewFinancial 新建财务报表对象
func NewFinancial(code string, reportDate string) *Financial {
	ymd := strings.Split(reportDate, "-")
	financial := Financial{Code: code, year: ymd[0], ReportDate: reportDate}
	switch ymd[1] {
	case "03":
		financial.reportType = "Q1"
	case "06":
		financial.reportType = "Q2"
	case "09":
		financial.reportType = "Q3"
	case "12":
		financial.reportType = "Q4"
	default:
		financial.reportType = "O"
	}
	return &financial
}

// InitData 初始化数据库数据
func (financial *Financial) InitData() {
	sql := "SELECT COUNT(*) AS cnt FROM financial WHERE code = ? AND year = ? AND report_date = ?"
	args := []interface{}{financial.Code, financial.year, financial.ReportDate}
	result := db.ExecSQL(sql, args...)

	if result[0]["cnt"].(int64) == 0 {
		sql = "INSERT INTO financial(code, year, report_date, report_type) VALUES(?, ?, ?, ?)"
		args = append(args, financial.reportType)
		db.ExecSQL(sql, args...)
	}
}

// UpdateData 更新数据库数据
func (financial *Financial) UpdateData() {
	sql, args := tools.MakeUpdateSqlAndArgs(financial)
	if sql == "" {
		return
	}
	db.ExecSQL(sql, args...)
}
