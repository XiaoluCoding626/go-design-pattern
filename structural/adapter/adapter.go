package adapter

import "fmt"

// IPlug 定义所有插头的通用接口
type IPlug interface {
	GetPin() int
}

// TwoPinPlug 两针插头的具体实现
type TwoPinPlug struct{}

// GetPin 返回两针插头的针数
func (t *TwoPinPlug) GetPin() int {
	return 2
}

// ThreePinPlug 三针插头的具体实现
type ThreePinPlug struct{}

// GetPin 返回三针插头的针数
func (t *ThreePinPlug) GetPin() int {
	return 3
}

// IPowerSocket 定义所有电源插座的通用接口
type IPowerSocket interface {
	Charge(p IPlug) // 统一的充电方法
}

// ThreePinSocket 三孔插座 - 被适配的目标类
type ThreePinSocket struct{}

// Charge 三孔插座只能为三针插头充电
func (s *ThreePinSocket) Charge(p IPlug) {
	if p.GetPin() != 3 {
		fmt.Println("三孔插座无法为非三针插头充电")
		return
	}
	fmt.Println("三孔插座正在为三针插头充电")
}

// TwoPinSocket 两孔插座
type TwoPinSocket struct{}

// Charge 两孔插座只能为两针插头充电
func (s *TwoPinSocket) Charge(p IPlug) {
	if p.GetPin() != 2 {
		fmt.Println("两孔插座无法为非两针插头充电")
		return
	}
	fmt.Println("两孔插座正在为两针插头充电")
}

// PowerAdapter 电源适配器 - 将两针插头适配到三孔插座
type PowerAdapter struct {
	socket *ThreePinSocket // 持有一个三孔插座的引用
}

// NewPowerAdapter 创建一个新的电源适配器
func NewPowerAdapter() *PowerAdapter {
	return &PowerAdapter{
		socket: &ThreePinSocket{},
	}
}

// Charge 适配器的充电方法 - 实现IPowerSocket接口
// 当接收到两针插头时，进行适配转换后使用三孔插座充电
func (a *PowerAdapter) Charge(p IPlug) {
	if p.GetPin() != 2 {
		fmt.Println("适配器只能适配两针插头")
		return
	}

	// 创建一个虚拟的三针插头，这是适配的核心逻辑
	fmt.Println("适配器正在将两针插头转换为三针插头")
	virtualThreePinPlug := &ThreePinPlug{}

	// 使用三孔插座为虚拟的三针插头充电
	a.socket.Charge(virtualThreePinPlug)
	fmt.Println("适配完成，两针插头已成功充电")
}
