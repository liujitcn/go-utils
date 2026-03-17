package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdhttp "net/http"
	"net/url"
	"strings"
	"time"
)

// Client HTTP 通用请求客户端。
type Client struct {
	httpClient     *stdhttp.Client
	baseURL        *url.URL
	defaultHeaders stdhttp.Header
}

// Response HTTP 响应结果。
type Response struct {
	Status     string
	StatusCode int
	Header     stdhttp.Header
	Body       []byte
}

// NewClient 创建 HTTP 通用请求客户端。
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		httpClient:     &stdhttp.Client{Timeout: 30 * time.Second},
		defaultHeaders: make(stdhttp.Header),
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(client)
	}
	return client
}

// Do 发送通用 HTTP 请求。
func (c *Client) Do(method, target string, opts ...RequestOption) (*Response, error) {
	reqOptions, err := c.buildRequestOptions(opts...)
	if err != nil {
		return nil, err
	}

	var requestURL string
	requestURL, err = c.buildURL(target, reqOptions.query)
	if err != nil {
		return nil, err
	}

	var req *stdhttp.Request
	req, err = stdhttp.NewRequestWithContext(reqOptions.context, method, requestURL, bytes.NewReader(reqOptions.body))
	if err != nil {
		return nil, err
	}
	req.Header = cloneHeader(c.defaultHeaders)
	mergeHeaders(req.Header, reqOptions.headers)

	var resp *stdhttp.Response
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body []byte
	body, err = readAndCloseBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     resp.Header.Clone(),
		Body:       body,
	}, nil
}

// Get 发送 GET 请求，并将响应体反序列化到目标对象。
func (c *Client) Get(target string, result any, opts ...RequestOption) error {
	return c.DoInto(stdhttp.MethodGet, target, result, opts...)
}

// Post 发送 POST 请求，并将响应体反序列化到目标对象。
func (c *Client) Post(target string, result any, opts ...RequestOption) error {
	return c.DoInto(stdhttp.MethodPost, target, result, opts...)
}

// Put 发送 PUT 请求，并将响应体反序列化到目标对象。
func (c *Client) Put(target string, result any, opts ...RequestOption) error {
	return c.DoInto(stdhttp.MethodPut, target, result, opts...)
}

// Patch 发送 PATCH 请求，并将响应体反序列化到目标对象。
func (c *Client) Patch(target string, result any, opts ...RequestOption) error {
	return c.DoInto(stdhttp.MethodPatch, target, result, opts...)
}

// Delete 发送 DELETE 请求，并将响应体反序列化到目标对象。
func (c *Client) Delete(target string, result any, opts ...RequestOption) error {
	return c.DoInto(stdhttp.MethodDelete, target, result, opts...)
}

// Head 发送 HEAD 请求，并将响应体反序列化到目标对象。
func (c *Client) Head(target string, result any, opts ...RequestOption) error {
	return c.DoInto(stdhttp.MethodHead, target, result, opts...)
}

// Options 发送 OPTIONS 请求，并将响应体反序列化到目标对象。
func (c *Client) Options(target string, result any, opts ...RequestOption) error {
	return c.DoInto(stdhttp.MethodOptions, target, result, opts...)
}

// DecodeJSON 将响应体反序列化到目标对象。
func (r *Response) DecodeJSON(target any) error {
	if target == nil {
		return fmt.Errorf("http: decode target is nil")
	}
	if len(r.Body) == 0 {
		return fmt.Errorf("http: response body is empty")
	}
	return json.Unmarshal(r.Body, target)
}

// String 返回响应体字符串。
func (r *Response) String() string {
	return string(r.Body)
}

// DoInto 发送请求，并将响应体反序列化到目标对象。
func (c *Client) DoInto(method, target string, result any, opts ...RequestOption) error {
	var response *Response
	var err error
	response, err = c.Do(method, target, opts...)
	if err != nil {
		return err
	}
	if result == nil || len(response.Body) == 0 {
		return nil
	}
	return response.DecodeJSON(result)
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
		headers: make(stdhttp.Header),
		query:   make(url.Values),
	}
}

// cloneHeader 复制请求头，避免共享底层数据。
func cloneHeader(header stdhttp.Header) stdhttp.Header {
	if header == nil {
		return make(stdhttp.Header)
	}
	return header.Clone()
}

// readAndCloseBody 读取响应体内容，并确保在所有路径中关闭响应流。
func readAndCloseBody(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	// 统一在辅助方法内关闭响应体，避免调用方遗漏关闭导致资源泄漏告警。
	defer func() {
		_ = body.Close()
	}()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// cloneURL 复制 URL，避免修改客户端默认配置。
func cloneURL(raw *url.URL) *url.URL {
	if raw == nil {
		return nil
	}
	cloned := *raw
	return &cloned
}

// mergeHeaders 合并请求头，后写入的值覆盖前面的值。
func mergeHeaders(dst, src stdhttp.Header) {
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
