package http

import (
	"context"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"net/url"
	"time"
)

// ClientOption HTTP 客户端配置项。
type ClientOption func(*Client)

// RequestOption HTTP 请求配置项。
type RequestOption func(*requestOptions) error

// requestOptions 单次请求配置。
type requestOptions struct {
	context context.Context
	headers stdhttp.Header
	query   url.Values
	body    []byte
}

// WithHTTPClient 设置自定义 HTTP 客户端。
func WithHTTPClient(client *stdhttp.Client) ClientOption {
	return func(c *Client) {
		if client == nil {
			return
		}
		c.httpClient = client
	}
}

// WithTimeout 设置 HTTP 客户端超时时间。
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if timeout <= 0 {
			return
		}
		c.httpClient.Timeout = timeout
	}
}

// WithBaseURL 设置请求基础地址。
func WithBaseURL(raw string) ClientOption {
	return func(c *Client) {
		if raw == "" {
			return
		}
		parsed, err := url.Parse(raw)
		if err != nil {
			return
		}
		c.baseURL = parsed
	}
}

// WithDefaultHeader 设置客户端默认请求头。
func WithDefaultHeader(key, value string) ClientOption {
	return func(c *Client) {
		if key == "" {
			return
		}
		c.defaultHeaders.Set(key, value)
	}
}

// WithDefaultHeaders 批量设置客户端默认请求头。
func WithDefaultHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for key, value := range headers {
			if key == "" {
				continue
			}
			c.defaultHeaders.Set(key, value)
		}
	}
}

// WithContext 设置单次请求上下文。
func WithContext(ctx context.Context) RequestOption {
	return func(opts *requestOptions) error {
		if ctx == nil {
			return nil
		}
		opts.context = ctx
		return nil
	}
}

// WithHeader 设置单次请求头。
func WithHeader(key, value string) RequestOption {
	return func(opts *requestOptions) error {
		if key == "" {
			return nil
		}
		opts.headers.Set(key, value)
		return nil
	}
}

// WithHeaders 批量设置单次请求头。
func WithHeaders(headers map[string]string) RequestOption {
	return func(opts *requestOptions) error {
		for key, value := range headers {
			if key == "" {
				continue
			}
			opts.headers.Set(key, value)
		}
		return nil
	}
}

// WithQuery 设置单个查询参数。
func WithQuery(key, value string) RequestOption {
	return func(opts *requestOptions) error {
		if key == "" {
			return nil
		}
		opts.query.Add(key, value)
		return nil
	}
}

// WithQueries 批量设置查询参数。
func WithQueries(values map[string]string) RequestOption {
	return func(opts *requestOptions) error {
		for key, value := range values {
			if key == "" {
				continue
			}
			opts.query.Add(key, value)
		}
		return nil
	}
}

// WithBodyBytes 设置字节请求体。
func WithBodyBytes(body []byte) RequestOption {
	return func(opts *requestOptions) error {
		if body == nil {
			opts.body = nil
			return nil
		}
		opts.body = append([]byte(nil), body...)
		return nil
	}
}

// WithBodyString 设置字符串请求体。
func WithBodyString(body string) RequestOption {
	return func(opts *requestOptions) error {
		opts.body = []byte(body)
		return nil
	}
}

// WithJSONBody 设置 JSON 请求体，并自动写入 Content-Type。
func WithJSONBody(body any) RequestOption {
	return func(opts *requestOptions) error {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		opts.body = data
		opts.headers.Set("Content-Type", "application/json")
		return nil
	}
}

// WithFormBody 设置表单请求体，并自动写入 Content-Type。
func WithFormBody(values url.Values) RequestOption {
	return func(opts *requestOptions) error {
		if values == nil {
			return nil
		}
		opts.body = []byte(values.Encode())
		opts.headers.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	}
}

// WithContentType 设置请求体类型。
func WithContentType(contentType string) RequestOption {
	return func(opts *requestOptions) error {
		if contentType == "" {
			return nil
		}
		opts.headers.Set("Content-Type", contentType)
		return nil
	}
}

// WithBearerToken 设置 Bearer 鉴权请求头。
func WithBearerToken(token string) RequestOption {
	return func(opts *requestOptions) error {
		if token == "" {
			return fmt.Errorf("http: bearer token is empty")
		}
		opts.headers.Set("Authorization", "Bearer "+token)
		return nil
	}
}
