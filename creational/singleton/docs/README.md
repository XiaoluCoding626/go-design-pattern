# 单例模式 (Singleton Pattern)

## 概述

单例模式是一种创建型设计模式，它确保一个类只有一个实例，并提供一个全局访问点来访问该实例。单例模式是最简单也是最常用的设计模式之一。

## 目录结构

```
singleton/
├── eager_singleton.go       # 饿汉式单例模式实现
├── eager_singleton_test.go  # 饿汉式单例模式测试
├── lazy_singleton.go        # 懒汉式单例模式实现
├── lazy_singleton_test.go   # 懒汉式单例模式测试
└── README.md                # 本文件
```

## 关键特性

- **单一实例**：全局只存在一个实例，避免重复创建
- **全局访问点**：提供了获取该实例的全局访问方法
- **延迟初始化**：懒汉式实现通常使用懒加载方式创建实例
- **线程安全**：多线程环境下能安全工作

## 实现方式

本项目提供了两种单例模式的Go语言实现:

### 1. 懒汉式实现 (Lazy Singleton)

懒汉式单例只有在首次请求时才创建实例。

```go
import "sync"

var (
    lazyInstance *Lazy
    lazyOnce sync.Once
)

func GetLazy() *Lazy {
    lazyOnce.Do(func() {
        lazyInstance = &Lazy{}
    })
    return lazyInstance
}
```

**适用场景**：实例占用资源较多，且不一定被使用时

**优点**：
- 第一次使用时才初始化，避免资源浪费
- 使用`sync.Once`确保线程安全且只初始化一次

### 2. 饿汉式实现 (Eager Singleton)

饿汉式单例在显式初始化时创建实例。

```go
var eagerInstance *Eager

func InitEager(count int) {
    eagerInstance = &Eager{count: count}
}

func GetEager() *Eager {
    return eagerInstance
}
```

**适用场景**：需要在程序启动时就确保实例存在，或需要外部提供初始化参数

**优点**：
- 实现简单，没有线程安全问题
- 可以控制初始化时机和参数

## 在并发环境中的表现

### 懒汉式单例

懒汉式单例使用`sync.Once`确保在多goroutine访问时也只初始化一次：

```go
// 测试例证明即使500个goroutine并发访问
// 所有goroutine获取的都是同一个实例
func TestConcurrentGetLazy(t *testing.T) {
    // 500个协程同时请求单例
    instances := [500]*Lazy{}
    for i := 0; i < 500; i++ {
        go func(index int) {
            instances[index] = GetLazy()
        }(i)
    }
    
    // 验证所有实例都是同一个
    firstInstance := instances[0]
    for i := 1; i < 500; i++ {
        assert.Same(firstInstance, instances[i])
    }
}
```

### 饿汉式单例

饿汉式单例在方法访问线程安全，但如果需要修改实例内部状态，需要实现额外的线程安全机制：

```go
type Eager struct {
    count int
    mu sync.Mutex // 保护count的并发访问
}

func (e *Eager) Increase() {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.count++
}

func (e *Eager) GetCount() int {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.count
}
```

## 使用方法

```go
// 懒汉式单例 - 自动初始化
instance := singleton.GetLazy()
instance.HelloWorld()

// 饿汉式单例 - 需要显式初始化
singleton.InitEager(10) // 设置初始值
instance := singleton.GetEager()
count := instance.GetCount() // 10
instance.Increase()
count = instance.GetCount() // 11
```

## 适用场景

- **配置管理**：需要全局访问的配置信息
- **资源池**：数据库连接池、线程池等需要集中管理的资源
- **日志记录器**：全局统一的日志记录
- **设备访问**：打印机、显卡等硬件访问控制
- **缓存**：全局统一的缓存管理

## 最佳实践

1. **选择合适的实现**：
   - 懒汉式用于资源占用较大、不一定使用的情况
   - 饿汉式用于需要确保初始化成功或需要参数的情况

2. **线程安全考虑**：
   - 实现单例时必须考虑并发安全
   - 访问和修改单例状态时也要考虑线程安全

3. **避免滥用**：
   - 单例模式引入了全局状态，过度使用会增加系统耦合度
   - 考虑使用依赖注入作为替代方案

4. **测试友好**：
   - 考虑提供重置单例的方法（仅用于测试）
   - 或者使用接口以便于模拟和测试

## 优缺点

### 优点
- 确保一个类只有一个实例
- 提供对实例的全局访问点
- 控制实例创建的时机和方式

### 缺点
- 引入全局状态，增加系统耦合度
- 可能导致并发问题
- 不适合需要多实例配置的场景
- 可能使单元测试变得困难

## 相关模式

- **工厂方法模式**：可以使用工厂实现单例的变体
- **享元模式**：与单例类似，但允许有限数量的实例
