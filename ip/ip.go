package ip

import (
	"net"
	"net/http"
	"strings"
)

const (
	headerKeyXForwardedFor = "X-Forwarded-For"
	headerKeyXRealIP       = "X-Real-IP"
)

// GetClientRealIP 获取客户端真实 IP。
func GetClientRealIP(request *http.Request) string {
	if request == nil {
		return ""
	}

	// 优先读取代理链路中的真实客户端 IP。
	forwardedFor := request.Header.Get(headerKeyXForwardedFor)
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		for _, item := range ips {
			clientIP := strings.TrimSpace(item)
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
	}

	// 反向代理场景再尝试读取 X-Real-IP。
	realIP := strings.TrimSpace(request.Header.Get(headerKeyXRealIP))
	if net.ParseIP(realIP) != nil {
		return realIP
	}

	return GetIPFromRemoteAddr(request.RemoteAddr)
}

// GetIPFromRemoteAddr 从 RemoteAddr 中提取 IP。
func GetIPFromRemoteAddr(hostAddress string) string {
	hostAddress = strings.TrimSpace(hostAddress)
	if hostAddress == "" {
		return ""
	}

	host := hostAddress
	if strings.Contains(hostAddress, ":") {
		var err error
		host, _, err = net.SplitHostPort(hostAddress)
		if err != nil {
			return ""
		}
	}

	if net.ParseIP(host) == nil {
		return ""
	}

	return host
}
