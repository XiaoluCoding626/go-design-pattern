# 中介者模式（Mediator Pattern）

## 1. 概述

中介者模式是一种行为型设计模式，它通过提供一个中央控制对象（中介者）来管理对象间的交互，从而减少了对象之间的直接引用，降低了系统组件的耦合度。

这种模式特别适合于组件之间交互复杂、依赖关系错综复杂的系统，通过将复杂的交互逻辑集中到中介者对象中，使各个组件能够专注于自身的业务逻辑，而不必关心与其他组件的交互细节。

## 2. 核心组件

### 2.1 中介者（Mediator）

- **定义**：中介者接口定义了组件之间交互的协议。
- **职责**：协调各个组件之间的通信，处理消息路由和分发。

### 2.2 具体中介者（Concrete Mediator）

- **定义**：实现中介者接口的具体类。
- **职责**：维护组件的注册信息，实现消息的分发逻辑。

### 2.3 参与者（Colleague）

- **定义**：通过中介者与其他组件进行交互的对象。
- **职责**：发送和接收消息，但不直接与其他参与者交互。

### 2.4 具体参与者（Concrete Colleague）

- **定义**：实现参与者接口的具体类。
- **职责**：执行自身的业务逻辑，通过中介者与其他参与者交互。

## 3. 实现概述

本实现是一个模拟聊天室的中介者模式示例，具有以下特点：

### 3.1 消息系统

- **消息类型**：支持文本消息、命令消息和通知消息。
- **消息结构**：包含类型、内容、发送者、接收者和时间戳。
- **消息路由**：支持广播和定向消息。

### 3.2 参与者类型

- **用户（User）**：表示人类用户，可以发送和接收各类消息。
- **机器人（Bot）**：自动化参与者，可以识别和响应特定命令。
- **基础参与者**：提供通用功能，减少代码重复。

### 3.3 聊天室功能

- **参与者注册/注销**：动态管理聊天室成员。
- **消息分发**：根据消息类型和目标接收者路由消息。
- **错误处理**：处理参与者不存在等异常情况。

## 4. 实现详解

### 4.1 消息和消息类型

```go
// MessageType 定义可交换的不同消息类型
type MessageType int

const (
    TextMessage         MessageType = iota // 文本消息
    CommandMessage                         // 命令消息
    NotificationMessage                    // 通知消息
)

// Message 表示带有元数据的通信对象
type Message struct {
    Type      MessageType // 消息类型
    Content   string      // 消息内容
    Sender    string      // 发送者ID
    Recipient string      // 接收者ID（空字符串表示广播给所有人）
    Timestamp time.Time   // 时间戳
}
```

### 4.2 中介者接口和实现

```go
// Mediator 定义通信协调的接口
type Mediator interface {
    Register(colleague Colleague)   // 注册参与者
    Unregister(colleague Colleague) // 注销参与者
    Send(message Message)           // 发送消息
}

// ChatRoom 是实现 Mediator 接口的具体中介者
type ChatRoom struct {
    name       string               // 聊天室名称
    colleagues map[string]Colleague // 参与者映射表
}
```

### 4.3 参与者接口和实现

```go
// Colleague 定义通过中介者通信的参与者的接口
type Colleague interface {
    GetID() string                                                  // 获取ID
    GetName() string                                                // 获取名称
    Send(content string, messageType MessageType, recipient string) // 发送消息
    Receive(message Message)                                        // 接收消息
    SetMediator(mediator Mediator)                                  // 设置中介者
}
```

## 5. 使用示例

### 5.1 创建并设置中介者和参与者

```go
// 创建聊天室（中介者）
chatRoom := NewChatRoom("设计模式讨论组")

// 创建用户（参与者）
alice := NewUser("u1", "爱丽丝", "管理员")
bob := NewUser("u2", "鲍勃", "开发者")

// 创建机器人参与者
helpBot := NewBot("b1", "帮助机器人", "!")

// 在中介者中注册
chatRoom.Register(alice)
chatRoom.Register(bob)
chatRoom.Register(helpBot)

// 为每个参与者设置中介者
alice.SetMediator(chatRoom)
bob.SetMediator(chatRoom)
helpBot.SetMediator(chatRoom)
```

### 5.2 消息交互

```go
// 广播消息
alice.Send("大家好！", TextMessage, "")

// 直接消息
bob.Send("嗨，爱丽丝，收到你的消息了", TextMessage, "u1")

// 命令消息
charlie.Send("!help", CommandMessage, "b1")

// 通知
helpBot.Send("系统每日备份已安排", NotificationMessage, "")
```

### 5.3 动态管理参与者

```go
// 注销参与者
chatRoom.Unregister(charlie)

// 新参与者加入
newUser := NewUser("u3", "新用户", "访客")
chatRoom.Register(newUser)
newUser.SetMediator(chatRoom)
```

## 6. 优势和适用场景

### 6.1 优势

1. **降低耦合度**：参与者无需直接引用其他参与者，只需与中介者交互。
2. **集中控制**：交互逻辑集中在中介者中，便于管理和修改。
3. **简化组件**：参与者的实现更加简单，职责更加单一。
4. **提高可扩展性**：添加新的参与者不需要修改现有代码，只需注册到中介者。

### 6.2 适用场景

1. **组件之间存在复杂的通信**：当系统中多个组件需要以不同方式进行通信。
2. **希望组件之间松散耦合**：减少组件之间的直接依赖，提高系统的可维护性。
3. **需要集中管理交互规则**：当组件之间的交互规则复杂并且可能发生变化。
4. **通信逻辑需要复用**：跨多个组件的通信逻辑可以集中在中介者中重用。

### 6.3 实际应用示例

1. **聊天应用**：如本例所示，管理用户和消息的路由。
2. **GUI框架**：协调控件之间的交互，如按钮点击影响文本框内容。
3. **航空管制系统**：塔台（中介者）协调多架飞机（参与者）的起降。
4. **多人游戏**：游戏服务器作为中介者协调玩家之间的交互。
5. **中央事件总线**：在前端框架中管理组件间的事件通信。

## 7. 注意事项

1. **避免中介者过于复杂**：随着系统扩大，中介者可能变得过于庞大，此时可考虑划分多个中介者。
2. **性能考虑**：所有通信都经过中介者可能造成性能瓶颈，需权衡设计。
3. **错误处理**：中介者应妥善处理通信异常，如接收者不存在的情况。
4. **线程安全**：在并发环境中使用中介者需考虑同步问题。

## 8. 总结

中介者模式通过引入一个中心化的协调对象，有效解决了对象之间的复杂依赖关系问题。本实现通过一个聊天室系统展示了中介者模式在消息传递系统中的应用，包括不同类型的消息、多种参与者角色以及灵活的消息路由机制。

这种模式特别适合于通信复杂的系统，可以显著降低系统的耦合度，提高可维护性和可扩展性。然而，随着系统规模的扩大，也需注意避免中介者本身变得过于复杂和臃肿。