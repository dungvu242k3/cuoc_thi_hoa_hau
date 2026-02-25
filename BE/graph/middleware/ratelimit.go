/*
Package middleware bảo vệ ứng dụng.
Tác dụng của ratelimit.go:
- Chống bị DDoS bằng cách đếm số request trên mỗi IP lưu vào Redis.
- Vượt mức (vd: 100 req/min) -> Chặn cổ trả về lỗi 429 Too Many Requests.
*/
package middleware

import (
	"context"
	"cuoc_thi_hoa_hau/internal/types"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func RateLimitMiddleware(cache types.CacheService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = strings.Split(forwarded, ",")[0]
			}

			key := fmt.Sprintf("rate_limit:global:%s", ip)

			limit := 100
			window := time.Minute

			count, err := cache.Incr(context.Background(), key)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				cache.Set(context.Background(), key, 1, window)
			}

			if count > int64(limit) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
