package public

// 指数样本接口
const UrlIndexSample = "https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com/static/html/csindex/public/uploads/file/autofile/cons/%scons.xls"

// 行业分类接口（证券会：1、中证：2）
const UrlCategory = "https://www.csindex.com.cn/csindex-home/dataServer/queryCsiPeIndustryBytradeDate?classType=%s"

// 行业下的股票信息接口
const UrlCategoryStock = "https://www.csindex.com.cn/csindex-home/dataServer/queryCsiPeSecurity?classType=%s&CicsCode=&level=&isAll=true"

// 股票基本信息接口
const UrlStockBaseInfo = "https://emweb.securities.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s"

// 股票所属公司主营业务接口
const UrlStockMainBusiness = "https://datacenter.eastmoney.com/securities/api/data/get?type=RPT_F10_ORG_BASICINFO&sty=MAIN_BUSINESS&filter=(SECURITY_CODE=%s)"

// 股票所属公司类型接口
const UrlStockCompanyType = "https://datacenter.eastmoney.com/securities/api/data/v1/get?reportName=RPT_F10_ORG_BASICINFO&columns=ORG_TYPE,ORG_TYPE_CODE&filter=(SECURITY_CODE=%s)"

// ----- 财务报表 -----

// 现金流量表数据接口
const UrlCashFlowSheet = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/xjllbAjaxNew?companyType=%s&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// 资产负债表数据接口
const UrlBalanceSheet = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/zcfzbAjaxNew?companyType=%s&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// 利润表数据接口
const UrlIncomeSheet = "https://emweb.securities.eastmoney.com/NewFinanceAnalysis/lrbAjaxNew?companyType=%s&reportDateType=0&reportType=1&dates=%s&code=%s%s"

// 分红数据接口
const UrlDividend = "https://datacenter.eastmoney.com/securities/api/data/v1/get?reportName=RPT_F10_DIVIDEND_COMPRE&columns=TOTAL_DIVIDEND,STATISTICS_YEAR&filter=(SECURITY_CODE=%s)"
