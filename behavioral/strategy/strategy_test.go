package strategy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试加法策略
func TestAddStrategy(t *testing.T) {
	add := &Add{}
	result, err := add.Do(5, 3)

	assert.NoError(t, err, "加法策略不应返回错误")
	assert.Equal(t, 8, result, "加法策略应返回正确结果")
}

// 测试减法策略
func TestSubtractStrategy(t *testing.T) {
	subtract := &Subtract{}
	result, err := subtract.Do(5, 3)

	assert.NoError(t, err, "减法策略不应返回错误")
	assert.Equal(t, 2, result, "减法策略应返回正确结果")
}

// 测试乘法策略
func TestMultiplyStrategy(t *testing.T) {
	multiply := &Multiply{}
	result, err := multiply.Do(5, 3)

	assert.NoError(t, err, "乘法策略不应返回错误")
	assert.Equal(t, 15, result, "乘法策略应返回正确结果")
}

// 测试除法策略
func TestDivideStrategy(t *testing.T) {
	divide := &Divide{}

	// 正常情况
	result, err := divide.Do(10, 2)
	assert.NoError(t, err, "正常除法不应返回错误")
	assert.Equal(t, 5, result, "除法策略应返回正确结果")

	// 除以零的情况
	_, err = divide.Do(10, 0)
	assert.Equal(t, ErrDivideByZero, err, "除以零应返回特定错误")
}

// 测试操作者
func TestOperator(t *testing.T) {
	// 使用加法策略
	operator := NewOperator(&Add{})
	result, err := operator.Calculate(5, 3)
	assert.NoError(t, err, "使用加法策略的操作者不应返回错误")
	assert.Equal(t, 8, result, "使用加法策略的操作者应返回正确结果")

	// 切换到减法策略
	operator.SetStrategy(&Subtract{})
	result, err = operator.Calculate(5, 3)
	assert.NoError(t, err, "使用减法策略的操作者不应返回错误")
	assert.Equal(t, 2, result, "使用减法策略的操作者应返回正确结果")

	// 切换到乘法策略
	operator.SetStrategy(&Multiply{})
	result, err = operator.Calculate(5, 3)
	assert.NoError(t, err, "使用乘法策略的操作者不应返回错误")
	assert.Equal(t, 15, result, "使用乘法策略的操作者应返回正确结果")

	// 切换到除法策略
	operator.SetStrategy(&Divide{})
	result, err = operator.Calculate(10, 2)
	assert.NoError(t, err, "使用除法策略的操作者不应返回错误")
	assert.Equal(t, 5, result, "使用除法策略的操作者应返回正确结果")

	// 测试除以零的情况
	_, err = operator.Calculate(10, 0)
	assert.Equal(t, ErrDivideByZero, err, "使用除法策略的操作者在除以零时应返回特定错误")
}

// 测试未设置策略的情况
func TestOperatorWithNoStrategy(t *testing.T) {
	operator := &Operator{} // 不设置策略
	_, err := operator.Calculate(5, 3)
	assert.EqualError(t, err, "未设置策略", "未设置策略的操作者应返回错误信息")
}

// 测试示例
func ExampleOperator() {
	// 创建操作者并设置初始策略为加法
	operator := NewOperator(&Add{})

	// 使用加法策略
	result, _ := operator.Calculate(5, 3)
	fmt.Println("5 + 3 =", result)

	// 切换到减法策略
	operator.SetStrategy(&Subtract{})
	result, _ = operator.Calculate(5, 3)
	fmt.Println("5 - 3 =", result)

	// 切换到乘法策略
	operator.SetStrategy(&Multiply{})
	result, _ = operator.Calculate(5, 3)
	fmt.Println("5 * 3 =", result)

	// 切换到除法策略
	operator.SetStrategy(&Divide{})
	result, _ = operator.Calculate(6, 3)
	fmt.Println("6 / 3 =", result)

	// Output:
	// 5 + 3 = 8
	// 5 - 3 = 2
	// 5 * 3 = 15
	// 6 / 3 = 2
}
