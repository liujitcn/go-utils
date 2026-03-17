package http

import "sync"

var (
	defaultClientOnce sync.Once
	defaultClient     *Client
	defaultClientMu   sync.RWMutex
)

// Init 初始化包级默认 HTTP 客户端，首次调用生效。
func Init(opts ...ClientOption) *Client {
	defaultClientOnce.Do(func() {
		defaultClientMu.Lock()
		defer defaultClientMu.Unlock()
		defaultClient = NewClient(opts...)
	})
	return Default()
}

// Default 返回包级默认 HTTP 客户端。
func Default() *Client {
	defaultClientOnce.Do(func() {
		defaultClientMu.Lock()
		defer defaultClientMu.Unlock()
		defaultClient = NewClient()
	})

	defaultClientMu.RLock()
	defer defaultClientMu.RUnlock()
	return defaultClient
}

// SetDefaultClient 替换包级默认 HTTP 客户端。
func SetDefaultClient(client *Client) {
	defaultClientOnce.Do(func() {})
	if client == nil {
		client = NewClient()
	}

	defaultClientMu.Lock()
	defer defaultClientMu.Unlock()
	defaultClient = client
}

// Do 使用包级默认客户端发送通用 HTTP 请求。
func Do(method, target string, opts ...RequestOption) (*Response, error) {
	return Default().Do(method, target, opts...)
}

// Get 使用包级默认客户端发送 GET 请求，并将响应体反序列化到目标对象。
func Get(target string, result any, opts ...RequestOption) error {
	return Default().Get(target, result, opts...)
}

// Post 使用包级默认客户端发送 POST 请求，并将响应体反序列化到目标对象。
func Post(target string, result any, opts ...RequestOption) error {
	return Default().Post(target, result, opts...)
}

// Put 使用包级默认客户端发送 PUT 请求，并将响应体反序列化到目标对象。
func Put(target string, result any, opts ...RequestOption) error {
	return Default().Put(target, result, opts...)
}

// Patch 使用包级默认客户端发送 PATCH 请求，并将响应体反序列化到目标对象。
func Patch(target string, result any, opts ...RequestOption) error {
	return Default().Patch(target, result, opts...)
}

// Delete 使用包级默认客户端发送 DELETE 请求，并将响应体反序列化到目标对象。
func Delete(target string, result any, opts ...RequestOption) error {
	return Default().Delete(target, result, opts...)
}

// Head 使用包级默认客户端发送 HEAD 请求，并将响应体反序列化到目标对象。
func Head(target string, result any, opts ...RequestOption) error {
	return Default().Head(target, result, opts...)
}

// Options 使用包级默认客户端发送 OPTIONS 请求，并将响应体反序列化到目标对象。
func Options(target string, result any, opts ...RequestOption) error {
	return Default().Options(target, result, opts...)
}
