package strategy

import (
	"fmt"
)

// 策略模式：定义一系列算法，将每个算法封装起来，并使它们可以互换。
// 策略模式使算法可以独立于使用它的客户端而变化。

// IStrategy 定义所有支持算法的接口
type IStrategy interface {
	// Do 执行具体的算法，并返回结果和可能的错误
	Do(a, b int) (int, error)
}

// Add 实现加法策略
type Add struct{}

// Do 执行加法操作
func (*Add) Do(a, b int) (int, error) {
	return a + b, nil
}

// Subtract 实现减法策略
type Subtract struct{}

// Do 执行减法操作
func (*Subtract) Do(a, b int) (int, error) {
	return a - b, nil
}

// Multiply 实现乘法策略
type Multiply struct{}

// Do 执行乘法操作
func (*Multiply) Do(a, b int) (int, error) {
	return a * b, nil
}

// Divide 实现除法策略
type Divide struct{}

// Do 执行除法操作，并处理除零错误
func (*Divide) Do(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

// ErrDivideByZero 当尝试除以零时返回此错误
var ErrDivideByZero = fmt.Errorf("不能除以零")

// Operator 是使用策略的上下文
type Operator struct {
	strategy IStrategy
}

// SetStrategy 更改操作者使用的策略
func (o *Operator) SetStrategy(strategy IStrategy) {
	o.strategy = strategy
}

// Calculate 将当前策略应用于给定的操作数
func (o *Operator) Calculate(a, b int) (int, error) {
	if o.strategy == nil {
		return 0, fmt.Errorf("未设置策略")
	}
	return o.strategy.Do(a, b)
}

// NewOperator 创建一个带有指定策略的新操作者
func NewOperator(strategy IStrategy) *Operator {
	return &Operator{strategy: strategy}
}
