package ip

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetClientRealIP 验证真实 IP 提取优先级。
func TestGetClientRealIP(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
	assert.NoError(t, err)

	request.Header.Set(headerKeyXForwardedFor, "invalid, 10.0.0.1, 8.8.8.8")
	request.Header.Set(headerKeyXRealIP, "1.1.1.1")
	request.RemoteAddr = "127.0.0.1:8080"

	assert.Equal(t, "10.0.0.1", GetClientRealIP(request))
}

// TestGetClientRealIPFallback 验证降级读取逻辑。
func TestGetClientRealIPFallback(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
	assert.NoError(t, err)

	request.Header.Set(headerKeyXRealIP, "1.1.1.1")
	request.RemoteAddr = "127.0.0.1:8080"

	assert.Equal(t, "1.1.1.1", GetClientRealIP(request))

	request.Header.Del(headerKeyXRealIP)
	assert.Equal(t, "127.0.0.1", GetClientRealIP(request))
}

// TestGetIPFromRemoteAddr 验证 RemoteAddr 解析。
func TestGetIPFromRemoteAddr(t *testing.T) {
	assert.Equal(t, "127.0.0.1", GetIPFromRemoteAddr("127.0.0.1:8080"))
	assert.Equal(t, "127.0.0.1", GetIPFromRemoteAddr("127.0.0.1"))
	assert.Equal(t, "", GetIPFromRemoteAddr("invalid"))
}
