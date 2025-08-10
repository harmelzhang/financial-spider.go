package cmd

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"harmel.cn/financial/internal/public"
	"harmel.cn/financial/internal/spider"
)

type CommandMainSpiderInput struct {
	g.Meta `name:"spider" brief:"start spider server"`
}

type CommandMainSpiderOutput struct{}

func (c *CommandMain) Spider(ctx context.Context, in CommandMainSpiderInput) (out *CommandMainSpiderOutput, err error) {
	rootDir, err := os.Getwd()
	if err != nil {
		g.Log("spider").Errorf(ctx, "get root dir failed, err is %v", err)
		return
	}

	// 系统变量
	public.SpiderTaskIntervalDays = g.Cfg().MustGet(ctx, "spider.taskIntervalDays").Int64()
	public.SpiderExecutorPoolSize = g.Cfg().MustGet(ctx, "spider.executorPoolSize").Int()
	public.SpiderTimtout = g.Cfg().MustGet(ctx, "spider.timeout").Int()
	public.QueryReportPageSize = g.Cfg().MustGet(ctx, "spider.queryReportPageSize").Int()

	// 启动爬虫管理器
	spiderManager := spider.NewSpiderManager(rootDir)
	err = spiderManager.Start(ctx)

	return
}
