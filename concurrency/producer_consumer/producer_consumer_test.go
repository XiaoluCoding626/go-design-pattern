package producer_consumer

import (
	"fmt"
	"testing"
)

// TestProducerConsumer 测试生产者 - 消费者模式的基本功能
func TestProducerConsumer(t *testing.T) {
	pc := NewProducerConsumer(5)

	pc.Produce(5)
	pc.Consume()

	pc.Wait()
}

// TestMultipleConsumers 测试单个生产者与多个消费者的情况
func TestMultipleConsumers(t *testing.T) {
	pc := NewProducerConsumer(10)

	pc.Produce(10)

	// 启动多个消费者
	for i := 0; i < 3; i++ {
		pc.Consume()
	}

	pc.Wait()
}

// ExampleProducerConsumer 演示生产者 - 消费者模式
func ExampleProducerConsumer() {
	// 创建一个缓冲区大小为 5 的生产者 - 消费者实例
	pc := NewProducerConsumer(5)

	// 启动一个生产者，生成 5 个元素
	pc.Produce(5)

	// 启动一个消费者，处理这些元素
	pc.Consume()

	// 等待所有操作完成
	pc.Wait()

	// 注意：由于并发的原因，此示例的输出是不确定的，
	// 所以我们不提供预期输出的注释。
}

// DemoSequential 为文档目的顺序展示生产者 - 消费者的流程
func DemoSequential() {
	// 创建一个带缓冲的通道
	queue := make(chan int, 5)

	// 顺序生产元素
	for i := 0; i < 5; i++ {
		fmt.Println("正在生产数据:", i)
		queue <- i
	}
	close(queue)

	// 顺序消费元素
	for data := range queue {
		fmt.Println("正在消费数据:", data)
	}

	// 输出:
	// 正在生产数据: 0
	// 正在生产数据: 1
	// 正在生产数据: 2
	// 正在生产数据: 3
	// 正在生产数据: 4
	// 正在消费数据: 0
	// 正在消费数据: 1
	// 正在消费数据: 2
	// 正在消费数据: 3
	// 正在消费数据: 4
}
