package stringcase

// KebabCase 将字符串转换为短横线命名法。
func KebabCase(s string) string {
	return delimiterCase(s, '-', false)
}

// UpperKebabCase 将字符串转换为大写短横线命名法。
func UpperKebabCase(s string) string {
	return delimiterCase(s, '-', true)
}
