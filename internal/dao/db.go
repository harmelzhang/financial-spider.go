package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// 数据库访问对象
func DB(ctx context.Context, table string) *gdb.Model {
	return g.DB().Model(table).Safe().Ctx(ctx)
}
