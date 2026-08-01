# translator

`translator` 提供统一的文本翻译接口，具体厂商以独立 Go 模块提供，调用方只需引入实际使用的 Provider。

## 模块

| Provider | 模块 | 实现方式 | 批量翻译 |
|---|---|---|---|
| Alibaba | `github.com/liujitcn/go-utils/translator/alibaba` | 阿里云官方 SDK | 否 |
| Baidu | `github.com/liujitcn/go-utils/translator/baidu` | 百度官方 HTTP 签名协议 | 否 |
| Google V1 | `github.com/liujitcn/go-utils/translator/google` | 非官方 `client=gtx` 接口 | 否 |
| Google V2/V3 | `github.com/liujitcn/go-utils/translator/google` | Google Cloud 官方 SDK | 否 |
| Volc | `github.com/liujitcn/go-utils/translator/volc` | 火山引擎官方 Go SDK 签名 | 是 |

Google V1 无需凭证，但不是受支持的公开服务，没有 SLA，可能被限流或随时变更。生产环境优先选择官方接口。

## 公共接口

```go
type Translator interface {
	Translate(
		ctx context.Context,
		source string,
		sourceLang string,
		targetLang string,
	) (string, error)
}

type BatchTranslator interface {
	TranslateBatch(
		ctx context.Context,
		sources []string,
		sourceLang string,
		targetLang string,
	) ([]string, error)
}
```

不同厂商的语言代码并不完全一致，应以对应厂商文档为准。调用方应设置超时，并在业务层处理限流、重试、缓存和占位符保护；本模块不会自动修改待翻译正文。

## 使用

以阿里云为例：

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/liujitcn/go-utils/translator/alibaba"
)

func main() {
	client, err := alibaba.NewTranslator("access-key-id", "access-key-secret")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.Translate(ctx, "Hello", "en", "zh")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
```

各 Provider 的配置和注意事项见：

- [Alibaba](alibaba/README.md)
- [Baidu](baidu/README.md)
- [Google](google/README.md)
- [Volc](volc/README.md)

## 上游与修复

代码迁移自 `tx7do/go-utils` 的 `67ec3305fb0494389482534a6692b71d4e388786` 版本，并完成以下调整：

- 公共接口增加 `context.Context`，Volc 批量接口同步支持取消。
- 修复 Alibaba 空响应解引用导致的 panic，并把截止时间映射为 SDK 超时。
- 修复 Baidu 可预测盐值、多段结果只返回第一段、请求无法取消的问题。
- 修复 Google V1 强制类型断言 panic、V2 忽略源语言和空结果越界、V3 缺少 `Parent`、客户端并发初始化竞态，并提供 `Close`。
- Volc 改用官方 Go SDK 签名，移除凭证和正文调试输出，修复响应层级、翻译字段、显式凭证优先级及短 AccessKey 脱敏越界。
