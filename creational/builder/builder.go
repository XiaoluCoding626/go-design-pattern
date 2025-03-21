package builder

import (
	"errors"
	"fmt"
)

// CarType 定义汽车类型
type CarType string

const (
	SedanType  CarType = "轿车"
	SUVType    CarType = "SUV"
	SportType  CarType = "跑车"
	LuxuryType CarType = "豪华车"
)

// ICar 汽车接口，定义汽车应该具备的能力
type ICar interface {
	Speed() int                            // 获取最大速度
	Brand() string                         // 获取品牌
	Type() CarType                         // 获取汽车类型
	Brief()                                // 打印汽车简介
	GetAttributes() map[string]interface{} // 获取所有属性
}

// ICarBuilder 汽车建造者接口，定义建造一辆车所需的步骤
type ICarBuilder interface {
	SetType(carType CarType) ICarBuilder                          // 设置车型
	SetWheel(size int, brand string) ICarBuilder                  // 设置车轮
	SetEngine(engine string, power int) ICarBuilder               // 设置引擎
	SetSpeed(max int) ICarBuilder                                 // 设置最大速度
	SetBrand(brand string) ICarBuilder                            // 设置品牌
	SetColor(color string) ICarBuilder                            // 设置颜色
	SetSeats(seats int) ICarBuilder                               // 设置座位数
	SetFuelType(fuelType string) ICarBuilder                      // 设置燃料类型
	AddFeature(featureName string, value interface{}) ICarBuilder // 添加特性
	Reset() ICarBuilder                                           // 重置构建器
	Build() (ICar, error)                                         // 构建汽车
}

// Car 具体的汽车结构体
type Car struct {
	carType    CarType                // 汽车类型
	wheelSize  int                    // 车轮尺寸
	wheelBrand string                 // 车轮品牌
	engine     string                 // 引擎型号
	power      int                    // 引擎功率(马力)
	maxSpeed   int                    // 最大速度(公里/小时)
	brandName  string                 // 品牌名称
	color      string                 // 颜色
	seats      int                    // 座位数
	fuelType   string                 // 燃料类型
	features   map[string]interface{} // 额外特性
}

// Speed 返回汽车最大速度
func (c *Car) Speed() int {
	return c.maxSpeed
}

// Brand 返回汽车品牌
func (c *Car) Brand() string {
	return c.brandName
}

// Type 返回汽车类型
func (c *Car) Type() CarType {
	return c.carType
}

// Brief 打印汽车简介
func (c *Car) Brief() {
	fmt.Printf("这是一辆%s的%s\n", c.brandName, c.carType)
	fmt.Printf("车轮: %d英寸 %s品牌\n", c.wheelSize, c.wheelBrand)
	fmt.Printf("引擎: %s (%d马力)\n", c.engine, c.power)
	fmt.Printf("最大速度: %d公里/小时\n", c.maxSpeed)
	fmt.Printf("颜色: %s\n", c.color)
	fmt.Printf("座位数: %d\n", c.seats)
	fmt.Printf("燃料类型: %s\n", c.fuelType)

	if len(c.features) > 0 {
		fmt.Println("额外特性:")
		for name, value := range c.features {
			fmt.Printf("  - %s: %v\n", name, value)
		}
	}
}

// GetAttributes 返回汽车的所有属性
func (c *Car) GetAttributes() map[string]interface{} {
	return map[string]interface{}{
		"type":       c.carType,
		"wheelSize":  c.wheelSize,
		"wheelBrand": c.wheelBrand,
		"engine":     c.engine,
		"power":      c.power,
		"maxSpeed":   c.maxSpeed,
		"brand":      c.brandName,
		"color":      c.color,
		"seats":      c.seats,
		"fuelType":   c.fuelType,
		"features":   c.features,
	}
}

// CarBuilder 汽车建造者具体实现
type CarBuilder struct {
	car *Car // 正在构建的汽车
}

// NewCarBuilder 创建新的汽车建造者实例
func NewCarBuilder() ICarBuilder {
	builder := &CarBuilder{}
	builder.Reset()
	return builder
}

// SetType 设置汽车类型
func (b *CarBuilder) SetType(carType CarType) ICarBuilder {
	b.car.carType = carType
	return b
}

// SetWheel 设置车轮大小和品牌
func (b *CarBuilder) SetWheel(size int, brand string) ICarBuilder {
	b.car.wheelSize = size
	b.car.wheelBrand = brand
	return b
}

// SetEngine 设置引擎型号和功率
func (b *CarBuilder) SetEngine(engine string, power int) ICarBuilder {
	b.car.engine = engine
	b.car.power = power
	return b
}

// SetSpeed 设置最大速度
func (b *CarBuilder) SetSpeed(max int) ICarBuilder {
	b.car.maxSpeed = max
	return b
}

// SetBrand 设置品牌
func (b *CarBuilder) SetBrand(brand string) ICarBuilder {
	b.car.brandName = brand
	return b
}

// SetColor 设置颜色
func (b *CarBuilder) SetColor(color string) ICarBuilder {
	b.car.color = color
	return b
}

// SetSeats 设置座位数
func (b *CarBuilder) SetSeats(seats int) ICarBuilder {
	b.car.seats = seats
	return b
}

// SetFuelType 设置燃料类型
func (b *CarBuilder) SetFuelType(fuelType string) ICarBuilder {
	b.car.fuelType = fuelType
	return b
}

// AddFeature 添加特性
func (b *CarBuilder) AddFeature(featureName string, value interface{}) ICarBuilder {
	b.car.features[featureName] = value
	return b
}

// Reset 重置构建器
func (b *CarBuilder) Reset() ICarBuilder {
	b.car = &Car{
		features: make(map[string]interface{}),
	}
	return b
}

// Build 构建并返回汽车
func (b *CarBuilder) Build() (ICar, error) {
	// 验证必要的组件是否已设置
	if b.car.carType == "" {
		return nil, errors.New("必须设置汽车类型")
	}
	if b.car.wheelSize == 0 {
		return nil, errors.New("必须设置车轮尺寸")
	}
	if b.car.engine == "" {
		return nil, errors.New("必须设置引擎型号")
	}
	if b.car.maxSpeed == 0 {
		return nil, errors.New("必须设置最大速度")
	}
	if b.car.brandName == "" {
		return nil, errors.New("必须设置品牌")
	}

	// 创建一个新的汽车实例，避免修改正在构建的实例
	car := &Car{
		carType:    b.car.carType,
		wheelSize:  b.car.wheelSize,
		wheelBrand: b.car.wheelBrand,
		engine:     b.car.engine,
		power:      b.car.power,
		maxSpeed:   b.car.maxSpeed,
		brandName:  b.car.brandName,
		color:      b.car.color,
		seats:      b.car.seats,
		fuelType:   b.car.fuelType,
		features:   make(map[string]interface{}),
	}

	// 复制特性
	for k, v := range b.car.features {
		car.features[k] = v
	}

	// 设置默认值
	if car.color == "" {
		car.color = "白色"
	}
	if car.seats == 0 {
		car.seats = 5
	}
	if car.fuelType == "" {
		car.fuelType = "汽油"
	}

	return car, nil
}

// Director 指导者，负责使用建造者创建特定类型的汽车
type Director struct {
	builder ICarBuilder
}

// NewDirector 创建新的指导者
func NewDirector(builder ICarBuilder) *Director {
	return &Director{
		builder: builder,
	}
}

// ChangeBuilder 更改建造者
func (d *Director) ChangeBuilder(builder ICarBuilder) {
	d.builder = builder
}

// BuildSedan 构建轿车
func (d *Director) BuildSedan(brand string) (ICar, error) {
	return d.builder.Reset().
		SetType(SedanType).
		SetWheel(17, "米其林").
		SetEngine("2.0L 涡轮增压", 180).
		SetSpeed(220).
		SetBrand(brand).
		SetSeats(5).
		SetColor("银色").
		SetFuelType("汽油").
		AddFeature("自动驾驶", "辅助").
		AddFeature("导航系统", true).
		Build()
}

// BuildSUV 构建SUV
func (d *Director) BuildSUV(brand string) (ICar, error) {
	return d.builder.Reset().
		SetType(SUVType).
		SetWheel(19, "固特异").
		SetEngine("2.5L V6", 220).
		SetSpeed(200).
		SetBrand(brand).
		SetSeats(7).
		SetColor("黑色").
		SetFuelType("柴油").
		AddFeature("四驱系统", true).
		AddFeature("越野模式", "高级").
		Build()
}

// BuildSportsCar 构建跑车
func (d *Director) BuildSportsCar(brand string) (ICar, error) {
	return d.builder.Reset().
		SetType(SportType).
		SetWheel(21, "倍耐力").
		SetEngine("4.0L V8 双涡轮", 580).
		SetSpeed(330).
		SetBrand(brand).
		SetSeats(2).
		SetColor("红色").
		SetFuelType("高级汽油").
		AddFeature("碳纤维车身", true).
		AddFeature("弹射起步", true).
		AddFeature("活跃悬挂", "赛道模式").
		Build()
}

// BuildLuxuryCar 构建豪华车
func (d *Director) BuildLuxuryCar(brand string) (ICar, error) {
	return d.builder.Reset().
		SetType(LuxuryType).
		SetWheel(20, "马牌").
		SetEngine("3.0L 直列六缸 混合动力", 400).
		SetSpeed(250).
		SetBrand(brand).
		SetSeats(5).
		SetColor("深蓝色").
		SetFuelType("混合动力").
		AddFeature("真皮内饰", "Nappa皮革").
		AddFeature("按摩座椅", true).
		AddFeature("环绕音响", "Burmester").
		AddFeature("自动泊车", "全景影像").
		Build()
}
