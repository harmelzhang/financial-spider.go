package financial

// ----- 三大报表报告期 -----

// QueryBalanceSheetReportDateUrl 查询资产负债表报告期接口
const QueryBalanceSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/zcfzbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"

// QueryIncomeSheetReportDateUrl 查询利润表报告期接口
const QueryIncomeSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/lrbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"

// QueryCashFlowSheetReportDateUrl 查询现金流量表报告期接口
const QueryCashFlowSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/xjllbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"

// QueryReportDateUrl 查询所有报告期（未使用）
const QueryReportDateUrl = "https://datacenter.eastmoney.com/securities/api/data/get?type=RPT_F10_FINANCE_GINCOME&sty=REPORT_DATE&filter=(SECURITY_CODE=%s)"

// ----- 财务数据接口 -----

// QueryCashFlowSheetUrl 查询现金流量表数据接口
const QueryCashFlowSheetUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/xjllbAjaxNew?companyType=%s&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// QueryBalanceSheetUrl 查询资产负债表数据接口
const QueryBalanceSheetUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/zcfzbAjaxNew?companyType=%s&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// QueryIncomeSheetUrl 查询利润表数据接口
const QueryIncomeSheetUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/lrbAjaxNew?companyType=%s&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// QueryDividendUrl 查询分红数据接口
const QueryDividendUrl = "https://datacenter.eastmoney.com/securities/api/data/v1/get?reportName=RPT_F10_DIVIDEND_COMPRE&columns=TOTAL_DIVIDEND,STATISTICS_YEAR&filter=(SECURITY_CODE=%s)"

// QueryPageSize 每次取财报数据时，传入的报告期数，最大值为5
const QueryPageSize = 5
