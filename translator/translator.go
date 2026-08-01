// Package translator 定义文本翻译模块的公共接口。
package translator

import "context"

// Translator 定义单文本翻译接口。
type Translator interface {
	// Translate 翻译单段文本。
	Translate(ctx context.Context, source, sourceLang, targetLang string) (string, error)
}

// BatchTranslator 定义支持单次批量请求的翻译接口。
type BatchTranslator interface {
	// TranslateBatch 通过单次请求批量翻译文本。
	TranslateBatch(ctx context.Context, sources []string, sourceLang, targetLang string) ([]string, error)
}
