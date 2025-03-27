package bridge

import "fmt"

// Device 表示设备的接口，这是"实现部分"的接口
type Device interface {
	TurnOn()         // 开启设备
	TurnOff()        // 关闭设备
	SetVolume(int)   // 设置音量
	GetName() string // 获取设备名称
}

// RemoteControl 表示遥控器的抽象，这是"抽象部分"的基础
type RemoteControl interface {
	PowerOn()    // 开启电源
	PowerOff()   // 关闭电源
	VolumeUp()   // 提高音量
	VolumeDown() // 降低音量
}

// TV 电视机实现了Device接口
type TV struct {
	name   string
	isOn   bool
	volume int
}

// NewTV 创建一个新的电视机
func NewTV(name string) *TV {
	return &TV{
		name:   name,
		isOn:   false,
		volume: 10,
	}
}

// TurnOn 开启电视机
func (t *TV) TurnOn() {
	t.isOn = true
	fmt.Printf("%s 电视机打开了，当前音量：%d\n", t.name, t.volume)
}

// TurnOff 关闭电视机
func (t *TV) TurnOff() {
	t.isOn = false
	fmt.Printf("%s 电视机关闭了\n", t.name)
}

// SetVolume 设置电视机音量
func (t *TV) SetVolume(volume int) {
	if volume < 0 {
		volume = 0
	} else if volume > 100 {
		volume = 100
	}
	t.volume = volume
	fmt.Printf("%s 电视机音量设置为：%d\n", t.name, t.volume)
}

// GetName 获取电视机名称
func (t *TV) GetName() string {
	return t.name
}

// Radio 收音机实现了Device接口
type Radio struct {
	name   string
	isOn   bool
	volume int
}

// NewRadio 创建一个新的收音机
func NewRadio(name string) *Radio {
	return &Radio{
		name:   name,
		isOn:   false,
		volume: 5,
	}
}

// TurnOn 开启收音机
func (r *Radio) TurnOn() {
	r.isOn = true
	fmt.Printf("%s 收音机打开了，当前音量：%d\n", r.name, r.volume)
}

// TurnOff 关闭收音机
func (r *Radio) TurnOff() {
	r.isOn = false
	fmt.Printf("%s 收音机关闭了\n", r.name)
}

// SetVolume 设置收音机音量
func (r *Radio) SetVolume(volume int) {
	if volume < 0 {
		volume = 0
	} else if volume > 100 {
		volume = 100
	}
	r.volume = volume
	fmt.Printf("%s 收音机音量设置为：%d\n", r.name, r.volume)
}

// GetName 获取收音机名称
func (r *Radio) GetName() string {
	return r.name
}

// BaseRemoteControl 是所有遥控器的基础实现
type BaseRemoteControl struct {
	device Device // 持有对Device的引用——这是桥接模式的核心
	volume int    // 当前音量
}

// NewBaseRemoteControl 创建一个新的基础遥控器
func NewBaseRemoteControl(device Device) *BaseRemoteControl {
	return &BaseRemoteControl{
		device: device,
		volume: 10,
	}
}

// PowerOn 开启设备
func (r *BaseRemoteControl) PowerOn() {
	r.device.TurnOn()
}

// PowerOff 关闭设备
func (r *BaseRemoteControl) PowerOff() {
	r.device.TurnOff()
}

// VolumeUp 提高音量
func (r *BaseRemoteControl) VolumeUp() {
	r.volume += 10
	r.device.SetVolume(r.volume)
}

// VolumeDown 降低音量
func (r *BaseRemoteControl) VolumeDown() {
	r.volume -= 10
	if r.volume < 0 {
		r.volume = 0
	}
	r.device.SetVolume(r.volume)
}

// StandardRemoteControl 标准遥控器扩展了基础遥控器
type StandardRemoteControl struct {
	*BaseRemoteControl
}

// NewStandardRemoteControl 创建一个新的标准遥控器
func NewStandardRemoteControl(device Device) *StandardRemoteControl {
	return &StandardRemoteControl{
		BaseRemoteControl: NewBaseRemoteControl(device),
	}
}

// AdvancedRemoteControl 高级遥控器扩展了基础遥控器，添加了额外功能
type AdvancedRemoteControl struct {
	*BaseRemoteControl
}

// NewAdvancedRemoteControl 创建一个新的高级遥控器
func NewAdvancedRemoteControl(device Device) *AdvancedRemoteControl {
	return &AdvancedRemoteControl{
		BaseRemoteControl: NewBaseRemoteControl(device),
	}
}

// Mute 静音功能（高级遥控器特有）
func (a *AdvancedRemoteControl) Mute() {
	a.device.SetVolume(0)
	fmt.Printf("静音 %s\n", a.device.GetName())
}

// MaxVolume 最大音量功能（高级遥控器特有）
func (a *AdvancedRemoteControl) MaxVolume() {
	a.device.SetVolume(100)
	fmt.Printf("将 %s 音量调到最大\n", a.device.GetName())
}
