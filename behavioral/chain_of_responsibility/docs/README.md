# 责任链模式 (Chain of Responsibility Pattern)

## 简介

责任链模式是一种行为设计模式，它通过为请求创建一个接收者对象的链，使多个对象都有机会处理请求。沿着这条链传递请求，直到有一个对象能够处理它为止。

责任链模式将请求的发送者与接收者解耦，使得多个对象都有机会处理请求，而不需要显式地指定接收者。

## 结构

![责任链模式结构](https://refactoringguru.cn/images/patterns/diagrams/chain-of-responsibility/structure-2x.png)

责任链模式主要有以下几个角色：

1. **处理者接口/抽象类 (Handler)**: 定义一个处理请求的接口，通常包含设置继任者和处理请求的方法。
2. **具体处理者 (Concrete Handler)**: 实现处理者接口，处理它负责的请求并可以访问它的后继者。
3. **客户端 (Client)**: 向链上的具体处理者对象提交请求。

## 代码实现

本项目实现了一个基于金额审批的责任链模式示例。

### 核心组件

#### 1. 处理者接口

```go
type IApprover interface {
    // SetNext 设置下一个审批人
    SetNext(approver IApprover) IApprover
    // Approve 批准金额
    Approve(amount float64) ApprovalResult
    // GetName 获取审批人名称
    GetName() string
}
```

#### 2. 审批结果结构体

```go
type ApprovalResult struct {
    Approved bool   // 是否批准
    Approver string // 审批人
    Message  string // 审批消息
}
```

#### 3. 基础处理者类

```go
type BaseApprover struct {
    name  string
    limit float64
    next  IApprover
}
```

#### 4. 具体处理者类

- **经理 (Manager)**: 可审批1000元以下的请求
- **总监 (Director)**: 可审批5000元以下的请求
- **CFO**: 可审批20000元以下的请求

### 工作流程

1. 客户端创建一个责任链，例如：经理 -> 总监 -> CFO
2. 客户端将请求（金额）提交到链的第一个处理者（经理）
3. 经理检查金额是否在自己的权限范围内：
   - 如果是，处理请求并返回结果
   - 如果不是，将请求传递给下一个处理者（总监）
4. 总监重复相同的检查流程
5. CFO作为链中的最后一个处理者，如果金额也超出其权限范围，则拒绝请求

## 使用示例

### 创建并使用标准责任链

```go
// 创建责任链
chain := CreateApprovalChain()

// 提交不同金额的请求
result1 := chain.Approve(500)  // 由经理处理
result2 := chain.Approve(3000) // 由总监处理
result3 := chain.Approve(12000) // 由CFO处理
result4 := chain.Approve(30000) // 超出所有处理者的权限，被拒绝

fmt.Println(result1.Message) // 输出: 经理批准了 500.00 元的请求
```

### 创建自定义责任链

责任链的顺序可以根据需要灵活设置：

```go
// 创建一个具有特殊顺序的责任链：CFO -> 经理 -> 总监
cfo := NewCFO(20000)
manager := NewManager(1000)
director := NewDirector(5000)

cfo.SetNext(manager).SetNext(director)

// 提交请求
result := cfo.Approve(3000)
```

## 优点

1. **减少耦合度**: 发送者与接收者之间的耦合度降低。
2. **灵活性增强**: 可以动态地改变处理请求的顺序或增减处理者。
3. **符合开闭原则**: 可以在不修改现有代码的情况下增加新的处理者。
4. **单一职责原则**: 每个处理者只需关心自己能处理的请求类型。

## 缺点

1. **不保证被处理**: 如果没有任何处理者能处理请求，请求可能会被丢弃。
2. **调试复杂**: 由于请求在链中的传递路径不明确，调试可能会变得困难。
3. **延迟增加**: 请求需要经过多个对象，可能导致处理延迟增加。

## 适用场景

1. **多个对象可以处理同一请求，但处理程序和顺序在运行时才能确定**。
2. **需要按照特定顺序向多个处理者发送请求**。
3. **处理者集合需要动态变化**。

## 总结

责任链模式是一种强大的行为设计模式，它允许我们构建一系列可以处理请求的对象，而不必显式指定接收者。这种模式在处理具有不同权限级别的系统（如审批流程）中特别有用，使系统更加灵活和可扩展。

本项目中的责任链实现演示了如何处理具有不同权限级别的金额审批场景，展示了模式的核心优势：灵活的请求处理流程和良好的解耦。