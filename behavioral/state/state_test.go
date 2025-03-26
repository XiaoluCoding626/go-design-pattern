package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试初始化自动售货机
func TestNewVendingMachine(t *testing.T) {
	// 创建带有库存的售货机
	inventory := map[string]int{
		"A1": 5,
		"B2": 3,
		"C3": 0,
	}
	vm := NewVendingMachine(inventory)

	// 验证初始状态
	assert.Equal(t, "等待投币", vm.GetCurrentState())
	assert.Equal(t, 5, vm.GetProductCount("A1"))
	assert.Equal(t, 3, vm.GetProductCount("B2"))
	assert.Equal(t, 0, vm.GetProductCount("C3"))

	// 创建无库存的售货机
	emptyVM := NewVendingMachine(map[string]int{})
	assert.Equal(t, "商品售罄", emptyVM.GetCurrentState())
}

// 测试完整购买流程
func TestPurchaseFlow(t *testing.T) {
	inventory := map[string]int{"A1": 1}
	vm := NewVendingMachine(inventory)

	// 尝试在投币前选择商品
	err := vm.SelectProduct("A1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "请先投币")

	// 投币
	err = vm.InsertCoin()
	assert.NoError(t, err)
	assert.Equal(t, "已投币，等待选择", vm.GetCurrentState())

	// 重复投币
	err = vm.InsertCoin()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已投币")

	// 选择商品
	err = vm.SelectProduct("A1")
	assert.NoError(t, err)
	assert.Equal(t, "已选择商品，等待确认", vm.GetCurrentState())

	// 出货
	err = vm.Dispense()
	assert.NoError(t, err)
	assert.Equal(t, 0, vm.GetProductCount("A1"))
	assert.Equal(t, "商品售罄", vm.GetCurrentState()) // 商品售罄

	// 尝试再次购买
	err = vm.InsertCoin()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "所有商品已售罄")
}

// 测试退款流程
func TestRefundFlow(t *testing.T) {
	inventory := map[string]int{"A1": 1}
	vm := NewVendingMachine(inventory)

	// 未投币状态下尝试退款
	err := vm.Refund()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "您未投币")

	// 投币
	err = vm.InsertCoin()
	assert.NoError(t, err)

	// 投币后退款
	err = vm.Refund()
	assert.NoError(t, err)
	assert.Equal(t, "等待投币", vm.GetCurrentState())

	// 再次投币并选择商品
	err = vm.InsertCoin()
	assert.NoError(t, err)
	err = vm.SelectProduct("A1")
	assert.NoError(t, err)

	// 选择商品后退款
	err = vm.Refund()
	assert.NoError(t, err)
	assert.Equal(t, "等待投币", vm.GetCurrentState())
}

// 测试售罄商品
func TestSoldOutProduct(t *testing.T) {
	inventory := map[string]int{"A1": 0, "B2": 1}
	vm := NewVendingMachine(inventory)

	// 投币
	err := vm.InsertCoin()
	assert.NoError(t, err)

	// 尝试选择售罄商品
	err = vm.SelectProduct("A1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已售罄")
	assert.Equal(t, "已投币，等待选择", vm.GetCurrentState())

	// 选择有库存的商品
	err = vm.SelectProduct("B2")
	assert.NoError(t, err)
}

// 测试不存在的商品
func TestNonExistentProduct(t *testing.T) {
	inventory := map[string]int{"A1": 1}
	vm := NewVendingMachine(inventory)

	// 投币
	vm.InsertCoin()

	// 尝试选择不存在的商品
	err := vm.SelectProduct("X9")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已售罄") // 不存在的商品被当作已售罄处理
}

// 示例：展示状态模式的使用
func ExampleVendingMachine() {
	// 创建自动售货机并设置库存
	inventory := map[string]int{
		"A1": 3, // 可乐
		"B2": 2, // 薯片
	}
	vm := NewVendingMachine(inventory)

	// 购买流程展示
	vm.InsertCoin()
	vm.SelectProduct("A1")
	vm.Dispense()

	// 再次购买
	vm.InsertCoin()
	vm.SelectProduct("B2")
	vm.Refund() // 决定退款不购买

	// Output:
	// 投币成功！请选择商品
	// 已选择商品 A1
	// 商品 A1 已出货，谢谢购买！
	// 投币成功！请选择商品
	// 已选择商品 B2
	// 退币成功
}
