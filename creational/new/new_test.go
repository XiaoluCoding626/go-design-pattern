package new

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

// 浮点数比较函数，允许一定的误差范围
func floatEqual(a, b float64) bool {
	epsilon := 0.00001 // 允许的误差
	return math.Abs(a-b) < epsilon
}

// 测试基本构造函数NewProduct
func TestNewProduct(t *testing.T) {
	// 正常情况
	p, err := NewProduct("手机", 1999.99)
	if err != nil {
		t.Fatalf("创建有效商品时出错: %v", err)
	}
	if p.GetName() != "手机" {
		t.Errorf("商品名称应为 '手机', 实际为: %s", p.GetName())
	}
	if !floatEqual(p.GetPrice(), 1999.99) {
		t.Errorf("商品价格应为 1999.99, 实际为: %.2f", p.GetPrice())
	}
	if p.GetCategory() != "未分类" {
		t.Errorf("默认商品类别应为 '未分类', 实际为: %s", p.GetCategory())
	}
	if p.GetStock() != 0 {
		t.Errorf("默认商品库存应为 0, 实际为: %d", p.GetStock())
	}
	if !floatEqual(p.GetDiscount(), 0) {
		t.Errorf("默认商品折扣百分比应为 0, 实际为: %.1f", p.GetDiscount())
	}
	if p.ID == "" {
		t.Error("商品ID不应为空")
	}
	if p.CreatedAt.After(time.Now()) || p.CreatedAt.Before(time.Now().Add(-time.Minute)) {
		t.Error("商品创建时间不在合理范围内")
	}

	// 测试参数验证
	testCases := []struct {
		name     string
		prodName string
		price    float64
		wantErr  string
	}{
		{"空名称", "", 100, "商品名称不能为空"},
		{"零价格", "电视", 0, "商品价格必须大于零"},
		{"负价格", "电视", -10, "商品价格必须大于零"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewProduct(tc.prodName, tc.price)
			if err == nil {
				t.Error("应该返回错误，但没有")
			} else if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("错误信息应包含 '%s', 实际为: %s", tc.wantErr, err)
			}
		})
	}
}

// 测试带折扣的构造函数
func TestNewDiscountedProduct(t *testing.T) {
	// 正常情况
	p, err := NewDiscountedProduct("笔记本", 6999.99, 15)
	if err != nil {
		t.Fatalf("创建折扣商品时出错: %v", err)
	}

	expectedPrice := 6999.99 * 0.85 // 15% 折扣
	if !floatEqual(p.GetPrice(), expectedPrice) {
		t.Errorf("折扣后价格应为 %.2f, 实际为: %.2f", expectedPrice, p.GetPrice())
	}
	if !floatEqual(p.GetDiscount(), 15) {
		t.Errorf("折扣百分比应为 15, 实际为: %.1f", p.GetDiscount())
	}

	// 测试参数验证
	testCases := []struct {
		name     string
		prodName string
		price    float64
		discount float64
		wantErr  string
	}{
		{"负折扣", "电视", 100, -10, "折扣百分比必须在0到100之间"},
		{"过大折扣", "电视", 100, 120, "折扣百分比必须在0到100之间"},
		{"无效价格", "电视", -10, 20, "商品价格必须大于零"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewDiscountedProduct(tc.prodName, tc.price, tc.discount)
			if err == nil {
				t.Error("应该返回错误，但没有")
			} else if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("错误信息应包含 '%s', 实际为: %s", tc.wantErr, err)
			}
		})
	}
}

// 测试带库存的构造函数
func TestNewProductInStock(t *testing.T) {
	// 正常情况
	p, err := NewProductInStock("耳机", 299.99, 50)
	if err != nil {
		t.Fatalf("创建带库存商品时出错: %v", err)
	}

	if p.GetStock() != 50 {
		t.Errorf("库存应为 50, 实际为: %d", p.GetStock())
	}

	// 测试参数验证
	_, err = NewProductInStock("键盘", 399.99, -10)
	if err == nil {
		t.Error("负库存应该返回错误，但没有")
	} else if !strings.Contains(err.Error(), "初始库存不能为负数") {
		t.Errorf("错误信息应包含 '初始库存不能为负数', 实际为: %s", err)
	}
}

// 测试完整构造函数
func TestNewProductComplete(t *testing.T) {
	// 正常情况
	p, err := NewProductComplete("平板电脑", 3499.99, "电子产品", 25, 10)
	if err != nil {
		t.Fatalf("创建完整商品时出错: %v", err)
	}

	expectedPrice := 3499.99 * 0.9 // 10% 折扣
	if p.GetName() != "平板电脑" {
		t.Errorf("商品名称应为 '平板电脑', 实际为: %s", p.GetName())
	}
	if !floatEqual(p.GetPrice(), expectedPrice) {
		t.Errorf("折扣后价格应为 %.2f, 实际为: %.2f", expectedPrice, p.GetPrice())
	}
	if p.GetCategory() != "电子产品" {
		t.Errorf("商品类别应为 '电子产品', 实际为: %s", p.GetCategory())
	}
	if p.GetStock() != 25 {
		t.Errorf("库存应为 25, 实际为: %d", p.GetStock())
	}

	// 测试参数验证
	testCases := []struct {
		name     string
		prodName string
		price    float64
		category string
		stock    int
		discount float64
		wantErr  string
	}{
		{"空类别", "电视", 100, "", 10, 5, "商品类别不能为空"},
		{"负库存", "电视", 100, "电器", -5, 10, "初始库存不能为负数"},
		{"无效折扣", "电视", 100, "电器", 10, 110, "折扣百分比必须在0到100之间"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewProductComplete(tc.prodName, tc.price, tc.category, tc.stock, tc.discount)
			if err == nil {
				t.Error("应该返回错误，但没有")
			} else if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("错误信息应包含 '%s', 实际为: %s", tc.wantErr, err)
			}
		})
	}
}

// 测试链式方法
func TestChainMethods(t *testing.T) {
	// 创建基本商品并使用链式方法设置属性
	p, err := NewProduct("鼠标", 99.99)
	if err != nil {
		t.Fatalf("创建基本商品时出错: %v", err)
	}

	p.WithCategory("电脑配件").WithStock(100).WithDiscount(20)

	if p.GetCategory() != "电脑配件" {
		t.Errorf("类别应为 '电脑配件', 实际为: %s", p.GetCategory())
	}
	if p.GetStock() != 100 {
		t.Errorf("库存应为 100, 实际为: %d", p.GetStock())
	}
	if !floatEqual(p.GetDiscount(), 20) {
		t.Errorf("折扣应为 20, 实际为: %.1f", p.GetDiscount())
	}

	expectedPrice := 99.99 * 0.8 // 20% 折扣
	if !floatEqual(p.GetPrice(), expectedPrice) {
		t.Errorf("折扣后价格应为 %.2f, 实际为: %.2f", expectedPrice, p.GetPrice())
	}

	// 测试无效参数处理
	originalCategory := p.GetCategory()
	p.WithCategory("")
	if p.GetCategory() != originalCategory {
		t.Errorf("空类别不应更改原值，应为 '%s', 实际为: %s", originalCategory, p.GetCategory())
	}

	originalStock := p.GetStock()
	p.WithStock(-10)
	if p.GetStock() != originalStock {
		t.Errorf("负库存不应更改原值，应为 %d, 实际为: %d", originalStock, p.GetStock())
	}

	originalDiscount := p.GetDiscount()
	p.WithDiscount(110)
	if !floatEqual(p.GetDiscount(), originalDiscount) {
		t.Errorf("无效折扣不应更改原值，应为 %.1f, 实际为: %.1f", originalDiscount, p.GetDiscount())
	}
}

// 测试获取属性方法
func TestGetters(t *testing.T) {
	p, err := NewProductComplete("显示器", 1299.99, "电脑外设", 30, 15)
	if err != nil {
		t.Fatalf("创建测试商品时出错: %v", err)
	}

	expectedTests := []struct {
		name     string
		got      interface{}
		expected interface{}
		isFloat  bool
	}{
		{"GetName", p.GetName(), "显示器", false},
		{"GetOriginalPrice", p.GetOriginalPrice(), 1299.99, true},
		{"GetPrice", p.GetPrice(), 1299.99 * 0.85, true}, // 15% 折扣
		{"GetCategory", p.GetCategory(), "电脑外设", false},
		{"GetStock", p.GetStock(), 30, false},
		{"GetDiscount", p.GetDiscount(), 15.0, true},
	}

	for _, test := range expectedTests {
		t.Run(test.name, func(t *testing.T) {
			if test.isFloat {
				// 浮点数比较
				gotFloat, ok := test.got.(float64)
				if !ok {
					t.Fatalf("%s: 无法将结果转换为float64", test.name)
				}
				expectedFloat, ok := test.expected.(float64)
				if !ok {
					t.Fatalf("%s: 无法将期望值转换为float64", test.name)
				}
				if !floatEqual(gotFloat, expectedFloat) {
					t.Errorf("%s: 期望 %v, 实际 %v", test.name, expectedFloat, gotFloat)
				}
			} else {
				// 直接比较
				if test.got != test.expected {
					t.Errorf("%s: 期望 %v, 实际 %v", test.name, test.expected, test.got)
				}
			}
		})
	}
}

// 测试库存修改方法
func TestStockMethods(t *testing.T) {
	p, _ := NewProductInStock("键盘", 199.99, 50)

	// 测试增加库存
	err := p.AddStock(20)
	if err != nil {
		t.Errorf("增加库存出错: %v", err)
	}
	if p.GetStock() != 70 {
		t.Errorf("增加库存后，应为 70, 实际为: %d", p.GetStock())
	}

	// 测试减少库存
	err = p.ReduceStock(30)
	if err != nil {
		t.Errorf("减少库存出错: %v", err)
	}
	if p.GetStock() != 40 {
		t.Errorf("减少库存后，应为 40, 实际为: %d", p.GetStock())
	}

	// 测试库存不足
	err = p.ReduceStock(50)
	if err == nil {
		t.Error("库存不足应返回错误，但没有")
	} else if !strings.Contains(err.Error(), "库存不足") {
		t.Errorf("错误信息应包含 '库存不足', 实际为: %s", err)
	}

	// 测试参数验证
	err = p.AddStock(-10)
	if err == nil {
		t.Error("负增量应返回错误，但没有")
	}

	err = p.ReduceStock(-5)
	if err == nil {
		t.Error("负减量应返回错误，但没有")
	}
}

// 测试折扣应用方法
func TestApplyDiscount(t *testing.T) {
	p, _ := NewProduct("扫地机器人", 2499.99)

	// 测试应用折扣
	err := p.ApplyDiscount(30)
	if err != nil {
		t.Errorf("应用折扣出错: %v", err)
	}

	expectedPrice := 2499.99 * 0.7 // 30% 折扣
	if !floatEqual(p.GetPrice(), expectedPrice) {
		t.Errorf("折扣后价格应为 %.2f, 实际为: %.2f", expectedPrice, p.GetPrice())
	}

	// 测试参数验证
	err = p.ApplyDiscount(-10)
	if err == nil {
		t.Error("负折扣百分比应返回错误，但没有")
	}

	err = p.ApplyDiscount(120)
	if err == nil {
		t.Error("过大折扣百分比应返回错误，但没有")
	}
}

// 测试String方法
func TestString(t *testing.T) {
	// 测试无折扣商品
	p1, _ := NewProduct("咖啡机", 899.99)
	str1 := p1.String()

	requiredParts := []string{
		"商品: 咖啡机",
		"价格: ¥899.99",
		"库存: 0",
		"类别: 未分类",
	}

	for _, part := range requiredParts {
		if !strings.Contains(str1, part) {
			t.Errorf("String()输出应包含 '%s', 实际输出: %s", part, str1)
		}
	}

	// 测试带折扣商品
	p2, _ := NewDiscountedProduct("微波炉", 699.99, 25)
	str2 := p2.String()

	if !strings.Contains(str2, "折扣: 25.0%") {
		t.Errorf("折扣商品的String()输出应包含折扣信息, 实际输出: %s", str2)
	}

	expectedDiscountedPrice := 699.99 * 0.75 // 25% 折扣
	discountedPriceStr := fmt.Sprintf("折后价: ¥%.2f", expectedDiscountedPrice)
	if !strings.Contains(str2, discountedPriceStr) {
		t.Errorf("折扣商品的String()输出应包含折后价格 '%s', 实际输出: %s",
			discountedPriceStr, str2)
	}
}

// 测试克隆方法
func TestClone(t *testing.T) {
	// 创建一个完整的商品
	original, _ := NewProductComplete("相机", 4999.99, "数码产品", 15, 10)

	// 克隆商品
	clone := original.Clone()

	// 验证基本属性复制正确
	if clone.GetName() != original.GetName() {
		t.Errorf("克隆名称应为 '%s', 实际为: %s", original.GetName(), clone.GetName())
	}
	if !floatEqual(clone.GetPrice(), original.GetPrice()) {
		t.Errorf("克隆价格应为 %.2f, 实际为: %.2f", original.GetPrice(), clone.GetPrice())
	}
	if clone.GetCategory() != original.GetCategory() {
		t.Errorf("克隆类别应为 '%s', 实际为: %s", original.GetCategory(), clone.GetCategory())
	}
	if clone.GetStock() != original.GetStock() {
		t.Errorf("克隆库存应为 %d, 实际为: %d", original.GetStock(), clone.GetStock())
	}
	if !floatEqual(clone.GetDiscount(), original.GetDiscount()) {
		t.Errorf("克隆折扣应为 %.1f, 实际为: %.1f", original.GetDiscount(), clone.GetDiscount())
	}

	// 验证ID和CreatedAt是不同的
	if clone.ID == original.ID {
		t.Error("克隆和原始商品的ID应该不同")
	}
	if clone.CreatedAt.Equal(original.CreatedAt) {
		t.Error("克隆和原始商品的创建时间应该不同")
	}

	// 修改原始商品属性，确认克隆不受影响
	original.AddStock(10)
	original.ApplyDiscount(20)

	if clone.GetStock() == original.GetStock() {
		t.Error("修改原始商品库存后，克隆库存不应跟着变化")
	}
	if floatEqual(clone.GetDiscount(), original.GetDiscount()) {
		t.Error("修改原始商品折扣后，克隆折扣不应跟着变化")
	}
}

// 集成测试
func TestIntegration(t *testing.T) {
	// 模拟商品上架、销售和促销场景

	// 1. 创建商品并入库
	phone, err := NewProductComplete("智能手机", 5999.99, "电子产品", 100, 0)
	if err != nil {
		t.Fatalf("创建商品失败: %v", err)
	}

	// 2. 销售10个商品
	err = phone.ReduceStock(10)
	if err != nil {
		t.Fatalf("减少库存失败: %v", err)
	}

	if phone.GetStock() != 90 {
		t.Errorf("销售10个后，库存应为90, 实际为: %d", phone.GetStock())
	}

	// 3. 促销活动，应用15%折扣
	err = phone.ApplyDiscount(15)
	if err != nil {
		t.Fatalf("应用折扣失败: %v", err)
	}

	expectedPrice := 5999.99 * 0.85
	if !floatEqual(phone.GetPrice(), expectedPrice) {
		t.Errorf("促销后价格应为 %.2f, 实际为: %.2f", expectedPrice, phone.GetPrice())
	}

	// 4. 获取原价（用于财务记录）
	originalPrice := phone.GetOriginalPrice()
	if !floatEqual(originalPrice, 5999.99) {
		t.Errorf("原价应为 5999.99, 实际为: %.2f", originalPrice)
	}

	// 5. 补充库存
	err = phone.AddStock(50)
	if err != nil {
		t.Fatalf("增加库存失败: %v", err)
	}

	if phone.GetStock() != 140 {
		t.Errorf("补充库存后，总库存应为140, 实际为: %d", phone.GetStock())
	}

	// 6. 克隆商品创建新款
	newPhone := phone.Clone()
	newPhone.WithCategory("新品电子").WithDiscount(5) // 5%折扣

	if newPhone.GetCategory() != "新品电子" {
		t.Errorf("新商品类别应为 '新品电子', 实际为: %s", newPhone.GetCategory())
	}

	newExpectedPrice := 5999.99 * 0.95 // 5% 折扣
	if !floatEqual(newPhone.GetPrice(), newExpectedPrice) {
		t.Errorf("新商品价格应为 %.2f, 实际为: %.2f", newExpectedPrice, newPhone.GetPrice())
	}

	// 确认原商品折扣未变
	if !floatEqual(phone.GetDiscount(), 15) {
		t.Errorf("原商品折扣应保持 15%%, 实际为: %.1f%%", phone.GetDiscount())
	}
}

// 示例测试，展示New模式的常见用法
func ExampleProduct() {
	// 创建基本商品
	laptop, _ := NewProduct("笔记本电脑", 6999.99)

	// 使用链式方法设置额外属性
	laptop.WithCategory("电脑").WithStock(10).WithDiscount(10)

	// 打印商品信息
	fmt.Println(laptop)

	// 创建折扣商品
	phone, _ := NewDiscountedProduct("智能手机", 4999.99, 20)
	fmt.Printf("手机原价: ¥%.2f, 折后价: ¥%.2f\n",
		phone.GetOriginalPrice(), phone.GetPrice())

	// 模拟销售
	phone.ReduceStock(1)
	fmt.Printf("销售一部手机后，库存: %d\n", phone.GetStock())
}
