// Package alibaba 提供阿里云机器翻译适配器。
package alibaba

import (
	"context"
	"fmt"
	"time"

	alimt "github.com/alibabacloud-go/alimt-20190107/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

const defaultRegionID = "cn-hangzhou"

type translateClient interface {
	// TranslateGeneralWithOptions 调用阿里云通用文本翻译接口。
	TranslateGeneralWithOptions(
		request *alimt.TranslateGeneralRequest,
		runtime *service.RuntimeOptions,
	) (*alimt.TranslateGeneralResponse, error)
}

// Translator 是阿里云机器翻译适配器。
type Translator struct {
	client   translateClient
	regionID string
}

// NewTranslator 创建阿里云机器翻译适配器。
func NewTranslator(accessKeyID, accessKeySecret string, opts ...Option) (*Translator, error) {
	if accessKeyID == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("create Alibaba translator: access key ID and secret are required")
	}

	translator := &Translator{regionID: defaultRegionID}
	for _, opt := range opts {
		opt(translator)
	}
	if translator.regionID == "" {
		return nil, fmt.Errorf("create Alibaba translator: region ID is required")
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fmt.Sprintf("mt.%s.aliyuncs.com", translator.regionID)),
	}

	var client *alimt.Client
	var err error
	client, err = alimt.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create Alibaba translator: %w", err)
	}

	translator.client = client
	return translator, nil
}

// Translate 翻译单段文本。
func (translator *Translator) Translate(
	ctx context.Context,
	source, sourceLang, targetLang string,
) (string, error) {
	var err error
	err = ctx.Err()
	if err != nil {
		return "", err
	}
	if source == "" {
		return "", fmt.Errorf("Alibaba translate: source cannot be empty")
	}

	request := &alimt.TranslateGeneralRequest{
		SourceLanguage: tea.String(sourceLang),
		TargetLanguage: tea.String(targetLang),
		SourceText:     tea.String(source),
		FormatType:     tea.String("text"),
		Scene:          tea.String("general"),
	}
	runtime := &service.RuntimeOptions{}
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		if timeout <= 0 {
			return "", context.DeadlineExceeded
		}
		timeoutMilliseconds := int(timeout.Milliseconds())
		if timeoutMilliseconds == 0 {
			timeoutMilliseconds = 1
		}
		runtime.ConnectTimeout = tea.Int(timeoutMilliseconds)
		runtime.ReadTimeout = tea.Int(timeoutMilliseconds)
	}

	var response *alimt.TranslateGeneralResponse
	response, err = translator.client.TranslateGeneralWithOptions(request, runtime)
	if err != nil {
		return "", fmt.Errorf("Alibaba translate: %w", err)
	}
	err = ctx.Err()
	if err != nil {
		return "", err
	}
	if response == nil || response.Body == nil {
		return "", fmt.Errorf("Alibaba translate: empty response")
	}
	if response.Body.Code != nil && tea.Int32Value(response.Body.Code) != 200 {
		return "", fmt.Errorf(
			"Alibaba translate: code=%d, message=%s",
			tea.Int32Value(response.Body.Code),
			tea.StringValue(response.Body.Message),
		)
	}
	if response.Body.Data == nil || response.Body.Data.Translated == nil {
		return "", fmt.Errorf("Alibaba translate: empty translation")
	}

	return tea.StringValue(response.Body.Data.Translated), nil
}
