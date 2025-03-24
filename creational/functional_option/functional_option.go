// Package functionaloption 展示了Go中函数选项模式的实现
// 函数选项模式允许使用灵活、可扩展且易用的API创建配置复杂对象
package functionaloption

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// HTTPClientOptions 包含HTTP客户端的所有可配置选项
type HTTPClientOptions struct {
	Timeout            time.Duration                              // 请求超时时间
	KeepAlive          time.Duration                              // 保持连接的时间
	MaxIdleConns       int                                        // 最大空闲连接数
	IdleConnTimeout    time.Duration                              // 空闲连接超时
	TLSConfig          *tls.Config                                // TLS配置
	Transport          *http.Transport                            // 自定义传输配置
	Proxy              func(*http.Request) (*url.URL, error)      // 代理设置
	CheckRedirect      func(*http.Request, []*http.Request) error // 重定向策略
	Jar                http.CookieJar                             // Cookie处理
	MaxConnsPerHost    int                                        // 每个主机最大连接数
	DisableKeepAlives  bool                                       // 是否禁用长连接
	DisableCompression bool                                       // 是否禁用压缩
	RetryMax           int                                        // 最大重试次数
	RetryWaitMin       time.Duration                              // 重试最小等待时间
	RetryWaitMax       time.Duration                              // 重试最大等待时间
}

// defaultHTTPClientOptions 返回具有合理默认值的配置
func defaultHTTPClientOptions() HTTPClientOptions {
	return HTTPClientOptions{
		Timeout:           30 * time.Second,
		KeepAlive:         30 * time.Second,
		MaxIdleConns:      100,
		IdleConnTimeout:   90 * time.Second,
		MaxConnsPerHost:   10,
		DisableKeepAlives: false,
		RetryMax:          0, // 默认不重试
	}
}

// Option 定义修改HTTPClientOptions的函数类型
type Option func(*HTTPClientOptions)

// WithTimeout 设置HTTP请求的超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *HTTPClientOptions) {
		if timeout > 0 {
			o.Timeout = timeout
		}
	}
}

// WithKeepAlive 设置TCP保持连接的时间
func WithKeepAlive(keepAlive time.Duration) Option {
	return func(o *HTTPClientOptions) {
		if keepAlive >= 0 {
			o.KeepAlive = keepAlive
		}
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(o *HTTPClientOptions) {
		if maxIdleConns >= 0 {
			o.MaxIdleConns = maxIdleConns
		}
	}
}

// WithIdleConnTimeout 设置空闲连接超时时间
func WithIdleConnTimeout(idleConnTimeout time.Duration) Option {
	return func(o *HTTPClientOptions) {
		if idleConnTimeout > 0 {
			o.IdleConnTimeout = idleConnTimeout
		}
	}
}

// WithTLSConfig 设置TLS配置
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(o *HTTPClientOptions) {
		if tlsConfig != nil {
			o.TLSConfig = tlsConfig
		}
	}
}

// WithProxy 设置代理
func WithProxy(proxy func(*http.Request) (*url.URL, error)) Option {
	return func(o *HTTPClientOptions) {
		if proxy != nil {
			o.Proxy = proxy
		}
	}
}

// WithProxyURL 通过URL字符串设置代理
func WithProxyURL(proxyURL string) Option {
	return func(o *HTTPClientOptions) {
		if proxyURL != "" {
			proxyFunc := func(_ *http.Request) (*url.URL, error) {
				return url.Parse(proxyURL)
			}
			o.Proxy = proxyFunc
		}
	}
}

// WithCheckRedirect 设置重定向策略
func WithCheckRedirect(checkRedirect func(*http.Request, []*http.Request) error) Option {
	return func(o *HTTPClientOptions) {
		if checkRedirect != nil {
			o.CheckRedirect = checkRedirect
		}
	}
}

// WithCookieJar 设置Cookie处理
func WithCookieJar(jar http.CookieJar) Option {
	return func(o *HTTPClientOptions) {
		if jar != nil {
			o.Jar = jar
		}
	}
}

// WithMaxConnsPerHost 设置每个主机的最大连接数
func WithMaxConnsPerHost(maxConnsPerHost int) Option {
	return func(o *HTTPClientOptions) {
		if maxConnsPerHost > 0 {
			o.MaxConnsPerHost = maxConnsPerHost
		}
	}
}

// WithDisableKeepAlives 设置是否禁用长连接
func WithDisableKeepAlives(disable bool) Option {
	return func(o *HTTPClientOptions) {
		o.DisableKeepAlives = disable
	}
}

// WithDisableCompression 设置是否禁用压缩
func WithDisableCompression(disable bool) Option {
	return func(o *HTTPClientOptions) {
		o.DisableCompression = disable
	}
}

// WithRetry 配置重试策略
func WithRetry(maxRetries int, minWait, maxWait time.Duration) Option {
	return func(o *HTTPClientOptions) {
		if maxRetries > 0 {
			o.RetryMax = maxRetries
		}
		if minWait > 0 {
			o.RetryWaitMin = minWait
		}
		if maxWait > 0 {
			o.RetryWaitMax = maxWait
		}
	}
}

// WithCustomTransport 设置自定义传输配置
func WithCustomTransport(transport *http.Transport) Option {
	return func(o *HTTPClientOptions) {
		if transport != nil {
			o.Transport = transport
		}
	}
}

// NewHTTPClient 使用功能选项模式创建并配置HTTP客户端
//
// 示例:
//
//	client := NewHTTPClient(
//	    WithTimeout(5 * time.Second),
//	    WithMaxIdleConns(20),
//	    WithProxyURL("http://proxy.example.com:8080"),
//	    WithRetry(3, 100*time.Millisecond, 2*time.Second),
//	)
func NewHTTPClient(opts ...Option) *http.Client {
	// 使用默认选项作为起点
	options := defaultHTTPClientOptions()

	// 应用所有提供的选项
	for _, opt := range opts {
		opt(&options)
	}

	// 创建传输配置
	transport := options.Transport
	if transport == nil {
		transport = &http.Transport{
			Proxy: options.Proxy,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // 连接超时
				KeepAlive: options.KeepAlive,
			}).DialContext,
			MaxIdleConns:           options.MaxIdleConns,
			IdleConnTimeout:        options.IdleConnTimeout,
			TLSClientConfig:        options.TLSConfig,
			MaxConnsPerHost:        options.MaxConnsPerHost,
			DisableKeepAlives:      options.DisableKeepAlives,
			DisableCompression:     options.DisableCompression,
			MaxResponseHeaderBytes: 1 << 20, // 1 MB
		}
	}

	// 创建HTTP客户端
	client := &http.Client{
		Transport:     transport,
		CheckRedirect: options.CheckRedirect,
		Jar:           options.Jar,
		Timeout:       options.Timeout,
	}

	// 如果配置了重试，可以添加重试逻辑
	// 注意：这里只是一个示意，实际的重试可能需要更复杂的逻辑
	if options.RetryMax > 0 {
		// 真实实现中，这里应该替换为实际的重试包装逻辑
		// 这可能涉及到使用自定义RoundTripper或请求拦截器
	}

	return client
}

// ConfigureHTTPClient 使用选项配置现有的HTTP客户端
// 这对于需要修改但不想完全替换的客户端很有用
func ConfigureHTTPClient(client *http.Client, opts ...Option) *http.Client {
	if client == nil {
		return NewHTTPClient(opts...)
	}

	// 从现有客户端提取选项
	options := defaultHTTPClientOptions()

	if client.Timeout > 0 {
		options.Timeout = client.Timeout
	}

	if transport, ok := client.Transport.(*http.Transport); ok {
		if transport.IdleConnTimeout > 0 {
			options.IdleConnTimeout = transport.IdleConnTimeout
		}
		if transport.MaxIdleConns > 0 {
			options.MaxIdleConns = transport.MaxIdleConns
		}
		options.DisableKeepAlives = transport.DisableKeepAlives
		options.DisableCompression = transport.DisableCompression
		options.MaxConnsPerHost = transport.MaxConnsPerHost
		options.Proxy = transport.Proxy
		options.TLSConfig = transport.TLSClientConfig
	}

	options.CheckRedirect = client.CheckRedirect
	options.Jar = client.Jar

	// 应用新选项
	for _, opt := range opts {
		opt(&options)
	}

	// 更新客户端配置
	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.MaxIdleConns = options.MaxIdleConns
		transport.IdleConnTimeout = options.IdleConnTimeout
		transport.DisableKeepAlives = options.DisableKeepAlives
		transport.DisableCompression = options.DisableCompression
		transport.MaxConnsPerHost = options.MaxConnsPerHost
		transport.Proxy = options.Proxy
		transport.TLSClientConfig = options.TLSConfig

		// 创建一个新的DialContext以应用KeepAlive设置
		transport.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: options.KeepAlive,
		}).DialContext
	} else {
		client.Transport = &http.Transport{
			Proxy: options.Proxy,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: options.KeepAlive,
			}).DialContext,
			MaxIdleConns:           options.MaxIdleConns,
			IdleConnTimeout:        options.IdleConnTimeout,
			TLSClientConfig:        options.TLSConfig,
			MaxConnsPerHost:        options.MaxConnsPerHost,
			DisableKeepAlives:      options.DisableKeepAlives,
			DisableCompression:     options.DisableCompression,
			MaxResponseHeaderBytes: 1 << 20,
		}
	}

	client.Timeout = options.Timeout
	client.CheckRedirect = options.CheckRedirect
	client.Jar = options.Jar

	return client
}
