package middleware

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"
	"fmt"
	"net/http"
	"strings"
)

type contextKey struct {
	name string
}

var UserCtxKey = &contextKey{"user"}

func Middleware(tokenProvider port.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated access (resolvers will check if needed)
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate Bearer format
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := tokenProvider.Validate(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Extract user info
			userID, _ := claims["user_id"].(string)
			roleStr, _ := claims["role"].(string)
			role := domain.Role(roleStr)

			authClaims := &domain.AuthClaims{
				UserID: userID,
				Role:   role,
			}

			// Put into context
			ctx := context.WithValue(r.Context(), UserCtxKey, authClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ForContext(ctx context.Context) string {
	raw, _ := ctx.Value(UserCtxKey).(*domain.AuthClaims)
	if raw == nil {
		return ""
	}
	return raw.UserID
}

func RequireRole(ctx context.Context, role domain.Role) error {
	raw, _ := ctx.Value(UserCtxKey).(*domain.AuthClaims)
	if raw == nil {
		return fmt.Errorf("unauthorized")
	}
	if raw.Role != role && raw.Role != domain.RoleAdmin { // Admin bypass
		return fmt.Errorf("forbidden: requires %s", role)
	}
	return nil
}
