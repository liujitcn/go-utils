# Volc Translator

火山引擎机器翻译 Provider，使用官方 Go SDK 完成请求签名，支持单文本和单请求批量翻译。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/translator/volc@latest
~~~

## 使用

~~~go
client, err := volc.NewTranslator(
	"access-key",
	"secret-key",
	volc.WithRegion("cn-north-1"),
	volc.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
if err != nil {
	return err
}

result, err := client.Translate(ctx, "你好", "zh", "en")
results, err := client.TranslateBatch(ctx, []string{"你好", "世界"}, "zh", "en")
~~~

默认区域为 cn-north-1。WithHTTPClient 可设置代理、连接池或自定义超时，传入的客户端由调用方管理。批量请求返回结果顺序与输入顺序一致；sourceLang 传 auto 时会省略源语言，由服务端自动识别。

每个输入文本都必须非空，批量限制和服务配额以[火山引擎机器翻译文档](https://www.volcengine.com/docs/4640)为准。
