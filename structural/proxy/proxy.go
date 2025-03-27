package proxy

import (
	"fmt"
	"time"
)

// IBuyCar 定义了被代理对象和代理对象共同实现的接口
// 这是代理模式的核心 - 代理和被代理对象实现相同的接口
type IBuyCar interface {
	BuyCar() error
	GetCarInfo() string
}

// RealBuyer 是实际买车的人（被代理对象 - RealSubject）
type RealBuyer struct {
	Name  string
	Money float64
}

// NewRealBuyer 创建实际买车人的实例
func NewRealBuyer(name string, money float64) *RealBuyer {
	return &RealBuyer{
		Name:  name,
		Money: money,
	}
}

// BuyCar 实现了IBuyCar接口的方法
func (r *RealBuyer) BuyCar() error {
	if r.Money < 100000 {
		return fmt.Errorf("余额不足，无法购买汽车")
	}
	fmt.Printf("<%s> 成功购买了一辆汽车，花费了 ¥%.2f\n", r.Name, 100000.0)
	r.Money -= 100000
	return nil
}

// GetCarInfo 获取车辆信息
func (r *RealBuyer) GetCarInfo() string {
	return "标准汽车型号XYZ"
}

// 4S店代理 - 基本代理，提供额外服务
type FourSProxy struct {
	realBuyer IBuyCar
	services  []string
	fee       float64
}

// NewFourSProxy 创建4S店代理实例
func NewFourSProxy(buyer IBuyCar) *FourSProxy {
	return &FourSProxy{
		realBuyer: buyer,
		services:  []string{"上牌服务", "汽车注册", "保险办理"},
		fee:       5000,
	}
}

// BuyCar 代理实现的购车方法，添加了额外的服务
func (f *FourSProxy) BuyCar() error {
	fmt.Println("=== 通过4S店代理购车开始 ===")

	// 代理前的操作
	fmt.Println("1. 从制造商订购汽车到4S店")
	fmt.Println("2. 准备购车文件")

	// 调用实际对象的方法
	if err := f.realBuyer.BuyCar(); err != nil {
		fmt.Printf("购车失败: %s\n", err)
		return err
	}

	// 代理后的增强操作
	fmt.Println("提供额外服务:")
	for i, service := range f.services {
		fmt.Printf("  %d. %s\n", i+1, service)
	}

	fmt.Printf("收取服务费: ¥%.2f\n", f.fee)
	fmt.Println("=== 通过4S店代理购车完成 ===")
	return nil
}

// GetCarInfo 代理获取车辆信息的方法
func (f *FourSProxy) GetCarInfo() string {
	// 可以添加额外的车辆信息或修改返回内容
	baseInfo := f.realBuyer.GetCarInfo()
	return baseInfo + " (通过4S店提供)"
}

// VirtualBuyerProxy 虚拟代理 - 延迟创建被代理对象，节约资源
type VirtualBuyerProxy struct {
	name      string
	money     float64
	realBuyer *RealBuyer
}

// NewVirtualBuyerProxy 创建虚拟代理实例
func NewVirtualBuyerProxy(name string, money float64) *VirtualBuyerProxy {
	return &VirtualBuyerProxy{
		name:  name,
		money: money,
		// realBuyer 初始为nil，等需要时才创建
	}
}

// BuyCar 虚拟代理实现，延迟创建被代理对象
func (v *VirtualBuyerProxy) BuyCar() error {
	fmt.Println("=== 通过虚拟代理购车开始 ===")
	fmt.Println("准备创建实际购买者...")

	// 延迟初始化 - 仅在首次调用时创建实际对象
	if v.realBuyer == nil {
		fmt.Println("首次调用，创建实际购买者")
		v.realBuyer = NewRealBuyer(v.name, v.money)
	} else {
		fmt.Println("复用已有的实际购买者")
	}

	err := v.realBuyer.BuyCar()
	fmt.Println("=== 通过虚拟代理购车结束 ===")
	return err
}

// GetCarInfo 获取车辆信息
func (v *VirtualBuyerProxy) GetCarInfo() string {
	if v.realBuyer == nil {
		v.realBuyer = NewRealBuyer(v.name, v.money)
	}
	return v.realBuyer.GetCarInfo() + " (虚拟代理提供)"
}

// ProtectionProxy 保护代理 - 控制对资源的访问权限
type ProtectionProxy struct {
	realBuyer IBuyCar
	isVIP     bool
}

// NewProtectionProxy 创建保护代理
func NewProtectionProxy(buyer IBuyCar, isVIP bool) *ProtectionProxy {
	return &ProtectionProxy{
		realBuyer: buyer,
		isVIP:     isVIP,
	}
}

// BuyCar 保护代理实现，加入权限控制
func (p *ProtectionProxy) BuyCar() error {
	fmt.Println("=== 通过保护代理购车开始 ===")

	// 权限检查
	if !p.isVIP {
		fmt.Println("权限不足: 仅VIP客户可以通过此渠道购车")
		return fmt.Errorf("权限不足: 需要VIP权限")
	}

	fmt.Println("VIP客户，权限验证通过")
	err := p.realBuyer.BuyCar()

	if err == nil {
		fmt.Println("VIP客户专享折扣已应用")
	}

	fmt.Println("=== 通过保护代理购车结束 ===")
	return err
}

// GetCarInfo 获取车辆信息
func (p *ProtectionProxy) GetCarInfo() string {
	if !p.isVIP {
		return "基础车辆信息 (需VIP权限查看详细配置)"
	}
	return p.realBuyer.GetCarInfo() + " (VIP专享配置)"
}

// LoggingProxy 日志代理 - 记录操作日志
type LoggingProxy struct {
	realBuyer IBuyCar
}

// NewLoggingProxy 创建日志代理
func NewLoggingProxy(buyer IBuyCar) *LoggingProxy {
	return &LoggingProxy{
		realBuyer: buyer,
	}
}

// BuyCar 日志代理实现，添加日志记录
func (l *LoggingProxy) BuyCar() error {
	fmt.Println("=== 日志记录: 购车操作开始 ===")
	startTime := time.Now()

	fmt.Printf("[%s] 购车请求已接收\n", startTime.Format("2006-01-02 15:04:05"))

	err := l.realBuyer.BuyCar()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	if err != nil {
		fmt.Printf("[%s] 购车失败: %s\n", endTime.Format("2006-01-02 15:04:05"), err)
	} else {
		fmt.Printf("[%s] 购车成功\n", endTime.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("操作耗时: %v\n", duration)
	fmt.Println("=== 日志记录: 购车操作结束 ===")
	return err
}

// GetCarInfo 获取车辆信息并记录日志
func (l *LoggingProxy) GetCarInfo() string {
	fmt.Printf("[%s] 获取车辆信息\n", time.Now().Format("2006-01-02 15:04:05"))
	return l.realBuyer.GetCarInfo()
}

// CachedBuyerProxy 缓存代理 - 缓存重复请求的结果
type CachedBuyerProxy struct {
	realBuyer IBuyCar
	carInfo   string
	cached    bool
}

// NewCachedBuyerProxy 创建缓存代理
func NewCachedBuyerProxy(buyer IBuyCar) *CachedBuyerProxy {
	return &CachedBuyerProxy{
		realBuyer: buyer,
		cached:    false,
	}
}

// BuyCar 实现购车方法，不支持缓存
func (c *CachedBuyerProxy) BuyCar() error {
	fmt.Println("=== 通过缓存代理购车开始 ===")
	fmt.Println("购车操作无法缓存，正在执行实际购车...")
	err := c.realBuyer.BuyCar()
	fmt.Println("=== 通过缓存代理购车结束 ===")
	return err
}

// GetCarInfo 获取车辆信息，支持缓存
func (c *CachedBuyerProxy) GetCarInfo() string {
	if c.cached {
		fmt.Println("从缓存获取车辆信息")
		return c.carInfo + " (缓存)"
	}

	fmt.Println("首次获取车辆信息，将结果缓存")
	c.carInfo = c.realBuyer.GetCarInfo()
	c.cached = true
	return c.carInfo
}
