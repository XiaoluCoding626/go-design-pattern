# New 模式（构造函数模式）

## 简介

New 模式是 Go 语言中常用的一种创建对象的设计模式，也被称为构造函数模式（Constructor Pattern）。这种模式遵循 Go 语言的惯例，通过提供以 `New` 开头的工厂函数来创建和初始化结构体实例，而不是直接由调用者创建。

虽然 New 模式不是 GoF（Gang of Four）设计模式中的一员，但它是 Go 语言特有的一种实践模式，在标准库和优质第三方库中被广泛应用。

## 设计思想

New 模式的核心思想是：

1. **封装复杂的初始化逻辑**：隐藏内部细节，提供简洁的 API。
2. **验证输入参数**：确保创建的对象始终处于有效状态。
3. **提供合理的默认值**：简化客户端代码，减少必需参数。
4. **返回接口而非具体类型**：增强代码的灵活性（当适用时）。
5. **处理错误**：返回错误而不是创建无效对象。

## 代码实现

### 基本结构

```go
type Product struct {
    name     string    // 私有字段，强制使用构造函数
    price    float64
    // 其他字段...
}

// NewProduct 创建并返回一个新的 Product 实例
func NewProduct(name string, price float64) (*Product, error) {
    // 验证参数
    if name == "" {
        return nil, errors.New("商品名称不能为空")
    }
    if price <= 0 {
        return nil, errors.New("商品价格必须大于零")
    }
    
    // 创建实例并设置默认值
    return &Product{
        name:  name,
        price: price,
        // 设置其他默认值...
    }, nil
}
```

### 多种构造函数

为不同的创建场景提供专门的构造函数是一种常见实践：

```go
// NewDiscountedProduct 创建带有折扣的商品
func NewDiscountedProduct(name string, price float64, discountPercent float64) (*Product, error) {
    // 参数验证...
    
    // 创建基础产品
    product, err := NewProduct(name, price)
    if err != nil {
        return nil, err
    }
    
    // 设置折扣
    product.discount = (100 - discountPercent) / 100
    return product, nil
}
```

### 链式方法设置可选参数

结合函数选项模式（Functional Options Pattern）处理可选配置：

```go
// WithCategory 设置商品类别
func (p *Product) WithCategory(category string) *Product {
    if category != "" {
        p.category = category
    }
    return p
}

// 使用示例
product, _ := NewProduct("手机", 1999.99)
product.WithCategory("电子产品").WithStock(100)
```

## 何时使用 New 模式

New 模式特别适合以下场景：

1. **结构体需要验证**：确保创建的对象始终有效。
2. **初始化逻辑复杂**：涉及计算、默认值设置或资源分配。
3. **私有字段保护**：强制通过构造函数创建对象。
4. **多步骤初始化**：需要多个步骤才能完成对象初始化。
5. **不同的创建变体**：根据不同需求创建同一类型的不同变体。

## 最佳实践

### 1. 命名约定

- 使用 `New[Type]` 作为主构造函数的名称
- 对于特殊变体，使用 `New[Type][Variant]` 格式

```go
NewProduct(...)        // 主要构造函数
NewProductInStock(...) // 特定变体
```

### 2. 错误处理

始终返回错误，而不是恐慌（panic）：

```go
func NewProduct(name string, price float64) (*Product, error) {
    if price < 0 {
        return nil, errors.New("价格不能为负")
    }
    // ...
}
```

### 3. 使用指针接收者方法

对于链式配置方法，使用指针接收者以便修改对象：

```go
func (p *Product) WithDiscount(discount float64) *Product {
    // 修改 p 的状态
    return p  // 返回接收者以支持链式调用
}
```

### 4. 默认值

提供合理的默认值，减少客户端代码的负担：

```go
func NewProduct(name string, price float64) (*Product, error) {
    return &Product{
        name:      name,
        price:     price,
        createdAt: time.Now(),
        category:  "未分类", // 默认分类
        discount:  1.0,     // 默认无折扣
    }, nil
}
```

### 5. 隐藏实现细节

将不需要导出的字段设为小写，只暴露必要的 API：

```go
type Product struct {
    name     string  // 私有，不导出
    price    float64 // 私有，不导出
    ID       string  // 公开，导出
}

// 提供访问方法
func (p *Product) GetName() string {
    return p.name
}
```

## 与其他模式的比较

### New 模式 vs 工厂方法模式

- **New 模式**：专注于单一类型的对象创建，通常返回具体类型
- **工厂方法**：创建一族相关对象，通常返回接口类型

### New 模式 vs 构建器模式

- **New 模式**：适用于简单到中等复杂度的对象创建
- **构建器模式**：更适合创建非常复杂、有多个可选配置的对象

## 代码示例与解析

以下是一个完整的示例，展示了 New 模式的实现：

```go
// Product 表示一个商品
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

// 使用链式方法设置可选参数
func (p *Product) WithCategory(category string) *Product {
    if category != "" {
        p.category = category
    }
    return p
}

// 使用示例
product, err := NewProduct("笔记本电脑", 6999.99)
if err != nil {
    log.Fatal(err)
}
product.WithCategory("电子产品").WithStock(50).WithDiscount(10)
```

## 结语

New 模式是 Go 语言中创建和初始化对象的优雅解决方案。通过遵循这种模式，你可以编写出更加健壮、可维护和易于使用的代码。它虽然简单，但在正确应用时非常强大，能够有效地封装复杂的初始化逻辑，提供参数验证，并确保对象始终处于有效状态。

在实际项目中，New 模式往往会与其他设计模式（如工厂方法、构建器或原型模式）结合使用，以满足更复杂的对象创建需求。