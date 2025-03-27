package adapter

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 捕获标准输出的辅助函数
func captureOutput(f func()) string {
	// 保存原始的标准输出
	oldStdout := os.Stdout

	// 创建一个管道
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行函数
	f()

	// 恢复原始的标准输出
	os.Stdout = oldStdout
	w.Close()

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

// 测试插头类型
func TestPlugs(t *testing.T) {
	assert := assert.New(t)

	// 测试两针插头
	twoPin := &TwoPinPlug{}
	assert.Equal(2, twoPin.GetPin(), "两针插头的针数应为2")

	// 测试三针插头
	threePin := &ThreePinPlug{}
	assert.Equal(3, threePin.GetPin(), "三针插头的针数应为3")
}

// 测试三孔插座
func TestThreePinSocket(t *testing.T) {
	assert := assert.New(t)
	socket := &ThreePinSocket{}
	twoPin := &TwoPinPlug{}
	threePin := &ThreePinPlug{}

	// 测试三孔插座拒绝两针插头
	output := captureOutput(func() {
		socket.Charge(twoPin)
	})
	assert.Contains(output, "三孔插座无法为非三针插头充电", "三孔插座应该拒绝两针插头")

	// 测试三孔插座接受三针插头
	output = captureOutput(func() {
		socket.Charge(threePin)
	})
	assert.Contains(output, "三孔插座正在为三针插头充电", "三孔插座应该接受三针插头")
}

// 测试两孔插座
func TestTwoPinSocket(t *testing.T) {
	assert := assert.New(t)
	socket := &TwoPinSocket{}
	twoPin := &TwoPinPlug{}
	threePin := &ThreePinPlug{}

	// 测试两孔插座接受两针插头
	output := captureOutput(func() {
		socket.Charge(twoPin)
	})
	assert.Contains(output, "两孔插座正在为两针插头充电", "两孔插座应该接受两针插头")

	// 测试两孔插座拒绝三针插头
	output = captureOutput(func() {
		socket.Charge(threePin)
	})
	assert.Contains(output, "两孔插座无法为非两针插头充电", "两孔插座应该拒绝三针插头")
}

// 测试电源适配器 - 适配器模式的核心测试
func TestPowerAdapter(t *testing.T) {
	assert := assert.New(t)
	adapter := NewPowerAdapter()
	twoPin := &TwoPinPlug{}
	threePin := &ThreePinPlug{}

	// 测试适配器接受两针插头
	output := captureOutput(func() {
		adapter.Charge(twoPin)
	})
	assert.Contains(output, "适配器正在将两针插头转换为三针插头", "适配器应该接受两针插头并进行转换")
	assert.Contains(output, "三孔插座正在为三针插头充电", "适配后应该使用三孔插座充电")
	assert.Contains(output, "适配完成，两针插头已成功充电", "适配器应该完成充电过程")

	// 测试适配器拒绝三针插头
	output = captureOutput(func() {
		adapter.Charge(threePin)
	})
	assert.Contains(output, "适配器只能适配两针插头", "适配器应该拒绝三针插头")
}

// 测试适配器作为电源插座接口的实现
func TestAdapterAsIPowerSocket(t *testing.T) {
	assert := assert.New(t)
	var socket IPowerSocket = NewPowerAdapter() // 多态测试
	twoPin := &TwoPinPlug{}

	output := captureOutput(func() {
		socket.Charge(twoPin)
	})
	assert.Contains(output, "适配器正在将两针插头转换为三针插头", "通过IPowerSocket接口，适配器应该正常工作")
	assert.Contains(output, "三孔插座正在为三针插头充电", "通过IPowerSocket接口，适配器应该使三孔插座正常工作")
}

// 集成测试：模拟现实场景
func TestRealWorldScenario(t *testing.T) {
	// 创建设备
	twoPin := &TwoPinPlug{}
	threePin := &ThreePinPlug{}
	threePinSocket := &ThreePinSocket{}
	adapter := NewPowerAdapter()

	t.Run("不使用适配器的情况", func(t *testing.T) {
		// 对每个子测试使用新的断言对象，避免变量遮蔽问题
		a := assert.New(t)

		// 三针插头可以直接使用三孔插座
		output := captureOutput(func() {
			threePinSocket.Charge(threePin)
		})
		a.Contains(output, "三孔插座正在为三针插头充电", "三针插头应该可以直接使用三孔插座")

		// 两针插头不能直接使用三孔插座
		output = captureOutput(func() {
			threePinSocket.Charge(twoPin)
		})
		a.Contains(output, "三孔插座无法为非三针插头充电", "两针插头不应该可以直接使用三孔插座")
	})

	t.Run("使用适配器的情况", func(t *testing.T) {
		// 对每个子测试使用新的断言对象，避免变量遮蔽问题
		a := assert.New(t)

		// 两针插头通过适配器可以使用三孔插座提供的电
		output := captureOutput(func() {
			adapter.Charge(twoPin)
		})
		a.Contains(output, "适配器正在将两针插头转换为三针插头", "应显示适配过程")
		a.Contains(output, "三孔插座正在为三针插头充电", "应使用三孔插座充电")
		a.Contains(output, "适配完成，两针插头已成功充电", "充电应成功完成")
	})
}

// 示例测试 - 展示适配器模式的使用方法
func ExamplePowerAdapter() {
	// 创建组件
	twoPin := &TwoPinPlug{}
	threePinSocket := &ThreePinSocket{}
	adapter := NewPowerAdapter()

	// 尝试直接使用三孔插座 - 失败
	threePinSocket.Charge(twoPin)

	// 使用适配器 - 成功
	adapter.Charge(twoPin)

	// Output:
	// 三孔插座无法为非三针插头充电
	// 适配器正在将两针插头转换为三针插头
	// 三孔插座正在为三针插头充电
	// 适配完成，两针插头已成功充电
}
