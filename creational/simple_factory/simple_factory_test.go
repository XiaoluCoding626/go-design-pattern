package simple_factory

import (
	"strings"
	"testing"
)

// TestShapeTypes 测试形状类型的字符串表示
func TestShapeTypes(t *testing.T) {
	tests := []struct {
		shapeType ShapeType
		expected  string
	}{
		{ShapeTypeUnknown, "Unknown"},
		{ShapeTypeCircle, "Circle"},
		{ShapeTypeRectangle, "Rectangle"},
		{ShapeTypeTriangle, "Triangle"},
	}

	for _, test := range tests {
		if test.shapeType.String() != test.expected {
			t.Errorf("ShapeType %d 的字符串表示应为 %s，实际为 %s",
				test.shapeType, test.expected, test.shapeType.String())
		}
	}
}

// TestCircleImplementation 测试圆形的实现
func TestCircleImplementation(t *testing.T) {
	circle := NewCircle(5.0)

	// 测试类型
	if circle.GetType() != ShapeTypeCircle {
		t.Errorf("Circle 类型应为 %v，实际为 %v", ShapeTypeCircle, circle.GetType())
	}

	// 测试绘制方法
	expected := "Drawing Circle with radius 5.00"
	if circle.Draw() != expected {
		t.Errorf("Circle.Draw() 应返回 '%s'，实际返回 '%s'", expected, circle.Draw())
	}
}

// TestRectangleImplementation 测试矩形的实现
func TestRectangleImplementation(t *testing.T) {
	rectangle := NewRectangle(10.0, 20.0)

	// 测试类型
	if rectangle.GetType() != ShapeTypeRectangle {
		t.Errorf("Rectangle 类型应为 %v，实际为 %v", ShapeTypeRectangle, rectangle.GetType())
	}

	// 测试绘制方法
	expected := "Drawing Rectangle with width 10.00 and height 20.00"
	if rectangle.Draw() != expected {
		t.Errorf("Rectangle.Draw() 应返回 '%s'，实际返回 '%s'", expected, rectangle.Draw())
	}
}

// TestTriangleImplementation 测试三角形的实现
func TestTriangleImplementation(t *testing.T) {
	triangle := NewTriangle(3.0, 4.0, 5.0)

	// 测试类型
	if triangle.GetType() != ShapeTypeTriangle {
		t.Errorf("Triangle 类型应为 %v，实际为 %v", ShapeTypeTriangle, triangle.GetType())
	}

	// 测试绘制方法
	expected := "Drawing Triangle with sides 3.00, 4.00, 5.00"
	if triangle.Draw() != expected {
		t.Errorf("Triangle.Draw() 应返回 '%s'，实际返回 '%s'", expected, triangle.Draw())
	}
}

// TestCreateShapeByType 测试通过类型创建形状
func TestCreateShapeByType(t *testing.T) {
	factory := NewShapeFactory()

	// 测试创建圆形
	circle, err := factory.CreateShape(ShapeTypeCircle, 5.0)
	if err != nil {
		t.Errorf("创建圆形时发生错误: %v", err)
	}
	if circle.GetType() != ShapeTypeCircle {
		t.Errorf("创建的形状应该是圆形，实际是 %v", circle.GetType())
	}
	if !strings.Contains(circle.Draw(), "radius 5.00") {
		t.Errorf("圆形绘制方法应包含半径信息，实际输出: %s", circle.Draw())
	}

	// 测试创建矩形
	rectangle, err := factory.CreateShape(ShapeTypeRectangle, 10.0, 20.0)
	if err != nil {
		t.Errorf("创建矩形时发生错误: %v", err)
	}
	if rectangle.GetType() != ShapeTypeRectangle {
		t.Errorf("创建的形状应该是矩形，实际是 %v", rectangle.GetType())
	}
	if !strings.Contains(rectangle.Draw(), "width 10.00 and height 20.00") {
		t.Errorf("矩形绘制方法应包含宽高信息，实际输出: %s", rectangle.Draw())
	}

	// 测试创建三角形
	triangle, err := factory.CreateShape(ShapeTypeTriangle, 3.0, 4.0, 5.0)
	if err != nil {
		t.Errorf("创建三角形时发生错误: %v", err)
	}
	if triangle.GetType() != ShapeTypeTriangle {
		t.Errorf("创建的形状应该是三角形，实际是 %v", triangle.GetType())
	}
	if !strings.Contains(triangle.Draw(), "sides 3.00, 4.00, 5.00") {
		t.Errorf("三角形绘制方法应包含边长信息，实际输出: %s", triangle.Draw())
	}
}

// TestCreateShapeByName 测试通过名称创建形状
func TestCreateShapeByName(t *testing.T) {
	factory := NewShapeFactory()

	// 测试通过不同形式的名称创建形状
	testCases := []struct {
		name     string
		params   []float64
		expected ShapeType
	}{
		{"circle", []float64{2.5}, ShapeTypeCircle},
		{"Circle", []float64{3.5}, ShapeTypeCircle},
		{"rectangle", []float64{4.0, 5.0}, ShapeTypeRectangle},
		{"Rectangle", []float64{6.0, 7.0}, ShapeTypeRectangle},
		{"triangle", []float64{3.0, 4.0, 5.0}, ShapeTypeTriangle},
		{"Triangle", []float64{5.0, 12.0, 13.0}, ShapeTypeTriangle},
	}

	for _, tc := range testCases {
		shape, err := factory.CreateShapeByName(tc.name, tc.params...)
		if err != nil {
			t.Errorf("通过名称 '%s' 创建形状时发生错误: %v", tc.name, err)
			continue
		}
		if shape.GetType() != tc.expected {
			t.Errorf("通过名称 '%s' 创建的形状类型应为 %v，实际为 %v",
				tc.name, tc.expected, shape.GetType())
		}
	}
}

// TestErrorCases 测试错误情况
func TestErrorCases(t *testing.T) {
	factory := NewShapeFactory()

	// 测试无效的形状类型
	_, err := factory.CreateShape(ShapeType(99))
	if err == nil {
		t.Error("使用无效的形状类型应返回错误，但未返回")
	}

	// 测试无效的形状名称
	_, err = factory.CreateShapeByName("invalid_shape")
	if err == nil {
		t.Error("使用无效的形状名称应返回错误，但未返回")
	}

	// 测试参数不足的情况
	testCases := []struct {
		shapeType ShapeType
		params    []float64
	}{
		{ShapeTypeCircle, []float64{}},
		{ShapeTypeRectangle, []float64{1.0}},
		{ShapeTypeTriangle, []float64{1.0, 2.0}},
	}

	for _, tc := range testCases {
		_, err := factory.CreateShape(tc.shapeType, tc.params...)
		if err == nil {
			t.Errorf("创建类型 %v 的形状时参数不足应返回错误，但未返回", tc.shapeType)
		}
	}
}

// TestBackwardCompatibility 测试向后兼容的 NewShape 函数
func TestBackwardCompatibility(t *testing.T) {
	// 测试圆形
	circle := NewShape("circle")
	if circle == nil {
		t.Error("NewShape('circle') 不应返回 nil")
	} else if circle.GetType() != ShapeTypeCircle {
		t.Errorf("NewShape('circle') 应返回圆形，但返回了 %v", circle.GetType())
	}

	// 测试矩形
	rectangle := NewShape("rectangle")
	if rectangle == nil {
		t.Error("NewShape('rectangle') 不应返回 nil")
	} else if rectangle.GetType() != ShapeTypeRectangle {
		t.Errorf("NewShape('rectangle') 应返回矩形，但返回了 %v", rectangle.GetType())
	}

	// 测试无效类型
	unknown := NewShape("invalid")
	if unknown != nil {
		t.Errorf("NewShape('invalid') 应返回 nil，但返回了 %v", unknown)
	}
}

// TestInterfaceCompliance 测试所有形状都正确实现了 Shape 接口
func TestInterfaceCompliance(t *testing.T) {
	// 创建一个 Shape 接口类型的变量，并尝试赋值各种形状
	var shape Shape

	// 测试圆形
	shape = NewCircle(1.0)
	if shape.GetType() != ShapeTypeCircle {
		t.Error("Circle 未正确实现 Shape 接口")
	}

	// 测试矩形
	shape = NewRectangle(1.0, 2.0)
	if shape.GetType() != ShapeTypeRectangle {
		t.Error("Rectangle 未正确实现 Shape 接口")
	}

	// 测试三角形
	shape = NewTriangle(1.0, 2.0, 3.0)
	if shape.GetType() != ShapeTypeTriangle {
		t.Error("Triangle 未正确实现 Shape 接口")
	}
}
