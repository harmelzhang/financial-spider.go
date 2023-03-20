package vo

type Category struct {
	Id       string     `json:"catId"`
	Name     string     `json:"industryName"`
	Level    string     `json:"level"`
	ParentId string     `json:"parentCid"`
	Children []Category `json:"children"`
}

type CategoryData struct {
	Date    string                `json:"tradeDate"`
	MapList map[string][]Category `json:"csipeindustryMapList"`
}

type CategoryResult struct {
	Code    string       `json:"code"`
	Success bool         `json:"success"`
	Message string       `json:"msg"`
	Data    CategoryData `json:"data"`
}
