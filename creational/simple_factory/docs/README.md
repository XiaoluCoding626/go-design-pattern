# 简单工厂模式 (Simple Factory Pattern)

## 概述

简单工厂模式不属于GoF的23种设计模式，但它是一种常用的对象创建模式。该模式通过一个工厂类或函数，根据参数的不同返回不同类型的实例，从而将客户端与对象的创建过程解耦。

## 目录结构

```
simple_factory/
├── simple_factory.go      # 简单工厂实现
├── simple_factory_test.go # 测试代码
└── README.md              # 本文档
```

## UML 类图

```
┌────────────────┐      creates      ┌────────────────┐
│  ShapeFactory  │ ─────────────────>│     Shape      │
└────────────────┘                   │    (Interface) │
        │                            └────────┬───────┘
        │                                     │
        │                                     │
        │                                     │
        │                      ┌──────────────┼──────────────┐
        │                      │              │              │
        │                      │              │              │
        │                 ┌────▼────┐    ┌────▼────┐    ┌────▼────┐
        └────────────────>│ Circle  │    │Rectangle│    │Triangle │
                          └─────────┘    └─────────┘    └─────────┘
```

## 核心组件

### 1. 产品接口 (Shape)

所有由工厂创建的对象都实现相同的接口：

```go
type Shape interface {
    Draw() string
    GetType() ShapeType
}
```

### 2. 具体产品类

- **Circle**: 实现了圆形
- **Rectangle**: 实现了矩形
- **Triangle**: 实现了三角形

### 3. 工厂类 (ShapeFactory)

工厂类提供了创建不同形状的方法：

```go
// 通过类型创建形状
factory.CreateShape(ShapeTypeCircle, 5.0)

// 通过名称创建形状
factory.CreateShapeByName("rectangle", 10.0, 20.0)
```

## 使用方法

### 1. 创建工厂实例

```go
factory := simple_factory.NewShapeFactory()
```

### 2. 通过形状类型创建对象

```go
// 创建半径为5的圆形
circle, err := factory.CreateShape(simple_factory.ShapeTypeCircle, 5.0)
if err != nil {
    // 处理错误
}
fmt.Println(circle.Draw()) // 输出: Drawing Circle with radius 5.00

// 创建宽10高20的矩形
rectangle, err := factory.CreateShape(simple_factory.ShapeTypeRectangle, 10.0, 20.0)
if err != nil {
    // 处理错误
}
fmt.Println(rectangle.Draw()) // 输出: Drawing Rectangle with width 10.00 and height 20.00

// 创建三条边长为3,4,5的三角形
triangle, err := factory.CreateShape(simple_factory.ShapeTypeTriangle, 3.0, 4.0, 5.0)
if err != nil {
    // 处理错误
}
fmt.Println(triangle.Draw()) // 输出: Drawing Triangle with sides 3.00, 4.00, 5.00
```

### 3. 通过形状名称创建对象

```go
// 创建半径为3.5的圆形
circle, err := factory.CreateShapeByName("Circle", 3.5)
if err != nil {
    // 处理错误
}

// 名称不区分大小写
rectangle, err := factory.CreateShapeByName("rectangle", 4.0, 5.0)
if err != nil {
    // 处理错误
}
```

### 4. 向后兼容的方法 (不推荐用于新代码)

```go
// 使用旧版工厂函数创建默认配置的形状
circle := simple_factory.NewShape("circle") // 创建半径为1.0的圆
```

## 适用场景

简单工厂模式适用于以下场景：

1. **创建对象的过程相对简单**：不需要复杂的构建过程
2. **客户端无需关心对象的具体创建细节**：只需告诉工厂需要什么类型的对象
3. **系统中处理对象的方式大致相同**：所有对象实现相同的接口
4. **需要在运行时动态决定创建哪种对象**：基于配置或用户输入等创建对象

## 优点

- **封装变化点**：将对象的创建逻辑集中在一个地方，当需要新增或修改对象类型时只需修改工厂类
- **减少重复代码**：统一的对象创建过程避免了在多处编写相似代码
- **隐藏实现细节**：客户端无需了解对象的构造过程，只需知道如何使用
- **参数验证集中处理**：在工厂中统一处理参数验证和异常情况

## 缺点

- **工厂类职责过重**：随着系统扩展，工厂类可能变得庞大
- **扩展困难**：每次添加新产品时都需要修改工厂类代码，违反开闭原则
- **无法实现复杂的创建逻辑**：对于需要复杂初始化过程的对象不适用

## 与其他创建型模式的区别

- **工厂方法模式**：使用继承而非集中的工厂类，每种产品由专门的工厂子类创建
- **抽象工厂模式**：创建一系列相关对象，而非单一对象
- **构建器模式**：专注于分步骤构建复杂对象

## 最佳实践

1. **适当的错误处理**：工厂方法应返回错误而不是简单返回nil
2. **清晰的类型命名**：使用明确的类型名称和枚举值
3. **考虑扩展性**：设计易于扩展的工厂结构
4. **避免过度使用**：简单场景下直接创建对象可能更清晰

## 代码测试覆盖

项目中的测试文件 `simple_factory_test.go` 提供了全面的测试：

- 所有形状类型的创建和验证
- 参数验证和错误处理
- 不同调用方式的测试
- 向后兼容性验证
- 接口实现的完整性验证

运行测试：

```bash
cd simple_factory
go test -v
```