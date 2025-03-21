package singleton

import "sync"

// 单例模式：饿汉式实现
//
// 这种实现方式在显式初始化时就创建单例实例。
// 与懒汉式相比，饿汉式提供了对初始化时机和参数的更多控制。

// 全局单例实例 - 在启动期间初始化
var eagerInstance *Eager

// Eager 表示一个早期初始化的单例实例。
type Eager struct {
	count int

	mu sync.Mutex // 添加互斥锁保护对 count 这个共享变量的并发访问
}

// Increase 增加 Eager 单例中的计数器值。
// 此方法是线程安全的。
func (e *Eager) Increase() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.count++
}

// GetCount 返回当前计数值。
// 此方法是线程安全的。
func (e *Eager) GetCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.count
}

// InitEager 使用指定的计数值初始化单例实例。
// 使用显式初始化函数而非 init() 提供了更好的
// 可见性和对单例创建时机与方式的控制。
func InitEager(count int) {
	eagerInstance = &Eager{count: count}
}

// GetEager 返回 Eager 的单例实例。
// 由于实例在任何并发访问之前就已初始化，因此该方法对读取操作是线程安全的。
// 注意：如果在调用 InitEager 之前调用此方法，将返回 nil。
func GetEager() *Eager {
	return eagerInstance
}
