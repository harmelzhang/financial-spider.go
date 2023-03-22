package stock

// FetchStockBaseInfoUrl 查询股票基本信息接口
const FetchStockBaseInfoUrl = "https://emweb.securities.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s"

// FetchStockMainBusinessUrl 查询股票所属公司主营业务接口
const FetchStockMainBusinessUrl = "https://datacenter.eastmoney.com/securities/api/data/get?type=RPT_F10_ORG_BASICINFO&sty=MAIN_BUSINESS&filter=(SECURITY_CODE=%s)"

// ----- 三大报表报告期 -----

// FetchBalanceSheetReportDateUrl 查询资产负债表报告期接口
const FetchBalanceSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/zcfzbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"

// FetchIncomeSheetReportDateUrl 查询利润表报告期接口
const FetchIncomeSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/lrbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"

// FetchCashFlowSheetReportDateUrl 查询现金流量表报告期接口
const FetchCashFlowSheetReportDateUrl = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/xjllbDateAjaxNew?companyType=4&reportDateType=0&code=%s%s"
