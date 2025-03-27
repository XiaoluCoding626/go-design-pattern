# 装饰器模式（Decorator Pattern）

## 介绍

装饰器模式是一种结构型设计模式，允许你通过将对象放入包含行为的特殊封装对象中来为原对象动态添加新的行为。装饰器提供了子类继承的替代方案，可以在不改变原始类（或其客户代码）的情况下，动态地扩展对象的功能。

## 意图

装饰器模式的主要意图是：
- 在不修改现有对象结构的情况下动态添加职责
- 比继承更灵活地扩展功能
- 将对象功能分解成不同的装饰层，按需组合

## 结构

本实现中的装饰器模式包含以下核心元素：

### 组件结构
1. **Component（组件接口）**：定义被装饰对象和装饰器共有的接口
2. **ConcreteComponent（具体组件）**：需要被装饰的原始对象，实现Component接口
3. **BaseDecorator（基础装饰器）**：所有装饰器的共同父类，实现Component接口并包含对另一个Component的引用
4. **具体装饰器**：添加特定功能的装饰器，如化妆装饰器和配饰装饰器

### UML类图

```
              ┌───────────────┐
              │  «interface»  │
              │   Component   │
              └───────┬───────┘
                      │
                      │ implements
          ┌───────────┼───────────┐
          │           │           │
┌─────────▼─────────┐ │ ┌─────────▼─────────┐
│ ConcreteComponent │ │ │   BaseDecorator   │◄─────┐
└───────────────────┘ │ └─────────┬─────────┘     │
                      │           │               │ has-a
                      │           │ extends       │
                      │ ┌─────────┴─────────┐     │
                      │ │                   │     │
              ┌───────┴─┴───────┐ ┌─────────▼─────────┐
              │ MakeupDecorator │ │ AccessoryDecorator │
              └────────┬────────┘ └─────────┬─────────┘
                       │                    │
       ┌───────────────┼──────────────┐    │
       │               │              │    │
┌──────▼──────┐ ┌──────▼──────┐ ┌────▼────┐ ┌────▼────┐
│ Foundation  │ │  Lipstick   │ │ Eyeshadow│ │Necklace │ ...
│ Decorator   │ │ Decorator   │ │Decorator │ │Decorator│
└─────────────┘ └─────────────┘ └─────────┘ └─────────┘
```

## 实现

本实现采用了多层次的装饰器模式：

1. **基础接口和类**
   - `Component` 接口定义了所有组件的共同方法 `Show()`
   - `ConcreteComponent` 是被装饰的基础组件
   - `BaseDecorator` 是所有装饰器的基类

2. **分类装饰器**
   - `MakeupDecorator` 是化妆类装饰器的抽象基类，使用"【】"包裹的格式
   - `AccessoryDecorator` 是配饰类装饰器的抽象基类，使用"+"连接的格式

3. **具体装饰器**
   - 化妆类：`FoundationDecorator`、`LipstickDecorator`、`EyeshadowDecorator`
   - 配饰类：`NecklaceDecorator`、`EarringsDecorator`

## 代码示例

### 创建基础组件

```go
// 创建一个基础组件
person := NewConcreteComponent("素颜")
```

### 使用单一装饰器

```go
// 应用粉底装饰
withFoundation := NewFoundationDecorator(person)
fmt.Println(withFoundation.Show())  // 输出: "打粉底【素颜】"

// 应用耳环装饰
withEarrings := NewEarringsDecorator(person)
fmt.Println(withEarrings.Show())  // 输出: "素颜 + 耳环"
```

### 组合多个装饰器

```go
// 应用多层装饰
fullMakeup := NewLipstickDecorator(
    NewEyeshadowDecorator(
        NewFoundationDecorator(person),
    ),
)
fmt.Println(fullMakeup.Show())  // 输出: "涂口红【画眼影【打粉底【素颜】】】"

// 混合不同类型的装饰器
withMakeupAndAccessories := NewNecklaceDecorator(
    NewEarringsDecorator(
        NewLipstickDecorator(person),
    ),
)
fmt.Println(withMakeupAndAccessories.Show())  // 输出: "涂口红【素颜】 + 耳环 + 项链"
```

## 优点

1. **比静态继承更灵活**：可以动态组合多个装饰器，实现功能叠加
2. **遵循单一职责原则**：将不同的功能分散到不同的类中
3. **运行时添加/移除职责**：无需修改现有代码即可扩展功能
4. **避免功能齐全的超类**：通过组合小型类实现复杂功能

## 缺点

1. **难以移除特定装饰器**：必须重建整个装饰器栈
2. **实现顺序很重要**：装饰器嵌套顺序不同会产生不同结果
3. **初始化代码复杂**：嵌套多个装饰器的代码可读性较低

## 应用场景

1. **需要动态透明地添加职责**：不想通过继承生成很多子类
2. **功能需要按特定顺序组合**：装饰器允许精确控制功能添加顺序
3. **核心功能简单但需要多种额外功能**：通过不同装饰器组合实现
4. **可插拔式功能扩展**：允许用户选择需要的功能组合

## 与其他模式的关系

- **适配器模式**：改变接口，而装饰器模式不改变接口只增强功能
- **组合模式**：装饰器像组合模式但通常只有一个子组件
- **策略模式**：装饰器改变对象外表，策略模式改变内部逻辑
- **代理模式**：代理控制对对象的访问，装饰器为对象添加功能

## 实际应用

- Java I/O 流库大量使用装饰器模式
- Web应用中的请求处理中间件
- GUI组件功能增强
- 缓存、日志、事务等横切关注点的实现