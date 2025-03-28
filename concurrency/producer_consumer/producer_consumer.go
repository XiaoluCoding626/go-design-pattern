package producer_consumer

import (
	"fmt"
	"sync"
)

// ProducerConsumer 实现了生产者 - 消费者模式
type ProducerConsumer struct {
	queue chan int
	wg    sync.WaitGroup
}

// NewProducerConsumer 使用指定的缓冲区大小创建一个新实例
func NewProducerConsumer(bufferSize int) *ProducerConsumer {
	return &ProducerConsumer{
		queue: make(chan int, bufferSize),
	}
}

// Produce 生成数据并将其发送到队列中
func (pc *ProducerConsumer) Produce(count int) {
	pc.wg.Add(1)
	go func() {
		defer pc.wg.Done()
		defer close(pc.queue) // 生产完成时关闭队列

		for i := 0; i < count; i++ {
			fmt.Println("正在生产数据:", i)
			pc.queue <- i
		}
	}()
}

// Consume 处理队列中的数据
func (pc *ProducerConsumer) Consume() {
	pc.wg.Add(1)
	go func() {
		defer pc.wg.Done()
		for data := range pc.queue {
			fmt.Println("正在消费数据:", data)
		}
	}()
}

// Wait 阻塞直到所有生产和消费操作完成
func (pc *ProducerConsumer) Wait() {
	pc.wg.Wait()
}
