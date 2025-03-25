package observer

import (
	"fmt"
	"sync"
	"time"
)

// StockEvent 表示股票事件数据
type StockEvent struct {
	Symbol    string    // 股票代码
	Price     float64   // 当前价格
	PrevPrice float64   // 前一个价格
	Timestamp time.Time // 时间戳
}

// ChangePercent 返回价格变动百分比
func (e StockEvent) ChangePercent() float64 {
	if e.PrevPrice == 0 {
		return 0
	}
	return (e.Price - e.PrevPrice) / e.PrevPrice * 100
}

// IsUp 返回股票是否上涨
func (e StockEvent) IsUp() bool {
	return e.Price > e.PrevPrice
}

// IsPriceChange 返回价格变动幅度是否超过阈值
func (e StockEvent) IsPriceChange(threshold float64) bool {
	change := e.ChangePercent()
	return change >= threshold || change <= -threshold
}

// String 格式化打印事件信息
func (e StockEvent) String() string {
	var direction string
	if e.IsUp() {
		direction = "上涨"
	} else if e.Price < e.PrevPrice {
		direction = "下跌"
	} else {
		direction = "持平"
	}

	percent := e.ChangePercent()

	return fmt.Sprintf("%s: %.2f -> %.2f (%.2f%% %s)",
		e.Symbol, e.PrevPrice, e.Price, percent, direction)
}

// Subject 定义了主题接口
type Subject interface {
	Register(observer Observer)                   // 注册观察者
	Deregister(observer Observer)                 // 注销观察者
	Notify(event StockEvent, message string)      // 通知所有观察者
	NotifyAsync(event StockEvent, message string) // 异步通知所有观察者
	HasObserver(observer Observer) bool           // 检查观察者是否已注册
	CountObservers() int                          // 获取观察者数量
}

// Observer 定义了观察者接口
type Observer interface {
	Update(event StockEvent, message string) // 接收更新
	GetID() string                           // 获取观察者标识
}

// StockMarket 具体主题，实现了 Subject 接口
type StockMarket struct {
	observers []Observer         // 观察者列表
	stocks    map[string]float64 // 股票价格映射表
	mutex     sync.RWMutex       // 保证线程安全
}

// NewStockMarket 创建一个新的股票市场
func NewStockMarket() *StockMarket {
	return &StockMarket{
		observers: make([]Observer, 0),
		stocks:    make(map[string]float64),
	}
}

// Register 实现注册观察者
func (s *StockMarket) Register(observer Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查是否已注册
	if s.HasObserverUnsafe(observer) {
		fmt.Printf("观察者 %s 已经注册\n", observer.GetID())
		return
	}
	s.observers = append(s.observers, observer)
	fmt.Printf("观察者 %s 已注册到股票市场\n", observer.GetID())
}

// Deregister 实现注销观察者
func (s *StockMarket) Deregister(observer Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, obs := range s.observers {
		if obs.GetID() == observer.GetID() {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			fmt.Printf("观察者 %s 已从股票市场注销\n", observer.GetID())
			return
		}
	}
}

// HasObserverUnsafe 检查观察者是否已注册（非线程安全，只在加锁后使用）
func (s *StockMarket) HasObserverUnsafe(observer Observer) bool {
	for _, obs := range s.observers {
		if obs.GetID() == observer.GetID() {
			return true
		}
	}
	return false
}

// HasObserver 线程安全地检查观察者是否已注册
func (s *StockMarket) HasObserver(observer Observer) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.HasObserverUnsafe(observer)
}

// CountObservers 获取观察者数量
func (s *StockMarket) CountObservers() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.observers)
}

// Notify 通知所有观察者（同步）
func (s *StockMarket) Notify(event StockEvent, message string) {
	s.mutex.RLock()
	observers := make([]Observer, len(s.observers))
	copy(observers, s.observers)
	s.mutex.RUnlock()

	fmt.Printf("\n【市场公告】%s\n", message)
	fmt.Printf("股票行情: %s\n", event.String())

	for _, observer := range observers {
		observer.Update(event, message)
	}
}

// NotifyAsync 异步通知所有观察者
func (s *StockMarket) NotifyAsync(event StockEvent, message string) {
	s.mutex.RLock()
	observers := make([]Observer, len(s.observers))
	copy(observers, s.observers)
	s.mutex.RUnlock()

	fmt.Printf("\n【市场公告】%s\n", message)
	fmt.Printf("股票行情: %s\n", event.String())

	var wg sync.WaitGroup
	for _, observer := range observers {
		wg.Add(1)
		go func(o Observer) {
			defer wg.Done()
			o.Update(event, message)
		}(observer)
	}

	// 可以选择等待所有通知完成或不等待
	// wg.Wait()
}

// UpdateStockPrice 更新股票价格并通知观察者
func (s *StockMarket) UpdateStockPrice(symbol string, newPrice float64, message string, notifyThreshold float64) {
	s.mutex.Lock()
	prevPrice, exists := s.stocks[symbol]
	if !exists {
		prevPrice = 0
	}
	s.stocks[symbol] = newPrice
	s.mutex.Unlock()

	event := StockEvent{
		Symbol:    symbol,
		Price:     newPrice,
		PrevPrice: prevPrice,
		Timestamp: time.Now(),
	}

	// 只有价格变动超过阈值时才通知
	if !exists || event.IsPriceChange(notifyThreshold) {
		s.Notify(event, message)
	}
}

// GetStockPrice 获取股票价格
func (s *StockMarket) GetStockPrice(symbol string) (float64, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	price, exists := s.stocks[symbol]
	return price, exists
}

// 观察者类型

// InvestorType 表示投资者类型
type InvestorType int

const (
	Conservative InvestorType = iota // 保守型
	Moderate                         // 稳健型
	Aggressive                       // 激进型
)

// Investor 实现了 Observer 接口的投资者
type Investor struct {
	id             string         // 投资者ID
	name           string         // 投资者名称
	investorType   InvestorType   // 投资者类型
	buyThreshold   float64        // 买入阈值
	sellThreshold  float64        // 卖出阈值
	currentHolding map[string]int // 当前持股
}

// NewInvestor 创建一个新的投资者
func NewInvestor(id, name string, investorType InvestorType) *Investor {
	var buyThreshold, sellThreshold float64

	switch investorType {
	case Conservative:
		buyThreshold = 5.0   // 股票上涨5%才买入
		sellThreshold = -2.0 // 股票下跌2%就卖出
	case Moderate:
		buyThreshold = 3.0   // 股票上涨3%才买入
		sellThreshold = -5.0 // 股票下跌5%才卖出
	case Aggressive:
		buyThreshold = 0.5    // 股票上涨0.5%就买入
		sellThreshold = -10.0 // 股票下跌10%才卖出
	}

	return &Investor{
		id:             id,
		name:           name,
		investorType:   investorType,
		buyThreshold:   buyThreshold,
		sellThreshold:  sellThreshold,
		currentHolding: make(map[string]int),
	}
}

// Update 实现了 Observer 接口的更新方法
func (i *Investor) Update(event StockEvent, message string) {
	changePercent := event.ChangePercent()

	var action string
	switch {
	case changePercent >= i.buyThreshold:
		// 上涨超过买入阈值，买入
		quantity := i.decideQuantity(event, true)
		i.currentHolding[event.Symbol] += quantity
		action = fmt.Sprintf("买入 %d 股 %s", quantity, event.Symbol)
	case changePercent <= i.sellThreshold:
		// 下跌超过卖出阈值，卖出
		quantity := i.decideQuantity(event, false)
		if quantity > i.currentHolding[event.Symbol] {
			quantity = i.currentHolding[event.Symbol]
		}
		i.currentHolding[event.Symbol] -= quantity
		action = fmt.Sprintf("卖出 %d 股 %s", quantity, event.Symbol)
	default:
		// 不满足交易条件，观望
		action = "观望行情"
	}

	fmt.Printf("%s(%s): %s [持股: %d]\n",
		i.name, i.typeString(), action, i.currentHolding[event.Symbol])
}

// 根据投资者类型和事件决定交易数量
func (i *Investor) decideQuantity(event StockEvent, isBuying bool) int {
	var baseQuantity int

	// 基础交易量根据投资者类型决定
	switch i.investorType {
	case Conservative:
		baseQuantity = 100
	case Moderate:
		baseQuantity = 200
	case Aggressive:
		baseQuantity = 500
	}

	// 根据价格变动幅度调整交易量
	changePercent := event.ChangePercent()
	if changePercent < 0 {
		changePercent = -changePercent
	}

	// 变动越大，交易量越大
	multiplier := 1.0 + (changePercent / 10.0)
	return int(float64(baseQuantity) * multiplier)
}

// GetID 实现 Observer 接口的 GetID 方法
func (i *Investor) GetID() string {
	return i.id
}

// 返回投资者类型的字符串表示
func (i *Investor) typeString() string {
	switch i.investorType {
	case Conservative:
		return "保守型"
	case Moderate:
		return "稳健型"
	case Aggressive:
		return "激进型"
	default:
		return "未知类型"
	}
}

// MarketAnalyst 另一种观察者，市场分析师
type MarketAnalyst struct {
	id      string
	name    string
	company string
}

// NewMarketAnalyst 创建一个新的市场分析师
func NewMarketAnalyst(id, name, company string) *MarketAnalyst {
	return &MarketAnalyst{
		id:      id,
		name:    name,
		company: company,
	}
}

// Update 实现了 Observer 接口的更新方法
func (a *MarketAnalyst) Update(event StockEvent, message string) {
	var analysis string

	// 根据价格变动提供分析
	changePercent := event.ChangePercent()
	switch {
	case changePercent > 5:
		analysis = "市场过热，建议获利了结"
	case changePercent > 2:
		analysis = "短期上升趋势，可逢低买入"
	case changePercent < -5:
		analysis = "市场恐慌，优质股可以逐步建仓"
	case changePercent < -2:
		analysis = "短期下跌趋势，建议观望"
	default:
		analysis = "市场波动不大，维持原有策略"
	}

	fmt.Printf("%s分析师(%s): %s\n", a.name, a.company, analysis)
}

// GetID 实现 Observer 接口的 GetID 方法
func (a *MarketAnalyst) GetID() string {
	return a.id
}
