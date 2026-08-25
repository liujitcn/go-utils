package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client HTTP 通用请求客户端。
type Client struct {
	httpClient     *http.Client
	baseURL        *url.URL
	defaultHeaders http.Header
}

// NewClient 创建 HTTP 通用请求客户端。
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		httpClient:     &http.Client{Timeout: 30 * time.Second},
		defaultHeaders: make(http.Header),
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(client)
	}
	return client
}

// buildRequestOptions 构建单次请求配置。
func (c *Client) buildRequestOptions(opts ...RequestOption) (*requestOptions, error) {
	reqOptions := defaultRequestOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(reqOptions); err != nil {
			return nil, err
		}
	}
	return reqOptions, nil
}

// buildHTTPRequest 根据请求配置创建标准库请求对象。
func (c *Client) buildHTTPRequest(method, requestURL string, reqOptions *requestOptions) (*http.Request, error) {
	req, err := http.NewRequestWithContext(reqOptions.context, method, requestURL, bytes.NewReader(reqOptions.body))
	if err != nil {
		return nil, err
	}
	req.Header = cloneHeader(c.defaultHeaders)
	mergeHeaders(req.Header, reqOptions.headers)
	return req, nil
}

// buildURL 构建最终请求地址，并合并查询参数。
func (c *Client) buildURL(target string, query url.Values) (string, error) {
	if strings.TrimSpace(target) == "" {
		if c.baseURL == nil {
			return "", fmt.Errorf("http: request url is empty")
		}
	}

	var targetURL *url.URL
	var err error
	if strings.TrimSpace(target) == "" {
		targetURL = cloneURL(c.baseURL)
	} else {
		targetURL, err = url.Parse(target)
		if err != nil {
			return "", err
		}
	}

	if c.baseURL != nil && !targetURL.IsAbs() {
		targetURL = c.baseURL.ResolveReference(targetURL)
	}
	if targetURL == nil {
		return "", fmt.Errorf("http: request url is invalid")
	}

	values := targetURL.Query()
	for key, items := range query {
		for _, item := range items {
			values.Add(key, item)
		}
	}
	targetURL.RawQuery = values.Encode()
	return targetURL.String(), nil
}

// defaultRequestOptions 返回请求默认配置。
func defaultRequestOptions() *requestOptions {
	return &requestOptions{
		context: context.Background(),
		headers: make(http.Header),
		query:   make(url.Values),
	}
}

// cloneHeader 复制请求头，避免共享底层数据。
func cloneHeader(header http.Header) http.Header {
	if header == nil {
		return make(http.Header)
	}
	return header.Clone()
}

// cloneURL 复制 URL，避免修改客户端默认配置。
func cloneURL(raw *url.URL) *url.URL {
	if raw == nil {
		return nil
	}
	cloned := new(url.URL)
	*cloned = *raw
	return cloned
}

// mergeHeaders 合并请求头，后写入的值覆盖前面的值。
func mergeHeaders(dst, src http.Header) {
	for key, values := range src {
		if len(values) == 0 {
			continue
		}
		dst.Del(key)
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
