package singleton

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// 并发测试的协程数量
	concurrentCount = 500
)

// TestGetLazy 验证懒汉式单例的基本行为
// 测试目标: 确保多次调用 GetLazy() 返回的是同一个实例
func TestGetLazy(t *testing.T) {
	assert := assert.New(t)

	// 首次获取单例实例
	instance1 := GetLazy()
	// 再次获取单例实例
	instance2 := GetLazy()

	// 验证两个实例是同一个对象的引用
	assert.Same(instance1, instance2, "GetLazy() 应该始终返回相同的实例")

	// 验证单例实例可以正常工作
	// 这里只是调用方法，没有实际的断言，属于烟雾测试(smoke test)
	instance1.HelloWorld()
}

// TestConcurrentGetLazy 验证在并发环境下懒汉式单例的线程安全性
// 测试目标: 确保在多协程同时请求单例时，所有协程获取的都是同一个实例
func TestConcurrentGetLazy(t *testing.T) {
	assert := assert.New(t)

	// 用于等待所有协程完成
	wg := sync.WaitGroup{}
	wg.Add(concurrentCount)

	// 用于存储所有协程获取的实例
	instances := [concurrentCount]*Lazy{}

	// 启动多个协程同时获取单例
	for i := 0; i < concurrentCount; i++ {
		go func(index int) {
			defer wg.Done()
			instances[index] = GetLazy()
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()

	// 验证所有协程获取的都是同一个实例
	firstInstance := instances[0]
	for i := 1; i < concurrentCount; i++ {
		assert.Same(firstInstance, instances[i],
			"在并发环境下，GetLazy() 应该返回相同的实例")
	}
}
