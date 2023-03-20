package models

import (
	"financial-spider.go/utils/db"
)

// Category 行业分类
type Category struct {
	Type         string // 分类类型（证券会、中证）
	Id           string // 行业ID
	Name         string // 名称
	Level        string // 层级
	DisplayOrder int    // 显示顺序
	ParentId     string // 父分类ID
}

// IntoDb 更新数据库
func (category *Category) IntoDb() {
	sql := "INSERT INTO category(type, id, name, level,  display_order, parent_id) VALUES(?, ?, ?, ?, ?, ?)"
	args := []interface{}{category.Type, category.Id, category.Name, category.Level, category.DisplayOrder, category.ParentId}
	if category.ParentId == "" {
		args[len(args)-1] = nil
	}
	db.ExecSQL(sql, args...)
}
