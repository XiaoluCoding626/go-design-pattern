# 原型模式 (Prototype Pattern)

## 简介

原型模式是一种创建型设计模式，它允许你通过复制现有对象来创建新对象，而无需依赖它们所属的类。当创建对象的成本较高或者创建过程复杂时，原型模式特别有用。

在本项目中，我们通过几何形状（圆形、矩形、三角形）的例子来展示原型模式的实现。每种形状都可以被克隆（复制）以创建新的形状实例，而无需了解形状的具体类型或重新执行繁琐的初始化过程。

## 设计思想

![原型模式UML图](https://refactoringguru.cn/images/patterns/diagrams/prototype/structure.png)

原型模式的核心组件：

1. **原型接口(Prototype)**: 声明克隆方法的接口。在我们的例子中是`Shape`接口。
2. **具体原型(Concrete Prototype)**: 实现克隆方法的具体类。在我们的例子中是`Circle`、`Rectangle`和`Triangle`。
3. **客户端(Client)**: 通过调用原型的克隆方法来创建新对象。在我们的例子中，可以是任何使用形状对象的代码。
4. **原型管理器(可选)**: 用于管理和检索可用原型。在我们的例子中是`ShapeCache`。

## 代码实现

### Shape 接口

```go
type Shape interface {
    Clone() Shape         // 浅克隆
    DeepClone() Shape     // 深克隆
    GetType() string      // 获取形状类型
    GetColor() Color      // 获取颜色
    SetColor(color Color) // 设置颜色
    GetArea() float64     // 计算面积
    String() string       // 字符串表示
}
```

### 具体形状类

我们实现了三种形状：

1. **Circle (圆形)**
```go
type Circle struct {
    BaseShape
    Radius float64
    Center *Point
}
```

2. **Rectangle (矩形)**
```go
type Rectangle struct {
    BaseShape
    Width    float64
    Height   float64
    Position *Point
}
```

3. **Triangle (三角形)**
```go
type Triangle struct {
    BaseShape
    A, B, C *Point // 三个顶点
}
```

### 浅克隆与深克隆

原型模式提供了两种克隆对象的方式：

1. **浅克隆 (Shallow Clone)**
   
   浅克隆只复制对象本身，不复制对象内部的引用类型。例如，Circle的浅克隆会创建一个新的Circle对象，但新旧对象会共享同一个Center指针。

   ```go
   func (c *Circle) Clone() Shape {
       return &Circle{
           BaseShape: BaseShape{
               Type:  c.Type,
               Color: c.Color,
           },
           Radius: c.Radius,
           Center: c.Center, // 共享同一个指针
       }
   }
   ```

2. **深克隆 (Deep Clone)**
   
   深克隆会复制对象本身及其所有引用的对象，创建一个完全独立的副本。例如，Circle的深克隆会创建新的Circle和Point对象。

   ```go
   func (c *Circle) DeepClone() Shape {
       return &Circle{
           BaseShape: BaseShape{
               Type:  c.Type,
               Color: c.Color,
           },
           Radius: c.Radius,
           Center: &Point{
               X: c.Center.X,
               Y: c.Center.Y,
           },
       }
   }
   ```

### 原型管理器

`ShapeCache` 类用于存储和管理原型对象，它提供了一个中央化的机制来管理不同类型的形状原型：

```go
type ShapeCache struct {
    shapes map[string]Shape
    mu     sync.RWMutex // 用于线程安全
}
```

核心功能包括：
- 添加形状原型到缓存
- 从缓存中获取形状的克隆
- 预加载常用形状

## 使用示例

### 基本使用

```go
// 创建原始形状
circle := NewCircle(10, 0, 0)
circle.SetColor(Red)

// 通过克隆创建新形状
clonedCircle := circle.Clone().(*Circle)
deepClonedCircle := circle.DeepClone().(*Circle)

// 修改原始形状不会影响深克隆
circle.Radius = 20
circle.Center.X = 15
fmt.Println(circle)         // 半径=20, 中心=(15,0)
fmt.Println(clonedCircle)   // 半径=10, 中心=(15,0) - 注意中心点改变了，因为是浅克隆
fmt.Println(deepClonedCircle) // 半径=10, 中心=(0,0) - 完全独立
```

### 使用原型管理器

```go
// 创建并初始化形状缓存
cache := NewShapeCache()
cache.LoadCache()

// 从缓存获取预定义的形状
redCircle := cache.Get("redCircle").(*Circle)
rectangle := cache.Get("rectangle").(*Rectangle)
triangle := cache.Get("triangle").(*Triangle)

// 使用获取的形状
fmt.Println(redCircle.GetArea())
fmt.Println(rectangle)
fmt.Println(triangle)
```

## 优点

1. **避免子类泛滥**: 原型模式让你能够复制现有对象，而无需创建新的子类。
2. **减少重复的初始化代码**: 通过克隆预初始化的对象，可以避免复杂的对象创建过程。
3. **动态配置对象**: 你可以在运行时决定需要实例化哪些类。
4. **构造复杂对象更简单**: 对于那些构造过程复杂的对象，克隆可能比通过构造函数创建更简单。

## 缺点

1. **克隆复杂对象困难**: 对于包含循环引用或私有字段的复杂对象，实现克隆可能会很困难。
2. **深克隆与浅克隆的选择**: 需要慎重选择使用深克隆还是浅克隆，不同的选择可能导致不同的行为。

## 适用场景

1. **创建成本高的对象**: 当对象的创建成本高昂或复杂时，使用原型模式可以提高性能。
2. **需要保存对象状态**: 当需要保存对象的状态并在以后创建相似对象时。
3. **避免构造函数限制**: 某些语言或框架中，构造函数可能有限制；原型模式可以绕过这些限制。
4. **减少子类数量**: 当系统需要独立于它所创建的产品时，可以使用原型模式而不是工厂模式。

## 与其他模式的关系

1. **抽象工厂模式**: 抽象工厂模式通常使用原型模式来实现其工厂。
2. **命令模式**: 可以使用原型模式保存命令的历史记录。
3. **组合模式**: 在组合模式中，克隆可用于复制复杂的组件结构。
4. **备忘录模式**: 原型模式可用于实现备忘录模式，保存对象的状态快照。

## 结语

原型模式是一种简单但强大的设计模式，特别适合需要创建许多相似对象的场景。通过将克隆操作封装在原型对象中，客户端代码可以更加简洁，系统也变得更加灵活。在我们的几何形状示例中，原型模式使得创建不同类型和属性的形状变得更加简单和统一。