package command

import (
	"fmt"
	"strings"
)

// Command 接口定义了命令的执行和撤销方法
type Command interface {
	Execute() error
	Undo() error
	Name() string
}

// Device 表示可以接收命令的设备接口
type Device interface {
	On() error
	Off() error
	GetName() string
}

// Light 表示灯的接收者
type Light struct {
	name  string
	isOn  bool
	level int // 亮度级别
}

// NewLight 创建一个新的灯
func NewLight(name string) *Light {
	return &Light{
		name:  name,
		isOn:  false,
		level: 0,
	}
}

// On 打开灯的操作
func (l *Light) On() error {
	if l.isOn {
		return fmt.Errorf("%s 已经是开启状态", l.name)
	}
	l.isOn = true
	l.level = 100
	fmt.Printf("%s 已打开\n", l.name)
	return nil
}

// Off 关闭灯的操作
func (l *Light) Off() error {
	if !l.isOn {
		return fmt.Errorf("%s 已经是关闭状态", l.name)
	}
	l.isOn = false
	l.level = 0
	fmt.Printf("%s 已关闭\n", l.name)
	return nil
}

// GetName 返回设备名称
func (l *Light) GetName() string {
	return l.name
}

// SetLevel 设置灯的亮度
func (l *Light) SetLevel(level int) error {
	if level < 0 || level > 100 {
		return fmt.Errorf("亮度必须在0-100之间")
	}
	if level > 0 && !l.isOn {
		l.isOn = true
	} else if level == 0 && l.isOn {
		l.isOn = false
	}
	l.level = level
	fmt.Printf("%s 亮度设置为 %d%%\n", l.name, level)
	return nil
}

// TV 表示电视的接收者
type TV struct {
	name    string
	isOn    bool
	volume  int
	channel int
}

// NewTV 创建一个新的电视
func NewTV(name string) *TV {
	return &TV{
		name:    name,
		isOn:    false,
		volume:  50,
		channel: 1,
	}
}

// On 打开电视的操作
func (t *TV) On() error {
	if t.isOn {
		return fmt.Errorf("%s 已经是开启状态", t.name)
	}
	t.isOn = true
	fmt.Printf("%s 已打开, 音量: %d, 频道: %d\n", t.name, t.volume, t.channel)
	return nil
}

// Off 关闭电视的操作
func (t *TV) Off() error {
	if !t.isOn {
		return fmt.Errorf("%s 已经是关闭状态", t.name)
	}
	t.isOn = false
	fmt.Printf("%s 已关闭\n", t.name)
	return nil
}

// GetName 返回设备名称
func (t *TV) GetName() string {
	return t.name
}

// SetVolume 设置电视音量
func (t *TV) SetVolume(volume int) error {
	if !t.isOn {
		return fmt.Errorf("%s 处于关闭状态，无法调整音量", t.name)
	}
	if volume < 0 || volume > 100 {
		return fmt.Errorf("音量必须在0-100之间")
	}
	t.volume = volume
	fmt.Printf("%s 音量设置为 %d\n", t.name, volume)
	return nil
}

// SetChannel 设置电视频道
func (t *TV) SetChannel(channel int) error {
	if !t.isOn {
		return fmt.Errorf("%s 处于关闭状态，无法切换频道", t.name)
	}
	if channel < 1 {
		return fmt.Errorf("频道必须大于0")
	}
	t.channel = channel
	fmt.Printf("%s 切换到频道 %d\n", t.name, channel)
	return nil
}

// TurnOnCommand 表示开启设备命令
type TurnOnCommand struct {
	device Device
}

// NewTurnOnCommand 创建一个新的开启命令
func NewTurnOnCommand(device Device) *TurnOnCommand {
	return &TurnOnCommand{
		device: device,
	}
}

// Execute 执行开启命令
func (c *TurnOnCommand) Execute() error {
	return c.device.On()
}

// Undo 撤销开启命令
func (c *TurnOnCommand) Undo() error {
	return c.device.Off()
}

// Name 返回命令名称
func (c *TurnOnCommand) Name() string {
	return fmt.Sprintf("开启 %s", c.device.GetName())
}

// TurnOffCommand 表示关闭设备命令
type TurnOffCommand struct {
	device Device
}

// NewTurnOffCommand 创建一个新的关闭命令
func NewTurnOffCommand(device Device) *TurnOffCommand {
	return &TurnOffCommand{
		device: device,
	}
}

// Execute 执行关闭命令
func (c *TurnOffCommand) Execute() error {
	return c.device.Off()
}

// Undo 撤销关闭命令
func (c *TurnOffCommand) Undo() error {
	return c.device.On()
}

// Name 返回命令名称
func (c *TurnOffCommand) Name() string {
	return fmt.Sprintf("关闭 %s", c.device.GetName())
}

// SetLevelCommand 表示设置灯亮度的命令
type SetLevelCommand struct {
	light     *Light
	level     int
	prevLevel int
}

// NewSetLevelCommand 创建一个新的设置亮度命令
func NewSetLevelCommand(light *Light, level int) *SetLevelCommand {
	return &SetLevelCommand{
		light:     light,
		level:     level,
		prevLevel: light.level,
	}
}

// Execute 执行设置亮度命令
func (c *SetLevelCommand) Execute() error {
	c.prevLevel = c.light.level
	return c.light.SetLevel(c.level)
}

// Undo 撤销设置亮度命令
func (c *SetLevelCommand) Undo() error {
	return c.light.SetLevel(c.prevLevel)
}

// Name 返回命令名称
func (c *SetLevelCommand) Name() string {
	return fmt.Sprintf("设置 %s 亮度为 %d%%", c.light.name, c.level)
}

// MacroCommand 表示宏命令，可以执行多个命令
type MacroCommand struct {
	name     string
	commands []Command
}

// NewMacroCommand 创建一个新的宏命令
func NewMacroCommand(name string, commands []Command) *MacroCommand {
	return &MacroCommand{
		name:     name,
		commands: commands,
	}
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

// Undo 按相反顺序撤销所有命令
func (m *MacroCommand) Undo() error {
	for i := len(m.commands) - 1; i >= 0; i-- {
		if err := m.commands[i].Undo(); err != nil {
			return fmt.Errorf("撤销宏命令 %s 时出错: %s 失败: %w", m.name, m.commands[i].Name(), err)
		}
	}
	return nil
}

// Name 返回宏命令名称
func (m *MacroCommand) Name() string {
	return m.name
}

// RemoteControl 表示命令调用者（遥控器）
type RemoteControl struct {
	onCommands    []Command
	offCommands   []Command
	history       []Command
	maxHistoryLen int
}

// NewRemoteControl 创建一个新的遥控器
func NewRemoteControl(slots int) *RemoteControl {
	onCommands := make([]Command, slots)
	offCommands := make([]Command, slots)

	// 初始化为无操作命令
	for i := 0; i < slots; i++ {
		onCommands[i] = &NoOpCommand{}
		offCommands[i] = &NoOpCommand{}
	}

	return &RemoteControl{
		onCommands:    onCommands,
		offCommands:   offCommands,
		history:       make([]Command, 0),
		maxHistoryLen: 10,
	}
}

// SetCommand 设置遥控器按钮对应的命令
func (r *RemoteControl) SetCommand(slot int, onCommand Command, offCommand Command) error {
	if slot < 0 || slot >= len(r.onCommands) {
		return fmt.Errorf("无效的插槽编号: %d", slot)
	}

	r.onCommands[slot] = onCommand
	r.offCommands[slot] = offCommand
	return nil
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

// OffButtonPressed 按下关闭按钮
func (r *RemoteControl) OffButtonPressed(slot int) error {
	if slot < 0 || slot >= len(r.offCommands) {
		return fmt.Errorf("无效的插槽编号: %d", slot)
	}

	cmd := r.offCommands[slot]
	err := cmd.Execute()
	if err == nil {
		r.addToHistory(cmd)
	}
	return err
}

// addToHistory 添加命令到历史记录
func (r *RemoteControl) addToHistory(cmd Command) {
	r.history = append(r.history, cmd)
	if len(r.history) > r.maxHistoryLen {
		// 移除最旧的命令
		r.history = r.history[1:]
	}
}

// UndoLastCommand 撤销最后执行的命令
func (r *RemoteControl) UndoLastCommand() error {
	if len(r.history) == 0 {
		return fmt.Errorf("没有可撤销的命令")
	}

	lastIndex := len(r.history) - 1
	lastCmd := r.history[lastIndex]
	r.history = r.history[:lastIndex]

	return lastCmd.Undo()
}

// ShowHistory 展示命令历史记录
func (r *RemoteControl) ShowHistory() {
	if len(r.history) == 0 {
		fmt.Println("命令历史记录为空")
		return
	}

	fmt.Println("命令历史记录:")
	for i, cmd := range r.history {
		fmt.Printf("%d: %s\n", i+1, cmd.Name())
	}
}

// NoOpCommand 表示无操作命令
type NoOpCommand struct{}

func (c *NoOpCommand) Execute() error { return nil }
func (c *NoOpCommand) Undo() error    { return nil }
func (c *NoOpCommand) Name() string   { return "无操作" }

// String 返回遥控器描述
func (r *RemoteControl) String() string {
	var sb strings.Builder

	sb.WriteString("\n------ 遥控器 ------\n")
	for i := 0; i < len(r.onCommands); i++ {
		onName := r.onCommands[i].Name()
		offName := r.offCommands[i].Name()
		sb.WriteString(fmt.Sprintf("[%d] %-20s %-20s\n", i, onName, offName))
	}
	sb.WriteString("------------------\n")

	return sb.String()
}
