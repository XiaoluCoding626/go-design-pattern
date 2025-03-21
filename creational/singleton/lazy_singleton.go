package singleton

import (
	"fmt"
	"sync"
)

// 单例模式：懒汉式实现
//
// 这种实现方式只在首次请求时才创建单例实例。
// 懒汉式在实例未被立即使用时可以节省资源，
// 但需要线程安全机制来处理并发初始化的问题。
var (
	// 全局单例实例 - 首次访问时才初始化
	lazyInstance *Lazy
	// 使用 sync.Once 确保单例只被初始化一次
	lazyOnce sync.Once
)

// Lazy 表示一个在首次使用时才初始化的单例实例。
type Lazy struct{}

// HelloWorld 演示 Lazy 单例的一个方法。
func (l *Lazy) HelloWorld() {
	fmt.Println("hello world")
}

// GetLazy 返回 Lazy 的单例实例。
// 它使用 sync.Once 确保线程安全，防止在并发初始化时出现竞争条件。
//
// 注意：虽然懒汉式初始化可以节省资源，但它会使代码的行为
// 变得不那么可预测，因为初始化发生在不确定的时间点。
func GetLazy() *Lazy {
	lazyOnce.Do(func() {
		lazyInstance = &Lazy{}
	})

	return lazyInstance
}
