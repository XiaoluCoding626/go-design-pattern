# 访问者模式 (Visitor Pattern)

## 介绍

访问者模式是一种行为设计模式，它允许在不修改对象结构的情况下定义对象新的操作。这种模式适用于对象结构相对稳定，但需要经常添加新的操作的场景。

本实现通过"动物园游览系统"展示访问者模式，包含不同类型的景点（豹子馆、海豚馆、水族馆）和不同类型的游客（学生、普通游客、VIP游客）。每类游客访问不同景点时有不同的行为和票价计算方式。

## 设计结构

### 核心接口

- **Visitor（访问者）**: 定义对每种元素的访问操作
- **Scenery（场景元素）**: 定义接受访问者的接口
- **Zoo（对象结构）**: 包含多个场景元素，允许访问者访问其中的所有元素

### 关键组件

1. **访问者接口与实现**:
   - `Visitor`接口: 定义访问不同景点的方法
   - 具体访问者: `StudentVisitor`、`CommonVisitor`、`VIPVisitor`

2. **景点接口与实现**:
   - `Scenery`接口: 定义接受访问者的方法
   - 具体景点: `LeopardSpot`、`DolphinSpot`、`Aquarium`

3. **对象结构**:
   - `Zoo`: 管理多个景点，并提供让访问者访问所有景点的方法

## 类图关系

```
+---------------+     visits    +---------------+
|    Visitor    |<------------->|    Scenery    |
+---------------+               +---------------+
       ^                               ^
       |                               |
+------+------+               +--------+-------+
|             |               |                |
| StudentVisitor  CommonVisitor  VIPVisitor    |    LeopardSpot  DolphinSpot  Aquarium
+-------------+               +----------------+
                                    ^
                                    |
                               +----+----+
                               |   Zoo   |
                               +---------+
```

## 代码示例

### 访问者接口

```go
// Visitor 抽象访问者接口 - 定义对每种场景的访问方法
type Visitor interface {
    VisitLeopardSpot(leopard *LeopardSpot) 
    VisitDolphinSpot(dolphin *DolphinSpot)
    VisitAquarium(aquarium *Aquarium)
    GetTotalExpense() int
    GetVisitorType() string
}
```

### 景点接口

```go
// Scenery 场馆景点接口 - 定义场景对象的通用行为
type Scenery interface {
    Accept(visitor Visitor) // 接待访问者
    Price() int             // 基础票价
    GetName() string        // 获取景点名称
    GetDescription() string // 获取景点描述
}
```

### 具体访问者示例 - 学生访问者

```go
// StudentVisitor 学生访问者
type StudentVisitor struct {
    BaseVisitor
    hasStudentID bool // 是否持有学生证
}

// 计算学生折扣
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

// 其他访问方法...
```

### 具体景点示例 - 豹子馆

```go
// LeopardSpot 豹子馆实现
type LeopardSpot struct {
    description string
    basePrice   int
}

// Accept 实现Scenery接口的Accept方法
func (l *LeopardSpot) Accept(visitor Visitor) {
    visitor.VisitLeopardSpot(l)
}

// Price 基础票价
func (l *LeopardSpot) Price() int {
    return l.basePrice
}

// 其他方法...
```

### 对象结构 - 动物园

```go
// Zoo 动物园类 - 复合对象，包含多个景点
type Zoo struct {
    Name      string     // 动物园名称
    Sceneries []Scenery  // 动物园包含的景点
    OpenTime  *time.Time // 开放时间
}

// Accept 动物园接待游客，游客将参观所有景点
func (z *Zoo) Accept(v Visitor) {
    fmt.Printf("\n%s 欢迎 %s 游客参观！\n", z.Name, v.GetVisitorType())
    for _, scenery := range z.Sceneries {
        scenery.Accept(v)
    }
    fmt.Printf("%s 游客参观完成，总花费: %d 元\n", v.GetVisitorType(), v.GetTotalExpense())
}

// 其他方法...
```

## 使用场景

访问者模式适用于以下场景：

1. **对象结构稳定但操作多变**：景点类型相对固定，但不同游客的访问行为各异
2. **需要对不同类型对象执行不同操作**：学生、普通游客和VIP对景点有不同行为和票价计算方式
3. **需要对对象结构中的元素执行一系列不相关的操作**：计算票价、提供特殊服务等

## 优缺点

### 优点

1. **单一职责原则**：将对象结构与操作分离
2. **开闭原则**：添加新访问者不需要修改现有对象结构
3. **双分派机制**：根据访问者类型和元素类型决定具体行为

### 缺点

1. **元素变更困难**：添加新元素需要修改所有访问者
2. **访问者可能破坏封装**：访问者可能需要访问元素的内部状态
3. **复杂度增加**：模式引入了多个接口和类，使系统更加复杂

## 运行测试

该项目包含完整的测试用例，你可以通过以下命令运行测试：

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestStudentVisitor

# 运行基准测试
go test -bench=BenchmarkVisitorPerformance
```

## 总结

访问者模式为我们提供了一种优雅的方式，在不修改对象结构的情况下添加新的操作。在本例中，我们可以轻松添加新的游客类型（如"老年游客"或"团队游客"）而不需要修改任何景点代码。同样，如果游客类型稳定，我们也可以轻松添加新的景点而只需要在每个游客类型中实现相应的访问方法。