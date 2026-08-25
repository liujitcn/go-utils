package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kratos/kratos/v3/log"
)

// Response HTTP 响应结果。
type Response struct {
	Status     string
	StatusCode int
	Header     http.Header
	Body       []byte
}

// Do 发送通用 HTTP 请求，并记录请求地址与失败错误。
func (c *Client) Do(method, target string, opts ...RequestOption) (*Response, error) {
	reqOptions, err := c.buildRequestOptions(opts...)
	if err != nil {
		log.Error("http: build request options failed", "method", method, "target", target, "error", err)
		return nil, err
	}

	var requestURL string
	requestURL, err = c.buildURL(target, reqOptions.query)
	if err != nil {
		log.Error("http: build request url failed", "target", target, "error", err)
		return nil, err
	}

	// 记录最终请求地址，便于排查基础地址与查询参数合并后的实际请求。
	log.Info("http: request", "method", method, "url", requestURL)

	var req *http.Request
	req, err = c.buildHTTPRequest(method, requestURL, reqOptions)
	if err != nil {
		log.Error("http: create request failed", "method", method, "url", requestURL, "error", err)
		return nil, err
	}

	var resp *http.Response
	resp, err = c.httpClient.Do(req)
	if err != nil {
		// 请求发送失败时记录最终地址和底层错误，便于定位网络或协议问题。
		log.Error("http: request failed", "method", method, "url", requestURL, "error", err)
		return nil, err
	}

	var body []byte
	body, err = readAndCloseBody(resp.Body)
	if err != nil {
		log.Error("http: read response body failed", "method", method, "url", requestURL, "error", err)
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
	return c.DoInto(http.MethodGet, target, result, opts...)
}

// Post 发送 POST 请求，并将响应体反序列化到目标对象。
func (c *Client) Post(target string, result any, opts ...RequestOption) error {
	return c.DoInto(http.MethodPost, target, result, opts...)
}

// Put 发送 PUT 请求，并将响应体反序列化到目标对象。
func (c *Client) Put(target string, result any, opts ...RequestOption) error {
	return c.DoInto(http.MethodPut, target, result, opts...)
}

// Patch 发送 PATCH 请求，并将响应体反序列化到目标对象。
func (c *Client) Patch(target string, result any, opts ...RequestOption) error {
	return c.DoInto(http.MethodPatch, target, result, opts...)
}

// Delete 发送 DELETE 请求，并将响应体反序列化到目标对象。
func (c *Client) Delete(target string, result any, opts ...RequestOption) error {
	return c.DoInto(http.MethodDelete, target, result, opts...)
}

// Head 发送 HEAD 请求，并将响应体反序列化到目标对象。
func (c *Client) Head(target string, result any, opts ...RequestOption) error {
	return c.DoInto(http.MethodHead, target, result, opts...)
}

// Options 发送 OPTIONS 请求，并将响应体反序列化到目标对象。
func (c *Client) Options(target string, result any, opts ...RequestOption) error {
	return c.DoInto(http.MethodOptions, target, result, opts...)
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
