package factory_method

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// 用于捕获标准输出的辅助函数
func captureOutput(f func()) string {
	// 保存原始的标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行函数
	f()

	// 恢复原始的标准输出
	w.Close()
	os.Stdout = oldStdout

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

// 测试 FileLogger
func TestFileLogger(t *testing.T) {
	logger := &FileLogger{}
	output := captureOutput(func() {
		logger.Log("测试消息")
	})

	expected := "Log to file: 测试消息\n"
	if output != expected {
		t.Errorf("FileLogger.Log() 输出 = %q, 期望 %q", output, expected)
	}
}

// 测试 ConsoleLogger
func TestConsoleLogger(t *testing.T) {
	logger := &ConsoleLogger{}
	output := captureOutput(func() {
		logger.Log("测试消息")
	})

	expected := "Log to console: 测试消息\n"
	if output != expected {
		t.Errorf("ConsoleLogger.Log() 输出 = %q, 期望 %q", output, expected)
	}
}

// 测试 NetworkLogger
func TestNetworkLogger(t *testing.T) {
	logger := NewNetworkLogger("http://example.com/log")
	output := captureOutput(func() {
		logger.Log("测试消息")
	})

	expected := "Log to network endpoint http://example.com/log: 测试消息\n"
	if output != expected {
		t.Errorf("NetworkLogger.Log() 输出 = %q, 期望 %q", output, expected)
	}
}

// 测试 FileLoggerFactory
func TestFileLoggerFactory(t *testing.T) {
	factory := &FileLoggerFactory{}
	logger1 := factory.CreateLogger()
	logger2 := factory.CreateLogger()

	// 验证工厂创建的是正确类型的日志记录器
	_, ok := logger1.(*FileLogger)
	if !ok {
		t.Error("FileLoggerFactory.CreateLogger() 应该返回 *FileLogger")
	}

	// 验证单例模式工作正常
	if logger1 != logger2 {
		t.Error("FileLoggerFactory 应该返回相同的实例（单例模式）")
	}

	// 验证日志记录器工作正常
	output := captureOutput(func() {
		logger1.Log("测试消息")
	})

	expected := "Log to file: 测试消息\n"
	if output != expected {
		t.Errorf("从工厂创建的 FileLogger.Log() 输出 = %q, 期望 %q", output, expected)
	}
}

// 测试 ConsoleLoggerFactory
func TestConsoleLoggerFactory(t *testing.T) {
	factory := &ConsoleLoggerFactory{}
	logger1 := factory.CreateLogger()
	logger2 := factory.CreateLogger()

	// 验证工厂创建的是正确类型的日志记录器
	_, ok := logger1.(*ConsoleLogger)
	if !ok {
		t.Error("ConsoleLoggerFactory.CreateLogger() 应该返回 *ConsoleLogger")
	}

	// 验证单例模式工作正常
	if logger1 != logger2 {
		t.Error("ConsoleLoggerFactory 应该返回相同的实例（单例模式）")
	}

	// 验证日志记录器工作正常
	output := captureOutput(func() {
		logger1.Log("测试消息")
	})

	expected := "Log to console: 测试消息\n"
	if output != expected {
		t.Errorf("从工厂创建的 ConsoleLogger.Log() 输出 = %q, 期望 %q", output, expected)
	}
}

// 测试 NetworkLoggerFactory
func TestNetworkLoggerFactory(t *testing.T) {
	factory := NewNetworkLoggerFactory("http://example.com/log")
	logger := factory.CreateLogger()

	// 验证工厂创建的是正确类型的日志记录器
	netLogger, ok := logger.(*NetworkLogger)
	if !ok {
		t.Error("NetworkLoggerFactory.CreateLogger() 应该返回 *NetworkLogger")
	}

	// 验证 endpoint 设置正确
	if netLogger.endpoint != "http://example.com/log" {
		t.Errorf("NetworkLogger.endpoint = %q, 期望 %q", netLogger.endpoint, "http://example.com/log")
	}

	// 验证日志记录器工作正常
	output := captureOutput(func() {
		logger.Log("测试消息")
	})

	expected := "Log to network endpoint http://example.com/log: 测试消息\n"
	if output != expected {
		t.Errorf("从工厂创建的 NetworkLogger.Log() 输出 = %q, 期望 %q", output, expected)
	}
}

// 测试 GetLoggerFactory
func TestGetLoggerFactory(t *testing.T) {
	// 测试 FileLoggerFactory
	fileFactory, err := GetLoggerFactory(FileType, nil)
	if err != nil {
		t.Errorf("GetLoggerFactory(FileType) 返回错误: %v", err)
	}

	_, ok := fileFactory.(*FileLoggerFactory)
	if !ok {
		t.Error("GetLoggerFactory(FileType) 应该返回 *FileLoggerFactory")
	}

	// 测试 ConsoleLoggerFactory
	consoleFactory, err := GetLoggerFactory(ConsoleType, nil)
	if err != nil {
		t.Errorf("GetLoggerFactory(ConsoleType) 返回错误: %v", err)
	}

	_, ok = consoleFactory.(*ConsoleLoggerFactory)
	if !ok {
		t.Error("GetLoggerFactory(ConsoleType) 应该返回 *ConsoleLoggerFactory")
	}

	// 测试 NetworkLoggerFactory（带有有效配置）
	config := map[string]string{"endpoint": "http://example.com/log"}
	networkFactory, err := GetLoggerFactory(NetworkType, config)
	if err != nil {
		t.Errorf("GetLoggerFactory(NetworkType) 返回错误: %v", err)
	}

	netFactory, ok := networkFactory.(*NetworkLoggerFactory)
	if !ok {
		t.Error("GetLoggerFactory(NetworkType) 应该返回 *NetworkLoggerFactory")
	}

	if netFactory.endpoint != "http://example.com/log" {
		t.Errorf("NetworkLoggerFactory.endpoint = %q, 期望 %q", netFactory.endpoint, "http://example.com/log")
	}

	// 测试 NetworkLoggerFactory（缺少配置）
	_, err = GetLoggerFactory(NetworkType, map[string]string{})
	if err == nil {
		t.Error("GetLoggerFactory(NetworkType) 在缺少 endpoint 配置时应该返回错误")
	}

	if !strings.Contains(err.Error(), "endpoint配置") {
		t.Errorf("错误消息 = %q, 应该包含 'endpoint配置'", err.Error())
	}

	// 测试不支持的日志记录器类型
	_, err = GetLoggerFactory("unknown", nil)
	if err == nil {
		t.Error("GetLoggerFactory('unknown') 应该返回错误")
	}

	if !strings.Contains(err.Error(), "不支持的日志记录器类型") {
		t.Errorf("错误消息 = %q, 应该包含 '不支持的日志记录器类型'", err.Error())
	}
}

// 测试集成场景
func TestIntegration(t *testing.T) {
	// 创建并使用文件日志记录器
	fileFactory, _ := GetLoggerFactory(FileType, nil)
	fileLogger := fileFactory.CreateLogger()
	fileOutput := captureOutput(func() {
		fileLogger.Log("文件日志消息")
	})

	// 创建并使用控制台日志记录器
	consoleFactory, _ := GetLoggerFactory(ConsoleType, nil)
	consoleLogger := consoleFactory.CreateLogger()
	consoleOutput := captureOutput(func() {
		consoleLogger.Log("控制台日志消息")
	})

	// 创建并使用网络日志记录器
	config := map[string]string{"endpoint": "http://example.com/log"}
	networkFactory, _ := GetLoggerFactory(NetworkType, config)
	networkLogger := networkFactory.CreateLogger()
	networkOutput := captureOutput(func() {
		networkLogger.Log("网络日志消息")
	})

	// 验证所有输出
	expectedFileOutput := "Log to file: 文件日志消息\n"
	if fileOutput != expectedFileOutput {
		t.Errorf("文件日志输出 = %q, 期望 %q", fileOutput, expectedFileOutput)
	}

	expectedConsoleOutput := "Log to console: 控制台日志消息\n"
	if consoleOutput != expectedConsoleOutput {
		t.Errorf("控制台日志输出 = %q, 期望 %q", consoleOutput, expectedConsoleOutput)
	}

	expectedNetworkOutput := "Log to network endpoint http://example.com/log: 网络日志消息\n"
	if networkOutput != expectedNetworkOutput {
		t.Errorf("网络日志输出 = %q, 期望 %q", networkOutput, expectedNetworkOutput)
	}
}
