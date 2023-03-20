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

// 新的分类信息
var categorys = make([]models.Category, 0)

// 初始化爬取地址信息
func init() {
	for cType, _ := range cConfig.TypeNameMap {
		fetchUrls[cType] = fmt.Sprintf(cConfig.FetchCategoryUrl, cType)
	}
}

// 递归获取所有分类信息
func recursionCategorys(cType string, categoryVOs []vo.Category) {
	if len(categoryVOs) == 0 {
		return
	}
	for order, categoryVO := range categoryVOs {
		category := models.Category{
			Type:         cType,
			Id:           categoryVO.Id,
			Name:         categoryVO.Name,
			Level:        categoryVO.Level,
			DisplayOrder: order + 1,
			ParentId:     categoryVO.ParentId,
		}
		categorys = append(categorys, category)

		if len(categoryVO.Children) != 0 {
			recursionCategorys(cType, categoryVO.Children)
		}
	}
}

// FetchCategory 爬取分类信息
func FetchCategory() {
	for categoryType, url := range fetchUrls {
		log.Printf("爬取%s分类信息", cConfig.TypeNameMap[categoryType])

		categoryRes := vo.CategoryResult{}
		err := json.Unmarshal(http.Get(url), &categoryRes)
		if err != nil {
			log.Fatalf("解析JSON出错 : %s", err)
		}

		if categoryRes.Code == "200" && categoryRes.Success {
			recursionCategorys(cConfig.TypeNameMap[categoryType], categoryRes.Data.MapList["4"])
		} else {
			log.Fatalf("获取数据失败 > Code:%s, Msg: %s", categoryRes.Code, categoryRes.Message)
		}
	}

	fmt.Println(">>>", categorys)

	// 删除旧数据
	db.ExecSQL("DELETE FROM category")
	// 插入新数据
	for _, category := range categorys {
		category.IntoDb()
	}
}
