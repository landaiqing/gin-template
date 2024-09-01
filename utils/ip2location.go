package utils

import (
	"regexp"
	"strings"
)

func RemoveZeroAndAdjust(s string) string {
	// 正则表达式匹配 "|0|" 或 "|0" 或 "0|" 并替换为 "|"
	re := regexp.MustCompile(`(\|0|0\||0)`)
	result := re.ReplaceAllString(s, "|")

	// 移除可能出现的连续 "|"
	re = regexp.MustCompile(`\|+`)
	result = re.ReplaceAllString(result, "|")

	// 移除字符串开头和结尾可能出现的 "|"
	result = strings.Trim(result, "|")

	return result
}
