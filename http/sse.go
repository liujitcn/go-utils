package http

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultSSEEventName = "message"

// SSEEvent 表示一条 Server-Sent Events 消息。
type SSEEvent struct {
	Event string
	Data  string
	ID    string
	Retry time.Duration
}

// SSEStream 表示正在读取中的 Server-Sent Events 响应流。
type SSEStream struct {
	*StreamResponse

	reader      *bufio.Reader
	lastEventID string
	builder     sseEventBuilder
}

type sseEventBuilder struct {
	event    string
	data     []string
	retry    time.Duration
	hasRetry bool
}

// DoSSE 发送通用 HTTP 请求，并返回 SSE 事件流读取器。
func (c *Client) DoSSE(method, target string, opts ...RequestOption) (*SSEStream, error) {
	streamOptions := make([]RequestOption, 0, len(opts)+1)
	streamOptions = append(streamOptions, WithHeader("Accept", "text/event-stream"))
	streamOptions = append(streamOptions, opts...)

	response, err := c.DoStream(method, target, streamOptions...)
	if err != nil {
		return nil, err
	}

	return NewSSEStream(response), nil
}

// GetSSE 发送 GET 请求，并返回 SSE 事件流读取器。
func (c *Client) GetSSE(target string, opts ...RequestOption) (*SSEStream, error) {
	return c.DoSSE(http.MethodGet, target, opts...)
}

// PostSSE 发送 POST 请求，并返回 SSE 事件流读取器。
func (c *Client) PostSSE(target string, opts ...RequestOption) (*SSEStream, error) {
	return c.DoSSE(http.MethodPost, target, opts...)
}

// NewSSEStream 基于原始 HTTP 流式响应创建 SSE 事件流读取器。
func NewSSEStream(response *StreamResponse) *SSEStream {
	stream := &SSEStream{
		StreamResponse: response,
	}
	if response != nil && response.Body != nil {
		stream.reader = bufio.NewReader(response.Body)
	}
	return stream
}

// Next 读取下一条 SSE 消息；流结束时返回 io.EOF。
func (s *SSEStream) Next() (*SSEEvent, error) {
	if s == nil || s.StreamResponse == nil || s.Body == nil || s.reader == nil {
		return nil, fmt.Errorf("http: sse stream is nil")
	}

	for {
		line, err := s.readLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				event := s.builder.build(s.lastEventID)
				if event != nil {
					return event, nil
				}
			}
			return nil, err
		}

		if line == "" {
			event := s.builder.build(s.lastEventID)
			if event == nil {
				s.builder.reset()
				continue
			}
			return event, nil
		}

		s.applyLine(line)
	}
}

// Close 关闭 SSE 响应流。
func (s *SSEStream) Close() error {
	if s == nil || s.StreamResponse == nil {
		return nil
	}
	return s.StreamResponse.Close()
}

// DecodeJSON 将 SSE data 字段反序列化到目标对象。
func (e *SSEEvent) DecodeJSON(target any) error {
	if target == nil {
		return fmt.Errorf("http: decode target is nil")
	}
	if e == nil || e.Data == "" {
		return fmt.Errorf("http: sse event data is empty")
	}
	return json.Unmarshal([]byte(e.Data), target)
}

// IsDone 判断当前 SSE data 是否为 AI 流式接口常见的结束标记。
func (e *SSEEvent) IsDone() bool {
	if e == nil {
		return false
	}
	return strings.TrimSpace(e.Data) == "[DONE]"
}

// readLine 读取并归一化单行 SSE 内容。
func (s *SSEStream) readLine() (string, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && line != "" {
			return trimLineEnding(line), nil
		}
		return "", err
	}
	return trimLineEnding(line), nil
}

// applyLine 解析单行 SSE 字段并更新当前事件构建器。
func (s *SSEStream) applyLine(line string) {
	if strings.HasPrefix(line, ":") {
		return
	}

	field, value, hasValue := strings.Cut(line, ":")
	if hasValue && strings.HasPrefix(value, " ") {
		value = strings.TrimPrefix(value, " ")
	}

	// 根据 SSE 字段语义更新事件数据，未知字段按协议忽略。
	switch field {
	case "event":
		s.builder.event = value
	case "data":
		s.builder.data = append(s.builder.data, value)
	case "id":
		if !strings.Contains(value, "\x00") {
			s.lastEventID = value
		}
	case "retry":
		retry, ok := parseSSERetry(value)
		if ok {
			s.builder.retry = retry
			s.builder.hasRetry = true
		}
	}
}

// build 构建 SSE 消息；没有 data 字段时按协议不派发事件。
func (b *sseEventBuilder) build(lastEventID string) *SSEEvent {
	if len(b.data) == 0 {
		return nil
	}

	eventName := b.event
	if eventName == "" {
		eventName = defaultSSEEventName
	}

	event := &SSEEvent{
		Event: eventName,
		Data:  strings.Join(b.data, "\n"),
		ID:    lastEventID,
	}
	if b.hasRetry {
		event.Retry = b.retry
	}
	b.reset()
	return event
}

// reset 重置事件构建器，准备读取下一条 SSE 消息。
func (b *sseEventBuilder) reset() {
	b.event = ""
	b.data = nil
	b.retry = 0
	b.hasRetry = false
}

// parseSSERetry 解析 SSE retry 字段，单位为毫秒。
func parseSSERetry(value string) (time.Duration, bool) {
	retryMilliseconds, err := strconv.Atoi(value)
	if err != nil || retryMilliseconds < 0 {
		return 0, false
	}
	return time.Duration(retryMilliseconds) * time.Millisecond, true
}

// trimLineEnding 去除 SSE 行尾的换行符。
func trimLineEnding(line string) string {
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return line
}
