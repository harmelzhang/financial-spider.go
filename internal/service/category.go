package service

import (
	"context"

	"harmel.cn/financial/internal/dao"
	"harmel.cn/financial/internal/model"
)

type categoryService struct{}

// 行业分类业务逻辑处理对象
var CategoryService = new(categoryService)

// 插入记录
func (s *categoryService) Insert(ctx context.Context, entity *model.Category) (err error) {
	err = dao.CategoryDao.Insert(ctx, entity)
	return
}

// 删除所有数据
func (s *categoryService) DeleteAll(ctx context.Context) (err error) {
	err = dao.CategoryDao.DeleteAll(ctx)
	return
}

// 删除指定类型的数据
func (s *categoryService) DeleteByType(ctx context.Context, typeName string) (err error) {
	err = dao.CategoryDao.DeleteByType(ctx, typeName)
	return
}
