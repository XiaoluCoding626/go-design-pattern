// Package objectpool 实现了对象池设计模式，用于管理和重用昂贵的对象
package objectpool

import (
	"errors"
	"sync"
	"time"
)

// 定义常见错误
var (
	ErrPoolClosed        = errors.New("object pool is closed")
	ErrPoolTimeout       = errors.New("timeout waiting for object")
	ErrInvalidObject     = errors.New("invalid object returned to pool")
	ErrPoolAtMaxCapacity = errors.New("pool reached max capacity")
)

// Object 表示对象池中的对象接口
// 使用接口允许池可以管理任何类型的对象
type Object interface {
	// Reset 重置对象状态，为重用做准备
	Reset() error

	// Validate 验证对象是否有效/可重用
	Validate() bool

	// ID 返回对象的唯一标识符
	ID() int
}

// ObjectFactory 定义了用于创建新对象的工厂函数类型
type ObjectFactory func() (Object, error)

// PoolConfig 保存对象池的配置选项
type PoolConfig struct {
	// InitialSize 是池初始化时创建的对象数量
	InitialSize int

	// MaxSize 是池可以增长到的最大大小
	MaxSize int

	// MaxIdle 是允许保持空闲状态的最大对象数量
	MaxIdle int

	// Factory 用于创建新对象的工厂函数
	Factory ObjectFactory

	// MinEvictableIdleTime 是对象在被收回前可以空闲的最小时间
	MinEvictableIdleTime time.Duration

	// ValidationInterval 是验证空闲对象的时间间隔
	ValidationInterval time.Duration
}

// DefaultPoolConfig 返回具有合理默认值的池配置
func DefaultPoolConfig(factory ObjectFactory) PoolConfig {
	return PoolConfig{
		InitialSize:          5,
		MaxSize:              10,
		MaxIdle:              5,
		Factory:              factory,
		MinEvictableIdleTime: 5 * time.Minute,
		ValidationInterval:   30 * time.Second,
	}
}

// ObjectPool 表示对象池，管理对象的创建、获取和回收
type ObjectPool struct {
	// 池配置
	config PoolConfig

	// 可用对象的通道
	idle chan Object

	// 用于同步的互斥锁
	mu sync.Mutex

	// 所有对象的映射(包括活跃和空闲对象)
	objects map[int]poolObject

	// 跟踪最后一次归还的时间戳
	lastReturn map[int]time.Time

	// 跟踪正在使用的对象数量
	activeCount int

	// 控制后台清理的停止信号
	stopCleaner chan struct{}

	// 指示池是否已关闭
	closed bool

	// 统计信息
	stats PoolStats
}

// poolObject 表示对象池中的一个对象及其状态
type poolObject struct {
	obj    Object
	active bool
}

// PoolStats 记录池的使用统计信息
type PoolStats struct {
	// 创建的对象总数
	Created int

	// 获取的对象总数
	Acquired int

	// 释放的对象总数
	Released int

	// 丢弃的对象总数
	Destroyed int

	// 等待对象的总时间
	WaitTime time.Duration

	// 对象处于活跃状态的总时间
	ActiveTime time.Duration

	// 最大等待时间
	MaxWaitTime time.Duration

	// 池满导致等待的次数
	Waits int

	// 超时错误总数
	Timeouts int
}

// NewObjectPool 创建并初始化一个对象池
func NewObjectPool(config PoolConfig) (*ObjectPool, error) {
	if config.Factory == nil {
		return nil, errors.New("factory function required")
	}

	if config.InitialSize > config.MaxSize {
		config.InitialSize = config.MaxSize
	}

	if config.MaxIdle > config.MaxSize {
		config.MaxIdle = config.MaxSize
	}

	pool := &ObjectPool{
		config:      config,
		idle:        make(chan Object, config.MaxSize),
		objects:     make(map[int]poolObject),
		lastReturn:  make(map[int]time.Time),
		stopCleaner: make(chan struct{}),
	}

	// 初始化对象
	for i := 0; i < config.InitialSize; i++ {
		obj, err := config.Factory()
		if err != nil {
			// 如果创建失败，释放已创建的对象
			pool.Close()
			return nil, err
		}

		pool.idle <- obj
		pool.objects[obj.ID()] = poolObject{obj: obj, active: false}
		pool.stats.Created++
	}

	// 启动后台清理协程
	go pool.periodicCleaning()

	return pool, nil
}

// periodicCleaning 定期检查并清理空闲对象
func (p *ObjectPool) periodicCleaning() {
	ticker := time.NewTicker(p.config.ValidationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.evictExpiredObjects()
		case <-p.stopCleaner:
			return
		}
	}
}

// evictExpiredObjects 清除长时间未使用的空闲对象
func (p *ObjectPool) evictExpiredObjects() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	now := time.Now()
	idleCount := len(p.idle)

	// 如果空闲对象少于或等于MaxIdle,不执行清理
	if idleCount <= p.config.MaxIdle {
		return
	}

	// 尝试从通道获取对象并检查它们
	toRemove := idleCount - p.config.MaxIdle
	for i := 0; i < toRemove; i++ {
		select {
		case obj := <-p.idle:
			lastUsed, exists := p.lastReturn[obj.ID()]
			// 如果对象长时间未使用,或者无效,则销毁它
			if !exists || now.Sub(lastUsed) > p.config.MinEvictableIdleTime || !obj.Validate() {
				delete(p.objects, obj.ID())
				delete(p.lastReturn, obj.ID())
				p.stats.Destroyed++
			} else {
				// 对象仍然有效,放回通道
				p.idle <- obj
			}
		default:
			// 通道为空,停止清理
			return
		}
	}
}

// AcquireWithTimeout 尝试在指定的超时时间内从池中获取对象
func (p *ObjectPool) AcquireWithTimeout(timeout time.Duration) (Object, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}

	startTime := time.Now()

	// 尝试从空闲对象池获取
	select {
	case obj, ok := <-p.idle:
		if !ok {
			return nil, ErrPoolClosed
		}

		// 更新对象状态和统计信息
		p.mu.Lock()
		info := p.objects[obj.ID()]
		info.active = true
		p.objects[obj.ID()] = info
		p.activeCount++
		waitTime := time.Since(startTime)
		p.stats.WaitTime += waitTime
		p.stats.Acquired++
		if waitTime > p.stats.MaxWaitTime {
			p.stats.MaxWaitTime = waitTime
		}
		p.mu.Unlock()

		// 验证对象并在必要时重置
		if !obj.Validate() {
			p.discardObject(obj)
			return p.createNewObject()
		}

		return obj, nil

	case <-time.After(timeout):
		// 尝试创建新对象(如果池未满)
		p.mu.Lock()
		canCreate := len(p.objects) < p.config.MaxSize
		p.mu.Unlock()

		if canCreate {
			return p.createNewObject()
		}

		// 池已满且等待超时
		p.mu.Lock()
		p.stats.Timeouts++
		p.mu.Unlock()
		return nil, ErrPoolTimeout
	}
}

// createNewObject 创建一个新对象并添加到池中
func (p *ObjectPool) createNewObject() (Object, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 再次检查容量(避免竞态条件)
	if len(p.objects) >= p.config.MaxSize {
		p.stats.Waits++
		return nil, ErrPoolAtMaxCapacity
	}

	// 创建新对象
	obj, err := p.config.Factory()
	if err != nil {
		return nil, err
	}

	// 记录新对象
	p.objects[obj.ID()] = poolObject{obj: obj, active: true}
	p.activeCount++
	p.stats.Created++
	p.stats.Acquired++

	return obj, nil
}

// AcquireObject 从对象池中获取对象(默认使用1秒超时)
func (p *ObjectPool) AcquireObject() (Object, error) {
	return p.AcquireWithTimeout(1 * time.Second)
}

// ReleaseObject 将对象归还给对象池
func (p *ObjectPool) ReleaseObject(obj Object) error {
	if p.closed {
		return ErrPoolClosed
	}

	if obj == nil {
		return ErrInvalidObject
	}

	p.mu.Lock()
	// 检查对象是否属于这个池
	info, exists := p.objects[obj.ID()]
	if !exists || !info.active {
		p.mu.Unlock()
		return ErrInvalidObject
	}

	// 更新状态和统计信息
	info.active = false
	p.objects[obj.ID()] = info
	p.activeCount--
	p.lastReturn[obj.ID()] = time.Now()
	p.stats.Released++
	p.mu.Unlock()

	// 重置对象状态
	if err := obj.Reset(); err != nil {
		return p.discardObject(obj)
	}

	// 如果对象无效,丢弃它
	if !obj.Validate() {
		return p.discardObject(obj)
	}

	// 将对象归还到池中
	select {
	case p.idle <- obj:
		return nil
	default:
		// 如果通道已满,丢弃对象
		return p.discardObject(obj)
	}
}

// discardObject 从池中移除无效对象
func (p *ObjectPool) discardObject(obj Object) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.objects, obj.ID())
	delete(p.lastReturn, obj.ID())
	p.stats.Destroyed++
	return nil
}

// Close 关闭对象池,清理资源
func (p *ObjectPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.stopCleaner)

	// 清空通道
	close(p.idle)
	// 修复"declared and not used"错误: 使用匿名变量接收通道值
	for _ = range p.idle {
		p.stats.Destroyed++
	}

	// 清空映射
	p.objects = nil
	p.lastReturn = nil
}

// Status 返回池的当前状态信息
func (p *ObjectPool) Status() (active int, idle int, total int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.activeCount, len(p.idle), len(p.objects)
}

// Stats 返回池的统计信息
func (p *ObjectPool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.stats
}
