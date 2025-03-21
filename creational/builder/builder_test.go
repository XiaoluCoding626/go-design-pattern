package builder

import (
	"strings"
	"testing"
)

// 测试直接使用建造者创建汽车
func TestCarBuilderDirectUsage(t *testing.T) {
	builder := NewCarBuilder()

	car, err := builder.
		SetType(SportType).
		SetWheel(18, "布里奇斯通").
		SetEngine("3.0L V6", 300).
		SetSpeed(280).
		SetBrand("测试品牌").
		SetColor("蓝色").
		SetSeats(4).
		SetFuelType("汽油").
		AddFeature("天窗", true).
		AddFeature("导航", "高级版").
		Build()

	if err != nil {
		t.Fatalf("构建汽车失败: %v", err)
	}

	// 验证基本属性
	if car.Type() != SportType {
		t.Errorf("汽车类型错误: 得到 %v, 期望 %v", car.Type(), SportType)
	}
	if car.Speed() != 280 {
		t.Errorf("汽车速度错误: 得到 %v, 期望 %v", car.Speed(), 280)
	}
	if car.Brand() != "测试品牌" {
		t.Errorf("汽车品牌错误: 得到 %v, 期望 %v", car.Brand(), "测试品牌")
	}

	// 检查所有属性
	attrs := car.GetAttributes()
	if attrs["wheelSize"] != 18 {
		t.Errorf("车轮大小错误: 得到 %v, 期望 %v", attrs["wheelSize"], 18)
	}
	if attrs["wheelBrand"] != "布里奇斯通" {
		t.Errorf("车轮品牌错误: 得到 %v, 期望 %v", attrs["wheelBrand"], "布里奇斯通")
	}
	if attrs["engine"] != "3.0L V6" {
		t.Errorf("引擎错误: 得到 %v, 期望 %v", attrs["engine"], "3.0L V6")
	}
	if attrs["power"] != 300 {
		t.Errorf("功率错误: 得到 %v, 期望 %v", attrs["power"], 300)
	}
	if attrs["color"] != "蓝色" {
		t.Errorf("颜色错误: 得到 %v, 期望 %v", attrs["color"], "蓝色")
	}
	if attrs["seats"] != 4 {
		t.Errorf("座位数错误: 得到 %v, 期望 %v", attrs["seats"], 4)
	}
	if attrs["fuelType"] != "汽油" {
		t.Errorf("燃料类型错误: 得到 %v, 期望 %v", attrs["fuelType"], "汽油")
	}

	// 检查功能特性
	features := attrs["features"].(map[string]interface{})
	if features["天窗"] != true {
		t.Errorf("特性'天窗'错误: 得到 %v, 期望 %v", features["天窗"], true)
	}
	if features["导航"] != "高级版" {
		t.Errorf("特性'导航'错误: 得到 %v, 期望 %v", features["导航"], "高级版")
	}
}

// 测试缺少必要组件时的错误处理
func TestCarBuilderMissingComponents(t *testing.T) {
	builder := NewCarBuilder()

	// 测试缺少类型
	_, err := builder.
		SetWheel(18, "布里奇斯通").
		SetEngine("3.0L V6", 300).
		SetSpeed(280).
		SetBrand("测试品牌").
		Build()
	if err == nil || !strings.Contains(err.Error(), "汽车类型") {
		t.Error("期望因缺少汽车类型而失败，但未失败或错误消息不正确")
	}

	// 测试缺少车轮
	_, err = builder.Reset().
		SetType(SportType).
		SetEngine("3.0L V6", 300).
		SetSpeed(280).
		SetBrand("测试品牌").
		Build()
	if err == nil || !strings.Contains(err.Error(), "车轮尺寸") {
		t.Error("期望因缺少车轮尺寸而失败，但未失败或错误消息不正确")
	}

	// 测试缺少引擎
	_, err = builder.Reset().
		SetType(SportType).
		SetWheel(18, "布里奇斯通").
		SetSpeed(280).
		SetBrand("测试品牌").
		Build()
	if err == nil || !strings.Contains(err.Error(), "引擎") {
		t.Error("期望因缺少引擎而失败，但未失败或错误消息不正确")
	}

	// 测试缺少速度
	_, err = builder.Reset().
		SetType(SportType).
		SetWheel(18, "布里奇斯通").
		SetEngine("3.0L V6", 300).
		SetBrand("测试品牌").
		Build()
	if err == nil || !strings.Contains(err.Error(), "速度") {
		t.Error("期望因缺少最大速度而失败，但未失败或错误消息不正确")
	}

	// 测试缺少品牌
	_, err = builder.Reset().
		SetType(SportType).
		SetWheel(18, "布里奇斯通").
		SetEngine("3.0L V6", 300).
		SetSpeed(280).
		Build()
	if err == nil || !strings.Contains(err.Error(), "品牌") {
		t.Error("期望因缺少品牌而失败，但未失败或错误消息不正确")
	}
}

// 测试重置功能
func TestCarBuilderReset(t *testing.T) {
	builder := NewCarBuilder()

	// 首先部分构建一辆车
	builder.
		SetType(SportType).
		SetWheel(18, "布里奇斯通").
		SetEngine("3.0L V6", 300)

	// 然后重置
	builder.Reset()

	// 检查是否需要重新设置所有必要属性
	_, err := builder.
		SetType(SedanType).
		Build()
	if err == nil || !strings.Contains(err.Error(), "车轮尺寸") {
		t.Error("重置后应该需要重新设置所有必要属性")
	}

	// 完整构建一辆新车
	car, err := builder.
		SetType(SedanType).
		SetWheel(16, "米其林").
		SetEngine("1.8L 直列四缸", 140).
		SetSpeed(200).
		SetBrand("测试品牌").
		Build()

	if err != nil {
		t.Fatalf("重置后构建汽车失败: %v", err)
	}

	// 验证重置后设置的属性
	if car.Type() != SedanType {
		t.Errorf("重置后汽车类型错误: 得到 %v, 期望 %v", car.Type(), SedanType)
	}
	attrs := car.GetAttributes()
	if attrs["wheelSize"] != 16 {
		t.Errorf("重置后车轮大小错误: 得到 %v, 期望 %v", attrs["wheelSize"], 16)
	}
}

// 测试默认值设置
func TestCarBuilderDefaultValues(t *testing.T) {
	builder := NewCarBuilder()

	// 只设置必要的属性
	car, err := builder.
		SetType(SedanType).
		SetWheel(16, "通用").
		SetEngine("1.5L", 120).
		SetSpeed(180).
		SetBrand("经济型品牌").
		Build()

	if err != nil {
		t.Fatalf("构建汽车失败: %v", err)
	}

	// 检查默认值
	attrs := car.GetAttributes()
	if attrs["color"] != "白色" {
		t.Errorf("默认颜色错误: 得到 %v, 期望 %v", attrs["color"], "白色")
	}
	if attrs["seats"] != 5 {
		t.Errorf("默认座位数错误: 得到 %v, 期望 %v", attrs["seats"], 5)
	}
	if attrs["fuelType"] != "汽油" {
		t.Errorf("默认燃料类型错误: 得到 %v, 期望 %v", attrs["fuelType"], "汽油")
	}
}

// 测试Director
func TestDirector(t *testing.T) {
	builder := NewCarBuilder()
	director := NewDirector(builder)

	// 测试构建轿车
	sedan, err := director.BuildSedan("丰田")
	if err != nil {
		t.Fatalf("通过Director构建轿车失败: %v", err)
	}
	if sedan.Type() != SedanType || sedan.Brand() != "丰田" {
		t.Errorf("轿车属性错误: 类型=%v, 品牌=%v", sedan.Type(), sedan.Brand())
	}

	// 测试构建SUV
	suv, err := director.BuildSUV("本田")
	if err != nil {
		t.Fatalf("通过Director构建SUV失败: %v", err)
	}
	if suv.Type() != SUVType || suv.Brand() != "本田" {
		t.Errorf("SUV属性错误: 类型=%v, 品牌=%v", suv.Type(), suv.Brand())
	}
	suvAttrs := suv.GetAttributes()
	if suvAttrs["seats"] != 7 {
		t.Errorf("SUV座位数错误: 得到 %v, 期望 %v", suvAttrs["seats"], 7)
	}

	// 测试构建跑车
	sportsCar, err := director.BuildSportsCar("法拉利")
	if err != nil {
		t.Fatalf("通过Director构建跑车失败: %v", err)
	}
	if sportsCar.Type() != SportType || sportsCar.Brand() != "法拉利" {
		t.Errorf("跑车属性错误: 类型=%v, 品牌=%v", sportsCar.Type(), sportsCar.Brand())
	}
	sportsAttrs := sportsCar.GetAttributes()
	if sportsAttrs["maxSpeed"] != 330 {
		t.Errorf("跑车速度错误: 得到 %v, 期望 %v", sportsAttrs["maxSpeed"], 330)
	}

	// 测试构建豪华车
	luxuryCar, err := director.BuildLuxuryCar("奔驰")
	if err != nil {
		t.Fatalf("通过Director构建豪华车失败: %v", err)
	}
	if luxuryCar.Type() != LuxuryType || luxuryCar.Brand() != "奔驰" {
		t.Errorf("豪华车属性错误: 类型=%v, 品牌=%v", luxuryCar.Type(), luxuryCar.Brand())
	}
	luxuryAttrs := luxuryCar.GetAttributes()
	features := luxuryAttrs["features"].(map[string]interface{})
	if features["按摩座椅"] != true {
		t.Errorf("豪华车特性错误: 按摩座椅=%v, 期望 true", features["按摩座椅"])
	}
}

// 测试更改建造者
func TestDirectorChangeBuilder(t *testing.T) {
	originalBuilder := NewCarBuilder()
	director := NewDirector(originalBuilder)

	// 创建一个自定义建造者（这里只是为了测试，使用同一个建造者类型）
	newBuilder := NewCarBuilder()
	director.ChangeBuilder(newBuilder)

	// 确保可以正常使用新建造者
	car, err := director.BuildSedan("测试品牌")
	if err != nil {
		t.Fatalf("更改建造者后构建失败: %v", err)
	}
	if car.Brand() != "测试品牌" {
		t.Errorf("更改建造者后构建的汽车品牌错误: 得到 %v, 期望 %v", car.Brand(), "测试品牌")
	}
}

// 测试链式调用返回正确的建造者实例
func TestCarBuilderChaining(t *testing.T) {
	builder := NewCarBuilder()

	// 确保每个方法都返回相同的建造者实例
	b1 := builder.SetType(SedanType)
	b2 := b1.SetWheel(16, "通用")
	b3 := b2.SetEngine("2.0L", 150)

	if b1 != builder || b2 != builder || b3 != builder {
		t.Error("链式方法调用应该返回同一个建造者实例")
	}
}

// 集成测试：模拟实际使用场景
func TestIntegrationScenario(t *testing.T) {
	// 创建建造者和指导者
	builder := NewCarBuilder()
	director := NewDirector(builder)

	// 场景1：使用指导者创建预定义的汽车
	luxuryCar, _ := director.BuildLuxuryCar("宝马")

	// 验证预定义汽车的特性
	luxuryAttrs := luxuryCar.GetAttributes()
	if luxuryAttrs["type"] != LuxuryType || luxuryAttrs["brand"] != "宝马" {
		t.Errorf("预定义豪华车创建错误: 类型=%v, 品牌=%v", luxuryAttrs["type"], luxuryAttrs["brand"])
	}

	// 场景2：自定义创建一辆特殊车型
	customCar, _ := builder.Reset().
		SetType(SUVType).
		SetWheel(22, "固特异").
		SetEngine("3.5L 混合动力", 320).
		SetSpeed(240).
		SetBrand("路虎").
		SetColor("军绿色").
		SetSeats(5).
		SetFuelType("混合动力").
		AddFeature("空气悬挂", true).
		AddFeature("涉水能力", "900mm").
		AddFeature("车顶行李架", "铝合金").
		Build()

	// 验证自定义汽车的特性
	customAttrs := customCar.GetAttributes()
	features := customAttrs["features"].(map[string]interface{})

	if customAttrs["wheelSize"] != 22 ||
		customAttrs["color"] != "军绿色" ||
		features["涉水能力"] != "900mm" {
		t.Error("自定义车辆创建失败，属性不符合预期")
	}

	// 场景3：尝试使用相同的建造者创建不同类型的车
	builder.Reset()
	economyCar, _ := builder.
		SetType(SedanType).
		SetWheel(15, "普通品牌").
		SetEngine("1.6L 自然吸气", 110).
		SetSpeed(180).
		SetBrand("大众").
		SetColor("银色").
		Build()

	// 验证经济型车的特性
	economyAttrs := economyCar.GetAttributes()
	if economyAttrs["power"] != 110 || economyAttrs["wheelSize"] != 15 {
		t.Error("经济型车辆创建失败，属性不符合预期")
	}
}
