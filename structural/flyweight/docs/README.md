# 享元设计模式 (Flyweight Pattern)

## 简介

享元模式是一种结构型设计模式，主要用于减少创建对象的数量，以减少内存占用和提高性能。当系统中有大量相似对象时，此模式特别有用。

享元模式将对象的状态分为两部分：
- **内部状态（Intrinsic State）**：可以共享的、不变的状态
- **外部状态（Extrinsic State）**：根据具体场景变化的、不可共享的状态

## 核心概念

1. **享元接口（Flyweight）**：定义了享元类的公共方法
2. **具体享元类（Concrete Flyweight）**：实现享元接口，存储内部状态
3. **享元工厂（Flyweight Factory）**：创建并管理享元对象池
4. **客户端（Client）**：维护外部状态，使用享元对象

## 实现示例

本例通过一个游戏场景实现享元模式，其中游戏中的角色皮肤（Dress）作为可共享的享元对象。

### 类图结构

```
                    ┌─────────────┐
                    │    Dress    │
                    │   Interface │
                    └──────┬──────┘
                           │
                           │
              ┌────────────┴─────────────┐
              │                          │
     ┌────────▼─────────┐      ┌────────▼────────┐
     │  ConcreteDress   │      │   DressFactory  │
     │   Base Class     │      │                 │
     └────────┬─────────┘      └────────┬────────┘
              │                         │
              │                         │manages
    ┌─────────┼──────────┐              │
    │         │          │              │
┌───▼───┐ ┌───▼───┐ ┌────▼────┐    ┌───▼───┐     ┌───────┐
│Terrorist│ │Counter│ │ Elite  │    │ Game  │◄────│Player │
│ Dress  │ │Dress  │ │ Dress  │    │       │     │       │
└────────┘ └───────┘ └─────────┘    └───────┘     └───────┘
```

### 关键组件

#### 1. 享元接口和基类

```go
// Dress 是享元接口，定义了所有具体享元类需要实现的方法
type Dress interface {
    GetColor() string                                  // 获取皮肤颜色
    Display(playerID int, playerName string, x, y int) // 显示玩家信息（内部状态+外部状态）
}

// ConcreteDress 是具体享元类的实现基础，包含共享的内部状态
type ConcreteDress struct {
    color     string // 内部状态 - 颜色是可共享的
    textureID int    // 内部状态 - 纹理ID是可共享的
    meshType  string // 内部状态 - 网格类型是可共享的
}
```

#### 2. 具体享元类

```go
// TerroristDress 是恐怖分子皮肤类，实现了Dress接口
type TerroristDress struct {
    ConcreteDress
}

// CounterTerroristDress 是反恐精英皮肤类，实现了Dress接口
type CounterTerroristDress struct {
    ConcreteDress
}

// EliteDress 是精英特种部队皮肤类，实现了Dress接口
type EliteDress struct {
    ConcreteDress
}
```

#### 3. 享元工厂

```go
// DressFactory 是享元工厂，负责创建和管理享元对象
type DressFactory struct {
    dresses map[string]Dress // 享元对象池
    count   map[string]int   // 跟踪每种皮肤使用次数
}

// GetDress 根据类型获取享元对象，如果不存在则创建
func (f *DressFactory) GetDress(dressType string) (Dress, error) {
    // 检查是否已有此类皮肤对象，如有则复用
    if dress, exists := f.dresses[dressType]; exists {
        f.count[dressType]++
        return dress, nil
    }
    
    // 如果没有，则根据类型创建新的享元对象
    var dress Dress
    switch dressType {
    case TerroristDressType:
        dress = NewTerroristDress()
    case CounterTerroristDressType:
        dress = NewCounterTerroristDress()
    case EliteDressType:
        dress = NewEliteDress()
    default:
        return nil, fmt.Errorf("不支持的皮肤类型: %s", dressType)
    }
    
    // 将新创建的享元对象存入池中
    f.dresses[dressType] = dress
    f.count[dressType] = 1
    return dress, nil
}
```

#### 4. 客户端 - 玩家类

```go
// Player 表示游戏中的玩家，包含外部状态（extrinsic state）
type Player struct {
    id         int    // 外部状态 - 玩家ID是每个玩家特有的
    name       string // 外部状态 - 玩家名字是每个玩家特有的
    dress      Dress  // 引用享元对象（内部状态）
    playerType string // 外部状态 - 玩家类型
    x, y       int    // 外部状态 - 玩家位置坐标
}
```

#### 5. 游戏会话 - 上下文

```go
// Game 代表一个游戏会话，管理所有玩家
type Game struct {
    players   []*Player      // 所有玩家列表
    factory   *DressFactory  // 皮肤工厂
    teamCount map[string]int // 每个团队的玩家数量
}
```

## 使用示例

### 基本用法

```go
// 创建游戏会话
game := NewGame()

// 添加不同类型的玩家
game.AddTerroristPlayer("T1", 10, 20)
game.AddCounterTerroristPlayer("CT1", 30, 40)
game.AddElitePlayer("E1", 50, 60)

// 显示所有玩家
game.DisplayPlayers()

// 显示内存使用统计
game.DisplayMemoryUsage()
```

### 模拟游戏示例

```go
// 创建一个包含多个玩家的游戏
game := SimulateGame(5)

// 输出：
// 当前游戏中的所有玩家:
// 玩家 #1 (Player1) 使用 蓝色 皮肤 (纹理ID: 201, 网格类型: 城市战术) 位于坐标 (10,5)
// 玩家 #2 (Player2) 使用 黑色 皮肤 (纹理ID: 301, 网格类型: 高级战术) 位于坐标 (20,10)
// 玩家 #3 (Player3) 使用 红色 皮肤 (纹理ID: 101, 网格类型: 沙漠迷彩) 位于坐标 (30,15)
// 玩家 #4 (Player4) 使用 蓝色 皮肤 (纹理ID: 201, 网格类型: 城市战术) 位于坐标 (40,20)
// 玩家 #5 (Player5) 使用 黑色 皮肤 (纹理ID: 301, 网格类型: 高级战术) 位于坐标 (50,25)
//
// 内存使用统计:
// 总玩家数: 5
// 唯一皮肤对象数: 3
// 节省的对象数: 2
// 内存节省比例: 40.00%
//
// 各类皮肤使用统计:
// 恐怖分子皮肤: 被 1 名玩家使用
// 反恐精英皮肤: 被 2 名玩家使用
// 精英部队皮肤: 被 2 名玩家使用
```

## 享元模式的优势

1. **减少内存使用**：通过共享对象减少内存消耗，特别是在处理大量相似对象时
2. **提高性能**：减少对象创建和垃圾回收的开销
3. **支持大规模对象系统**：使得系统能够支持更多的对象

## 享元模式的注意事项

1. **状态划分**：需要正确区分内部状态和外部状态
2. **线程安全**：共享对象需要考虑并发安全问题
3. **复杂度增加**：实现可能增加代码复杂度

## 适用场景

享元模式在以下场景特别有用：

1. 应用程序使用大量相似对象
2. 对象存储开销大
3. 大多数对象状态可以设为外部状态
4. 移除了外部状态后，可以用较少的共享对象代替原有对象

## 在Go中实现享元模式的特点

1. **接口**：使用Go的接口定义享元对象的行为
2. **结构体嵌入**：通过结构体嵌入实现代码复用
3. **工厂方法**：使用工厂方法创建和管理享元对象
4. **Map作为对象池**：使用map存储和查找享元对象

## 性能对比

本实现包含两个基准测试，对比了使用享元模式和不使用享元模式的性能差异：

- `BenchmarkWithFlyweight`：使用享元模式创建1000个玩家
- `BenchmarkWithoutFlyweight`：不使用享元模式，每个玩家创建独立的皮肤对象

在大规模对象创建场景下，享元模式通常会表现出明显的性能优势和内存使用效率。
