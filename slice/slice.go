package slice

import (
	"golang.org/x/exp/constraints"
)

func Filter[T any](slice []T, predicate func(value T, index int, slice []T) bool) (filtered []T) {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			filtered = append(filtered, el)
		}
	}
	return filtered
}

func ForEach[T any](slice []T, function func(value T, index int, slice []T)) {
	for i, el := range slice {
		function(el, i, slice)
	}
}

func Map[T any, R any](slice []T, mapper func(value T, index int, slice []T) R) (mapped []R) {
	if len(slice) > 0 {
		mapped = make([]R, len(slice))
		for i, el := range slice {
			mapped[i] = mapper(el, i, slice)
		}
	}
	return mapped
}

func Reduce[T any, R any](slice []T, reducer func(acc R, value T, index int, slice []T) R, initial R) R {
	acc := initial
	for i, el := range slice {
		acc = reducer(acc, el, i, slice)
	}
	return acc
}

func Find[T any](slice []T, predicate func(value T, index int, slice []T) bool) *T {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			return &el
		}
	}
	return nil
}

func FindIndex[T any](slice []T, predicate func(value T, index int, slice []T) bool) int {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			return i
		}
	}
	return -1
}

func FindIndexOf[T comparable](slice []T, value T) int {
	for i, el := range slice {
		if el == value {
			return i
		}
	}
	return -1
}

func FindLastIndex[T any](slice []T, predicate func(value T, index int, slice []T) bool) int {
	for i := len(slice) - 1; i > 0; i-- {
		el := slice[i]
		if ok := predicate(el, i, slice); ok {
			return i
		}
	}
	return -1
}

func FindLastIndexOf[T comparable](slice []T, value T) int {
	for i := len(slice) - 1; i > 0; i-- {
		el := slice[i]
		if el == value {
			return i
		}
	}
	return -1
}

func FindIndexes[T any](slice []T, predicate func(value T, index int, slice []T) bool) []int {
	var indexes []int
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func FindIndexesOf[T comparable](slice []T, value T) []int {
	var indexes []int
	for i, el := range slice {
		if el == value {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func Includes[T comparable](slice []T, value T) bool {
	for _, el := range slice {
		if el == value {
			return true
		}
	}
	return false
}

func Some[T any](slice []T, predicate func(value T, index int, slice []T) bool) bool {
	for i, el := range slice {
		if ok := predicate(el, i, slice); ok {
			return true
		}
	}
	return false
}

func Every[T any](slice []T, predicate func(value T, index int, slice []T) bool) bool {
	for i, el := range slice {
		if ok := predicate(el, i, slice); !ok {
			return false
		}
	}
	return true
}

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

func Sum[T constraints.Complex | constraints.Integer | constraints.Float](slice []T) (result T) {
	for _, el := range slice {
		result += el
	}
	return result
}

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

func Insert[T any](slice []T, i int, value T) []T {
	if len(slice) == i {
		return append(slice, value)
	}
	slice = append(slice[:i+1], slice[i:]...)
	slice[i] = value
	return slice
}

func Copy[T any](slice []T) []T {
	duplicate := make([]T, len(slice), cap(slice))
	copy(duplicate, slice)
	return duplicate
}

func Intersection[T comparable](slices ...[]T) []T {
	possibleIntersections := map[T]int{}
	for i, slice := range slices {
		for _, el := range slice {
			if i == 0 {
				possibleIntersections[el] = 0
			} else if _, elementExists := possibleIntersections[el]; elementExists {
				possibleIntersections[el] = i
			}
		}
	}

	intersected := make([]T, 0)
	for _, el := range slices[0] {
		if lastVisitorIndex, exists := possibleIntersections[el]; exists && lastVisitorIndex == len(slices)-1 {
			intersected = append(intersected, el)
			delete(possibleIntersections, el)
		}
	}

	return intersected
}

func Difference[T comparable](slices ...[]T) []T {
	possibleDifferences := map[T]int{}
	nonDifferentElements := map[T]int{}

	for i, slice := range slices {
		for _, el := range slice {
			if lastVisitorIndex, elementExists := possibleDifferences[el]; elementExists && lastVisitorIndex != i {
				nonDifferentElements[el] = i
			} else if !elementExists {
				possibleDifferences[el] = i
			}
		}
	}

	differentElements := make([]T, 0)

	for _, slice := range slices {
		for _, el := range slice {
			if _, exists := nonDifferentElements[el]; !exists {
				differentElements = append(differentElements, el)
			}
		}
	}

	return differentElements
}

func Union[T comparable](slices ...[]T) []T {
	return Unique(Merge(slices...))
}

func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))

	itemCount := len(slice)
	middle := itemCount / 2
	result[middle] = slice[middle]

	for i := 0; i < middle; i++ {
		mirrorIdx := itemCount - i - 1
		result[i], result[mirrorIdx] = slice[mirrorIdx], slice[i]
	}
	return result
}

func Unique[T comparable](slice []T) []T {
	unique := make([]T, 0)
	visited := map[T]bool{}

	for _, value := range slice {
		if exists := visited[value]; !exists {
			unique = append(unique, value)
			visited[value] = true
		}
	}
	return unique
}

func Chunk[T any](input []T, size int) [][]T {
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
