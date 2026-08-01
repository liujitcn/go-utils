package baidu

import "net/http"

// Option 配置百度翻译器。
type Option func(*Translator)

// WithHTTPClient 设置请求客户端，调用方负责其生命周期。
func WithHTTPClient(client *http.Client) Option {
	return func(translator *Translator) {
		translator.client = client
	}
}
