# Alibaba Translator

阿里云机器翻译 Provider，使用阿里云官方 `alimt` SDK。

## 安装

```bash
go get github.com/liujitcn/go-utils/translator/alibaba@latest
```

## 使用

```go
client, err := alibaba.NewTranslator(
	"access-key-id",
	"access-key-secret",
	alibaba.WithRegionID("cn-hangzhou"),
)
if err != nil {
	return err
}

result, err := client.Translate(ctx, "Hello", "en", "zh")
```

默认区域为 `cn-hangzhou`，对应端点为 `mt.cn-hangzhou.aliyuncs.com`。SDK 本身不原生接收 `context.Context`，Provider 会在调用前后检查上下文，并把 Context 截止时间映射到 SDK 的连接和读取超时。

语言代码、服务开通和配额以[阿里云机器翻译文档](https://help.aliyun.com/product/30396.html)为准。
