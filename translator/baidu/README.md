# Baidu Translator

百度通用文本翻译 Provider，按百度开放平台协议生成签名并发送 HTTP 请求。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/translator/baidu@latest
~~~

## 使用

~~~go
client := baidu.NewTranslator(
	"app-id",
	"secret-key",
	baidu.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
result, err := client.Translate(ctx, "Hello", "en", "zh")
~~~

默认 HTTP 客户端超时为 30 秒；WithHTTPClient 可设置代理、连接池或自定义超时，传入的客户端由调用方管理。响应包含多段结果时，Provider 会按换行拼接所有译文。

语言代码、频率限制和错误码以[百度翻译开放平台文档](https://fanyi-api.baidu.com/doc/21)为准。
