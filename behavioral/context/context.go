package context

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// 定义上下文键类型，避免类型冲突
type contextKey string

// 定义上下文键常量
const (
	requestInfoKey contextKey = "requestInfo"
	requestIDKey   contextKey = "requestID"
	userTokenKey   contextKey = "userToken"
)

// RequestInfo 包含请求相关信息
type RequestInfo struct {
	Username  string
	IPAddress string
	Timestamp time.Time
}

// 自定义错误
var (
	ErrRequestCancelled = errors.New("request was cancelled")
	ErrRequestTimeout   = errors.New("request timed out")
	ErrUnauthorized     = errors.New("unauthorized request")
)

// WithRequestInfo 将请求信息添加到上下文中
func WithRequestInfo(ctx context.Context, info RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfoKey, info)
}

// GetRequestInfo 从上下文中获取请求信息
func GetRequestInfo(ctx context.Context) (RequestInfo, bool) {
	info, ok := ctx.Value(requestInfoKey).(RequestInfo)
	return info, ok
}

// WithRequestID 将请求ID添加到上下文中
func WithRequestID(ctx context.Context) context.Context {
	id := fmt.Sprintf("req-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID 从上下文中获取请求ID
func GetRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

// WithUserToken 将用户令牌添加到上下文中
func WithUserToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, userTokenKey, token)
}

// GetUserToken 从上下文中获取用户令牌
func GetUserToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(userTokenKey).(string)
	return token, ok
}

// ProcessRequest 处理请求的入口函数
// 使用超时控制和取消信号处理
func ProcessRequest(parentCtx context.Context, info RequestInfo, timeout time.Duration) error {
	// 1. 创建请求上下文
	ctx := WithRequestInfo(parentCtx, info)
	ctx = WithRequestID(ctx)

	// 2. 添加超时控制
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel() // 确保资源被释放

	// 3. 记录请求开始
	requestID, _ := GetRequestID(ctx)
	log.Printf("[%s] Starting request processing for user %s from %s",
		requestID, info.Username, info.IPAddress)

	// 4. 执行多阶段处理并传递上下文
	if err := validateRequest(ctx); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := processBusinessLogic(ctx); err != nil {
		return fmt.Errorf("business logic failed: %w", err)
	}

	if err := saveResults(ctx); err != nil {
		return fmt.Errorf("saving results failed: %w", err)
	}

	// 5. 记录请求完成
	log.Printf("[%s] Request processing completed successfully", requestID)
	return nil
}

// validateRequest 验证请求
func validateRequest(ctx context.Context) error {
	// 检查上下文是否已取消
	if err := ctx.Err(); err != nil {
		return mapContextError(err)
	}

	info, ok := GetRequestInfo(ctx)
	if !ok {
		return errors.New("request info not found in context")
	}

	// 模拟请求验证
	log.Printf("Validating request from user %s at IP %s", info.Username, info.IPAddress)

	// 模拟验证工作
	time.Sleep(200 * time.Millisecond)

	return nil
}

// processBusinessLogic 处理业务逻辑
func processBusinessLogic(ctx context.Context) error {
	// 创建一个工作组来执行并行任务
	var wg sync.WaitGroup
	errCh := make(chan error, 2) // 错误通道

	// 启动两个并行任务
	wg.Add(2)

	// 任务1: 数据处理
	go func() {
		defer wg.Done()
		if err := processData(ctx); err != nil {
			errCh <- err
		}
	}()

	// 任务2: 更新状态
	go func() {
		defer wg.Done()
		if err := updateStatus(ctx); err != nil {
			errCh <- err
		}
	}()

	// 等待任务完成
	wg.Wait()
	close(errCh)

	// 检查错误
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

// processData 处理数据
func processData(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return mapContextError(ctx.Err())
	case <-time.After(500 * time.Millisecond):
		requestID, _ := GetRequestID(ctx)
		info, _ := GetRequestInfo(ctx)
		log.Printf("[%s] Processed data for user %s", requestID, info.Username)
	}
	return nil
}

// updateStatus 更新状态
func updateStatus(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return mapContextError(ctx.Err())
	case <-time.After(300 * time.Millisecond):
		requestID, _ := GetRequestID(ctx)
		info, _ := GetRequestInfo(ctx)
		log.Printf("[%s] Updated status for user %s", requestID, info.Username)
	}
	return nil
}

// saveResults 保存结果
func saveResults(ctx context.Context) error {
	// 检查上下文是否已取消
	if err := ctx.Err(); err != nil {
		return mapContextError(err)
	}

	requestID, _ := GetRequestID(ctx)
	info, _ := GetRequestInfo(ctx)

	// 模拟保存结果
	log.Printf("[%s] Saving results for user %s", requestID, info.Username)
	time.Sleep(400 * time.Millisecond)

	return nil
}

// mapContextError 将context错误映射到自定义错误
func mapContextError(err error) error {
	switch err {
	case context.Canceled:
		return ErrRequestCancelled
	case context.DeadlineExceeded:
		return ErrRequestTimeout
	default:
		return err
	}
}
