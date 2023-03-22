package stock

// QueryStockBaseInfoUrl 查询股票基本信息接口
const QueryStockBaseInfoUrl = "https://emweb.securities.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s"

// QueryStockMainBusinessUrl 查询股票所属公司主营业务接口
const QueryStockMainBusinessUrl = "https://datacenter.eastmoney.com/securities/api/data/get?type=RPT_F10_ORG_BASICINFO&sty=MAIN_BUSINESS&filter=(SECURITY_CODE=%s)"

// ----- 三大报表报告期 -----

// QueryBalanceSheetReportDateUrl 查询资产负债表报告期接口
const QueryBalanceSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/zcfzbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"

// QueryIncomeSheetReportDateUrl 查询利润表报告期接口
const QueryIncomeSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/lrbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"

// QueryCashFlowSheetReportDateUrl 查询现金流量表报告期接口
const QueryCashFlowSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/xjllbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"
