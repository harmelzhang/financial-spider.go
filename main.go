package main

import (
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"harmel.cn/financial/internal/cmd"
)

func main() {
	cmd, err := gcmd.NewFromObject(cmd.CommandMain{})
	if err != nil {
		panic(err)
	}
	cmd.Run(gctx.New())
}
