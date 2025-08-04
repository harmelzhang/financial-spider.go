package service

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type indexSampleService struct{}

// 指数样本逻辑处理对象
var IndexSampleService = new(indexSampleService)

// 插入记录
func (isService *indexSampleService) Insert(ctx context.Context, entity *model.IndexSample) (err error) {
	_, err = model.DB(ctx, model.IndexSampleTableInfo.Table()).Insert(entity)
	return
}

// 删除所有数据
func (isService *indexSampleService) DeleteAll(ctx context.Context) (err error) {
	_, err = model.DB(ctx, model.IndexSampleTableInfo.Table()).Delete()
	return
}

// 删除指定类型的数据
func (isService *indexSampleService) DeleteByType(ctx context.Context, typeCode string) (err error) {
	_, err = model.DB(ctx, model.IndexSampleTableInfo.Table()).Where(model.IndexSampleTableInfo.Columns().TypeCode, typeCode).Delete()
	return
}
