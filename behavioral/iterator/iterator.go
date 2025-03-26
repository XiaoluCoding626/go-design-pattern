package iterator

// Iterator 定义迭代器接口，使用泛型支持不同类型
type Iterator[T any] interface {
	HasNext() bool      // 是否有下一个元素
	Next() (T, bool)    // 获取下一个元素，返回元素和是否成功
	Reset()             // 重置迭代器
	Current() (T, bool) // 获取当前元素
}

// ConcreteIterator 具体迭代器实现
type ConcreteIterator[T any] struct {
	index int // 迭代器当前位置
	data  []T // 数据集合
}

// NewConcreteIterator 创建新的具体迭代器实例
func NewConcreteIterator[T any](data []T) *ConcreteIterator[T] {
	return &ConcreteIterator[T]{index: 0, data: data}
}

// HasNext 实现迭代器接口，判断是否有下一个元素
func (it *ConcreteIterator[T]) HasNext() bool {
	return it.index < len(it.data)
}

// Next 实现迭代器接口，获取下一个元素
func (it *ConcreteIterator[T]) Next() (T, bool) {
	var zero T
	if !it.HasNext() {
		return zero, false
	}
	value := it.data[it.index]
	it.index++
	return value, true
}

// Reset 重置迭代器到初始位置
func (it *ConcreteIterator[T]) Reset() {
	it.index = 0
}

// Current 获取当前元素
func (it *ConcreteIterator[T]) Current() (T, bool) {
	var zero T
	if it.index <= 0 || it.index > len(it.data) {
		return zero, false
	}
	return it.data[it.index-1], true
}

// Aggregate 聚合对象接口
type Aggregate[T any] interface {
	CreateIterator() Iterator[T] // 创建迭代器
	Add(item T)                  // 添加元素
	Remove(index int) bool       // 移除元素
	Get(index int) (T, bool)     // 获取元素
	Count() int                  // 获取元素数量
}

// ConcreteAggregate 具体聚合对象实现
type ConcreteAggregate[T any] struct {
	data []T // 数据集合
}

// NewConcreteAggregate 创建新的具体聚合对象实例
func NewConcreteAggregate[T any]() *ConcreteAggregate[T] {
	return &ConcreteAggregate[T]{data: make([]T, 0)}
}

// CreateIterator 实现聚合对象接口，创建迭代器
func (a *ConcreteAggregate[T]) CreateIterator() Iterator[T] {
	return NewConcreteIterator[T](a.data)
}

// Add 添加元素到集合
func (a *ConcreteAggregate[T]) Add(item T) {
	a.data = append(a.data, item)
}

// Remove 从集合中移除指定索引的元素
func (a *ConcreteAggregate[T]) Remove(index int) bool {
	if index < 0 || index >= len(a.data) {
		return false
	}

	a.data = append(a.data[:index], a.data[index+1:]...)
	return true
}

// Get 获取指定索引的元素
func (a *ConcreteAggregate[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(a.data) {
		return zero, false
	}
	return a.data[index], true
}

// Count 获取集合中元素的数量
func (a *ConcreteAggregate[T]) Count() int {
	return len(a.data)
}
