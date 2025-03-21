package prototype

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"
)

// 定义颜色常量
type Color string

const (
	Red    Color = "红色"
	Green  Color = "绿色"
	Blue   Color = "蓝色"
	Yellow Color = "黄色"
	Black  Color = "黑色"
	White  Color = "白色"
)

// Shape 接口定义了克隆方法和其他公共方法
type Shape interface {
	Clone() Shape         // 浅克隆
	DeepClone() Shape     // 深克隆
	GetType() string      // 获取形状类型
	GetColor() Color      // 获取颜色
	SetColor(color Color) // 设置颜色
	GetArea() float64     // 计算面积
	String() string       // 字符串表示
}

// BaseShape 包含所有形状共有的属性
type BaseShape struct {
	Type  string
	Color Color
}

// 基础方法实现
func (b *BaseShape) GetType() string {
	return b.Type
}

func (b *BaseShape) GetColor() Color {
	return b.Color
}

func (b *BaseShape) SetColor(color Color) {
	b.Color = color
}

// Point 表示二维坐标点
type Point struct {
	X, Y float64
}

// Circle 结构体表示圆形
type Circle struct {
	BaseShape
	Radius float64
	Center *Point // 改为指针类型以支持真正的浅克隆
}

// NewCircle 创建新的圆形
func NewCircle(radius float64, x, y float64) *Circle {
	return &Circle{
		BaseShape: BaseShape{
			Type:  "圆形",
			Color: Blue,
		},
		Radius: radius,
		Center: &Point{X: x, Y: y}, // 创建指针
	}
}

// Clone 浅克隆实现
func (c *Circle) Clone() Shape {
	// 浅拷贝会共享Center指针
	return &Circle{
		BaseShape: BaseShape{
			Type:  c.Type,
			Color: c.Color,
		},
		Radius: c.Radius,
		Center: c.Center, // 共享同一个指针
	}
}

// DeepClone 深克隆实现
func (c *Circle) DeepClone() Shape {
	// 深拷贝创建新的Point实例
	return &Circle{
		BaseShape: BaseShape{
			Type:  c.Type,
			Color: c.Color,
		},
		Radius: c.Radius,
		Center: &Point{
			X: c.Center.X,
			Y: c.Center.Y,
		},
	}
}

// 另一种深克隆实现，使用序列化（适合更复杂的对象）
func (c *Circle) DeepCloneViaSerialization() (Shape, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(c)
	if err != nil {
		return nil, fmt.Errorf("序列化失败: %v", err)
	}

	var clone Circle
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&clone)
	if err != nil {
		return nil, fmt.Errorf("反序列化失败: %v", err)
	}

	return &clone, nil
}

// GetArea 计算圆的面积
func (c *Circle) GetArea() float64 {
	return 3.14159 * c.Radius * c.Radius
}

// String 返回圆的字符串表示
func (c *Circle) String() string {
	return fmt.Sprintf("%s[颜色=%s, 半径=%.2f, 中心=(%.1f,%.1f)]",
		c.Type, c.Color, c.Radius, c.Center.X, c.Center.Y)
}

// Rectangle 结构体表示矩形
type Rectangle struct {
	BaseShape
	Width    float64
	Height   float64
	Position *Point // 改为指针类型
}

// NewRectangle 创建新的矩形
func NewRectangle(width, height float64, x, y float64) *Rectangle {
	return &Rectangle{
		BaseShape: BaseShape{
			Type:  "矩形",
			Color: Red,
		},
		Width:    width,
		Height:   height,
		Position: &Point{X: x, Y: y}, // 创建指针
	}
}

// Clone 浅克隆实现
func (r *Rectangle) Clone() Shape {
	return &Rectangle{
		BaseShape: BaseShape{
			Type:  r.Type,
			Color: r.Color,
		},
		Width:    r.Width,
		Height:   r.Height,
		Position: r.Position, // 共享同一个指针
	}
}

// DeepClone 深克隆实现
func (r *Rectangle) DeepClone() Shape {
	return &Rectangle{
		BaseShape: BaseShape{
			Type:  r.Type,
			Color: r.Color,
		},
		Width:  r.Width,
		Height: r.Height,
		Position: &Point{
			X: r.Position.X,
			Y: r.Position.Y,
		},
	}
}

// GetArea 计算矩形的面积
func (r *Rectangle) GetArea() float64 {
	return r.Width * r.Height
}

// String 返回矩形的字符串表示
func (r *Rectangle) String() string {
	return fmt.Sprintf("%s[颜色=%s, 宽=%.2f, 高=%.2f, 位置=(%.1f,%.1f)]",
		r.Type, r.Color, r.Width, r.Height, r.Position.X, r.Position.Y)
}

// Triangle 结构体表示三角形
type Triangle struct {
	BaseShape
	A, B, C *Point // 改为指针类型
}

// NewTriangle 创建新的三角形
func NewTriangle(x1, y1, x2, y2, x3, y3 float64) *Triangle {
	return &Triangle{
		BaseShape: BaseShape{
			Type:  "三角形",
			Color: Green,
		},
		A: &Point{X: x1, Y: y1},
		B: &Point{X: x2, Y: y2},
		C: &Point{X: x3, Y: y3},
	}
}

// Clone 浅克隆实现
func (t *Triangle) Clone() Shape {
	return &Triangle{
		BaseShape: BaseShape{
			Type:  t.Type,
			Color: t.Color,
		},
		A: t.A,
		B: t.B,
		C: t.C,
	}
}

// DeepClone 深克隆实现
func (t *Triangle) DeepClone() Shape {
	return &Triangle{
		BaseShape: BaseShape{
			Type:  t.Type,
			Color: t.Color,
		},
		A: &Point{X: t.A.X, Y: t.A.Y},
		B: &Point{X: t.B.X, Y: t.B.Y},
		C: &Point{X: t.C.X, Y: t.C.Y},
	}
}

// GetArea 使用海伦公式计算三角形面积
func (t *Triangle) GetArea() float64 {
	// 计算三边长度
	a := distance(t.B, t.C)
	b := distance(t.A, t.C)
	c := distance(t.A, t.B)

	// 海伦公式
	s := (a + b + c) / 2
	return Sqrt(s * (s - a) * (s - b) * (s - c))
}

// 简单的平方根实现，避免导入math包
func Sqrt(x float64) float64 {
	// 使用牛顿迭代法计算平方根
	z := x / 2.0
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

// 计算两点间距离
func distance(p1, p2 *Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return Sqrt(dx*dx + dy*dy)
}

// String 返回三角形的字符串表示
func (t *Triangle) String() string {
	return fmt.Sprintf("%s[颜色=%s, 顶点A=(%.1f,%.1f), B=(%.1f,%.1f), C=(%.1f,%.1f)]",
		t.Type, t.Color, t.A.X, t.A.Y, t.B.X, t.B.Y, t.C.X, t.C.Y)
}

// ShapeCache 是原型管理器，用于存储和检索不同类型的原型
type ShapeCache struct {
	shapes map[string]Shape
	mu     sync.RWMutex // 用于线程安全
}

// NewShapeCache 创建新的形状缓存
func NewShapeCache() *ShapeCache {
	return &ShapeCache{
		shapes: make(map[string]Shape),
	}
}

// Add 添加形状到缓存
func (sc *ShapeCache) Add(id string, shape Shape) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	// 存储深克隆，避免外部修改影响原型
	sc.shapes[id] = shape.DeepClone()
}

// Get 获取形状的克隆
func (sc *ShapeCache) Get(id string) Shape {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	shape, ok := sc.shapes[id]
	if !ok {
		return nil
	}

	// 返回深克隆，避免修改原型
	return shape.DeepClone()
}

// LoadCache 预加载一些常用形状
func (sc *ShapeCache) LoadCache() {
	// 创建并存储基础形状
	circle := NewCircle(10, 5, 5)
	sc.Add("circle", circle)

	rectangle := NewRectangle(20, 10, 10, 20)
	sc.Add("rectangle", rectangle)

	triangle := NewTriangle(0, 0, 10, 0, 5, 10)
	sc.Add("triangle", triangle)

	// 添加一些特殊形状
	redCircle := NewCircle(15, 10, 10)
	redCircle.SetColor(Red)
	sc.Add("redCircle", redCircle)

	blueRectangle := NewRectangle(30, 5, 15, 25)
	blueRectangle.SetColor(Blue)
	sc.Add("blueRectangle", blueRectangle)
}

// GetShapeTypes 返回所有可用的形状类型
func (sc *ShapeCache) GetShapeTypes() []string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	types := make([]string, 0, len(sc.shapes))
	for key := range sc.shapes {
		types = append(types, key)
	}
	return types
}
