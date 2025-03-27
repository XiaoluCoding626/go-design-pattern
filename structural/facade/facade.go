package facade

import "fmt"

// 子系统组件 - 这些组件构成了复杂的子系统

// VegVendor 蔬菜供应商
type VegVendor struct{}

// Purchase 采购蔬菜食材
func (v *VegVendor) Purchase() {
	fmt.Println("采购新鲜蔬菜食材")
}

// MeatVendor 肉类供应商
type MeatVendor struct{}

// Purchase 采购肉类食材
func (m *MeatVendor) Purchase() {
	fmt.Println("采购优质肉类食材")
}

// Warehouse 食材仓库
type Warehouse struct{}

// Store 存储食材
func (w *Warehouse) Store(item string) {
	fmt.Printf("将%s存入仓库\n", item)
}

// Retrieve 取出食材
func (w *Warehouse) Retrieve(item string) {
	fmt.Printf("从仓库取出%s\n", item)
}

// Chef 厨师
type Chef struct{}

// PrepareIngredients 准备食材
func (c *Chef) PrepareIngredients() {
	fmt.Println("清洗并准备食材")
}

// Cook 烹饪食物
func (c *Chef) Cook() {
	fmt.Println("厨师烹饪美食")
}

// Cashier 收银员
type Cashier struct{}

// CollectPayment 收取付款
func (c *Cashier) CollectPayment() {
	fmt.Println("收取顾客付款")
}

// GenerateReceipt 生成收据
func (c *Cashier) GenerateReceipt() {
	fmt.Println("生成消费收据")
}

// Waiter 服务员
type Waiter struct{}

// TakeOrder 接收订单
func (w *Waiter) TakeOrder() {
	fmt.Println("服务员接收顾客点单")
}

// ServeFood 上菜
func (w *Waiter) ServeFood() {
	fmt.Println("服务员将食物送到顾客桌前")
}

// Cleaner 清洁工
type Cleaner struct{}

// CleanTable 清理餐桌
func (c *Cleaner) CleanTable() {
	fmt.Println("清理餐桌和餐具")
}

// CleanFloor 清扫地面
func (c *Cleaner) CleanFloor() {
	fmt.Println("清扫餐厅地面")
}

// RestaurantFacade 是餐厅外观类，为客户端提供简单统一的接口
type RestaurantFacade struct {
	vegVendor  *VegVendor
	meatVendor *MeatVendor
	warehouse  *Warehouse
	chef       *Chef
	cashier    *Cashier
	waiter     *Waiter
	cleaner    *Cleaner
}

// NewRestaurantFacade 创建一个餐厅外观实例
func NewRestaurantFacade() *RestaurantFacade {
	return &RestaurantFacade{
		vegVendor:  &VegVendor{},
		meatVendor: &MeatVendor{},
		warehouse:  &Warehouse{},
		chef:       &Chef{},
		cashier:    &Cashier{},
		waiter:     &Waiter{},
		cleaner:    &Cleaner{},
	}
}

// OrderFood 提供给客户端的点餐服务流程
func (f *RestaurantFacade) OrderFood() {
	// 隐藏内部系统的复杂性，提供简单的接口
	fmt.Println("\n===== 开始点餐服务 =====")
	f.waiter.TakeOrder()
	f.warehouse.Retrieve("需要的食材")
	f.chef.PrepareIngredients()
	f.chef.Cook()
	f.waiter.ServeFood()
	fmt.Println("===== 点餐服务完成 =====")
}

// PrepareIngredients 提供给管理者的食材准备流程
func (f *RestaurantFacade) PrepareIngredients() {
	// 隐藏采购和存储的复杂流程
	fmt.Println("\n===== 开始准备食材 =====")
	f.vegVendor.Purchase()
	f.meatVendor.Purchase()
	f.warehouse.Store("蔬菜")
	f.warehouse.Store("肉类")
	fmt.Println("===== 食材准备完成 =====")
}

// HandlePayment 提供给客户端的结账流程
func (f *RestaurantFacade) HandlePayment() {
	// 简化结账流程
	fmt.Println("\n===== 开始结账流程 =====")
	f.cashier.CollectPayment()
	f.cashier.GenerateReceipt()
	fmt.Println("===== 结账完成 =====")
}

// CleanupRestaurant 提供给管理者的清理餐厅流程
func (f *RestaurantFacade) CleanupRestaurant() {
	// 整合清理流程
	fmt.Println("\n===== 开始清理餐厅 =====")
	f.cleaner.CleanTable()
	f.cleaner.CleanFloor()
	fmt.Println("===== 餐厅清理完成 =====")
}

// CompleteServiceCycle 提供完整的服务周期，整合多个流程
func (f *RestaurantFacade) CompleteServiceCycle() {
	fmt.Println("\n======= 开始完整服务周期 =======")
	f.PrepareIngredients()
	f.OrderFood()
	f.HandlePayment()
	f.CleanupRestaurant()
	fmt.Println("======= 服务周期结束 =======")
}
