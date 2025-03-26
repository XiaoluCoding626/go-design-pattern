package iterator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcreteIterator(t *testing.T) {
	data := []string{"a", "b", "c"}
	iterator := NewConcreteIterator(data)

	// 测试HasNext
	assert.True(t, iterator.HasNext(), "新创建的迭代器应该有下一个元素")

	// 测试Next
	item, ok := iterator.Next()
	assert.True(t, ok, "Next()应返回成功状态")
	assert.Equal(t, "a", item, "Next()应返回第一个元素")

	item, ok = iterator.Next()
	assert.True(t, ok, "Next()应返回成功状态")
	assert.Equal(t, "b", item, "Next()应返回第二个元素")

	// 测试Current
	current, ok := iterator.Current()
	assert.True(t, ok, "Current()应返回成功状态")
	assert.Equal(t, "b", current, "Current()应返回当前元素")

	// 测试Reset
	iterator.Reset()
	assert.True(t, iterator.HasNext(), "重置后迭代器应该有下一个元素")

	item, ok = iterator.Next()
	assert.True(t, ok, "重置后Next()应返回成功状态")
	assert.Equal(t, "a", item, "重置后Next()应返回第一个元素")

	// 测试遍历结束后调用Next
	iterator = NewConcreteIterator([]string{})
	item, ok = iterator.Next()
	assert.False(t, ok, "空迭代器不应返回有效元素")
}

func TestConcreteAggregate(t *testing.T) {
	aggregate := NewConcreteAggregate[int]()

	// 测试初始状态
	assert.Equal(t, 0, aggregate.Count(), "新创建的聚合对象应该为空")

	// 测试Add方法
	aggregate.Add(10)
	aggregate.Add(20)
	aggregate.Add(30)

	assert.Equal(t, 3, aggregate.Count(), "添加3个元素后长度应为3")

	// 测试Get方法
	item, ok := aggregate.Get(1)
	assert.True(t, ok, "Get()应返回成功状态")
	assert.Equal(t, 20, item, "索引1处的值应为20")

	// 测试超出范围的Get
	_, ok = aggregate.Get(5)
	assert.False(t, ok, "获取超出范围的索引应该返回false")

	// 测试Remove方法
	ok = aggregate.Remove(1)
	assert.True(t, ok, "删除有效索引应该返回true")
	assert.Equal(t, 2, aggregate.Count(), "删除后长度应为2")

	item, ok = aggregate.Get(1)
	assert.True(t, ok, "Get()应返回成功状态")
	assert.Equal(t, 30, item, "删除索引1后，新的索引1处的值应为30")

	// 测试删除无效索引
	ok = aggregate.Remove(5)
	assert.False(t, ok, "删除无效索引应该返回false")

	// 测试CreateIterator
	iterator := aggregate.CreateIterator()
	assert.True(t, iterator.HasNext(), "非空聚合对象创建的迭代器应该有下一个元素")
}

func TestGenericTypes(t *testing.T) {
	// 测试字符串类型
	strAggregate := NewConcreteAggregate[string]()
	strAggregate.Add("test1")
	strAggregate.Add("test2")

	strIterator := strAggregate.CreateIterator()
	str, ok := strIterator.Next()
	assert.True(t, ok, "字符串迭代器Next()应返回成功状态")
	assert.Equal(t, "test1", str, "字符串迭代器应返回正确的值")

	// 测试结构体类型
	type Person struct {
		Name string
		Age  int
	}

	personAggregate := NewConcreteAggregate[Person]()
	personAggregate.Add(Person{Name: "张三", Age: 30})
	personAggregate.Add(Person{Name: "李四", Age: 25})

	personIterator := personAggregate.CreateIterator()
	person, ok := personIterator.Next()
	assert.True(t, ok, "结构体迭代器Next()应返回成功状态")
	assert.Equal(t, "张三", person.Name, "结构体迭代器应返回正确的Name字段")
	assert.Equal(t, 30, person.Age, "结构体迭代器应返回正确的Age字段")
}

// 使用迭代器模式的示例
func Example() {
	// 创建字符串集合
	strCollection := NewConcreteAggregate[string]()
	strCollection.Add("设计模式")
	strCollection.Add("迭代器模式")
	strCollection.Add("Go语言实现")
	strCollection.Add("优化版本")

	fmt.Println("字符串集合示例:")
	// 获取迭代器并遍历
	iterator := strCollection.CreateIterator()
	for iterator.HasNext() {
		item, _ := iterator.Next()
		fmt.Printf("- %s\n", item)
	}

	// 重置并再次遍历前两个元素
	iterator.Reset()
	count := 0
	fmt.Println("\n重置后再遍历前两个元素:")
	for iterator.HasNext() && count < 2 {
		item, _ := iterator.Next()
		fmt.Printf("- %s\n", item)
		count++
	}

	// 获取当前元素（第二个元素）
	if current, ok := iterator.Current(); ok {
		fmt.Printf("\n当前元素: %s\n", current)
	}

	// 数字集合示例
	fmt.Println("\n数字集合示例:")
	numCollection := NewConcreteAggregate[int]()
	for i := 1; i <= 5; i++ {
		numCollection.Add(i * 10)
	}

	// 移除第3个元素
	numCollection.Remove(2)

	// 遍历数字集合
	numIterator := numCollection.CreateIterator()
	for numIterator.HasNext() {
		item, _ := numIterator.Next()
		fmt.Printf("- %d\n", item)
	}

	// Output:
	// 字符串集合示例:
	// - 设计模式
	// - 迭代器模式
	// - Go语言实现
	// - 优化版本
	//
	// 重置后再遍历前两个元素:
	// - 设计模式
	// - 迭代器模式
	//
	// 当前元素: 迭代器模式
	//
	// 数字集合示例:
	// - 10
	// - 20
	// - 40
	// - 50
}
