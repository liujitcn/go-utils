# translator

translator 定义统一的文本翻译接口；每个服务商实现位于独立 Go 子模块，调用方只需引入实际使用的 Provider。

## 安装

公共接口：

~~~bash
go get github.com/liujitcn/go-utils/translator@latest
~~~

使用具体 Provider 时，再安装对应模块，例如：

~~~bash
go get github.com/liujitcn/go-utils/translator/volc@latest
~~~

## 公共接口

~~~go
type Translator interface {
	Translate(ctx context.Context, source, sourceLang, targetLang string) (string, error)
}

type BatchTranslator interface {
	TranslateBatch(ctx context.Context, sources []string, sourceLang, targetLang string) ([]string, error)
}
~~~

所有 Provider 都接收 context.Context。调用方应设置截止时间，并在业务层处理重试、限流、缓存、占位符保护和供应商错误。

## Provider

| Provider | 模块 | 凭证 | 批量接口 | 默认行为 |
| --- | --- | --- | --- | --- |
| Alibaba | translator/alibaba | AccessKey ID、AccessKey Secret | 否 | 区域 cn-hangzhou |
| Baidu | translator/baidu | App ID、Secret Key | 否 | HTTP 客户端超时 30 秒 |
| Google | translator/google | V1 无凭证；V2/V3 使用 API Key 或 ADC | 否 | 默认 V1 |
| Volc | translator/volc | Access Key、Secret Key | 是 | 区域 cn-north-1 |

语言代码和配额由服务商决定，传入值不会被本模块自动转换。空文本会被各 Provider 拒绝；Volc 的批量请求还会拒绝其中任意空元素。

## 使用

~~~go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/liujitcn/go-utils/translator/volc"
)

func main() {
	client, err := volc.NewTranslator("access-key", "secret-key")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.Translate(ctx, "你好", "zh", "en")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
~~~

各 Provider 的凭证、区域、版本和 HTTP 客户端配置见：

- [Alibaba](alibaba/README.md)
- [Baidu](baidu/README.md)
- [Google](google/README.md)
- [Volc](volc/README.md)

## 资源管理

Google V2/V3 客户端按首次请求延迟创建，使用结束后调用 Google Provider 的 Close 释放客户端。Baidu 和 Volc 接受调用方提供的 *http.Client，生命周期由调用方管理。
