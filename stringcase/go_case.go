package stringcase

import "strings"

type abbreviationSet struct {
	items map[string]struct{}
}

// CommonAbbreviations 定义 Go 标识符中需要保持全大写的常见缩写。
var CommonAbbreviations = abbreviationSet{items: map[string]struct{}{
	"ACL":   {},
	"API":   {},
	"ASCII": {},
	"CDN":   {},
	"COS":   {},
	"CPU":   {},
	"CRM":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"ERP":   {},
	"GPS":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IM":    {},
	"IP":    {},
	"JSON":  {},
	"JWT":   {},
	"LBS":   {},
	"LHS":   {},
	"LLM":   {},
	"MFA":   {},
	"MQ":    {},
	"OMS":   {},
	"OSS":   {},
	"OTP":   {},
	"POS":   {},
	"QPS":   {},
	"QR":    {},
	"RAM":   {},
	"RBAC":  {},
	"RHS":   {},
	"RPC":   {},
	"S3":    {},
	"SKU":   {},
	"SLA":   {},
	"SMTP":  {},
	"SMS":   {},
	"SPU":   {},
	"SSH":   {},
	"SSO":   {},
	"TLS":   {},
	"TTL":   {},
	"UID":   {},
	"UI":    {},
	"UUID":  {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"VM":    {},
	"WMS":   {},
	"XML":   {},
	"XSRF":  {},
	"XSS":   {},
}}

// Contains 判断值是否属于通用 Go 缩写，比较时忽略大小写。
func (s abbreviationSet) Contains(value string) bool {
	_, exists := s.items[strings.ToUpper(value)]
	return exists
}

// ToGoPascalCase 将输入转换为保留常见缩写的 Go 大驼峰标识符。
func ToGoPascalCase(input string) string {
	return goCamelCase(input, true)
}

// ToGoCamelCase 将输入转换为保留常见缩写的 Go 小驼峰标识符。
func ToGoCamelCase(input string) string {
	return goCamelCase(input, false)
}

func goCamelCase(input string, upperFirst bool) string {
	words := Split(input)
	var builder strings.Builder
	for index, word := range words {
		if word == "" {
			continue
		}
		upperWord := strings.ToUpper(word)
		if CommonAbbreviations.Contains(upperWord) {
			if index == 0 && !upperFirst {
				builder.WriteString(strings.ToLower(upperWord))
				continue
			}
			builder.WriteString(upperWord)
			continue
		}
		if index == 0 && !upperFirst {
			builder.WriteString(strings.ToLower(word))
			continue
		}
		builder.WriteString(ToPascalCase(word))
	}
	return builder.String()
}
