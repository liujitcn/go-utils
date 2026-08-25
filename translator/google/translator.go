// Package google 提供 Google V1、V2 和 V3 翻译适配器。
package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/translate"
	translateV3 "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/googleapis/gax-go/v2"
	httputil "github.com/liujitcn/go-utils/http"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

const (
	defaultV1Endpoint = "https://translate.googleapis.com/translate_a/single"
	maxResponseSize   = 1 << 20
)

type v2Client interface {
	// Translate 调用 Google V2 文本翻译接口。
	Translate(
		ctx context.Context,
		inputs []string,
		target language.Tag,
		options *translate.Options,
	) ([]translate.Translation, error)
	// Close 关闭 Google V2 客户端。
	Close() error
}

type v3Client interface {
	// TranslateText 调用 Google V3 文本翻译接口。
	TranslateText(
		ctx context.Context,
		request *translatepb.TranslateTextRequest,
		options ...gax.CallOption,
	) (*translatepb.TranslateTextResponse, error)
	// Close 关闭 Google V3 客户端。
	Close() error
}

// Translator 是 Google 翻译适配器。
type Translator struct {
	options  options
	clientMu sync.Mutex
	clientV2 v2Client
	clientV3 v3Client
}

// NewTranslator 创建 Google 翻译适配器，默认使用非官方 V1 接口。
func NewTranslator(opts ...Option) *Translator {
	translatorOptions := options{
		version:    "v1",
		location:   "global",
		httpClient: &http.Client{Timeout: 30 * time.Second},
		v1Endpoint: defaultV1Endpoint,
	}
	for _, opt := range opts {
		opt(&translatorOptions)
	}
	return &Translator{options: translatorOptions}
}

// Translate 按配置的 Google 版本翻译单段文本。
func (translator *Translator) Translate(
	ctx context.Context,
	source, sourceLang, targetLang string,
) (string, error) {
	switch translator.options.version {
	case "v1":
		return translator.TranslateV1(ctx, source, sourceLang, targetLang)
	case "v2":
		text, _, err := translator.TranslateV2(ctx, source, sourceLang, targetLang)
		return text, err
	case "v3":
		return translator.TranslateV3(ctx, source, sourceLang, targetLang)
	default:
		return "", fmt.Errorf("Google translate: unsupported version %q", translator.options.version)
	}
}

// TranslateV1 使用无需凭证的非官方 Google 接口翻译文本。
func (translator *Translator) TranslateV1(
	ctx context.Context,
	source, sourceLang, targetLang string,
) (string, error) {
	var err error
	err = ctx.Err()
	if err != nil {
		return "", err
	}
	if source == "" {
		return "", fmt.Errorf("Google V1 translate: source cannot be empty")
	}

	client := httputil.NewClient(httputil.WithHTTPClient(translator.options.httpClient))
	var response *httputil.Response
	response, err = client.Do(
		http.MethodGet,
		translator.options.v1Endpoint,
		httputil.WithContext(ctx),
		httputil.WithQuery("client", "gtx"),
		httputil.WithQuery("sl", sourceLang),
		httputil.WithQuery("tl", targetLang),
		httputil.WithQuery("dt", "t"),
		httputil.WithQuery("q", source),
	)
	if err != nil {
		return "", fmt.Errorf("Google V1 translate: send request: %w", err)
	}
	if len(response.Body) > maxResponseSize {
		return "", fmt.Errorf("Google V1 translate: response exceeds %d bytes", maxResponseSize)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Google V1 translate: HTTP status %d", response.StatusCode)
	}

	var translated string
	translated, err = decodeV1Response(response.Body)
	if err != nil {
		return "", fmt.Errorf("Google V1 translate: %w", err)
	}
	return translated, nil
}

// TranslateV2 使用 Google Cloud Translation Basic V2 翻译文本。
func (translator *Translator) TranslateV2(
	ctx context.Context,
	source, sourceLang, targetLang string,
) (string, language.Tag, error) {
	var err error
	err = ctx.Err()
	if err != nil {
		return "", language.Und, err
	}
	if source == "" {
		return "", language.Und, fmt.Errorf("Google V2 translate: source cannot be empty")
	}

	var targetLanguage language.Tag
	targetLanguage, err = language.Parse(targetLang)
	if err != nil {
		return "", language.Und, fmt.Errorf("Google V2 translate: parse target language: %w", err)
	}

	translateOptions := &translate.Options{Format: translate.Text}
	if sourceLang != "" && sourceLang != "auto" {
		translateOptions.Source, err = language.Parse(sourceLang)
		if err != nil {
			return "", language.Und, fmt.Errorf("Google V2 translate: parse source language: %w", err)
		}
	}

	var client v2Client
	client, err = translator.ensureV2Client(ctx)
	if err != nil {
		return "", language.Und, err
	}

	var translations []translate.Translation
	translations, err = client.Translate(ctx, []string{source}, targetLanguage, translateOptions)
	if err != nil {
		return "", language.Und, fmt.Errorf("Google V2 translate: %w", err)
	}
	if len(translations) != 1 {
		return "", language.Und, fmt.Errorf("Google V2 translate: expected one translation, got %d", len(translations))
	}

	detectedLanguage := translations[0].Source
	if detectedLanguage == language.Und && translateOptions.Source != language.Und {
		detectedLanguage = translateOptions.Source
	}
	return translations[0].Text, detectedLanguage, nil
}

// TranslateV3 使用 Google Cloud Translation Advanced V3 翻译文本。
func (translator *Translator) TranslateV3(
	ctx context.Context,
	source, sourceLang, targetLang string,
) (string, error) {
	var err error
	err = ctx.Err()
	if err != nil {
		return "", err
	}
	if source == "" {
		return "", fmt.Errorf("Google V3 translate: source cannot be empty")
	}

	parent := translator.options.parent
	if parent == "" {
		if translator.options.projectID == "" {
			return "", fmt.Errorf("Google V3 translate: project ID or parent is required")
		}
		parent = fmt.Sprintf(
			"projects/%s/locations/%s",
			translator.options.projectID,
			translator.options.location,
		)
	}

	request := &translatepb.TranslateTextRequest{
		Parent:             parent,
		TargetLanguageCode: targetLang,
		MimeType:           "text/plain",
		Contents:           []string{source},
	}
	if sourceLang != "" && sourceLang != "auto" {
		request.SourceLanguageCode = sourceLang
	}

	var client v3Client
	client, err = translator.ensureV3Client(ctx)
	if err != nil {
		return "", err
	}

	var response *translatepb.TranslateTextResponse
	response, err = client.TranslateText(ctx, request)
	if err != nil {
		return "", fmt.Errorf("Google V3 translate: %w", err)
	}
	if response == nil || len(response.GetTranslations()) != 1 {
		return "", fmt.Errorf("Google V3 translate: expected one translation")
	}
	return response.GetTranslations()[0].GetTranslatedText(), nil
}

// Close 关闭已创建的 Google 官方客户端。
func (translator *Translator) Close() error {
	translator.clientMu.Lock()
	defer translator.clientMu.Unlock()

	var closeErrors []error
	var closeErr error
	if translator.clientV2 != nil {
		closeErr = translator.clientV2.Close()
		if closeErr != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close Google V2 client: %w", closeErr))
		}
		translator.clientV2 = nil
	}
	if translator.clientV3 != nil {
		closeErr = translator.clientV3.Close()
		if closeErr != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close Google V3 client: %w", closeErr))
		}
		translator.clientV3 = nil
	}
	return errors.Join(closeErrors...)
}

// ensureV2Client 并发安全地创建或返回 Google V2 客户端。
func (translator *Translator) ensureV2Client(ctx context.Context) (v2Client, error) {
	translator.clientMu.Lock()
	defer translator.clientMu.Unlock()

	if translator.clientV2 != nil {
		return translator.clientV2, nil
	}

	var client *translate.Client
	var err error
	client, err = translate.NewClient(ctx, translator.googleClientOptions()...)
	if err != nil {
		return nil, fmt.Errorf("create Google V2 client: %w", err)
	}
	translator.clientV2 = client
	return translator.clientV2, nil
}

// ensureV3Client 并发安全地创建或返回 Google V3 客户端。
func (translator *Translator) ensureV3Client(ctx context.Context) (v3Client, error) {
	translator.clientMu.Lock()
	defer translator.clientMu.Unlock()

	if translator.clientV3 != nil {
		return translator.clientV3, nil
	}

	var client *translateV3.TranslationClient
	var err error
	client, err = translateV3.NewTranslationClient(ctx, translator.googleClientOptions()...)
	if err != nil {
		return nil, fmt.Errorf("create Google V3 client: %w", err)
	}
	translator.clientV3 = client
	return translator.clientV3, nil
}

// googleClientOptions 返回独立的 Google 客户端配置副本。
func (translator *Translator) googleClientOptions() []option.ClientOption {
	clientOptions := append([]option.ClientOption(nil), translator.options.clientOptions...)
	if translator.options.apiKey != "" {
		clientOptions = append(clientOptions, option.WithAPIKey(translator.options.apiKey))
	}
	return clientOptions
}

// decodeV1Response 解码 Google 非官方 V1 的异构数组响应。
func decodeV1Response(responseBody []byte) (string, error) {
	var payload []json.RawMessage
	var err error
	err = json.Unmarshal(responseBody, &payload)
	if err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if len(payload) == 0 || string(payload[0]) == "null" {
		return "", fmt.Errorf("empty response")
	}

	var segments [][]json.RawMessage
	err = json.Unmarshal(payload[0], &segments)
	if err != nil {
		return "", fmt.Errorf("decode segments: %w", err)
	}

	translatedSegments := make([]string, 0, len(segments))
	for _, segment := range segments {
		if len(segment) == 0 || string(segment[0]) == "null" {
			continue
		}

		var translatedSegment string
		err = json.Unmarshal(segment[0], &translatedSegment)
		if err != nil {
			return "", fmt.Errorf("decode translated segment: %w", err)
		}
		translatedSegments = append(translatedSegments, translatedSegment)
	}
	if len(translatedSegments) == 0 {
		return "", fmt.Errorf("empty translation")
	}
	return strings.Join(translatedSegments, ""), nil
}
