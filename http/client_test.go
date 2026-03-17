package http

import (
	"encoding/json"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClientGet 验证 GET 请求会合并基础地址、查询参数与请求头。
func TestClientGet(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		require.Equal(t, stdhttp.MethodGet, r.Method)
		require.Equal(t, "book", r.URL.Query().Get("keyword"))
		require.Equal(t, "10", r.URL.Query().Get("size"))
		require.Equal(t, "client", r.Header.Get("X-Default"))
		require.Equal(t, "request", r.Header.Get("X-Request"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithDefaultHeader("X-Default", "client"),
	)

	type responseBody struct {
		Message string `json:"message"`
	}

	var result responseBody
	err := client.Get(
		"/search",
		&result,
		WithQuery("keyword", "book"),
		WithQueries(map[string]string{"size": "10"}),
		WithHeader("X-Request", "request"),
	)
	require.NoError(t, err)
	require.Equal(t, "pong", result.Message)
}

// TestClientPostJSON 验证 POST 请求会发送 JSON 请求体。
func TestClientPostJSON(t *testing.T) {
	type requestBody struct {
		Name string `json:"name"`
	}

	type responseBody struct {
		OK bool `json:"ok"`
	}

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		require.Equal(t, stdhttp.MethodPost, r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		data, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload requestBody
		require.NoError(t, json.Unmarshal(data, &payload))
		require.Equal(t, "alice", payload.Name)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseBody{OK: true})
	}))
	defer server.Close()

	client := NewClient()

	var payload responseBody
	err := client.Post(server.URL, &payload, WithJSONBody(requestBody{Name: "alice"}))
	require.NoError(t, err)
	require.True(t, payload.OK)
}

// TestClientPostForm 验证 POST 请求会发送表单请求体。
func TestClientPostForm(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		require.NoError(t, r.ParseForm())
		require.Equal(t, "alice", r.Form.Get("name"))
		require.Equal(t, "admin", r.Form.Get("role"))
		w.WriteHeader(stdhttp.StatusCreated)
	}))
	defer server.Close()

	client := NewClient()

	var result struct{}
	err := client.Post(server.URL, &result, WithFormBody(url.Values{
		"name": []string{"alice"},
		"role": []string{"admin"},
	}))
	require.NoError(t, err)
}

// TestSingletonInit 验证默认客户端只在首次初始化时应用配置。
func TestSingletonInit(t *testing.T) {
	resetDefaultClientForTest()

	first := Init(WithTimeout(5), WithDefaultHeader("X-App", "first"))
	second := Init(WithTimeout(50), WithDefaultHeader("X-App", "second"))

	require.Same(t, first, second)
	require.Equal(t, "first", first.defaultHeaders.Get("X-App"))
	require.Equal(t, 5, int(first.httpClient.Timeout))
}

// TestSingletonStaticMethods 验证包级静态方法会复用默认客户端。
func TestSingletonStaticMethods(t *testing.T) {
	resetDefaultClientForTest()

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		require.Equal(t, "singleton", r.Header.Get("X-Default"))
		require.Equal(t, "1", r.URL.Query().Get("page"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "ok"})
	}))
	defer server.Close()

	Init(
		WithBaseURL(server.URL),
		WithDefaultHeader("X-Default", "singleton"),
	)

	var result map[string]string
	err := Get("/", &result, WithQuery("page", "1"))
	require.NoError(t, err)
	require.Equal(t, "ok", result["message"])
}

// TestSetDefaultClient 验证可以显式替换默认客户端实例。
func TestSetDefaultClient(t *testing.T) {
	resetDefaultClientForTest()

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		require.Equal(t, "replaced", r.Header.Get("X-Default"))
		w.WriteHeader(stdhttp.StatusNoContent)
	}))
	defer server.Close()

	SetDefaultClient(NewClient(
		WithBaseURL(server.URL),
		WithDefaultHeader("X-Default", "replaced"),
	))

	var result struct{}
	err := Head("/", &result)
	require.NoError(t, err)
}

// resetDefaultClientForTest 重置默认客户端，避免测试之间互相影响。
func resetDefaultClientForTest() {
	defaultClientMu.Lock()
	defer defaultClientMu.Unlock()
	defaultClient = nil
	defaultClientOnce = sync.Once{}
}
