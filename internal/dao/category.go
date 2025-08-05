package dao

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type categoryDao struct{}

// 行业分类数据访问层
var CategoryDao = new(categoryDao)

func (dao *categoryDao) Insert(ctx context.Context, entity *model.Category) (err error) {
	_, err = DB(ctx, model.CategoryTableInfo.Table()).Insert(entity)
	return
}

// 删除所有数据
func (dao *categoryDao) DeleteAll(ctx context.Context) (err error) {
	_, err = DB(ctx, model.CategoryTableInfo.Table()).Delete()
	return
}

// 删除指定类型的数据
func (dao *categoryDao) DeleteByType(ctx context.Context, typeName string) (err error) {
	_, err = DB(ctx, model.CategoryTableInfo.Table()).Where(model.CategoryTableInfo.Columns().Type, typeName).Delete()
	return
}
