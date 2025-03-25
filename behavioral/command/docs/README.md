# 命令模式 (Command Pattern)

## 简介

命令模式是一种行为型设计模式，它将请求封装为一个对象，从而使你可以用不同的请求对客户进行参数化，对请求排队或记录请求日志，以及支持可撤销的操作。

命令模式将"请求"封装成对象，以便使用不同的请求、队列或者日志来参数化其他对象，同时支持可撤销的操作。

## 结构

命令模式包含以下几个核心角色：

1. **命令接口 (Command)** - 声明执行操作的方法
2. **具体命令 (Concrete Command)** - 实现命令接口，负责连接接收者和动作
3. **调用者 (Invoker)** - 要求命令执行请求
4. **接收者 (Receiver)** - 知道如何实施与执行一个请求相关的操作
5. **客户端 (Client)** - 创建具体命令对象并设置其接收者

## 优点

- 将请求的发送者和接收者解耦
- 可以将命令对象存储在队列中
- 可以方便地实现撤销和重做功能
- 可以组合命令创建复合命令
- 符合开闭原则，可以很容易地增加新命令

## 缺点

- 可能导致系统中的类数量激增
- 每个命令都是一个类，可能会导致代码量增加

## 适用场景

- 需要抽象出待执行的动作，然后以参数的形式提供出来
- 需要在不同的时刻指定、排列和执行请求
- 需要支持撤销操作
- 需要支持事务操作

## 实现说明

本实现展示了一个家庭自动化系统的命令模式，包括：

1. 基本的命令接口和具体命令类
2. 设备接口和具体设备（灯和电视）
3. 带撤销功能的命令
4. 宏命令（组合多个命令）
5. 遥控器作为调用者，带有历史记录和撤销功能

### 命令接口

```go
// Command 接口定义了命令的执行和撤销方法
type Command interface {
    Execute() error
    Undo() error
    Name() string
}
```

### 设备接口

```go
// Device 表示可以接收命令的设备接口
type Device interface {
    On() error
    Off() error
    GetName() string
}
```

### 具体命令

```go
// TurnOnCommand 表示开启设备命令
type TurnOnCommand struct {
    device Device
}

// Execute 执行开启命令
func (c *TurnOnCommand) Execute() error {
    return c.device.On()
}

// Undo 撤销开启命令
func (c *TurnOnCommand) Undo() error {
    return c.device.Off()
}
```

### 宏命令

```go
// MacroCommand 表示宏命令，可以执行多个命令
type MacroCommand struct {
    name     string
    commands []Command
}

// Execute 执行所有命令
func (m *MacroCommand) Execute() error {
    for _, cmd := range m.commands {
        if err := cmd.Execute(); err != nil {
            return fmt.Errorf("执行宏命令 %s 时出错: %s 失败: %w", m.name, cmd.Name(), err)
        }
    }
    return nil
}
```

### 调用者（遥控器）

```go
// RemoteControl 表示命令调用者（遥控器）
type RemoteControl struct {
    onCommands    []Command
    offCommands   []Command
    history       []Command
    maxHistoryLen int
}

// OnButtonPressed 按下开启按钮
func (r *RemoteControl) OnButtonPressed(slot int) error {
    if slot < 0 || slot >= len(r.onCommands) {
        return fmt.Errorf("无效的插槽编号: %d", slot)
    }

    cmd := r.onCommands[slot]
    err := cmd.Execute()
    if err == nil {
        r.addToHistory(cmd)
    }
    return err
}
```

## 代码示例

### 基本用法

```go
// 创建设备
livingRoomLight := NewLight("客厅灯")
tv := NewTV("客厅电视")

// 创建命令
livingRoomLightOn := NewTurnOnCommand(livingRoomLight)
livingRoomLightOff := NewTurnOffCommand(livingRoomLight)
tvOn := NewTurnOnCommand(tv)
tvOff := NewTurnOffCommand(tv)

// 创建遥控器
remote := NewRemoteControl(2)
remote.SetCommand(0, livingRoomLightOn, livingRoomLightOff)
remote.SetCommand(1, tvOn, tvOff)

// 使用遥控器控制设备
remote.OnButtonPressed(0)  // 打开客厅灯
remote.OnButtonPressed(1)  // 打开电视
remote.OffButtonPressed(0) // 关闭客厅灯

// 撤销最后一个命令
remote.UndoLastCommand()   // 重新打开客厅灯
```

### 宏命令

```go
// 创建"回家"宏命令
comeHomeCommands := []Command{
    NewTurnOnCommand(livingRoomLight),
    NewTurnOnCommand(kitchenLight),
    NewTurnOnCommand(tv),
}
comeHomeMacro := NewMacroCommand("回家模式", comeHomeCommands)

// 执行宏命令
comeHomeMacro.Execute()   // 同时打开客厅灯、厨房灯和电视

// 撤销宏命令
comeHomeMacro.Undo()      // 按相反顺序关闭电视、厨房灯、客厅灯
```

## 测试说明

测试用例覆盖了以下几个方面：

1. 基本的命令执行和撤销功能
2. 各种设备的特定命令
3. 宏命令的执行和撤销
4. 遥控器的按钮控制和历史记录管理
5. 复杂的家庭自动化场景

可以使用以下命令运行测试：

```bash
cd /path/to/go-design-pattern/behavioral/command
go test -v
```

### 测试中的宏命令示例

```go
// 创建"晚上回家"宏命令
eveningCommands := []Command{
    NewTurnOnCommand(livingRoomLight),
    NewTurnOnCommand(kitchenLight),
    NewTurnOnCommand(tv),
}
eveningMacro := NewMacroCommand("晚上回家", eveningCommands)

// 测试"晚上回家"宏命令执行
output := captureOutput(func() {
    err := eveningMacro.Execute()
    assert.NoError(t, err)
})
assert.Contains(t, output, "客厅灯 已打开")
assert.Contains(t, output, "厨房灯 已打开")
assert.Contains(t, output, "客厅电视 已打开")
```

## 命令模式与其他模式的关系

- **命令模式和策略模式**：两者都可以参数化对象的行为，但策略模式通常只有一个上下文和许多策略，而命令模式为每个操作创建一个命令对象。
- **命令模式和备忘录模式**：可以结合使用来实现撤销功能。
- **命令模式和责任链模式**：可以将命令链接成一个链来处理请求。

## 总结

命令模式通过将请求封装为对象，实现了请求发送者和接收者的解耦。在本实现中，我们展示了命令模式如何用于家庭自动化系统，通过遥控器控制各种家用设备，并支持宏命令和撤销操作。

命令模式的关键在于封装"调用操作的请求"，使调用者与实现解耦，并且可以对请求进行排队、日志和撤销操作。