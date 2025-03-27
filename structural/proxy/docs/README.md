# 代理设计模式 (Proxy Pattern)

## 介绍

代理模式是一种结构型设计模式，它允许你提供一个替代品或占位符来控制对原始对象的访问。代理对象控制着对原始对象的访问，并允许在请求到达原始对象前后执行一些处理。

代理模式的核心思想是：**通过引入一个新的代理对象来间接访问目标对象，从而实现在不改变目标对象的前提下，增强其功能或控制其访问**。

![代理模式结构](https://refactoringguru.cn/images/patterns/diagrams/proxy/structure-indexed.png)

## 类图结构

```
┌───────────┐     ┌───────────┐
│           │     │           │
│  Client   │─────▶  IBuyCar  │
│           │     │ Interface │
└───────────┘     └─────┬─────┘
                        │
                        │implements
        ┌───────────────┼───────────────────┐
        │               │                   │
┌───────▼──────┐  ┌─────▼─────┐    ┌────────▼─────────┐
│              │  │           │    │                  │
│  RealBuyer   │◀─┤  Proxies  │────┤ Other Components │
│(RealSubject) │  │           │    │                  │
└──────────────┘  └───────────┘    └──────────────────┘
                        │
                        │
        ┌───────────────┴───────────────────────────────────────┐
        │               │               │           │           │
┌───────▼─────┐  ┌──────▼─────┐  ┌──────▼────┐ ┌────▼─────┐ ┌───▼────────┐
│             │  │            │  │           │ │          │ │            │
│ FourSProxy  │  │ VirtualProxy│  │Protection │ │ Logging │ │  Cached    │
│(Basic Proxy)│  │(Lazy Init)  │  │  Proxy   │ │  Proxy  │ │  Proxy     │
└─────────────┘  └────────────┘  └───────────┘ └──────────┘ └────────────┘
```

## 代理模式的类型

代理模式根据其用途和功能可分为多种类型:

### 1. 基本代理 (Basic Proxy)

在访问对象前后添加额外的处理逻辑，增强原有对象的功能。

### 2. 虚拟代理 (Virtual Proxy)

延迟创建开销大的对象，直到真正需要使用时才创建，从而优化资源使用。

### 3. 保护代理 (Protection Proxy)

控制对原始对象的访问权限，根据客户端权限决定是否允许访问。

### 4. 日志代理 (Logging Proxy)

记录对原始对象的访问日志，用于调试、审计等目的。

### 5. 缓存代理 (Cache Proxy)

存储耗时操作的结果，在需要相同结果时直接返回缓存，避免重复执行。

### 6. 远程代理 (Remote Proxy)

代表位于远程服务器上的对象，处理网络通信细节。

### 7. 智能引用代理 (Smart Reference Proxy)

在访问对象时进行额外的处理，如引用计数、加载时验证等。

## 代码实现

本示例通过汽车购买场景展示了五种不同类型的代理模式实现。

### 接口设计

```go
// IBuyCar 定义了被代理对象和代理对象共同实现的接口
type IBuyCar interface {
    BuyCar() error
    GetCarInfo() string
}
```

### 真实主体 (Real Subject)

```go
// RealBuyer 是实际买车的人（被代理对象）
type RealBuyer struct {
    Name  string
    Money float64
}

// BuyCar 实现了IBuyCar接口的方法
func (r *RealBuyer) BuyCar() error {
    if r.Money < 100000 {
        return fmt.Errorf("余额不足，无法购买汽车")
    }
    fmt.Printf("<%s> 成功购买了一辆汽车，花费了 ¥%.2f\n", r.Name, 100000.0)
    r.Money -= 100000
    return nil
}

// GetCarInfo 获取车辆信息
func (r *RealBuyer) GetCarInfo() string {
    return "标准汽车型号XYZ"
}
```

### 1. 基本代理 - 4S店代理

```go
// FourSProxy 是基本代理，提供额外服务
type FourSProxy struct {
    realBuyer IBuyCar
    services  []string
    fee       float64
}

// BuyCar 代理实现的购车方法，添加了额外的服务
func (f *FourSProxy) BuyCar() error {
    fmt.Println("=== 通过4S店代理购车开始 ===")
    
    // 代理前的操作
    fmt.Println("1. 从制造商订购汽车到4S店")
    fmt.Println("2. 准备购车文件")
    
    // 调用实际对象的方法
    if err := f.realBuyer.BuyCar(); err != nil {
        fmt.Printf("购车失败: %s\n", err)
        return err
    }
    
    // 代理后的增强操作
    fmt.Println("提供额外服务:")
    for i, service := range f.services {
        fmt.Printf("  %d. %s\n", i+1, service)
    }
    
    fmt.Printf("收取服务费: ¥%.2f\n", f.fee)
    fmt.Println("=== 通过4S店代理购车完成 ===")
    return nil
}
```

### 2. 虚拟代理 - 延迟初始化

```go
// VirtualBuyerProxy 虚拟代理 - 延迟创建被代理对象
type VirtualBuyerProxy struct {
    name      string
    money     float64
    realBuyer *RealBuyer  // 初始为nil
}

// BuyCar 虚拟代理实现，延迟创建被代理对象
func (v *VirtualBuyerProxy) BuyCar() error {
    fmt.Println("=== 通过虚拟代理购车开始 ===")
    fmt.Println("准备创建实际购买者...")
    
    // 延迟初始化 - 仅在首次调用时创建实际对象
    if v.realBuyer == nil {
        fmt.Println("首次调用，创建实际购买者")
        v.realBuyer = NewRealBuyer(v.name, v.money)
    } else {
        fmt.Println("复用已有的实际购买者")
    }
    
    err := v.realBuyer.BuyCar()
    fmt.Println("=== 通过虚拟代理购车结束 ===")
    return err
}
```

### 3. 保护代理 - 权限控制

```go
// ProtectionProxy 保护代理 - 控制对资源的访问权限
type ProtectionProxy struct {
    realBuyer IBuyCar
    isVIP     bool
}

// BuyCar 保护代理实现，加入权限控制
func (p *ProtectionProxy) BuyCar() error {
    fmt.Println("=== 通过保护代理购车开始 ===")
    
    // 权限检查
    if !p.isVIP {
        fmt.Println("权限不足: 仅VIP客户可以通过此渠道购车")
        return fmt.Errorf("权限不足: 需要VIP权限")
    }
    
    fmt.Println("VIP客户，权限验证通过")
    err := p.realBuyer.BuyCar()
    
    if err == nil {
        fmt.Println("VIP客户专享折扣已应用")
    }
    
    fmt.Println("=== 通过保护代理购车结束 ===")
    return err
}
```

### 4. 日志代理 - 记录操作

```go
// LoggingProxy 日志代理 - 记录操作日志
type LoggingProxy struct {
    realBuyer IBuyCar
}

// BuyCar 日志代理实现，添加日志记录
func (l *LoggingProxy) BuyCar() error {
    fmt.Println("=== 日志记录: 购车操作开始 ===")
    startTime := time.Now()
    
    fmt.Printf("[%s] 购车请求已接收\n", startTime.Format("2006-01-02 15:04:05"))
    
    err := l.realBuyer.BuyCar()
    
    endTime := time.Now()
    duration := endTime.Sub(startTime)
    
    if err != nil {
        fmt.Printf("[%s] 购车失败: %s\n", endTime.Format("2006-01-02 15:04:05"), err)
    } else {
        fmt.Printf("[%s] 购车成功\n", endTime.Format("2006-01-02 15:04:05"))
    }
    
    fmt.Printf("操作耗时: %v\n", duration)
    fmt.Println("=== 日志记录: 购车操作结束 ===")
    return err
}
```

### 5. 缓存代理 - 结果缓存

```go
// CachedBuyerProxy 缓存代理 - 缓存重复请求的结果
type CachedBuyerProxy struct {
    realBuyer IBuyCar
    carInfo   string
    cached    bool
}

// GetCarInfo 获取车辆信息，支持缓存
func (c *CachedBuyerProxy) GetCarInfo() string {
    if c.cached {
        fmt.Println("从缓存获取车辆信息")
        return c.carInfo + " (缓存)"
    }
    
    fmt.Println("首次获取车辆信息，将结果缓存")
    c.carInfo = c.realBuyer.GetCarInfo()
    c.cached = true
    return c.carInfo
}
```

## 使用示例

### 基本代理示例

```go
// 创建实际购买者
buyer := NewRealBuyer("张三", 150000)

// 创建4S店代理
proxy := NewFourSProxy(buyer)

// 通过代理购车
proxy.BuyCar()

// 输出:
// === 通过4S店代理购车开始 ===
// 1. 从制造商订购汽车到4S店
// 2. 准备购车文件
// <张三> 成功购买了一辆汽车，花费了 ¥100000.00
// 提供额外服务:
//   1. 上牌服务
//   2. 汽车注册
//   3. 保险办理
// 收取服务费: ¥5000.00
// === 通过4S店代理购车完成 ===
```

### 虚拟代理示例

```go
// 创建虚拟代理，此时不会创建实际对象
proxy := NewVirtualBuyerProxy("李四", 200000)

// 第一次调用时创建实际对象
proxy.BuyCar()

// 第二次调用时复用已有对象
proxy.BuyCar()
```

### 保护代理示例

```go
// 创建实际购买者
buyer := NewRealBuyer("王五", 300000)

// 创建VIP保护代理
vipProxy := NewProtectionProxy(buyer, true)
vipProxy.BuyCar()  // 成功

// 创建普通保护代理
normalProxy := NewProtectionProxy(buyer, false)
normalProxy.BuyCar()  // 失败，权限不足
```

### 日志代理示例

```go
buyer := NewRealBuyer("日志测试", 150000)
proxy := NewLoggingProxy(buyer)
proxy.BuyCar()
```

### 缓存代理示例

```go
buyer := NewRealBuyer("赵六", 150000)
proxy := NewCachedBuyerProxy(buyer)

// 获取车辆信息 - 第一次，将会缓存结果
info1 := proxy.GetCarInfo()

// 获取车辆信息 - 第二次，将使用缓存
info2 := proxy.GetCarInfo()
```

### 组合多个代理

代理模式的一个强大特性是可以组合多个代理，形成代理链：

```go
// 创建代理链：缓存代理 -> 保护代理 -> 日志代理 -> 4S店代理 -> 实际购买者
buyer := NewRealBuyer("复合代理客户", 200000)
fourSProxy := NewFourSProxy(buyer)
loggingProxy := NewLoggingProxy(fourSProxy)
protectionProxy := NewProtectionProxy(loggingProxy, true)
cachedProxy := NewCachedBuyerProxy(protectionProxy)

// 通过代理链购车
cachedProxy.BuyCar()
```

## 代理模式的优点

1. **单一职责原则**：代理类可以处理被代理对象的功能增强，使主体类专注于自身业务
2. **开闭原则**：无需修改原始类即可扩展其功能
3. **组合灵活**：多个代理可以组合使用，形成功能链
4. **远程处理**：可以隐藏远程对象的复杂性
5. **延迟加载**：使用虚拟代理延迟创建开销大的对象
6. **访问控制**：通过保护代理控制对敏感对象的访问
7. **缓存优化**：使用缓存代理提高系统性能

## 代理模式的缺点

1. **增加系统复杂度**：引入额外的代理类
2. **可能导致请求处理变慢**：增加额外的间接层
3. **接口一致性要求**：代理类和实际类必须实现相同的接口

## 适用场景

1. **延迟初始化** (虚拟代理)：优化系统资源，避免创建开销大的对象
2. **访问控制** (保护代理)：控制对敏感对象的访问权限
3. **本地执行远程服务** (远程代理)：隐藏远程对象的复杂性
4. **记录日志** (日志代理)：记录方法调用的日志信息
5. **缓存请求结果** (缓存代理)：存储计算结果并在请求相同结果时重用
6. **对现有对象的包装** (基本代理)：在不修改原始代码的情况下添加功能

## 与其他模式的关系

1. **装饰器模式**：两者结构类似，但意图不同。代理控制对对象的访问，而装饰器为对象添加新职责
2. **适配器模式**：适配器提供不同的接口，而代理提供相同的接口
3. **外观模式**：外观提供简化的接口，代理保持相同的接口但控制访问

## 总结

代理模式是一种强大的结构型设计模式，通过引入间接层控制对目标对象的访问。它可以用于多种场景，如延迟初始化、访问控制、日志记录和缓存等。通过正确使用代理模式，可以增强系统的灵活性、可维护性和性能。

在Go语言中，接口和组合机制使得代理模式的实现变得简单而灵活。本示例展示了五种常见的代理类型，以及如何组合多个代理形成代理链，以满足复杂的业务需求。