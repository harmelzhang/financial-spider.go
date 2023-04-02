package stock

// QueryStockBaseInfoUrl 查询股票基本信息接口
const QueryStockBaseInfoUrl = "https://emweb.securities.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s"

// QueryStockMainBusinessUrl 查询股票所属公司主营业务接口
const QueryStockMainBusinessUrl = "https://datacenter.eastmoney.com/securities/api/data/get?type=RPT_F10_ORG_BASICINFO&sty=MAIN_BUSINESS&filter=(SECURITY_CODE=%s)"

// QueryCompanyTypeUrl 查询公司类型接口
const QueryCompanyTypeUrl = "https://datacenter.eastmoney.com/securities/api/data/v1/get?reportName=RPT_F10_ORG_BASICINFO&columns=ORG_TYPE,ORG_TYPE_CODE&filter=(SECURITY_CODE=%s)"
