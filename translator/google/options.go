package google

import (
	"net/http"

	"google.golang.org/api/option"
)

// Option 配置 Google 翻译器。
type Option func(*options)

type options struct {
	version       string
	apiKey        string
	projectID     string
	location      string
	parent        string
	httpClient    *http.Client
	v1Endpoint    string
	clientOptions []option.ClientOption
}

// WithVersion 设置 Google 翻译版本，可选值为 v1、v2、v3。
func WithVersion(version string) Option {
	return func(options *options) {
		options.version = version
	}
}

// WithAPIKey 设置 Google Cloud API Key，V2 和 V3 未设置时使用 ADC。
func WithAPIKey(apiKey string) Option {
	return func(options *options) {
		options.apiKey = apiKey
	}
}

// WithApiKey 设置 Google Cloud API Key。
// Deprecated: 请使用 WithAPIKey。
func WithApiKey(apiKey string) Option {
	return WithAPIKey(apiKey)
}

// WithProjectID 设置 Google Cloud 项目 ID，V3 必填。
func WithProjectID(projectID string) Option {
	return func(options *options) {
		options.projectID = projectID
	}
}

// WithLocation 设置 Google Cloud V3 区域，默认使用 global。
func WithLocation(location string) Option {
	return func(options *options) {
		options.location = location
	}
}

// WithParent 直接设置 Google Cloud V3 Parent，设置后优先于项目和区域。
func WithParent(parent string) Option {
	return func(options *options) {
		options.parent = parent
	}
}

// WithHTTPClient 设置 Google V1 请求客户端，调用方负责其生命周期。
func WithHTTPClient(client *http.Client) Option {
	return func(options *options) {
		options.httpClient = client
	}
}

// WithClientOptions 追加 Google V2 和 V3 官方客户端配置。
func WithClientOptions(clientOptions ...option.ClientOption) Option {
	return func(options *options) {
		options.clientOptions = append(options.clientOptions, clientOptions...)
	}
}
