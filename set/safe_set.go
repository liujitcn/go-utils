package set

import "sync"

type threadSafeSet[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

var _ Set[int] = (*threadSafeSet[int])(nil)

// NewSet 创建线程安全集合，并写入初始元素。
func NewSet[T comparable](values ...T) Set[T] {
	collection := newThreadSafeSetWithSize[T](len(values))
	collection.Append(values...)
	return collection
}

// NewSetWithSize 按指定容量创建线程安全集合。
func NewSetWithSize[T comparable](size int) Set[T] {
	return newThreadSafeSetWithSize[T](size)
}

// NewSetFromMapKeys 根据 map 的键创建线程安全集合。
func NewSetFromMapKeys[T comparable, V any](values map[T]V) Set[T] {
	collection := newThreadSafeSetWithSize[T](len(values))
	for value := range values {
		collection.Add(value)
	}
	return collection
}

// NewThreadSafeSet 创建线程安全集合，并写入初始元素。
func NewThreadSafeSet[T comparable](values ...T) Set[T] {
	return NewSet(values...)
}

// NewThreadSafeSetWithSize 按指定容量创建线程安全集合。
func NewThreadSafeSetWithSize[T comparable](size int) Set[T] {
	return NewSetWithSize[T](size)
}

// NewThreadSafeSetFromMapKeys 根据 map 的键创建线程安全集合。
func NewThreadSafeSetFromMapKeys[T comparable, V any](values map[T]V) Set[T] {
	return NewSetFromMapKeys(values)
}

// NewSafe 创建线程安全集合，并写入初始元素。
func NewSafe[T comparable](values ...T) Set[T] {
	return NewSet(values...)
}

// NewSafeWithSize 按指定容量创建线程安全集合。
func NewSafeWithSize[T comparable](size int) Set[T] {
	return NewSetWithSize[T](size)
}

// SafeFromSlice 根据切片创建线程安全集合。
func SafeFromSlice[T comparable](values []T) Set[T] {
	return NewSet(values...)
}

// SafeFromMapKeys 根据 map 的键创建线程安全集合。
func SafeFromMapKeys[T comparable, V any](values map[T]V) Set[T] {
	return NewSetFromMapKeys(values)
}

// newThreadSafeSetWithSize 按指定容量创建线程安全集合实现。
func newThreadSafeSetWithSize[T comparable](size int) *threadSafeSet[T] {
	if size < 0 {
		size = 0
	}
	return &threadSafeSet[T]{items: make(map[T]struct{}, size)}
}

// Add 向集合中添加元素，并返回元素是否为新增。
func (s *threadSafeSet[T]) Add(value T) bool {
	if s == nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureItemsLocked(1)
	if _, exists := s.items[value]; exists {
		return false
	}
	s.items[value] = struct{}{}
	return true
}

// Append 向集合中批量添加元素，并返回新增元素数量。
func (s *threadSafeSet[T]) Append(values ...T) int {
	if s == nil || len(values) == 0 {
		return 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureItemsLocked(len(values))
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
func (s *threadSafeSet[T]) Cardinality() int {
	return s.Len()
}

// Len 返回集合元素数量。
func (s *threadSafeSet[T]) Len() int {
	if s == nil {
		return 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

// Clear 清空集合中的所有元素。
func (s *threadSafeSet[T]) Clear() {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for value := range s.items {
		delete(s.items, value)
	}
}

// Clone 复制集合，返回线程安全集合。
func (s *threadSafeSet[T]) Clone() Set[T] {
	return NewSet(s.ToSlice()...)
}

// Contains 判断集合是否包含全部指定元素。
func (s *threadSafeSet[T]) Contains(values ...T) bool {
	if len(values) == 0 {
		return true
	}
	if s == nil {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.items) == 0 {
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
func (s *threadSafeSet[T]) ContainsOne(value T) bool {
	if s == nil {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.items) == 0 {
		return false
	}
	_, exists := s.items[value]
	return exists
}

// ContainsAll 判断集合是否包含全部指定元素。
func (s *threadSafeSet[T]) ContainsAll(values ...T) bool {
	return s.Contains(values...)
}

// ContainsAny 判断集合是否包含任意一个指定元素。
func (s *threadSafeSet[T]) ContainsAny(values ...T) bool {
	if s == nil || len(values) == 0 {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.items) == 0 {
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
func (s *threadSafeSet[T]) ContainsAnyElement(other Set[T]) bool {
	if s == nil || s.IsEmpty() || setCardinality(other) == 0 {
		return false
	}

	if s.Cardinality() < other.Cardinality() {
		for _, value := range s.ToSlice() {
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
func (s *threadSafeSet[T]) Difference(other Set[T]) Set[T] {
	if s == nil || s.IsEmpty() {
		return NewSet[T]()
	}
	if setCardinality(other) == 0 {
		return s.Clone()
	}

	differenceSet := NewSetWithSize[T](s.Cardinality())
	for _, value := range s.ToSlice() {
		if !other.ContainsOne(value) {
			differenceSet.Add(value)
		}
	}
	return differenceSet
}

// Equal 判断两个集合是否元素完全一致。
func (s *threadSafeSet[T]) Equal(other Set[T]) bool {
	if s.Cardinality() != setCardinality(other) {
		return false
	}
	return s.IsSubset(other)
}

// Intersect 计算当前集合与另一个集合的交集。
func (s *threadSafeSet[T]) Intersect(other Set[T]) Set[T] {
	leftValues := s.ToSlice()
	rightValues := setValues(other)
	if len(leftValues) == 0 || len(rightValues) == 0 {
		return NewSet[T]()
	}

	if len(leftValues) > len(rightValues) {
		leftValues, rightValues = rightValues, leftValues
	}

	rightSet := NewThreadUnsafeSet(rightValues...)
	intersectionSet := NewSetWithSize[T](len(leftValues))
	for _, value := range leftValues {
		if rightSet.ContainsOne(value) {
			intersectionSet.Add(value)
		}
	}
	return intersectionSet
}

// Intersection 计算当前集合与另一个集合的交集。
func (s *threadSafeSet[T]) Intersection(other Set[T]) Set[T] {
	return s.Intersect(other)
}

// IsEmpty 判断集合是否为空。
func (s *threadSafeSet[T]) IsEmpty() bool {
	return s.Cardinality() == 0
}

// IsProperSubset 判断当前集合是否为另一个集合的真子集。
func (s *threadSafeSet[T]) IsProperSubset(other Set[T]) bool {
	return s.Cardinality() < setCardinality(other) && s.IsSubset(other)
}

// IsProperSuperset 判断当前集合是否为另一个集合的真超集。
func (s *threadSafeSet[T]) IsProperSuperset(other Set[T]) bool {
	return s.Cardinality() > setCardinality(other) && s.IsSuperset(other)
}

// IsSubset 判断当前集合是否为另一个集合的子集。
func (s *threadSafeSet[T]) IsSubset(other Set[T]) bool {
	if s == nil || s.IsEmpty() {
		return true
	}
	if s.Cardinality() > setCardinality(other) {
		return false
	}

	for _, value := range s.ToSlice() {
		if !other.ContainsOne(value) {
			return false
		}
	}
	return true
}

// IsSuperset 判断当前集合是否为另一个集合的超集。
func (s *threadSafeSet[T]) IsSuperset(other Set[T]) bool {
	if setCardinality(other) == 0 {
		return true
	}
	return other.IsSubset(s)
}

// Each 遍历集合元素，回调返回 true 时提前停止。
func (s *threadSafeSet[T]) Each(function func(value T) bool) {
	if function == nil {
		return
	}

	for _, value := range s.ToSlice() {
		if stop := function(value); stop {
			return
		}
	}
}

// Iter 返回可遍历集合元素的通道。
func (s *threadSafeSet[T]) Iter() <-chan T {
	return newIterChannel(s.ToSlice())
}

// Iterator 返回集合元素迭代器。
func (s *threadSafeSet[T]) Iterator() *Iterator[T] {
	return newIterator(s.ToSlice())
}

// Remove 从集合中删除元素。
func (s *threadSafeSet[T]) Remove(value T) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		return
	}
	delete(s.items, value)
}

// RemoveAll 从集合中批量删除元素。
func (s *threadSafeSet[T]) RemoveAll(values ...T) {
	if s == nil || len(values) == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, value := range values {
		delete(s.items, value)
	}
}

// String 返回集合的字符串表示。
func (s *threadSafeSet[T]) String() string {
	return formatSetString(s.ToSlice())
}

// SymmetricDifference 计算两个集合互不共有的元素集合。
func (s *threadSafeSet[T]) SymmetricDifference(other Set[T]) Set[T] {
	if other == nil {
		return s.Clone()
	}

	symmetricSet := NewSetWithSize[T](s.Cardinality() + other.Cardinality())
	for _, value := range s.ToSlice() {
		if !other.ContainsOne(value) {
			symmetricSet.Add(value)
		}
	}
	for _, value := range other.ToSlice() {
		if !s.ContainsOne(value) {
			symmetricSet.Add(value)
		}
	}
	return symmetricSet
}

// Union 计算当前集合与另一个集合的并集。
func (s *threadSafeSet[T]) Union(other Set[T]) Set[T] {
	unionSet := NewSetWithSize[T](s.Cardinality() + setCardinality(other))
	unionSet.Append(s.ToSlice()...)
	unionSet.Append(setValues(other)...)
	return unionSet
}

// Pop 删除并返回任意一个元素。
func (s *threadSafeSet[T]) Pop() (T, bool) {
	if s == nil {
		var zero T
		return zero, false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
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
func (s *threadSafeSet[T]) ToSlice() []T {
	if s == nil {
		return []T{}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	values := make([]T, 0, len(s.items))
	for value := range s.items {
		values = append(values, value)
	}
	return values
}

// ToSliceByOrder 按指定顺序将集合转换为切片，不在顺序切片中的元素会被忽略。
func (s *threadSafeSet[T]) ToSliceByOrder(order []T) []T {
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
func (s *threadSafeSet[T]) MarshalJSON() ([]byte, error) {
	return jsonMarshalSet(s)
}

// UnmarshalJSON 从 JSON 数组解码集合。
func (s *threadSafeSet[T]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errNilSafeSetUnmarshal
	}

	values, err := jsonUnmarshalValues[T](data)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[T]struct{}, len(values))
	for _, value := range values {
		s.items[value] = struct{}{}
	}
	return nil
}

// ensureItemsLocked 初始化底层 map，调用方必须已持有写锁。
func (s *threadSafeSet[T]) ensureItemsLocked(size int) {
	if s.items != nil {
		return
	}
	if size < 0 {
		size = 0
	}
	s.items = make(map[T]struct{}, size)
}
