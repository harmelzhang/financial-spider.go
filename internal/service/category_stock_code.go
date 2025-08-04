package service

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type categoryStockCodeService struct{}

// 行业分类与股票代码业务逻辑处理对象
var CategoryStockCodeService = new(categoryStockCodeService)

// 插入记录
func (cscService *categoryStockCodeService) Insert(ctx context.Context, entity *model.CategoryStockCode) (err error) {
	_, err = model.DB(ctx, model.CategoryStockCodeTableInfo.Table()).Insert(entity)
	return
}

// 删除所有数据
func (cscService *categoryStockCodeService) DeleteAll(ctx context.Context) (err error) {
	_, err = model.DB(ctx, model.CategoryStockCodeTableInfo.Table()).Delete()
	return
}

// 删除指定类型的数据
func (cscService *categoryStockCodeService) DeleteByType(ctx context.Context, typeName string) (err error) {
	_, err = model.DB(ctx, model.CategoryStockCodeTableInfo.Table()).Where(model.CategoryStockCodeTableInfo.Columns().Type, typeName).Delete()
	return
}
