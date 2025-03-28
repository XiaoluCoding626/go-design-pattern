package bounded_parallelism

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Task 表示一个可执行的任务
type Task[T any] struct {
	ID       string            // 任务标识符
	Execute  func() (T, error) // 任务执行函数
	Priority int               // 任务优先级（可选）
	Timeout  time.Duration     // 任务超时时间（可选）
}

// Result 表示任务执行的结果
type Result[T any] struct {
	TaskID    string    // 对应的任务ID
	Value     T         // 任务执行的返回值
	Err       error     // 任务执行中遇到的错误
	StartTime time.Time // 任务开始执行的时间
	EndTime   time.Time // 任务完成的时间
}

// BoundedExecutor 实现有界并行性模式，限制并发执行的任务数量
type BoundedExecutor[T any] struct {
	semaphore chan struct{}      // 信号量，用于限制并发数
	tasks     chan Task[T]       // 任务队列
	results   chan Result[T]     // 结果通道
	wg        sync.WaitGroup     // 等待所有工作完成
	ctx       context.Context    // 用于取消操作的上下文
	cancel    context.CancelFunc // 取消函数
	closed    bool               // 是否已关闭
	mu        sync.Mutex         // 保护 closed 字段的互斥锁
}

// NewBoundedExecutor 创建一个新的有界执行器
func NewBoundedExecutor[T any](maxConcurrent int, queueSize int) *BoundedExecutor[T] {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	if queueSize < 0 {
		queueSize = 0
	}

	ctx, cancel := context.WithCancel(context.Background())
	executor := &BoundedExecutor[T]{
		semaphore: make(chan struct{}, maxConcurrent),
		tasks:     make(chan Task[T], queueSize),
		results:   make(chan Result[T], queueSize),
		ctx:       ctx,
		cancel:    cancel,
		closed:    false,
	}

	// 启动工作池
	executor.startWorkers(maxConcurrent)
	return executor
}

// startWorkers 启动工作协程池
func (e *BoundedExecutor[T]) startWorkers(count int) {
	for i := 0; i < count; i++ {
		e.wg.Add(1)
		go func(workerID int) {
			defer e.wg.Done()
			for {
				select {
				case task, ok := <-e.tasks:
					if !ok {
						return // 任务通道已关闭，退出
					}
					e.executeTask(workerID, task)
				case <-e.ctx.Done():
					return // 上下文被取消，退出
				}
			}
		}(i + 1)
	}
}

// executeTask 执行单个任务并处理结果
func (e *BoundedExecutor[T]) executeTask(workerID int, task Task[T]) {
	e.semaphore <- struct{}{}        // 获取信号量
	defer func() { <-e.semaphore }() // 释放信号量

	var result Result[T]
	result.TaskID = task.ID
	result.StartTime = time.Now()

	fmt.Printf("工作者 %d 开始执行任务: %s\n", workerID, task.ID)

	// 执行任务，支持超时控制
	if task.Timeout > 0 {
		taskCtx, cancel := context.WithTimeout(e.ctx, task.Timeout)
		defer cancel()

		// 在单独的goroutine中执行任务
		done := make(chan struct{})
		go func() {
			result.Value, result.Err = task.Execute()
			close(done)
		}()

		// 等待任务完成或超时
		select {
		case <-done:
			// 任务正常完成
		case <-taskCtx.Done():
			result.Err = errors.New("任务执行超时")
		}
	} else {
		// 无超时的任务直接执行
		result.Value, result.Err = task.Execute()
	}

	result.EndTime = time.Now()

	// 安全地发送结果，防止因通道关闭导致panic
	sendResult := func() (sent bool) {
		// 使用recover捕获向已关闭通道发送的异常
		defer func() {
			if r := recover(); r != nil {
				// 通道已关闭，忽略错误
				sent = false
			}
		}()

		// 使用非阻塞发送尝试提交结果
		select {
		case e.results <- result:
			return true
		case <-e.ctx.Done():
			return false
		default:
			// 队列已满或已关闭的情况下，尝试阻塞发送
			select {
			case e.results <- result:
				return true
			case <-e.ctx.Done():
				return false
			}
		}
	}

	// 尝试发送结果
	sent := sendResult()

	fmt.Printf("工作者 %d 完成任务: %s, 耗时: %v, 结果已发送: %v\n",
		workerID, task.ID, result.EndTime.Sub(result.StartTime), sent)
}

// Submit 提交一个任务到执行队列
func (e *BoundedExecutor[T]) Submit(task Task[T]) error {
	// 检查执行器是否已关闭
	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		return errors.New("执行器已关闭")
	}
	e.mu.Unlock()

	// 使用非阻塞发送尝试提交任务
	select {
	case e.tasks <- task:
		return nil
	case <-e.ctx.Done():
		return errors.New("执行器已关闭")
	default:
		// 队列已满的情况
		// 阻塞发送，但仍然可以被取消
		select {
		case e.tasks <- task:
			return nil
		case <-e.ctx.Done():
			return errors.New("执行器已关闭")
		}
	}
}

// Results 返回结果通道，用于获取任务执行结果
func (e *BoundedExecutor[T]) Results() <-chan Result[T] {
	return e.results
}

// Shutdown 优雅关闭执行器，等待所有进行中的任务完成
func (e *BoundedExecutor[T]) Shutdown() {
	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		return
	}
	e.closed = true
	e.mu.Unlock()

	close(e.tasks) // 不再接受新任务
	e.wg.Wait()    // 等待所有工作者完成
	close(e.results)
}

// ShutdownNow 立即关闭执行器，取消所有进行中的任务
func (e *BoundedExecutor[T]) ShutdownNow() {
	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		return
	}
	e.closed = true
	e.mu.Unlock()

	e.cancel() // 取消上下文

	// 安全地关闭任务通道
	select {
	case _, ok := <-e.tasks:
		if ok {
			close(e.tasks) // 如果通道还没关闭，则关闭它
		}
	default:
		close(e.tasks) // 如果通道为空，则关闭它
	}

	// 此处不等待工作者完成，直接关闭结果通道会导致正在执行的任务发生panic
	// 但我们已经在executeTask中处理了这种可能性，所以可以安全地关闭结果通道
	close(e.results)
}

// RunExample 运行有界并行模式的示例
func RunExample() {
	// 创建有界执行器，最多允许3个并发任务，队列大小为10
	executor := NewBoundedExecutor[string](3, 10)

	// 提交任务
	for i := 1; i <= 10; i++ {
		taskID := fmt.Sprintf("Task-%d", i)
		task := Task[string]{
			ID: taskID,
			Execute: func() (string, error) {
				// 模拟任务执行
				time.Sleep(2 * time.Second)
				return fmt.Sprintf("结果-%s", taskID), nil
			},
			Timeout: 5 * time.Second,
		}

		executor.Submit(task)
	}

	// 启动收集结果的协程
	var resultWg sync.WaitGroup
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for result := range executor.Results() {
			if result.Err != nil {
				fmt.Printf("任务 %s 执行失败: %v\n", result.TaskID, result.Err)
			} else {
				fmt.Printf("任务 %s 执行成功: %v\n", result.TaskID, result.Value)
			}
		}
	}()

	// 等待一段时间后优雅关闭执行器
	time.Sleep(12 * time.Second)
	executor.Shutdown()

	// 等待结果处理完成
	resultWg.Wait()
	fmt.Println("所有任务已完成")
}
