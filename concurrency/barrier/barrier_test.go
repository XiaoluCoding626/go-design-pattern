package barrier

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestBarrierBasic 测试 Barrier 的基本功能
func TestBarrierBasic(t *testing.T) {
	const numWorkers = 3
	b := NewBarrier(numWorkers)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// 用于检查所有 goroutine 是否已经通过 barrier
	passed := make(chan bool, numWorkers)
	allReady := make(chan struct{})

	// 启动多个 goroutine
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			defer wg.Done()

			// 等待信号开始
			<-allReady

			t.Logf("工作者 %d 到达屏障", id)
			b.Wait()
			t.Logf("工作者 %d 通过屏障", id)

			passed <- true
		}(i)
	}

	// 释放所有 goroutine
	close(allReady)

	// 等待所有 goroutine 完成
	wg.Wait()
	close(passed)

	// 检查是否所有 goroutine 都通过了 barrier
	count := 0
	for range passed {
		count++
	}

	if count != numWorkers {
		t.Errorf("期望所有 %d 个工作者通过屏障，但只有 %d 个通过", numWorkers, count)
	}
}

// TestMultiplePhases 测试 Barrier 在多个阶段中的使用
func TestMultiplePhases(t *testing.T) {
	const numWorkers = 5
	const phases = 3

	b := NewBarrier(numWorkers)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// 跟踪每个工作者的阶段完成情况
	phaseComplete := make([][]bool, numWorkers)
	for i := range phaseComplete {
		phaseComplete[i] = make([]bool, phases)
	}

	// 保护对 phaseComplete 的访问
	var mutex sync.Mutex

	// 启动工作者
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			defer wg.Done()

			for phase := 0; phase < phases; phase++ {
				// 模拟工作
				time.Sleep(10 * time.Millisecond)

				// 到达屏障
				b.Wait()

				// 标记此阶段完成
				mutex.Lock()
				phaseComplete[id][phase] = true
				mutex.Unlock()

				// 给其他 goroutine 时间标记他们的状态
				time.Sleep(5 * time.Millisecond)

				// 确认所有工作者都完成了这个阶段
				mutex.Lock()
				for w := 0; w < numWorkers; w++ {
					if !phaseComplete[w][phase] {
						t.Errorf("阶段 %d: 工作者 %d 未能完成，但工作者 %d 已经继续",
							phase, w, id)
					}
				}
				mutex.Unlock()
			}
		}(i)
	}

	wg.Wait()
	// 检查所有阶段是否都已完成
	for w := 0; w < numWorkers; w++ {
		for p := 0; p < phases; p++ {
			if !phaseComplete[w][p] {
				t.Errorf("工作者 %d 未能完成阶段 %d", w, p)
			}
		}
	}
}

// TestBarrierReset 测试 Barrier 的重置功能
func TestBarrierReset(t *testing.T) {
	b := NewBarrier(2)

	// 第一个 goroutine 到达并等待
	var wg sync.WaitGroup
	wg.Add(1)

	arrived := make(chan struct{})
	go func() {
		defer wg.Done()

		close(arrived)
		// 这个应该被阻塞，因为只有1个工作者到达
		b.Wait()
	}()

	// 确保第一个 goroutine 已经到达
	<-arrived
	time.Sleep(100 * time.Millisecond)

	// 重置 barrier
	b.Reset()

	// 等待第一个 goroutine 完成
	wg.Wait()
}

// TestBarrierTimeout 测试带超时的等待功能
func TestBarrierTimeout(t *testing.T) {
	b := NewBarrier(2)

	// 尝试等待，但设置很短的超时
	result := b.WaitWithTimeout(50 * time.Millisecond)

	if result {
		t.Error("预期因超时而失败，但等待成功")
	}
}

// TestBarrierPanic 测试创建无效 Barrier 时的 panic 情况
func TestBarrierPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("预期参与者数量为 0 时会 panic，但未发生")
		}
	}()

	// 这应该导致 panic
	_ = NewBarrier(0)
}

// TestRunExample 测试完整的示例函数
func TestRunExample(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过耗时的示例测试")
	}

	RunExample() // 验证它是否不会崩溃
}

// BenchmarkBarrier 基准测试 Barrier 的性能
func BenchmarkBarrier(b *testing.B) {
	const numWorkers = 10

	for i := 0; i < b.N; i++ {
		barrier := NewBarrier(numWorkers)
		var wg sync.WaitGroup
		wg.Add(numWorkers)

		for j := 0; j < numWorkers; j++ {
			go func() {
				defer wg.Done()
				barrier.Wait()
			}()
		}

		wg.Wait()
	}
}

// BenchmarkBarrierMultiplePhases 测试多阶段 Barrier 的性能
func BenchmarkBarrierMultiplePhases(b *testing.B) {
	const numWorkers = 10
	const phases = 5

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		barrier := NewBarrier(numWorkers)
		var wg sync.WaitGroup
		wg.Add(numWorkers)

		for j := 0; j < numWorkers; j++ {
			go func() {
				defer wg.Done()

				for phase := 0; phase < phases; phase++ {
					barrier.Wait()
				}
			}()
		}

		wg.Wait()
	}
}

// ExampleBarrier 提供一个简单的使用示例
func ExampleBarrier() {
	// 创建一个有 3 个参与者的屏障
	b := NewBarrier(3)

	// 启动 3 个 goroutine
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 1; i <= 3; i++ {
		go func(id int) {
			defer wg.Done()

			// 第一阶段工作
			fmt.Printf("工作者 %d 完成第一阶段\n", id)

			b.Wait() // 所有工作者在此同步

			// 第二阶段工作
			fmt.Printf("工作者 %d 完成第二阶段\n", id)
		}(i)
	}

	wg.Wait()
	// 输出顺序可能会不同，故不提供预期输出
}
