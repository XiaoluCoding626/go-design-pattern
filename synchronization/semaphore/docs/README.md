# 信号量模式 (Semaphore Pattern)

## 简介

信号量是一种用于控制对共享资源访问的同步机制，它维护一个计数器，表示可用资源的数量。当资源被获取时，计数器减少；当资源被释放时，计数器增加。如果计数器为零，尝试获取资源的线程将被阻塞，直到其他线程释放资源。

本实现提供了两种信号量实现：
1. **标准信号量** - 用于控制对资源的并发访问数量
2. **带权重的信号量** - 支持为不同的操作分配不同的资源消耗权重

## 功能特性

- **资源计数管理** - 维护可用资源数量，自动跟踪已获取资源
- **阻塞与非阻塞操作** - 支持阻塞式获取和非阻塞式尝试获取
- **上下文支持** - 与Go的context包集成，支持超时和取消操作
- **批量操作** - 支持批量获取和释放资源
- **并发安全** - 所有操作都是线程安全的，适用于高并发环境
- **权重分配** - 带权重的信号量允许为不同操作分配不同的资源消耗

## 接口设计

### Semaphorer 接口

```go
// Semaphorer 定义了信号量应该具有的行为
type Semaphorer interface {
    // Acquire 尝试获取一个票证，可能会阻塞直到有票证可用或超时
    Acquire(ctx context.Context) error

    // TryAcquire 尝试非阻塞地获取一个票证，立即返回结果
    TryAcquire() bool

    // AcquireMany 尝试获取多个票证
    AcquireMany(n int, ctx context.Context) error

    // Release 释放一个已获取的票证
    Release() error

    // ReleaseMany 释放多个已获取的票证
    ReleaseMany(n int) error

    // Available 返回当前可用的票证数量
    Available() int

    // Size 返回信号量的总容量
    Size() int
}
```

### 信号量类型

#### 标准信号量 (Semaphore)

标准信号量管理固定数量的资源，每个资源权重相同。

```go
// Semaphore 实现了信号量设计模式
type Semaphore struct {
    tickets chan struct{} // 使用通道实现信号量机制
    size int             // 信号量的总容量
    mu sync.Mutex        // 保护计数器的互斥锁
    acquired int         // 已获取的票证数量
}
```

#### 带权重的信号量 (WeightedSemaphore)

带权重的信号量允许获取不同权重的资源，适用于资源消耗不均匀的场景。

```go
// WeightedSemaphore 实现了带权重的信号量
type WeightedSemaphore struct {
    capacity int64      // 总容量
    used int64          // 当前使用量
    mu sync.Mutex       // 保护并发访问
    cond *sync.Cond     // 当资源被释放时通知等待者
}
```

## 使用示例

### 基本用法

```go
// 创建一个容量为5的信号量
sem := semaphore.New(5)

// 获取一个资源
err := sem.Acquire(context.Background())
if err != nil {
    // 处理错误
}

// 执行需要保护的操作
// ...

// 释放资源
sem.Release()
```

### 带超时的获取

```go
// 创建一个容量为3的信号量
sem := semaphore.New(3)

// 尝试在100毫秒内获取资源
err := sem.AcquireWithTimeout(100 * time.Millisecond)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("获取资源超时")
    }
} else {
    // 成功获取资源，执行操作
    // ...
    sem.Release()
}
```

### 非阻塞尝试获取

```go
sem := semaphore.New(2)

// 非阻塞尝试获取资源
if sem.TryAcquire() {
    // 成功获取资源，执行操作
    // ...
    sem.Release()
} else {
    fmt.Println("资源当前不可用，执行替代操作")
}
```

### 批量操作

```go
sem := semaphore.New(10)
ctx := context.Background()

// 批量获取5个资源
err := sem.AcquireMany(5, ctx)
if err != nil {
    // 处理错误
} else {
    // 执行需要多个资源的操作
    // ...
    
    // 批量释放
    sem.ReleaseMany(5)
}
```

### 带权重的信号量

```go
// 创建总容量为100的带权重信号量
ws := semaphore.NewWeighted(100)
ctx := context.Background()

// 获取权重为30的资源
err := ws.Acquire(ctx, 30)
if err != nil {
    // 处理错误
}

// 获取权重为20的资源
err = ws.Acquire(ctx, 20)
if err != nil {
    // 处理错误
}

// 释放权重为30的资源
ws.Release(30)

// 检查可用资源
available := ws.Available() // 应为50
```

### 在函数退出时自动释放

```go
func processWithSemaphore(sem *semaphore.Semaphore) error {
    // 获取资源
    if err := sem.Acquire(context.Background()); err != nil {
        return err
    }
    
    // 确保在函数退出时释放资源
    defer sem.Release()
    
    // 执行受保护的操作
    // ...
    
    return nil
}
```

## 应用场景

信号量模式在以下场景特别有用：

1. **限制并发访问数**：控制同时访问共享资源（如数据库连接、文件句柄）的并发数
2. **实现资源池**：管理有限的资源池（如连接池、线程池）
3. **流量控制**：限制API请求速率，防止系统过载
4. **并行任务控制**：限制同时执行的并行任务数量
5. **带权重的资源分配**：为不同操作分配不同的资源消耗量

## 带权重信号量的使用场景

1. **数据库连接池**：不同查询消耗不同数量的连接资源
2. **混合工作负载**：同时处理轻量级和重量级操作
3. **资源消耗不均匀的应用**：操作消耗的资源数量差异较大

## 性能考量

- **信号量通道实现**：使用Go通道提供高效的阻塞和同步机制
- **非阻塞操作**：`TryAcquire`方法用于避免不必要的阻塞
- **上下文集成**：支持优雅取消和超时控制
- **批量操作优化**：`AcquireMany`和`ReleaseMany`方法优化了批量资源管理

## 注意事项

1. **避免死锁**：确保释放所有获取的资源，推荐使用defer语句
2. **防止资源泄露**：在错误处理路径上也要释放资源
3. **选择合适的容量**：信号量容量过小会导致系统瓶颈，过大则失去控制作用
4. **并发安全**：信号量操作是并发安全的，但受保护的资源可能需要额外的同步
5. **避免长时间持有资源**：尽量减少资源持有时间，提高系统吞吐量