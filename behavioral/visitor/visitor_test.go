package visitor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// captureOutput 捕获标准输出的辅助函数
func captureOutput(f func()) string {
	// 保存原始的标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行函数
	f()

	// 恢复标准输出并获取输出内容
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestSceneryCreation 测试景点创建和基本属性
func TestSceneryCreation(t *testing.T) {
	assert := assert.New(t)

	// 测试豹子馆
	leopard := NewLeopardSpot()
	assert.Equal("豹子馆", leopard.GetName(), "豹子馆名称错误")
	assert.Equal(25, leopard.Price(), "豹子馆票价错误")
	assert.Contains(leopard.GetDescription(), "豹科动物", "豹子馆描述错误")

	// 测试海豚馆（无表演）
	dolphinNoShow := NewDolphinSpot(false)
	assert.Equal("海豚馆", dolphinNoShow.GetName(), "海豚馆(无表演)名称错误")
	assert.Equal(30, dolphinNoShow.Price(), "海豚馆(无表演)票价错误")
	assert.False(dolphinNoShow.HasShow(), "海豚馆应该无表演")

	// 测试海豚馆（有表演）
	dolphinWithShow := NewDolphinSpot(true)
	assert.Contains(dolphinWithShow.GetName(), "含表演", "海豚馆(有表演)名称错误")
	assert.Equal(45, dolphinWithShow.Price(), "海豚馆(有表演)票价错误")
	assert.True(dolphinWithShow.HasShow(), "海豚馆应该有表演")

	// 测试水族馆（无VIP区）
	aquariumNoVip := NewAquarium(false)
	assert.Equal("水族馆", aquariumNoVip.GetName(), "水族馆(无VIP区)名称错误")
	assert.Equal(35, aquariumNoVip.Price(), "水族馆(无VIP区)票价错误")
	assert.False(aquariumNoVip.HasVipArea(), "水族馆应该无VIP区")

	// 测试水族馆（有VIP区）
	aquariumWithVip := NewAquarium(true)
	assert.Contains(aquariumWithVip.GetName(), "VIP区", "水族馆(有VIP区)名称错误")
	assert.Equal(50, aquariumWithVip.Price(), "水族馆(有VIP区)票价错误")
	assert.True(aquariumWithVip.HasVipArea(), "水族馆应该有VIP区")
}

// TestVisitorCreation 测试访问者创建和基本属性
func TestVisitorCreation(t *testing.T) {
	assert := assert.New(t)

	// 测试学生访问者（有学生证）
	studentWithID := NewStudentVisitor(true)
	assert.Equal("学生", studentWithID.GetVisitorType(), "学生访问者类型错误")
	assert.True(studentWithID.hasStudentID, "学生访问者应该有学生证")

	// 测试学生访问者（无学生证）
	studentNoID := NewStudentVisitor(false)
	assert.False(studentNoID.hasStudentID, "学生访问者应该无学生证")

	// 测试普通访问者（周末）
	commonWeekend := NewCommonVisitor(true)
	assert.Equal("普通", commonWeekend.GetVisitorType(), "普通访问者类型错误")
	assert.True(commonWeekend.isWeekend, "普通访问者应该在周末")

	// 测试普通访问者（非周末）
	commonWeekday := NewCommonVisitor(false)
	assert.False(commonWeekday.isWeekend, "普通访问者应该在非周末")

	// 测试VIP访问者（不同等级）
	vip1 := NewVIPVisitor(1)
	vip3 := NewVIPVisitor(3)
	vipInvalid := NewVIPVisitor(5) // 超出范围，应被调整为3

	assert.Contains(vip1.GetVisitorType(), "VIP-1", "VIP访问者类型错误")
	assert.Contains(vip3.GetVisitorType(), "VIP-3", "VIP访问者类型错误")
	assert.Contains(vipInvalid.GetVisitorType(), "VIP-3", "VIP访问者类型错误")
	assert.Equal(3, vipInvalid.vipLevel, "无效VIP等级应该被调整为3")
}

// TestStudentVisitor 测试学生访问者的行为
func TestStudentVisitor(t *testing.T) {
	assert := assert.New(t)

	studentWithID := NewStudentVisitor(true)
	studentNoID := NewStudentVisitor(false)

	leopard := NewLeopardSpot()
	dolphin := NewDolphinSpot(true)
	aquarium := NewAquarium(true)

	// 测试持有学生证的学生访问豹子馆（半价）
	output := captureOutput(func() {
		studentWithID.VisitLeopardSpot(leopard)
	})
	assert.Contains(output, "票价: 12", "学生(有证)访问豹子馆价格错误") // 半价，25/2 = 12.5，取整为12

	// 测试无学生证的学生访问海豚馆（8折）
	output = captureOutput(func() {
		studentNoID.VisitDolphinSpot(dolphin)
	})
	assert.Contains(output, "票价: 36", "学生(无证)访问海豚馆价格错误") // 8折，45*0.8 = 36

	// 检查累计费用
	studentWithID.VisitAquarium(aquarium)
	expectedExpense := 12 + 25 // 豹子馆12元 + 水族馆带VIP区50/2=25元
	assert.Equal(expectedExpense, studentWithID.GetTotalExpense(), "学生访问总费用错误")
}

// TestCommonVisitor 测试普通访问者的行为
func TestCommonVisitor(t *testing.T) {
	assert := assert.New(t)

	weekendVisitor := NewCommonVisitor(true)
	weekdayVisitor := NewCommonVisitor(false)

	leopard := NewLeopardSpot()
	dolphin := NewDolphinSpot(false)

	// 测试周末访问（价格上浮20%）
	output := captureOutput(func() {
		weekendVisitor.VisitLeopardSpot(leopard)
	})
	assert.Contains(output, "票价: 30", "普通访问者周末访问豹子馆价格错误") // 周末上浮，25*1.2 = 30

	// 测试工作日访问（原价）
	output = captureOutput(func() {
		weekdayVisitor.VisitDolphinSpot(dolphin)
	})
	assert.Contains(output, "票价: 30", "普通访问者工作日访问海豚馆价格错误") // 工作日原价 30
}

// TestVIPVisitor 测试VIP访问者的行为
func TestVIPVisitor(t *testing.T) {
	assert := assert.New(t)

	vip1 := NewVIPVisitor(1)
	vip2 := NewVIPVisitor(2)
	vip3 := NewVIPVisitor(3)

	aquarium := NewAquarium(true)

	// 测试不同VIP等级的折扣
	output1 := captureOutput(func() {
		vip1.VisitAquarium(aquarium)
	})
	output2 := captureOutput(func() {
		vip2.VisitAquarium(aquarium)
	})
	output3 := captureOutput(func() {
		vip3.VisitAquarium(aquarium)
	})

	assert.Contains(output1, "票价: 45", "VIP1访问水族馆价格错误") // VIP1级9折，50*0.9 = 45
	assert.Contains(output2, "票价: 40", "VIP2访问水族馆价格错误") // VIP2级8折，50*0.8 = 40
	assert.Contains(output3, "票价: 35", "VIP3访问水族馆价格错误") // VIP3级7折，50*0.7 = 35
}

// TestZooManagement 测试动物园景点管理
func TestZooManagement(t *testing.T) {
	assert := assert.New(t)

	zoo := NewZoo("野生动物园")

	// 测试添加景点
	output := captureOutput(func() {
		zoo.Add(NewLeopardSpot())
		zoo.Add(NewDolphinSpot(true))
	})

	assert.Contains(output, "新增景点: 豹子馆", "动物园添加景点输出错误")
	assert.Contains(output, "新增景点: 海豚馆", "动物园添加景点输出错误")
	assert.Equal(2, len(zoo.Sceneries), "动物园景点数量错误")
}

// TestZooAcceptVisitors 测试动物园接待访问者综合功能
func TestZooAcceptVisitors(t *testing.T) {
	assert := assert.New(t)

	zoo := NewZoo("野生动物园")
	zoo.Add(NewLeopardSpot())
	zoo.Add(NewDolphinSpot(true))
	zoo.Add(NewAquarium(true))

	student := NewStudentVisitor(true)
	vip := NewVIPVisitor(2)

	// 测试学生访问所有景点
	studentOutput := captureOutput(func() {
		zoo.Accept(student)
	})

	assert.Contains(studentOutput, "欢迎 学生 游客参观", "动物园欢迎语错误")

	expectedStudentExpense := 12 + 22 + 25 // 豹子馆半价12元 + 海豚馆半价22元 + 水族馆半价25元
	assert.Contains(studentOutput, fmt.Sprintf("总花费: %d 元", expectedStudentExpense),
		"学生总花费显示错误")

	// 测试VIP访问所有景点
	vipOutput := captureOutput(func() {
		zoo.Accept(vip)
	})

	expectedVIPExpense := 20 + 36 + 40 // 豹子馆8折20元 + 海豚馆8折36元 + 水族馆8折40元
	assert.Contains(vipOutput, fmt.Sprintf("总花费: %d 元", expectedVIPExpense),
		"VIP总花费显示错误")
}

// Example 提供访问者模式的使用示例
func Example() {
	// 创建动物园
	zoo := NewZoo("快乐动物园")

	// 添加景点
	zoo.Add(NewLeopardSpot())
	zoo.Add(NewDolphinSpot(true))
	zoo.Add(NewAquarium(false))

	// 创建不同类型的访问者
	student := NewStudentVisitor(true)
	vip := NewVIPVisitor(3)

	// 访问者参观动物园
	fmt.Println("\n==== 学生游客参观行程 ====")
	zoo.Accept(student)

	fmt.Println("\n==== VIP游客参观行程 ====")
	zoo.Accept(vip)

	// Output:
	// 动物园 快乐动物园 新增景点: 豹子馆
	// 动物园 快乐动物园 新增景点: 海豚馆(含表演)
	// 动物园 快乐动物园 新增景点: 水族馆
	//
	// ==== 学生游客参观行程 ====
	//
	// 快乐动物园 欢迎 学生 游客参观！
	// 学生游客参观豹子馆，详情: 观赏猎豹、美洲豹等各种豹科动物，票价: 12元 (原价: 25元)
	// 学生游客参观海豚馆(含表演)，详情: 观赏海豚并可能欣赏精彩表演，今日有精彩表演，票价: 22元 (原价: 45元)
	// 学生游客参观水族馆，详情: 欣赏各种海洋生物，票价: 17元 (原价: 35元)
	// 学生 游客参观完成，总花费: 51 元
	//
	// ==== VIP游客参观行程 ====
	//
	// 快乐动物园 欢迎 VIP-3 游客参观！
	// VIP-3游客参观豹子馆，详情: 观赏猎豹、美洲豹等各种豹科动物，享受专属讲解，票价: 17元 (原价: 25元)
	// VIP-3游客参观海豚馆(含表演)，详情: 观赏海豚并可能欣赏精彩表演，安排前排观看表演，票价: 31元 (原价: 45元)
	// VIP-3游客参观水族馆，详情: 欣赏各种海洋生物，票价: 24元 (原价: 35元)
	// VIP-3 游客参观完成，总花费: 72 元
}

// TestTicketRounding 测试票价计算中的舍入行为
func TestTicketRounding(t *testing.T) {
	assert := assert.New(t)

	student := NewStudentVisitor(true)

	// 创建一个特殊价格景点（测试奇数价格的舍入）
	oddPricedScenery := &LeopardSpot{
		description: "测试景点",
		basePrice:   23, // 奇数价格，除以2后为11.5
	}

	output := captureOutput(func() {
		student.VisitLeopardSpot(oddPricedScenery)
	})

	// 检查票价是否正确舍入，23/2=11.5，舍入为11
	assert.Contains(output, "票价: 11", "学生票价舍入错误")
}

// TestBenchmarkVisitorPerformance 测试访问者模式性能
func BenchmarkVisitorPerformance(b *testing.B) {
	zoo := NewZoo("基准测试动物园")
	zoo.Add(NewLeopardSpot())
	zoo.Add(NewDolphinSpot(true))
	zoo.Add(NewAquarium(true))

	student := NewStudentVisitor(true)
	common := NewCommonVisitor(false)
	vip := NewVIPVisitor(3)

	visitors := []Visitor{student, common, vip}

	// 禁止输出以避免影响基准测试结果
	oldStdout := os.Stdout
	os.Stdout = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 选择一个访问者
		visitor := visitors[i%len(visitors)]
		zoo.Accept(visitor)
	}

	// 恢复标准输出
	os.Stdout = oldStdout
}
