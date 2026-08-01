package volc

import "net/http"

// Option 配置火山引擎翻译器。
type Option func(*Translator)

// WithRegion 设置火山引擎区域，默认使用 cn-north-1。
func WithRegion(region string) Option {
	return func(translator *Translator) {
		translator.region = region
	}
}

// WithHTTPClient 设置请求客户端，调用方负责其生命周期。
func WithHTTPClient(client *http.Client) Option {
	return func(translator *Translator) {
		translator.httpClient = client
	}
}
