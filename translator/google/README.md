# Google Translator

Google Provider 同时保留 V1、V2 和 V3 三种模式。

## 安装

```bash
go get github.com/liujitcn/go-utils/translator/google@latest
```

## V1

```go
client := google.NewTranslator(google.WithVersion("v1"))
result, err := client.Translate(ctx, "Hello", "en", "zh-CN")
```

V1 使用非官方 `client=gtx` 接口，无需凭证，但没有 SLA，只适合开发、演示或可接受失败的低风险任务。

## V2

```go
client := google.NewTranslator(
	google.WithVersion("v2"),
	google.WithAPIKey("api-key"),
)
defer client.Close()

result, err := client.Translate(ctx, "Hello", "en", "zh-CN")
```

不设置 API Key 时，官方客户端使用 Application Default Credentials。`WithApiKey` 为兼容上游保留，新增代码使用 `WithAPIKey`。

## V3

```go
client := google.NewTranslator(
	google.WithVersion("v3"),
	google.WithProjectID("project-id"),
	google.WithLocation("global"),
)
defer client.Close()

result, err := client.Translate(ctx, "Hello", "en", "zh-CN")
```

V3 必须通过 `WithProjectID` 配置项目，或用 `WithParent` 直接传入 `projects/{project}/locations/{location}`。非 `global` 区域如需专用端点，可通过 `WithClientOptions(option.WithEndpoint(...))` 设置。

官方接口认证、区域和语言代码以 [Google Cloud Translation 文档](https://cloud.google.com/translate/docs)为准。
