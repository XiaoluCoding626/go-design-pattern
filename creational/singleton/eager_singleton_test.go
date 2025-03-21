package singleton

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 初始化单例实例
func TestMain(m *testing.M) {
	// 程序运行时初始化，设置初始计数为3
	InitEager(3)
	os.Exit(m.Run())
}

// 测试GetEager返回的是否为同一个单例实例
func TestGetEager(t *testing.T) {
	assert := assert.New(t)

	// 获取两次实例，应该是同一个对象
	instance1 := GetEager()
	instance2 := GetEager()

	// 确保返回的是同一个实例（指针相等）
	assert.Same(instance1, instance2, "应该返回相同的单例实例")

	// 验证初始化值正确
	assert.Equal(3, instance1.GetCount(), "初始计数值应为3")
}

// 测试对实例的修改会影响所有获取的引用
func TestEagerModification(t *testing.T) {
	assert := assert.New(t)

	// 获取两个引用
	instance1 := GetEager()
	instance2 := GetEager()

	// 通过第一个引用修改计数
	initialCount := instance1.GetCount()
	instance1.Increase()

	// 验证两个引用都能看到修改
	assert.Equal(initialCount+1, instance1.GetCount(), "实例1的计数应增加")
	assert.Equal(initialCount+1, instance2.GetCount(), "实例2应反映相同的计数变化")
}

// 测试并发环境下的单例行为
func TestEagerConcurrentAccess(t *testing.T) {
	assert := assert.New(t)

	// 重置单例以获得可预测的起点
	InitEager(0)

	// 创建多个goroutine同时访问和修改单例
	const goroutineCount = 100
	var wg sync.WaitGroup

	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()
			instance := GetEager()
			instance.Increase() // 增加计数
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 验证最终计数是否正确
	finalInstance := GetEager()
	assert.Equal(goroutineCount, finalInstance.GetCount(), "最终计数应等于goroutine数量")
}
