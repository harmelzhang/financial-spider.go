package financial

// QueryCashFlowSheetUrl 查询现金流量表数据接口
const QueryCashFlowSheetUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/xjllbAjaxNew?companyType=4&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// QueryBalanceSheetUrl 查询资产负债表数据接口
const QueryBalanceSheetUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/zcfzbAjaxNew?companyType=4&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// QueryIncomeSheetUrl 查询利润表数据接口
const QueryIncomeSheetUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/lrbAjaxNew?companyType=4&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// QueryDividendUrl 查询分红数据接口
const QueryDividendUrl = "https://datacenter.eastmoney.com/securities/api/data/v1/get?reportName=RPT_F10_DIVIDEND_COMPRE&columns=TOTAL_DIVIDEND,STATISTICS_YEAR&filter=(SECURITY_CODE=%s)"

// QueryPageSize 每次取财报数据时，传入的报告期数，最大值为5
const QueryPageSize = 5
