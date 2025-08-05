package dao

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type indexSampleDao struct{}

// 指数样本数据访问层
var IndexSampleDao = new(indexSampleDao)

// 插入记录
func (dao *indexSampleDao) Insert(ctx context.Context, entity *model.IndexSample) (err error) {
	_, err = DB(ctx, model.IndexSampleTableInfo.Table()).Insert(entity)
	return
}

// 删除所有数据
func (dao *indexSampleDao) DeleteAll(ctx context.Context) (err error) {
	_, err = DB(ctx, model.IndexSampleTableInfo.Table()).Delete()
	return
}

// 删除指定类型的数据
func (dao *indexSampleDao) DeleteByType(ctx context.Context, typeCode string) (err error) {
	_, err = DB(ctx, model.IndexSampleTableInfo.Table()).Where(model.IndexSampleTableInfo.Columns().TypeCode, typeCode).Delete()
	return
}
