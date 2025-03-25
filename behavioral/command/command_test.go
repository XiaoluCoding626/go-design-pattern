package command

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// captureOutput 捕获标准输出的辅助函数
func captureOutput(fn func()) string {
	// 保存原始的标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行函数
	fn()

	// 恢复标准输出并获取捕获的内容
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// 测试灯的基本命令
func TestLightCommands(t *testing.T) {
	light := NewLight("客厅灯")

	// 测试开灯命令
	onCommand := NewTurnOnCommand(light)
	output := captureOutput(func() {
		err := onCommand.Execute()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已打开")

	// 测试重复开灯命令应该返回错误
	err := onCommand.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已经是开启状态")

	// 测试关灯命令
	offCommand := NewTurnOffCommand(light)
	output = captureOutput(func() {
		err := offCommand.Execute()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已关闭")

	// 测试撤销命令
	output = captureOutput(func() {
		err := offCommand.Undo()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已打开")
}

// 测试电视的基本命令
func TestTVCommands(t *testing.T) {
	tv := NewTV("客厅电视")

	// 测试开电视命令
	onCommand := NewTurnOnCommand(tv)
	output := captureOutput(func() {
		err := onCommand.Execute()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅电视 已打开")
	assert.Contains(t, output, "音量: 50")

	// 测试关电视命令
	offCommand := NewTurnOffCommand(tv)
	output = captureOutput(func() {
		err := offCommand.Execute()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅电视 已关闭")
}

// 测试灯的亮度设置命令
func TestSetLevelCommand(t *testing.T) {
	light := NewLight("卧室灯")

	// 先开灯
	onCommand := NewTurnOnCommand(light)
	onCommand.Execute()

	// 测试设置亮度命令
	levelCommand := NewSetLevelCommand(light, 50)
	output := captureOutput(func() {
		err := levelCommand.Execute()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "卧室灯 亮度设置为 50%")

	// 测试撤销亮度设置
	output = captureOutput(func() {
		err := levelCommand.Undo()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "卧室灯 亮度设置为 100%")

	// 测试亮度设置超出范围
	invalidLevelCommand := NewSetLevelCommand(light, 150)
	err := invalidLevelCommand.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "亮度必须在0-100之间")
}

// 测试宏命令
func TestMacroCommand(t *testing.T) {
	livingRoomLight := NewLight("客厅灯")
	kitchenLight := NewLight("厨房灯")
	tv := NewTV("客厅电视")

	// 创建"回家"宏命令（开灯、开电视）
	comeHomeCommands := []Command{
		NewTurnOnCommand(livingRoomLight),
		NewTurnOnCommand(kitchenLight),
		NewTurnOnCommand(tv),
	}
	comeHomeMacro := NewMacroCommand("回家模式", comeHomeCommands)

	// 测试"回家"宏命令执行
	output := captureOutput(func() {
		err := comeHomeMacro.Execute()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已打开")
	assert.Contains(t, output, "厨房灯 已打开")
	assert.Contains(t, output, "客厅电视 已打开")

	// 测试宏命令的撤销
	output = captureOutput(func() {
		err := comeHomeMacro.Undo()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅电视 已关闭")
	assert.Contains(t, output, "厨房灯 已关闭")
	assert.Contains(t, output, "客厅灯 已关闭")
}

// 测试遥控器基本功能
func TestRemoteControl(t *testing.T) {
	remote := NewRemoteControl(3)
	livingRoomLight := NewLight("客厅灯")
	kitchenLight := NewLight("厨房灯")
	tv := NewTV("客厅电视")

	// 设置遥控器的按钮
	err := remote.SetCommand(0, NewTurnOnCommand(livingRoomLight), NewTurnOffCommand(livingRoomLight))
	assert.NoError(t, err)

	err = remote.SetCommand(1, NewTurnOnCommand(kitchenLight), NewTurnOffCommand(kitchenLight))
	assert.NoError(t, err)

	err = remote.SetCommand(2, NewTurnOnCommand(tv), NewTurnOffCommand(tv))
	assert.NoError(t, err)

	// 测试遥控器的字符串表示
	remoteStr := remote.String()
	assert.Contains(t, remoteStr, "开启 客厅灯")
	assert.Contains(t, remoteStr, "关闭 客厅灯")

	// 测试按钮操作
	output := captureOutput(func() {
		err := remote.OnButtonPressed(0)
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已打开")

	output = captureOutput(func() {
		err := remote.OffButtonPressed(0)
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已关闭")

	// 测试无效的插槽
	err = remote.OnButtonPressed(10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的插槽编号")
}

// 测试遥控器的历史记录和撤销功能
func TestRemoteControlHistory(t *testing.T) {
	remote := NewRemoteControl(2)
	livingRoomLight := NewLight("客厅灯")
	kitchenLight := NewLight("厨房灯")

	// 设置遥控器的按钮
	remote.SetCommand(0, NewTurnOnCommand(livingRoomLight), NewTurnOffCommand(livingRoomLight))
	remote.SetCommand(1, NewTurnOnCommand(kitchenLight), NewTurnOffCommand(kitchenLight))

	// 执行一系列命令
	remote.OnButtonPressed(0)  // 开客厅灯
	remote.OnButtonPressed(1)  // 开厨房灯
	remote.OffButtonPressed(0) // 关客厅灯

	// 测试历史记录显示
	output := captureOutput(func() {
		remote.ShowHistory()
	})
	assert.Contains(t, output, "开启 客厅灯")
	assert.Contains(t, output, "开启 厨房灯")
	assert.Contains(t, output, "关闭 客厅灯")

	// 测试撤销最后一条命令（关客厅灯）
	output = captureOutput(func() {
		err := remote.UndoLastCommand()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已打开")

	// 再次撤销（开厨房灯）
	output = captureOutput(func() {
		err := remote.UndoLastCommand()
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "厨房灯 已关闭")

	// 历史记录为空时撤销应该返回错误
	remote.UndoLastCommand() // 再撤销一次清空历史
	err := remote.UndoLastCommand()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "没有可撤销的命令")
}

// 测试复杂场景：家庭自动化
func TestHomeAutomation(t *testing.T) {
	remote := NewRemoteControl(4)
	livingRoomLight := NewLight("客厅灯")
	kitchenLight := NewLight("厨房灯")
	bedroomLight := NewLight("卧室灯")
	tv := NewTV("客厅电视")

	// 设置遥控器按钮 0-2 为单个设备控制
	remote.SetCommand(0, NewTurnOnCommand(livingRoomLight), NewTurnOffCommand(livingRoomLight))
	remote.SetCommand(1, NewTurnOnCommand(tv), NewTurnOffCommand(tv))

	// 创建"晚上回家"宏命令
	eveningCommands := []Command{
		NewTurnOnCommand(livingRoomLight),
		NewTurnOnCommand(kitchenLight),
		NewTurnOnCommand(tv),
	}
	eveningMacro := NewMacroCommand("晚上回家", eveningCommands)

	// 创建"睡觉时间"宏命令
	bedtimeCommands := []Command{
		NewTurnOffCommand(livingRoomLight),
		NewTurnOffCommand(kitchenLight),
		NewTurnOffCommand(tv),
		NewTurnOnCommand(bedroomLight),
		NewSetLevelCommand(bedroomLight, 30),
	}
	bedtimeMacro := NewMacroCommand("睡觉时间", bedtimeCommands)

	// 设置宏命令到遥控器
	remote.SetCommand(2, eveningMacro, &NoOpCommand{})
	remote.SetCommand(3, bedtimeMacro, &NoOpCommand{})

	// 执行"晚上回家"场景
	output := captureOutput(func() {
		err := remote.OnButtonPressed(2)
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已打开")
	assert.Contains(t, output, "厨房灯 已打开")
	assert.Contains(t, output, "客厅电视 已打开")

	// 执行"睡觉时间"场景
	output = captureOutput(func() {
		err := remote.OnButtonPressed(3)
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "客厅灯 已关闭")
	assert.Contains(t, output, "厨房灯 已关闭")
	assert.Contains(t, output, "客厅电视 已关闭")
	assert.Contains(t, output, "卧室灯 已打开")
	assert.Contains(t, output, "卧室灯 亮度设置为 30%")
}
