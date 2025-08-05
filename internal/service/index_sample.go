package service

import (
	"context"

	"harmel.cn/financial/internal/dao"
	"harmel.cn/financial/internal/model"
)

type indexSampleService struct{}

// 指数样本逻辑处理对象
var IndexSampleService = new(indexSampleService)

// 插入记录
func (s *indexSampleService) Insert(ctx context.Context, entity *model.IndexSample) (err error) {
	err = dao.IndexSampleDao.Insert(ctx, entity)
	return
}

// 删除所有数据
func (s *indexSampleService) DeleteAll(ctx context.Context) (err error) {
	err = dao.IndexSampleDao.DeleteAll(ctx)
	return
}

// 删除指定类型的数据
func (s *indexSampleService) DeleteByType(ctx context.Context, typeCode string) (err error) {
	err = dao.IndexSampleDao.DeleteByType(ctx, typeCode)
	return
}
