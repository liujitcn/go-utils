// Package set 提供泛型集合能力。
package set

import (
	"encoding/json"
	"errors"
	"maps"
)

// Set 表示一组不重复的可比较元素，覆盖 golang-set/v2 v2.7.0 接口并保留本仓库既有便捷方法。
type Set[T comparable] interface {
	// Add 向集合中添加元素，并返回元素是否为新增。
	Add(value T) bool
	// Append 向集合中批量添加元素，并返回新增元素数量。
	Append(values ...T) int
	// Cardinality 返回集合元素数量。
	Cardinality() int
	// Clear 清空集合中的所有元素。
	Clear()
	// Clone 复制集合，返回与原集合并发安全属性一致的新集合。
	Clone() Set[T]
	// Contains 判断集合是否包含全部指定元素。
	Contains(values ...T) bool
	// ContainsOne 判断集合是否包含指定元素。
	ContainsOne(value T) bool
	// ContainsAll 判断集合是否包含全部指定元素。
	ContainsAll(values ...T) bool
	// ContainsAny 判断集合是否包含任意一个指定元素。
	ContainsAny(values ...T) bool
	// ContainsAnyElement 判断当前集合与另一个集合是否存在任意相同元素。
	ContainsAnyElement(other Set[T]) bool
	// Difference 计算当前集合相对另一个集合的差集。
	Difference(other Set[T]) Set[T]
	// Equal 判断两个集合是否元素完全一致。
	Equal(other Set[T]) bool
	// Intersect 计算当前集合与另一个集合的交集。
	Intersect(other Set[T]) Set[T]
	// Intersection 计算当前集合与另一个集合的交集。
	Intersection(other Set[T]) Set[T]
	// IsEmpty 判断集合是否为空。
	IsEmpty() bool
	// IsProperSubset 判断当前集合是否为另一个集合的真子集。
	IsProperSubset(other Set[T]) bool
	// IsProperSuperset 判断当前集合是否为另一个集合的真超集。
	IsProperSuperset(other Set[T]) bool
	// IsSubset 判断当前集合是否为另一个集合的子集。
	IsSubset(other Set[T]) bool
	// IsSuperset 判断当前集合是否为另一个集合的超集。
	IsSuperset(other Set[T]) bool
	// Each 遍历集合元素，回调返回 true 时提前停止。
	Each(function func(value T) bool)
	// Iter 返回可遍历集合元素的通道。
	Iter() <-chan T
	// Iterator 返回集合元素迭代器。
	Iterator() *Iterator[T]
	// Len 返回集合元素数量。
	Len() int
	// Remove 从集合中删除元素。
	Remove(value T)
	// RemoveAll 从集合中批量删除元素。
	RemoveAll(values ...T)
	// String 返回集合的字符串表示。
	String() string
	// SymmetricDifference 计算两个集合互不共有的元素集合。
	SymmetricDifference(other Set[T]) Set[T]
	// Union 计算当前集合与另一个集合的并集。
	Union(other Set[T]) Set[T]
	// Pop 删除并返回任意一个元素。
	Pop() (T, bool)
	// ToSlice 将集合转换为切片，返回顺序不保证稳定。
	ToSlice() []T
	// ToSliceByOrder 按指定顺序将集合转换为切片，不在顺序切片中的元素会被忽略。
	ToSliceByOrder(order []T) []T
	// MarshalJSON 将集合编码为 JSON 数组。
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON 从 JSON 数组解码集合。
	UnmarshalJSON(data []byte) error
}

type threadUnsafeSet[T comparable] struct {
	items map[T]struct{}
}

var _ Set[int] = (*threadUnsafeSet[int])(nil)

// New 创建非线程安全集合，并写入初始元素。
func New[T comparable](values ...T) Set[T] {
	return NewThreadUnsafeSet(values...)
}

// NewWithSize 按指定容量创建非线程安全集合。
func NewWithSize[T comparable](size int) Set[T] {
	return NewThreadUnsafeSetWithSize[T](size)
}

// FromSlice 根据切片创建非线程安全集合。
func FromSlice[T comparable](values []T) Set[T] {
	return New(values...)
}

// FromMapKeys 根据 map 的键创建非线程安全集合。
func FromMapKeys[T comparable, V any](values map[T]V) Set[T] {
	return NewThreadUnsafeSetFromMapKeys(values)
}

// NewThreadUnsafeSet 创建非线程安全集合，并写入初始元素。
func NewThreadUnsafeSet[T comparable](values ...T) Set[T] {
	collection := newThreadUnsafeSetWithSize[T](len(values))
	collection.Append(values...)
	return collection
}

// NewThreadUnsafeSetWithSize 按指定容量创建非线程安全集合。
func NewThreadUnsafeSetWithSize[T comparable](size int) Set[T] {
	return newThreadUnsafeSetWithSize[T](size)
}

// NewThreadUnsafeSetFromMapKeys 根据 map 的键创建非线程安全集合。
func NewThreadUnsafeSetFromMapKeys[T comparable, V any](values map[T]V) Set[T] {
	collection := newThreadUnsafeSetWithSize[T](len(values))
	for value := range maps.Keys(values) {
		collection.Add(value)
	}
	return collection
}

// newThreadUnsafeSetWithSize 按指定容量创建非线程安全集合实现。
func newThreadUnsafeSetWithSize[T comparable](size int) *threadUnsafeSet[T] {
	if size < 0 {
		size = 0
	}
	return &threadUnsafeSet[T]{items: make(map[T]struct{}, size)}
}

// Add 向集合中添加元素，并返回元素是否为新增。
func (s *threadUnsafeSet[T]) Add(value T) bool {
	if s == nil {
		return false
	}
	s.ensureItems(1)

	if _, exists := s.items[value]; exists {
		return false
	}
	s.items[value] = struct{}{}
	return true
}

// Append 向集合中批量添加元素，并返回新增元素数量。
func (s *threadUnsafeSet[T]) Append(values ...T) int {
	if s == nil || len(values) == 0 {
		return 0
	}
	s.ensureItems(len(values))

	addedCount := 0
	for _, value := range values {
		if _, exists := s.items[value]; exists {
			continue
		}
		s.items[value] = struct{}{}
		addedCount++
	}
	return addedCount
}

// Cardinality 返回集合元素数量。
func (s *threadUnsafeSet[T]) Cardinality() int {
	return s.Len()
}

// Len 返回集合元素数量。
func (s *threadUnsafeSet[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.items)
}

// Clear 清空集合中的所有元素。
func (s *threadUnsafeSet[T]) Clear() {
	if s == nil {
		return
	}
	for value := range s.items {
		delete(s.items, value)
	}
}

// Clone 复制集合，返回非线程安全集合。
func (s *threadUnsafeSet[T]) Clone() Set[T] {
	cloned := newThreadUnsafeSetWithSize[T](s.Cardinality())
	if s == nil {
		return cloned
	}

	for value := range s.items {
		cloned.add(value)
	}
	return cloned
}

// Contains 判断集合是否包含全部指定元素。
func (s *threadUnsafeSet[T]) Contains(values ...T) bool {
	if len(values) == 0 {
		return true
	}
	if s == nil || len(s.items) == 0 {
		return false
	}

	for _, value := range values {
		if _, exists := s.items[value]; !exists {
			return false
		}
	}
	return true
}

// ContainsOne 判断集合是否包含指定元素。
func (s *threadUnsafeSet[T]) ContainsOne(value T) bool {
	if s == nil || len(s.items) == 0 {
		return false
	}
	_, exists := s.items[value]
	return exists
}

// ContainsAll 判断集合是否包含全部指定元素。
func (s *threadUnsafeSet[T]) ContainsAll(values ...T) bool {
	return s.Contains(values...)
}

// ContainsAny 判断集合是否包含任意一个指定元素。
func (s *threadUnsafeSet[T]) ContainsAny(values ...T) bool {
	if s == nil || len(values) == 0 || len(s.items) == 0 {
		return false
	}

	for _, value := range values {
		if _, exists := s.items[value]; exists {
			return true
		}
	}
	return false
}

// ContainsAnyElement 判断当前集合与另一个集合是否存在任意相同元素。
func (s *threadUnsafeSet[T]) ContainsAnyElement(other Set[T]) bool {
	if s == nil || s.IsEmpty() || setCardinality(other) == 0 {
		return false
	}

	if s.Cardinality() < other.Cardinality() {
		for value := range s.items {
			if other.ContainsOne(value) {
				return true
			}
		}
		return false
	}

	for _, value := range other.ToSlice() {
		if s.ContainsOne(value) {
			return true
		}
	}
	return false
}

// Difference 计算当前集合相对另一个集合的差集。
func (s *threadUnsafeSet[T]) Difference(other Set[T]) Set[T] {
	if s == nil || s.IsEmpty() {
		return NewThreadUnsafeSet[T]()
	}
	if setCardinality(other) == 0 {
		return s.Clone()
	}

	differenceSet := newThreadUnsafeSetWithSize[T](s.Cardinality())
	for value := range s.items {
		if !other.ContainsOne(value) {
			differenceSet.add(value)
		}
	}
	return differenceSet
}

// Equal 判断两个集合是否元素完全一致。
func (s *threadUnsafeSet[T]) Equal(other Set[T]) bool {
	if s.Cardinality() != setCardinality(other) {
		return false
	}
	return s.IsSubset(other)
}

// Intersect 计算当前集合与另一个集合的交集。
func (s *threadUnsafeSet[T]) Intersect(other Set[T]) Set[T] {
	if s == nil || s.IsEmpty() || setCardinality(other) == 0 {
		return NewThreadUnsafeSet[T]()
	}

	leftValues := s.ToSlice()
	rightValues := other.ToSlice()
	if len(leftValues) > len(rightValues) {
		leftValues, rightValues = rightValues, leftValues
	}

	rightSet := NewThreadUnsafeSet(rightValues...)
	intersectionSet := newThreadUnsafeSetWithSize[T](len(leftValues))
	for _, value := range leftValues {
		if rightSet.ContainsOne(value) {
			intersectionSet.add(value)
		}
	}
	return intersectionSet
}

// Intersection 计算当前集合与另一个集合的交集。
func (s *threadUnsafeSet[T]) Intersection(other Set[T]) Set[T] {
	return s.Intersect(other)
}

// IsEmpty 判断集合是否为空。
func (s *threadUnsafeSet[T]) IsEmpty() bool {
	return s.Cardinality() == 0
}

// IsProperSubset 判断当前集合是否为另一个集合的真子集。
func (s *threadUnsafeSet[T]) IsProperSubset(other Set[T]) bool {
	return s.Cardinality() < setCardinality(other) && s.IsSubset(other)
}

// IsProperSuperset 判断当前集合是否为另一个集合的真超集。
func (s *threadUnsafeSet[T]) IsProperSuperset(other Set[T]) bool {
	return s.Cardinality() > setCardinality(other) && s.IsSuperset(other)
}

// IsSubset 判断当前集合是否为另一个集合的子集。
func (s *threadUnsafeSet[T]) IsSubset(other Set[T]) bool {
	if s == nil || len(s.items) == 0 {
		return true
	}
	if s.Cardinality() > setCardinality(other) {
		return false
	}

	for value := range s.items {
		if !other.ContainsOne(value) {
			return false
		}
	}
	return true
}

// IsSuperset 判断当前集合是否为另一个集合的超集。
func (s *threadUnsafeSet[T]) IsSuperset(other Set[T]) bool {
	if setCardinality(other) == 0 {
		return true
	}
	return other.IsSubset(s)
}

// Each 遍历集合元素，回调返回 true 时提前停止。
func (s *threadUnsafeSet[T]) Each(function func(value T) bool) {
	if s == nil || function == nil {
		return
	}

	for value := range s.items {
		if stop := function(value); stop {
			return
		}
	}
}

// Iter 返回可遍历集合元素的通道。
func (s *threadUnsafeSet[T]) Iter() <-chan T {
	return newIterChannel(s.ToSlice())
}

// Iterator 返回集合元素迭代器。
func (s *threadUnsafeSet[T]) Iterator() *Iterator[T] {
	return newIterator(s.ToSlice())
}

// Remove 从集合中删除元素。
func (s *threadUnsafeSet[T]) Remove(value T) {
	if s == nil || len(s.items) == 0 {
		return
	}
	delete(s.items, value)
}

// RemoveAll 从集合中批量删除元素。
func (s *threadUnsafeSet[T]) RemoveAll(values ...T) {
	if s == nil || len(values) == 0 {
		return
	}

	for _, value := range values {
		delete(s.items, value)
	}
}

// String 返回集合的字符串表示。
func (s *threadUnsafeSet[T]) String() string {
	return formatSetString(s.ToSlice())
}

// SymmetricDifference 计算两个集合互不共有的元素集合。
func (s *threadUnsafeSet[T]) SymmetricDifference(other Set[T]) Set[T] {
	if other == nil {
		return s.Clone()
	}

	symmetricSet := newThreadUnsafeSetWithSize[T](s.Cardinality() + other.Cardinality())
	for _, value := range s.ToSlice() {
		if !other.ContainsOne(value) {
			symmetricSet.add(value)
		}
	}
	for _, value := range other.ToSlice() {
		if !s.ContainsOne(value) {
			symmetricSet.add(value)
		}
	}
	return symmetricSet
}

// Union 计算当前集合与另一个集合的并集。
func (s *threadUnsafeSet[T]) Union(other Set[T]) Set[T] {
	unionSet := newThreadUnsafeSetWithSize[T](s.Cardinality() + setCardinality(other))
	if s != nil {
		for value := range s.items {
			unionSet.add(value)
		}
	}
	for _, value := range setValues(other) {
		unionSet.add(value)
	}
	return unionSet
}

// Pop 删除并返回任意一个元素。
func (s *threadUnsafeSet[T]) Pop() (T, bool) {
	if s == nil || len(s.items) == 0 {
		var zero T
		return zero, false
	}

	for value := range s.items {
		delete(s.items, value)
		return value, true
	}

	var zero T
	return zero, false
}

// ToSlice 将集合转换为切片，返回顺序不保证稳定。
func (s *threadUnsafeSet[T]) ToSlice() []T {
	values := make([]T, 0, s.Cardinality())
	if s == nil {
		return values
	}

	for value := range s.items {
		values = append(values, value)
	}
	return values
}

// ToSliceByOrder 按指定顺序将集合转换为切片，不在顺序切片中的元素会被忽略。
func (s *threadUnsafeSet[T]) ToSliceByOrder(order []T) []T {
	values := make([]T, 0, s.Cardinality())
	if s == nil || len(order) == 0 {
		return values
	}

	visited := NewThreadUnsafeSetWithSize[T](len(order))
	for _, value := range order {
		if s.ContainsOne(value) && visited.Add(value) {
			values = append(values, value)
		}
	}
	return values
}

// MarshalJSON 将集合编码为 JSON 数组。
func (s *threadUnsafeSet[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
}

// UnmarshalJSON 从 JSON 数组解码集合。
func (s *threadUnsafeSet[T]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("set: nil threadUnsafeSet 无法反序列化 JSON")
	}

	var values []T
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	s.items = make(map[T]struct{}, len(values))
	s.Append(values...)
	return nil
}

// add 向集合中添加元素，不返回新增状态。
func (s *threadUnsafeSet[T]) add(value T) {
	if s == nil {
		return
	}
	s.ensureItems(1)
	s.items[value] = struct{}{}
}

// ensureItems 初始化底层 map，保证零值集合也可写入。
func (s *threadUnsafeSet[T]) ensureItems(size int) {
	if s.items != nil {
		return
	}
	if size < 0 {
		size = 0
	}
	s.items = make(map[T]struct{}, size)
}
