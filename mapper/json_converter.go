package mapper

import (
	"encoding/json"

	"github.com/jinzhu/copier"
)

// JSONTypeConverter JSON 字符串与对象转换器
type JSONTypeConverter[T any] struct {
}

// NewJSONTypeConverter 创建 JSON 字符串与对象转换器
func NewJSONTypeConverter[T any]() *JSONTypeConverter[T] {
	return &JSONTypeConverter[T]{}
}

// ToEntity 将 JSON 字符串转换为对象
func (m *JSONTypeConverter[T]) ToEntity(dto *string) *T {
	if dto == nil {
		return nil
	}

	var entity T
	if err := json.Unmarshal([]byte(*dto), &entity); err != nil {
		return nil
	}
	return &entity
}

// ToDTO 将对象转换为 JSON 字符串
func (m *JSONTypeConverter[T]) ToDTO(entity *T) *string {
	if entity == nil {
		return nil
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return nil
	}
	return new(string(data))
}

// NewConverterPair 创建 JSON 字符串与对象的双向转换器
func (m *JSONTypeConverter[T]) NewConverterPair() []copier.TypeConverter {
	fromFn := func(src T) string {
		value := m.ToDTO(&src)
		if value == nil {
			return ""
		}
		return *value
	}
	toFn := func(src string) T {
		value := m.ToEntity(&src)
		if value == nil {
			var zero T
			return zero
		}
		return *value
	}

	var zero T
	return NewGenericTypeConverterPair(zero, "", fromFn, toFn)
}
