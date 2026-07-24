package trans

import "time"

// Enum 返回枚举值指针。
func Enum[T any](v T) *T {
	return &v
}

// EnumValue 解引用枚举指针，空指针时返回 nil。
func EnumValue[T *any](v *T) T {
	if v == nil {
		return nil
	}
	return *v
}

// String 返回 string 指针。
func String(a string) *string {
	return &a
}

// StringValue 解引用 string 指针，空指针时返回零值。
func StringValue(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

// Int 返回 int 指针。
func Int(a int) *int {
	return &a
}

// IntValue 解引用 int 指针，空指针时返回零值。
func IntValue(a *int) int {
	if a == nil {
		return 0
	}
	return *a
}

// Int8 返回 int8 指针。
func Int8(a int8) *int8 {
	return &a
}

// Int8Value 解引用 int8 指针，空指针时返回零值。
func Int8Value(a *int8) int8 {
	if a == nil {
		return 0
	}
	return *a
}

// Int16 返回 int16 指针。
func Int16(a int16) *int16 {
	return &a
}

// Int16Value 解引用 int16 指针，空指针时返回零值。
func Int16Value(a *int16) int16 {
	if a == nil {
		return 0
	}
	return *a
}

// Int32 返回 int32 指针。
func Int32(a int32) *int32 {
	return &a
}

// Int32Value 解引用 int32 指针，空指针时返回零值。
func Int32Value(a *int32) int32 {
	if a == nil {
		return 0
	}
	return *a
}

// Int64 返回 int64 指针。
func Int64(a int64) *int64 {
	return &a
}

// Int64Value 解引用 int64 指针，空指针时返回零值。
func Int64Value(a *int64) int64 {
	if a == nil {
		return 0
	}
	return *a
}

// Bool 返回 bool 指针。
func Bool(a bool) *bool {
	return &a
}

// BoolValue 解引用 bool 指针，空指针时返回零值。
func BoolValue(a *bool) bool {
	if a == nil {
		return false
	}
	return *a
}

// Uint 返回 uint 指针。
func Uint(a uint) *uint {
	return &a
}

// UintValue 解引用 uint 指针，空指针时返回零值。
func UintValue(a *uint) uint {
	if a == nil {
		return 0
	}
	return *a
}

// Uint8 返回 uint8 指针。
func Uint8(a uint8) *uint8 {
	return &a
}

// Uint8Value 解引用 uint8 指针，空指针时返回零值。
func Uint8Value(a *uint8) uint8 {
	if a == nil {
		return 0
	}
	return *a
}

// Uint16 返回 uint16 指针。
func Uint16(a uint16) *uint16 {
	return &a
}

// Uint16Value 解引用 uint16 指针，空指针时返回零值。
func Uint16Value(a *uint16) uint16 {
	if a == nil {
		return 0
	}
	return *a
}

// Uint32 返回 uint32 指针。
func Uint32(a uint32) *uint32 {
	return &a
}

// Uint32Value 解引用 uint32 指针，空指针时返回零值。
func Uint32Value(a *uint32) uint32 {
	if a == nil {
		return 0
	}
	return *a
}

// Uint64 返回 uint64 指针。
func Uint64(a uint64) *uint64 {
	return &a
}

// Uint64Value 解引用 uint64 指针，空指针时返回零值。
func Uint64Value(a *uint64) uint64 {
	if a == nil {
		return 0
	}
	return *a
}

// Float32 返回 float32 指针。
func Float32(a float32) *float32 {
	return &a
}

// Float32Value 解引用 float32 指针，空指针时返回零值。
func Float32Value(a *float32) float32 {
	if a == nil {
		return 0
	}
	return *a
}

// Float64 返回 float64 指针。
func Float64(a float64) *float64 {
	return &a
}

// Float64Value 解引用 float64 指针，空指针时返回零值。
func Float64Value(a *float64) float64 {
	if a == nil {
		return 0
	}
	return *a
}

// Time 返回 time.Time 指针。
func Time(a time.Time) *time.Time {
	return &a
}

// TimeValue 解引用 time.Time 指针，空指针时返回当前时间。
func TimeValue(a *time.Time) time.Time {
	if a == nil {
		return time.Now()
	}
	return *a
}

// IntSlice 将 int 切片转换为 int 指针切片。
func IntSlice(a []int) []*int {
	if a == nil {
		return nil
	}
	res := make([]*int, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// IntValueSlice 将 int 指针切片转换为 int 切片。
func IntValueSlice(a []*int) []int {
	if a == nil {
		return nil
	}
	res := make([]int, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int8Slice 将 int8 切片转换为 int8 指针切片。
func Int8Slice(a []int8) []*int8 {
	if a == nil {
		return nil
	}
	res := make([]*int8, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int8ValueSlice 将 int8 指针切片转换为 int8 切片。
func Int8ValueSlice(a []*int8) []int8 {
	if a == nil {
		return nil
	}
	res := make([]int8, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int16Slice 将 int16 切片转换为 int16 指针切片。
func Int16Slice(a []int16) []*int16 {
	if a == nil {
		return nil
	}
	res := make([]*int16, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int16ValueSlice 将 int16 指针切片转换为 int16 切片。
func Int16ValueSlice(a []*int16) []int16 {
	if a == nil {
		return nil
	}
	res := make([]int16, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int32Slice 将 int32 切片转换为 int32 指针切片。
func Int32Slice(a []int32) []*int32 {
	if a == nil {
		return nil
	}
	res := make([]*int32, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int32ValueSlice 将 int32 指针切片转换为 int32 切片。
func Int32ValueSlice(a []*int32) []int32 {
	if a == nil {
		return nil
	}
	res := make([]int32, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int64Slice 将 int64 切片转换为 int64 指针切片。
func Int64Slice(a []int64) []*int64 {
	if a == nil {
		return nil
	}
	res := make([]*int64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int64ValueSlice 将 int64 指针切片转换为 int64 切片。
func Int64ValueSlice(a []*int64) []int64 {
	if a == nil {
		return nil
	}
	res := make([]int64, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// UintSlice 将 uint 切片转换为 uint 指针切片。
func UintSlice(a []uint) []*uint {
	if a == nil {
		return nil
	}
	res := make([]*uint, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// UintValueSlice 将 uint 指针切片转换为 uint 切片。
func UintValueSlice(a []*uint) []uint {
	if a == nil {
		return nil
	}
	res := make([]uint, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint8Slice 将 uint8 切片转换为 uint8 指针切片。
func Uint8Slice(a []uint8) []*uint8 {
	if a == nil {
		return nil
	}
	res := make([]*uint8, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint8ValueSlice 将 uint8 指针切片转换为 uint8 切片。
func Uint8ValueSlice(a []*uint8) []uint8 {
	if a == nil {
		return nil
	}
	res := make([]uint8, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint16Slice 将 uint16 切片转换为 uint16 指针切片。
func Uint16Slice(a []uint16) []*uint16 {
	if a == nil {
		return nil
	}
	res := make([]*uint16, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint16ValueSlice 将 uint16 指针切片转换为 uint16 切片。
func Uint16ValueSlice(a []*uint16) []uint16 {
	if a == nil {
		return nil
	}
	res := make([]uint16, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint32Slice 将 uint32 切片转换为 uint32 指针切片。
func Uint32Slice(a []uint32) []*uint32 {
	if a == nil {
		return nil
	}
	res := make([]*uint32, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint32ValueSlice 将 uint32 指针切片转换为 uint32 切片。
func Uint32ValueSlice(a []*uint32) []uint32 {
	if a == nil {
		return nil
	}
	res := make([]uint32, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint64Slice 将 uint64 切片转换为 uint64 指针切片。
func Uint64Slice(a []uint64) []*uint64 {
	if a == nil {
		return nil
	}
	res := make([]*uint64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint64ValueSlice 将 uint64 指针切片转换为 uint64 切片。
func Uint64ValueSlice(a []*uint64) []uint64 {
	if a == nil {
		return nil
	}
	res := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Float32Slice 将 float32 切片转换为 float32 指针切片。
func Float32Slice(a []float32) []*float32 {
	if a == nil {
		return nil
	}
	res := make([]*float32, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Float32ValueSlice 将 float32 指针切片转换为 float32 切片。
func Float32ValueSlice(a []*float32) []float32 {
	if a == nil {
		return nil
	}
	res := make([]float32, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Float64Slice 将 float64 切片转换为 float64 指针切片。
func Float64Slice(a []float64) []*float64 {
	if a == nil {
		return nil
	}
	res := make([]*float64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Float64ValueSlice 将 float64 指针切片转换为 float64 切片。
func Float64ValueSlice(a []*float64) []float64 {
	if a == nil {
		return nil
	}
	res := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// StringSlice 将 string 切片转换为 string 指针切片。
func StringSlice(a []string) []*string {
	if a == nil {
		return nil
	}
	res := make([]*string, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// StringSliceValue 将 string 指针切片转换为 string 切片。
func StringSliceValue(a []*string) []string {
	if a == nil {
		return nil
	}
	res := make([]string, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// BoolSlice 将 bool 切片转换为 bool 指针切片。
func BoolSlice(a []bool) []*bool {
	if a == nil {
		return nil
	}
	res := make([]*bool, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// BoolSliceValue 将 bool 指针切片转换为 bool 切片。
func BoolSliceValue(a []*bool) []bool {
	if a == nil {
		return nil
	}
	res := make([]bool, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

type mapKeyValueType interface {
	~string | ~bool |
		~int | ~int8 | ~int16 | ~int32 |
		~uint | ~uint8 | ~uint16 | ~uint32 |
		~float32 | ~float64
}

// MapKeys 获取map的键
func MapKeys[TKey mapKeyValueType, TValue mapKeyValueType](source map[TKey]TValue) []TKey {
	var target []TKey
	for k := range source {
		target = append(target, k)
	}
	return target
}

// MapValues 获取map的值
func MapValues[TKey mapKeyValueType, TValue mapKeyValueType](source map[TKey]TValue) []TValue {
	var target []TValue
	for _, v := range source {
		target = append(target, v)
	}
	return target
}
