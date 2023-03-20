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

// IntoDb 插入数据库
func (category *Category) IntoDb() {
	sql := "INSERT INTO category(type, id, name, level,  display_order, parent_id) VALUES(?, ?, ?, ?, ?, ?)"
	args := []interface{}{category.Type, category.Id, category.Name, category.Level, category.DisplayOrder, category.ParentId}
	if category.ParentId == "" {
		args[len(args)-1] = nil
	}
	db.ExecSQL(sql, args...)
}

// CategoryStockCode 分类与股票代码之间的关系
type CategoryStockCode struct {
	Type       string // 分类类型（证券会、中证）
	CategoryId string // 分类ID
	StockCode  string // 股票代码
}

// IntoDb 插入数据库
func (csc *CategoryStockCode) IntoDb() {
	sql := "INSERT INTO category_stock_code(type, category_id, stock_code) VALUES(?, ?, ?)"
	args := []interface{}{csc.Type, csc.CategoryId, csc.StockCode}
	db.ExecSQL(sql, args...)
}
