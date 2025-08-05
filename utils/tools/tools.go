package tools

import (
	"github.com/mozillazg/go-pinyin"
)

// 获取拼音首字母
func PinyinFirstWord(words string) string {
	result := ""
	pyArgs := pinyin.NewArgs()
	for _, word := range pinyin.Pinyin(words, pyArgs) {
		for _, w := range word {
			result += w[:1]
		}
	}
	return result
}
