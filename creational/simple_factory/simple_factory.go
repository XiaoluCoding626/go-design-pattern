// Package simple_factory 实现简单工厂设计模式
//
// 简单工厂模式不属于GoF的23种设计模式，但它是一种常用的对象创建模式。
// 它提供了一个工厂类/函数，根据参数的不同返回不同的实例，使得
// 客户端与对象的创建过程解耦。
package simple_factory

import (
	"fmt"
)

// ShapeType 定义了可创建的形状类型
type ShapeType int

// 支持的形状类型常量
const (
	ShapeTypeUnknown   ShapeType = iota // 未知形状类型
	ShapeTypeCircle                     // 圆形
	ShapeTypeRectangle                  // 矩形
	ShapeTypeTriangle                   // 三角形
)

// String 返回ShapeType的字符串表示
func (s ShapeType) String() string {
	return [...]string{"Unknown", "Circle", "Rectangle", "Triangle"}[s]
}

// Shape 定义了所有形状必须实现的接口
type Shape interface {
	Draw() string       // 绘制形状
	GetType() ShapeType // 获取形状类型
}

// BaseShape 提供所有形状的通用功能
type BaseShape struct {
	shapeType ShapeType
}

// GetType 返回形状的类型
func (b *BaseShape) GetType() ShapeType {
	return b.shapeType
}

// Circle 实现圆形
type Circle struct {
	BaseShape
	radius float64 // 添加圆的半径属性
}

// NewCircle 创建一个新的圆形
func NewCircle(radius float64) *Circle {
	return &Circle{
		BaseShape: BaseShape{shapeType: ShapeTypeCircle},
		radius:    radius,
	}
}

// Draw 实现Shape接口的Draw方法
func (c *Circle) Draw() string {
	return fmt.Sprintf("Drawing Circle with radius %.2f", c.radius)
}

// Rectangle 实现矩形
type Rectangle struct {
	BaseShape
	width, height float64 // 添加矩形的宽度和高度属性
}

// NewRectangle 创建一个新的矩形
func NewRectangle(width, height float64) *Rectangle {
	return &Rectangle{
		BaseShape: BaseShape{shapeType: ShapeTypeRectangle},
		width:     width,
		height:    height,
	}
}

// Draw 实现Shape接口的Draw方法
func (r *Rectangle) Draw() string {
	return fmt.Sprintf("Drawing Rectangle with width %.2f and height %.2f", r.width, r.height)
}

// Triangle 实现三角形
type Triangle struct {
	BaseShape
	a, b, c float64 // 添加三角形的三条边长度
}

// NewTriangle 创建一个新的三角形
func NewTriangle(a, b, c float64) *Triangle {
	return &Triangle{
		BaseShape: BaseShape{shapeType: ShapeTypeTriangle},
		a:         a,
		b:         b,
		c:         c,
	}
}

// Draw 实现Shape接口的Draw方法
func (t *Triangle) Draw() string {
	return fmt.Sprintf("Drawing Triangle with sides %.2f, %.2f, %.2f", t.a, t.b, t.c)
}

// ShapeFactory 定义了形状工厂结构体
type ShapeFactory struct{}

// NewShapeFactory 创建一个新的形状工厂
func NewShapeFactory() *ShapeFactory {
	return &ShapeFactory{}
}

// CreateShape 根据形状类型创建具体的形状实例
// 第一个参数是形状类型，后续参数是创建形状所需的参数
func (f *ShapeFactory) CreateShape(shapeType ShapeType, params ...float64) (Shape, error) {
	switch shapeType {
	case ShapeTypeCircle:
		if len(params) < 1 {
			return nil, fmt.Errorf("创建圆形需要指定半径")
		}
		return NewCircle(params[0]), nil
	case ShapeTypeRectangle:
		if len(params) < 2 {
			return nil, fmt.Errorf("创建矩形需要指定宽度和高度")
		}
		return NewRectangle(params[0], params[1]), nil
	case ShapeTypeTriangle:
		if len(params) < 3 {
			return nil, fmt.Errorf("创建三角形需要指定三条边长度")
		}
		return NewTriangle(params[0], params[1], params[2]), nil
	default:
		return nil, fmt.Errorf("不支持的形状类型: %v", shapeType)
	}
}

// CreateShapeByName 根据形状名称创建具体的形状实例
func (f *ShapeFactory) CreateShapeByName(shapeName string, params ...float64) (Shape, error) {
	switch shapeName {
	case "circle", "Circle":
		return f.CreateShape(ShapeTypeCircle, params...)
	case "rectangle", "Rectangle":
		return f.CreateShape(ShapeTypeRectangle, params...)
	case "triangle", "Triangle":
		return f.CreateShape(ShapeTypeTriangle, params...)
	default:
		return nil, fmt.Errorf("不支持的形状名称: %s", shapeName)
	}
}

// 为向后兼容保留的函数
// NewShape 是原始简单工厂函数，根据形状类型名创建形状
// 注意：这个函数仅为兼容旧代码保留，新代码应使用ShapeFactory
func NewShape(shapeName string) Shape {
	factory := NewShapeFactory()
	switch shapeName {
	case "circle":
		shape, _ := factory.CreateShape(ShapeTypeCircle, 1.0) // 默认半径1.0
		return shape
	case "rectangle":
		shape, _ := factory.CreateShape(ShapeTypeRectangle, 1.0, 1.0) // 默认1.0x1.0
		return shape
	default:
		return nil // 为保持兼容性，保留返回nil的行为
	}
}
