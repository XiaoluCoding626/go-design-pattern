# 策略模式 (Strategy Pattern)

## 简介

策略模式是一种行为设计模式，它定义了一系列算法，将每个算法封装起来，并使它们可以互换。策略模式使算法可以独立于使用它的客户端而变化。

## 结构

策略模式主要包含以下几个角色：

1. **策略接口（Strategy）**：定义了所有支持的算法的通用接口
2. **具体策略（Concrete Strategy）**：实现了策略接口的具体算法
3. **上下文（Context）**：持有一个策略对象的引用，并在需要时调用策略对象的方法

## 优点

- 可以在运行时切换算法（策略）
- 将算法的实现和使用分离
- 消除了大量的条件判断
- 符合开闭原则，增加新的策略不需要修改上下文代码

## 缺点

- 客户端必须了解所有的策略，以便选择合适的策略
- 增加了对象的数量
- 策略和上下文之间可能存在通信开销

## 适用场景

- 当需要使用不同的算法变体，并且可以在运行时切换算法时
- 当有许多相似的类，它们只在行为上有所不同时
- 当算法的使用涉及复杂的条件语句时

## 示例实现

在本例中，我们实现了一个简单的计算器，可以使用不同的策略（加法、减法、乘法、除法）来处理两个数字：

### 策略接口

```go
// IStrategy 定义所有支持算法的接口
type IStrategy interface {
    // Do 执行具体的算法，并返回结果和可能的错误
    Do(a, b int) (int, error)
}
```

### 具体策略实现

```go
// Add 实现加法策略
type Add struct{}

// Do 执行加法操作
func (*Add) Do(a, b int) (int, error) {
    return a + b, nil
}

// 更多策略实现：Subtract, Multiply, Divide
```

### 上下文实现

```go
// Operator 是使用策略的上下文
type Operator struct {
    strategy IStrategy
}

// SetStrategy 更改操作者使用的策略
func (o *Operator) SetStrategy(strategy IStrategy) {
    o.strategy = strategy
}

// Calculate 将当前策略应用于给定的操作数
func (o *Operator) Calculate(a, b int) (int, error) {
    if o.strategy == nil {
        return 0, fmt.Errorf("未设置策略")
    }
    return o.strategy.Do(a, b)
}
```

## 使用示例

```go
// 创建操作者并设置初始策略为加法
operator := NewOperator(&Add{})

// 使用加法策略
result, _ := operator.Calculate(5, 3)
fmt.Println("5 + 3 =", result)  // 输出：5 + 3 = 8

// 切换到减法策略
operator.SetStrategy(&Subtract{})
result, _ = operator.Calculate(5, 3)
fmt.Println("5 - 3 =", result)  // 输出：5 - 3 = 2

// 切换到乘法策略
operator.SetStrategy(&Multiply{})
result, _ = operator.Calculate(5, 3)
fmt.Println("5 * 3 =", result)  // 输出：5 * 3 = 15

// 切换到除法策略
operator.SetStrategy(&Divide{})
result, _ = operator.Calculate(6, 3)
fmt.Println("6 / 3 =", result)  // 输出：6 / 3 = 2
```

## 错误处理

该实现包含了错误处理，例如处理除以零的情况：

```go
// Divide 实现除法策略
type Divide struct{}

// Do 执行除法操作，并处理除零错误
func (*Divide) Do(a, b int) (int, error) {
    if b == 0 {
        return 0, ErrDivideByZero
    }
    return a / b, nil
}

// ErrDivideByZero 当尝试除以零时返回此错误
var ErrDivideByZero = fmt.Errorf("不能除以零")
```

## 运行测试

可以使用以下命令运行测试：

```bash
cd /path/to/go-design-pattern/behavioral/strategy
go test -v
```

测试用例覆盖了所有策略的基本功能，以及错误处理和边界情况。

## 扩展思考

1. 如何扩展更多的算法策略？比如添加取模、幂运算等。
2. 如何在不同业务场景中应用策略模式？例如支付系统中的不同支付方式、排序算法中的不同排序策略等。
3. 策略模式和状态模式的区别是什么？（提示：二者结构相似，但意图不同）

## 参考资料

- 《Design Patterns: Elements of Reusable Object-Oriented Software》（设计模式：可复用面向对象软件的基础）
- 《Head First Design Patterns》（深入浅出设计模式）