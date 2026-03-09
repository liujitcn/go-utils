package slice

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyInt int

type Pluckable struct {
	Code  string
	Value string
}

var numerals = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var numeralsWithUserDefinedType = []MyInt{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var days = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var lastNames = []string{"Jacobs", "Vin", "Jacobs", "Smith"}

func TestFilter(t *testing.T) {
	expectedResult := []int{0, 2, 4, 6, 8}
	actualResult := Filter(numerals, func(value int, _ int, _ []int) bool {
		return value%2 == 0
	})
	assert.Equal(t, expectedResult, actualResult)
}

func TestForEach(t *testing.T) {
	result := 0
	ForEach(numerals, func(value int, _ int, _ []int) {
		result += value
	})
	assert.Equal(t, 45, result)
}

func TestMap(t *testing.T) {
	expectedResult := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	actualResult := Map(numerals, func(value int, _ int, _ []int) string {
		return strconv.Itoa(value)
	})
	assert.Equal(t, expectedResult, actualResult)

	assert.Nil(t, Map([]int{}, func(_ int, _ int, _ []int) string {
		return ""
	}))

	assert.Nil(t, Map([]int(nil), func(_ int, _ int, _ []int) string {
		return ""
	}))
}

func TestReduce(t *testing.T) {
	expectedResult := map[string]string{"result": "0123456789"}
	actualResult := Reduce(
		numerals,
		func(acc map[string]string, cur int, _ int, _ []int) map[string]string {
			acc["result"] += strconv.Itoa(cur)
			return acc
		},
		map[string]string{"result": ""},
	)
	assert.Equal(t, expectedResult, actualResult)
}

func TestFind(t *testing.T) {
	expectedResult := "Wednesday"
	actualResult := Find(days, func(value string, index int, slice []string) bool {
		return strings.Contains(value, "Wed")
	})
	assert.Equal(t, expectedResult, *actualResult)
	assert.Nil(t, Find(days, func(value string, index int, slice []string) bool {
		return strings.Contains(value, "Rishon")
	}))
}

func TestFindIndex(t *testing.T) {
	expectedResult := 3
	actualResult := FindIndex(days, func(value string, index int, slice []string) bool {
		return strings.Contains(value, "Wed")
	})
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, -1, FindIndex(days, func(value string, index int, slice []string) bool {
		return strings.Contains(value, "Rishon")
	}))
}

func TestFindIndexOf(t *testing.T) {
	expectedResult := 3
	actualResult := FindIndexOf(days, "Wednesday")
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, -1, FindIndexOf(days, "Rishon"))
}

func TestFindLastIndex(t *testing.T) {
	expectedResult := 2
	actualResult := FindLastIndex(lastNames, func(value string, index int, slice []string) bool {
		return value == "Jacobs"
	})
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, -1, FindLastIndex(lastNames, func(value string, index int, slice []string) bool {
		return value == "Hamudi"
	}))
}

func TestFindLastIndexOf(t *testing.T) {
	expectedResult := 2
	actualResult := FindLastIndexOf(lastNames, "Jacobs")
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, -1, FindLastIndexOf(lastNames, "Hamudi"))
}

func TestFindIndexes(t *testing.T) {
	expectedResult := []int{0, 2}
	actualResult := FindIndexes(lastNames, func(value string, index int, slice []string) bool {
		return value == "Jacobs"
	})
	assert.Equal(t, expectedResult, actualResult)
	assert.Nil(t, FindIndexes(lastNames, func(value string, index int, slice []string) bool {
		return value == "Hamudi"
	}))
}

func TestFindIndexesOf(t *testing.T) {
	expectedResult := []int{0, 2}
	actualResult := FindIndexesOf(lastNames, "Jacobs")
	assert.Equal(t, expectedResult, actualResult)
	assert.Nil(t, FindIndexesOf(lastNames, "Hamudi"))
}

func TestIncludes(t *testing.T) {
	assert.True(t, Includes(numerals, 1))
	assert.False(t, Includes(numerals, 11))
}

func TestAny(t *testing.T) {
	assert.True(t, Some(numerals, func(value int, _ int, _ []int) bool {
		return value%5 == 0
	}))
	assert.False(t, Some(numerals, func(value int, _ int, _ []int) bool {
		return value == 11
	}))
}

func TestAll(t *testing.T) {
	assert.True(t, Every([]int{1, 1, 1}, func(value int, _ int, _ []int) bool {
		return value == 1
	}))
	assert.False(t, Every([]int{1, 1, 1, 2}, func(value int, _ int, _ []int) bool {
		return value == 1
	}))
}

func TestMerge(t *testing.T) {
	result := Merge(numerals[:5], numerals[5:])
	assert.Equal(t, numerals, result)

	assert.Nil(t, Merge([]int(nil), []int(nil)))
	assert.Nil(t, Merge([]int{}, []int{}))
}

func TestSum(t *testing.T) {
	result := Sum(numerals)
	assert.Equal(t, 45, result)
}

func TestSum2(t *testing.T) {
	result := Sum(numeralsWithUserDefinedType)
	assert.Equal(t, MyInt(45), result)
}

func TestRemove(t *testing.T) {
	testSlice := []int{1, 2, 3}
	result := Remove(testSlice, 1)
	assert.Equal(t, []int{1, 3}, result)
	assert.Equal(t, []int{1, 2, 3}, testSlice)
	result = Remove(result, 1)
	assert.Equal(t, []int{1}, result)
	result = Remove(result, 3)
	assert.Equal(t, []int{1}, result)
	result = Remove(result, 0)
	assert.Equal(t, []int{}, result)
	result = Remove(result, 1)
	assert.Equal(t, []int{}, result)
}

func TestCopy(t *testing.T) {
	testSlice := []int{1, 2, 3}
	copiedSlice := Copy(testSlice)
	copiedSlice[0] = 2
	assert.NotEqual(t, testSlice, copiedSlice)
}

func TestInsert(t *testing.T) {
	testSlice := []int{1, 2}
	result := Insert(testSlice, 0, 3)
	assert.Equal(t, []int{3, 1, 2}, result)
	assert.NotEqual(t, testSlice, result)
	assert.Equal(t, []int{1, 3, 2}, Insert(testSlice, 1, 3))
	assert.Equal(t, []int{1, 2, 3}, Insert(testSlice, 2, 3))
}

func TestIntersection(t *testing.T) {
	expectedResult := []int{3, 4, 5}

	first := []int{1, 2, 3, 4, 5}
	second := []int{2, 3, 4, 5, 6}
	third := []int{3, 4, 5, 6, 7}

	assert.Equal(t, expectedResult, Intersection(first, second, third))
}

func TestDifference(t *testing.T) {
	expectedResult := []int{1, 7}

	first := []int{1, 2, 3, 4, 5}
	second := []int{2, 3, 4, 5, 6}
	third := []int{3, 4, 5, 6, 7}

	assert.Equal(t, expectedResult, Difference(first, second, third))
}

func TestUnion(t *testing.T) {
	expectedResult := []int{1, 2, 3, 4, 5, 6, 7}

	first := []int{1, 2, 3, 4, 5}
	second := []int{2, 3, 4, 5, 6}
	third := []int{3, 4, 5, 6, 7}

	assert.Equal(t, expectedResult, Union(first, second, third))
}

func TestReverse(t *testing.T) {
	expectedResult := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	assert.Equal(t, expectedResult, Reverse(numerals))
	// ensure does not modify the original
	assert.Equal(t, expectedResult, Reverse(numerals))

	// test basic odd length case
	expectedResult = []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	assert.Equal(t, expectedResult, Reverse(numerals[1:]))
}

func TestUnique(t *testing.T) {
	duplicates := []int{6, 6, 6, 9, 0, 0, 0}
	expectedResult := []int{6, 9, 0}
	assert.Equal(t, expectedResult, Unique(duplicates))
	// Ensure original is unaltered
	assert.NotEqual(t, expectedResult, duplicates)
}

func TestChunk(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}, Chunk(numbers, 2))
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}, Chunk(numbers, 3))
}

func TestPluck(t *testing.T) {
	items := []Pluckable{
		{
			Code:  "azer",
			Value: "Azer",
		},
		{
			Code:  "tyuio",
			Value: "Tyuio",
		},
	}

	assert.Equal(t, []string{"azer", "tyuio"}, Pluck(items, func(item Pluckable) *string {
		return &item.Code
	}))
	assert.Equal(t, []string{"Azer", "Tyuio"}, Pluck(items, func(item Pluckable) *string {
		return &item.Value
	}))
}

func TestFlatten(t *testing.T) {
	items := [][]int{
		{1, 2, 3, 4},
		{5, 6},
		{7, 8},
		{9, 10, 11},
	}

	flattened := Flatten(items)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, flattened)

	assert.Nil(t, Flatten([][]int{}))
	assert.Nil(t, Flatten([][]int(nil)))

	assert.Nil(t, Flatten([][]int{{}, {}}))
	assert.Nil(t, Flatten([][]int{nil, nil}))

	assert.Nil(t, Flatten([][]int{{}, nil}))
}
