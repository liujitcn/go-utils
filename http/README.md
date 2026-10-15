# http

http 是独立 Go 子模块，基于标准库 net/http 提供普通请求、流式响应和 Server-Sent Events（SSE）读取能力。实例客户端和包级默认客户端均可使用。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/http@latest
~~~

## 创建客户端

~~~go
client := http.NewClient(
	http.WithBaseURL("https://api.example.com"),
	http.WithDefaultHeaders(map[string]string{"X-App": "demo"}),
	http.WithTimeout(15*time.Second),
	http.WithMaxResponseBodyBytes(8 << 20),
)
~~~

默认 HTTP 客户端超时为 30 秒，普通请求日志默认开启。WithBaseURL 接收绝对地址，普通请求的相对 target 会基于该地址解析；target 也可以直接传绝对 URL。

客户端配置项：

- WithHTTPClient：替换底层标准库客户端，生命周期由调用方负责。
- WithTimeout：设置底层客户端总超时；非正数不生效。
- WithNoTimeout：关闭底层客户端总超时，适合长连接。
- WithMaxResponseBodyBytes：限制普通请求响应体读取大小，非正数表示不限制；超过限制返回 ErrResponseBodyTooLarge。
- WithRequestLog：开关普通请求的完整 URL 日志。
- WithDefaultHeader、WithDefaultHeaders：设置默认请求头。

## 普通请求

~~~go
var result struct {
	Message string `json:"message"`
}
err := client.Get("/ping", &result,
	http.WithQuery("lang", "zh-CN"),
	http.WithBearerToken(token),
)
~~~

支持 Do、Get、Post、Put、Patch、Delete、Head、Options。Get 等便捷方法会把非空响应体 JSON 解码到 result；result 为 nil 时只发送请求。Do 返回 Response，包含 Status、StatusCode、Header 和 Body，可用 DecodeJSON 或 String 读取。

单次请求配置项包括 WithContext、WithHeader、WithHeaders、WithQuery、WithQueries、WithJSONBody、WithFormBody、WithBodyBytes、WithBodyString、WithContentType 和 WithBearerToken。JSON 与表单请求体会自动设置对应 Content-Type。

## 流式响应

~~~go
stream, err := client.PostStream("/v1/stream", http.WithJSONBody(map[string]any{
	"stream": true,
}))
if err != nil {
	return err
}
defer stream.Close()

buf := make([]byte, 4096)
_, err = stream.Body.Read(buf)
~~~

DoStream、GetStream、PostStream、PutStream、PatchStream、DeleteStream 返回未读取的 StreamResponse。响应体由调用方读取并关闭；流式请求会复制底层客户端并将总超时设为 0，单次取消仍应通过 WithContext 传入上下文。

## SSE

~~~go
stream, err := client.GetSSE("/events", http.WithContext(ctx))
if err != nil {
	return err
}
defer stream.Close()

for {
	event, err := stream.Next()
	if errors.Is(err, io.EOF) {
		break
	}
	if err != nil {
		return err
	}
	if event.IsDone() {
		break
	}
	fmt.Println(event.Event, event.ID, event.Data)
}
~~~

DoSSE、GetSSE、PostSSE 会自动请求 text/event-stream。SSEEvent 提供 Event、Data、ID、Retry 字段；多行 data 会按 SSE 规范使用换行拼接，DecodeJSON 可直接解析 Data，IsDone 用于识别常见的 [DONE] 结束标记。

## 包级默认客户端

~~~go
http.Init(http.WithBaseURL("https://example.com"))
var result map[string]any
if err := http.Get("/health", &result); err != nil {
	return err
}
~~~

Init 只在包级默认客户端首次初始化时应用传入配置；Default 返回当前默认客户端。需要在初始化后替换实例时使用 SetDefaultClient。

## 使用边界

- 普通请求只在读取响应体时应用 MaxResponseBodyBytes；流式和 SSE 响应由调用方自行控制读取量。
- 默认日志包含最终完整 URL，查询参数可能含敏感信息；需要避免记录时使用 WithRequestLog(false)。
- 客户端不会自动按 HTTP 状态码返回错误，调用方应检查 Response.StatusCode 或流式响应的状态码。
- 所有返回的响应体都需要按对应模式处理，流式响应和 SSE 必须调用 Close。

## 测试

~~~bash
go test ./...
~~~
