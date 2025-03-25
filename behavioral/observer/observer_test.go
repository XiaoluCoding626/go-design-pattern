package observer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// captureOutput 捕获标准输出的辅助函数
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestStockEvent 测试股票事件的基本功能
func TestStockEvent(t *testing.T) {
	assert := assert.New(t)

	event := StockEvent{
		Symbol:    "AAPL",
		Price:     150.0,
		PrevPrice: 145.0,
		Timestamp: time.Now(),
	}

	// 测试价格变动百分比计算
	expectedChangePercent := (150.0 - 145.0) / 145.0 * 100
	assert.Equal(expectedChangePercent, event.ChangePercent(), "ChangePercent 计算错误")

	// 测试涨跌判断
	assert.True(event.IsUp(), "IsUp 判断错误，应该上涨但返回下跌")

	// 测试价格变动阈值判断
	assert.True(event.IsPriceChange(2.0), "IsPriceChange 判断错误，价格变动超过阈值但未能正确判断")  // 阈值2.0%，变动约3.45%
	assert.False(event.IsPriceChange(5.0), "IsPriceChange 判断错误，价格变动未超过阈值但判断为超过") // 阈值5.0%，变动约3.45%

	// 测试字符串表示
	str := event.String()
	assert.Contains(str, "AAPL", "String 方法输出应包含股票代码")
	assert.Contains(str, "上涨", "String 方法输出应包含涨跌方向")
}

// TestStockMarket 测试股票市场的基本功能
func TestStockMarket(t *testing.T) {
	assert := assert.New(t)
	market := NewStockMarket()

	// 测试初始状态
	assert.Equal(0, market.CountObservers(), "新建市场应该没有观察者")

	// 创建观察者
	investor := NewInvestor("inv1", "张三", Moderate)

	// 测试注册观察者
	output := captureOutput(func() {
		market.Register(investor)
	})
	assert.Contains(output, "已注册到股票市场", "观察者注册输出不正确")
	assert.Equal(1, market.CountObservers(), "注册后市场应该有1个观察者")

	// 测试重复注册
	output = captureOutput(func() {
		market.Register(investor)
	})
	assert.Contains(output, "已经注册", "重复注册输出不正确")
	assert.Equal(1, market.CountObservers(), "重复注册后市场应该仍有1个观察者")

	// 测试检查观察者是否存在
	assert.True(market.HasObserver(investor), "HasObserver 方法未能识别已注册的观察者")

	// 测试注销观察者
	output = captureOutput(func() {
		market.Deregister(investor)
	})
	assert.Contains(output, "已从股票市场注销", "观察者注销输出不正确")
	assert.Equal(0, market.CountObservers(), "注销后市场应该没有观察者")
	assert.False(market.HasObserver(investor), "HasObserver 方法错误识别了已注销的观察者")
}

// TestStockPriceUpdate 测试股票价格更新和通知
func TestStockPriceUpdate(t *testing.T) {
	assert := assert.New(t)
	market := NewStockMarket()
	investor := NewInvestor("inv1", "张三", Moderate)
	analyst := NewMarketAnalyst("anl1", "李四", "某证券公司")

	market.Register(investor)
	market.Register(analyst)

	// 测试股票价格初始更新
	output := captureOutput(func() {
		market.UpdateStockPrice("AAPL", 150.0, "苹果公司股票价格更新", 0.1)
	})
	assert.Contains(output, "苹果公司股票价格更新", "价格更新通知输出不正确")
	assert.Contains(output, "AAPL", "通知中缺少股票代码")

	// 检查价格是否正确存储
	price, exists := market.GetStockPrice("AAPL")
	assert.True(exists, "应该存在股票AAPL但GetStockPrice返回不存在")
	assert.Equal(150.0, price, "股票价格存储不正确")

	// 测试阈值通知机制（小于阈值不应通知）
	output = captureOutput(func() {
		market.UpdateStockPrice("AAPL", 150.1, "苹果公司股票小幅变动", 1.0)
	})
	// 价格变动约0.067%，低于阈值1.0%，不应该有通知
	assert.NotContains(output, "苹果公司股票小幅变动", "低于阈值的价格变动不应通知")

	// 测试阈值通知机制（大于阈值应通知）
	output = captureOutput(func() {
		market.UpdateStockPrice("AAPL", 160.0, "苹果公司股票大幅上涨", 1.0)
	})
	// 价格变动约6.6%，高于阈值1.0%，应该通知
	assert.Contains(output, "苹果公司股票大幅上涨", "高于阈值的价格变动应该通知")
}

// TestDifferentInvestorTypes 测试不同类型投资者的反应
func TestDifferentInvestorTypes(t *testing.T) {
	assert := assert.New(t)
	market := NewStockMarket()

	// 创建三种类型的投资者
	conservativeInvestor := NewInvestor("con1", "保守张", Conservative)
	moderateInvestor := NewInvestor("mod1", "稳健王", Moderate)
	aggressiveInvestor := NewInvestor("agg1", "激进李", Aggressive)

	market.Register(conservativeInvestor)
	market.Register(moderateInvestor)
	market.Register(aggressiveInvestor)

	// 测试不同涨幅对不同投资者的影响

	// 小幅上涨2%（应该仅激进型买入）
	output := captureOutput(func() {
		market.UpdateStockPrice("GOOGL", 1000.0, "谷歌股票初始价格", 0.1)
		market.UpdateStockPrice("GOOGL", 1020.0, "谷歌股票小幅上涨2%", 0.1)
	})
	assert.Contains(output, "激进李(激进型): 买入", "激进型投资者应该在2%涨幅时买入")
	assert.NotContains(output, "保守张(保守型): 买入", "保守型投资者不应该在2%涨幅时买入")
	assert.NotContains(output, "稳健王(稳健型): 买入", "稳健型投资者不应该在2%涨幅时买入")

	// 中等上涨4%（应该激进型和稳健型买入）
	output = captureOutput(func() {
		market.UpdateStockPrice("GOOGL", 1060.0, "谷歌股票上涨4%", 0.1)
	})
	assert.Contains(output, "激进李(激进型): 买入", "激进型投资者应该在4%涨幅时买入")
	assert.Contains(output, "稳健王(稳健型): 买入", "稳健型投资者应该在4%涨幅时买入")
	assert.NotContains(output, "保守张(保守型): 买入", "保守型投资者不应该在4%涨幅时买入")

	// 大幅上涨8%（所有投资者都应该买入）
	output = captureOutput(func() {
		market.UpdateStockPrice("GOOGL", 1144.8, "谷歌股票大涨8%", 0.1)
	})
	assert.Contains(output, "激进李(激进型): 买入", "激进型投资者应该在8%涨幅时买入")
	assert.Contains(output, "稳健王(稳健型): 买入", "稳健型投资者应该在8%涨幅时买入")
	assert.Contains(output, "保守张(保守型): 买入", "保守型投资者应该在8%涨幅时买入")

	// 小幅下跌3%（应该仅保守型卖出）
	output = captureOutput(func() {
		market.UpdateStockPrice("GOOGL", 1110.456, "谷歌股票小跌3%", 0.1)
	})
	assert.Contains(output, "保守张(保守型): 卖出", "保守型投资者应该在3%跌幅时卖出")
	assert.NotContains(output, "稳健王(稳健型): 卖出", "稳健型投资者不应该在3%跌幅时卖出")
	assert.NotContains(output, "激进李(激进型): 卖出", "激进型投资者不应该在3%跌幅时卖出")
}

// TestMarketAnalyst 测试市场分析师的反应
func TestMarketAnalyst(t *testing.T) {
	assert := assert.New(t)
	market := NewStockMarket()
	analyst := NewMarketAnalyst("anl1", "李四", "某证券公司")
	market.Register(analyst)

	// 测试不同程度的价格变动对分析师的影响

	// 1. 小幅上涨1%
	output := captureOutput(func() {
		market.UpdateStockPrice("MSFT", 100.0, "微软股票初始价格", 0.1)
		market.UpdateStockPrice("MSFT", 101.0, "微软股票小涨1%", 0.1)
	})
	assert.Contains(output, "市场波动不大", "分析师对1%%涨幅的分析不符合预期")

	// 2. 中等上涨3%
	output = captureOutput(func() {
		market.UpdateStockPrice("MSFT", 104.03, "微软股票上涨3%", 0.1)
	})
	assert.Contains(output, "短期上升趋势", "分析师对3%%涨幅的分析不符合预期")

	// 3. 大幅上涨6%
	output = captureOutput(func() {
		market.UpdateStockPrice("MSFT", 110.2718, "微软股票大涨6%", 0.1)
	})
	assert.Contains(output, "市场过热", "分析师对6%%涨幅的分析不符合预期")

	// 4. 中等下跌3%
	output = captureOutput(func() {
		market.UpdateStockPrice("MSFT", 106.9636, "微软股票下跌3%", 0.1)
	})
	assert.Contains(output, "短期下跌趋势", "分析师对3%%跌幅的分析不符合预期")

	// 5. 大幅下跌7%
	output = captureOutput(func() {
		market.UpdateStockPrice("MSFT", 99.4761, "微软股票大跌7%", 0.1)
	})
	assert.Contains(output, "市场恐慌", "分析师对7%%跌幅的分析不符合预期")
}

// TestAsyncNotify 测试异步通知功能
func TestAsyncNotify(t *testing.T) {
	assert := assert.New(t)
	market := NewStockMarket()

	// 创建一个会延迟处理通知的观察者
	var wg sync.WaitGroup
	processTimes := make([]time.Time, 0, 3)
	mutex := sync.Mutex{}

	// 注册观察者 - 直接使用testObserver类型
	market.Register(&testObserver{
		id: "slow1",
		updateFn: func(event StockEvent, message string) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			mutex.Lock()
			processTimes = append(processTimes, time.Now())
			mutex.Unlock()
		},
	})

	market.Register(&testObserver{
		id: "slow2",
		updateFn: func(event StockEvent, message string) {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			mutex.Lock()
			processTimes = append(processTimes, time.Now())
			mutex.Unlock()
		},
	})

	market.Register(&testObserver{
		id: "slow3",
		updateFn: func(event StockEvent, message string) {
			defer wg.Done()
			time.Sleep(30 * time.Millisecond)
			mutex.Lock()
			processTimes = append(processTimes, time.Now())
			mutex.Unlock()
		},
	})

	// 测试异步通知
	event := StockEvent{
		Symbol:    "FB",
		Price:     300.0,
		PrevPrice: 290.0,
		Timestamp: time.Now(),
	}

	wg.Add(3) // 三个观察者
	start := time.Now()
	market.NotifyAsync(event, "Facebook股票更新")

	// 等待所有观察者处理完成
	wg.Wait()
	totalTime := time.Since(start)

	// 对于异步通知，总时间应该接近最长的单个观察者处理时间
	assert.Less(totalTime, 150*time.Millisecond, "异步通知总时间过长")

	// 确认所有观察者都收到了通知
	assert.Equal(3, len(processTimes), "预期有3个观察者处理完成通知")
}

// 用于测试异步通知的自定义观察者
type testObserver struct {
	id       string
	updateFn func(StockEvent, string)
}

func (o *testObserver) Update(event StockEvent, message string) {
	if o.updateFn != nil {
		o.updateFn(event, message)
	}
}

func (o *testObserver) GetID() string {
	return o.id
}

// TestTransactionQuantity 测试投资者的交易数量计算
func TestTransactionQuantity(t *testing.T) {
	assert := assert.New(t)
	conservativeInvestor := NewInvestor("con1", "保守张", Conservative)
	moderateInvestor := NewInvestor("mod1", "稳健王", Moderate)
	aggressiveInvestor := NewInvestor("agg1", "激进李", Aggressive)

	// 创建一个价格变动10%的事件
	event := StockEvent{
		Symbol:    "AMZN",
		Price:     3300.0,
		PrevPrice: 3000.0,
		Timestamp: time.Now(),
	}

	// 测试不同投资者的基础交易量
	conQuantity := conservativeInvestor.decideQuantity(event, true)
	modQuantity := moderateInvestor.decideQuantity(event, true)
	aggQuantity := aggressiveInvestor.decideQuantity(event, true)

	// 激进型投资者应该买入更多
	assert.Less(conQuantity, modQuantity, "保守型投资者交易量应小于稳健型")
	assert.Less(modQuantity, aggQuantity, "稳健型投资者交易量应小于激进型")

	// 测试价格变动对交易量的影响
	smallEvent := StockEvent{
		Symbol:    "AMZN",
		Price:     3030.0,
		PrevPrice: 3000.0,
		Timestamp: time.Now(),
	}

	smallQuantity := aggressiveInvestor.decideQuantity(smallEvent, true)
	largeQuantity := aggressiveInvestor.decideQuantity(event, true)

	// 价格变动越大，交易量应该越大
	assert.Less(smallQuantity, largeQuantity, "大幅变动的交易量应大于小幅变动")
}

// 集成测试：模拟股票市场场景
func TestStockMarketScenario(t *testing.T) {
	assert := assert.New(t)
	market := NewStockMarket()

	// 创建不同类型的观察者
	investor1 := NewInvestor("inv1", "保守张", Conservative)
	investor2 := NewInvestor("inv2", "稳健王", Moderate)
	investor3 := NewInvestor("inv3", "激进李", Aggressive)
	analyst := NewMarketAnalyst("anl1", "专家赵", "大摩投资")

	// 注册观察者
	market.Register(investor1)
	market.Register(investor2)
	market.Register(investor3)
	market.Register(analyst)

	// 验证观察者数量
	assert.Equal(4, market.CountObservers(), "注册后应有4个观察者")

	output := captureOutput(func() {
		// 模拟股票价格波动
		fmt.Println("\n=== 股市波动模拟开始 ===")

		// 初始价格
		market.UpdateStockPrice("TSLA", 700.0, "特斯拉股票初始价格", 0.1)

		// 适度上涨
		market.UpdateStockPrice("TSLA", 728.0, "特斯拉发布新产品", 0.1)

		// 大幅上涨
		market.UpdateStockPrice("TSLA", 800.0, "特斯拉季度盈利超预期", 0.1)

		// 小幅回调
		market.UpdateStockPrice("TSLA", 780.0, "获利回吐", 0.1)

		// 重大利空
		market.UpdateStockPrice("TSLA", 700.0, "特斯拉汽车发生安全事故", 0.1)

		fmt.Println("\n=== 股市波动模拟结束 ===")
	})

	// 检查输出中是否包含所有观察者的反应
	assert.Contains(output, "保守张", "输出中应包含保守型投资者的反应")
	assert.Contains(output, "稳健王", "输出中应包含稳健型投资者的反应")
	assert.Contains(output, "激进李", "输出中应包含激进型投资者的反应")
	assert.Contains(output, "专家赵", "输出中应包含市场分析师的反应")

	// 确认所有类型的通知都发出了
	assert.Contains(output, "买入", "输出中应包含买入操作")
	assert.Contains(output, "卖出", "输出中应包含卖出操作")
	assert.Contains(output, "观望", "输出中应包含观望操作")

	// 确认价格更新成功
	price, exists := market.GetStockPrice("TSLA")
	assert.True(exists, "应该存在股票TSLA")
	assert.Equal(700.0, price, "最终股票价格应为700.0")
}

// 基准测试：测试通知性能
func BenchmarkNotify(b *testing.B) {
	market := NewStockMarket()
	numObservers := 100

	// 创建多个观察者
	for i := 0; i < numObservers; i++ {
		market.Register(NewInvestor(fmt.Sprintf("inv%d", i),
			fmt.Sprintf("投资者%d", i),
			InvestorType(i%3)))
	}

	event := StockEvent{
		Symbol:    "INDEX",
		Price:     1000.0,
		PrevPrice: 980.0,
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		market.Notify(event, "指数更新")
	}
}

// 基准测试：测试异步通知性能
func BenchmarkNotifyAsync(b *testing.B) {
	market := NewStockMarket()
	numObservers := 100

	// 创建多个观察者
	for i := 0; i < numObservers; i++ {
		market.Register(NewInvestor(fmt.Sprintf("inv%d", i),
			fmt.Sprintf("投资者%d", i),
			InvestorType(i%3)))
	}

	event := StockEvent{
		Symbol:    "INDEX",
		Price:     1000.0,
		PrevPrice: 980.0,
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			market.NotifyAsync(event, "指数异步更新")
		}()
		wg.Wait()
	}
}
