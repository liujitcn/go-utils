// Package volc 提供火山引擎机器翻译适配器。
package volc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/volcengine/volc-sdk-golang/base"
)

const (
	translateAction  = "TranslateText"
	translateHost    = "translate.volcengineapi.com"
	translateRegion  = "cn-north-1"
	translateService = "translate"
	translateVersion = "2020-06-01"
)

type translateClient interface {
	// CtxJson 调用需要签名的火山引擎 JSON 接口。
	CtxJson(
		ctx context.Context,
		api string,
		query url.Values,
		body string,
	) ([]byte, int, error)
}

// Translator 是火山引擎机器翻译适配器。
type Translator struct {
	accessKey  string
	region     string
	httpClient *http.Client
	client     translateClient
}

// NewTranslator 创建火山引擎机器翻译适配器。
func NewTranslator(accessKey, secretKey string, opts ...Option) (*Translator, error) {
	secretKey = strings.TrimRight(secretKey, "\r\n")
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("create Volc translator: access key and secret key are required")
	}

	translator := &Translator{
		accessKey: accessKey,
		region:    translateRegion,
	}
	for _, opt := range opts {
		opt(translator)
	}

	serviceInfo := &base.ServiceInfo{
		Timeout: 30 * time.Second,
		Scheme:  "https",
		Host:    translateHost,
		Header:  http.Header{"Accept": []string{"application/json"}},
		Credentials: base.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			Region:          translator.region,
			Service:         translateService,
		},
	}
	apiInfo := map[string]*base.ApiInfo{
		translateAction: {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{translateAction},
				"Version": []string{translateVersion},
			},
		},
	}
	client := base.NewClient(serviceInfo, apiInfo)
	client.SetCredential(serviceInfo.Credentials)
	if translator.httpClient != nil {
		client.Client = translator.httpClient
	}
	translator.client = client
	return translator, nil
}

// Translate 翻译单段文本。
func (translator *Translator) Translate(
	ctx context.Context,
	source, sourceLang, targetLang string,
) (string, error) {
	if source == "" {
		return "", fmt.Errorf("Volc translate: source cannot be empty")
	}

	var translations []string
	var err error
	translations, err = translator.TranslateBatch(ctx, []string{source}, sourceLang, targetLang)
	if err != nil {
		return "", err
	}
	return translations[0], nil
}

// TranslateBatch 通过一个火山引擎请求批量翻译文本。
func (translator *Translator) TranslateBatch(
	ctx context.Context,
	sources []string,
	sourceLang, targetLang string,
) ([]string, error) {
	var err error
	err = ctx.Err()
	if err != nil {
		return nil, err
	}
	if len(sources) == 0 {
		return []string{}, nil
	}
	for index, source := range sources {
		if source == "" {
			return nil, fmt.Errorf("Volc translate: source at index %d cannot be empty", index)
		}
	}

	request := translateTextRequest{
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
		TextList:       sources,
	}
	if sourceLang == "auto" {
		request.SourceLanguage = ""
	}

	var response translateTextResponse
	response, err = translator.callAPI(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(response.TranslationList) != len(sources) {
		return nil, fmt.Errorf(
			"Volc translate: expected %d translations, got %d",
			len(sources),
			len(response.TranslationList),
		)
	}

	translations := make([]string, len(response.TranslationList))
	for index, item := range response.TranslationList {
		translations[index] = item.Text
	}
	return translations, nil
}

// String 返回不泄漏完整 AccessKey 的翻译器摘要。
func (translator *Translator) String() string {
	maskedAccessKey := strings.Repeat("*", len(translator.accessKey))
	if len(translator.accessKey) > 10 {
		maskedAccessKey = translator.accessKey[:6] + "****" + translator.accessKey[len(translator.accessKey)-4:]
	}
	return fmt.Sprintf("Translator{region: %s, accessKey: %s}", translator.region, maskedAccessKey)
}

// callAPI 调用火山引擎翻译接口并校验标准响应。
func (translator *Translator) callAPI(
	ctx context.Context,
	request translateTextRequest,
) (translateTextResponse, error) {
	var requestBody []byte
	var err error
	requestBody, err = json.Marshal(request)
	if err != nil {
		return translateTextResponse{}, fmt.Errorf("Volc translate: encode request: %w", err)
	}

	var responseBody []byte
	var statusCode int
	responseBody, statusCode, err = translator.client.CtxJson(
		ctx,
		translateAction,
		nil,
		string(requestBody),
	)
	if err != nil {
		return translateTextResponse{}, fmt.Errorf("Volc translate: send request: %w", err)
	}
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return translateTextResponse{}, fmt.Errorf("Volc translate: HTTP status %d", statusCode)
	}

	var response translateTextResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return translateTextResponse{}, fmt.Errorf("Volc translate: decode response: %w", err)
	}

	metadata := response.ResponseMetadata
	if response.ResponseMetaData.Error != nil {
		metadata = response.ResponseMetaData
	}
	if metadata.Error != nil {
		return translateTextResponse{}, fmt.Errorf(
			"Volc translate: code=%s, message=%s",
			metadata.Error.Code,
			metadata.Error.Message,
		)
	}
	return response, nil
}
