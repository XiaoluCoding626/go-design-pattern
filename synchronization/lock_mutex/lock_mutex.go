package lockmutex

import (
	"fmt"
	"sync"
)

// Mutex 定义互斥锁接口
// 通过接口可以方便地替换不同的互斥锁实现或用于模拟测试
type Mutex interface {
	Lock()   // 获取锁
	Unlock() // 释放锁
}

// StandardMutex 标准互斥锁实现，封装Go的sync.Mutex
type StandardMutex struct {
	mu sync.Mutex
}

// Lock 获取互斥锁
func (m *StandardMutex) Lock() {
	m.mu.Lock()
}

// Unlock 释放互斥锁
func (m *StandardMutex) Unlock() {
	m.mu.Unlock()
}

// Counter 计数器结构体，使用互斥锁保护共享数据
type Counter struct {
	mutex Mutex // 使用接口而非具体类型，提高灵活性
	count int   // 计数器的当前值
}

// NewCounter 创建一个新的计数器实例，使用标准互斥锁
func NewCounter() *Counter {
	return &Counter{
		mutex: &StandardMutex{},
	}
}

// NewCounterWithMutex 使用指定的互斥锁创建计数器
// 这允许注入自定义的锁实现，例如用于测试
func NewCounterWithMutex(mutex Mutex) *Counter {
	return &Counter{
		mutex: mutex,
	}
}

// Increment 增加计数器的方法，受互斥锁保护
func (c *Counter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock() // 确保在任何情况下都会释放锁

	c.count++
	fmt.Println("增加后的计数:", c.count)
}

// GetCount 获取当前计数值，同样受到互斥锁保护
func (c *Counter) GetCount() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.count
}

// SafeOperation 在锁的保护下执行操作
// 这种模式确保操作总是在锁的保护下执行，避免忘记加锁或解锁
func (c *Counter) SafeOperation(operation func(*Counter)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	operation(c)
}

// Reset 重置计数器值为0
func (c *Counter) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count = 0
	fmt.Println("计数器已重置")
}
