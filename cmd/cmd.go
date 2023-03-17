package cmd

import (
	isService "financial-spider.go/service/index_sample"
	"fmt"
)

// 处理函数
type handler func()

// 命令
type command struct {
	name    string  // 名称
	usage   string  // 用法
	handler handler // 处理函数
}

var (
	fetch = command{
		name:  "fetch",
		usage: "抓取网络数据",
		handler: func() {
			isService.Init()
		},
	}
	export = command{
		name:  "export",
		usage: "导出数据到本地",
		handler: func() {
			// TODO 后续会增加导出数据的功能
			fmt.Println(">>>> 导出数据")
		},
	}
)

// 目前支持的所有命令
var cmds = []command{fetch, export}

func Run(args ...string) {
	if len(args) == 1 {
		fmt.Println("USAGE")
		fmt.Printf("    %s\n", "financial COMMAND")
		fmt.Println("COMMAND")
		for _, cmd := range cmds {
			fmt.Printf("    %s\t%s\n", cmd.name, cmd.usage)
		}
		return
	}

	switch args[1] {
	case fetch.name:
		fetch.handler()
	case export.name:
		export.handler()
	default:
		fmt.Println("错误 : 不支持的参数")
	}
}
