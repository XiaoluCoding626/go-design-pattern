# 解释器模式（Interpreter Pattern）

解释器模式是一种行为设计模式，它用于定义一个语言的语法，以及解释该语言中的表达式。这种模式特别适用于需要解析简单语言或格式化文本的场景。

## 概述

解释器模式主要包括以下几个关键组件：

1. **抽象表达式（Abstract Expression）**：定义解释操作的接口
2. **终结符表达式（Terminal Expression）**：实现与语法中终结符相关的解释操作
3. **非终结符表达式（Non-terminal Expression）**：实现与语法中非终结符相关的解释操作
4. **上下文（Context）**：包含解释器需要的全局信息
5. **客户端（Client）**：构建抽象语法树并调用解释操作

## 实现详解

本项目实现了一个简单的算术表达式解释器，支持整数、变量、加减乘除模运算以及括号优先级。

### 关键组件

#### 1. 表达式接口 (Expression Interface)

所有表达式类型的基础接口：

```go
type Expression interface {
    Interpret(context *Context) (int, error)
    String() string
}
```

#### 2. 终结符表达式

- **NumberExpression**：表示数字字面量
- **VariableExpression**：表示变量引用

#### 3. 非终结符表达式

- **AddExpression**：表示加法运算
- **SubtractExpression**：表示减法运算
- **MultiplyExpression**：表示乘法运算
- **DivideExpression**：表示除法运算
- **ModuloExpression**：表示取模运算

#### 4. 上下文环境 (Context)

存储和管理变量及其值：

```go
type Context struct {
    variables map[string]int
}
```

#### 5. 解析器 (Parser)

负责将表达式字符串解析为抽象语法树：

```go
type Parser struct {
    context *Context
    tokens  []string
    pos     int
}
```

### 表达式解析过程

解析过程遵循以下步骤：

1. **词法分析**：将输入表达式字符串拆分为标记（tokens）
2. **语法分析**：根据语法规则构建抽象语法树
3. **解释执行**：遍历抽象语法树，计算表达式结果

### 运算符优先级

解释器实现了标准的运算符优先级：

1. 括号 `()`：最高优先级
2. 乘法 `*`、除法 `/`、取模 `%`：中等优先级
3. 加法 `+`、减法 `-`：最低优先级

## 使用示例

### 基本使用

```go
// 创建上下文并设置变量
context := NewContext()
context.SetVariable("x", 10)
context.SetVariable("y", 5)

// 计算表达式
result, err := Evaluate("x + y * 2", context)
if err != nil {
    log.Fatalf("表达式解析错误: %v", err)
}
fmt.Printf("结果: %d\n", result)  // 输出: 结果: 20
```

### 手动构建表达式树

```go
// 创建上下文
context := NewContext()
context.SetVariable("x", 5)
context.SetVariable("y", 7)

// 手动构建表达式树: (3 + x) * (y - 2)
three := NewNumberExpression(3)
x := NewVariableExpression("x")
y := NewVariableExpression("y")
two := NewNumberExpression(2)

add := NewAddExpression(three, x)
subtract := NewSubtractExpression(y, two)
multiply := NewMultiplyExpression(add, subtract)

// 计算结果
result, _ := multiply.Interpret(context)
fmt.Printf("(3 + x) * (y - 2) = %d\n", result)  // 输出: (3 + x) * (y - 2) = 40
```

## 设计考量

1. **错误处理**：通过返回错误值处理变量未定义、除零等异常
2. **可扩展性**：易于添加新的表达式类型（如幂运算）
3. **分离关注点**：词法分析、语法分析和解释执行分离
4. **运算符优先级**：通过递归下降解析器实现正确的运算符优先级

## 优缺点

### 优点

- 易于扩展语法规则
- 易于实现简单语言的解释器
- 将操作与表达式结构分离

### 缺点

- 复杂语法的解释器可能变得很庞大
- 对于非常复杂的语法，可能不如专用的解析器生成工具
- 性能可能不如编译型实现

## 应用场景

- 简单语言或DSL（领域特定语言）的解释器
- SQL解析器
- 正则表达式引擎
- 公式计算器
- 配置文件解析

## 总结

解释器模式是实现简单语言解释器的有效方式，它将语法规则表示为面向对象的类结构，使得语言能够被轻松地扩展和修改。本项目实现的算术表达式解释器展示了如何使用这种模式来解析和评估数学表达式。