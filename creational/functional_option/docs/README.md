# 函数选项模式 (Functional Options Pattern)

函数选项模式是Go语言中一种流行的设计模式，用于创建具有大量可选配置的对象。这种模式允许提供灵活、可扩展且易用的API，同时保持向后兼容性。

## 1. 模式概述

函数选项模式通过使用一系列可选的函数作为配置参数，来解决以下问题：

- **过多构造函数参数**：当对象配置项很多时，构造函数参数列表会变得非常长。
- **参数顺序依赖**：传统构造函数中参数顺序固定，容易出错。
- **默认值处理**：频繁需要向构造函数传递默认值或零值。
- **API兼容性**：添加新配置项可能破坏现有API。

函数选项模式为这些问题提供了优雅的解决方案。

## 2. 实现原理

函数选项模式的核心包括：

1. **选项函数类型**：定义一个函数类型，它接受配置对象并修改它。
2. **默认配置**：创建带有默认值的基础配置对象。
3. **选项函数**：实现各种选项函数，每个函数修改配置的特定部分。
4. **构造函数**：接受可变数量的选项函数，应用它们并创建最终对象。

## 3. 代码示例

本示例实现了一个HTTP客户端构建器，展示了函数选项模式的用法：

```go
// 创建默认HTTP客户端
client := NewHTTPClient()

// 创建带有特定配置的客户端
client = NewHTTPClient(
    WithTimeout(5 * time.Second),
    WithMaxIdleConns(20),
    WithProxyURL("http://proxy.example.com:8080"),
    WithDisableCompression(true),
)

// 配置现有客户端
existingClient := &http.Client{...}
updatedClient := ConfigureHTTPClient(existingClient,
    WithMaxIdleConns(50),
    WithIdleConnTimeout(30 * time.Second),
)
```

## 4. 优缺点

### 优点

- **灵活性**：使用者可以只指定需要的选项。
- **可读性**：选项函数名称具有自描述性，提高代码可读性。
- **可扩展性**：添加新选项不会破坏现有代码。
- **默认值**：可以轻松定义合理的默认值。
- **封装**：选项函数可以封装复杂配置逻辑。

### 缺点

- **初始复杂性**：模式实现比简单构造函数更复杂。
- **运行时开销**：每个选项都是函数调用，有轻微的性能消耗。
- **必选参数处理**：不适合有必选参数的情况（可结合其他模式解决）。

## 5. 适用场景

函数选项模式特别适用于：

- 构造具有多个可选配置的复杂对象
- 需要提供合理默认值的API
- 预计将来需要添加更多配置选项的系统
- 需要保持API向后兼容性的库或框架

## 6. 最佳实践

1. **使用描述性命名**：选项函数名应清晰表达其作用，通常以"With"开头。
2. **验证输入**：选项函数应验证参数有效性。
3. **提供合理默认值**：为所有配置项设定合理默认值。
4. **选项函数应是幂等的**：多次应用同一选项应有一致效果。
5. **考虑选项组**：提供组合多个相关选项的辅助函数。
6. **使用私有配置结构**：将实现细节封装在包内部。

## 7. 示例代码结构

本示例的实现包含：

- **选项类型定义**：`Option func(*HTTPClientOptions)`
- **配置结构**：`HTTPClientOptions`包含所有可配置项
- **默认配置**：`defaultHTTPClientOptions()`提供合理默认值
- **选项函数**：如`WithTimeout()`、`WithProxy()`等
- **构造函数**：`NewHTTPClient()`和`ConfigureHTTPClient()`

通过这种模式，我们可以构建灵活、可扩展且用户友好的API。

## 8. 参考

- [Dave Cheney - Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
- [Rob Pike - Self-referential functions and the design of options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html)