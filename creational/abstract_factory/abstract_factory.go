package abstract_factory

import (
	"fmt"
	"sync"
)

// DoorType 表示门的类型
type DoorType string

const (
	WoodenType DoorType = "wooden"
	MetalType  DoorType = "metal"
	GlassType  DoorType = "glass"
)

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

// DoorFactory 是抽象工厂接口，定义了创建门、门把手和门锁的方法
type DoorFactory interface {
	CreateDoor() Door
	CreateDoorHandle() DoorHandle
	CreateDoorLock() DoorLock
}

// ----- 木门产品族 -----

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

func (h *WoodenDoorHandle) Press() {
	fmt.Println("按下木门把手")
}

func (h *WoodenDoorHandle) Pull() {
	fmt.Println("拉动木门把手")
}

func (h *WoodenDoorHandle) GetMaterial() string {
	return "实木材质"
}

// WoodenDoorLock 是木门锁实现
type WoodenDoorLock struct{}

func (l *WoodenDoorLock) Lock() {
	fmt.Println("锁上木门锁")
}

func (l *WoodenDoorLock) Unlock() {
	fmt.Println("解锁木门锁")
}

func (l *WoodenDoorLock) GetSecurityLevel() int {
	return 1 // 安全级别低
}

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

// ----- 金属门产品族 -----

// MetalDoor 是金属门实现
type MetalDoor struct{}

func (d *MetalDoor) Open() {
	fmt.Println("金属门打开，发出沉重的声音")
}

func (d *MetalDoor) Close() {
	fmt.Println("金属门关闭，发出响亮的碰撞声")
}

func (d *MetalDoor) GetMaterial() string {
	return "钢铁材质"
}

// MetalDoorHandle 是金属门把手实现
type MetalDoorHandle struct{}

func (h *MetalDoorHandle) Press() {
	fmt.Println("按下金属门把手")
}

func (h *MetalDoorHandle) Pull() {
	fmt.Println("拉动金属门把手")
}

func (h *MetalDoorHandle) GetMaterial() string {
	return "不锈钢材质"
}

// MetalDoorLock 是金属门锁实现
type MetalDoorLock struct{}

func (l *MetalDoorLock) Lock() {
	fmt.Println("锁上金属安全锁")
}

func (l *MetalDoorLock) Unlock() {
	fmt.Println("解锁金属安全锁")
}

func (l *MetalDoorLock) GetSecurityLevel() int {
	return 3 // 安全级别高
}

// MetalDoorFactory 是金属门工厂，实现了 DoorFactory 接口
type MetalDoorFactory struct{}

func (f *MetalDoorFactory) CreateDoor() Door {
	return &MetalDoor{}
}

func (f *MetalDoorFactory) CreateDoorHandle() DoorHandle {
	return &MetalDoorHandle{}
}

func (f *MetalDoorFactory) CreateDoorLock() DoorLock {
	return &MetalDoorLock{}
}

// ----- 玻璃门产品族 -----

// GlassDoor 是玻璃门实现
type GlassDoor struct{}

func (d *GlassDoor) Open() {
	fmt.Println("玻璃门滑动打开")
}

func (d *GlassDoor) Close() {
	fmt.Println("玻璃门平稳关闭")
}

func (d *GlassDoor) GetMaterial() string {
	return "钢化玻璃材质"
}

// GlassDoorHandle 是玻璃门把手实现
type GlassDoorHandle struct{}

func (h *GlassDoorHandle) Press() {
	fmt.Println("按下玻璃门把手")
}

func (h *GlassDoorHandle) Pull() {
	fmt.Println("拉动玻璃门把手")
}

func (h *GlassDoorHandle) GetMaterial() string {
	return "铝合金材质"
}

// GlassDoorLock 是玻璃门锁实现
type GlassDoorLock struct{}

func (l *GlassDoorLock) Lock() {
	fmt.Println("锁上玻璃门电子锁")
}

func (l *GlassDoorLock) Unlock() {
	fmt.Println("解锁玻璃门电子锁")
}

func (l *GlassDoorLock) GetSecurityLevel() int {
	return 2 // 安全级别中等
}

// GlassDoorFactory 是玻璃门工厂，实现了 DoorFactory 接口
type GlassDoorFactory struct{}

func (f *GlassDoorFactory) CreateDoor() Door {
	return &GlassDoor{}
}

func (f *GlassDoorFactory) CreateDoorHandle() DoorHandle {
	return &GlassDoorHandle{}
}

func (f *GlassDoorFactory) CreateDoorLock() DoorLock {
	return &GlassDoorLock{}
}

// ----- 工厂创建器 -----

var (
	woodenFactory *WoodenDoorFactory
	metalFactory  *MetalDoorFactory
	glassFactory  *GlassDoorFactory
	once          sync.Once
)

// GetDoorFactory 根据指定的门类型返回相应的工厂实例
// 使用单例模式确保每种工厂只创建一个实例
func GetDoorFactory(doorType DoorType) (DoorFactory, error) {
	once.Do(func() {
		woodenFactory = &WoodenDoorFactory{}
		metalFactory = &MetalDoorFactory{}
		glassFactory = &GlassDoorFactory{}
	})

	switch doorType {
	case WoodenType:
		return woodenFactory, nil
	case MetalType:
		return metalFactory, nil
	case GlassType:
		return glassFactory, nil
	default:
		return nil, fmt.Errorf("不支持的门类型: %s", doorType)
	}
}

// DoorCreator 用于创建完整的门组件（门、把手、锁）
type DoorCreator struct {
	factory DoorFactory
}

// NewDoorCreator 创建一个新的门组件创建器
func NewDoorCreator(doorType DoorType) (*DoorCreator, error) {
	factory, err := GetDoorFactory(doorType)
	if err != nil {
		return nil, err
	}
	return &DoorCreator{factory: factory}, nil
}

// CreateCompleteDoor 创建一个完整的门系统并返回各组件
func (c *DoorCreator) CreateCompleteDoor() (Door, DoorHandle, DoorLock) {
	door := c.factory.CreateDoor()
	handle := c.factory.CreateDoorHandle()
	lock := c.factory.CreateDoorLock()
	return door, handle, lock
}
