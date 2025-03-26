package template_method

import "fmt"

// AbstractSoyaMilk 抽象基类，定义了制作豆浆的算法骨架
type AbstractSoyaMilk struct {
	// 嵌入接口，强制子类实现某些方法
	soyaMilkBehavior
}

// soyaMilkBehavior 子类需要实现的接口，包含所有可能被覆盖的方法
type soyaMilkBehavior interface {
	AddCondiment()                 // 添加配料，由子类实现
	CustomerWantsCondiments() bool // 是否需要调料，由子类决定
	Hook()                         // 额外的钩子方法，由子类按需实现
}

// Make 模板方法，定义了算法骨架
func (s *AbstractSoyaMilk) Make() {
	fmt.Println("=== 开始制作豆浆 ===")
	s.SelectBean()

	// 使用钩子方法判断是否需要添加配料
	// 注意：这里调用的是接口方法，会动态分派到子类实现
	if s.soyaMilkBehavior.CustomerWantsCondiments() {
		s.soyaMilkBehavior.AddCondiment()
	}

	s.Soak()
	s.Beat()
	s.soyaMilkBehavior.Hook() // 提供额外的钩子方法供子类扩展
	fmt.Println("=== 豆浆制作完成 ===")
}

// SelectBean 选择原料，具体方法
func (s *AbstractSoyaMilk) SelectBean() {
	fmt.Println("第 1 步：选择新鲜的豆子")
}

// Soak 浸泡，具体方法
func (s *AbstractSoyaMilk) Soak() {
	fmt.Println("第 3 步：豆子和配料开始浸泡 3 小时")
}

// Beat 榨汁，具体方法
func (s *AbstractSoyaMilk) Beat() {
	fmt.Println("第 4 步：豆子和配料放入豆浆机榨汁")
}

// 默认实现，供子类继承
func (s *AbstractSoyaMilk) CustomerWantsCondiments() bool {
	return true // 默认需要添加配料
}

// 默认空实现，供子类继承
func (s *AbstractSoyaMilk) Hook() {
	// 默认为空实现，子类可以按需覆盖
}

// RedBeanSoyaMilk 红豆豆浆，具体子类
type RedBeanSoyaMilk struct {
	AbstractSoyaMilk
}

// NewRedBeanSoyaMilk 创建红豆豆浆实例
func NewRedBeanSoyaMilk() *RedBeanSoyaMilk {
	milk := &RedBeanSoyaMilk{}
	// 设置AbstractSoyaMilk中的接口为当前实例，确保动态分派
	milk.soyaMilkBehavior = milk
	return milk
}

// AddCondiment 实现具体的配料添加
func (r *RedBeanSoyaMilk) AddCondiment() {
	fmt.Println("第 2 步：加入上好的红豆")
}

// PeanutSoyaMilk 花生豆浆，具体子类
type PeanutSoyaMilk struct {
	AbstractSoyaMilk
}

// NewPeanutSoyaMilk 创建花生豆浆实例
func NewPeanutSoyaMilk() *PeanutSoyaMilk {
	milk := &PeanutSoyaMilk{}
	milk.soyaMilkBehavior = milk
	return milk
}

// AddCondiment 实现具体的配料添加
func (p *PeanutSoyaMilk) AddCondiment() {
	fmt.Println("第 2 步：加入上好的花生")
}

// Hook 覆盖钩子方法，添加一些额外步骤
func (p *PeanutSoyaMilk) Hook() {
	fmt.Println("第 5 步：花生豆浆完成后撒一些花生碎")
}

// PureSoyaMilk 纯豆浆，不需要配料的子类
type PureSoyaMilk struct {
	AbstractSoyaMilk
}

// NewPureSoyaMilk 创建纯豆浆实例
func NewPureSoyaMilk() *PureSoyaMilk {
	milk := &PureSoyaMilk{}
	milk.soyaMilkBehavior = milk
	return milk
}

// AddCondiment 纯豆浆不添加配料 (即使不会被调用也需要实现接口)
func (p *PureSoyaMilk) AddCondiment() {
	// 空实现，因为纯豆浆不添加配料
}

// CustomerWantsCondiments 覆盖钩子方法，明确表示不添加配料
func (p *PureSoyaMilk) CustomerWantsCondiments() bool {
	return false // 纯豆浆不需要配料
}
