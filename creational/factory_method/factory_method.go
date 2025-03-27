package factory_method

import (
	"fmt"
	"sync"
)

// Logger 是日志记录器接口
type Logger interface {
	Log(message string)
}

// FileLogger 实现了记录到文件的日志记录器
type FileLogger struct{}

func (f *FileLogger) Log(message string) {
	fmt.Println("Log to file: " + message)
}

// ConsoleLogger 实现了记录到控制台的日志记录器
type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(message string) {
	fmt.Println("Log to console: " + message)
}

// NetworkLogger 是一个新的实现，用于记录到网络
type NetworkLogger struct {
	endpoint string
}

func NewNetworkLogger(endpoint string) *NetworkLogger {
	return &NetworkLogger{endpoint: endpoint}
}

func (n *NetworkLogger) Log(message string) {
	fmt.Printf("Log to network endpoint %s: %s\n", n.endpoint, message)
}

// LoggerFactory 是定义工厂方法的接口
type LoggerFactory interface {
	CreateLogger() Logger
}

// FileLoggerFactory 创建 FileLogger 实例
type FileLoggerFactory struct {
	instance Logger
	once     sync.Once
}

func (f *FileLoggerFactory) CreateLogger() Logger {
	// 使用单例模式和懒初始化
	f.once.Do(func() {
		f.instance = &FileLogger{}
	})
	return f.instance
}

// ConsoleLoggerFactory 创建 ConsoleLogger 实例
type ConsoleLoggerFactory struct {
	instance Logger
	once     sync.Once
}

func (c *ConsoleLoggerFactory) CreateLogger() Logger {
	// 使用单例模式和懒初始化
	c.once.Do(func() {
		c.instance = &ConsoleLogger{}
	})
	return c.instance
}

// NetworkLoggerFactory 创建 NetworkLogger 实例
type NetworkLoggerFactory struct {
	endpoint string
}

func NewNetworkLoggerFactory(endpoint string) *NetworkLoggerFactory {
	return &NetworkLoggerFactory{endpoint: endpoint}
}

func (n *NetworkLoggerFactory) CreateLogger() Logger {
	return NewNetworkLogger(n.endpoint)
}

// LoggerType 表示要创建的日志记录器类型
type LoggerType string

const (
	FileType    LoggerType = "file"
	ConsoleType LoggerType = "console"
	NetworkType LoggerType = "network"
)

// GetLoggerFactory 根据类型返回合适的工厂
func GetLoggerFactory(loggerType LoggerType, config map[string]string) (LoggerFactory, error) {
	switch loggerType {
	case FileType:
		return &FileLoggerFactory{}, nil
	case ConsoleType:
		return &ConsoleLoggerFactory{}, nil
	case NetworkType:
		endpoint, ok := config["endpoint"]
		if !ok {
			return nil, fmt.Errorf("网络日志记录器需要endpoint配置")
		}
		return NewNetworkLoggerFactory(endpoint), nil
	default:
		return nil, fmt.Errorf("不支持的日志记录器类型: %s", loggerType)
	}
}
