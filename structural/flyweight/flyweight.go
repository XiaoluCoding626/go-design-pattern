package flyweight

import (
	"fmt"
	"strconv"
)

// Dress 是享元接口，定义了所有具体享元类需要实现的方法
// 它包含的是内部状态（intrinsic state）- 可以被多个对象共享的状态
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

// GetColor 返回皮肤颜色
func (d *ConcreteDress) GetColor() string {
	return d.color
}

// Display 使用内部状态和外部状态显示玩家信息
func (d *ConcreteDress) Display(playerID int, playerName string, x, y int) {
	fmt.Printf("玩家 #%d (%s) 使用 %s 皮肤 (纹理ID: %d, 网格类型: %s) 位于坐标 (%d,%d)\n",
		playerID, playerName, d.color, d.textureID, d.meshType, x, y)
}

// TerroristDress 是恐怖分子皮肤类，实现了Dress接口
type TerroristDress struct {
	ConcreteDress
}

// NewTerroristDress 创建一个新的恐怖分子皮肤享元对象
func NewTerroristDress() *TerroristDress {
	return &TerroristDress{
		ConcreteDress: ConcreteDress{
			color:     "红色",
			textureID: 101,
			meshType:  "沙漠迷彩",
		},
	}
}

// CounterTerroristDress 是反恐精英皮肤类，实现了Dress接口
type CounterTerroristDress struct {
	ConcreteDress
}

// NewCounterTerroristDress 创建一个新的反恐精英皮肤享元对象
func NewCounterTerroristDress() *CounterTerroristDress {
	return &CounterTerroristDress{
		ConcreteDress: ConcreteDress{
			color:     "蓝色",
			textureID: 201,
			meshType:  "城市战术",
		},
	}
}

// EliteDress 是精英特种部队皮肤类，实现了Dress接口（新增类型）
type EliteDress struct {
	ConcreteDress
}

// NewEliteDress 创建一个新的精英特种部队皮肤享元对象
func NewEliteDress() *EliteDress {
	return &EliteDress{
		ConcreteDress: ConcreteDress{
			color:     "黑色",
			textureID: 301,
			meshType:  "高级战术",
		},
	}
}

// 定义皮肤类型常量
const (
	TerroristDressType        = "T"  // 恐怖分子皮肤
	CounterTerroristDressType = "CT" // 反恐精英皮肤
	EliteDressType            = "E"  // 精英部队皮肤
)

// DressFactory 是享元工厂，负责创建和管理享元对象
type DressFactory struct {
	dresses map[string]Dress // 享元对象池
	count   map[string]int   // 跟踪每种皮肤使用次数
}

// NewDressFactory 创建并初始化一个新的皮肤工厂
func NewDressFactory() *DressFactory {
	return &DressFactory{
		dresses: make(map[string]Dress),
		count:   make(map[string]int),
	}
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

// GetDressCount 返回指定类型皮肤的使用次数
func (f *DressFactory) GetDressCount(dressType string) int {
	return f.count[dressType]
}

// GetTotalDressCount 返回所有皮肤对象的总数
func (f *DressFactory) GetTotalDressCount() int {
	return len(f.dresses)
}

// GetDressUsageStats 返回所有皮肤的使用统计
func (f *DressFactory) GetDressUsageStats() map[string]int {
	return f.count
}

// Player 表示游戏中的玩家，包含外部状态（extrinsic state）
type Player struct {
	id         int    // 外部状态 - 玩家ID是每个玩家特有的
	name       string // 外部状态 - 玩家名字是每个玩家特有的
	dress      Dress  // 引用享元对象（内部状态）
	playerType string // 外部状态 - 玩家类型
	x, y       int    // 外部状态 - 玩家位置坐标
}

// NewPlayer 创建一个新的玩家对象
func NewPlayer(id int, name, playerType, dressType string, factory *DressFactory, x, y int) (*Player, error) {
	dress, err := factory.GetDress(dressType)
	if err != nil {
		return nil, err
	}

	return &Player{
		id:         id,
		name:       name,
		dress:      dress,
		playerType: playerType,
		x:          x,
		y:          y,
	}, nil
}

// Display 显示玩家信息，结合内部和外部状态
func (p *Player) Display() {
	p.dress.Display(p.id, p.name, p.x, p.y)
}

// Game 代表一个游戏会话，管理所有玩家
type Game struct {
	players   []*Player      // 所有玩家列表
	factory   *DressFactory  // 皮肤工厂
	teamCount map[string]int // 每个团队的玩家数量
}

// NewGame 创建一个新的游戏实例
func NewGame() *Game {
	return &Game{
		players:   make([]*Player, 0),
		factory:   NewDressFactory(),
		teamCount: make(map[string]int),
	}
}

// AddPlayer 向游戏中添加新玩家
func (g *Game) AddPlayer(name, teamType string, x, y int) error {
	var dressType string

	// 根据团队类型选择合适的皮肤类型
	switch teamType {
	case "Terrorist":
		dressType = TerroristDressType
	case "CounterTerrorist":
		dressType = CounterTerroristDressType
	case "Elite":
		dressType = EliteDressType
	default:
		return fmt.Errorf("未知的团队类型: %s", teamType)
	}

	// 更新团队计数
	if _, exists := g.teamCount[teamType]; !exists {
		g.teamCount[teamType] = 0
	}
	g.teamCount[teamType]++

	// 创建玩家ID
	playerID := len(g.players) + 1

	// 创建玩家
	player, err := NewPlayer(playerID, name, teamType, dressType, g.factory, x, y)
	if err != nil {
		return err
	}

	// 添加到玩家列表
	g.players = append(g.players, player)
	return nil
}

// DisplayPlayers 显示所有玩家信息
func (g *Game) DisplayPlayers() {
	fmt.Println("\n当前游戏中的所有玩家:")
	for _, player := range g.players {
		player.Display()
	}
}

// DisplayMemoryUsage 显示内存使用情况，展示享元模式的节省效果
func (g *Game) DisplayMemoryUsage() {
	totalPlayers := len(g.players)
	uniqueDresses := g.factory.GetTotalDressCount()
	savedObjects := totalPlayers - uniqueDresses

	fmt.Printf("\n内存使用统计:\n")
	fmt.Printf("总玩家数: %d\n", totalPlayers)
	fmt.Printf("唯一皮肤对象数: %d\n", uniqueDresses)
	fmt.Printf("节省的对象数: %d\n", savedObjects)

	if totalPlayers > 0 {
		savingPercentage := float64(savedObjects) / float64(totalPlayers) * 100
		fmt.Printf("内存节省比例: %.2f%%\n", savingPercentage)
	}

	fmt.Println("\n各类皮肤使用统计:")
	for dressType, count := range g.factory.GetDressUsageStats() {
		var typeName string
		switch dressType {
		case TerroristDressType:
			typeName = "恐怖分子皮肤"
		case CounterTerroristDressType:
			typeName = "反恐精英皮肤"
		case EliteDressType:
			typeName = "精英部队皮肤"
		}
		fmt.Printf("%s: 被 %d 名玩家使用\n", typeName, count)
	}
}

// AddTerroristPlayer 添加恐怖分子玩家的便捷方法
func (g *Game) AddTerroristPlayer(name string, x, y int) error {
	return g.AddPlayer(name, "Terrorist", x, y)
}

// AddCounterTerroristPlayer 添加反恐精英玩家的便捷方法
func (g *Game) AddCounterTerroristPlayer(name string, x, y int) error {
	return g.AddPlayer(name, "CounterTerrorist", x, y)
}

// AddElitePlayer 添加精英特种部队玩家的便捷方法
func (g *Game) AddElitePlayer(name string, x, y int) error {
	return g.AddPlayer(name, "Elite", x, y)
}

// SimulateGame 模拟游戏示例
func SimulateGame(playerCount int) *Game {
	game := NewGame()

	// 添加玩家
	for i := 1; i <= playerCount; i++ {
		teamType := i % 3
		playerName := "Player" + strconv.Itoa(i)
		x, y := i*10, i*5

		var err error
		switch teamType {
		case 0:
			err = game.AddTerroristPlayer(playerName, x, y)
		case 1:
			err = game.AddCounterTerroristPlayer(playerName, x, y)
		case 2:
			err = game.AddElitePlayer(playerName, x, y)
		}

		if err != nil {
			fmt.Printf("添加玩家失败: %v\n", err)
		}
	}

	return game
}
