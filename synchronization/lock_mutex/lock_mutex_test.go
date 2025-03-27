package lockmutex

import (
	"fmt"
	"sync"
	"testing"
)

// 测试Counter的基本功能
func TestCounterBasic(t *testing.T) {
	counter := NewCounter()

	// 测试初始值
	if counter.GetCount() != 0 {
		t.Errorf("初始计数应为0，但得到: %d", counter.GetCount())
	}

	// 测试增加操作
	counter.Increment()
	if counter.GetCount() != 1 {
		t.Errorf("增加后计数应为1，但得到: %d", counter.GetCount())
	}

	// 测试重置操作
	counter.Reset()
	if counter.GetCount() != 0 {
		t.Errorf("重置后计数应为0，但得到: %d", counter.GetCount())
	}
}

// 测试并发增加场景，验证互斥锁是否有效保护共享数据
func TestConcurrentIncrement(t *testing.T) {
	counter := NewCounter()
	numGoroutines := 100
	var wg sync.WaitGroup

	// 启动多个goroutine同时增加计数器
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 验证结果
	if count := counter.GetCount(); count != numGoroutines {
		t.Errorf("并发增加后计数应为%d，但得到: %d", numGoroutines, count)
	}
}

// 测试SafeOperation功能
func TestSafeOperation(t *testing.T) {
	counter := NewCounter()

	// 使用SafeOperation增加计数
	counter.SafeOperation(func(c *Counter) {
		c.count += 10 // 直接访问count，因为在SafeOperation内部已加锁
	})

	if counter.GetCount() != 10 {
		t.Errorf("SafeOperation后计数应为10，但得到: %d", counter.GetCount())
	}

	// 测试多个操作
	counter.SafeOperation(func(c *Counter) {
		c.count *= 2 // 乘2
		c.count -= 5 // 减5
	})

	if counter.GetCount() != 15 {
		t.Errorf("多个操作后计数应为15，但得到: %d", counter.GetCount())
	}
}

// 模拟互斥锁，用于测试依赖注入
type MockMutex struct {
	lockCalled   bool
	unlockCalled bool
}

func (m *MockMutex) Lock() {
	m.lockCalled = true
}

func (m *MockMutex) Unlock() {
	m.unlockCalled = true
}

// 测试依赖注入自定义互斥锁
func TestCustomMutex(t *testing.T) {
	mockMutex := &MockMutex{}
	counter := NewCounterWithMutex(mockMutex)

	counter.Increment()

	// 验证Lock和Unlock被调用
	if !mockMutex.lockCalled {
		t.Error("Lock方法未被调用")
	}

	if !mockMutex.unlockCalled {
		t.Error("Unlock方法未被调用")
	}
}

// 测试竞态条件 - 故意不使用互斥锁以演示问题
func TestRaceCondition(t *testing.T) {
	// 仅用于演示，跳过实际测试
	if testing.Short() {
		t.Skip("跳过竞态条件演示")
	}

	// 创建一个没有锁保护的计数器
	type UnsafeCounter struct {
		count int
	}

	unsafeCounter := &UnsafeCounter{}
	numGoroutines := 1000
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsafeCounter.count++ // 没有锁保护，会导致竞态条件
		}()
	}

	wg.Wait()
	fmt.Printf("不安全计数器的最终值: %d (预期: %d)\n", unsafeCounter.count, numGoroutines)

	// 这个测试很可能失败，因为没有使用互斥锁保护
	if unsafeCounter.count != numGoroutines {
		t.Logf("竞态条件导致计数不准确: %d != %d", unsafeCounter.count, numGoroutines)
	}
}

// 测试递增和重置交替操作的并发安全性
func TestConcurrentIncrementAndReset(t *testing.T) {
	counter := NewCounter()
	var wg sync.WaitGroup

	// 启动10个goroutine进行递增
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				counter.Increment()
			}
		}()
	}

	// 启动5个goroutine进行重置
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				counter.Reset()
			}
		}()
	}

	wg.Wait()
	// 注意：由于并发操作，最终结果是不确定的
	// 但重要的是操作应该是原子的，不会导致竞态条件
	t.Logf("并发递增和重置后的最终计数: %d", counter.GetCount())
}

// 示例展示如何使用互斥锁模式
func ExampleCounter() {
	counter := NewCounter()

	// 基本递增
	counter.Increment()
	counter.Increment()

	// 获取当前值
	fmt.Printf("当前计数: %d\n", counter.GetCount())

	// 安全操作
	counter.SafeOperation(func(c *Counter) {
		c.count += 3
		fmt.Printf("在安全操作内部的计数: %d\n", c.count)
	})

	// 重置计数器
	counter.Reset()

	fmt.Printf("重置后的计数: %d\n", counter.GetCount())

	// Output:
	// 增加后的计数: 1
	// 增加后的计数: 2
	// 当前计数: 2
	// 在安全操作内部的计数: 5
	// 计数器已重置
	// 重置后的计数: 0
}

// 基准测试
func BenchmarkIncrement(b *testing.B) {
	counter := NewCounter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.SafeOperation(func(c *Counter) {
			c.count++
		})
	}
}

func BenchmarkConcurrentIncrement(b *testing.B) {
	counter := NewCounter()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}
