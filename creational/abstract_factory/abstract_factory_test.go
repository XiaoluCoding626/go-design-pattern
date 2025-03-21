package abstractfactory

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

	// 执行函数
	f()

	// 恢复原始的标准输出
	w.Close()
	os.Stdout = oldStdout

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

// 测试木门
func TestWoodenDoor(t *testing.T) {
	door := &WoodenDoor{}

	// 测试开门
	output := captureOutput(door.Open)
	if !strings.Contains(output, "木门打开") {
		t.Errorf("WoodenDoor.Open() 输出 = %q, 应包含 '木门打开'", output)
	}

	// 测试关门
	output = captureOutput(door.Close)
	if !strings.Contains(output, "木门关闭") {
		t.Errorf("WoodenDoor.Close() 输出 = %q, 应包含 '木门关闭'", output)
	}

	// 测试材质
	if material := door.GetMaterial(); material != "实木材质" {
		t.Errorf("WoodenDoor.GetMaterial() = %q, 期望 '实木材质'", material)
	}
}

// 测试木门把手
func TestWoodenDoorHandle(t *testing.T) {
	handle := &WoodenDoorHandle{}

	// 测试按下
	output := captureOutput(handle.Press)
	if !strings.Contains(output, "按下木门把手") {
		t.Errorf("WoodenDoorHandle.Press() 输出 = %q, 应包含 '按下木门把手'", output)
	}

	// 测试拉动
	output = captureOutput(handle.Pull)
	if !strings.Contains(output, "拉动木门把手") {
		t.Errorf("WoodenDoorHandle.Pull() 输出 = %q, 应包含 '拉动木门把手'", output)
	}

	// 测试材质
	if material := handle.GetMaterial(); material != "实木材质" {
		t.Errorf("WoodenDoorHandle.GetMaterial() = %q, 期望 '实木材质'", material)
	}
}

// 测试木门锁
func TestWoodenDoorLock(t *testing.T) {
	lock := &WoodenDoorLock{}

	// 测试上锁
	output := captureOutput(lock.Lock)
	if !strings.Contains(output, "锁上木门锁") {
		t.Errorf("WoodenDoorLock.Lock() 输出 = %q, 应包含 '锁上木门锁'", output)
	}

	// 测试解锁
	output = captureOutput(lock.Unlock)
	if !strings.Contains(output, "解锁木门锁") {
		t.Errorf("WoodenDoorLock.Unlock() 输出 = %q, 应包含 '解锁木门锁'", output)
	}

	// 测试安全级别
	if level := lock.GetSecurityLevel(); level != 1 {
		t.Errorf("WoodenDoorLock.GetSecurityLevel() = %d, 期望 1", level)
	}
}

// 测试金属门
func TestMetalDoor(t *testing.T) {
	door := &MetalDoor{}

	// 测试开门
	output := captureOutput(door.Open)
	if !strings.Contains(output, "金属门打开") {
		t.Errorf("MetalDoor.Open() 输出 = %q, 应包含 '金属门打开'", output)
	}

	// 测试关门
	output = captureOutput(door.Close)
	if !strings.Contains(output, "金属门关闭") {
		t.Errorf("MetalDoor.Close() 输出 = %q, 应包含 '金属门关闭'", output)
	}

	// 测试材质
	if material := door.GetMaterial(); material != "钢铁材质" {
		t.Errorf("MetalDoor.GetMaterial() = %q, 期望 '钢铁材质'", material)
	}
}

// 测试金属门把手
func TestMetalDoorHandle(t *testing.T) {
	handle := &MetalDoorHandle{}

	// 测试按下
	output := captureOutput(handle.Press)
	if !strings.Contains(output, "按下金属门把手") {
		t.Errorf("MetalDoorHandle.Press() 输出 = %q, 应包含 '按下金属门把手'", output)
	}

	// 测试拉动
	output = captureOutput(handle.Pull)
	if !strings.Contains(output, "拉动金属门把手") {
		t.Errorf("MetalDoorHandle.Pull() 输出 = %q, 应包含 '拉动金属门把手'", output)
	}

	// 测试材质
	if material := handle.GetMaterial(); material != "不锈钢材质" {
		t.Errorf("MetalDoorHandle.GetMaterial() = %q, 期望 '不锈钢材质'", material)
	}
}

// 测试金属门锁
func TestMetalDoorLock(t *testing.T) {
	lock := &MetalDoorLock{}

	// 测试上锁
	output := captureOutput(lock.Lock)
	if !strings.Contains(output, "锁上金属安全锁") {
		t.Errorf("MetalDoorLock.Lock() 输出 = %q, 应包含 '锁上金属安全锁'", output)
	}

	// 测试解锁
	output = captureOutput(lock.Unlock)
	if !strings.Contains(output, "解锁金属安全锁") {
		t.Errorf("MetalDoorLock.Unlock() 输出 = %q, 应包含 '解锁金属安全锁'", output)
	}

	// 测试安全级别
	if level := lock.GetSecurityLevel(); level != 3 {
		t.Errorf("MetalDoorLock.GetSecurityLevel() = %d, 期望 3", level)
	}
}

// 测试玻璃门
func TestGlassDoor(t *testing.T) {
	door := &GlassDoor{}

	// 测试开门
	output := captureOutput(door.Open)
	if !strings.Contains(output, "玻璃门滑动") {
		t.Errorf("GlassDoor.Open() 输出 = %q, 应包含 '玻璃门滑动'", output)
	}

	// 测试关门
	output = captureOutput(door.Close)
	if !strings.Contains(output, "玻璃门平稳") {
		t.Errorf("GlassDoor.Close() 输出 = %q, 应包含 '玻璃门平稳'", output)
	}

	// 测试材质
	if material := door.GetMaterial(); material != "钢化玻璃材质" {
		t.Errorf("GlassDoor.GetMaterial() = %q, 期望 '钢化玻璃材质'", material)
	}
}

// 测试玻璃门把手
func TestGlassDoorHandle(t *testing.T) {
	handle := &GlassDoorHandle{}

	// 测试按下
	output := captureOutput(handle.Press)
	if !strings.Contains(output, "按下玻璃门把手") {
		t.Errorf("GlassDoorHandle.Press() 输出 = %q, 应包含 '按下玻璃门把手'", output)
	}

	// 测试拉动
	output = captureOutput(handle.Pull)
	if !strings.Contains(output, "拉动玻璃门把手") {
		t.Errorf("GlassDoorHandle.Pull() 输出 = %q, 应包含 '拉动玻璃门把手'", output)
	}

	// 测试材质
	if material := handle.GetMaterial(); material != "铝合金材质" {
		t.Errorf("GlassDoorHandle.GetMaterial() = %q, 期望 '铝合金材质'", material)
	}
}

// 测试玻璃门锁
func TestGlassDoorLock(t *testing.T) {
	lock := &GlassDoorLock{}

	// 测试上锁
	output := captureOutput(lock.Lock)
	if !strings.Contains(output, "锁上玻璃门电子锁") {
		t.Errorf("GlassDoorLock.Lock() 输出 = %q, 应包含 '锁上玻璃门电子锁'", output)
	}

	// 测试解锁
	output = captureOutput(lock.Unlock)
	if !strings.Contains(output, "解锁玻璃门电子锁") {
		t.Errorf("GlassDoorLock.Unlock() 输出 = %q, 应包含 '解锁玻璃门电子锁'", output)
	}

	// 测试安全级别
	if level := lock.GetSecurityLevel(); level != 2 {
		t.Errorf("GlassDoorLock.GetSecurityLevel() = %d, 期望 2", level)
	}
}

// 测试木门工厂
func TestWoodenDoorFactory(t *testing.T) {
	factory := &WoodenDoorFactory{}

	// 测试创建门
	door := factory.CreateDoor()
	if _, ok := door.(*WoodenDoor); !ok {
		t.Error("WoodenDoorFactory.CreateDoor() 应该返回 *WoodenDoor")
	}

	// 测试创建门把手
	handle := factory.CreateDoorHandle()
	if _, ok := handle.(*WoodenDoorHandle); !ok {
		t.Error("WoodenDoorFactory.CreateDoorHandle() 应该返回 *WoodenDoorHandle")
	}

	// 测试创建门锁
	lock := factory.CreateDoorLock()
	if _, ok := lock.(*WoodenDoorLock); !ok {
		t.Error("WoodenDoorFactory.CreateDoorLock() 应该返回 *WoodenDoorLock")
	}
}

// 测试金属门工厂
func TestMetalDoorFactory(t *testing.T) {
	factory := &MetalDoorFactory{}

	// 测试创建门
	door := factory.CreateDoor()
	if _, ok := door.(*MetalDoor); !ok {
		t.Error("MetalDoorFactory.CreateDoor() 应该返回 *MetalDoor")
	}

	// 测试创建门把手
	handle := factory.CreateDoorHandle()
	if _, ok := handle.(*MetalDoorHandle); !ok {
		t.Error("MetalDoorFactory.CreateDoorHandle() 应该返回 *MetalDoorHandle")
	}

	// 测试创建门锁
	lock := factory.CreateDoorLock()
	if _, ok := lock.(*MetalDoorLock); !ok {
		t.Error("MetalDoorFactory.CreateDoorLock() 应该返回 *MetalDoorLock")
	}
}

// 测试玻璃门工厂
func TestGlassDoorFactory(t *testing.T) {
	factory := &GlassDoorFactory{}

	// 测试创建门
	door := factory.CreateDoor()
	if _, ok := door.(*GlassDoor); !ok {
		t.Error("GlassDoorFactory.CreateDoor() 应该返回 *GlassDoor")
	}

	// 测试创建门把手
	handle := factory.CreateDoorHandle()
	if _, ok := handle.(*GlassDoorHandle); !ok {
		t.Error("GlassDoorFactory.CreateDoorHandle() 应该返回 *GlassDoorHandle")
	}

	// 测试创建门锁
	lock := factory.CreateDoorLock()
	if _, ok := lock.(*GlassDoorLock); !ok {
		t.Error("GlassDoorFactory.CreateDoorLock() 应该返回 *GlassDoorLock")
	}
}

// 测试GetDoorFactory函数
func TestGetDoorFactory(t *testing.T) {
	// 测试获取木门工厂
	factory1, err := GetDoorFactory(WoodenType)
	if err != nil {
		t.Errorf("GetDoorFactory(WoodenType) 返回错误: %v", err)
	}

	_, ok := factory1.(*WoodenDoorFactory)
	if !ok {
		t.Error("GetDoorFactory(WoodenType) 应该返回 *WoodenDoorFactory")
	}

	// 测试获取金属门工厂
	factory2, err := GetDoorFactory(MetalType)
	if err != nil {
		t.Errorf("GetDoorFactory(MetalType) 返回错误: %v", err)
	}

	_, ok = factory2.(*MetalDoorFactory)
	if !ok {
		t.Error("GetDoorFactory(MetalType) 应该返回 *MetalDoorFactory")
	}

	// 测试获取玻璃门工厂
	factory3, err := GetDoorFactory(GlassType)
	if err != nil {
		t.Errorf("GetDoorFactory(GlassType) 返回错误: %v", err)
	}

	_, ok = factory3.(*GlassDoorFactory)
	if !ok {
		t.Error("GetDoorFactory(GlassType) 应该返回 *GlassDoorFactory")
	}

	// 测试单例模式
	factory1Again, _ := GetDoorFactory(WoodenType)
	if factory1 != factory1Again {
		t.Error("单例模式失败：GetDoorFactory(WoodenType) 返回的不是同一个实例")
	}

	// 测试不支持的门类型
	_, err = GetDoorFactory("unsupported")
	if err == nil {
		t.Error("GetDoorFactory('unsupported') 应该返回错误")
	}

	if !strings.Contains(err.Error(), "不支持的门类型") {
		t.Errorf("错误消息 = %q, 应该包含 '不支持的门类型'", err.Error())
	}
}

// 测试DoorCreator
func TestDoorCreator(t *testing.T) {
	// 测试创建木门创建器
	creator1, err := NewDoorCreator(WoodenType)
	if err != nil {
		t.Errorf("NewDoorCreator(WoodenType) 返回错误: %v", err)
	}

	// 测试创建完整的门系统
	door1, handle1, lock1 := creator1.CreateCompleteDoor()

	_, ok := door1.(*WoodenDoor)
	if !ok {
		t.Error("木门创建器创建的门应该是 *WoodenDoor")
	}

	_, ok = handle1.(*WoodenDoorHandle)
	if !ok {
		t.Error("木门创建器创建的把手应该是 *WoodenDoorHandle")
	}

	_, ok = lock1.(*WoodenDoorLock)
	if !ok {
		t.Error("木门创建器创建的锁应该是 *WoodenDoorLock")
	}

	// 测试创建金属门创建器
	creator2, _ := NewDoorCreator(MetalType)
	door2, handle2, lock2 := creator2.CreateCompleteDoor()

	_, ok = door2.(*MetalDoor)
	if !ok {
		t.Error("金属门创建器创建的门应该是 *MetalDoor")
	}

	_, ok = handle2.(*MetalDoorHandle)
	if !ok {
		t.Error("金属门创建器创建的把手应该是 *MetalDoorHandle")
	}

	_, ok = lock2.(*MetalDoorLock)
	if !ok {
		t.Error("金属门创建器创建的锁应该是 *MetalDoorLock")
	}

	// 测试创建玻璃门创建器
	creator3, _ := NewDoorCreator(GlassType)
	door3, handle3, lock3 := creator3.CreateCompleteDoor()

	_, ok = door3.(*GlassDoor)
	if !ok {
		t.Error("玻璃门创建器创建的门应该是 *GlassDoor")
	}

	_, ok = handle3.(*GlassDoorHandle)
	if !ok {
		t.Error("玻璃门创建器创建的把手应该是 *GlassDoorHandle")
	}

	_, ok = lock3.(*GlassDoorLock)
	if !ok {
		t.Error("玻璃门创建器创建的锁应该是 *GlassDoorLock")
	}

	// 测试创建不支持的门类型创建器
	_, err = NewDoorCreator("unsupported")
	if err == nil {
		t.Error("NewDoorCreator('unsupported') 应该返回错误")
	}
}

// 集成测试场景
func TestIntegrationScenario(t *testing.T) {
	// 1. 使用工厂创建木门产品族并验证行为
	woodenCreator, _ := NewDoorCreator(WoodenType)
	woodenDoor, woodenHandle, woodenLock := woodenCreator.CreateCompleteDoor()

	output := captureOutput(func() {
		woodenDoor.Open()
		woodenHandle.Press()
		woodenLock.Lock()
	})

	if !strings.Contains(output, "木门打开") ||
		!strings.Contains(output, "按下木门把手") ||
		!strings.Contains(output, "锁上木门锁") {
		t.Error("木门产品族的行为不正确")
	}

	// 2. 使用工厂创建金属门产品族并验证行为
	metalCreator, _ := NewDoorCreator(MetalType)
	metalDoor, metalHandle, metalLock := metalCreator.CreateCompleteDoor()

	output = captureOutput(func() {
		metalDoor.Close()
		metalHandle.Pull()
		metalLock.Unlock()
	})

	if !strings.Contains(output, "金属门关闭") ||
		!strings.Contains(output, "拉动金属门把手") ||
		!strings.Contains(output, "解锁金属安全锁") {
		t.Error("金属门产品族的行为不正确")
	}

	// 3. 验证安全级别
	if woodenLock.GetSecurityLevel() >= metalLock.GetSecurityLevel() {
		t.Error("木门锁的安全级别不应高于或等于金属门锁")
	}
}
