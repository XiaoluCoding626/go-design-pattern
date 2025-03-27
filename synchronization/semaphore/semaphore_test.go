package semaphore

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试创建信号量
func TestNew(t *testing.T) {
	// 测试普通信号量
	s := New(5)
	assert.NotNil(t, s, "信号量不应为nil")
	assert.Equal(t, 5, s.Size(), "信号量大小应为5")

	// 测试权重信号量
	ws := NewWeighted(10)
	assert.NotNil(t, ws, "权重信号量不应为nil")
	assert.Equal(t, int64(10), ws.Size(), "权重信号量大小应为10")

	// 测试负值处理 - 应该至少创建容量为1的信号量
	s = New(-1)
	assert.Equal(t, 1, s.Size(), "负值容量应被修正为1")
}

// 测试基本的获取和释放操作
func TestBasicAcquireRelease(t *testing.T) {
	s := New(3)

	// 初始化信号量
	s.initialize()
	assert.Equal(t, 3, s.Available(), "初始可用票证应为3")

	// 获取一个票证
	ctx := context.Background()
	err := s.Acquire(ctx)
	assert.NoError(t, err, "获取票证不应有错误")
	assert.Equal(t, 2, s.Available(), "获取后可用票证应为2")

	// 再获取一个
	err = s.Acquire(ctx)
	assert.NoError(t, err, "获取第二个票证不应有错误")
	assert.Equal(t, 1, s.Available(), "获取后可用票证应为1")

	// 释放一个
	err = s.Release()
	assert.NoError(t, err, "释放票证不应有错误")
	assert.Equal(t, 2, s.Available(), "释放后可用票证应为2")
}

// 测试非阻塞获取
func TestTryAcquire(t *testing.T) {
	s := New(2)
	s.initialize()

	// 第一个票证应该可以获取
	success := s.TryAcquire()
	assert.True(t, success, "首次尝试获取应成功")

	// 第二个也应该可以获取
	success = s.TryAcquire()
	assert.True(t, success, "第二次尝试获取应成功")

	// 第三个应该失败
	success = s.TryAcquire()
	assert.False(t, success, "第三次尝试获取应失败")

	// 释放一个后再尝试
	err := s.Release()
	assert.NoError(t, err)
	success = s.TryAcquire()
	assert.True(t, success, "释放后尝试获取应成功")
}

// 测试带超时的获取
func TestAcquireWithTimeout(t *testing.T) {
	s := New(1)
	s.initialize()

	// 先获取唯一的票证
	ctx := context.Background()
	err := s.Acquire(ctx)
	assert.NoError(t, err)

	// 尝试带超时获取，应该超时
	err = s.AcquireWithTimeout(50 * time.Millisecond)
	assert.Error(t, err, "应该因超时而失败")
	assert.Contains(t, err.Error(), "deadline", "应该是超时错误")

	// 释放后再尝试获取，应该成功
	err = s.Release()
	assert.NoError(t, err)
	err = s.AcquireWithTimeout(50 * time.Millisecond)
	assert.NoError(t, err, "释放后获取应成功")
}

// 测试批量获取和释放
func TestBulkOperations(t *testing.T) {
	s := New(5)
	s.initialize()

	ctx := context.Background()

	// 批量获取3个票证
	err := s.AcquireMany(3, ctx)
	assert.NoError(t, err, "批量获取应成功")
	assert.Equal(t, 2, s.Available(), "剩余票证应为2")

	// 尝试批量获取4个(超过可用数)，应该失败
	err = s.AcquireMany(4, ctx)
	assert.Error(t, err, "超量获取应失败")
	assert.Equal(t, 2, s.Available(), "失败后可用票证仍为2")

	// 批量释放2个
	err = s.ReleaseMany(2)
	assert.NoError(t, err, "批量释放应成功")
	assert.Equal(t, 4, s.Available(), "释放后可用票证应为4")

	// 尝试释放超过已获取的数量
	err = s.ReleaseMany(2)
	assert.Error(t, err, "超量释放应失败")
	assert.Equal(t, ErrIllegalRelease, err, "应返回非法释放错误")
}

// 测试等待所有票证返回
func TestWaitAll(t *testing.T) {
	s := New(3)
	s.initialize()

	ctx := context.Background()

	// 获取所有票证
	_ = s.AcquireMany(3, ctx)
	assert.Equal(t, 0, s.Available(), "所有票证都被获取")

	// 在另一个goroutine中释放票证
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = s.ReleaseMany(3)
	}()

	// 等待所有票证返回
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err := s.WaitAll(ctx)
	assert.NoError(t, err, "应成功等待所有票证")
	assert.Equal(t, 3, s.Available(), "所有票证应该都已返回")
}

// 测试Context取消
func TestContextCancellation(t *testing.T) {
	s := New(1)
	s.initialize()

	// 获取唯一的票证
	_ = s.Acquire(context.Background())

	// 创建一个将取消的context
	ctx, cancel := context.WithCancel(context.Background())

	// 在短时间后取消context
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// 尝试获取票证，应该因context取消而失败
	err := s.Acquire(ctx)
	assert.Error(t, err, "获取应因context取消而失败")
	assert.Equal(t, context.Canceled, err, "应返回context取消错误")
}

// 测试并发获取和释放
func TestConcurrentOperations(t *testing.T) {
	s := New(50)
	s.initialize()

	const goroutines = 100
	const opsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// 启动多个goroutine并发获取和释放票证
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for j := 0; j < opsPerGoroutine; j++ {
				// 获取票证
				err := s.Acquire(ctx)
				assert.NoError(t, err)

				// 模拟一些工作
				time.Sleep(1 * time.Millisecond)

				// 释放票证
				err = s.Release()
				assert.NoError(t, err)
			}
		}()
	}

	wg.Wait()

	// 所有操作完成后，所有票证应该都已返回
	assert.Equal(t, 50, s.Available(), "所有票证应该返回")
}

// 测试错误处理 - 非法释放
func TestIllegalRelease(t *testing.T) {
	s := New(3)
	s.initialize()

	// 不获取直接释放，应该失败
	err := s.Release()
	assert.Error(t, err, "未获取就释放应失败")
	assert.Equal(t, ErrIllegalRelease, err, "应返回非法释放错误")

	// 获取1个，尝试释放2个
	ctx := context.Background()
	_ = s.Acquire(ctx)
	err = s.ReleaseMany(2)
	assert.Error(t, err, "释放超过获取数量应失败")
	assert.Equal(t, ErrIllegalRelease, err, "应返回非法释放错误")
}

// 测试带权重的信号量
func TestWeightedSemaphore(t *testing.T) {
	ws := NewWeighted(10)
	ctx := context.Background()

	// 获取权重为3的资源
	err := ws.Acquire(ctx, 3)
	assert.NoError(t, err, "获取权重为3的资源应成功")
	assert.Equal(t, int64(7), ws.Available(), "剩余可用资源应为7")

	// 再获取权重为4的资源
	err = ws.Acquire(ctx, 4)
	assert.NoError(t, err, "获取权重为4的资源应成功")
	assert.Equal(t, int64(3), ws.Available(), "剩余可用资源应为3")

	// 尝试获取权重为5的资源，超过可用资源
	success := ws.TryAcquire(5)
	assert.False(t, success, "获取超过可用资源应失败")

	// 尝试获取权重为11的资源，超过总容量
	err = ws.Acquire(ctx, 11)
	assert.Error(t, err, "获取超过总容量应失败")

	// 释放权重为3的资源
	ws.Release(3)
	assert.Equal(t, int64(6), ws.Available(), "释放后可用资源应为6")

	// 现在可以获取权重为5的资源了
	success = ws.TryAcquire(5)
	assert.True(t, success, "释放后获取应成功")
	assert.Equal(t, int64(1), ws.Available(), "剩余可用资源应为1")
}

// 测试权重信号量的并发操作
func TestConcurrentWeightedOperations(t *testing.T) {
	ws := NewWeighted(100)

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			weight := int64(id%5 + 1) // 使用1-5的权重

			ctx := context.Background()
			err := ws.Acquire(ctx, weight)
			assert.NoError(t, err)

			// 模拟一些工作
			time.Sleep(10 * time.Millisecond)

			// 释放资源
			ws.Release(weight)
		}(i)
	}

	wg.Wait()

	// 所有操作完成后，所有资源应该都已返回
	assert.Equal(t, int64(100), ws.Available(), "所有资源应该返回")
}

// 测试带超时的权重信号量操作
func TestWeightedSemaphoreTimeout(t *testing.T) {
	ws := NewWeighted(10)

	// 获取大部分资源
	err := ws.Acquire(context.Background(), 8)
	assert.NoError(t, err)

	// 尝试在有限时间内获取超过剩余资源的权重
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = ws.Acquire(ctx, 5)
	assert.Error(t, err, "超时获取应失败")
	assert.Contains(t, err.Error(), "deadline", "应是超时错误")

	// 释放一些资源
	ws.Release(4)

	// 现在应该可以获取权重为3的资源
	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = ws.Acquire(ctx, 3)
	assert.NoError(t, err, "释放后获取应成功")
}

// 测试信号量在零或负面大小的情况
func TestZeroOrNegativeSize(t *testing.T) {
	// 零大小 - 应该创建容量为1的信号量
	s := New(0)
	assert.Equal(t, 1, s.Size(), "零大小应被修正为1")

	// 负大小 - 也应该创建容量为1的信号量
	s = New(-5)
	assert.Equal(t, 1, s.Size(), "负大小应被修正为1")

	// 权重信号量的零大小
	ws := NewWeighted(0)
	assert.Equal(t, int64(0), ws.Size(), "权重信号量允许零大小")

	// 尝试获取权重为1的资源，应该失败
	success := ws.TryAcquire(1)
	assert.False(t, success, "零容量信号量不应允许获取")
}

// 测试在退出前释放所有资源
func TestReleaseBeforeExit(t *testing.T) {
	s := New(5)
	s.initialize()
	ctx := context.Background()

	// 获取3个票证
	_ = s.AcquireMany(3, ctx)

	// 使用defer确保在函数退出前释放所有资源
	defer func() {
		_ = s.ReleaseMany(3)
	}()

	// 模拟一些操作
	time.Sleep(10 * time.Millisecond)

	// 函数结束前，延迟函数会释放资源
}

// 作为一种并发安全的计数器使用
func TestSemaphoreAsCounter(t *testing.T) {
	// 创建一个大容量的信号量作为计数器
	counter := New(1000)
	counter.initialize()

	const operations = 100
	var wg sync.WaitGroup
	wg.Add(operations)

	// 并发增加计数
	for i := 0; i < operations; i++ {
		go func() {
			defer wg.Done()
			_ = counter.Acquire(context.Background())
		}()
	}

	wg.Wait()

	// 验证计数
	assert.Equal(t, 900, counter.Available(), "计数器应记录100个操作")
}
