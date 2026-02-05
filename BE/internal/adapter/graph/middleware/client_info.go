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
		// IP
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = strings.Split(forwarded, ",")[0]
		}
		// Clean port if present (IPv4:port or [IPv6]:port)
		// But RemoteAddr is usually "ip:port", XFF is usually just field.
		// If RemoteAddr, we should strip port.
		// Simple approach:
		if strings.Contains(ip, ":") && !strings.Contains(ip, "]") {
			// likely ipv4:port
			parts := strings.Split(ip, ":")
			if len(parts) == 2 {
				ip = parts[0]
			}
		}
		// If IPv6 [::1]:port ... complicated. For now trust simple extraction or XFF.

		// UA
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
