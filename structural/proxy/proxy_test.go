package proxy

import (
	"bytes"
	"fmt"
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

// 测试RealBuyer
func TestRealBuyer(t *testing.T) {
	t.Run("足够的钱能购买汽车", func(t *testing.T) {
		buyer := NewRealBuyer("张三", 150000)

		output := captureOutput(func() {
			err := buyer.BuyCar()
			if err != nil {
				t.Errorf("购车应该成功，但出现错误: %v", err)
			}
		})

		expected := "成功购买了一辆汽车"
		if !strings.Contains(output, expected) {
			t.Errorf("期望输出包含 '%s'，但得到: %s", expected, output)
		}

		if buyer.Money != 50000 {
			t.Errorf("购车后余额应为50000，但得到: %.2f", buyer.Money)
		}
	})

	t.Run("钱不够无法购买汽车", func(t *testing.T) {
		buyer := NewRealBuyer("李四", 80000)
		err := buyer.BuyCar()

		if err == nil {
			t.Error("钱不够时应返回错误，但没有")
		}

		if !strings.Contains(err.Error(), "余额不足") {
			t.Errorf("错误信息应包含'余额不足'，但得到: %v", err)
		}

		if buyer.Money != 80000 {
			t.Errorf("购车失败后余额应保持不变，但变为: %.2f", buyer.Money)
		}
	})

	t.Run("获取车辆信息", func(t *testing.T) {
		buyer := NewRealBuyer("王五", 100000)
		info := buyer.GetCarInfo()

		expected := "标准汽车型号XYZ"
		if info != expected {
			t.Errorf("期望车辆信息为 '%s'，但得到: %s", expected, info)
		}
	})
}

// 测试4S店代理
func TestFourSProxy(t *testing.T) {
	t.Run("通过4S店购车成功", func(t *testing.T) {
		buyer := NewRealBuyer("赵六", 200000)
		proxy := NewFourSProxy(buyer)

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err != nil {
				t.Errorf("购车应该成功，但出现错误: %v", err)
			}
		})

		expectedPhrases := []string{
			"通过4S店代理购车开始",
			"从制造商订购汽车到4S店",
			"准备购车文件",
			"成功购买了一辆汽车",
			"提供额外服务",
			"上牌服务",
			"汽车注册",
			"保险办理",
			"收取服务费: ¥5000.00",
			"通过4S店代理购车完成",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("输出应包含 '%s'，但未找到", phrase)
			}
		}
	})

	t.Run("通过4S店购车失败", func(t *testing.T) {
		buyer := NewRealBuyer("孙七", 50000)
		proxy := NewFourSProxy(buyer)

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err == nil {
				t.Error("余额不足时应返回错误，但没有")
			}
		})

		if !strings.Contains(output, "购车失败") {
			t.Errorf("输出应包含购车失败信息，但未找到")
		}
	})

	t.Run("获取车辆信息", func(t *testing.T) {
		buyer := NewRealBuyer("周八", 100000)
		proxy := NewFourSProxy(buyer)
		info := proxy.GetCarInfo()

		expected := "标准汽车型号XYZ (通过4S店提供)"
		if info != expected {
			t.Errorf("期望车辆信息为 '%s'，但得到: %s", expected, info)
		}
	})
}

// 测试虚拟代理
func TestVirtualBuyerProxy(t *testing.T) {
	t.Run("延迟创建实际购买者", func(t *testing.T) {
		proxy := NewVirtualBuyerProxy("吴九", 150000)

		// 在调用BuyCar前，realBuyer应该是nil
		if proxy.realBuyer != nil {
			t.Error("调用方法前，realBuyer不应被创建")
		}

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err != nil {
				t.Errorf("购车应该成功，但出现错误: %v", err)
			}
		})

		// 调用BuyCar后，realBuyer应该被创建
		if proxy.realBuyer == nil {
			t.Error("调用方法后，realBuyer应该被创建")
		}

		if !strings.Contains(output, "首次调用") {
			t.Errorf("首次调用时应提示首次调用，输出: %s", output)
		}

		// 再次调用，应该复用已有对象
		output = captureOutput(func() {
			proxy.BuyCar()
		})

		if !strings.Contains(output, "复用已有") {
			t.Errorf("再次调用时应提示复用对象，输出: %s", output)
		}
	})
}

// 测试保护代理
func TestProtectionProxy(t *testing.T) {
	t.Run("VIP客户可以购车", func(t *testing.T) {
		buyer := NewRealBuyer("VIP客户", 200000)
		proxy := NewProtectionProxy(buyer, true)

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err != nil {
				t.Errorf("VIP客户应能购车，但出现错误: %v", err)
			}
		})

		if !strings.Contains(output, "VIP客户，权限验证通过") {
			t.Errorf("输出应包含权限验证通过信息，输出: %s", output)
		}

		if !strings.Contains(output, "VIP客户专享折扣") {
			t.Errorf("输出应包含VIP折扣信息，输出: %s", output)
		}
	})

	t.Run("非VIP客户无法购车", func(t *testing.T) {
		buyer := NewRealBuyer("普通客户", 200000)
		proxy := NewProtectionProxy(buyer, false)

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err == nil {
				t.Error("非VIP客户应无法购车，但没有错误返回")
			}
		})

		if !strings.Contains(output, "权限不足") {
			t.Errorf("输出应包含权限不足信息，输出: %s", output)
		}
	})

	t.Run("获取车辆信息", func(t *testing.T) {
		buyer := NewRealBuyer("测试用户", 100000)

		// VIP用户
		vipProxy := NewProtectionProxy(buyer, true)
		vipInfo := vipProxy.GetCarInfo()
		if !strings.Contains(vipInfo, "VIP专享配置") {
			t.Errorf("VIP用户应能查看专享配置，但得到: %s", vipInfo)
		}

		// 非VIP用户
		normalProxy := NewProtectionProxy(buyer, false)
		normalInfo := normalProxy.GetCarInfo()
		if !strings.Contains(normalInfo, "基础车辆信息") {
			t.Errorf("非VIP用户应只能查看基础信息，但得到: %s", normalInfo)
		}
	})
}

// 测试日志代理
func TestLoggingProxy(t *testing.T) {
	t.Run("记录成功购车操作", func(t *testing.T) {
		buyer := NewRealBuyer("日志测试", 150000)
		proxy := NewLoggingProxy(buyer)

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err != nil {
				t.Errorf("购车应该成功，但出现错误: %v", err)
			}
		})

		expectedPhrases := []string{
			"日志记录: 购车操作开始",
			"购车请求已接收",
			"购车成功",
			"操作耗时",
			"日志记录: 购车操作结束",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("输出应包含 '%s'，但未找到", phrase)
			}
		}
	})

	t.Run("记录失败购车操作", func(t *testing.T) {
		buyer := NewRealBuyer("资金不足", 50000)
		proxy := NewLoggingProxy(buyer)

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err == nil {
				t.Error("资金不足应返回错误，但没有")
			}
		})

		if !strings.Contains(output, "购车失败") {
			t.Errorf("输出应包含购车失败信息，但未找到")
		}
	})

	t.Run("记录获取车辆信息", func(t *testing.T) {
		buyer := NewRealBuyer("信息查询", 100000)
		proxy := NewLoggingProxy(buyer)

		output := captureOutput(func() {
			info := proxy.GetCarInfo()
			if info != "标准汽车型号XYZ" {
				t.Errorf("车辆信息不正确，得到: %s", info)
			}
		})

		if !strings.Contains(output, "获取车辆信息") {
			t.Errorf("输出应包含获取信息的日志，但未找到")
		}
	})
}

// 测试缓存代理
func TestCachedBuyerProxy(t *testing.T) {
	t.Run("缓存车辆信息", func(t *testing.T) {
		buyer := NewRealBuyer("缓存测试", 100000)
		proxy := NewCachedBuyerProxy(buyer)

		// 首次获取信息
		output1 := captureOutput(func() {
			info := proxy.GetCarInfo()
			if !strings.Contains(info, "标准汽车型号XYZ") {
				t.Errorf("车辆信息不正确，得到: %s", info)
			}
		})

		if !strings.Contains(output1, "首次获取车辆信息") {
			t.Errorf("首次获取应显示缓存信息，但输出: %s", output1)
		}

		// 再次获取信息，应该使用缓存
		output2 := captureOutput(func() {
			info := proxy.GetCarInfo()
			if !strings.Contains(info, "(缓存)") {
				t.Errorf("应返回缓存信息，但得到: %s", info)
			}
		})

		if !strings.Contains(output2, "从缓存获取") {
			t.Errorf("再次获取应使用缓存，但输出: %s", output2)
		}
	})

	t.Run("购车操作不支持缓存", func(t *testing.T) {
		buyer := NewRealBuyer("缓存测试", 150000)
		proxy := NewCachedBuyerProxy(buyer)

		output := captureOutput(func() {
			err := proxy.BuyCar()
			if err != nil {
				t.Errorf("购车应该成功，但出现错误: %v", err)
			}
		})

		if !strings.Contains(output, "购车操作无法缓存") {
			t.Errorf("应提示购车无法缓存，但输出: %s", output)
		}
	})
}

// 组合多个代理的测试
func TestProxyChain(t *testing.T) {
	buyer := NewRealBuyer("复合代理测试", 150000)

	// 创建代理链: 日志代理 -> 4S店代理 -> 实际购买者
	fourSProxy := NewFourSProxy(buyer)
	loggingProxy := NewLoggingProxy(fourSProxy)

	output := captureOutput(func() {
		err := loggingProxy.BuyCar()
		if err != nil {
			t.Errorf("代理链购车应该成功，但出现错误: %v", err)
		}
	})

	// 验证日志代理的输出
	if !strings.Contains(output, "日志记录: 购车操作开始") {
		t.Errorf("输出应包含日志代理信息，但未找到")
	}

	// 验证4S店代理的输出
	if !strings.Contains(output, "通过4S店代理购车开始") {
		t.Errorf("输出应包含4S店代理信息，但未找到")
	}

	// 验证实际购买者的输出
	if !strings.Contains(output, "成功购买了一辆汽车") {
		t.Errorf("输出应包含实际购买信息，但未找到")
	}
}

// 以下是示例函数，用于展示代理模式的使用方法

// 示例: 基本代理(4S店)的使用
func ExampleFourSProxy() {
	// 创建实际购买者
	buyer := NewRealBuyer("张三", 150000)

	// 创建4S店代理
	proxy := NewFourSProxy(buyer)

	// 通过代理购车
	proxy.BuyCar()

	// Output:
	// === 通过4S店代理购车开始 ===
	// 1. 从制造商订购汽车到4S店
	// 2. 准备购车文件
	// <张三> 成功购买了一辆汽车，花费了 ¥100000.00
	// 提供额外服务:
	//   1. 上牌服务
	//   2. 汽车注册
	//   3. 保险办理
	// 收取服务费: ¥5000.00
	// === 通过4S店代理购车完成 ===
}

// 示例: 虚拟代理的使用
func ExampleVirtualBuyerProxy() {
	// 创建虚拟代理，此时不会创建实际对象
	proxy := NewVirtualBuyerProxy("李四", 200000)

	// 第一次调用时创建实际对象
	proxy.BuyCar()

	// 第二次调用时复用已有对象
	proxy.BuyCar()

	// Output:
	// === 通过虚拟代理购车开始 ===
	// 准备创建实际购买者...
	// 首次调用，创建实际购买者
	// <李四> 成功购买了一辆汽车，花费了 ¥100000.00
	// === 通过虚拟代理购车结束 ===
	// === 通过虚拟代理购车开始 ===
	// 准备创建实际购买者...
	// 复用已有的实际购买者
	// <李四> 成功购买了一辆汽车，花费了 ¥100000.00
	// === 通过虚拟代理购车结束 ===
}

// 示例: 保护代理的使用
func ExampleProtectionProxy() {
	// 创建实际购买者
	buyer := NewRealBuyer("王五", 300000)

	// 创建VIP保护代理
	vipProxy := NewProtectionProxy(buyer, true)
	fmt.Println("--- VIP客户尝试购车 ---")
	vipProxy.BuyCar()

	// 创建普通保护代理
	normalProxy := NewProtectionProxy(buyer, false)
	fmt.Println("\n--- 普通客户尝试购车 ---")
	normalProxy.BuyCar()

	// Output:
	// --- VIP客户尝试购车 ---
	// === 通过保护代理购车开始 ===
	// VIP客户，权限验证通过
	// <王五> 成功购买了一辆汽车，花费了 ¥100000.00
	// VIP客户专享折扣已应用
	// === 通过保护代理购车结束 ===
	//
	// --- 普通客户尝试购车 ---
	// === 通过保护代理购车开始 ===
	// 权限不足: 仅VIP客户可以通过此渠道购车
}

// 示例: 缓存代理的使用
func ExampleCachedBuyerProxy() {
	// 创建实际购买者
	buyer := NewRealBuyer("赵六", 150000)

	// 创建缓存代理
	proxy := NewCachedBuyerProxy(buyer)

	// 获取车辆信息 - 第一次，将会缓存结果
	fmt.Println("第一次获取车辆信息:")
	info1 := proxy.GetCarInfo()
	fmt.Println("结果:", info1)

	// 获取车辆信息 - 第二次，将使用缓存
	fmt.Println("\n第二次获取车辆信息:")
	info2 := proxy.GetCarInfo()
	fmt.Println("结果:", info2)

	// Output:
	// 第一次获取车辆信息:
	// 首次获取车辆信息，将结果缓存
	// 结果: 标准汽车型号XYZ
	//
	// 第二次获取车辆信息:
	// 从缓存获取车辆信息
	// 结果: 标准汽车型号XYZ (缓存)
}

// 示例: 组合多个代理
func DemoProxyChain() {
	// 创建实际购买者
	buyer := NewRealBuyer("复合代理客户", 200000)

	// 创建代理链：缓存代理 -> 保护代理 -> 日志代理 -> 4S店代理 -> 实际购买者
	fourSProxy := NewFourSProxy(buyer)
	loggingProxy := NewLoggingProxy(fourSProxy)
	protectionProxy := NewProtectionProxy(loggingProxy, true)
	cachedProxy := NewCachedBuyerProxy(protectionProxy)

	// 通过代理链获取车辆信息（第一次）
	fmt.Println("=== 第一次获取车辆信息 ===")
	info1 := cachedProxy.GetCarInfo()
	fmt.Println("结果:", info1)

	// 通过代理链获取车辆信息（第二次，使用缓存）
	fmt.Println("\n=== 第二次获取车辆信息 ===")
	info2 := cachedProxy.GetCarInfo()
	fmt.Println("结果:", info2)

	// 通过代理链购车
	fmt.Println("\n=== 通过代理链购车 ===")
	cachedProxy.BuyCar()

	// 不会生成完全一致的Output，因为有时间戳
	// 所以这里仅作为示例，不作为测试
}
