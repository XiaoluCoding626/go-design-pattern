# 备忘录模式（Memento Pattern）

## 概述

备忘录模式（Memento Pattern）是一种行为设计模式，它允许在不破坏对象封装性的前提下捕获并恢复对象的内部状态。该模式常用于实现撤销/重做功能或保存对象状态快照。

## 主要组成部分

备忘录模式由以下三个主要角色组成：

1. **发起人（Originator）**：负责创建备忘录并使用备忘录恢复自身状态
   - 在本实现中由 `Document` 结构体担任

2. **备忘录（Memento）**：存储发起人的内部状态
   - 在本实现中由 `Memento` 接口和 `documentMemento` 结构体担任

3. **管理者（Caretaker）**：负责保存备忘录，但不会修改备忘录内容
   - 在本实现中由 `Caretaker` 结构体担任

## 实现原理

### 备忘录（Memento）

我们使用接口和私有结构体实现备忘录，确保只有发起人（Document）可以访问备忘录的完整内容，而其他对象只能持有但不能查看或修改其内容：

```go
// Memento 定义备忘录接口
type Memento interface {
    // GetState 是一个空方法，仅用于限制访问
}

// documentMemento 具体的备忘录实现，保存文档的完整状态
type documentMemento struct {
    title string
    body  string
}
```

### 发起人（Originator）

发起人能够创建备忘录以及从备忘录中恢复状态：

```go
// Document 定义文档结构体
type Document struct {
    title string
    body  string
}

// CreateMemento 创建文档的备忘录
func (d *Document) CreateMemento() Memento {
    return &documentMemento{
        title: d.title,
        body:  d.body,
    }
}

// RestoreFromMemento 从备忘录恢复文档状态
func (d *Document) RestoreFromMemento(m Memento) {
    if memento, ok := m.(*documentMemento); ok {
        d.title = memento.title
        d.body = memento.body
    }
}
```

### 管理者（Caretaker）

管理者负责管理备忘录的历史记录，并提供撤销/重做功能：

```go
// Caretaker 定义备忘录管理者结构体
type Caretaker struct {
    document    *Document
    mementos    []Memento
    currentPos  int
    maxMementos int
}
```

## 功能特点

本实现具有以下功能：

1. **完整的撤销/重做功能**：支持多级操作的撤销和重做
2. **状态分支处理**：当在历史中间状态执行新操作时，自动删除该点之后的历史记录
3. **历史记录限制**：可配置最大历史记录数，防止内存无限增长
4. **完整状态保存**：备忘录保存文档的完整状态（标题和内容）

## 使用示例

```go
// 创建文档
document := NewDocument("设计模式示例")

// 创建备忘录管理者
editor := NewCaretaker(document, 10)

// 进行一系列操作
editor.Append("这是第一行内容。")
editor.Append("这是第二行内容。")

// 撤销操作
editor.Undo()

// 重做操作
editor.Redo()

// 删除文档内容
editor.Delete()

// 撤销删除
editor.Undo()
```

## 优势与适用场景

### 优势

1. **保持封装性**：备忘录模式可在不暴露对象实现细节的情况下保存和恢复对象状态
2. **简化发起人**：将状态存储和恢复的责任转移给备忘录，简化了发起人的结构
3. **维护历史记录**：提供了一种简单的方式来维护对象的历史状态

### 适用场景

1. **需要实现撤销/重做功能**的应用，如文本编辑器、图形编辑器等
2. **需要保存对象状态快照**以便后续恢复的场景
3. **需要在不破坏对象封装的前提下访问对象内部状态**的场景

## 总结

备忘录模式通过保存对象的状态快照，提供了一种优雅的方式来实现撤销/重做功能。在本实现中，我们通过明确区分发起人、备忘录和管理者的职责，确保了设计的清晰性和可维护性。同时，增加了历史记录管理和状态分支处理等实用功能，使得该实现能够满足实际开发中的各种需求。