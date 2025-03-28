package monitor

import (
	"sync"
	"testing"
	"time"
)

// TestBasicFunctionality tests the core operations of the Monitor
func TestBasicFunctionality(t *testing.T) {
	m := NewMonitor()

	// Test initial state
	val, valid := m.GetValue()
	if val != 0 || valid {
		t.Errorf("初始状态应为 (0, false)，但得到 (%d, %t)", val, valid)
	}

	// Test Produce
	m.Produce(42)
	val, valid = m.GetValue()
	if val != 42 || !valid {
		t.Errorf("调用 Produce(42) 后，GetValue() 应返回 (42, true)，但得到 (%d, %t)", val, valid)
	}

	// Test Consume
	consumedVal := m.Consume()
	if consumedVal != 42 {
		t.Errorf("应消费值 42，但得到 %d", consumedVal)
	}

	// Test state after consumption
	_, valid = m.GetValue()
	if valid {
		t.Error("消费后 valid 应为 false")
	}
}

// TestConsumerWaitsForProducer verifies that Consume blocks until data is available
func TestConsumerWaitsForProducer(t *testing.T) {
	m := NewMonitor()

	// Channel to collect results
	resultChan := make(chan int)

	// Start consumer in goroutine
	go func() {
		resultChan <- m.Consume()
	}()

	// Short delay to ensure consumer is waiting
	time.Sleep(50 * time.Millisecond)

	// Produce data
	expectedVal := 99
	m.Produce(expectedVal)

	// Check result with timeout
	select {
	case result := <-resultChan:
		if result != expectedVal {
			t.Errorf("消费者应接收 %d，但得到 %d", expectedVal, result)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("消费者未在预期时间内接收到数据")
	}
}

// TestReset verifies the Reset functionality
func TestReset(t *testing.T) {
	m := NewMonitor()

	// Setup state
	m.Produce(42)

	// Reset
	m.Reset()

	// Verify state was reset
	val, valid := m.GetValue()
	if val != 0 || valid {
		t.Errorf("重置后应为 (0, false)，但得到 (%d, %t)", val, valid)
	}
}

// TestExecuteProtected verifies protected execution
func TestExecuteProtected(t *testing.T) {
	m := NewMonitor()
	executed := false

	m.ExecuteProtected(func() {
		executed = true
		m.data = 77
		m.valid = true
	})

	if !executed {
		t.Error("保护执行函数未被调用")
	}

	// Verify state changes
	val, valid := m.GetValue()
	if val != 77 || !valid {
		t.Errorf("保护执行后状态应为 (77, true)，但得到 (%d, %t)", val, valid)
	}
}

// TestMultipleConsumers tests that multiple consumers can receive values
func TestMultipleConsumers(t *testing.T) {
	m := NewMonitor()
	const numConsumers = 3

	var wg sync.WaitGroup
	wg.Add(numConsumers)

	// Track consumed values
	consumedValues := make([]int, 0, numConsumers)
	var mu sync.Mutex

	// Start consumers
	for i := 0; i < numConsumers; i++ {
		go func() {
			defer wg.Done()
			val := m.Consume()

			mu.Lock()
			consumedValues = append(consumedValues, val)
			mu.Unlock()
		}()
	}

	// Give consumers time to start waiting
	time.Sleep(50 * time.Millisecond)

	// Produce values
	expectedVal := 42
	for i := 0; i < numConsumers; i++ {
		m.Produce(expectedVal)
		// Small delay between productions
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for all consumers
	wg.Wait()

	// Verify results
	if len(consumedValues) != numConsumers {
		t.Errorf("应有 %d 个消费者接收到值，但有 %d 个", numConsumers, len(consumedValues))
	}

	for i, val := range consumedValues {
		if val != expectedVal {
			t.Errorf("消费者 %d 应接收值 %d，但接收到 %d", i, expectedVal, val)
		}
	}
}

// TestProducerConsumerSequence tests a sequence of produce/consume operations
func TestProducerConsumerSequence(t *testing.T) {
	m := NewMonitor()

	// Series of values to test
	testValues := []int{1, 2, 3, 4, 5}

	for _, expected := range testValues {
		// Produce value
		m.Produce(expected)

		// Consume value
		actual := m.Consume()

		if actual != expected {
			t.Errorf("应消费 %d，但得到 %d", expected, actual)
		}
	}
}

// TestConcurrentOperations tests monitor under concurrent load
func TestConcurrentOperations(t *testing.T) {
	m := NewMonitor()
	const operations = 50 // 减少操作数量，避免过载

	var wg sync.WaitGroup
	wg.Add(operations)

	// 跟踪已消费的值
	consumed := make(chan int, operations)

	// 启动消费者
	for i := 0; i < operations; i++ {
		go func() {
			defer wg.Done()
			val := m.Consume()
			consumed <- val
		}()
	}

	// 确保消费者开始等待
	time.Sleep(100 * time.Millisecond)

	// 生产值
	for i := 1; i <= operations; i++ {
		m.Produce(i)
		// 生产之间添加小延迟，确保信号能够正确传递
		time.Sleep(2 * time.Millisecond)
	}

	// 等待所有消费者完成
	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		// 测试通过
		close(consumed)
		count := 0
		for range consumed {
			count++
		}
		if count != operations {
			t.Errorf("期望消费 %d 个值，但得到 %d 个", operations, count)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("并发测试超时")
	}
}
