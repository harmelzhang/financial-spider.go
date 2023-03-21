package category

import (
	"encoding/json"
	cConfig "financial-spider.go/config/category"
	"financial-spider.go/models"
	"financial-spider.go/models/vo"
	"financial-spider.go/utils/db"
	"financial-spider.go/utils/http"
	"fmt"
	"log"
)

// 待爬取的地址
var fetchUrls = map[cConfig.Type]string{}

// 初始化爬取地址信息
func init() {
	for cType, _ := range cConfig.TypeNameMap {
		fetchUrls[cType] = fmt.Sprintf(cConfig.FetchCategoryUrl, cType)
	}
}

// 递归获取所有分类信息
func recursionCategorys(typeName string, categoryVOs []vo.Category) {
	if len(categoryVOs) == 0 {
		return
	}
	for order, categoryVO := range categoryVOs {
		category := models.Category{
			Type:         typeName,
			Id:           categoryVO.Id,
			Name:         categoryVO.Name,
			Level:        categoryVO.Level,
			DisplayOrder: order + 1,
			ParentId:     categoryVO.ParentId,
		}
		// 插入新数据
		category.IntoDb()

		if len(categoryVO.Children) != 0 {
			recursionCategorys(typeName, categoryVO.Children)
		}
	}
}

// FetchCategory 爬取分类信息
func FetchCategory() {
	for categoryType, url := range fetchUrls {
		typeName := cConfig.TypeNameMap[categoryType]
		log.Printf("爬取%s分类信息", typeName)

		categoryRes := vo.CategoryResult{}
		err := json.Unmarshal(http.Get(url), &categoryRes)
		if err != nil {
			log.Fatalf("解析JSON出错 : %s", err)
		}

		if categoryRes.Code == "200" && categoryRes.Success {
			// 删除旧数据
			db.ExecSQL("DELETE FROM category WHERE type = ?", typeName)
			recursionCategorys(typeName, categoryRes.Data.MapList["4"])
		} else {
			log.Fatalf("获取数据失败 > Code:%s, Msg: %s", categoryRes.Code, categoryRes.Message)
		}

		log.Printf("爬取%s分类下的股票代码与分类之间的关系", typeName)
		findStockCodesByCategoeyType(categoryType)
	}
}

// FindAllStockCodes 查询所有股票代码
func FindAllStockCodes() ([]string, int) {
	data := db.ExecSQL("SELECT DISTINCT stock_code FROM category_stock_code ORDER BY stock_code")
	result := make([]string, 0)
	for _, item := range data {
		result = append(result, item["stock_code"].(string))
	}
	return result, len(result)
}

// 根据分类类型查询股票代码
func findStockCodesByCategoeyType(cType cConfig.Type) {
	url := fmt.Sprintf(cConfig.FetchStockCodeUrl, cType)

	stockCodesRes := vo.StockCodeResult{}
	err := json.Unmarshal(http.Get(url), &stockCodesRes)
	if err != nil {
		log.Fatalf("解析JSON出错 : %s", err)
	}

	if stockCodesRes.Code == "200" && stockCodesRes.Success {
		typeName := cConfig.TypeNameMap[cType]

		// 删除旧数据
		db.ExecSQL("DELETE FROM category_stock_code WHERE type = ?", typeName)

		for _, stockCodeVO := range stockCodesRes.Data.List {
			csc := models.CategoryStockCode{
				Type:      typeName,
				StockCode: stockCodeVO.Code,
			}
			if stockCodeVO.CicsLeve1Code != "" { // 中证（四级）
				csc.CategoryId = stockCodeVO.CicsLeve4Code
			} else { // 证券会（两级）
				// TODO 证券会暂时没对新三板股票进行分类，后续待优化
				if stockCodeVO.CsrcLeve2Code == "" {
					continue
				}
				csc.CategoryId = stockCodeVO.CsrcLeve2Code
			}
			// 插入新数据
			csc.IntoDb()
		}
	} else {
		log.Fatalf("获取数据失败 > Code:%s, Msg: %s", stockCodesRes.Code, stockCodesRes.Message)
	}
}
