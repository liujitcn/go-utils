package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClientLimitsResponseBody 验证普通请求在超过配置上限时返回稳定错误。
func TestClientLimitsResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		var err error
		_, err = w.Write([]byte("12345"))
		if err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client := NewClient(WithMaxResponseBodyBytes(4), WithRequestLog(false))
	_, err := client.Do(http.MethodGet, server.URL)
	if !errors.Is(err, ErrResponseBodyTooLarge) {
		t.Fatalf("响应体超限错误 = %v", err)
	}
}

// TestClientKeepsUnlimitedResponseDefault 验证未配置上限时保持兼容并读取完整响应。
func TestClientKeepsUnlimitedResponseDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		var err error
		_, err = w.Write([]byte("12345"))
		if err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client := NewClient(WithRequestLog(false))
	response, err := client.Do(http.MethodGet, server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(response.Body) != "12345" {
		t.Fatalf("响应体 = %q", response.Body)
	}
}
