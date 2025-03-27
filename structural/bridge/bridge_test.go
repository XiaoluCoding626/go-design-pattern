package bridge

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

// 测试设备接口实现
func TestDevices(t *testing.T) {
	// 创建并测试电视机
	t.Run("TV Device", func(t *testing.T) {
		assert := assert.New(t)
		tv := NewTV("Sony")

		assert.Equal("Sony", tv.GetName())

		// 测试开启电视
		output := captureOutput(func() {
			tv.TurnOn()
		})
		assert.Contains(output, "Sony 电视机打开了")
		assert.Contains(output, "当前音量：10")

		// 测试设置音量
		output = captureOutput(func() {
			tv.SetVolume(50)
		})
		assert.Contains(output, "Sony 电视机音量设置为：50")

		// 测试关闭电视
		output = captureOutput(func() {
			tv.TurnOff()
		})
		assert.Contains(output, "Sony 电视机关闭了")
	})

	// 创建并测试收音机
	t.Run("Radio Device", func(t *testing.T) {
		assert := assert.New(t)
		radio := NewRadio("Philips")

		assert.Equal("Philips", radio.GetName())

		// 测试开启收音机
		output := captureOutput(func() {
			radio.TurnOn()
		})
		assert.Contains(output, "Philips 收音机打开了")
		assert.Contains(output, "当前音量：5")

		// 测试设置音量
		output = captureOutput(func() {
			radio.SetVolume(25)
		})
		assert.Contains(output, "Philips 收音机音量设置为：25")

		// 测试关闭收音机
		output = captureOutput(func() {
			radio.TurnOff()
		})
		assert.Contains(output, "Philips 收音机关闭了")
	})
}

// 测试标准遥控器
func TestStandardRemoteControl(t *testing.T) {
	tv := NewTV("Samsung")
	radio := NewRadio("JBL")

	// 使用标准遥控器控制电视
	t.Run("Standard Remote with TV", func(t *testing.T) {
		assert := assert.New(t)
		remote := NewStandardRemoteControl(tv)

		// 测试开机
		output := captureOutput(func() {
			remote.PowerOn()
		})
		assert.Contains(output, "Samsung 电视机打开了")

		// 测试提高音量
		output = captureOutput(func() {
			remote.VolumeUp()
		})
		assert.Contains(output, "Samsung 电视机音量设置为：20")

		// 测试降低音量
		output = captureOutput(func() {
			remote.VolumeDown()
		})
		assert.Contains(output, "Samsung 电视机音量设置为：10")

		// 测试关机
		output = captureOutput(func() {
			remote.PowerOff()
		})
		assert.Contains(output, "Samsung 电视机关闭了")
	})

	// 使用标准遥控器控制收音机
	t.Run("Standard Remote with Radio", func(t *testing.T) {
		assert := assert.New(t)
		remote := NewStandardRemoteControl(radio)

		// 测试开机
		output := captureOutput(func() {
			remote.PowerOn()
		})
		assert.Contains(output, "JBL 收音机打开了")

		// 测试提高音量
		output = captureOutput(func() {
			remote.VolumeUp()
		})
		assert.Contains(output, "JBL 收音机音量设置为：20")

		// 测试降低音量
		output = captureOutput(func() {
			remote.VolumeDown()
		})
		assert.Contains(output, "JBL 收音机音量设置为：10")

		// 测试关机
		output = captureOutput(func() {
			remote.PowerOff()
		})
		assert.Contains(output, "JBL 收音机关闭了")
	})
}

// 测试高级遥控器
func TestAdvancedRemoteControl(t *testing.T) {
	tv := NewTV("LG")
	radio := NewRadio("Bose")

	// 使用高级遥控器控制电视
	t.Run("Advanced Remote with TV", func(t *testing.T) {
		assert := assert.New(t)
		remote := NewAdvancedRemoteControl(tv)

		// 测试基本功能
		output := captureOutput(func() {
			remote.PowerOn()
		})
		assert.Contains(output, "LG 电视机打开了")

		// 测试静音（高级功能）
		output = captureOutput(func() {
			remote.Mute()
		})
		assert.Contains(output, "LG 电视机音量设置为：0")
		assert.Contains(output, "静音 LG")

		// 测试最大音量（高级功能）
		output = captureOutput(func() {
			remote.MaxVolume()
		})
		assert.Contains(output, "LG 电视机音量设置为：100")
		assert.Contains(output, "将 LG 音量调到最大")

		// 测试关机
		output = captureOutput(func() {
			remote.PowerOff()
		})
		assert.Contains(output, "LG 电视机关闭了")
	})

	// 使用高级遥控器控制收音机
	t.Run("Advanced Remote with Radio", func(t *testing.T) {
		assert := assert.New(t)
		remote := NewAdvancedRemoteControl(radio)

		// 测试基本功能
		output := captureOutput(func() {
			remote.PowerOn()
		})
		assert.Contains(output, "Bose 收音机打开了")

		// 测试静音（高级功能）
		output = captureOutput(func() {
			remote.Mute()
		})
		assert.Contains(output, "Bose 收音机音量设置为：0")
		assert.Contains(output, "静音 Bose")

		// 测试最大音量（高级功能）
		output = captureOutput(func() {
			remote.MaxVolume()
		})
		assert.Contains(output, "Bose 收音机音量设置为：100")
		assert.Contains(output, "将 Bose 音量调到最大")

		// 测试关机
		output = captureOutput(func() {
			remote.PowerOff()
		})
		assert.Contains(output, "Bose 收音机关闭了")
	})
}

// 测试桥接模式的核心特性：设备和遥控器可以独立变化
func TestBridgePattern(t *testing.T) {
	assert := assert.New(t)
	devices := []Device{
		NewTV("TCL"),
		NewRadio("Sony"),
	}

	// 同一个遥控器可以控制不同类型的设备
	t.Run("Same remote with different devices", func(t *testing.T) {
		for _, device := range devices {
			remote := NewStandardRemoteControl(device)
			name := device.GetName()

			output := captureOutput(func() {
				remote.PowerOn()
				remote.VolumeUp()
				remote.PowerOff()
			})

			assert.Contains(output, name)
			assert.Contains(output, "打开了")
			assert.Contains(output, "音量设置为")
			assert.Contains(output, "关闭了")
		}
	})

	// 同一个设备可以被不同类型的遥控器控制
	t.Run("Same device with different remotes", func(t *testing.T) {
		tv := NewTV("Sharp")
		remotes := []RemoteControl{
			NewStandardRemoteControl(tv),
			NewAdvancedRemoteControl(tv),
		}

		for i, remote := range remotes {
			output := captureOutput(func() {
				remote.PowerOn()
				remote.VolumeUp()
				remote.PowerOff()
			})

			assert.Contains(output, "Sharp")
			assert.Contains(output, "电视机打开了")
			assert.Contains(output, "音量设置为")
			assert.Contains(output, "电视机关闭了")

			// 验证独立变化
			if i == 1 {
				advancedRemote := remote.(*AdvancedRemoteControl)
				output := captureOutput(func() {
					advancedRemote.Mute()
				})
				assert.Contains(output, "静音 Sharp")
			}
		}
	})
}

// 示例测试
func ExampleStandardRemoteControl() {
	tv := NewTV("Example TV")
	remote := NewStandardRemoteControl(tv)

	remote.PowerOn()
	remote.VolumeUp()
	remote.PowerOff()

	// Output:
	// Example TV 电视机打开了，当前音量：10
	// Example TV 电视机音量设置为：20
	// Example TV 电视机关闭了
}

func ExampleAdvancedRemoteControl() {
	radio := NewRadio("Example Radio")
	remote := NewAdvancedRemoteControl(radio)

	remote.PowerOn()
	remote.Mute()
	remote.MaxVolume()
	remote.PowerOff()

	// Output:
	// Example Radio 收音机打开了，当前音量：5
	// Example Radio 收音机音量设置为：0
	// 静音 Example Radio
	// Example Radio 收音机音量设置为：100
	// 将 Example Radio 音量调到最大
	// Example Radio 收音机关闭了
}
