package barrier

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Barrier 实现了障碍同步模式，允许一组goroutine在某个同步点等待，
// 直到所有goroutine都到达该点后才能继续执行
type Barrier struct {
	participants int        // 参与者总数
	count        int        // 当前已到达的参与者数
	generation   int        // 当前同步轮次
	mutex        sync.Mutex // 保护内部状态的互斥锁
	cond         *sync.Cond // 用于通知等待goroutine的条件变量
}

// NewBarrier 创建一个新的barrier，指定参与同步的goroutine数量
func NewBarrier(participants int) *Barrier {
	if participants <= 0 {
		panic("barrier参与者数量必须大于0")
	}
	b := &Barrier{
		participants: participants,
		count:        0,
		generation:   0,
	}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

// Wait 使当前goroutine在barrier处等待，直到所有参与者都到达
func (b *Barrier) Wait() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	generation := b.generation // 记录当前轮次

	b.count++
	if b.count == b.participants {
		// 最后一个到达的goroutine触发唤醒操作
		b.count = 0        // 重置计数器
		b.generation++     // 增加轮次
		b.cond.Broadcast() // 唤醒所有等待的goroutine
	} else {
		// 等待直到所有参与者都到达（通过轮次变化来检测）
		for generation == b.generation {
			b.cond.Wait()
		}
	}
}

// WaitWithTimeout 带超时版本的Wait
func (b *Barrier) WaitWithTimeout(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})

	go func() {
		b.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true // 成功
	case <-ctx.Done():
		return false // 超时
	}
}

// Reset 重置barrier状态，允许重新使用
func (b *Barrier) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.count = 0
	b.generation++
	b.cond.Broadcast() // 唤醒所有等待的goroutine
}

// Worker 演示使用barrier的工作函数
func Worker(id int, phase1, phase2 *Barrier) {
	fmt.Printf("工作者 %d 开始执行\n", id)

	// 模拟第一阶段工作
	for i := 0; i < 3; i++ {
		fmt.Printf("工作者 %d 正在执行第一阶段\n", id)
		time.Sleep(time.Millisecond * 100) // 模拟工作
	}

	fmt.Printf("工作者 %d 到达第一个同步点\n", id)
	phase1.Wait() // 在第一个barrier处同步
	fmt.Printf("工作者 %d 通过第一个同步点\n", id)

	// 模拟第二阶段工作
	for i := 0; i < 3; i++ {
		fmt.Printf("工作者 %d 正在执行第二阶段\n", id)
		time.Sleep(time.Millisecond * 100) // 模拟工作
	}

	fmt.Printf("工作者 %d 到达第二个同步点\n", id)
	phase2.Wait() // 在第二个barrier处同步
	fmt.Printf("工作者 %d 通过第二个同步点\n", id)

	fmt.Printf("工作者 %d 完成所有工作\n", id)
}

// RunExample 运行barrier模式示例
func RunExample() {
	numWorkers := 3
	phase1 := NewBarrier(numWorkers)
	phase2 := NewBarrier(numWorkers)

	// 启动工作者
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			defer wg.Done()
			Worker(id, phase1, phase2)
		}(i)
	}

	// 等待所有工作者完成
	wg.Wait()
	fmt.Println("所有工作者已完成工作")
}
