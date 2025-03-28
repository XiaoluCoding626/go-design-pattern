package monitor

import (
	"sync"
)

// Monitor 实现监视器模式，提供对共享资源的互斥访问和条件同步
type Monitor struct {
	mutex sync.Mutex // 互斥锁保护共享资源
	cond  *sync.Cond // 条件变量用于线程协作
	data  int        // 共享数据
	valid bool       // 数据有效性标志
}

// NewMonitor 创建监视器实例
func NewMonitor() *Monitor {
	m := &Monitor{}
	m.cond = sync.NewCond(&m.mutex)
	return m
}

// Produce 向监视器写入数据
func (m *Monitor) Produce(value int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 更新共享数据
	m.data = value
	m.valid = true

	// 通知所有等待的消费者
	m.cond.Broadcast()
}

// Consume 从监视器读取数据
func (m *Monitor) Consume() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 等待数据变为有效
	for !m.valid {
		m.cond.Wait()
	}

	// 获取数据并重置有效性标志
	value := m.data
	m.valid = false

	return value
}

// Reset 重置监视器状态
func (m *Monitor) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = 0
	m.valid = false
}

// ExecuteProtected 在互斥保护的情况下执行自定义函数
func (m *Monitor) ExecuteProtected(action func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	action()
}

// GetValue 安全地获取当前值
func (m *Monitor) GetValue() (int, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.data, m.valid
}
