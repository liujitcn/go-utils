// Package baidu 提供百度通用文本翻译适配器。
package baidu

import (
	"context"
	"crypto/md5" // #nosec G501 百度翻译签名协议固定使用 MD5。
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baiduEndpoint   = "https://fanyi-api.baidu.com/api/trans/vip/translate"
	maxResponseSize = 1 << 20
)

type baiduResponse struct {
	TransResult []struct {
		Source      string `json:"src"`
		Destination string `json:"dst"`
	} `json:"trans_result"`
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// Translator 是百度通用文本翻译适配器。
type Translator struct {
	appID     string
	secretKey string
	client    *http.Client
	endpoint  string
}

// NewTranslator 创建百度通用文本翻译适配器。
func NewTranslator(appID, secretKey string, opts ...Option) *Translator {
	translator := &Translator{
		appID:     appID,
		secretKey: secretKey,
		client:    &http.Client{Timeout: 30 * time.Second},
		endpoint:  baiduEndpoint,
	}
	for _, opt := range opts {
		opt(translator)
	}
	return translator
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
		return "", fmt.Errorf("Baidu translate: source cannot be empty")
	}
	if translator.appID == "" || translator.secretKey == "" {
		return "", fmt.Errorf("Baidu translate: app ID and secret key are required")
	}

	var salt string
	salt, err = newSalt()
	if err != nil {
		return "", fmt.Errorf("Baidu translate: generate salt: %w", err)
	}

	form := url.Values{}
	form.Set("q", source)
	form.Set("from", sourceLang)
	form.Set("to", targetLang)
	form.Set("appid", translator.appID)
	form.Set("salt", salt)
	form.Set("sign", translator.generateSign(source, salt))

	var request *http.Request
	request, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		translator.endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("Baidu translate: create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response *http.Response
	response, err = translator.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("Baidu translate: send request: %w", err)
	}

	var responseBody []byte
	responseBody, err = io.ReadAll(io.LimitReader(response.Body, maxResponseSize))
	closeErr := response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("Baidu translate: read response: %w", err)
	}
	if closeErr != nil {
		return "", fmt.Errorf("Baidu translate: close response: %w", closeErr)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Baidu translate: HTTP status %d", response.StatusCode)
	}

	var result baiduResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", fmt.Errorf("Baidu translate: decode response: %w", err)
	}
	if result.ErrorCode != "" {
		return "", fmt.Errorf("Baidu translate: code=%s, message=%s", result.ErrorCode, result.ErrorMsg)
	}
	if len(result.TransResult) == 0 {
		return "", fmt.Errorf("Baidu translate: empty translation")
	}

	translations := make([]string, 0, len(result.TransResult))
	for _, item := range result.TransResult {
		translations = append(translations, item.Destination)
	}
	return strings.Join(translations, "\n"), nil
}

// generateSign 生成百度翻译协议要求的请求签名。
func (translator *Translator) generateSign(source, salt string) string {
	content := translator.appID + source + salt + translator.secretKey
	digest := md5.Sum([]byte(content)) // #nosec G401 百度翻译签名协议固定使用 MD5。
	return hex.EncodeToString(digest[:])
}

// newSalt 生成不可预测的请求盐值。
func newSalt() (string, error) {
	buffer := make([]byte, 16)
	var err error
	_, err = cryptorand.Read(buffer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}
