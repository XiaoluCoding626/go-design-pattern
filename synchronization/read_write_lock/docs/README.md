# 读写锁模式 (Read-Write Lock Pattern)

## 简介

读写锁是一种特殊的锁机制，它允许多个读操作并发执行，但写操作必须互斥。这种机制在读多写少的场景下能显著提高系统的并发性能。

本实现提供了一个完整的读写锁模式，包括锁接口定义、标准实现以及保护共享数据的封装。具有以下特点：
- 提供统一的接口，方便替换不同实现
- 支持基本的读写锁操作
- 支持非阻塞的尝试获取锁
- 支持带超时的尝试获取锁
- 支持回调风格的读写操作

## 功能特性

- **基本读写锁操作**：支持常规的读锁/写锁获取与释放
- **非阻塞操作**：通过`TryReadLock`和`TryWriteLock`方法实现非阻塞的锁获取
- **超时机制**：支持在指定时间内尝试获取锁
- **读写回调**：使用回调函数简化读写操作
- **依赖注入**：通过接口设计支持测试和不同实现的替换

## 接口设计

### RWLocker 接口

```go
// RWLocker 定义了读写锁接口
type RWLocker interface {
    ReadLock()                                  // 获取读锁
    ReadUnlock()                                // 释放读锁
    WriteLock()                                 // 获取写锁
    WriteUnlock()                               // 释放写锁
    TryReadLock() bool                          // 尝试获取读锁，不阻塞
    TryWriteLock() bool                         // 尝试获取写锁，不阻塞
    TryReadLockWithTimeout(time.Duration) bool  // 带超时的尝试获取读锁
    TryWriteLockWithTimeout(time.Duration) bool // 带超时的尝试获取写锁
}
```

### StandardRWLock 实现

`StandardRWLock`是对Go标准库中`sync.RWMutex`的封装，实现了`RWLocker`接口：

```go
// StandardRWLock 标准读写锁实现
type StandardRWLock struct {
    rwMutex sync.RWMutex
}
```

### Data 类型

`Data`类型演示了如何使用读写锁保护共享数据：

```go
// Data 表示包含读写锁保护的共享数据
type Data struct {
    locker RWLocker // 使用接口允许注入不同的读写锁实现
    value  int      // 数据值
}
```

## 使用示例

### 基本读写操作

```go
// 创建一个新的数据实例
data := NewData()

// 写入数据
data.Write(100)

// 读取数据
value := data.Read()
fmt.Printf("读取的值: %d\n", value)
```

### 尝试读取/写入

```go
// 尝试读取，不阻塞
if val, ok := data.TryRead(); ok {
    fmt.Printf("成功读取: %d\n", val)
} else {
    fmt.Println("读取失败，当前有写锁")
}

// 尝试写入，不阻塞
if ok := data.TryWrite(200); ok {
    fmt.Println("成功写入")
} else {
    fmt.Println("写入失败，当前有读锁或写锁")
}
```

### 超时机制

```go
// 尝试在100毫秒内获取读锁
if val, ok := data.ReadWithTimeout(100 * time.Millisecond); ok {
    fmt.Printf("成功读取: %d\n", val)
} else {
    fmt.Println("读取超时")
}

// 尝试在100毫秒内获取写锁
if ok := data.WriteWithTimeout(300, 100 * time.Millisecond); ok {
    fmt.Println("成功写入")
} else {
    fmt.Println("写入超时")
}
```

### 回调函数

```go
// 在读锁保护下执行操作
data.ReadWithCallback(func(val int) {
    fmt.Printf("当前值: %d\n", val)
    // 可以执行其他只读操作...
})

// 在写锁保护下执行多个操作
data.WriteWithCallback(func(d *Data) {
    d.value = d.value * 2
    // 可以执行其他修改操作...
})
```

### 读取-修改-写入模式

```go
// 读取值，计算新值，然后写入
data.ReadWriteWithCallback(func(val int) int {
    return val * 2 // 返回新值
})
```

### 依赖注入

```go
// 创建自定义锁实现
customLocker := MyCustomRWLocker{}

// 使用自定义锁创建数据实例
data := NewDataWithLocker(customLocker)
```

## 使用场景

读写锁特别适合以下场景：

1. **读多写少的数据结构**：例如配置信息、缓存等
2. **需要并发读取但写入较少的共享资源**：例如数据库连接池
3. **读操作远多于写操作且读操作耗时长的情况**：例如大型数据分析

## 性能考量

- 读写锁比互斥锁有更大的开销，只有在读操作显著多于写操作时才有性能优势
- 短时间的临界区可能不适合使用读写锁，因为锁的开销可能超过临界区执行时间
- 在高并发环境中，频繁的锁争用可能导致性能下降

## 注意事项

- 读锁不会阻塞其他读操作，但会阻塞写操作
- 写锁会阻塞所有的读操作和写操作
- 不要在持有读锁的情况下尝试获取写锁，可能导致死锁
- `ReadWriteWithCallback`方法不是原子操作，中间会释放读锁再获取写锁