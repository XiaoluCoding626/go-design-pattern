package registry

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试服务结构体
type TestService struct {
	Name string
}

func (s *TestService) GetName() string {
	return s.Name
}

// 测试创建新的注册表实例
func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.services)
	assert.NotNil(t, registry.factories)
}

// 测试获取全局单例
func TestGetRegistry(t *testing.T) {
	registry1 := GetRegistry()
	registry2 := GetRegistry()

	assert.NotNil(t, registry1)
	assert.Same(t, registry1, registry2, "全局注册表应该是单例")
}

// 测试注册和获取服务
func TestRegisterAndGet(t *testing.T) {
	registry := NewRegistry()
	service := &TestService{Name: "TestService"}

	// 注册服务
	err := registry.Register("test", service)
	assert.NoError(t, err)

	// 获取服务
	result, err := registry.Get("test")
	assert.NoError(t, err)
	assert.Same(t, service, result)

	// 类型断言
	testService, ok := result.(*TestService)
	assert.True(t, ok)
	assert.Equal(t, "TestService", testService.Name)
}

// 测试注册工厂函数和延迟初始化
func TestRegisterFactory(t *testing.T) {
	registry := NewRegistry()
	initialized := false

	// 注册工厂函数
	err := registry.RegisterFactory("lazy", func() interface{} {
		initialized = true
		return &TestService{Name: "LazyService"}
	})

	assert.NoError(t, err)
	assert.False(t, initialized, "工厂函数不应该在注册时执行")

	// 第一次获取服务，应该初始化
	result, err := registry.Get("lazy")
	assert.NoError(t, err)
	assert.True(t, initialized, "工厂函数应该在第一次获取时执行")

	// 确认服务类型正确
	service, ok := result.(*TestService)
	assert.True(t, ok)
	assert.Equal(t, "LazyService", service.Name)

	// 第二次获取应该返回已缓存的实例
	initialized = false
	result2, err := registry.Get("lazy")
	assert.NoError(t, err)
	assert.False(t, initialized, "工厂函数不应该在第二次获取时执行")
	assert.Same(t, result, result2, "应该返回缓存的同一实例")
}

// 测试MustGet方法
func TestMustGet(t *testing.T) {
	registry := NewRegistry()
	service := &TestService{Name: "TestService"}
	registry.Register("test", service)

	// 正常情况
	result := registry.MustGet("test")
	assert.Same(t, service, result)

	// 服务不存在，应当panic
	assert.Panics(t, func() {
		registry.MustGet("nonexistent")
	})
}

// 测试Has方法
func TestHas(t *testing.T) {
	registry := NewRegistry()

	// 未注册的服务
	assert.False(t, registry.Has("test"))

	// 注册实例后检查
	registry.Register("test", &TestService{})
	assert.True(t, registry.Has("test"))

	// 注册工厂后检查
	registry.RegisterFactory("lazy", func() interface{} {
		return &TestService{}
	})
	assert.True(t, registry.Has("lazy"))
}

// 测试Unregister方法
func TestUnregister(t *testing.T) {
	registry := NewRegistry()

	registry.Register("test", &TestService{})
	assert.True(t, registry.Has("test"))

	registry.Unregister("test")
	assert.False(t, registry.Has("test"))

	// 测试注销工厂
	registry.RegisterFactory("lazy", func() interface{} {
		return &TestService{}
	})
	assert.True(t, registry.Has("lazy"))

	registry.Unregister("lazy")
	assert.False(t, registry.Has("lazy"))
}

// 测试Clear方法
func TestClear(t *testing.T) {
	registry := NewRegistry()

	registry.Register("test1", &TestService{})
	registry.Register("test2", &TestService{})
	registry.RegisterFactory("lazy", func() interface{} {
		return &TestService{}
	})

	assert.True(t, registry.Has("test1"))
	assert.True(t, registry.Has("test2"))
	assert.True(t, registry.Has("lazy"))

	registry.Clear()

	assert.False(t, registry.Has("test1"))
	assert.False(t, registry.Has("test2"))
	assert.False(t, registry.Has("lazy"))
}

// 测试Keys方法
func TestKeys(t *testing.T) {
	registry := NewRegistry()

	registry.Register("test1", &TestService{})
	registry.Register("test2", &TestService{})
	registry.RegisterFactory("lazy", func() interface{} {
		return &TestService{}
	})

	keys := registry.Keys()

	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "test1")
	assert.Contains(t, keys, "test2")
	assert.Contains(t, keys, "lazy")

	// 获取懒加载服务
	registry.Get("lazy")

	// 应该还是返回3个键，不会重复计数已实例化的工厂服务
	keys = registry.Keys()
	assert.Len(t, keys, 3)
}

// 测试错误情况：注册nil服务
func TestRegisterNil(t *testing.T) {
	registry := NewRegistry()

	err := registry.Register("test", nil)
	assert.Error(t, err)

	err = registry.RegisterFactory("lazy", nil)
	assert.Error(t, err)
}

// 测试错误情况：重复注册
func TestDuplicateRegistration(t *testing.T) {
	registry := NewRegistry()

	registry.Register("test", &TestService{})

	// 尝试重复注册
	err := registry.Register("test", &TestService{})
	assert.Error(t, err)

	// 尝试用工厂注册已有键
	err = registry.RegisterFactory("test", func() interface{} {
		return &TestService{}
	})
	assert.Error(t, err)

	// 注册工厂
	registry.RegisterFactory("lazy", func() interface{} {
		return &TestService{}
	})

	// 尝试重复注册工厂
	err = registry.RegisterFactory("lazy", func() interface{} {
		return &TestService{}
	})
	assert.Error(t, err)
}

// 测试错误情况：获取未注册的服务
func TestGetUnregisteredService(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Get("nonexistent")
	assert.Error(t, err)
}

// 测试工厂返回nil的情况
func TestFactoryReturnsNil(t *testing.T) {
	registry := NewRegistry()

	registry.RegisterFactory("nil", func() interface{} {
		return nil
	})

	_, err := registry.Get("nil")
	assert.Error(t, err)
}

// 并发测试
func TestConcurrentAccess(t *testing.T) {
	registry := NewRegistry()
	const goroutines = 100

	// 并发注册不同服务
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("service%d", id)
			registry.Register(key, &TestService{Name: key})
		}(i)
	}

	wg.Wait()

	// 验证所有服务都已注册
	for i := 0; i < goroutines; i++ {
		key := fmt.Sprintf("service%d", i)
		assert.True(t, registry.Has(key))
	}

	// 测试并发获取同一个懒加载服务
	registry.Clear()
	initCount := 0
	var mu sync.Mutex

	registry.RegisterFactory("concurrent", func() interface{} {
		mu.Lock()
		initCount++
		mu.Unlock()
		return &TestService{Name: "ConcurrentService"}
	})

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, _ = registry.Get("concurrent")
		}()
	}

	wg.Wait()

	// 验证工厂函数只被调用了一次
	assert.Equal(t, 1, initCount)
}

// 示例
func ExampleRegistry() {
	// 创建新注册表
	registry := NewRegistry()

	// 注册直接实例化的服务
	registry.Register("service", &TestService{Name: "Example"})

	// 注册懒加载服务
	registry.RegisterFactory("lazy", func() interface{} {
		return &TestService{Name: "LazyExample"}
	})

	// 获取服务
	service, _ := registry.Get("service")
	fmt.Println(service.(*TestService).GetName())

	// 获取懒加载服务
	lazyService, _ := registry.Get("lazy")
	fmt.Println(lazyService.(*TestService).GetName())

	// Output:
	// Example
	// LazyExample
}
