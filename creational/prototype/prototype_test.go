package prototype

import (
	"fmt"
	"testing"
)

// 测试形状创建函数
func TestNewShapes(t *testing.T) {
	// 测试圆形
	circle := NewCircle(5.0, 10, 20)
	if circle.Radius != 5.0 || circle.Center.X != 10 || circle.Center.Y != 20 {
		t.Errorf("NewCircle创建的圆形属性错误: %v", circle)
	}
	if circle.GetType() != "圆形" || circle.GetColor() != Blue {
		t.Errorf("NewCircle创建的圆形基本属性错误: 类型=%s, 颜色=%s", circle.GetType(), circle.GetColor())
	}

	// 测试矩形
	rectangle := NewRectangle(15.0, 10.0, 5, 5)
	if rectangle.Width != 15.0 || rectangle.Height != 10.0 || rectangle.Position.X != 5 || rectangle.Position.Y != 5 {
		t.Errorf("NewRectangle创建的矩形属性错误: %v", rectangle)
	}
	if rectangle.GetType() != "矩形" || rectangle.GetColor() != Red {
		t.Errorf("NewRectangle创建的矩形基本属性错误: 类型=%s, 颜色=%s", rectangle.GetType(), rectangle.GetColor())
	}

	// 测试三角形
	triangle := NewTriangle(0, 0, 10, 0, 5, 10)
	if triangle.A.X != 0 || triangle.A.Y != 0 || triangle.B.X != 10 || triangle.B.Y != 0 || triangle.C.X != 5 || triangle.C.Y != 10 {
		t.Errorf("NewTriangle创建的三角形属性错误: %v", triangle)
	}
	if triangle.GetType() != "三角形" || triangle.GetColor() != Green {
		t.Errorf("NewTriangle创建的三角形基本属性错误: 类型=%s, 颜色=%s", triangle.GetType(), triangle.GetColor())
	}
}

// 测试圆形的浅克隆和深克隆
func TestCircleClone(t *testing.T) {
	original := NewCircle(10.0, 5, 5)
	original.SetColor(Red)

	// 测试浅克隆
	shallowClone := original.Clone().(*Circle)

	// 验证基本属性是否正确复制
	if shallowClone.Radius != original.Radius ||
		shallowClone.GetType() != original.GetType() ||
		shallowClone.GetColor() != original.GetColor() {
		t.Error("圆形浅克隆复制基本属性失败")
	}

	// 修改原始对象属性
	original.Radius = 20.0
	original.SetColor(Blue)
	original.Center.X = 15

	// 验证浅克隆对象的基本属性不受影响
	if shallowClone.Radius != 10.0 || shallowClone.GetColor() != Red {
		t.Error("浅克隆的基本属性被原始对象的修改影响了")
	}

	// 验证引用类型属性在浅克隆中是共享的
	if shallowClone.Center.X != 15 {
		t.Error("浅克隆应该共享Center引用，但没有")
	}

	// 测试深克隆
	original = NewCircle(10.0, 5, 5)
	original.SetColor(Red)
	deepClone := original.DeepClone().(*Circle)

	// 验证基本属性是否正确复制
	if deepClone.Radius != original.Radius ||
		deepClone.GetType() != original.GetType() ||
		deepClone.GetColor() != original.GetColor() {
		t.Error("圆形深克隆复制基本属性失败")
	}

	// 修改原始对象的引用类型属性
	original.Center.X = 25

	// 验证深克隆对象的引用类型属性不受影响
	if deepClone.Center.X != 5 {
		t.Error("深克隆的引用类型属性被原始对象的修改影响了")
	}
}

// 测试矩形的浅克隆和深克隆
func TestRectangleClone(t *testing.T) {
	original := NewRectangle(20.0, 10.0, 5, 5)
	original.SetColor(Blue)

	// 测试浅克隆
	shallowClone := original.Clone().(*Rectangle)

	// 验证基本属性是否正确复制
	if shallowClone.Width != original.Width ||
		shallowClone.Height != original.Height ||
		shallowClone.GetColor() != original.GetColor() {
		t.Error("矩形浅克隆复制基本属性失败")
	}

	// 修改原始对象属性
	original.Width = 30.0
	original.SetColor(Yellow)
	original.Position.Y = 15

	// 验证浅克隆对象的基本属性不受影响
	if shallowClone.Width != 20.0 || shallowClone.GetColor() != Blue {
		t.Error("浅克隆的基本属性被原始对象的修改影响了")
	}

	// 验证引用类型属性在浅克隆中是共享的
	if shallowClone.Position.Y != 15 {
		t.Error("浅克隆应该共享Position引用，但没有")
	}

	// 测试深克隆
	original = NewRectangle(20.0, 10.0, 5, 5)
	deepClone := original.DeepClone().(*Rectangle)

	// 修改原始对象的引用类型属性
	original.Position.X = 25

	// 验证深克隆对象的引用类型属性不受影响
	if deepClone.Position.X != 5 {
		t.Error("深克隆的引用类型属性被原始对象的修改影响了")
	}
}

// 测试三角形的浅克隆和深克隆
func TestTriangleClone(t *testing.T) {
	original := NewTriangle(0, 0, 10, 0, 5, 10)

	// 测试浅克隆
	shallowClone := original.Clone().(*Triangle)

	// 验证基本属性是否正确复制
	if shallowClone.A.X != original.A.X ||
		shallowClone.B.X != original.B.X ||
		shallowClone.C.Y != original.C.Y {
		t.Error("三角形浅克隆复制基本属性失败")
	}

	// 修改原始对象的引用类型属性
	original.A.X = 5

	// 验证引用类型属性在浅克隆中是共享的
	if shallowClone.A.X != 5 {
		t.Error("浅克隆应该共享Point引用，但没有")
	}

	// 测试深克隆
	original = NewTriangle(0, 0, 10, 0, 5, 10)
	deepClone := original.DeepClone().(*Triangle)

	// 修改原始对象的引用类型属性
	original.B.Y = 5

	// 验证深克隆对象的引用类型属性不受影响
	if deepClone.B.Y != 0 {
		t.Error("深克隆的引用类型属性被原始对象的修改影响了")
	}
}

// 测试序列化深克隆
func TestDeepCloneViaSerialization(t *testing.T) {
	original := NewCircle(15.0, 10, 20)
	original.SetColor(Yellow)

	clone, err := original.DeepCloneViaSerialization()
	if err != nil {
		t.Fatalf("序列化克隆失败: %v", err)
	}

	circleClone := clone.(*Circle)

	// 验证属性是否正确复制
	if circleClone.Radius != original.Radius ||
		circleClone.GetColor() != original.GetColor() ||
		circleClone.Center.X != original.Center.X {
		t.Error("序列化深克隆复制属性失败")
	}

	// 修改原始对象
	original.Center.X = 50
	original.SetColor(Green)

	// 验证克隆对象不受影响
	if circleClone.Center.X != 10 || circleClone.GetColor() != Yellow {
		t.Error("序列化深克隆对象被原始对象的修改影响了")
	}
}

// 测试面积计算
func TestGetArea(t *testing.T) {
	// 测试圆形面积
	circle := NewCircle(10, 0, 0)
	expectedCircleArea := 3.14159 * 100 // π * r²
	if !floatEqual(circle.GetArea(), expectedCircleArea, 0.0001) {
		t.Errorf("圆形面积计算错误: 期望 %f, 得到 %f", expectedCircleArea, circle.GetArea())
	}

	// 测试矩形面积
	rectangle := NewRectangle(10, 20, 0, 0)
	expectedRectArea := 200.0 // width * height
	if !floatEqual(rectangle.GetArea(), expectedRectArea, 0.0001) {
		t.Errorf("矩形面积计算错误: 期望 %f, 得到 %f", expectedRectArea, rectangle.GetArea())
	}

	// 测试三角形面积
	// 创建一个3-4-5的直角三角形
	triangle := NewTriangle(0, 0, 3, 0, 0, 4)
	expectedTriArea := 6.0 // 1/2 * base * height
	if !floatEqual(triangle.GetArea(), expectedTriArea, 0.0001) {
		t.Errorf("三角形面积计算错误: 期望 %f, 得到 %f", expectedTriArea, triangle.GetArea())
	}
}

// 测试颜色管理
func TestColorManagement(t *testing.T) {
	circle := NewCircle(10, 0, 0)

	// 默认颜色应该是Blue
	if circle.GetColor() != Blue {
		t.Errorf("默认颜色错误: 期望 %s, 得到 %s", Blue, circle.GetColor())
	}

	// 测试设置颜色
	circle.SetColor(Red)
	if circle.GetColor() != Red {
		t.Errorf("设置颜色后错误: 期望 %s, 得到 %s", Red, circle.GetColor())
	}

	// 测试克隆后颜色
	clone := circle.Clone()
	if clone.GetColor() != Red {
		t.Errorf("克隆后颜色错误: 期望 %s, 得到 %s", Red, clone.GetColor())
	}

	// 修改原对象颜色不应影响克隆
	circle.SetColor(Green)
	if clone.GetColor() != Red {
		t.Errorf("原对象修改后克隆颜色错误: 期望 %s, 得到 %s", Red, clone.GetColor())
	}
}

// 测试ShapeCache
func TestShapeCache(t *testing.T) {
	cache := NewShapeCache()

	// 测试添加和获取
	circle := NewCircle(10, 0, 0)
	cache.Add("smallCircle", circle)

	retrieved := cache.Get("smallCircle")
	if retrieved == nil {
		t.Fatal("从缓存获取形状失败")
	}

	retrievedCircle, ok := retrieved.(*Circle)
	if !ok {
		t.Fatal("从缓存获取的对象类型错误")
	}

	if retrievedCircle.Radius != 10 {
		t.Errorf("从缓存获取的对象属性错误: 期望半径 %f, 得到 %f", 10.0, retrievedCircle.Radius)
	}

	// 测试获取一个不存在的形状
	nonExistent := cache.Get("nonExistent")
	if nonExistent != nil {
		t.Error("获取不存在的形状应该返回nil")
	}

	// 测试预加载和类型列表
	cache.LoadCache()
	types := cache.GetShapeTypes()
	if len(types) != 6 { // 之前添加了1个，预加载了5个，总共6个
		t.Errorf("预加载后形状数量错误: 期望 %d, 得到 %d", 6, len(types))
	}

	// 测试从预加载的缓存中获取
	preloaded := cache.Get("redCircle")
	if preloaded == nil {
		t.Fatal("无法获取预加载的形状")
	}

	preloadedCircle := preloaded.(*Circle)
	if preloadedCircle.GetColor() != Red {
		t.Errorf("预加载形状颜色错误: 期望 %s, 得到 %s", Red, preloadedCircle.GetColor())
	}

	// 测试修改原型后对缓存的影响
	circle.SetColor(Yellow)
	afterModification := cache.Get("smallCircle")
	afterModCircle := afterModification.(*Circle)
	if afterModCircle.GetColor() == Yellow {
		t.Error("修改原型后不应影响缓存中获取的副本")
	}
}

// 测试浅克隆和深克隆的区别
func TestShallowVsDeepClone(t *testing.T) {
	// 创建一个原始圆形
	original := NewCircle(5, 10, 20)

	// 浅克隆
	shallow := original.Clone().(*Circle)

	// 深克隆
	deep := original.DeepClone().(*Circle)

	// 修改原始对象的中心点
	original.Center.X = 30

	// 验证浅克隆的中心点也被修改了
	if shallow.Center.X != 30 {
		t.Errorf("浅克隆的中心点X应该跟随原始对象变化为30，但得到%f", shallow.Center.X)
	}

	// 验证深克隆的中心点没有被修改
	if deep.Center.X != 10 {
		t.Errorf("深克隆的中心点X不应变化，应保持为10，但得到%f", deep.Center.X)
	}
}

// 测试String方法输出
func TestString(t *testing.T) {
	circle := NewCircle(5, 10, 15)
	circleStr := circle.String()
	expected := "圆形[颜色=蓝色, 半径=5.00, 中心=(10.0,15.0)]"
	if circleStr != expected {
		t.Errorf("Circle.String()输出错误: \n期望: %s\n得到: %s", expected, circleStr)
	}

	rectangle := NewRectangle(10, 20, 5, 15)
	rectangleStr := rectangle.String()
	expected = "矩形[颜色=红色, 宽=10.00, 高=20.00, 位置=(5.0,15.0)]"
	if rectangleStr != expected {
		t.Errorf("Rectangle.String()输出错误: \n期望: %s\n得到: %s", expected, rectangleStr)
	}

	triangle := NewTriangle(0, 0, 10, 0, 5, 10)
	triangleStr := triangle.String()
	expected = "三角形[颜色=绿色, 顶点A=(0.0,0.0), B=(10.0,0.0), C=(5.0,10.0)]"
	if triangleStr != expected {
		t.Errorf("Triangle.String()输出错误: \n期望: %s\n得到: %s", expected, triangleStr)
	}
}

// 集成测试：原型模式真实使用场景
func TestPrototypePatternUsage(t *testing.T) {
	// 创建一个形状缓存
	cache := NewShapeCache()
	cache.LoadCache()

	// 从缓存中获取形状并克隆它们
	circle1 := cache.Get("circle").(*Circle)
	circle2 := cache.Get("circle").(*Circle)

	// 确认获取的是不同的对象
	if circle1 == circle2 {
		t.Error("原型模式应该返回不同的对象实例")
	}

	// 修改其中一个对象不应影响另一个
	circle1.SetColor(Yellow)
	circle1.Radius = 20

	if circle2.GetColor() == Yellow || circle2.Radius == 20 {
		t.Error("修改克隆对象应该不影响其他克隆")
	}

	// 使用原型创建一组形状
	fmt.Println("创建一组形状：")
	shapes := []Shape{
		cache.Get("redCircle"),
		cache.Get("rectangle"),
		cache.Get("triangle"),
	}

	// 计算所有形状的面积总和
	totalArea := 0.0
	for _, shape := range shapes {
		fmt.Println(shape.String())
		totalArea += shape.GetArea()
	}

	fmt.Printf("总面积: %.2f\n", totalArea)

	// 使用浅克隆和深克隆展示区别
	original := NewCircle(10, 5, 5)

	shallowClone := original.Clone().(*Circle)
	deepClone := original.DeepClone().(*Circle)

	fmt.Println("\n修改前:")
	fmt.Println("原始:", original)
	fmt.Println("浅克隆:", shallowClone)
	fmt.Println("深克隆:", deepClone)

	// 修改原始对象
	original.Radius = 15
	original.Center.X = 20

	fmt.Println("\n修改后:")
	fmt.Println("原始:", original)
	fmt.Println("浅克隆:", shallowClone)
	fmt.Println("深克隆:", deepClone)
}

// 辅助函数：比较浮点数是否相等
func floatEqual(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}
