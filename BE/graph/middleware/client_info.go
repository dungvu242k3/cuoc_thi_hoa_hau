/*
Package middleware hỗ trợ thu thập kiểm toán hệ thống.
Tác dụng của client_info.go:
- Trích xuất IP thật đằng sau Proxy/Load Balancer (qua X-Forwarded-For) và UserAgent.
- Phục vụ cho tính năng Audit (Ví dụ: ghi lại IP của Giám khảo khi chấm điểm).
*/
package middleware

import (
	"context"
	"net/http"
	"strings"
)

type clientInfoKey struct{}

var ClientInfoCtxKey = &clientInfoKey{}

type ClientInfo struct {
	IP        string
	UserAgent string
}

func ClientInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = strings.Split(forwarded, ",")[0]
		}
		if strings.Contains(ip, ":") && !strings.Contains(ip, "]") {
			parts := strings.Split(ip, ":")
			if len(parts) == 2 {
				ip = parts[0]
			}
		}

		ua := r.UserAgent()

		info := &ClientInfo{
			IP:        strings.TrimSpace(ip),
			UserAgent: ua,
		}

		ctx := context.WithValue(r.Context(), ClientInfoCtxKey, info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClientInfo(ctx context.Context) (string, string) {
	raw, _ := ctx.Value(ClientInfoCtxKey).(*ClientInfo)
	if raw == nil {
		return "", ""
	}
	return raw.IP, raw.UserAgent
}
