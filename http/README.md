# http

`http` 提供基于 options 模式封装的通用 HTTP 请求客户端，支持实例模式与包级单例静态调用，并在方法内部直接将响应体反序列化到结果对象。

## 安装

```bash
go get github.com/liujitcn/go-utils/http@latest
```

## 功能

- 支持 `Do`、`Get`、`Post`、`Put`、`Patch`、`Delete`、`Head`、`Options`
- 支持 `DoStream`、`GetStream`、`PostStream` 等 HTTP Stream 请求
- 支持 `DoSSE`、`GetSSE`、`PostSSE` 解析 `text/event-stream`
- 支持 `WithBaseURL`、`WithTimeout`、`WithNoTimeout`、`WithDefaultHeader`
- 支持 `WithHeader`、`WithQuery`、`WithJSONBody`、`WithFormBody`
- `Do` 会通过 Kratos log 记录最终请求地址，请求失败时会记录对应错误信息
- `DoStream` 与 SSE 请求会自动关闭客户端总超时，长连接取消建议使用 `WithContext`
- 支持包级单例初始化与静态请求方法

## 实例模式

```go
package main

import (
	"fmt"

	httputil "github.com/liujitcn/go-utils/http"
)

func main() {
	type UserResponse struct {
		Name string `json:"name"`
	}

	var result UserResponse

	client := httputil.NewClient(
		httputil.WithBaseURL("https://example.com"),
		httputil.WithDefaultHeader("X-App", "shop"),
	)

	err := client.Get("/user", &result, httputil.WithQuery("lang", "zh-CN"))
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Name)
}
```

## HTTP Stream

```go
package main

import (
	"bufio"
	"fmt"

	httputil "github.com/liujitcn/go-utils/http"
)

func main() {
	client := httputil.NewClient(
		httputil.WithBaseURL("https://api.example.com"),
	)

	stream, err := client.PostStream(
		"/v1/chat/completions",
		httputil.WithJSONBody(map[string]any{
			"stream": true,
		}),
	)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	reader := bufio.NewReader(stream.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(line)
	}
}
```

## SSE

```go
package main

import (
	"fmt"
	"io"

	httputil "github.com/liujitcn/go-utils/http"
)

func main() {
	client := httputil.NewClient(
		httputil.WithBaseURL("https://api.example.com"),
	)

	stream, err := client.PostSSE(
		"/v1/chat/completions",
		httputil.WithJSONBody(map[string]any{
			"stream": true,
		}),
		httputil.WithBearerToken("token"),
	)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	for {
		event, err := stream.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if event.IsDone() {
			break
		}
		fmt.Println(event.Data)
	}
}
```

## 单例静态模式

```go
package main

import (
	"fmt"

	httputil "github.com/liujitcn/go-utils/http"
)

func main() {
	type WxAccessToken struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	var token WxAccessToken

	httputil.Init(
		httputil.WithBaseURL("https://example.com"),
		httputil.WithDefaultHeader("X-App", "shop"),
	)

	err := httputil.Get("/cgi-bin/token", &token, httputil.WithQuery("grant_type", "client_credential"))
	if err != nil {
		panic(err)
	}

	fmt.Println(token.AccessToken)
}
```

## 测试

```bash
go test ./...
```
