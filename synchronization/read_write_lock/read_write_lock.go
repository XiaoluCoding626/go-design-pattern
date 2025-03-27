package read_write_lock

import (
	"sync"
	"time"
)

// RWLocker 定义了读写锁接口
// 通过接口可以方便地替换不同的读写锁实现或用于模拟测试
type RWLocker interface {
	ReadLock()                                  // 获取读锁
	ReadUnlock()                                // 释放读锁
	WriteLock()                                 // 获取写锁
	WriteUnlock()                               // 释放写锁
	TryReadLock() bool                          // 尝试获取读锁，不阻塞
	TryWriteLock() bool                         // 尝试获取写锁，不阻塞
	TryReadLockWithTimeout(time.Duration) bool  // 带超时的尝试获取读锁
	TryWriteLockWithTimeout(time.Duration) bool // 带超时的尝试获取写锁
}

// StandardRWLock 标准读写锁实现，封装Go的sync.RWMutex
type StandardRWLock struct {
	rwMutex sync.RWMutex
}

// NewStandardRWLock 创建一个新的标准读写锁
func NewStandardRWLock() *StandardRWLock {
	return &StandardRWLock{}
}

// ReadLock 获取读锁
func (l *StandardRWLock) ReadLock() {
	l.rwMutex.RLock()
}

// ReadUnlock 释放读锁
func (l *StandardRWLock) ReadUnlock() {
	l.rwMutex.RUnlock()
}

// WriteLock 获取写锁
func (l *StandardRWLock) WriteLock() {
	l.rwMutex.Lock()
}

// WriteUnlock 释放写锁
func (l *StandardRWLock) WriteUnlock() {
	l.rwMutex.Unlock()
}

// TryReadLock 尝试获取读锁，不阻塞，若获取成功则返回true
func (l *StandardRWLock) TryReadLock() bool {
	return l.rwMutex.TryRLock()
}

// TryWriteLock 尝试获取写锁，不阻塞，若获取成功则返回true
func (l *StandardRWLock) TryWriteLock() bool {
	return l.rwMutex.TryLock()
}

// TryReadLockWithTimeout 尝试在指定时间内获取读锁
func (l *StandardRWLock) TryReadLockWithTimeout(timeout time.Duration) bool {
	success := make(chan bool, 1)

	go func() {
		success <- l.rwMutex.TryRLock()
	}()

	select {
	case result := <-success:
		return result
	case <-time.After(timeout):
		return false
	}
}

// TryWriteLockWithTimeout 尝试在指定时间内获取写锁
func (l *StandardRWLock) TryWriteLockWithTimeout(timeout time.Duration) bool {
	success := make(chan bool, 1)

	go func() {
		success <- l.rwMutex.TryLock()
	}()

	select {
	case result := <-success:
		return result
	case <-time.After(timeout):
		return false
	}
}

// Data 表示包含读写锁保护的共享数据
type Data struct {
	locker RWLocker // 使用接口允许注入不同的读写锁实现
	value  int      // 数据值
}

// NewData 创建一个新的数据实例，使用标准读写锁
func NewData() *Data {
	return &Data{
		locker: NewStandardRWLock(),
	}
}

// NewDataWithLocker 使用指定的读写锁创建数据实例
func NewDataWithLocker(locker RWLocker) *Data {
	return &Data{
		locker: locker,
	}
}

// Read 读取数据值，使用读锁保证并发安全
func (d *Data) Read() int {
	d.locker.ReadLock()
	defer d.locker.ReadUnlock()

	return d.value
}

// TryRead 尝试读取数据值，不阻塞
// 如果当前有写锁，则返回false和0值
func (d *Data) TryRead() (int, bool) {
	if !d.locker.TryReadLock() {
		return 0, false
	}
	defer d.locker.ReadUnlock()

	return d.value, true
}

// ReadWithTimeout 尝试在指定时间内读取数据
func (d *Data) ReadWithTimeout(timeout time.Duration) (int, bool) {
	if !d.locker.TryReadLockWithTimeout(timeout) {
		return 0, false
	}
	defer d.locker.ReadUnlock()

	return d.value, true
}

// Write 写入数据值，使用写锁保证并发安全
func (d *Data) Write(val int) bool {
	d.locker.WriteLock()
	defer d.locker.WriteUnlock()

	d.value = val
	return true
}

// TryWrite 尝试写入数据，不阻塞
// 如果当前有其他读锁或写锁，则返回false
func (d *Data) TryWrite(val int) bool {
	if !d.locker.TryWriteLock() {
		return false
	}
	defer d.locker.WriteUnlock()

	d.value = val
	return true
}

// WriteWithTimeout 尝试在指定时间内写入数据
func (d *Data) WriteWithTimeout(val int, timeout time.Duration) bool {
	if !d.locker.TryWriteLockWithTimeout(timeout) {
		return false
	}
	defer d.locker.WriteUnlock()

	d.value = val
	return true
}

// ReadWithCallback 在读锁保护下执行自定义读操作
func (d *Data) ReadWithCallback(callback func(val int)) {
	d.locker.ReadLock()
	defer d.locker.ReadUnlock()

	callback(d.value)
}

// WriteWithCallback 在写锁保护下执行自定义写操作
func (d *Data) WriteWithCallback(callback func(d *Data)) {
	d.locker.WriteLock()
	defer d.locker.WriteUnlock()

	callback(d)
}

// ReadWriteWithCallback 先获取读锁执行读操作，然后升级为写锁执行写操作
// 注意：这个方法不是原子的，不是真正的锁升级，中间会释放读锁
func (d *Data) ReadWriteWithCallback(readCallback func(val int) int) {
	// 先读取
	val := d.Read()

	// 根据读取结果计算新值
	newVal := readCallback(val)

	// 写入新值
	d.Write(newVal)
}
