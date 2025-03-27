package decorator

import "fmt"

// Component 是组件的基础接口，定义了可以被装饰的对象的行为
type Component interface {
	Show() string // 返回字符串而非直接打印，提高灵活性和可测试性
}

// ConcreteComponent 是具体的组件实现
type ConcreteComponent struct {
	name string
}

// NewConcreteComponent 创建一个具体组件实例
func NewConcreteComponent(name string) *ConcreteComponent {
	return &ConcreteComponent{name: name}
}

// Show 实现 Component 接口
func (c *ConcreteComponent) Show() string {
	return c.name
}

// BaseDecorator 是所有装饰器的基类
type BaseDecorator struct {
	component Component
}

// NewBaseDecorator 创建一个基础装饰器
func NewBaseDecorator(component Component) *BaseDecorator {
	return &BaseDecorator{component: component}
}

// Show 实现 Component 接口，委托给被装饰的组件
func (d *BaseDecorator) Show() string {
	return d.component.Show()
}

// MakeupDecorator 是化妆装饰器的抽象基类
type MakeupDecorator struct {
	BaseDecorator
	decorationType string
}

// NewMakeupDecorator 创建一个化妆装饰器
func NewMakeupDecorator(component Component, decorationType string) *MakeupDecorator {
	return &MakeupDecorator{
		BaseDecorator:  BaseDecorator{component: component},
		decorationType: decorationType,
	}
}

// Show 使用模板方法模式实现装饰逻辑
func (m *MakeupDecorator) Show() string {
	return fmt.Sprintf("%s【%s】", m.decorationType, m.BaseDecorator.Show())
}

// FoundationDecorator 是粉底装饰器
type FoundationDecorator struct {
	*MakeupDecorator
}

// NewFoundationDecorator 创建一个粉底装饰器
func NewFoundationDecorator(component Component) *FoundationDecorator {
	return &FoundationDecorator{
		MakeupDecorator: NewMakeupDecorator(component, "打粉底"),
	}
}

// LipstickDecorator 是口红装饰器
type LipstickDecorator struct {
	*MakeupDecorator
}

// NewLipstickDecorator 创建一个口红装饰器
func NewLipstickDecorator(component Component) *LipstickDecorator {
	return &LipstickDecorator{
		MakeupDecorator: NewMakeupDecorator(component, "涂口红"),
	}
}

// EyeshadowDecorator 是眼影装饰器 (新增)
type EyeshadowDecorator struct {
	*MakeupDecorator
}

// NewEyeshadowDecorator 创建一个眼影装饰器
func NewEyeshadowDecorator(component Component) *EyeshadowDecorator {
	return &EyeshadowDecorator{
		MakeupDecorator: NewMakeupDecorator(component, "画眼影"),
	}
}

// AccessoryDecorator 是配饰装饰器的抽象基类 (新增，展示不同类型的装饰器)
type AccessoryDecorator struct {
	BaseDecorator
	accessoryType string
}

// NewAccessoryDecorator 创建一个配饰装饰器
func NewAccessoryDecorator(component Component, accessoryType string) *AccessoryDecorator {
	return &AccessoryDecorator{
		BaseDecorator: BaseDecorator{component: component},
		accessoryType: accessoryType,
	}
}

// Show 使用不同的装饰风格
func (a *AccessoryDecorator) Show() string {
	return fmt.Sprintf("%s + %s", a.BaseDecorator.Show(), a.accessoryType)
}

// NecklaceDecorator 是项链装饰器 (新增)
type NecklaceDecorator struct {
	*AccessoryDecorator
}

// NewNecklaceDecorator 创建一个项链装饰器
func NewNecklaceDecorator(component Component) *NecklaceDecorator {
	return &NecklaceDecorator{
		AccessoryDecorator: NewAccessoryDecorator(component, "项链"),
	}
}

// EarringsDecorator 是耳环装饰器 (新增)
type EarringsDecorator struct {
	*AccessoryDecorator
}

// NewEarringsDecorator 创建一个耳环装饰器
func NewEarringsDecorator(component Component) *EarringsDecorator {
	return &EarringsDecorator{
		AccessoryDecorator: NewAccessoryDecorator(component, "耳环"),
	}
}

// DisplayComponent 用于展示组件和装饰后的结果
func DisplayComponent(component Component) {
	fmt.Println(component.Show())
}
