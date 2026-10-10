package stringcase

import (
	"regexp"
	"strings"
	"unicode"
)

// ReplaceNonAlphanumeric 将非字母数字字符替换为指定内容。
func ReplaceNonAlphanumeric(s string, replacement string) string {
	if replacement == "" {
		replacement = "_"
	}
	// 使用正则表达式匹配非英文字母和数字的字符
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	// 替换为指定字符
	return re.ReplaceAllString(s, replacement)
}

// SplitByNonAlphanumeric 按非字母数字字符拆分字符串。
func SplitByNonAlphanumeric(input string) []string {
	var builder strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ') // 将非英文字符和数字的字符替换为空格
		}
	}
	processedInput := builder.String()
	return strings.Fields(processedInput) // 使用空格分割字符串
}

// SplitAndKeepDelimiters 拆分字符串并保留分隔符。
func SplitAndKeepDelimiters(input string) []string {
	var result []string
	var builder strings.Builder

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else {
			if builder.Len() > 0 {
				result = append(result, builder.String())
				builder.Reset()
			}
			result = append(result, string(r)) // 保留分隔符
		}
	}

	if builder.Len() > 0 {
		result = append(result, builder.String())
	}

	return result
}

// ContainsFn 使用自定义比较函数判断切片是否包含指定值。
func ContainsFn[T any](slice []T, value T, predicate func(got, want T) bool) bool {
	for _, item := range slice {
		if predicate(item, value) {
			return true
		}
	}
	return false
}
