package slice

import (
	"sort"

	"github.com/liujitcn/go-utils/str"
)

func EqualByString(ids1Str, ids2Str string) bool {
	ids1 := str.ConvertJsonStringToInt64Array(ids1Str)
	ids2 := str.ConvertJsonStringToInt64Array(ids2Str)
	return Equal(ids1, ids2)
}

func Equal(ids1, ids2 []int64) bool {
	if len(ids1) != len(ids2) {
		return false
	}
	// 排序
	sort.Slice(ids1, func(i, j int) bool {
		return ids1[i] < ids1[j]
	})
	sort.Slice(ids2, func(i, j int) bool {
		return ids2[i] < ids2[j]
	})
	for i, item := range ids1 {
		if item != ids2[i] {
			return false
		}
	}
	return true
}
