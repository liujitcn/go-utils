# Google Translator

Google Provider 同时支持 V1、V2 和 V3 三种模式。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/translator/google@latest
~~~

## V1

~~~go
client := google.NewTranslator(google.WithVersion("v1"))
result, err := client.Translate(ctx, "Hello", "en", "zh-CN")
~~~

V1 使用无需凭证的非官方接口，没有 SLA，适合开发、演示或可以接受失败的低风险任务。生产环境应评估服务商的官方接口和合规要求。

## V2

~~~go
client := google.NewTranslator(
	google.WithVersion("v2"),
	google.WithAPIKey("api-key"),
)
defer client.Close()

result, detected, err := client.TranslateV2(ctx, "Hello", "en", "zh-CN")
~~~

未设置 API Key 时，官方客户端使用 Application Default Credentials。TranslateV2 除译文外还返回检测到的源语言。

## V3

~~~go
client := google.NewTranslator(
	google.WithVersion("v3"),
	google.WithProjectID("project-id"),
	google.WithLocation("global"),
)
defer client.Close()

result, err := client.TranslateV3(ctx, "Hello", "en", "zh-CN")
~~~

V3 必须配置 WithProjectID，或通过 WithParent 直接传入 projects/{project}/locations/{location}。非 global 区域的专用端点可通过 WithClientOptions(option.WithEndpoint(...)) 设置。

WithHTTPClient 只影响 V1；WithClientOptions 用于 V2/V3；WithApiKey 是兼容旧调用的别名，新代码使用 WithAPIKey。

认证、区域和语言代码以 [Google Cloud Translation 文档](https://cloud.google.com/translate/docs)为准。
