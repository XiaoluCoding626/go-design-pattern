# 抽象工厂模式 (Abstract Factory Pattern)

## 简介

抽象工厂模式是一种创建型设计模式，它提供一个接口来创建一系列相关或相互依赖的对象，而无需指定它们的具体类。抽象工厂模式通常也被称为"工厂的工厂"，它将一组对象的创建过程抽象化。

在本例中，我们实现了一个门系统，可以创建不同材质（木质、金属、玻璃）的门及其配件（把手、锁）。

## 结构

![抽象工厂模式结构图](https://upload.wikimedia.org/wikipedia/commons/thumb/9/9d/Abstract_factory_UML.svg/700px-Abstract_factory_UML.svg.png)

### 核心组件

1. **抽象工厂（Abstract Factory）**：声明了一组创建产品的方法，每个方法对应一种产品类型
2. **具体工厂（Concrete Factory）**：实现了抽象工厂的方法，创建特定风格的产品
3. **抽象产品（Abstract Product）**：为产品声明接口
4. **具体产品（Concrete Product）**：实现抽象产品接口，代表特定风格的产品

## 代码实现

### 抽象产品接口

我们定义了三种抽象产品接口：

```go
// Door 是门接口
type Door interface {
    Open()
    Close()
    GetMaterial() string
}

// DoorHandle 是门把手接口
type DoorHandle interface {
    Press()
    Pull()
    GetMaterial() string
}

// DoorLock 是门锁接口
type DoorLock interface {
    Lock()
    Unlock()
    GetSecurityLevel() int
}
```

### 抽象工厂接口

```go
// DoorFactory 是抽象工厂接口，定义了创建门、门把手和门锁的方法
type DoorFactory interface {
    CreateDoor() Door
    CreateDoorHandle() DoorHandle
    CreateDoorLock() DoorLock
}
```

### 具体产品族

我们实现了三个产品族：木门、金属门和玻璃门，每个产品族包含门、门把手和门锁三种产品。

例如，木门产品族：

```go
// WoodenDoor 是木门实现
type WoodenDoor struct{}

func (d *WoodenDoor) Open() {
    fmt.Println("木门打开，发出吱呀声")
}

func (d *WoodenDoor) Close() {
    fmt.Println("木门关闭，发出砰的一声")
}

func (d *WoodenDoor) GetMaterial() string {
    return "实木材质"
}

// WoodenDoorHandle 是木门把手实现
type WoodenDoorHandle struct{}

// ...

// WoodenDoorLock 是木门锁实现
type WoodenDoorLock struct{}

// ...
```

### 具体工厂

每个产品族都有一个对应的具体工厂：

```go
// WoodenDoorFactory 是木门工厂，实现了 DoorFactory 接口
type WoodenDoorFactory struct{}

func (f *WoodenDoorFactory) CreateDoor() Door {
    return &WoodenDoor{}
}

func (f *WoodenDoorFactory) CreateDoorHandle() DoorHandle {
    return &WoodenDoorHandle{}
}

func (f *WoodenDoorFactory) CreateDoorLock() DoorLock {
    return &WoodenDoorLock{}
}
```

### 工厂创建器

为了方便客户端使用，我们提供了一个工厂创建器：

```go
// GetDoorFactory 根据指定的门类型返回相应的工厂实例
func GetDoorFactory(doorType DoorType) (DoorFactory, error) {
    // ...
}

// DoorCreator 用于创建完整的门组件（门、把手、锁）
type DoorCreator struct {
    factory DoorFactory
}

func NewDoorCreator(doorType DoorType) (*DoorCreator, error) {
    // ...
}

func (c *DoorCreator) CreateCompleteDoor() (Door, DoorHandle, DoorLock) {
    // ...
}
```

## 使用示例

```go
func main() {
    // 创建一个木门系统
    creator, err := NewDoorCreator(WoodenType)
    if err != nil {
        fmt.Println("错误:", err)
        return
    }
    
    door, handle, lock := creator.CreateCompleteDoor()
    
    // 使用木门系统
    fmt.Println("木门材质:", door.GetMaterial())
    door.Open()
    handle.Press()
    lock.Lock()
    
    // 创建一个金属门系统
    metalCreator, _ := NewDoorCreator(MetalType)
    metalDoor, metalHandle, metalLock := metalCreator.CreateCompleteDoor()
    
    // 使用金属门系统
    fmt.Println("金属门材质:", metalDoor.GetMaterial())
    fmt.Println("金属门锁安全级别:", metalLock.GetSecurityLevel())
    metalDoor.Close()
    metalHandle.Pull()
}
```

## 优点

1. **确保产品兼容性**：抽象工厂保证了一个产品族的所有产品都能适当地一起工作
2. **隔离具体类**：客户端代码与具体产品类分离，只与抽象接口交互
3. **单一职责原则**：将产品创建代码集中在一个地方，使得代码更容易维护
4. **开闭原则**：引入新产品族不需要修改现有代码，只需添加新工厂和产品类

## 缺点

1. **增加复杂性**：相比于简单工厂，抽象工厂模式更为复杂
2. **难以扩展产品种类**：增加新产品需要修改抽象工厂接口及所有实现类
3. **代码重复**：不同产品族可能有很多相似的实现代码

## 应用场景

1. **系统需要独立于其产品的创建、组合和表示**
2. **系统要由多个产品系列中的一个来配置**
3. **相关产品对象必须一起使用，设计约束需要强化这种约束**
4. **想要提供一个产品库，只想显示它们的接口而不是实现**

## 与其他模式的关系

1. **工厂方法模式**：抽象工厂通常基于一组工厂方法实现
2. **单例模式**：抽象工厂经常与单例模式结合，确保每种工厂类型只有一个实例
3. **建造者模式**：当创建复杂对象时，抽象工厂可以与建造者模式结合使用
4. **原型模式**：当产品的实例化成本高时，抽象工厂可以使用原型模式实现

## 总结

抽象工厂模式适用于需要创建一系列相关对象的场景，尤其是当这些对象的创建逻辑需要与使用逻辑分离时。在我们的示例中，不同材质的门系统（木门、金属门、玻璃门）可以通过相应的工厂进行创建，而客户端代码不需要知道具体如何创建这些对象。