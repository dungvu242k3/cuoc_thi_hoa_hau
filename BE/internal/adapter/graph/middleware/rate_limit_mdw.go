package middleware

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/port"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func RateLimitMiddleware(cache port.CacheService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Identifier: IP or UserID if logged in
			// Simple approach: Use IP for now, or X-Forwarded-For
			ip := r.RemoteAddr
			// If behind proxy like Nginx/LoadBalancer:
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = strings.Split(forwarded, ",")[0]
			}

			// Key: rate_limit:global:<ip>
			key := fmt.Sprintf("rate_limit:global:%s", ip)

			// Allow 100 requests per minute per IP
			limit := 100
			window := time.Minute

			count, err := cache.Incr(context.Background(), key)
			if err != nil {
				// Fail open if cache is down? or log and proceed
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
