// Package semaphore 实现了信号量设计模式，提供资源访问控制机制。
package semaphore

import (
	"context"
	"errors"
	"sync"
	"time"
)

// 定义错误常量
var (
	// ErrNoTickets 表示在超时时间内无法获取票证
	ErrNoTickets = errors.New("无法获取信号量票证")

	// ErrIllegalRelease 表示在没有票证的情况下尝试释放票证
	ErrIllegalRelease = errors.New("非法释放信号量票证")
)

// Semaphorer 定义了信号量应该具有的行为
type Semaphorer interface {
	// Acquire 尝试获取一个票证，可能会阻塞直到有票证可用或超时
	Acquire(ctx context.Context) error

	// TryAcquire 尝试非阻塞地获取一个票证，立即返回结果
	TryAcquire() bool

	// AcquireMany 尝试获取多个票证
	AcquireMany(n int, ctx context.Context) error

	// Release 释放一个已获取的票证
	Release() error

	// ReleaseMany 释放多个已获取的票证
	ReleaseMany(n int) error

	// Available 返回当前可用的票证数量
	Available() int

	// Size 返回信号量的总容量
	Size() int
}

// Semaphore 实现了信号量设计模式
type Semaphore struct {
	// 信号量的通道实现，空结构体是为了节省内存
	tickets chan struct{}

	// 信号量的最大容量
	size int

	// 用于保护计数器的互斥锁
	mu sync.Mutex

	// 已获取的票证数量
	acquired int
}

// New 创建一个新的信号量，指定票证总数
func New(size int) *Semaphore {
	if size <= 0 {
		size = 1 // 确保至少有一个票证
	}

	s := &Semaphore{
		tickets: make(chan struct{}, size),
		size:    size,
	}
	s.initialize() // 初始化填充通道
	return s
}

// initialize 确保信号量通道被填充到容量
func (s *Semaphore) initialize() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 清空通道，以防重新初始化
	for len(s.tickets) > 0 {
		<-s.tickets
	}

	// 填充通道到最大容量
	for i := 0; i < s.size; i++ {
		s.tickets <- struct{}{}
	}

	// 重置已获取的票证数量
	s.acquired = 0
}

// Acquire 尝试获取一个票证，如果无法立即获取，则阻塞等待
// 如果提供的context被取消，则返回context的错误
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case <-s.tickets:
		s.mu.Lock()
		s.acquired++
		s.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire 尝试非阻塞地获取一个票证，立即返回结果
func (s *Semaphore) TryAcquire() bool {
	select {
	case <-s.tickets:
		s.mu.Lock()
		s.acquired++
		s.mu.Unlock()
		return true
	default:
		return false
	}
}

// AcquireWithTimeout 尝试在指定超时时间内获取一个票证
func (s *Semaphore) AcquireWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.Acquire(ctx)
}

// AcquireMany 尝试获取多个票证
func (s *Semaphore) AcquireMany(n int, ctx context.Context) error {
	if n <= 0 {
		return nil
	}

	// 检查是否有足够的票证可用（非阻塞检查）
	s.mu.Lock()
	if len(s.tickets) < n {
		s.mu.Unlock()
		return ErrNoTickets
	}
	s.mu.Unlock()

	// 用于跟踪已获取的票证
	acquired := 0

	// 创建一个新的上下文，用于在失败时取消
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建错误通道
	errCh := make(chan error, 1)

	// 尝试获取票证
	go func() {
		for i := 0; i < n; i++ {
			select {
			case <-s.tickets:
				s.mu.Lock()
				s.acquired++
				acquired++
				s.mu.Unlock()
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
		errCh <- nil
	}()

	// 等待获取完成或出错
	err := <-errCh
	if err != nil {
		// 如果出错，释放已获取的票证
		for i := 0; i < acquired; i++ {
			s.Release()
		}
		return err
	}

	return nil
}

// Release 释放一个已获取的票证
func (s *Semaphore) Release() error {
	s.mu.Lock()
	if s.acquired <= 0 {
		s.mu.Unlock()
		return ErrIllegalRelease
	}
	s.acquired--
	s.mu.Unlock()

	// 归还票证到池中
	select {
	case s.tickets <- struct{}{}:
		return nil
	default:
		// 通道已满，这不应该发生
		s.mu.Lock()
		s.acquired++ // 回滚计数器
		s.mu.Unlock()
		return errors.New("信号量内部错误：通道已满")
	}
}

// ReleaseMany 释放多个已获取的票证
func (s *Semaphore) ReleaseMany(n int) error {
	if n <= 0 {
		return nil
	}

	s.mu.Lock()
	if s.acquired < n {
		s.mu.Unlock()
		return ErrIllegalRelease
	}
	s.acquired -= n
	s.mu.Unlock()

	// 归还所有票证
	for i := 0; i < n; i++ {
		select {
		case s.tickets <- struct{}{}:
			// 成功归还
		default:
			// 通道已满，这不应该发生
			s.mu.Lock()
			s.acquired++ // 回滚计数器
			s.mu.Unlock()
			return errors.New("信号量内部错误：通道已满")
		}
	}
	return nil
}

// Available 返回当前可用的票证数量
func (s *Semaphore) Available() int {
	return len(s.tickets)
}

// Size 返回信号量的总容量
func (s *Semaphore) Size() int {
	return s.size
}

// WaitAll 等待信号量恢复到完全可用状态
func (s *Semaphore) WaitAll(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		s.mu.Lock()
		acquired := s.acquired
		s.mu.Unlock()

		if acquired == 0 {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// 继续检查
		}
	}
}

// WeightedSemaphore 实现了带权重的信号量
type WeightedSemaphore struct {
	// 总容量
	capacity int64

	// 当前使用量
	used int64

	// 保护并发访问
	mu sync.Mutex

	// 当资源被释放时通知等待者
	cond *sync.Cond
}

// NewWeighted 创建一个新的带权重的信号量
func NewWeighted(capacity int64) *WeightedSemaphore {
	ws := &WeightedSemaphore{
		capacity: capacity,
	}
	ws.cond = sync.NewCond(&ws.mu)
	return ws
}

// Acquire 尝试获取指定权重的资源
func (ws *WeightedSemaphore) Acquire(ctx context.Context, weight int64) error {
	if weight <= 0 {
		return nil
	}

	if weight > ws.capacity {
		return errors.New("请求的权重超过总容量")
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	// 如果没有足够的资源，等待
	for ws.used+weight > ws.capacity {
		// 创建一个通道用于监听context取消
		done := make(chan struct{})

		// 在goroutine中监听context取消
		go func() {
			select {
			case <-ctx.Done():
				ws.cond.Broadcast() // 唤醒所有等待者
				close(done)
			case <-done:
				// 被正常唤醒或函数返回
			}
		}()

		// 等待资源可用
		ws.cond.Wait()

		// 关闭done通道，避免goroutine泄漏
		select {
		case <-done:
			// 已经关闭
		default:
			close(done)
		}

		// 检查context是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 继续尝试获取
		}
	}

	// 获取资源
	ws.used += weight
	return nil
}

// TryAcquire 尝试非阻塞地获取指定权重的资源
func (ws *WeightedSemaphore) TryAcquire(weight int64) bool {
	if weight <= 0 {
		return true
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.used+weight <= ws.capacity {
		ws.used += weight
		return true
	}
	return false
}

// Release 释放指定权重的资源
func (ws *WeightedSemaphore) Release(weight int64) {
	if weight <= 0 {
		return
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	if weight > ws.used {
		weight = ws.used // 不能释放超过已使用的资源
	}

	ws.used -= weight
	ws.cond.Broadcast() // 通知等待的goroutines
}

// Available 返回当前可用的资源量
func (ws *WeightedSemaphore) Available() int64 {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.capacity - ws.used
}

// Size 返回总容量
func (ws *WeightedSemaphore) Size() int64 {
	return ws.capacity
}
