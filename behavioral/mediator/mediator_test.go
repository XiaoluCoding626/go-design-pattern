package mediator

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试辅助结构 - 用于捕获输出消息
type MessageCollector struct {
	BaseColleague
	receivedMessages []Message
}

func NewMessageCollector(id string, name string) *MessageCollector {
	return &MessageCollector{
		BaseColleague: BaseColleague{
			id:   id,
			name: name,
		},
		receivedMessages: make([]Message, 0),
	}
}

func (mc *MessageCollector) GetID() string {
	return mc.id
}

func (mc *MessageCollector) GetName() string {
	return mc.name
}

func (mc *MessageCollector) Send(content string, messageType MessageType, recipient string) {
	if mc.mediator == nil {
		return
	}
	message := Message{
		Type:      messageType,
		Content:   content,
		Sender:    mc.id,
		Recipient: recipient,
		Timestamp: time.Now(),
	}
	mc.mediator.Send(message)
}

func (mc *MessageCollector) Receive(message Message) {
	mc.receivedMessages = append(mc.receivedMessages, message)
}

func (mc *MessageCollector) SetMediator(mediator Mediator) {
	mc.mediator = mediator
}

func (mc *MessageCollector) GetMessages() []Message {
	return mc.receivedMessages
}

func (mc *MessageCollector) CountMessagesOfType(messageType MessageType) int {
	count := 0
	for _, msg := range mc.receivedMessages {
		if msg.Type == messageType {
			count++
		}
	}
	return count
}

func (mc *MessageCollector) HasMessageFrom(sender string) bool {
	for _, msg := range mc.receivedMessages {
		if msg.Sender == sender {
			return true
		}
	}
	return false
}

// 测试用例
func TestMediatorBasicFunctionality(t *testing.T) {
	// 创建聊天室（中介者）
	chatRoom := NewChatRoom("设计模式讨论组")

	// 创建用户（参与者）
	alice := NewUser("u1", "爱丽丝", "管理员")
	bob := NewUser("u2", "鲍勃", "开发者")
	charlie := NewUser("u3", "查理", "测试员")

	// 创建机器人参与者
	helpBot := NewBot("b1", "帮助机器人", "!")

	// 在中介者中注册
	chatRoom.Register(alice)
	chatRoom.Register(bob)
	chatRoom.Register(charlie)
	chatRoom.Register(helpBot)

	// 为每个参与者设置中介者
	alice.SetMediator(chatRoom)
	bob.SetMediator(chatRoom)
	charlie.SetMediator(chatRoom)
	helpBot.SetMediator(chatRoom)

	// 添加消息收集器用于测试断言
	collector := NewMessageCollector("collector", "消息收集器")
	chatRoom.Register(collector)
	collector.SetMediator(chatRoom)

	// 测试广播消息
	alice.Send("大家好！", TextMessage, "")
	time.Sleep(50 * time.Millisecond)

	// 断言收集器应该收到消息
	assert.True(t, collector.HasMessageFrom("u1"), "收集器应该收到来自爱丽丝的消息")
	assert.Equal(t, 1, collector.CountMessagesOfType(TextMessage), "应该收到一条文本消息")

	// 测试直接消息
	bob.Send("嗨，爱丽丝，收到你的消息了", TextMessage, "u1")
	time.Sleep(50 * time.Millisecond)

	// 收集器不应该收到直接消息
	assert.Equal(t, 1, collector.CountMessagesOfType(TextMessage), "不应该收到直接消息")

	// 测试命令消息
	charlie.Send("!help", CommandMessage, "b1")
	time.Sleep(50 * time.Millisecond)

	// 测试通知
	helpBot.Send("系统每日备份已安排", NotificationMessage, "")
	time.Sleep(50 * time.Millisecond)

	// 断言收集器收到通知
	assert.Equal(t, 1, collector.CountMessagesOfType(NotificationMessage), "应该收到一条通知消息")

	// 测试注销
	chatRoom.Unregister(charlie)
	time.Sleep(50 * time.Millisecond)

	// 参与者离开后的消息
	alice.Send("查理离开了讨论组", TextMessage, "")
	time.Sleep(50 * time.Millisecond)

	// 断言消息数量增加
	assert.Equal(t, 2, collector.CountMessagesOfType(TextMessage), "应该收到两条文本消息")
}

// 测试错误处理和边界情况
func TestMediatorErrorHandling(t *testing.T) {
	chatRoom := NewChatRoom("错误处理测试")

	alice := NewUser("u1", "爱丽丝", "管理员")
	bob := NewUser("u2", "鲍勃", "开发者")

	chatRoom.Register(alice)
	alice.SetMediator(chatRoom)

	// 测试发送消息到不存在的接收者
	alice.Send("你好，不存在的用户", TextMessage, "non-existent")
	time.Sleep(50 * time.Millisecond)

	// bob未设置中介者，应该会报错
	bob.Send("没有中介者", TextMessage, "u1")
	time.Sleep(50 * time.Millisecond)

	// 注册bob但不设置中介者
	chatRoom.Register(bob)

	// 注销不存在的参与者
	nonExistentColleague := NewUser("non-existent", "不存在", "无")
	chatRoom.Unregister(nonExistentColleague)
	time.Sleep(50 * time.Millisecond)

	// 空消息内容
	alice.Send("", TextMessage, "")
	time.Sleep(50 * time.Millisecond)
}

// 测试机器人的命令响应
func TestBotCommandResponses(t *testing.T) {
	chatRoom := NewChatRoom("机器人命令测试")

	bot := NewBot("b1", "机器人", "!")

	// 创建一个消息收集器来发送命令并捕获机器人的回复
	collector := NewMessageCollector("collector", "收集器")

	chatRoom.Register(bot)
	chatRoom.Register(collector)

	bot.SetMediator(chatRoom)
	collector.SetMediator(chatRoom)

	// 发送无效命令（不使用前缀）
	collector.Send("help", CommandMessage, "b1")
	time.Sleep(50 * time.Millisecond)

	// 发送有效命令 - 现在由收集器发送，这样它就能接收回复
	collector.Send("!help", CommandMessage, "b1")
	time.Sleep(100 * time.Millisecond)

	// 验证机器人是否响应命令
	messageFound := false
	for _, msg := range collector.GetMessages() {
		if msg.Type == NotificationMessage && msg.Sender == "b1" &&
			strings.Contains(msg.Content, "正在处理命令: !help") {
			messageFound = true
			break
		}
	}

	assert.True(t, messageFound, "机器人应该回复命令消息")
}

// 测试复杂交互场景
func TestComplexInteractions(t *testing.T) {
	chatRoom := NewChatRoom("复杂交互测试")

	// 创建多个不同角色的用户
	admin := NewUser("a1", "管理员", "系统管理")
	user1 := NewUser("u1", "用户1", "普通用户")
	user2 := NewUser("u2", "用户2", "普通用户")
	moderator := NewUser("m1", "版主", "内容审核")
	infoBot := NewBot("b1", "信息机器人", "?")

	// 注册并设置中介者
	for _, c := range []Colleague{admin, user1, user2, moderator, infoBot} {
		chatRoom.Register(c)
		c.SetMediator(chatRoom)
	}

	// 模拟复杂的群组交互
	admin.Send("欢迎来到聊天室", NotificationMessage, "")
	time.Sleep(50 * time.Millisecond)

	user1.Send("大家好", TextMessage, "")
	time.Sleep(50 * time.Millisecond)

	user2.Send("?时间", CommandMessage, "b1")
	time.Sleep(50 * time.Millisecond)

	moderator.Send("请遵守聊天规则", TextMessage, "")
	time.Sleep(50 * time.Millisecond)

	// 私聊
	user1.Send("你好，管理员", TextMessage, "a1")
	time.Sleep(50 * time.Millisecond)

	// 注销和重新加入
	chatRoom.Unregister(user2)
	time.Sleep(50 * time.Millisecond)

	// 新用户加入
	newUser := NewUser("u3", "新用户", "访客")
	chatRoom.Register(newUser)
	newUser.SetMediator(chatRoom)
	time.Sleep(50 * time.Millisecond)

	newUser.Send("我是新来的", TextMessage, "")
	time.Sleep(50 * time.Millisecond)

	// 此测试主要检验在复杂交互场景下没有崩溃或异常
}

// 基准测试，评估性能
func BenchmarkMediator(b *testing.B) {
	chatRoom := NewChatRoom("性能测试聊天室")

	// 创建10个用户
	users := make([]*User, 10)
	for i := 0; i < 10; i++ {
		users[i] = NewUser(fmt.Sprintf("u%d", i), fmt.Sprintf("用户%d", i), "测试用户")
		chatRoom.Register(users[i])
		users[i].SetMediator(chatRoom)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每个用户发送一条消息
		senderIndex := i % 10
		users[senderIndex].Send(
			fmt.Sprintf("这是第%d条消息", i),
			TextMessage,
			"", // 广播
		)
	}
}
