# http

`http` 提供基于 options 模式封装的通用 HTTP 请求客户端，支持实例模式与包级单例静态调用，并在方法内部直接将响应体反序列化到结果对象。

## 安装

```bash
go get github.com/liujitcn/go-utils/http@latest
```

## 功能

- 支持 `Do`、`Get`、`Post`、`Put`、`Patch`、`Delete`、`Head`、`Options`
- 支持 `WithBaseURL`、`WithTimeout`、`WithDefaultHeader`
- 支持 `WithHeader`、`WithQuery`、`WithJSONBody`、`WithFormBody`
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
