package registry

import (
	"fmt"
	"sync"
)

// ServiceCreator 定义了创建服务实例的函数类型
type ServiceCreator func() interface{}

// Registry 定义注册表结构
type Registry struct {
	mutex     sync.RWMutex              // 用于并发安全
	services  map[string]interface{}    // 存储已实例化的服务
	factories map[string]ServiceCreator // 存储服务工厂函数
}

// NewRegistry 创建一个新的注册表实例
func NewRegistry() *Registry {
	return &Registry{
		services:  make(map[string]interface{}),
		factories: make(map[string]ServiceCreator),
	}
}

// 全局单例注册表实例
var (
	globalRegistry *Registry
	once           sync.Once
)

// GetRegistry 获取全局注册表单例
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = NewRegistry()
	})
	return globalRegistry
}

// Register 方法用于向注册表中注册已实例化的对象
func (r *Registry) Register(key string, service interface{}) error {
	if service == nil {
		return fmt.Errorf("不能注册nil服务")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.services[key]; exists {
		return fmt.Errorf("服务 '%s' 已经注册", key)
	}

	r.services[key] = service
	return nil
}

// RegisterFactory 注册一个服务创建工厂函数，推迟实例化到首次使用时
func (r *Registry) RegisterFactory(key string, creator ServiceCreator) error {
	if creator == nil {
		return fmt.Errorf("不能注册nil创建函数")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.services[key]; exists {
		return fmt.Errorf("服务 '%s' 已经注册", key)
	}

	if _, exists := r.factories[key]; exists {
		return fmt.Errorf("服务工厂 '%s' 已经注册", key)
	}

	r.factories[key] = creator
	return nil
}

// Get 方法用于从注册表中检索对象
func (r *Registry) Get(key string) (interface{}, error) {
	r.mutex.RLock()
	service, exists := r.services[key]
	r.mutex.RUnlock()

	if exists {
		return service, nil
	}

	// 检查是否有工厂可以创建此服务
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 二次检查，确保没有在获取锁期间被其他goroutine创建
	if service, exists := r.services[key]; exists {
		return service, nil
	}

	if factory, exists := r.factories[key]; exists {
		// 延迟实例化
		service = factory()
		if service == nil {
			return nil, fmt.Errorf("工厂方法返回nil对象")
		}
		r.services[key] = service
		return service, nil
	}

	return nil, fmt.Errorf("服务 '%s' 未注册", key)
}

// MustGet 获取服务，如果服务不存在则panic
func (r *Registry) MustGet(key string) interface{} {
	service, err := r.Get(key)
	if err != nil {
		panic(err)
	}
	return service
}

// GetTyped 获取指定类型的服务
func (r *Registry) GetTyped(key string, target interface{}) error {
	service, err := r.Get(key)
	if err != nil {
		return err
	}

	// 需要反射来将service赋值到target，此处简化实现
	// 实际使用时可以使用类似 *target = *service.(*TargetType) 的方式
	// 或者使用反射库进行类型安全的转换
	targetPtr, ok := target.(*interface{})
	if !ok {
		return fmt.Errorf("target必须是指针类型")
	}
	*targetPtr = service
	return nil
}

// Unregister 从注册表中删除服务
func (r *Registry) Unregister(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.services, key)
	delete(r.factories, key)
}

// Has 检查服务是否已注册
func (r *Registry) Has(key string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, existsService := r.services[key]
	_, existsFactory := r.factories[key]
	return existsService || existsFactory
}

// Clear 清空所有已注册的服务
func (r *Registry) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.services = make(map[string]interface{})
	r.factories = make(map[string]ServiceCreator)
}

// Keys 返回所有已注册的服务键
func (r *Registry) Keys() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	keys := make([]string, 0, len(r.services)+len(r.factories))

	for k := range r.services {
		keys = append(keys, k)
	}

	// 只添加尚未实例化的工厂键
	for k := range r.factories {
		if _, exists := r.services[k]; !exists {
			keys = append(keys, k)
		}
	}

	return keys
}
