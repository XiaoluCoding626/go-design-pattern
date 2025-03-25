package visitor

import (
	"fmt"
	"time"
)

// Visitor 抽象访问者接口 - 定义对每种场景的访问方法
type Visitor interface {
	VisitLeopardSpot(leopard *LeopardSpot) // 参观猎豹馆
	VisitDolphinSpot(dolphin *DolphinSpot) // 参观海豚馆
	VisitAquarium(aquarium *Aquarium)      // 参观水族馆
	GetTotalExpense() int                  // 获取总花费
	GetVisitorType() string                // 获取访问者类型
}

// Scenery 场馆景点接口 - 定义场景对象的通用行为
type Scenery interface {
	Accept(visitor Visitor) // 接待访问者
	Price() int             // 基础票价
	GetName() string        // 获取景点名称
	GetDescription() string // 获取景点描述
}

// Zoo 动物园类 - 复合对象，包含多个景点
type Zoo struct {
	Name      string     // 动物园名称
	Sceneries []Scenery  // 动物园包含的景点
	OpenTime  *time.Time // 开放时间
}

// NewZoo 创建一个新的动物园
func NewZoo(name string) *Zoo {
	now := time.Now()
	return &Zoo{
		Name:      name,
		Sceneries: make([]Scenery, 0),
		OpenTime:  &now,
	}
}

// Add 给动物园添加景点
func (z *Zoo) Add(scenery Scenery) {
	z.Sceneries = append(z.Sceneries, scenery)
	fmt.Printf("动物园 %s 新增景点: %s\n", z.Name, scenery.GetName())
}

// Accept 动物园接待游客，游客将参观所有景点
func (z *Zoo) Accept(v Visitor) {
	fmt.Printf("\n%s 欢迎 %s 游客参观！\n", z.Name, v.GetVisitorType())
	for _, scenery := range z.Sceneries {
		scenery.Accept(v)
	}
	fmt.Printf("%s 游客参观完成，总花费: %d 元\n", v.GetVisitorType(), v.GetTotalExpense())
}

// LeopardSpot 豹子馆实现
type LeopardSpot struct {
	description string
	basePrice   int
}

// NewLeopardSpot 创建豹子馆
func NewLeopardSpot() *LeopardSpot {
	return &LeopardSpot{
		description: "观赏猎豹、美洲豹等各种豹科动物",
		basePrice:   25,
	}
}

// Accept 实现Scenery接口的Accept方法
func (l *LeopardSpot) Accept(visitor Visitor) {
	visitor.VisitLeopardSpot(l)
}

// Price 基础票价
func (l *LeopardSpot) Price() int {
	return l.basePrice
}

// GetName 获取景点名称
func (l *LeopardSpot) GetName() string {
	return "豹子馆"
}

// GetDescription 获取景点描述
func (l *LeopardSpot) GetDescription() string {
	return l.description
}

// DolphinSpot 海豚馆实现
type DolphinSpot struct {
	description string
	basePrice   int
	hasShow     bool // 是否有表演
}

// NewDolphinSpot 创建海豚馆
func NewDolphinSpot(hasShow bool) *DolphinSpot {
	price := 30
	if hasShow {
		price = 45
	}
	return &DolphinSpot{
		description: "观赏海豚并可能欣赏精彩表演",
		basePrice:   price,
		hasShow:     hasShow,
	}
}

// Accept 实现Scenery接口的Accept方法
func (d *DolphinSpot) Accept(visitor Visitor) {
	visitor.VisitDolphinSpot(d)
}

// Price 海豚馆票价
func (d *DolphinSpot) Price() int {
	return d.basePrice
}

// GetName 获取景点名称
func (d *DolphinSpot) GetName() string {
	name := "海豚馆"
	if d.hasShow {
		name += "(含表演)"
	}
	return name
}

// GetDescription 获取景点描述
func (d *DolphinSpot) GetDescription() string {
	return d.description
}

// HasShow 检查是否有表演
func (d *DolphinSpot) HasShow() bool {
	return d.hasShow
}

// Aquarium 水族馆实现
type Aquarium struct {
	description string
	basePrice   int
	vipArea     bool // 是否包含VIP区域
}

// NewAquarium 创建水族馆
func NewAquarium(vipArea bool) *Aquarium {
	price := 35
	if vipArea {
		price = 50
	}
	return &Aquarium{
		description: "欣赏各种海洋生物",
		basePrice:   price,
		vipArea:     vipArea,
	}
}

// Accept 实现Scenery接口的Accept方法
func (a *Aquarium) Accept(visitor Visitor) {
	visitor.VisitAquarium(a)
}

// Price 水族馆票价
func (a *Aquarium) Price() int {
	return a.basePrice
}

// GetName 获取景点名称
func (a *Aquarium) GetName() string {
	name := "水族馆"
	if a.vipArea {
		name += "(含VIP区)"
	}
	return name
}

// GetDescription 获取景点描述
func (a *Aquarium) GetDescription() string {
	return a.description
}

// HasVipArea 检查是否有VIP区
func (a *Aquarium) HasVipArea() bool {
	return a.vipArea
}

// BaseVisitor 基础访问者，包含共享的功能
type BaseVisitor struct {
	totalExpense int    // 总花费
	visitorType  string // 访问者类型
}

// GetTotalExpense 获取总花费
func (bv *BaseVisitor) GetTotalExpense() int {
	return bv.totalExpense
}

// GetVisitorType 获取访问者类型
func (bv *BaseVisitor) GetVisitorType() string {
	return bv.visitorType
}

// StudentVisitor 学生访问者
type StudentVisitor struct {
	BaseVisitor
	hasStudentID bool // 是否持有学生证
}

// NewStudentVisitor 创建一个学生访问者
func NewStudentVisitor(hasStudentID bool) *StudentVisitor {
	return &StudentVisitor{
		BaseVisitor: BaseVisitor{
			totalExpense: 0,
			visitorType:  "学生",
		},
		hasStudentID: hasStudentID,
	}
}

// calculateDiscount 计算学生折扣
func (s *StudentVisitor) calculateDiscount(originalPrice int) int {
	if s.hasStudentID {
		return originalPrice / 2 // 持有学生证半价
	}
	return int(float64(originalPrice) * 0.8) // 无学生证8折
}

// VisitLeopardSpot 学生访问豹子馆
func (s *StudentVisitor) VisitLeopardSpot(leopard *LeopardSpot) {
	price := s.calculateDiscount(leopard.Price())
	s.totalExpense += price
	fmt.Printf("学生游客参观%s，详情: %s，票价: %d元 (原价: %d元)\n",
		leopard.GetName(), leopard.GetDescription(), price, leopard.Price())
}

// VisitDolphinSpot 学生访问海豚馆
func (s *StudentVisitor) VisitDolphinSpot(dolphin *DolphinSpot) {
	price := s.calculateDiscount(dolphin.Price())
	s.totalExpense += price
	showInfo := ""
	if dolphin.HasShow() {
		showInfo = "，今日有精彩表演"
	}
	fmt.Printf("学生游客参观%s，详情: %s%s，票价: %d元 (原价: %d元)\n",
		dolphin.GetName(), dolphin.GetDescription(), showInfo, price, dolphin.Price())
}

// VisitAquarium 学生访问水族馆
func (s *StudentVisitor) VisitAquarium(aquarium *Aquarium) {
	price := s.calculateDiscount(aquarium.Price())
	s.totalExpense += price
	vipInfo := ""
	if aquarium.HasVipArea() {
		vipInfo = "，包含VIP珍稀鱼类区域"
	}
	fmt.Printf("学生游客参观%s，详情: %s%s，票价: %d元 (原价: %d元)\n",
		aquarium.GetName(), aquarium.GetDescription(), vipInfo, price, aquarium.Price())
}

// CommonVisitor 普通游客
type CommonVisitor struct {
	BaseVisitor
	isWeekend bool // 是否周末参观
}

// NewCommonVisitor 创建一个普通访问者
func NewCommonVisitor(isWeekend bool) *CommonVisitor {
	return &CommonVisitor{
		BaseVisitor: BaseVisitor{
			totalExpense: 0,
			visitorType:  "普通",
		},
		isWeekend: isWeekend,
	}
}

// calculatePrice 计算普通游客票价（周末可能上浮）
func (c *CommonVisitor) calculatePrice(originalPrice int) int {
	if c.isWeekend {
		return int(float64(originalPrice) * 1.2) // 周末上浮20%
	}
	return originalPrice
}

// VisitLeopardSpot 普通游客访问豹子馆
func (c *CommonVisitor) VisitLeopardSpot(leopard *LeopardSpot) {
	price := c.calculatePrice(leopard.Price())
	c.totalExpense += price
	fmt.Printf("普通游客参观%s，详情: %s，票价: %d元\n",
		leopard.GetName(), leopard.GetDescription(), price)
}

// VisitDolphinSpot 普通游客访问海豚馆
func (c *CommonVisitor) VisitDolphinSpot(dolphin *DolphinSpot) {
	price := c.calculatePrice(dolphin.Price())
	c.totalExpense += price
	showInfo := ""
	if dolphin.HasShow() {
		showInfo = "，今日有精彩表演"
	}
	fmt.Printf("普通游客参观%s，详情: %s%s，票价: %d元\n",
		dolphin.GetName(), dolphin.GetDescription(), showInfo, price)
}

// VisitAquarium 普通游客访问水族馆
func (c *CommonVisitor) VisitAquarium(aquarium *Aquarium) {
	price := c.calculatePrice(aquarium.Price())
	c.totalExpense += price
	vipInfo := ""
	if aquarium.HasVipArea() {
		vipInfo = "，包含VIP珍稀鱼类区域"
	}
	fmt.Printf("普通游客参观%s，详情: %s%s，票价: %d元\n",
		aquarium.GetName(), aquarium.GetDescription(), vipInfo, price)
}

// VIPVisitor VIP游客
type VIPVisitor struct {
	BaseVisitor
	vipLevel int // VIP等级 1-3
}

// NewVIPVisitor 创建一个VIP访问者
func NewVIPVisitor(vipLevel int) *VIPVisitor {
	if vipLevel < 1 {
		vipLevel = 1
	} else if vipLevel > 3 {
		vipLevel = 3
	}
	return &VIPVisitor{
		BaseVisitor: BaseVisitor{
			totalExpense: 0,
			visitorType:  fmt.Sprintf("VIP-%d", vipLevel),
		},
		vipLevel: vipLevel,
	}
}

// calculateDiscount 计算VIP折扣
func (v *VIPVisitor) calculateDiscount(originalPrice int) int {
	discount := 1.0
	switch v.vipLevel {
	case 1:
		discount = 0.9 // 9折
	case 2:
		discount = 0.8 // 8折
	case 3:
		discount = 0.7 // 7折
	}
	return int(float64(originalPrice) * discount)
}

// VisitLeopardSpot VIP游客访问豹子馆
func (v *VIPVisitor) VisitLeopardSpot(leopard *LeopardSpot) {
	price := v.calculateDiscount(leopard.Price())
	v.totalExpense += price
	fmt.Printf("VIP-%d游客参观%s，详情: %s，享受专属讲解，票价: %d元 (原价: %d元)\n",
		v.vipLevel, leopard.GetName(), leopard.GetDescription(), price, leopard.Price())
}

// VisitDolphinSpot VIP游客访问海豚馆
func (v *VIPVisitor) VisitDolphinSpot(dolphin *DolphinSpot) {
	price := v.calculateDiscount(dolphin.Price())
	v.totalExpense += price
	showInfo := ""
	if dolphin.HasShow() {
		showInfo = "，安排前排观看表演"
	}
	fmt.Printf("VIP-%d游客参观%s，详情: %s%s，票价: %d元 (原价: %d元)\n",
		v.vipLevel, dolphin.GetName(), dolphin.GetDescription(), showInfo, price, dolphin.Price())
}

// VisitAquarium VIP游客访问水族馆
func (v *VIPVisitor) VisitAquarium(aquarium *Aquarium) {
	price := v.calculateDiscount(aquarium.Price())
	v.totalExpense += price
	vipInfo := ""
	if aquarium.HasVipArea() {
		vipInfo = "，专享VIP区域导览"
	}
	fmt.Printf("VIP-%d游客参观%s，详情: %s%s，票价: %d元 (原价: %d元)\n",
		v.vipLevel, aquarium.GetName(), aquarium.GetDescription(), vipInfo, price, aquarium.Price())
}
