# Volc Translator

火山引擎机器翻译 Provider，使用火山引擎官方 Go SDK 完成请求签名，支持单文本和批量翻译。

## 安装

```bash
go get github.com/liujitcn/go-utils/translator/volc@latest
```

## 使用

```go
client, err := volc.NewTranslator(
	"access-key",
	"secret-key",
	volc.WithRegion("cn-north-1"),
)
if err != nil {
	return err
}

result, err := client.Translate(ctx, "你好", "zh", "en")
results, err := client.TranslateBatch(ctx, []string{"你好", "世界"}, "zh", "en")
```

传入 `auto` 作为源语言时，请求会省略 `SourceLanguage`，由服务端自动识别。可通过 `volc.WithHTTPClient` 设置代理、连接池或自定义超时，传入的客户端由调用方管理。

语言代码、批量限制和服务配额以[火山引擎机器翻译文档](https://www.volcengine.com/docs/4640)为准。
