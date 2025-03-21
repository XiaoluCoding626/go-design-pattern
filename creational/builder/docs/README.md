# 建造者模式 (Builder Pattern)

## 简介

建造者模式是一种创建型设计模式，它允许你分步骤创建复杂对象。该模式允许你使用相同的创建代码生成不同类型和形式的对象。

在本例中，我们通过汽车制造过程来展示建造者模式。汽车是一个复杂对象，由多个组件组成（车轮、引擎、车身等），通过建造者模式，我们可以逐步配置这些组件，最终构建出完整的汽车。

## 结构

![建造者模式结构图](https://upload.wikimedia.org/wikipedia/commons/thumb/f/f3/Builder_UML_class_diagram.svg/700px-Builder_UML_class_diagram.svg.png)

### 核心组件

1. **产品（Product）**：复杂对象本身，在我们的例子中就是汽车（Car）
2. **建造者接口（Builder）**：声明创建产品各个部件的方法，在我们的例子中是ICarBuilder接口
3. **具体建造者（Concrete Builder）**：实现建造者接口的具体类，在我们的例子中是CarBuilder
4. **指导者（Director）**：定义使用建造者创建产品的顺序，在我们的例子中是Director类

## 代码实现

### 产品

我们的产品是汽车（Car），它实现了ICar接口：

```go
// ICar 汽车接口，定义汽车应该具备的能力
type ICar interface {
    Speed() int         // 获取最大速度
    Brand() string      // 获取品牌
    Type() CarType      // 获取汽车类型
    Brief()             // 打印汽车简介
    GetAttributes() map[string]interface{} // 获取所有属性
}
```

### 建造者接口

建造者接口定义了创建Car对象各个部分的方法：

```go
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
```

### 具体建造者

CarBuilder是建造者接口的具体实现：

```go
// CarBuilder 汽车建造者具体实现
type CarBuilder struct {
    car *Car // 正在构建的汽车
}

func NewCarBuilder() ICarBuilder {
    builder := &CarBuilder{}
    builder.Reset()
    return builder
}

// 实现各个部件的设置方法，例如：
func (b *CarBuilder) SetType(carType CarType) ICarBuilder {
    b.car.carType = carType
    return b
}

// ... 其他方法的实现 ...

func (b *CarBuilder) Build() (ICar, error) {
    // 验证必要的组件是否已设置
    if b.car.carType == "" {
        return nil, errors.New("必须设置汽车类型")
    }
    // ... 其他验证 ...
    
    // 创建一个新的汽车实例
    car := &Car{
        // 复制各个属性
    }
    
    // 设置默认值
    if car.color == "" {
        car.color = "白色"
    }
    // ... 其他默认值 ...
    
    return car, nil
}
```

### 指导者

指导者负责使用建造者按照特定顺序创建汽车：

```go
// Director 指导者，负责使用建造者创建特定类型的汽车
type Director struct {
    builder ICarBuilder
}

// 例如创建轿车的方法
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

// ... 其他车型的创建方法 ...
```

## 使用示例

下面是一个使用建造者模式创建汽车的示例：

```go
func main() {
    // 1. 使用Director创建预定义的汽车类型
    builder := NewCarBuilder()
    director := NewDirector(builder)
    
    sedan, err := director.BuildSedan("丰田")
    if err != nil {
        fmt.Println("创建轿车失败:", err)
        return
    }
    
    fmt.Println("已创建一辆轿车:")
    sedan.Brief()
    
    // 2. 直接使用Builder自定义汽车
    superCar, err := builder.Reset().
        SetType(SportType).
        SetWheel(21, "倍耐力").
        SetEngine("6.0L V12", 700).
        SetSpeed(350).
        SetBrand("法拉利").
        SetColor("红色").
        SetSeats(2).
        SetFuelType("高级汽油").
        AddFeature("赛道模式", true).
        AddFeature("陶瓷刹车", true).
        Build()
        
    if err != nil {
        fmt.Println("创建超跑失败:", err)
        return
    }
    
    fmt.Println("\n已创建一辆超级跑车:")
    superCar.Brief()
}
```

输出示例:
```
已创建一辆轿车:
这是一辆丰田的轿车
车轮: 17英寸 米其林品牌
引擎: 2.0L 涡轮增压 (180马力)
最大速度: 220公里/小时
颜色: 银色
座位数: 5
燃料类型: 汽油
额外特性:
  - 自动驾驶: 辅助
  - 导航系统: true

已创建一辆超级跑车:
这是一辆法拉利的跑车
车轮: 21英寸 倍耐力品牌
引擎: 6.0L V12 (700马力)
最大速度: 350公里/小时
颜色: 红色
座位数: 2
燃料类型: 高级汽油
额外特性:
  - 赛道模式: true
  - 陶瓷刹车: true
```

## 优点

1. **分步创建复杂对象**：可以逐步构建对象，轻松控制创建过程
2. **代码复用**：同一个建造者可以用于创建不同的产品表示
3. **单一职责原则**：将复杂构造代码从产品的业务逻辑中分离出来
4. **灵活性**：客户端代码无需了解产品内部结构即可创建复杂对象

## 缺点

1. **代码量增加**：需要创建多个新类，代码复杂度增加
2. **与产品结构紧密耦合**：每当产品发生变化，建造者也需要相应调整
3. **可能引入不必要的复杂性**：对于简单对象，使用建造者模式可能过于繁琐

## 适用场景

1. **创建复杂对象**：当对象的构建过程很复杂，包含多个步骤或部件时
2. **需要创建不同表示**：当需要使用相同的构建过程创建不同的表示时
3. **需要控制构建顺序**：当对象的创建过程需要特定的步骤顺序时
4. **参数过多**：当构造函数有太多参数，导致调用难以阅读和维护时

## 与其他模式的关系

1. **抽象工厂模式**：通常与建造者模式一起使用，指导者可以使用抽象工厂创建产品的部件
2. **工厂方法模式**：建造者模式注重逐步构建复杂对象，而工厂方法注重通过继承创建对象
3. **原型模式**：建造者聚焦于分步构建，而原型注重通过克隆创建复杂对象

## 总结

建造者模式是创建复杂对象的有力工具，特别适合那些有多个配置选项和构建步骤的对象。在我们的汽车示例中，我们能够以流畅的API创建各种不同类型的汽车，同时保持代码的可维护性和灵活性。
