package facade

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// 捕获标准输出的辅助函数
func captureOutput(f func()) string {
	// 保存原始的标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行需要捕获输出的函数
	f()

	// 恢复原始的标准输出
	w.Close()
	os.Stdout = oldStdout

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// 测试创建外观实例
func TestNewRestaurantFacade(t *testing.T) {
	facade := NewRestaurantFacade()

	if facade == nil {
		t.Error("NewRestaurantFacade() 返回了 nil")
	}

	// 验证所有子系统组件已初始化
	if facade.vegVendor == nil {
		t.Error("vegVendor 未初始化")
	}
	if facade.meatVendor == nil {
		t.Error("meatVendor 未初始化")
	}
	if facade.warehouse == nil {
		t.Error("warehouse 未初始化")
	}
	if facade.chef == nil {
		t.Error("chef 未初始化")
	}
	if facade.cashier == nil {
		t.Error("cashier 未初始化")
	}
	if facade.waiter == nil {
		t.Error("waiter 未初始化")
	}
	if facade.cleaner == nil {
		t.Error("cleaner 未初始化")
	}
}

// 测试点餐流程
func TestOrderFood(t *testing.T) {
	facade := NewRestaurantFacade()
	output := captureOutput(func() {
		facade.OrderFood()
	})

	// 验证输出中包含关键步骤信息
	expectedSteps := []string{
		"开始点餐服务",
		"服务员接收顾客点单",
		"从仓库取出需要的食材",
		"清洗并准备食材",
		"厨师烹饪美食",
		"服务员将食物送到顾客桌前",
		"点餐服务完成",
	}

	for _, step := range expectedSteps {
		if !strings.Contains(output, step) {
			t.Errorf("OrderFood() 输出中缺少预期步骤: %s", step)
		}
	}
}

// 测试准备食材流程
func TestPrepareIngredients(t *testing.T) {
	facade := NewRestaurantFacade()
	output := captureOutput(func() {
		facade.PrepareIngredients()
	})

	expectedSteps := []string{
		"开始准备食材",
		"采购新鲜蔬菜食材",
		"采购优质肉类食材",
		"将蔬菜存入仓库",
		"将肉类存入仓库",
		"食材准备完成",
	}

	for _, step := range expectedSteps {
		if !strings.Contains(output, step) {
			t.Errorf("PrepareIngredients() 输出中缺少预期步骤: %s", step)
		}
	}
}

// 测试结账流程
func TestHandlePayment(t *testing.T) {
	facade := NewRestaurantFacade()
	output := captureOutput(func() {
		facade.HandlePayment()
	})

	expectedSteps := []string{
		"开始结账流程",
		"收取顾客付款",
		"生成消费收据",
		"结账完成",
	}

	for _, step := range expectedSteps {
		if !strings.Contains(output, step) {
			t.Errorf("HandlePayment() 输出中缺少预期步骤: %s", step)
		}
	}
}

// 测试清理餐厅流程
func TestCleanupRestaurant(t *testing.T) {
	facade := NewRestaurantFacade()
	output := captureOutput(func() {
		facade.CleanupRestaurant()
	})

	expectedSteps := []string{
		"开始清理餐厅",
		"清理餐桌和餐具",
		"清扫餐厅地面",
		"餐厅清理完成",
	}

	for _, step := range expectedSteps {
		if !strings.Contains(output, step) {
			t.Errorf("CleanupRestaurant() 输出中缺少预期步骤: %s", step)
		}
	}
}

// 测试完整服务周期
func TestCompleteServiceCycle(t *testing.T) {
	facade := NewRestaurantFacade()
	output := captureOutput(func() {
		facade.CompleteServiceCycle()
	})

	// 只检查周期开始和结束标记，因为其他流程已在各自的测试中验证
	expectedMarkers := []string{
		"开始完整服务周期",
		"开始准备食材",
		"开始点餐服务",
		"开始结账流程",
		"开始清理餐厅",
		"服务周期结束",
	}

	for _, marker := range expectedMarkers {
		if !strings.Contains(output, marker) {
			t.Errorf("CompleteServiceCycle() 输出中缺少预期标记: %s", marker)
		}
	}
}

// 示例：使用外观模式简化点餐流程
func ExampleRestaurantFacade_OrderFood() {
	facade := NewRestaurantFacade()
	facade.OrderFood()
	// Output:
	//
	// ===== 开始点餐服务 =====
	// 服务员接收顾客点单
	// 从仓库取出需要的食材
	// 清洗并准备食材
	// 厨师烹饪美食
	// 服务员将食物送到顾客桌前
	// ===== 点餐服务完成 =====
}

// 示例：使用外观模式处理支付流程
func ExampleRestaurantFacade_HandlePayment() {
	facade := NewRestaurantFacade()
	facade.HandlePayment()
	// Output:
	//
	// ===== 开始结账流程 =====
	// 收取顾客付款
	// 生成消费收据
	// ===== 结账完成 =====
}

// 示例：外观模式在食材准备中的应用
func ExampleRestaurantFacade_PrepareIngredients() {
	facade := NewRestaurantFacade()
	facade.PrepareIngredients()
	// Output:
	//
	// ===== 开始准备食材 =====
	// 采购新鲜蔬菜食材
	// 采购优质肉类食材
	// 将蔬菜存入仓库
	// 将肉类存入仓库
	// ===== 食材准备完成 =====
}

// 示例：外观模式在清理流程中的应用
func ExampleRestaurantFacade_CleanupRestaurant() {
	facade := NewRestaurantFacade()
	facade.CleanupRestaurant()
	// Output:
	//
	// ===== 开始清理餐厅 =====
	// 清理餐桌和餐具
	// 清扫餐厅地面
	// ===== 餐厅清理完成 =====
}

// 表格驱动测试 - 测试所有外观方法是否能正常执行
func TestAllFacadeMethods(t *testing.T) {
	facade := NewRestaurantFacade()

	tests := []struct {
		name   string
		method func()
	}{
		{"OrderFood", facade.OrderFood},
		{"PrepareIngredients", facade.PrepareIngredients},
		{"HandlePayment", facade.HandlePayment},
		{"CleanupRestaurant", facade.CleanupRestaurant},
		{"CompleteServiceCycle", facade.CompleteServiceCycle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 只检查方法是否能执行而不会崩溃
			output := captureOutput(tt.method)
			if len(output) == 0 {
				t.Errorf("%s() 没有产生任何输出", tt.name)
			}
		})
	}
}
