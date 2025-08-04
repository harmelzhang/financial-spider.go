package service

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type categoryService struct{}

// 行业分类业务逻辑处理对象
var CategoryService = new(categoryService)

// 插入记录
func (cService *categoryService) Insert(ctx context.Context, entity *model.Category) (err error) {
	_, err = model.DB(ctx, model.CategoryTableInfo.Table()).Insert(entity)
	return
}

// 删除所有数据
func (cService *categoryService) DeleteAll(ctx context.Context) (err error) {
	_, err = model.DB(ctx, model.CategoryTableInfo.Table()).Delete()
	return
}

// 删除指定类型的数据
func (cService *categoryService) DeleteByType(ctx context.Context, typeName string) (err error) {
	_, err = model.DB(ctx, model.CategoryTableInfo.Table()).Where(model.CategoryTableInfo.Columns().Type, typeName).Delete()
	return
}
