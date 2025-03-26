package mediator

import (
	"fmt"
	"time"
)

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

// NewChatRoom 创建一个新的聊天室中介者
func NewChatRoom(name string) *ChatRoom {
	return &ChatRoom{
		name:       name,
		colleagues: make(map[string]Colleague),
	}
}

// Register 将参与者添加到中介者的注册表中
func (c *ChatRoom) Register(colleague Colleague) {
	c.colleagues[colleague.GetID()] = colleague
	fmt.Printf("[%s] %s 已加入聊天室\n", c.name, colleague.GetName())
}

// Unregister 从中介者的注册表中移除参与者
func (c *ChatRoom) Unregister(colleague Colleague) {
	if _, exists := c.colleagues[colleague.GetID()]; exists {
		delete(c.colleagues, colleague.GetID())
		fmt.Printf("[%s] %s 已离开聊天室\n", c.name, colleague.GetName())
	}
}

// Send 将消息分发给适当的接收者
func (c *ChatRoom) Send(message Message) {
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	// 记录消息
	switch message.Type {
	case TextMessage:
		fmt.Printf("[%s] 来自 %s 的消息: %s\n", c.name, message.Sender, message.Content)
	case CommandMessage:
		fmt.Printf("[%s] 来自 %s 的命令: %s\n", c.name, message.Sender, message.Content)
	case NotificationMessage:
		fmt.Printf("[%s] 通知: %s\n", c.name, message.Content)
	}

	// 将消息发送给适当的接收者
	if message.Recipient != "" {
		// 发送直接消息给特定接收者
		if recipient, exists := c.colleagues[message.Recipient]; exists {
			recipient.Receive(message)
		} else {
			fmt.Printf("[%s] 错误: 接收者 %s 未找到\n", c.name, message.Recipient)
		}
	} else {
		// 广播消息给除发送者外的所有参与者
		for id, colleague := range c.colleagues {
			if id != message.Sender {
				colleague.Receive(message)
			}
		}
	}
}

// Colleague 定义通过中介者通信的参与者的接口
type Colleague interface {
	GetID() string                                                  // 获取ID
	GetName() string                                                // 获取名称
	Send(content string, messageType MessageType, recipient string) // 发送消息
	Receive(message Message)                                        // 接收消息
	SetMediator(mediator Mediator)                                  // 设置中介者
}

// BaseColleague 为不同参与者类型提供通用功能
type BaseColleague struct {
	id       string   // 唯一标识符
	name     string   // 名称
	mediator Mediator // 中介者引用
}

// User 是表示聊天用户的具体参与者
type User struct {
	BaseColleague
	role string // 用户角色
}

// NewUser 创建一个新的用户参与者
func NewUser(id string, name string, role string) *User {
	return &User{
		BaseColleague: BaseColleague{
			id:   id,
			name: name,
		},
		role: role,
	}
}

// GetID 返回用户的唯一标识符
func (u *User) GetID() string {
	return u.id
}

// GetName 返回用户的名称
func (u *User) GetName() string {
	return u.name
}

// Send 创建消息并通过中介者发送
func (u *User) Send(content string, messageType MessageType, recipient string) {
	if u.mediator == nil {
		fmt.Printf("错误: %s 没有中介者，无法发送消息\n", u.name)
		return
	}

	message := Message{
		Type:      messageType,
		Content:   content,
		Sender:    u.id,
		Recipient: recipient,
		Timestamp: time.Now(),
	}

	u.mediator.Send(message)
}

// Receive 处理接收到的消息
func (u *User) Receive(message Message) {
	switch message.Type {
	case TextMessage:
		fmt.Printf("[%s (%s)] 收到来自 %s 的消息: %s\n",
			u.name, u.role, message.Sender, message.Content)
	case CommandMessage:
		fmt.Printf("[%s (%s)] 收到来自 %s 的命令: %s\n",
			u.name, u.role, message.Sender, message.Content)
	case NotificationMessage:
		fmt.Printf("[%s (%s)] 收到通知: %s\n",
			u.name, u.role, message.Content)
	}
}

// SetMediator 为该参与者设置中介者
func (u *User) SetMediator(mediator Mediator) {
	u.mediator = mediator
}

// Bot 是另一种具体参与者，表示自动化参与者
type Bot struct {
	BaseColleague
	commandPrefix string // 命令前缀
}

// NewBot 创建一个新的机器人参与者
func NewBot(id string, name string, commandPrefix string) *Bot {
	return &Bot{
		BaseColleague: BaseColleague{
			id:   id,
			name: name,
		},
		commandPrefix: commandPrefix,
	}
}

// GetID 返回机器人的唯一标识符
func (b *Bot) GetID() string {
	return b.id
}

// GetName 返回机器人的名称
func (b *Bot) GetName() string {
	return b.name
}

// Send 创建消息并通过中介者发送
func (b *Bot) Send(content string, messageType MessageType, recipient string) {
	if b.mediator == nil {
		fmt.Printf("错误: %s 没有中介者，无法发送消息\n", b.name)
		return
	}

	// 机器人通常发送命令或通知
	if messageType == TextMessage {
		messageType = NotificationMessage
	}

	message := Message{
		Type:      messageType,
		Content:   content,
		Sender:    b.id,
		Recipient: recipient,
		Timestamp: time.Now(),
	}

	b.mediator.Send(message)
}

// Receive 处理接收到的消息并自动响应命令
func (b *Bot) Receive(message Message) {
	// 机器人可以响应命令
	if message.Type == CommandMessage && len(message.Content) > 0 {
		if message.Content[0] == b.commandPrefix[0] {
			response := fmt.Sprintf("正在处理命令: %s", message.Content)
			b.Send(response, NotificationMessage, message.Sender)
		}
	} else if message.Type == TextMessage {
		fmt.Printf("[%s (机器人)] 收到来自 %s 的消息: %s\n",
			b.name, message.Sender, message.Content)
	}
}

// SetMediator 为该参与者设置中介者
func (b *Bot) SetMediator(mediator Mediator) {
	b.mediator = mediator
}
