package template_method

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 辅助函数，用于捕获标准输出
func captureOutput(f func()) string {
	// 保存原始的标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行函数
	f()

	// 恢复原始的标准输出并获取输出内容
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

// 测试红豆豆浆的制作过程
func TestRedBeanSoyaMilk(t *testing.T) {
	milk := NewRedBeanSoyaMilk()
	output := captureOutput(func() {
		milk.Make()
	})

	// 验证输出中包含必要的步骤
	assert.Contains(t, output, "开始制作豆浆")
	assert.Contains(t, output, "第 1 步：选择新鲜的豆子")
	assert.Contains(t, output, "第 2 步：加入上好的红豆")
	assert.Contains(t, output, "第 3 步：豆子和配料开始浸泡")
	assert.Contains(t, output, "第 4 步：豆子和配料放入豆浆机榨汁")
	assert.Contains(t, output, "豆浆制作完成")

	// 验证步骤顺序正确
	steps := []string{
		"开始制作豆浆",
		"第 1 步",
		"第 2 步",
		"第 3 步",
		"第 4 步",
		"豆浆制作完成",
	}

	lastIndex := 0
	for _, step := range steps {
		index := strings.Index(output, step)
		assert.True(t, index >= lastIndex, "步骤顺序不正确: %s", step)
		lastIndex = index
	}
}

// 测试花生豆浆的制作过程，包括额外的钩子方法
func TestPeanutSoyaMilk(t *testing.T) {
	milk := NewPeanutSoyaMilk()
	output := captureOutput(func() {
		milk.Make()
	})

	// 添加测试输出调试，帮助查看实际输出
	t.Logf("花生豆浆实际输出:\n%s", output)

	// 验证基本步骤
	assert.Contains(t, output, "开始制作豆浆")
	assert.Contains(t, output, "第 1 步：选择新鲜的豆子")
	assert.Contains(t, output, "第 2 步：加入上好的花生")
	assert.Contains(t, output, "第 3 步：豆子和配料开始浸泡")
	assert.Contains(t, output, "第 4 步：豆子和配料放入豆浆机榨汁")
	assert.Contains(t, output, "第 5 步：花生豆浆完成后撒一些花生碎")
	assert.Contains(t, output, "豆浆制作完成")
}

// 测试纯豆浆，验证钩子方法能正确跳过添加配料步骤
func TestPureSoyaMilk(t *testing.T) {
	milk := NewPureSoyaMilk()
	output := captureOutput(func() {
		milk.Make()
	})

	// 验证基本步骤
	assert.Contains(t, output, "开始制作豆浆")
	assert.Contains(t, output, "第 1 步：选择新鲜的豆子")
	assert.Contains(t, output, "第 3 步：豆子和配料开始浸泡")
	assert.Contains(t, output, "第 4 步：豆子和配料放入豆浆机榨汁")
	assert.Contains(t, output, "豆浆制作完成")

	// 验证没有添加配料的步骤
	assert.NotContains(t, output, "第 2 步：加入")
}

// 示例函数，展示模板方法模式的使用
func Example() {
	// 创建并制作红豆豆浆
	fmt.Println("制作红豆豆浆：")
	redBeanMilk := NewRedBeanSoyaMilk()
	redBeanMilk.Make()

	fmt.Println()

	// 创建并制作花生豆浆
	fmt.Println("制作花生豆浆：")
	peanutMilk := NewPeanutSoyaMilk()
	peanutMilk.Make()

	fmt.Println()

	// 创建并制作纯豆浆
	fmt.Println("制作纯豆浆：")
	pureMilk := NewPureSoyaMilk()
	pureMilk.Make()

	// Output:
	// 制作红豆豆浆：
	// === 开始制作豆浆 ===
	// 第 1 步：选择新鲜的豆子
	// 第 2 步：加入上好的红豆
	// 第 3 步：豆子和配料开始浸泡 3 小时
	// 第 4 步：豆子和配料放入豆浆机榨汁
	// === 豆浆制作完成 ===
	//
	// 制作花生豆浆：
	// === 开始制作豆浆 ===
	// 第 1 步：选择新鲜的豆子
	// 第 2 步：加入上好的花生
	// 第 3 步：豆子和配料开始浸泡 3 小时
	// 第 4 步：豆子和配料放入豆浆机榨汁
	// 第 5 步：花生豆浆完成后撒一些花生碎
	// === 豆浆制作完成 ===
	//
	// 制作纯豆浆：
	// === 开始制作豆浆 ===
	// 第 1 步：选择新鲜的豆子
	// 第 3 步：豆子和配料开始浸泡 3 小时
	// 第 4 步：豆子和配料放入豆浆机榨汁
	// === 豆浆制作完成 ===
}

// 测试自定义钩子方法
type CustomSoyaMilk struct {
	AbstractSoyaMilk
	wantsCondiments bool
}

func NewCustomSoyaMilk(wantsCondiments bool) *CustomSoyaMilk {
	milk := &CustomSoyaMilk{
		wantsCondiments: wantsCondiments,
	}
	milk.soyaMilkBehavior = milk // 使用新的接口名称
	return milk
}

func (c *CustomSoyaMilk) AddCondiment() {
	fmt.Println("第 2 步：加入自定义配料")
}

// 覆盖钩子方法，根据wantsCondiments决定是否添加配料
func (c *CustomSoyaMilk) CustomerWantsCondiments() bool {
	return c.wantsCondiments
}

// 实现Hook方法，完成接口
func (c *CustomSoyaMilk) Hook() {
	// 空实现
}

// 修复测试 - 当参数为false时不应添加配料
func TestCustomHook(t *testing.T) {
	// 测试需要配料的情况
	withCondiments := NewCustomSoyaMilk(true)
	output1 := captureOutput(func() {
		withCondiments.Make()
	})
	t.Logf("需要配料的实际输出:\n%s", output1)
	assert.Contains(t, output1, "第 2 步：加入自定义配料")

	// 测试不需要配料的情况
	withoutCondiments := NewCustomSoyaMilk(false)
	output2 := captureOutput(func() {
		withoutCondiments.Make()
	})
	t.Logf("不需要配料的实际输出:\n%s", output2)
	assert.NotContains(t, output2, "第 2 步：加入自定义配料")
}
