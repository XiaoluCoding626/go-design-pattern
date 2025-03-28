package condition_variable

import (
	"fmt"
	"sync"
)

// Data 存储数据并提供同步机制
type Data struct {
	ready bool
	cond  *sync.Cond
}

// NewData 创建并初始化Data实例
func NewData() *Data {
	return &Data{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// WaitForReady 阻塞直到数据准备就绪
func (d *Data) WaitForReady() {
	d.cond.L.Lock()
	defer d.cond.L.Unlock()

	for !d.ready {
		d.cond.Wait()
	}

	fmt.Println("Data is ready")
}

// SetReady 设置数据为就绪状态并通知等待者
func (d *Data) SetReady() {
	d.cond.L.Lock()
	defer d.cond.L.Unlock()

	d.ready = true
	d.cond.Signal()
}
