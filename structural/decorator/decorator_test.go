package decorator

import (
	"fmt"
	"testing"
)

// TestConcreteComponent 测试基础组件功能
func TestConcreteComponent(t *testing.T) {
	component := NewConcreteComponent("素颜")
	if result := component.Show(); result != "素颜" {
		t.Errorf("ConcreteComponent.Show() = %v, 期望 %v", result, "素颜")
	}
}

// TestSingleDecorator 测试单个装饰器
func TestSingleDecorator(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(Component) Component
		expected string
	}{
		{
			name: "粉底装饰器",
			setup: func(c Component) Component {
				return NewFoundationDecorator(c)
			},
			expected: "打粉底【素颜】",
		},
		{
			name: "口红装饰器",
			setup: func(c Component) Component {
				return NewLipstickDecorator(c)
			},
			expected: "涂口红【素颜】",
		},
		{
			name: "眼影装饰器",
			setup: func(c Component) Component {
				return NewEyeshadowDecorator(c)
			},
			expected: "画眼影【素颜】",
		},
		{
			name: "项链装饰器",
			setup: func(c Component) Component {
				return NewNecklaceDecorator(c)
			},
			expected: "素颜 + 项链",
		},
		{
			name: "耳环装饰器",
			setup: func(c Component) Component {
				return NewEarringsDecorator(c)
			},
			expected: "素颜 + 耳环",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := NewConcreteComponent("素颜")
			decorated := tt.setup(base)
			if result := decorated.Show(); result != tt.expected {
				t.Errorf("%s.Show() = %v, 期望 %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestMultipleDecorators 测试多层嵌套装饰器
func TestMultipleDecorators(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(Component) Component
		expected string
	}{
		{
			name: "粉底+口红",
			setup: func(c Component) Component {
				return NewLipstickDecorator(NewFoundationDecorator(c))
			},
			expected: "涂口红【打粉底【素颜】】",
		},
		{
			name: "粉底+眼影+口红",
			setup: func(c Component) Component {
				return NewLipstickDecorator(NewEyeshadowDecorator(NewFoundationDecorator(c)))
			},
			expected: "涂口红【画眼影【打粉底【素颜】】】",
		},
		{
			name: "完整妆容+配饰",
			setup: func(c Component) Component {
				withMakeup := NewLipstickDecorator(NewEyeshadowDecorator(NewFoundationDecorator(c)))
				return NewEarringsDecorator(NewNecklaceDecorator(withMakeup))
			},
			expected: "涂口红【画眼影【打粉底【素颜】】】 + 项链 + 耳环",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := NewConcreteComponent("素颜")
			decorated := tt.setup(base)
			if result := decorated.Show(); result != tt.expected {
				t.Errorf("%s.Show() = %v, 期望 %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestMixedDecorators 测试混合不同类型的装饰器
func TestMixedDecorators(t *testing.T) {
	// 先配饰后化妆
	t.Run("先配饰后化妆", func(t *testing.T) {
		base := NewConcreteComponent("素颜")
		withAccessories := NewNecklaceDecorator(base)
		withMakeup := NewFoundationDecorator(withAccessories)

		expected := "打粉底【素颜 + 项链】"
		if result := withMakeup.Show(); result != expected {
			t.Errorf("得到: %v, 期望: %v", result, expected)
		}
	})

	// 交替使用不同装饰器
	t.Run("交替使用装饰器", func(t *testing.T) {
		base := NewConcreteComponent("素颜")
		step1 := NewFoundationDecorator(base)
		step2 := NewNecklaceDecorator(step1)
		step3 := NewLipstickDecorator(step2)
		step4 := NewEarringsDecorator(step3)

		expected := "涂口红【打粉底【素颜】 + 项链】 + 耳环"
		if result := step4.Show(); result != expected {
			t.Errorf("得到: %v, 期望: %v", result, expected)
		}
	})
}

// ExampleFoundationDecorator 展示基本的装饰器用法
func ExampleFoundationDecorator() {
	girl := NewConcreteComponent("素颜")
	withFoundation := NewFoundationDecorator(girl)
	DisplayComponent(withFoundation)
	// Output: 打粉底【素颜】
}

// DemoMultipleDecorators 展示多个装饰器组合使用
func DemoMultipleDecorators() {
	girl := NewConcreteComponent("素颜")

	// 添加多层妆容
	withMakeup := NewLipstickDecorator(
		NewEyeshadowDecorator(
			NewFoundationDecorator(girl),
		),
	)

	DisplayComponent(withMakeup)
	// Output: 涂口红【画眼影【打粉底【素颜】】】
}

// DemoMixedDecorators 展示混合使用不同类型的装饰器
func DemoMixedDecorators() {
	person := NewConcreteComponent("基本形象")

	// 添加妆容和配饰
	enhanced := NewNecklaceDecorator(
		NewLipstickDecorator(
			NewFoundationDecorator(person),
		),
	)

	fmt.Println(enhanced.Show())
	// Output: 涂口红【打粉底【基本形象】】 + 项链
}

// BenchmarkDecoratorChain 测试装饰器链的性能
func BenchmarkDecoratorChain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		base := NewConcreteComponent("素颜")
		decorated := NewEarringsDecorator(
			NewNecklaceDecorator(
				NewLipstickDecorator(
					NewEyeshadowDecorator(
						NewFoundationDecorator(base),
					),
				),
			),
		)
		_ = decorated.Show()
	}
}

// TestDisplayComponent 测试展示功能
func TestDisplayComponent(t *testing.T) {
	// 这个测试主要是为了覆盖DisplayComponent函数
	// 因为它直接打印到标准输出，我们只是确保它不会崩溃
	component := NewConcreteComponent("测试组件")
	DisplayComponent(component) // 应该打印 "测试组件"
}

// TestBaseDecorator 测试基础装饰器
func TestBaseDecorator(t *testing.T) {
	base := NewConcreteComponent("基础组件")
	decorator := NewBaseDecorator(base)

	if result := decorator.Show(); result != "基础组件" {
		t.Errorf("BaseDecorator.Show() = %v, 期望 %v", result, "基础组件")
	}
}
