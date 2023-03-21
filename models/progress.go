package models

import (
	"encoding/json"
	"financial-spider.go/utils/tools"
	"log"
)

// Progress 进度信息
type Progress struct {
	Done  bool     // 是否完成
	Codes []string // 股票代码
	Time  int64    // 插入时间
}

// NewProgress 构建新的配置对象
func NewProgress() Progress {
	progress := Progress{Codes: make([]string, 0)}
	return progress
}

// Load 加载配置文件
func (progress *Progress) Load(path string) {
	content := tools.ReadFile(path)
	err := json.Unmarshal(content, progress)
	if err != nil {
		log.Fatalf("解析配置文件出错 : %s", err)
	}
}

// Write 写配置文件
func (progress *Progress) Write(path string) {
	content, err := json.Marshal(progress)
	if err != nil {
		log.Fatalf("序列化配置文件出错 : %s", err)
	}
	tools.WriteFile(path, content)
}
