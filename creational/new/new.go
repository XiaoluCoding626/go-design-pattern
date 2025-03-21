// Package new 展示 Go 语言中的 New 模式（构造函数模式）
package new

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Product 表示一个商品
// 注意：某些字段是小写私有的，强制用户通过构造函数创建实例
type Product struct {
	name      string    // 商品名称（私有）
	price     float64   // 商品价格（私有）
	ID        string    // 商品ID（公开）
	category  string    // 商品类别（私有）
	CreatedAt time.Time // 创建时间（公开）
	stock     int       // 库存数量（私有）
	discount  float64   // 折扣（私有）
}

// NewProduct 创建并返回一个基本的商品实例
// 这是主要的构造函数，要求提供必要的名称和价格参数
func NewProduct(name string, price float64) (*Product, error) {
	// 验证参数
	if name == "" {
		return nil, errors.New("商品名称不能为空")
	}
	if price <= 0 {
		return nil, errors.New("商品价格必须大于零")
	}

	// 创建并初始化商品
	p := &Product{
		name:      name,
		price:     price,
		ID:        generateID(name),
		CreatedAt: time.Now(),
		stock:     0,     // 默认库存为0
		discount:  1.0,   // 默认无折扣
		category:  "未分类", // 默认分类
	}

	return p, nil
}

// NewDiscountedProduct 创建带有折扣的商品
// 这是一个特殊用途的构造函数
func NewDiscountedProduct(name string, price float64, discountPercent float64) (*Product, error) {
	// 验证折扣参数
	if discountPercent < 0 || discountPercent > 100 {
		return nil, errors.New("折扣百分比必须在0到100之间")
	}

	// 先创建基本商品
	product, err := NewProduct(name, price)
	if err != nil {
		return nil, err
	}

	// 设置折扣（折扣以小数表示，例如：20%折扣 = 0.8）
	product.discount = (100 - discountPercent) / 100
	return product, nil
}

// NewProductInStock 创建一个有初始库存的商品
func NewProductInStock(name string, price float64, initialStock int) (*Product, error) {
	// 验证库存参数
	if initialStock < 0 {
		return nil, errors.New("初始库存不能为负数")
	}

	// 先创建基本商品
	product, err := NewProduct(name, price)
	if err != nil {
		return nil, err
	}

	// 设置初始库存
	product.stock = initialStock
	return product, nil
}

// NewProductComplete 创建一个完整配置的商品（所有参数）
func NewProductComplete(name string, price float64, category string,
	initialStock int, discountPercent float64) (*Product, error) {

	// 验证类别
	if category == "" {
		return nil, errors.New("商品类别不能为空")
	}

	// 验证折扣参数
	if discountPercent < 0 || discountPercent > 100 {
		return nil, errors.New("折扣百分比必须在0到100之间")
	}

	// 验证库存参数
	if initialStock < 0 {
		return nil, errors.New("初始库存不能为负数")
	}

	// 创建基本商品
	product, err := NewProduct(name, price)
	if err != nil {
		return nil, err
	}

	// 设置所有附加参数
	product.category = category
	product.stock = initialStock
	product.discount = (100 - discountPercent) / 100

	return product, nil
}

// WithCategory 是一个链式方法，用于设置商品类别
// 演示了 Functional Options 模式与 New 模式的结合
func (p *Product) WithCategory(category string) *Product {
	if category != "" {
		p.category = category
	}
	return p
}

// WithStock 是一个链式方法，用于设置商品库存
func (p *Product) WithStock(stock int) *Product {
	if stock >= 0 {
		p.stock = stock
	}
	return p
}

// WithDiscount 是一个链式方法，用于设置商品折扣
func (p *Product) WithDiscount(discountPercent float64) *Product {
	if discountPercent >= 0 && discountPercent <= 100 {
		p.discount = (100 - discountPercent) / 100
	}
	return p
}

// 获取商品属性的方法

// GetName 返回商品名称
func (p *Product) GetName() string {
	return p.name
}

// GetPrice 返回商品当前价格（考虑折扣）
func (p *Product) GetPrice() float64 {
	return p.price * p.discount
}

// GetOriginalPrice 返回商品原价
func (p *Product) GetOriginalPrice() float64 {
	return p.price
}

// GetCategory 返回商品类别
func (p *Product) GetCategory() string {
	return p.category
}

// GetStock 返回当前库存
func (p *Product) GetStock() int {
	return p.stock
}

// GetDiscount 返回折扣百分比
func (p *Product) GetDiscount() float64 {
	return (1 - p.discount) * 100
}

// 商品状态修改方法

// AddStock 增加库存数量
func (p *Product) AddStock(amount int) error {
	if amount < 0 {
		return errors.New("增加的库存数量不能为负")
	}
	p.stock += amount
	return nil
}

// ReduceStock 减少库存数量
func (p *Product) ReduceStock(amount int) error {
	if amount < 0 {
		return errors.New("减少的库存数量不能为负")
	}
	if p.stock < amount {
		return errors.New("库存不足")
	}
	p.stock -= amount
	return nil
}

// ApplyDiscount 应用折扣到商品
func (p *Product) ApplyDiscount(discountPercent float64) error {
	if discountPercent < 0 || discountPercent > 100 {
		return errors.New("折扣百分比必须在0到100之间")
	}
	p.discount = (100 - discountPercent) / 100
	return nil
}

// String 实现 Stringer 接口，提供友好的字符串表示
func (p *Product) String() string {
	discountInfo := ""
	if p.discount < 1.0 {
		discountInfo = fmt.Sprintf(" (折扣: %.1f%%，折后价: ¥%.2f)",
			(1-p.discount)*100, p.price*p.discount)
	}

	return fmt.Sprintf("商品: %s (ID: %s)\n"+
		"类别: %s\n"+
		"价格: ¥%.2f%s\n"+
		"库存: %d\n"+
		"创建时间: %s",
		p.name, p.ID,
		p.category,
		p.price, discountInfo,
		p.stock,
		p.CreatedAt.Format("2006-01-02 15:04:05"))
}

// Clone 创建并返回当前商品的一个深拷贝
// 展示了 New 模式与原型模式的结合
func (p *Product) Clone() *Product {
	return &Product{
		name:      p.name,
		price:     p.price,
		ID:        generateID(p.name), // 生成新ID
		category:  p.category,
		CreatedAt: time.Now(), // 创建时间更新
		stock:     p.stock,
		discount:  p.discount,
	}
}

// 辅助函数

// generateID 基于名称、当前时间和随机数生成一个唯一ID
func generateID(name string) string {
	timestamp := time.Now().UnixNano() / 1000000 // 毫秒时间戳
	// 添加6位随机数以确保唯一性，即使在同一毫秒内生成多个ID
	randomPart := rand.Intn(1000000)

	if len(name) > 3 {
		return fmt.Sprintf("%s-%d-%06d", name[:3], timestamp, randomPart)
	}
	return fmt.Sprintf("%s-%d-%06d", name, timestamp, randomPart)
}

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}
