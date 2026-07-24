package set

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

var errNilSafeSetUnmarshal = errors.New("set: nil threadSafeSet 无法反序列化 JSON")

// Elements 返回可用于 range 的元素迭代函数。
func Elements[T comparable](collection Set[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		if collection == nil {
			return
		}

		collection.Each(func(value T) bool {
			return !yield(value)
		})
	}
}

// Sorted 返回集合元素的升序切片。
func Sorted[T cmp.Ordered](collection Set[T]) []T {
	values := setValues(collection)
	slices.Sort(values)
	return values
}

// formatSetString 将元素快照格式化为稳定的集合字符串。
func formatSetString[T comparable](values []T) string {
	if len(values) == 0 {
		return "Set{}"
	}

	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, fmt.Sprint(value))
	}
	slices.Sort(parts)

	return "Set{" + strings.Join(parts, ", ") + "}"
}

// setCardinality 安全读取集合元素数量，nil 集合按空集合处理。
func setCardinality[T comparable](collection Set[T]) int {
	if collection == nil {
		return 0
	}
	return collection.Cardinality()
}

// setValues 安全读取集合元素快照，nil 集合按空集合处理。
func setValues[T comparable](collection Set[T]) []T {
	if collection == nil {
		return []T{}
	}
	return collection.ToSlice()
}

// jsonMarshalSet 将集合编码为 JSON 数组。
func jsonMarshalSet[T comparable](collection Set[T]) ([]byte, error) {
	return json.Marshal(setValues(collection))
}

// jsonUnmarshalValues 从 JSON 数组中解码集合元素。
func jsonUnmarshalValues[T comparable](data []byte) ([]T, error) {
	var values []T
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
