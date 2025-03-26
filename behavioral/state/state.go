package state

import "fmt"

// VendingMachineState 定义了自动售货机的状态接口
type VendingMachineState interface {
	InsertCoin() error               // 投币
	SelectProduct(code string) error // 选择商品
	Dispense() error                 // 出货
	Refund() error                   // 退款
	GetStateName() string            // 获取状态名称
}

// VendingMachine 自动售货机类，维护当前状态
type VendingMachine struct {
	// 各种状态
	noMoneyState         VendingMachineState
	hasCoinState         VendingMachineState
	productSelectedState VendingMachineState
	soldOutState         VendingMachineState

	// 当前状态
	currentState VendingMachineState

	// 商品库存
	inventory map[string]int
	// 当前选中的商品
	selectedProduct string
}

// NewVendingMachine 创建并初始化售货机
func NewVendingMachine(inventory map[string]int) *VendingMachine {
	machine := &VendingMachine{
		inventory: inventory,
	}

	// 初始化各种状态
	noMoneyState := &NoMoneyState{machine: machine}
	hasCoinState := &HasCoinState{machine: machine}
	productSelectedState := &ProductSelectedState{machine: machine}
	soldOutState := &SoldOutState{machine: machine}

	// 设置各个状态
	machine.noMoneyState = noMoneyState
	machine.hasCoinState = hasCoinState
	machine.productSelectedState = productSelectedState
	machine.soldOutState = soldOutState

	// 设置初始状态
	if machine.hasAnyProducts() {
		machine.currentState = noMoneyState
	} else {
		machine.currentState = soldOutState
	}

	return machine
}

// 检查是否有任何商品库存
func (vm *VendingMachine) hasAnyProducts() bool {
	for _, count := range vm.inventory {
		if count > 0 {
			return true
		}
	}
	return false
}

// SetState 设置售货机的当前状态
func (vm *VendingMachine) SetState(state VendingMachineState) {
	vm.currentState = state
}

// GetCurrentState 获取当前状态
func (vm *VendingMachine) GetCurrentState() string {
	return vm.currentState.GetStateName()
}

// InsertCoin 用户投币操作
func (vm *VendingMachine) InsertCoin() error {
	return vm.currentState.InsertCoin()
}

// SelectProduct 用户选择商品
func (vm *VendingMachine) SelectProduct(code string) error {
	return vm.currentState.SelectProduct(code)
}

// Dispense 售货机出货
func (vm *VendingMachine) Dispense() error {
	return vm.currentState.Dispense()
}

// Refund 退款操作
func (vm *VendingMachine) Refund() error {
	return vm.currentState.Refund()
}

// GetInventory 获取商品库存
func (vm *VendingMachine) GetInventory() map[string]int {
	return vm.inventory
}

// ReduceInventory 减少商品库存
func (vm *VendingMachine) ReduceInventory(code string) {
	if count, exists := vm.inventory[code]; exists && count > 0 {
		vm.inventory[code] = count - 1
	}
}

// HasInventory 检查特定商品是否有库存
func (vm *VendingMachine) HasInventory(code string) bool {
	if count, exists := vm.inventory[code]; exists && count > 0 {
		return true
	}
	return false
}

// GetProductCount 获取特定商品的库存数量
func (vm *VendingMachine) GetProductCount(code string) int {
	if count, exists := vm.inventory[code]; exists {
		return count
	}
	return 0
}

// SetSelectedProduct 设置当前选中的商品
func (vm *VendingMachine) SetSelectedProduct(code string) {
	vm.selectedProduct = code
}

// GetSelectedProduct 获取当前选中的商品
func (vm *VendingMachine) GetSelectedProduct() string {
	return vm.selectedProduct
}

// NoMoneyState 未投币状态
type NoMoneyState struct {
	machine *VendingMachine
}

func (s *NoMoneyState) InsertCoin() error {
	fmt.Println("投币成功！请选择商品")
	s.machine.SetState(s.machine.hasCoinState)
	return nil
}

func (s *NoMoneyState) SelectProduct(code string) error {
	return fmt.Errorf("请先投币")
}

func (s *NoMoneyState) Dispense() error {
	return fmt.Errorf("请先投币并选择商品")
}

func (s *NoMoneyState) Refund() error {
	return fmt.Errorf("您未投币，无法退款")
}

func (s *NoMoneyState) GetStateName() string {
	return "等待投币"
}

// HasCoinState 已投币状态
type HasCoinState struct {
	machine *VendingMachine
}

func (s *HasCoinState) InsertCoin() error {
	return fmt.Errorf("已投币，请选择商品或退币")
}

func (s *HasCoinState) SelectProduct(code string) error {
	if !s.machine.HasInventory(code) {
		return fmt.Errorf("商品 %s 已售罄", code)
	}

	fmt.Printf("已选择商品 %s\n", code)
	s.machine.SetSelectedProduct(code)
	s.machine.SetState(s.machine.productSelectedState)
	return nil
}

func (s *HasCoinState) Dispense() error {
	return fmt.Errorf("请先选择商品")
}

func (s *HasCoinState) Refund() error {
	fmt.Println("退币成功")
	s.machine.SetState(s.machine.noMoneyState)
	return nil
}

func (s *HasCoinState) GetStateName() string {
	return "已投币，等待选择"
}

// ProductSelectedState 已选择商品状态
type ProductSelectedState struct {
	machine *VendingMachine
}

func (s *ProductSelectedState) InsertCoin() error {
	return fmt.Errorf("已投币并选择了商品，请确认购买或退币")
}

func (s *ProductSelectedState) SelectProduct(code string) error {
	return fmt.Errorf("已选择商品，请确认购买或退币")
}

func (s *ProductSelectedState) Dispense() error {
	code := s.machine.GetSelectedProduct()
	if code == "" {
		return fmt.Errorf("未选择商品")
	}

	// 减少库存
	s.machine.ReduceInventory(code)
	fmt.Printf("商品 %s 已出货，谢谢购买！\n", code)

	// 判断是否还有库存，决定下一个状态
	if s.machine.hasAnyProducts() {
		s.machine.SetState(s.machine.noMoneyState)
	} else {
		s.machine.SetState(s.machine.soldOutState)
	}

	// 清除选中的商品
	s.machine.SetSelectedProduct("")
	return nil
}

func (s *ProductSelectedState) Refund() error {
	fmt.Println("退币成功")
	s.machine.SetSelectedProduct("")
	s.machine.SetState(s.machine.noMoneyState)
	return nil
}

func (s *ProductSelectedState) GetStateName() string {
	return "已选择商品，等待确认"
}

// SoldOutState 售罄状态
type SoldOutState struct {
	machine *VendingMachine
}

func (s *SoldOutState) InsertCoin() error {
	return fmt.Errorf("所有商品已售罄，无法投币")
}

func (s *SoldOutState) SelectProduct(code string) error {
	return fmt.Errorf("所有商品已售罄")
}

func (s *SoldOutState) Dispense() error {
	return fmt.Errorf("所有商品已售罄")
}

func (s *SoldOutState) Refund() error {
	return fmt.Errorf("未投币，无法退款")
}

func (s *SoldOutState) GetStateName() string {
	return "商品售罄"
}
