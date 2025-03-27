# 注册表模式 (Registry Pattern)

## 简介

注册表模式是一种创建型设计模式，它提供了一个全局访问点来存储、检索和管理对象实例。这种模式常用于需要集中管理应用程序组件或服务的场景，尤其是在依赖注入和服务定位方面。

注册表模式作为一种服务定位器(Service Locator)的实现，为应用程序提供了一种松耦合的方式来管理服务依赖关系。

## 结构和组件

```mermaid
classDiagram
    class Registry {
        -services: map[string]interface{}
        -factories: map[string]ServiceCreator
        -mutex: sync.RWMutex
        +NewRegistry() *Registry
        +GetRegistry() *Registry
        +Register(key string, service interface{}) error
        +RegisterFactory(key string, creator ServiceCreator) error
        +Get(key string) (interface{}, error)
        +MustGet(key string) interface{}
        +GetTyped(key string, target interface{}) error
        +Unregister(key string)
        +Has(key string) bool
        +Clear()
        +Keys() []string
    }
    
    class ServiceCreator {
        <<interface>>
        +Create() interface{}
    }
    
    class Service {
        <<interface>>
    }
    
    Registry --> ServiceCreator : 使用
    ServiceCreator --> Service : 创建
    Registry --> Service : 管理
```

注册表模式的主要组件包括：

1. **注册表（Registry）**：管理服务实例的中央存储库
2. **服务创建器（ServiceCreator）**：负责创建服务实例的工厂函数
3. **服务（Service）**：被注册和管理的实际组件

## 实现特性

我们的注册表模式实现具有以下特性：

### 核心功能

- **服务注册与检索**：支持通过唯一键注册和获取服务
- **懒加载**：支持通过工厂函数延迟初始化服务，直到首次请求时才创建实例
- **单例模式**：提供全局单一注册表实例

### 高级特性

- **线程安全**：使用互斥锁确保并发环境下的安全操作
- **错误处理**：提供全面的错误检查和报告
- **服务管理**：支持检查、移除和清理已注册的服务
- **双重检查锁定**：优化并发性能的获取服务实现

## 使用场景

注册表模式适用于以下场景：

1. **依赖管理**：在大型应用程序中管理组件依赖关系
2. **服务定位**：允许组件发现和使用其他服务而不需要直接依赖
3. **资源池管理**：管理可重用资源，如数据库连接或线程池
4. **延迟初始化**：需要推迟昂贵资源的创建到首次使用时
5. **测试模拟**：在测试环境中轻松替换服务实现

## 代码示例

### 基本使用

```go
// 创建注册表实例或获取全局单例
registry := registry.GetRegistry()

// 注册已实例化的服务
service := &UserService{}
registry.Register("userService", service)

// 获取服务
if userService, err := registry.Get("userService"); err == nil {
    // 使用服务
    userService.(*UserService).FindUser(123)
}
```

### 使用工厂函数实现懒加载

```go
// 注册工厂函数，延迟初始化到首次使用时
registry.RegisterFactory("expensiveService", func() interface{} {
    return NewExpensiveService()
})

// 首次获取时才会创建实例
service, _ := registry.Get("expensiveService")
```

### 类型安全的服务获取

```go
var userService *UserService
err := registry.GetTyped("userService", &userService)
if err == nil {
    userService.FindUser(123)
}
```

### 服务管理

```go
// 检查服务是否存在
if registry.Has("logService") {
    // 服务存在
}

// 移除服务
registry.Unregister("tempService")

// 获取所有注册的服务键
keys := registry.Keys()
```

## 优点

1. **减少耦合**：组件之间通过注册表间接交互，而不是直接依赖
2. **中央配置**：提供统一的服务配置和访问点
3. **灵活性**：可以在运行时动态替换服务实现
4. **延迟加载**：支持懒加载，提高应用启动性能
5. **资源管理**：集中管理服务生命周期

## 缺点

1. **隐藏依赖**：可能隐藏组件的实际依赖关系，使代码难以理解
2. **全局状态**：引入了全局状态，可能使测试变得困难
3. **类型安全**：使用interface{}需要类型断言，可能导致运行时错误

## 最佳实践

1. **有意识地使用**：了解这种模式的优缺点，在合适的场景下使用
2. **避免过度使用**：不要将其作为依赖注入的完全替代
3. **使用工厂函数**：优先使用懒加载工厂函数而不是直接注册实例
4. **错误处理**：始终检查Get方法的错误返回值
5. **命名标准**：为服务使用一致的命名约定

## 与其他模式的关系

- **服务定位器**：注册表是服务定位器模式的一种实现
- **单例模式**：注册表通常作为单例实现
- **工厂模式**：注册表使用工厂函数进行延迟初始化
- **依赖注入**：注册表可以与依赖注入结合使用，但不应完全替代它

## 总结

注册表模式提供了一种集中管理服务实例的方法，特别适合需要灵活组件配置和服务发现的复杂应用程序。虽然它有助于减少组件间的直接耦合，但应谨慎使用，避免创建难以测试和维护的全局状态。

在Go语言中，注册表模式可以与接口和类型系统结合，创建类型安全且功能强大的服务管理解决方案。