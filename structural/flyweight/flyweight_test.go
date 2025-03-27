package flyweight

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

// 捕获标准输出的辅助函数
func captureOutput(f func()) string {
	// 保存原始的标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行需要捕获输出的函数
	f()

	// 恢复原始的标准输出
	w.Close()
	os.Stdout = oldStdout

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// 模拟不使用享元模式的对象创建，用于对比
type NonFlyweightDress struct {
	color      string
	textureID  int
	meshType   string
	playerID   int
	playerName string
	x, y       int
}

// TestDressCreation 测试享元对象的创建
func TestDressCreation(t *testing.T) {
	tests := []struct {
		name           string
		createDress    func() Dress
		expectedColor  string
		expectedOutput string
	}{
		{
			name:           "TerroristDress",
			createDress:    func() Dress { return NewTerroristDress() },
			expectedColor:  "红色",
			expectedOutput: "红色.*101.*沙漠迷彩",
		},
		{
			name:           "CounterTerroristDress",
			createDress:    func() Dress { return NewCounterTerroristDress() },
			expectedColor:  "蓝色",
			expectedOutput: "蓝色.*201.*城市战术",
		},
		{
			name:           "EliteDress",
			createDress:    func() Dress { return NewEliteDress() },
			expectedColor:  "黑色",
			expectedOutput: "黑色.*301.*高级战术",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dress := test.createDress()

			// 测试颜色值
			if color := dress.GetColor(); color != test.expectedColor {
				t.Errorf("期望颜色为 %s，实际得到 %s", test.expectedColor, color)
			}

			// 测试显示方法
			output := captureOutput(func() {
				dress.Display(1, "TestPlayer", 10, 20)
			})

			match, _ := regexp.MatchString(test.expectedOutput, output)
			if !match {
				t.Errorf("期望输出包含 %s，实际输出为 %s", test.expectedOutput, output)
			}
		})
	}
}

// TestDressFactory 测试享元工厂的对象共享功能
func TestDressFactory(t *testing.T) {
	factory := NewDressFactory()

	// 测试创建新对象
	dress1, err := factory.GetDress(TerroristDressType)
	if err != nil {
		t.Fatalf("获取 TerroristDress 失败: %v", err)
	}

	if count := factory.GetDressCount(TerroristDressType); count != 1 {
		t.Errorf("期望计数为 1，实际为 %d", count)
	}

	// 测试复用已有对象
	dress2, _ := factory.GetDress(TerroristDressType)

	// 验证返回的是同一个对象实例（享元模式的核心）
	if dress1 != dress2 {
		t.Error("工厂没有返回同一个对象实例，享元模式失效")
	}

	if count := factory.GetDressCount(TerroristDressType); count != 2 {
		t.Errorf("期望计数为 2，实际为 %d", count)
	}

	// 测试获取不支持的皮肤类型
	_, err = factory.GetDress("InvalidType")
	if err == nil {
		t.Error("期望获取无效皮肤类型时返回错误，但没有")
	}

	// 测试总对象计数
	factory.GetDress(CounterTerroristDressType)
	factory.GetDress(EliteDressType)

	if count := factory.GetTotalDressCount(); count != 3 {
		t.Errorf("期望总皮肤对象数为 3，实际为 %d", count)
	}
}

// TestPlayer 测试玩家对象创建与显示
func TestPlayer(t *testing.T) {
	factory := NewDressFactory()

	player, err := NewPlayer(1, "TestPlayer", "Terrorist", TerroristDressType, factory, 10, 20)
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}

	output := captureOutput(func() {
		player.Display()
	})

	expectedOutputParts := []string{
		"玩家 #1",
		"TestPlayer",
		"红色",
		"纹理ID: 101",
		"沙漠迷彩",
		"坐标 (10,20)",
	}

	for _, part := range expectedOutputParts {
		if !strings.Contains(output, part) {
			t.Errorf("玩家显示输出应包含 '%s'，但输出为: %s", part, output)
		}
	}
}

// TestGame 测试游戏会话功能
func TestGame(t *testing.T) {
	game := NewGame()

	// 测试添加不同类型的玩家
	err1 := game.AddTerroristPlayer("T1", 10, 20)
	err2 := game.AddCounterTerroristPlayer("CT1", 30, 40)
	err3 := game.AddElitePlayer("E1", 50, 60)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("添加玩家失败: %v, %v, %v", err1, err2, err3)
	}

	if len(game.players) != 3 {
		t.Errorf("期望有 3 个玩家，实际有 %d 个", len(game.players))
	}

	// 测试添加无效团队类型
	err := game.AddPlayer("Invalid", "InvalidTeam", 0, 0)
	if err == nil {
		t.Error("期望添加无效团队类型时返回错误，但没有")
	}

	// 测试工厂状态
	factory := game.factory
	if factory.GetTotalDressCount() != 3 {
		t.Errorf("期望有 3 种皮肤，实际有 %d 种", factory.GetTotalDressCount())
	}

	// 添加更多相同类型的玩家，验证享元对象复用
	game.AddTerroristPlayer("T2", 15, 25)
	game.AddTerroristPlayer("T3", 16, 26)

	if factory.GetDressCount(TerroristDressType) != 3 {
		t.Errorf("期望恐怖分子皮肤使用 3 次，实际为 %d 次",
			factory.GetDressCount(TerroristDressType))
	}

	// 测试显示所有玩家
	output := captureOutput(func() {
		game.DisplayPlayers()
	})

	expectedPlayers := []string{"T1", "CT1", "E1", "T2", "T3"}
	for _, player := range expectedPlayers {
		if !strings.Contains(output, player) {
			t.Errorf("玩家显示输出应包含 '%s'，但没有找到", player)
		}
	}
}

// TestMemoryUsage 测试内存使用统计功能
func TestMemoryUsage(t *testing.T) {
	// 创建一个有 15 个玩家的游戏
	game := SimulateGame(15)

	output := captureOutput(func() {
		game.DisplayMemoryUsage()
	})

	// 验证内存使用统计
	if !strings.Contains(output, "总玩家数: 15") {
		t.Errorf("内存使用统计应显示总玩家数为 15")
	}

	if !strings.Contains(output, "唯一皮肤对象数: 3") {
		t.Errorf("内存使用统计应显示唯一皮肤对象数为 3")
	}

	// 节省的对象数应该是 15 - 3 = 12
	if !strings.Contains(output, "节省的对象数: 12") {
		t.Errorf("内存使用统计应显示节省的对象数为 12")
	}
}

// BenchmarkWithFlyweight 测试使用享元模式时的性能
func BenchmarkWithFlyweight(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game := NewGame()
		// 添加大量玩家，测试享元模式下的性能
		for j := 0; j < 1000; j++ {
			teamType := j % 3
			playerName := fmt.Sprintf("Player%d", j)
			x, y := j*10, j*5

			switch teamType {
			case 0:
				game.AddTerroristPlayer(playerName, x, y)
			case 1:
				game.AddCounterTerroristPlayer(playerName, x, y)
			case 2:
				game.AddElitePlayer(playerName, x, y)
			}
		}
	}
}

// BenchmarkWithoutFlyweight 测试不使用享元模式时的性能，用于对比
func BenchmarkWithoutFlyweight(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 不使用享元模式，直接创建每个对象的全部数据
		dresses := make([]*NonFlyweightDress, 0, 1000)

		for j := 0; j < 1000; j++ {
			teamType := j % 3
			playerName := fmt.Sprintf("Player%d", j)
			x, y := j*10, j*5

			var dress *NonFlyweightDress

			switch teamType {
			case 0:
				dress = &NonFlyweightDress{
					color:      "红色",
					textureID:  101,
					meshType:   "沙漠迷彩",
					playerID:   j,
					playerName: playerName,
					x:          x,
					y:          y,
				}
			case 1:
				dress = &NonFlyweightDress{
					color:      "蓝色",
					textureID:  201,
					meshType:   "城市战术",
					playerID:   j,
					playerName: playerName,
					x:          x,
					y:          y,
				}
			case 2:
				dress = &NonFlyweightDress{
					color:      "黑色",
					textureID:  301,
					meshType:   "高级战术",
					playerID:   j,
					playerName: playerName,
					x:          x,
					y:          y,
				}
			}

			dresses = append(dresses, dress)
		}
	}
}

// 示例：展示享元模式使用
func ExampleSimulateGame() {
	// 创建一个有 5 个玩家的游戏
	game := SimulateGame(5)

	// 显示所有玩家信息
	game.DisplayPlayers()

	// 获取内存使用信息但不直接显示(因为map迭代顺序不确定)
	totalPlayers := len(game.players)
	uniqueDresses := game.factory.GetTotalDressCount()
	savedObjects := totalPlayers - uniqueDresses
	savingPercentage := float64(savedObjects) / float64(totalPlayers) * 100

	// 按照固定顺序显示内存统计信息
	fmt.Printf("\n内存使用统计:\n")
	fmt.Printf("总玩家数: %d\n", totalPlayers)
	fmt.Printf("唯一皮肤对象数: %d\n", uniqueDresses)
	fmt.Printf("节省的对象数: %d\n", savedObjects)
	fmt.Printf("内存节省比例: %.2f%%\n", savingPercentage)

	// 按照固定顺序显示皮肤使用统计
	fmt.Println("\n各类皮肤使用统计:")
	fmt.Printf("恐怖分子皮肤: 被 %d 名玩家使用\n",
		game.factory.GetDressCount(TerroristDressType))
	fmt.Printf("反恐精英皮肤: 被 %d 名玩家使用\n",
		game.factory.GetDressCount(CounterTerroristDressType))
	fmt.Printf("精英部队皮肤: 被 %d 名玩家使用\n",
		game.factory.GetDressCount(EliteDressType))

	// Output:
	//
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
}

// 示例：演示享元工厂如何复用对象
func ExampleDressFactory_GetDress() {
	factory := NewDressFactory()

	// 第一次获取恐怖分子皮肤 - 创建新对象
	dress1, _ := factory.GetDress(TerroristDressType)
	fmt.Printf("创建第一个恐怖分子皮肤: %s\n", dress1.GetColor())

	// 第二次获取恐怖分子皮肤 - 复用已有对象
	dress2, _ := factory.GetDress(TerroristDressType)
	fmt.Printf("获取第二个恐怖分子皮肤: %s\n", dress2.GetColor())

	// 显示使用统计
	fmt.Printf("恐怖分子皮肤使用次数: %d\n", factory.GetDressCount(TerroristDressType))
	fmt.Printf("唯一皮肤对象总数: %d\n", factory.GetTotalDressCount())

	// Output:
	// 创建第一个恐怖分子皮肤: 红色
	// 获取第二个恐怖分子皮肤: 红色
	// 恐怖分子皮肤使用次数: 2
	// 唯一皮肤对象总数: 1
}

// TestPlayerCreationError 测试创建玩家时的错误处理
func TestPlayerCreationError(t *testing.T) {
	factory := NewDressFactory()

	_, err := NewPlayer(1, "ErrorPlayer", "Terrorist", "InvalidType", factory, 0, 0)
	if err == nil {
		t.Error("使用无效皮肤类型创建玩家时应返回错误")
	}
}

// TestAddMultiplePlayers 测试添加大量玩家
func TestAddMultiplePlayers(t *testing.T) {
	game := NewGame()
	playerCount := 100

	// 添加大量玩家
	for i := 0; i < playerCount; i++ {
		teamType := i % 3
		playerName := fmt.Sprintf("Player%d", i)
		x, y := i, i

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
			t.Errorf("添加玩家 %s 失败: %v", playerName, err)
		}
	}

	// 验证玩家数量
	if len(game.players) != playerCount {
		t.Errorf("期望有 %d 个玩家，实际有 %d 个", playerCount, len(game.players))
	}

	// 验证对象共享
	factory := game.factory
	if factory.GetTotalDressCount() != 3 {
		t.Errorf("期望有 3 种皮肤，实际有 %d 种", factory.GetTotalDressCount())
	}

	// 计算各类型玩家数量
	expectedTCount := playerCount / 3
	if playerCount%3 > 0 {
		expectedTCount++
	}

	if factory.GetDressCount(TerroristDressType) < expectedTCount {
		t.Errorf("期望至少有 %d 名恐怖分子玩家，实际有 %d 名",
			expectedTCount, factory.GetDressCount(TerroristDressType))
	}
}
