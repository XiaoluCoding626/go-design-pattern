package read_write_lock

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 测试基本的读写操作
func TestBasicReadWrite(t *testing.T) {
	data := NewData()

	// 默认应该是0
	if got := data.Read(); got != 0 {
		t.Errorf("初始值应为0，但得到: %v", got)
	}

	// 写入100
	data.Write(100)

	// 验证读取结果
	if got := data.Read(); got != 100 {
		t.Errorf("期望值为100，但得到: %v", got)
	}
}

// 测试尝试读取功能
func TestTryRead(t *testing.T) {
	data := NewData()
	data.Write(42)

	// 应该能够成功读取
	val, ok := data.TryRead()
	if !ok {
		t.Error("TryRead应该成功，但失败了")
	}
	if val != 42 {
		t.Errorf("期望值为42，但得到: %v", val)
	}
}

// 测试尝试写入功能
func TestTryWrite(t *testing.T) {
	data := NewData()

	// 应该能够成功写入
	ok := data.TryWrite(50)
	if !ok {
		t.Error("TryWrite应该成功，但失败了")
	}

	if got := data.Read(); got != 50 {
		t.Errorf("期望值为50，但得到: %v", got)
	}
}

// 测试超时读取
func TestReadWithTimeout(t *testing.T) {
	data := NewData()
	data.Write(30)

	// 应该能够在超时前读取
	val, ok := data.ReadWithTimeout(100 * time.Millisecond)
	if !ok {
		t.Error("ReadWithTimeout应该成功，但失败了")
	}
	if val != 30 {
		t.Errorf("期望值为30，但得到: %v", val)
	}
}

// 测试超时写入
func TestWriteWithTimeout(t *testing.T) {
	data := NewData()

	// 应该能够在超时前写入
	ok := data.WriteWithTimeout(70, 100*time.Millisecond)
	if !ok {
		t.Error("WriteWithTimeout应该成功，但失败了")
	}

	if got := data.Read(); got != 70 {
		t.Errorf("期望值为70，但得到: %v", got)
	}
}

// 测试读写冲突
func TestReadWriteContention(t *testing.T) {
	// 创建一个共享锁来模拟冲突
	sharedLocker := NewStandardRWLock()
	data1 := NewDataWithLocker(sharedLocker)
	data2 := NewDataWithLocker(sharedLocker)

	// 获取data1的写锁
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// 持有写锁
		data1.WriteWithCallback(func(d *Data) {
			// 当data1持有写锁时，data2应该无法获取读锁或写锁
			time.Sleep(100 * time.Millisecond)
		})
	}()

	// 给写锁时间获取
	time.Sleep(10 * time.Millisecond)

	// 尝试读取data2，应该会失败因为data1持有写锁
	_, ok := data2.TryRead()
	if ok {
		t.Error("当另一个实例持有相同锁的写锁时，TryRead应该失败")
	}

	// 尝试写入data2，应该会失败因为data1持有写锁
	ok = data2.TryWrite(100)
	if ok {
		t.Error("当另一个实例持有相同锁的写锁时，TryWrite应该失败")
	}

	// 等待写锁释放
	wg.Wait()

	// 现在应该可以成功读取和写入
	_, ok = data2.TryRead()
	if !ok {
		t.Error("当写锁释放后，TryRead应该成功")
	}

	ok = data2.TryWrite(100)
	if !ok {
		t.Error("当写锁释放后，TryWrite应该成功")
	}
}

// 测试读写回调
func TestReadWriteCallbacks(t *testing.T) {
	data := NewData()
	data.Write(10)

	// 测试读回调
	var readValue int
	data.ReadWithCallback(func(val int) {
		readValue = val
	})

	if readValue != 10 {
		t.Errorf("读回调中的值应为10，但得到: %v", readValue)
	}

	// 测试写回调
	data.WriteWithCallback(func(d *Data) {
		d.value = 20
	})

	if got := data.Read(); got != 20 {
		t.Errorf("写回调后的值应为20，但得到: %v", got)
	}

	// 测试读写回调
	data.ReadWriteWithCallback(func(val int) int {
		return val * 2
	})

	if got := data.Read(); got != 40 {
		t.Errorf("读写回调后的值应为40，但得到: %v", got)
	}
}

// 测试并发读取
func TestConcurrentReads(t *testing.T) {
	data := NewData()
	data.Write(100)

	const readers = 100
	var wg sync.WaitGroup
	wg.Add(readers)

	// 启动多个并发读取goroutine
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			val := data.Read()
			if val != 100 {
				t.Errorf("并发读取值应为100，但得到: %v", val)
			}
		}()
	}

	wg.Wait()
}

// 测试读写互斥
func TestReadWriteMutualExclusion(t *testing.T) {
	data := NewData()

	// 原子变量用来检测读写互斥
	var readsDuringWrite atomic.Int64

	// 用channel作为同步机制
	writingStarted := make(chan struct{})
	readingDone := make(chan struct{})
	writingDone := make(chan struct{})

	// 第一个值
	data.Write(1)

	// 启动一个写操作
	go func() {
		// 先通知我们即将开始写入
		close(writingStarted)

		// 执行写操作，在写锁保护下等待足够长时间
		data.WriteWithCallback(func(d *Data) {
			// 在写锁中暂停一段时间，给读取goroutines足够的时间尝试读取
			time.Sleep(200 * time.Millisecond)
			d.value = 2
		})

		// 写入完成
		close(writingDone)
	}()

	// 等待写操作开始
	<-writingStarted
	time.Sleep(10 * time.Millisecond) // 确保写锁已经获取

	// 启动多个读操作并尝试在写操作期间读取
	var wg sync.WaitGroup
	const readers = 10
	wg.Add(readers)

	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()

			// 持续尝试读取，直到写操作完成
			for {
				select {
				case <-writingDone:
					// 写操作完成，退出循环
					return
				default:
					// 尝试读取，用TryRead避免阻塞
					if val, ok := data.TryRead(); ok {
						readsDuringWrite.Add(1)
						if val == 2 {
							// 已经读取到新值，说明写操作已完成
							return
						}
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}

	// 等待所有读操作完成
	wg.Wait()
	close(readingDone)

	// 确保读操作完成
	<-readingDone

	if reads := readsDuringWrite.Load(); reads > 0 {
		t.Errorf("不应该有读操作在写操作期间执行，但检测到%d次读操作", reads)
	}

	if val := data.Read(); val != 2 {
		t.Errorf("写入后的值应为2，但得到: %v", val)
	}
}

// 模拟读写锁用于测试
type MockRWLocker struct {
	readLockCalled    bool
	readUnlockCalled  bool
	writeLockCalled   bool
	writeUnlockCalled bool
	shouldFailTryLock bool
}

func (m *MockRWLocker) ReadLock()                                  { m.readLockCalled = true }
func (m *MockRWLocker) ReadUnlock()                                { m.readUnlockCalled = true }
func (m *MockRWLocker) WriteLock()                                 { m.writeLockCalled = true }
func (m *MockRWLocker) WriteUnlock()                               { m.writeUnlockCalled = true }
func (m *MockRWLocker) TryReadLock() bool                          { return !m.shouldFailTryLock }
func (m *MockRWLocker) TryWriteLock() bool                         { return !m.shouldFailTryLock }
func (m *MockRWLocker) TryReadLockWithTimeout(time.Duration) bool  { return !m.shouldFailTryLock }
func (m *MockRWLocker) TryWriteLockWithTimeout(time.Duration) bool { return !m.shouldFailTryLock }

// 测试依赖注入功能
func TestDependencyInjection(t *testing.T) {
	mock := &MockRWLocker{}
	data := NewDataWithLocker(mock)

	data.Read()
	if !mock.readLockCalled {
		t.Error("读锁应该被调用")
	}
	if !mock.readUnlockCalled {
		t.Error("读锁解除应该被调用")
	}

	mock.readLockCalled = false
	mock.readUnlockCalled = false

	data.Write(123)
	if !mock.writeLockCalled {
		t.Error("写锁应该被调用")
	}
	if !mock.writeUnlockCalled {
		t.Error("写锁解除应该被调用")
	}
}

// 测试锁定失败的情况
func TestLockFailures(t *testing.T) {
	mock := &MockRWLocker{shouldFailTryLock: true}
	data := NewDataWithLocker(mock)

	// TryRead 应该失败
	_, ok := data.TryRead()
	if ok {
		t.Error("TryRead应该失败，但成功了")
	}

	// TryWrite 应该失败
	ok = data.TryWrite(100)
	if ok {
		t.Error("TryWrite应该失败，但成功了")
	}

	// ReadWithTimeout 应该失败
	_, ok = data.ReadWithTimeout(10 * time.Millisecond)
	if ok {
		t.Error("ReadWithTimeout应该失败，但成功了")
	}

	// WriteWithTimeout 应该失败
	ok = data.WriteWithTimeout(100, 10*time.Millisecond)
	if ok {
		t.Error("WriteWithTimeout应该失败，但成功了")
	}
}

// 模拟复杂场景：读多写少的数据缓存
func TestReadHeavyCache(t *testing.T) {
	data := NewData()
	const iterations = 1000
	const readers = 10
	const writers = 2

	var wg sync.WaitGroup
	wg.Add(readers + writers)

	// 启动多个读取者
	for i := 0; i < readers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				val := data.Read()
				if val < 0 {
					t.Errorf("读取到了无效值: %v", val)
				}
				// 偶尔使用超时读取
				if j%10 == 0 {
					if val, ok := data.ReadWithTimeout(5 * time.Millisecond); ok {
						if val < 0 {
							t.Errorf("超时读取到了无效值: %v", val)
						}
					}
				}
			}
		}(i)
	}

	// 启动几个写入者
	for i := 0; i < writers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations/10; j++ { // 写操作比读操作少
				success := data.Write(j + 1)
				if !success {
					t.Error("写入应该成功，但失败了")
				}
				time.Sleep(1 * time.Millisecond) // 稍微降低写入频率
			}
		}(i)
	}

	wg.Wait()
}

// 基准测试 - 读操作
func BenchmarkRead(b *testing.B) {
	data := NewData()
	data.Write(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Read()
	}
}

// 基准测试 - 写操作
func BenchmarkWrite(b *testing.B) {
	data := NewData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Write(i)
	}
}

// 基准测试 - 并发读操作
func BenchmarkConcurrentRead(b *testing.B) {
	data := NewData()
	data.Write(42)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data.Read()
		}
	})
}

// 基准测试 - 并发写操作
func BenchmarkConcurrentWrite(b *testing.B) {
	data := NewData()
	counter := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data.Write(counter)
			counter++
		}
	})
}

// 基准测试 - 读写混合操作
func BenchmarkMixedReadWrite(b *testing.B) {
	data := NewData()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			if counter%10 == 0 { // 10%的写操作
				data.Write(counter)
			} else {
				data.Read()
			}
			counter++
		}
	})
}

// 示例 - 演示读写锁的基本用法
func ExampleData_basic() {
	// 创建一个新的数据实例
	data := NewData()

	// 写入数据
	data.Write(100)

	// 读取数据
	value := data.Read()
	fmt.Printf("读取的值: %d\n", value)

	// 输出:
	// 读取的值: 100
}

// 示例 - 使用回调函数处理复杂操作
func ExampleData_callbacks() {
	data := NewData()
	data.Write(10)

	// 在读锁保护下执行操作
	data.ReadWithCallback(func(val int) {
		fmt.Printf("当前值: %d\n", val)
	})

	// 在写锁保护下执行多个操作
	data.WriteWithCallback(func(d *Data) {
		d.value = d.value * 2
		// 可以执行其他操作...
	})

	// 读取修改后的值
	fmt.Printf("修改后的值: %d\n", data.Read())

	// 输出:
	// 当前值: 10
	// 修改后的值: 20
}

// 示例 - 读取-计算-写入模式
func ExampleData_readModifyWrite() {
	data := NewData()
	data.Write(5)

	// 执行读取-计算-写入操作
	data.ReadWriteWithCallback(func(val int) int {
		// 根据读取的值计算新值
		return val * val // 平方
	})

	fmt.Printf("平方后的值: %d\n", data.Read())

	// 输出:
	// 平方后的值: 25
}
