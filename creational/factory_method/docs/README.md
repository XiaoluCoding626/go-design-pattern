# 工厂方法模式 (Factory Method Pattern)

## 简介

工厂方法模式是一种创建型设计模式，它定义了一个创建对象的接口，但由子类决定实例化的对象类型。工厂方法让一个类的实例化延迟到其子类。

在本示例中，我们实现了一个日志系统，允许应用程序将日志记录到不同的目标（文件、控制台和网络）。

## 结构

![工厂方法模式结构图](https://upload.wikimedia.org/wikipedia/commons/thumb/a/a3/FactoryMethod.svg/400px-FactoryMethod.svg.png)

### 核心组成部分

- **产品接口（Product）**: 定义工厂方法创建的对象的接口
- **具体产品（Concrete Product）**: 实现产品接口的具体类
- **创建者接口（Creator）**: 声明工厂方法，返回产品对象
- **具体创建者（Concrete Creator）**: 重写工厂方法以返回具体产品实例

## 代码实现

### 产品层次结构

我们定义了一个 `Logger` 接口作为产品接口，以及三种具体产品：

```go
// Logger 是日志记录器接口
type Logger interface {
    Log(message string)
}

// FileLogger 实现了记录到文件的日志记录器
type FileLogger struct{}

// ConsoleLogger 实现了记录到控制台的日志记录器
type ConsoleLogger struct{}

// NetworkLogger 是一个用于记录到网络的日志记录器
type NetworkLogger struct {
    endpoint string
}
```

### 工厂层次结构

我们定义了 `LoggerFactory` 接口作为创建者接口，以及三种具体创建者：

```go
// LoggerFactory 是定义工厂方法的接口
type LoggerFactory interface {
    CreateLogger() Logger
}

// FileLoggerFactory 创建 FileLogger 实例
type FileLoggerFactory struct {
    instance Logger
    once     sync.Once
}

// ConsoleLoggerFactory 创建 ConsoleLogger 实例
type ConsoleLoggerFactory struct {
    instance Logger
    once     sync.Once
}

// NetworkLoggerFactory 创建 NetworkLogger 实例
type NetworkLoggerFactory struct {
    endpoint string
}
```

### 工厂方法实现

每个具体工厂都实现了 `CreateLogger()` 方法：

```go
// FileLoggerFactory 的工厂方法
func (f *FileLoggerFactory) CreateLogger() Logger {
    // 使用单例模式和懒初始化
    f.once.Do(func() {
        f.instance = &FileLogger{}
    })
    return f.instance
}

// ConsoleLoggerFactory 的工厂方法
func (c *ConsoleLoggerFactory) CreateLogger() Logger {
    // 使用单例模式和懒初始化
    c.once.Do(func() {
        c.instance = &ConsoleLogger{}
    })
    return c.instance
}

// NetworkLoggerFactory 的工厂方法
func (n *NetworkLoggerFactory) CreateLogger() Logger {
    return NewNetworkLogger(n.endpoint)
}
```

### 工厂创建辅助函数

为了方便使用，我们提供了一个 `GetLoggerFactory` 函数：

```go
func GetLoggerFactory(loggerType LoggerType, config map[string]string) (LoggerFactory, error) {
    switch loggerType {
    case FileType:
        return &FileLoggerFactory{}, nil
    case ConsoleType:
        return &ConsoleLoggerFactory{}, nil
    case NetworkType:
        endpoint, ok := config["endpoint"]
        if !ok {
            return nil, fmt.Errorf("网络日志记录器需要endpoint配置")
        }
        return NewNetworkLoggerFactory(endpoint), nil
    default:
        return nil, fmt.Errorf("不支持的日志记录器类型: %s", loggerType)
    }
}
```

## 使用示例

```go
func main() {
    // 创建文件日志记录器
    fileFactory, _ := GetLoggerFactory(FileType, nil)
    fileLogger := fileFactory.CreateLogger()
    fileLogger.Log("这是一条文件日志")
    
    // 创建控制台日志记录器
    consoleFactory, _ := GetLoggerFactory(ConsoleType, nil)
    consoleLogger := consoleFactory.CreateLogger()
    consoleLogger.Log("这是一条控制台日志")
    
    // 创建网络日志记录器
    config := map[string]string{"endpoint": "http://example.com/log"}
    networkFactory, _ := GetLoggerFactory(NetworkType, config)
    networkLogger := networkFactory.CreateLogger()
    networkLogger.Log("这是一条网络日志")
}
```

## 优点

1. **遵循开闭原则**：可以引入新的产品类型，而无需修改现有代码
2. **遵循单一职责原则**：将产品创建代码放在程序的一个地方
3. **可以实现延迟初始化**：如本例中使用懒加载的单例模式
4. **松耦合**：创建者和具体产品解耦

## 缺点

1. **增加复杂性**：需要引入许多新的子类
2. **客户端必须了解不同的产品子类**：为了正确使用工厂

## 适用场景

1. 当事先不知道需要使用哪种具体产品时
2. 当系统需要独立于它所创建的产品时
3. 当需要为产品类提供扩展时

## 与其他模式的关系

- **抽象工厂模式**：工厂方法通常在抽象工厂的实现中使用
- **模板方法模式**：工厂方法是模板方法的一种特殊形式
- **原型模式**：不需要子类化创建者时，可以使用原型模式代替工厂方法