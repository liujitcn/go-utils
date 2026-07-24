package stringcase

import (
	"regexp"
	"strings"
	"unicode"
)

// isLower 判断字符是否为小写字母。
func isLower(ch rune) bool {
	return ch >= 'a' && ch <= 'z'
}

// toLower 将字符转换为小写字母。
func toLower(ch rune) rune {
	if ch >= 'A' && ch <= 'Z' {
		return ch + 32
	}
	return ch
}

// isUpper 判断字符是否为大写字母。
func isUpper(ch rune) bool {
	return ch >= 'A' && ch <= 'Z'
}

// toUpper 将字符转换为大写字母。
func toUpper(ch rune) rune {
	if ch >= 'a' && ch <= 'z' {
		return ch - 32
	}
	return ch
}

// isSpace 判断字符是否为空白字符。
func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// isDigit 判断字符是否为数字。
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isDelimiter 判断字符是否为分隔符。
func isDelimiter(ch rune) bool {
	return ch == '-' || ch == '_' || isSpace(ch)
}

type iterFunc func(prev, curr, next rune)

// stringIter 逐个遍历字符串中的相邻字符。
func stringIter(s string, callback iterFunc) {
	var prev rune
	var curr rune
	for _, next := range s {
		if curr == 0 {
			prev = curr
			curr = next
			continue
		}

		callback(prev, curr, next)

		prev = curr
		curr = next
	}

	if len(s) > 0 {
		callback(prev, curr, 0)
	}
}

// isAlpha 判断字符是否为英文字母。
func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

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

// isUpperCaseWord 判断单词是否全部为大写字母。
func isUpperCaseWord(word string) bool {
	for _, r := range word {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
