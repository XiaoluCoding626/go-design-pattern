package context

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试上下文值传递函数
func TestContextValues(t *testing.T) {
	t.Run("RequestInfo", func(t *testing.T) {
		ctx := context.Background()
		info := RequestInfo{
			Username:  "testuser",
			IPAddress: "127.0.0.1",
			Timestamp: time.Now(),
		}

		// 测试添加RequestInfo到上下文
		enrichedCtx := WithRequestInfo(ctx, info)
		assert.NotEqual(t, enrichedCtx, ctx, "WithRequestInfo应当返回新的上下文")

		// 测试从上下文获取RequestInfo
		retrievedInfo, ok := GetRequestInfo(enrichedCtx)
		assert.True(t, ok, "应能够从上下文获取RequestInfo")
		assert.Equal(t, info.Username, retrievedInfo.Username, "用户名应匹配")
		assert.Equal(t, info.IPAddress, retrievedInfo.IPAddress, "IP地址应匹配")

		// 测试从空上下文获取
		_, ok = GetRequestInfo(context.Background())
		assert.False(t, ok, "不应能从空上下文获取RequestInfo")
	})

	t.Run("RequestID", func(t *testing.T) {
		ctx := context.Background()

		// 测试添加RequestID
		enrichedCtx := WithRequestID(ctx)
		assert.NotEqual(t, enrichedCtx, ctx, "WithRequestID应当返回新的上下文")

		// 测试获取RequestID
		id, ok := GetRequestID(enrichedCtx)
		assert.True(t, ok, "应能够从上下文获取RequestID")
		assert.NotEmpty(t, id, "RequestID不应为空")

		// 确保每次生成的ID都不同
		anotherCtx := WithRequestID(ctx)
		anotherId, ok := GetRequestID(anotherCtx)
		assert.True(t, ok, "应能够从上下文获取RequestID")
		assert.NotEqual(t, id, anotherId, "每个RequestID应该是唯一的")
	})

	t.Run("UserToken", func(t *testing.T) {
		ctx := context.Background()
		token := "auth-token-123"

		// 测试添加UserToken
		enrichedCtx := WithUserToken(ctx, token)

		// 测试获取UserToken
		retrievedToken, ok := GetUserToken(enrichedCtx)
		assert.True(t, ok, "应能够从上下文获取UserToken")
		assert.Equal(t, token, retrievedToken, "获取到的UserToken应与存入的匹配")
	})
}

// 测试错误映射函数
func TestMapContextError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "Context canceled",
			err:      context.Canceled,
			expected: ErrRequestCancelled,
		},
		{
			name:     "Context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: ErrRequestTimeout,
		},
		{
			name:     "Other error",
			err:      errors.New("some other error"),
			expected: errors.New("some other error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mappedErr := mapContextError(tt.err)
			assert.Equal(t, tt.expected.Error(), mappedErr.Error(), "错误映射应正确")
		})
	}
}

// 测试请求处理成功流程
func TestProcessRequest_Success(t *testing.T) {
	ctx := context.Background()
	info := RequestInfo{
		Username:  "testuser",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// 使用足够长的超时以确保请求能完成
	err := ProcessRequest(ctx, info, 5*time.Second)
	assert.NoError(t, err, "请求处理应该成功")
}

// 测试请求超时
func TestProcessRequest_Timeout(t *testing.T) {
	ctx := context.Background()
	info := RequestInfo{
		Username:  "testuser",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// 使用极短的超时，确保超时发生
	err := ProcessRequest(ctx, info, 1*time.Nanosecond)
	assert.ErrorIs(t, err, ErrRequestTimeout, "应返回超时错误")
}

// 测试请求取消
func TestProcessRequest_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	info := RequestInfo{
		Username:  "testuser",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// 在另一个goroutine中取消请求
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// 处理请求，应该被取消
	err := ProcessRequest(ctx, info, 5*time.Second)
	assert.ErrorIs(t, err, ErrRequestCancelled, "应返回取消错误")
}

// 测试信息缺失的情况
func TestValidateRequest_MissingInfo(t *testing.T) {
	ctx := context.Background()
	// 没有添加RequestInfo

	err := validateRequest(ctx)
	assert.Error(t, err, "当RequestInfo缺失时应返回错误")
}

// 测试业务逻辑的并发执行
func TestProcessBusinessLogic(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx)
	ctx = WithRequestInfo(ctx, RequestInfo{Username: "testuser", IPAddress: "127.0.0.1"})

	// 业务逻辑处理应成功完成
	err := processBusinessLogic(ctx)
	assert.NoError(t, err, "业务逻辑处理应成功")

	// 测试取消的情况
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel() // 立即取消

	err = processBusinessLogic(ctxWithCancel)
	assert.ErrorIs(t, err, ErrRequestCancelled, "取消上下文后应返回取消错误")
}

// 测试业务逻辑和子任务的取消传播
func TestCancellationPropagation(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestInfo(ctx, RequestInfo{Username: "testuser", IPAddress: "127.0.0.1"})
	ctx = WithRequestID(ctx)

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(ctx)

	// 启动一个goroutine在短时间后取消上下文
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	// 验证processData函数能正确响应取消
	err := processData(ctx)
	assert.ErrorIs(t, err, ErrRequestCancelled, "processData应该检测到取消信号")

	// 验证updateStatus函数能正确响应取消
	err = updateStatus(ctx)
	assert.ErrorIs(t, err, ErrRequestCancelled, "updateStatus应该检测到取消信号")
}

// 集成测试 - 模拟完整请求流程
func TestCompleteRequestFlow(t *testing.T) {
	// 创建基础上下文
	baseCtx := context.Background()

	// 准备请求信息
	requestInfo := RequestInfo{
		Username:  "integration_test_user",
		IPAddress: "10.0.0.1",
		Timestamp: time.Now(),
	}

	// 测试成功流程
	t.Run("SuccessfulRequest", func(t *testing.T) {
		err := ProcessRequest(baseCtx, requestInfo, 3*time.Second)
		assert.NoError(t, err, "完整请求流程应成功")
	})

	// 测试取消流程
	t.Run("CancelledRequest", func(t *testing.T) {
		ctx, cancel := context.WithCancel(baseCtx)

		// 在短暂延迟后取消
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		err := ProcessRequest(ctx, requestInfo, 3*time.Second)
		assert.ErrorIs(t, err, ErrRequestCancelled, "应返回取消错误")
	})

	// 测试超时流程
	t.Run("TimeoutRequest", func(t *testing.T) {
		// 使用极短的超时
		err := ProcessRequest(baseCtx, requestInfo, 1*time.Nanosecond)
		assert.ErrorIs(t, err, ErrRequestTimeout, "应返回超时错误")
	})
}

// 模拟上下文链
func TestContextChain(t *testing.T) {
	// 创建基本上下文，然后是包含令牌的上下文，然后是包含请求信息的上下文
	baseCtx := context.Background()
	tokenCtx := WithUserToken(baseCtx, "auth-token-123")
	requestCtx := WithRequestInfo(tokenCtx, RequestInfo{
		Username:  "chainuser",
		IPAddress: "192.168.1.1",
		Timestamp: time.Now(),
	})
	idCtx := WithRequestID(requestCtx)

	// 验证我们可以从最终上下文获取所有值
	token, ok := GetUserToken(idCtx)
	assert.True(t, ok, "应能从链式上下文获取用户令牌")
	assert.Equal(t, "auth-token-123", token, "令牌值应正确")

	info, ok := GetRequestInfo(idCtx)
	assert.True(t, ok, "应能从链式上下文获取请求信息")
	assert.Equal(t, "chainuser", info.Username, "用户名应正确")

	_, ok = GetRequestID(idCtx)
	assert.True(t, ok, "应能从链式上下文获取请求ID")
}

// 基准测试 - 测量处理请求的性能
func BenchmarkProcessRequest(b *testing.B) {
	ctx := context.Background()
	info := RequestInfo{
		Username:  "benchuser",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	for i := 0; i < b.N; i++ {
		_ = ProcessRequest(ctx, info, 5*time.Second)
	}
}

// 提供一个简单的演示函数
func Example() {
	// 创建基础上下文
	baseCtx := context.Background()

	// 准备请求信息
	requestInfo := RequestInfo{
		Username:  "user123",
		IPAddress: "192.168.1.100",
		Timestamp: time.Now(),
	}

	// 处理请求，设置3秒超时
	err := ProcessRequest(baseCtx, requestInfo, 3*time.Second)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	log.Println("Request completed successfully")

	// 演示取消请求
	cancelCtx, cancel := context.WithCancel(baseCtx)
	go func() {
		time.Sleep(500 * time.Millisecond)
		log.Println("Cancelling request...")
		cancel()
	}()

	err = ProcessRequest(cancelCtx, requestInfo, 3*time.Second)
	if errors.Is(err, ErrRequestCancelled) {
		log.Println("Request was cancelled as expected")
	} else if err != nil {
		log.Printf("Unexpected error: %v", err)
	}
}
