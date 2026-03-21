package mapper

import (
	strutil "github.com/liujitcn/go-utils/string"

	"github.com/jinzhu/copier"
)

// Int64ArrayTypeConverter int64 数组与 JSON 字符串转换器
type Int64ArrayTypeConverter struct {
}

// NewInt64ArrayTypeConverter 创建 int64 数组与 JSON 字符串转换器
func NewInt64ArrayTypeConverter() *Int64ArrayTypeConverter {
	return &Int64ArrayTypeConverter{}
}

// ToEntity 将 JSON 字符串转换为 int64 数组
func (m *Int64ArrayTypeConverter) ToEntity(dto *string) *[]int64 {
	if dto == nil {
		return nil
	}

	value := strutil.ConvertJsonStringToInt64Array(*dto)
	return &value
}

// ToDTO 将 int64 数组转换为 JSON 字符串
func (m *Int64ArrayTypeConverter) ToDTO(entity *[]int64) *string {
	if entity == nil {
		return nil
	}

	value := strutil.ConvertInt64ArrayToString(*entity)
	return &value
}

// NewConverterPair 创建 int64 数组与 JSON 字符串的双向转换器
func (m *Int64ArrayTypeConverter) NewConverterPair() []copier.TypeConverter {
	fromFn := func(src []int64) string {
		value := m.ToDTO(&src)
		if value == nil {
			return "[]"
		}
		return *value
	}
	toFn := func(src string) []int64 {
		value := m.ToEntity(&src)
		if value == nil {
			return []int64{}
		}
		return *value
	}

	return NewGenericTypeConverterPair([]int64{}, "", fromFn, toFn)
}

// StringArrayTypeConverter string 数组与 JSON 字符串转换器
type StringArrayTypeConverter struct {
}

// NewStringArrayTypeConverter 创建 string 数组与 JSON 字符串转换器
func NewStringArrayTypeConverter() *StringArrayTypeConverter {
	return &StringArrayTypeConverter{}
}

// ToEntity 将 JSON 字符串转换为 string 数组
func (m *StringArrayTypeConverter) ToEntity(dto *string) *[]string {
	if dto == nil {
		return nil
	}

	value := strutil.ConvertJsonStringToStringArray(*dto)
	return &value
}

// ToDTO 将 string 数组转换为 JSON 字符串
func (m *StringArrayTypeConverter) ToDTO(entity *[]string) *string {
	if entity == nil {
		return nil
	}

	value := strutil.ConvertStringArrayToString(*entity)
	return &value
}

// NewConverterPair 创建 string 数组与 JSON 字符串的双向转换器
func (m *StringArrayTypeConverter) NewConverterPair() []copier.TypeConverter {
	fromFn := func(src []string) string {
		value := m.ToDTO(&src)
		if value == nil {
			return "[]"
		}
		return *value
	}
	toFn := func(src string) []string {
		value := m.ToEntity(&src)
		if value == nil {
			return []string{}
		}
		return *value
	}

	return NewGenericTypeConverterPair([]string{}, "", fromFn, toFn)
}
