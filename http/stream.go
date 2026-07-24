package http

import (
	"io"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v3/log"
)

// StreamResponse HTTP 流式响应结果。
type StreamResponse struct {
	Status     string
	StatusCode int
	Header     stdhttp.Header
	Body       io.ReadCloser
}

// DoStream 发送通用 HTTP 请求，并返回未读取的响应流。
func (c *Client) DoStream(method, target string, opts ...RequestOption) (*StreamResponse, error) {
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

	// 流式请求也记录最终地址，方便排查 SSE 或 chunk stream 的真实请求参数。
	log.Info("http: stream request", "method", method, "url", requestURL)

	var req *stdhttp.Request
	req, err = c.buildHTTPRequest(method, requestURL, reqOptions)
	if err != nil {
		log.Error("http: create stream request failed", "method", method, "url", requestURL, "error", err)
		return nil, err
	}

	var resp *stdhttp.Response
	resp, err = c.streamHTTPClient().Do(req)
	if err != nil {
		log.Error("http: stream request failed", "method", method, "url", requestURL, "error", err)
		return nil, err
	}

	return &StreamResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     resp.Header.Clone(),
		Body:       resp.Body,
	}, nil
}

// GetStream 发送 GET 请求，并返回未读取的响应流。
func (c *Client) GetStream(target string, opts ...RequestOption) (*StreamResponse, error) {
	return c.DoStream(stdhttp.MethodGet, target, opts...)
}

// PostStream 发送 POST 请求，并返回未读取的响应流。
func (c *Client) PostStream(target string, opts ...RequestOption) (*StreamResponse, error) {
	return c.DoStream(stdhttp.MethodPost, target, opts...)
}

// PutStream 发送 PUT 请求，并返回未读取的响应流。
func (c *Client) PutStream(target string, opts ...RequestOption) (*StreamResponse, error) {
	return c.DoStream(stdhttp.MethodPut, target, opts...)
}

// PatchStream 发送 PATCH 请求，并返回未读取的响应流。
func (c *Client) PatchStream(target string, opts ...RequestOption) (*StreamResponse, error) {
	return c.DoStream(stdhttp.MethodPatch, target, opts...)
}

// DeleteStream 发送 DELETE 请求，并返回未读取的响应流。
func (c *Client) DeleteStream(target string, opts ...RequestOption) (*StreamResponse, error) {
	return c.DoStream(stdhttp.MethodDelete, target, opts...)
}

// Close 关闭 HTTP 响应流。
func (r *StreamResponse) Close() error {
	if r == nil || r.Body == nil {
		return nil
	}
	return r.Body.Close()
}

// streamHTTPClient 返回关闭总超时的客户端副本，避免长连接被默认 Timeout 截断。
func (c *Client) streamHTTPClient() *stdhttp.Client {
	if c == nil || c.httpClient == nil {
		return &stdhttp.Client{}
	}
	client := *c.httpClient
	client.Timeout = 0
	return &client
}
