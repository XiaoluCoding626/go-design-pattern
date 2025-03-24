// Package functionaloption 的测试文件，验证函数选项模式实现的正确性
package functionaloption

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"
)

// 测试创建默认的HTTP客户端
func TestNewHTTPClientDefault(t *testing.T) {
	client := NewHTTPClient()

	// 验证默认超时
	if client.Timeout != 30*time.Second {
		t.Errorf("默认超时应为30秒，实际为%v", client.Timeout)
	}

	// 验证其他默认值（通过类型断言获取transport配置）
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("无法获取transport配置")
	}

	// 验证默认连接池大小
	if transport.MaxIdleConns != 100 {
		t.Errorf("默认最大空闲连接数应为100，实际为%d", transport.MaxIdleConns)
	}

	// 验证默认空闲超时
	if transport.IdleConnTimeout != 90*time.Second {
		t.Errorf("默认空闲连接超时应为90秒，实际为%v", transport.IdleConnTimeout)
	}
}

// 测试超时选项
func TestWithTimeout(t *testing.T) {
	client := NewHTTPClient(WithTimeout(5 * time.Second))
	if client.Timeout != 5*time.Second {
		t.Errorf("期望超时为5秒，实际为%v", client.Timeout)
	}

	// 测试零值或负值处理
	client = NewHTTPClient(WithTimeout(0))
	if client.Timeout != 30*time.Second {
		t.Errorf("无效超时值应被忽略，期望为30秒，实际为%v", client.Timeout)
	}

	client = NewHTTPClient(WithTimeout(-5 * time.Second))
	if client.Timeout != 30*time.Second {
		t.Errorf("负超时值应被忽略，期望为30秒，实际为%v", client.Timeout)
	}
}

// 测试代理设置
func TestWithProxyURL(t *testing.T) {
	proxyURL := "http://proxy.example.com:8080"
	client := NewHTTPClient(WithProxyURL(proxyURL))

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("无法获取transport配置")
	}

	if transport.Proxy == nil {
		t.Fatalf("代理函数未被设置")
	}

	// 创建测试请求以验证代理配置
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	url, err := transport.Proxy(req)
	if err != nil {
		t.Fatalf("代理函数执行错误: %v", err)
	}

	if url.String() != proxyURL {
		t.Errorf("代理URL应为%s，实际为%s", proxyURL, url.String())
	}

	// 测试空代理URL
	client = NewHTTPClient(WithProxyURL(""))
	transport, _ = client.Transport.(*http.Transport)

	// 默认的Transport在未设置代理时Proxy字段为nil，所以我们应该验证这一点
	// 当传递空字符串时，WithProxyURL应该什么都不做，保持默认nil值
	if transport.Proxy != nil {
		t.Errorf("空代理URL不应设置代理函数，Proxy应为nil")
	}
}

// 测试TLS配置
func TestWithTLSConfig(t *testing.T) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	client := NewHTTPClient(WithTLSConfig(tlsConfig))

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("无法获取transport配置")
	}

	if transport.TLSClientConfig == nil {
		t.Fatalf("TLS配置未被设置")
	}

	if transport.TLSClientConfig.InsecureSkipVerify != true {
		t.Errorf("InsecureSkipVerify应为true，实际为%v",
			transport.TLSClientConfig.InsecureSkipVerify)
	}

	if transport.TLSClientConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion应为TLS1.2，实际为%v",
			transport.TLSClientConfig.MinVersion)
	}
}

// 测试连接设置
func TestConnectionOptions(t *testing.T) {
	client := NewHTTPClient(
		WithMaxIdleConns(200),
		WithIdleConnTimeout(2*time.Minute),
		WithKeepAlive(45*time.Second),
		WithMaxConnsPerHost(20),
		WithDisableKeepAlives(true),
		WithDisableCompression(true),
	)

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("无法获取transport配置")
	}

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"MaxIdleConns", transport.MaxIdleConns, 200},
		{"IdleConnTimeout", transport.IdleConnTimeout, 2 * time.Minute},
		{"MaxConnsPerHost", transport.MaxConnsPerHost, 20},
		{"DisableKeepAlives", transport.DisableKeepAlives, true},
		{"DisableCompression", transport.DisableCompression, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.expected {
				t.Errorf("%s: 期望 %v, 实际 %v", test.name, test.expected, test.got)
			}
		})
	}
}

// 测试多选项组合
func TestMultipleOptions(t *testing.T) {
	jar, _ := cookiejar.New(nil)

	client := NewHTTPClient(
		WithTimeout(10*time.Second),
		WithMaxIdleConns(50),
		WithDisableKeepAlives(true),
		WithProxyURL("http://proxy.test.com:8888"),
		WithCookieJar(jar),
	)

	// 验证组合配置是否正确应用
	if client.Timeout != 10*time.Second {
		t.Errorf("超时应为10秒，实际为%v", client.Timeout)
	}

	if client.Jar != jar {
		t.Errorf("Cookie jar未正确设置")
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("无法获取transport配置")
	}

	if transport.MaxIdleConns != 50 {
		t.Errorf("MaxIdleConns应为50，实际为%d", transport.MaxIdleConns)
	}

	if !transport.DisableKeepAlives {
		t.Errorf("DisableKeepAlives应为true")
	}
}

// 测试重定向策略
func TestWithCheckRedirect(t *testing.T) {
	// 创建一个总是返回错误的重定向策略
	alwaysErrorPolicy := func(req *http.Request, via []*http.Request) error {
		return &url.Error{
			Op:  "Get",
			URL: req.URL.String(),
			Err: http.ErrUseLastResponse,
		}
	}

	client := NewHTTPClient(WithCheckRedirect(alwaysErrorPolicy))

	if client.CheckRedirect == nil {
		t.Fatalf("重定向策略未被设置")
	}
}

// 测试配置现有客户端
func TestConfigureHTTPClient(t *testing.T) {
	// 创建原始客户端
	original := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 配置现有客户端
	updated := ConfigureHTTPClient(original,
		WithTimeout(15*time.Second),
		WithMaxIdleConns(200),
	)

	// 验证是否返回同一个客户端实例
	if updated != original {
		t.Errorf("应返回同一个客户端实例")
	}

	// 验证配置是否已应用
	if updated.Timeout != 15*time.Second {
		t.Errorf("超时应为15秒，实际为%v", updated.Timeout)
	}

	transport, ok := updated.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("无法获取transport配置")
	}

	if transport.MaxIdleConns != 200 {
		t.Errorf("MaxIdleConns应为200，实际为%d", transport.MaxIdleConns)
	}
}

// 测试对空客户端的处理
func TestConfigureNilClient(t *testing.T) {
	client := ConfigureHTTPClient(nil, WithTimeout(25*time.Second))

	if client == nil {
		t.Fatalf("应创建新的客户端实例")
	}

	if client.Timeout != 25*time.Second {
		t.Errorf("超时应为25秒，实际为%v", client.Timeout)
	}
}

// 测试自定义Transport
func TestWithCustomTransport(t *testing.T) {
	customTransport := &http.Transport{
		MaxIdleConns:    300,
		IdleConnTimeout: 3 * time.Minute,
	}

	client := NewHTTPClient(WithCustomTransport(customTransport))

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("无法获取transport配置")
	}

	if transport != customTransport {
		t.Errorf("自定义Transport未被正确设置")
	}

	if transport.MaxIdleConns != 300 {
		t.Errorf("MaxIdleConns应为300，实际为%d", transport.MaxIdleConns)
	}

	if transport.IdleConnTimeout != 3*time.Minute {
		t.Errorf("IdleConnTimeout应为3分钟，实际为%v", transport.IdleConnTimeout)
	}
}

// 示例代码
func Example() {
	// 创建带有多个选项的HTTP客户端
	client := NewHTTPClient(
		WithTimeout(5*time.Second),
		WithMaxIdleConns(20),
		WithProxyURL("http://proxy.example.com:8080"),
		WithDisableCompression(true),
		WithRetry(3, 100*time.Millisecond, 2*time.Second),
	)

	// 使用客户端发起请求
	_, _ = client.Get("https://example.com")
}

// 示例：配置现有客户端
func ExampleConfigureHTTPClient() {
	// 假设这是从其他地方获取的客户端
	existingClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 配置现有客户端
	updatedClient := ConfigureHTTPClient(existingClient,
		WithMaxIdleConns(50),
		WithIdleConnTimeout(30*time.Second),
		WithDisableKeepAlives(true),
	)

	// 使用更新后的客户端
	_, _ = updatedClient.Get("https://example.com")
}
