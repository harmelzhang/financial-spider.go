package stock

// FetchStockBaseInfoUrl 查询股票基本信息接口
const FetchStockBaseInfoUrl = "https://emweb.securities.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s"

// FetchStockMainBusinessUrl 查询股票所属公司主营业务接口
const FetchStockMainBusinessUrl = "https://datacenter.eastmoney.com/securities/api/data/get?type=RPT_F10_ORG_BASICINFO&sty=MAIN_BUSINESS&filter=(SECURITY_CODE=%s)"
