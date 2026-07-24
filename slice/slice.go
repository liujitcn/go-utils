package slice

import (
	"github.com/liujitcn/go-utils/set"
	"golang.org/x/exp/constraints"
)

// Filter 按条件过滤切片元素。
func Filter[T any](slice []T, predicate func(value T, index int, slice []T) bool) (filtered []T) {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			filtered = append(filtered, el)
		}
	}
	return filtered
}

// ForEach 遍历切片中的每个元素。
func ForEach[T any](slice []T, function func(value T, index int, slice []T)) {
	for i, el := range slice {
		function(el, i, slice)
	}
}

// Map 将切片元素映射为新的切片。
func Map[T any, R any](slice []T, mapper func(value T, index int, slice []T) R) (mapped []R) {
	if len(slice) > 0 {
		mapped = make([]R, len(slice))
		for i, el := range slice {
			mapped[i] = mapper(el, i, slice)
		}
	}
	return mapped
}

// Reduce 对切片元素进行归约计算。
func Reduce[T any, R any](slice []T, reducer func(acc R, value T, index int, slice []T) R, initial R) R {
	acc := initial
	for i, el := range slice {
		acc = reducer(acc, el, i, slice)
	}
	return acc
}

// Find 查找第一个满足条件的元素。
func Find[T any](slice []T, predicate func(value T, index int, slice []T) bool) *T {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			return &el
		}
	}
	return nil
}

// FindIndex 查找第一个满足条件元素的索引。
func FindIndex[T any](slice []T, predicate func(value T, index int, slice []T) bool) int {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			return i
		}
	}
	return -1
}

// FindIndexOf 查找指定值首次出现的索引。
func FindIndexOf[T comparable](slice []T, value T) int {
	for i, el := range slice {
		if el == value {
			return i
		}
	}
	return -1
}

// FindLastIndex 从后向前查找第一个满足条件元素的索引。
func FindLastIndex[T any](slice []T, predicate func(value T, index int, slice []T) bool) int {
	for i := len(slice) - 1; i >= 0; i-- {
		el := slice[i]
		if ok := predicate(el, i, slice); ok {
			return i
		}
	}
	return -1
}

// FindLastIndexOf 从后向前查找指定值首次出现的索引。
func FindLastIndexOf[T comparable](slice []T, value T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		el := slice[i]
		if el == value {
			return i
		}
	}
	return -1
}

// FindIndexes 查找所有满足条件元素的索引。
func FindIndexes[T any](slice []T, predicate func(value T, index int, slice []T) bool) []int {
	var indexes []int
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// FindIndexesOf 查找指定值出现的所有索引。
func FindIndexesOf[T comparable](slice []T, value T) []int {
	var indexes []int
	for i, el := range slice {
		if el == value {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// Includes 判断切片是否包含指定值。
func Includes[T comparable](slice []T, value T) bool {
	for _, el := range slice {
		if el == value {
			return true
		}
	}
	return false
}

// Some 判断切片中是否至少有一个元素满足条件。
func Some[T any](slice []T, predicate func(value T, index int, slice []T) bool) bool {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			return true
		}
	}
	return false
}

// Every 判断切片中的元素是否全部满足条件。
func Every[T any](slice []T, predicate func(value T, index int, slice []T) bool) bool {
	for i, el := range slice {
		if ok := predicate(el, i, slice); !ok {
			return false
		}
	}
	return true
}

// Merge 合并多个切片。
func Merge[T any](slices ...[]T) (mergedSlice []T) {
	if len(slices) > 0 {
		mergedSliceCap := 0

		for _, slice := range slices {
			mergedSliceCap += len(slice)
		}

		if mergedSliceCap > 0 {
			mergedSlice = make([]T, 0, mergedSliceCap)

			for _, slice := range slices {
				mergedSlice = append(mergedSlice, slice...)
			}
		}
	}
	return mergedSlice
}

// Sum 计算切片元素之和。
func Sum[T constraints.Complex | constraints.Integer | constraints.Float](slice []T) (result T) {
	for _, el := range slice {
		result += el
	}
	return result
}

// Remove 删除指定下标的元素并返回新切片。
func Remove[T any](slice []T, i int) []T {
	if len(slice) == 0 || i > len(slice)-1 {
		return slice
	}
	copied := Copy(slice)
	if i == 0 {
		return copied[1:]
	}
	if i != len(copied)-1 {
		return append(copied[:i], copied[i+1:]...)
	}
	return copied[:i]
}

// Insert 在指定下标插入元素。
func Insert[T any](slice []T, i int, value T) []T {
	if len(slice) == i {
		return append(slice, value)
	}
	slice = append(slice[:i+1], slice[i:]...)
	slice[i] = value
	return slice
}

// Copy 复制切片。
func Copy[T any](slice []T) []T {
	duplicate := make([]T, len(slice), cap(slice))
	copy(duplicate, slice)
	return duplicate
}

// Intersection 计算多个切片的交集。
func Intersection[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return nil
	}

	intersectionSet := set.FromSlice(slices[0])
	for _, currentSlice := range slices[1:] {
		intersectionSet = intersectionSet.Intersection(set.FromSlice(currentSlice))
		if intersectionSet.IsEmpty() {
			return make([]T, 0)
		}
	}
	return intersectionSet.ToSliceByOrder(slices[0])
}

// Difference 计算多个切片的差异元素。
func Difference[T comparable](slices ...[]T) []T {
	firstIndexes := make(map[T]int)
	repeatedValues := set.New[T]()

	for i, currentSlice := range slices {
		currentValues := set.NewWithSize[T](len(currentSlice))
		for _, value := range currentSlice {
			if !currentValues.Add(value) {
				continue
			}

			firstIndex, exists := firstIndexes[value]
			if exists && firstIndex != i {
				repeatedValues.Add(value)
				continue
			}
			if !exists {
				firstIndexes[value] = i
			}
		}
	}

	differentElements := make([]T, 0)
	for _, currentSlice := range slices {
		for _, value := range currentSlice {
			if !repeatedValues.Contains(value) {
				differentElements = append(differentElements, value)
			}
		}
	}
	return differentElements
}

// Union 计算多个切片的并集并去重。
func Union[T comparable](slices ...[]T) []T {
	unionSize := 0
	for _, currentSlice := range slices {
		unionSize += len(currentSlice)
	}

	unioned := make([]T, 0, unionSize)
	visited := set.NewWithSize[T](unionSize)
	for _, currentSlice := range slices {
		for _, value := range currentSlice {
			if visited.Add(value) {
				unioned = append(unioned, value)
			}
		}
	}
	return unioned
}

// Reverse 反转切片元素顺序。
func Reverse[T any](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	result := make([]T, len(slice))
	for i, value := range slice {
		result[len(slice)-1-i] = value
	}
	return result
}

// Unique 对切片元素去重并保持首次出现顺序。
func Unique[T comparable](slice []T) []T {
	unique := make([]T, 0, len(slice))
	visited := set.NewWithSize[T](len(slice))

	for _, value := range slice {
		if visited.Add(value) {
			unique = append(unique, value)
		}
	}
	return unique
}

// Chunk 按指定大小拆分切片。
func Chunk[T any](input []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T

	for i := 0; i < len(input); i += size {
		end := i + size
		if end > len(input) {
			end = len(input)
		}
		chunks = append(chunks, input[i:end])
	}
	return chunks
}

// Pluck 从对象切片中提取指定字段值。
func Pluck[I any, O any](input []I, getter func(I) *O) []O {
	var output []O

	for _, item := range input {
		field := getter(item)

		if field != nil {
			output = append(output, *field)
		}
	}

	return output
}

// Flatten 将二维切片拍平成一维切片。
func Flatten[I any](input [][]I) (output []I) {
	if len(input) > 0 {
		var outputSize int

		for _, item := range input {
			outputSize += len(item)
		}

		if outputSize > 0 {
			output = make([]I, 0, outputSize)

			for _, item := range input {
				output = append(output, item...)
			}
		}
	}
	return output
}
