package cmd

import (
	"financial-spider.go/config"
	"financial-spider.go/models"
	cService "financial-spider.go/service/category"
	isService "financial-spider.go/service/index_sample"
	sService "financial-spider.go/service/stock"
	"financial-spider.go/utils/tools"
	"fmt"
	"log"
	"time"
)

// 处理函数
type handler func([]string)

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
		handler: func(args []string) {
			isService.QueryIndexSample()
			cService.QueryCategory()

			stockCodes, total := cService.FindAllStockCodes()

			progress := models.NewProgress()

			// 初始化配置文件
			if tools.FileIsExist(config.ProgressFileName) {
				progress.Load(config.ProgressFileName)
			} else {
				progress.Write(config.ProgressFileName)
			}

			// 如果到了五月一日，全部重跑（年报全部出了）
			if time.Now().Format("01-02") == "05-01" {
				progress = models.NewProgress()
			}

			// 如果上次成功了，判断时间是否大于配置天数
			if progress.Done {
				if time.Now().Unix()-progress.Time >= config.TaskIntervalDay*24*3600 {
					progress = models.NewProgress()
				} else {
					log.Printf("任务结束 : 离上次任务成功结束时间小于%d天", config.TaskIntervalDay)
					return
				}
			}

			log.Println("查询股票基本信息和对应公司财报数据")
			for i, code := range stockCodes {
				log.Printf("任务进度 : %d / %d", i+1, total)

				if tools.IndexOf(progress.Codes, code) != -1 {
					log.Println("跳过已查询的公司")
					continue
				}

				// 查询数据
				stock := sService.QueryStockBaseInfo(code)
				sService.QueryStockFinancialData(stock, args)

				progress.Codes = append(progress.Codes, code)
				progress.Time = time.Now().Unix()
				progress.Write(config.ProgressFileName)
			}

			progress.Done = true
			progress.Time = time.Now().Unix()
			progress.Write(config.ProgressFileName)

			log.Println("任务结束！")

		},
	}
	export = command{
		name:  "export",
		usage: "导出数据到本地",
		handler: func(args []string) {
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
		fetch.handler(args[2:])
	case export.name:
		export.handler(args[2:])
	default:
		fmt.Println("错误 : 不支持的参数")
	}
}
