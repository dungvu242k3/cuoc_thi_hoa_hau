/*
Package middleware chứa các bộ lọc (Filters) chặn trước cửa các API.
Tác dụng của auth.go:
- Chụp lấy chuẩn Authorization Bearer Token trên Request Header.
- Giải mã và nhét thông tin (UserID, Role) vào r.Context() để các Resolvers bên trong có thể lấy ra xài.
*/
package middleware

import (
	"context"
	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"
	"fmt"
	"net/http"
	"strings"
)

type contextKey struct {
	name string
}

var UserCtxKey = &contextKey{"user"}

func Middleware(tokenProvider types.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

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

			ctx := context.WithValue(r.Context(), UserCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ForContext(ctx context.Context) string {
	raw, _ := ctx.Value(UserCtxKey).(*model.AuthClaims)
	if raw == nil {
		return ""
	}
	return raw.UserID
}

// RequireRole là một bảo vệ cấp 2 (Guard). Service/Resolver có thể gọi cái này để chặn đứng mấy thanh niên
// đang mượn token nhưng không phải quyền (VD: thí sinh múa rìu đòi vào API của Ban Giám Khảo).
func RequireRole(ctx context.Context, role model.Role) error {
	raw, _ := ctx.Value(UserCtxKey).(*model.AuthClaims)
	if raw == nil {
		return fmt.Errorf("unauthorized")
	}
	if raw.Role != string(role) && raw.Role != string(model.RoleAdmin) {
		return fmt.Errorf("forbidden: requires %s", role)
	}
	return nil
}
func RequirePermission(ctx context.Context, requiredPerm string) error {
	raw, _ := ctx.Value(UserCtxKey).(*model.AuthClaims)
	if raw == nil {
		return fmt.Errorf("unauthorized")
	}

	if model.HasPermission(raw.Role, requiredPerm) {
		return nil
	}

	return fmt.Errorf("forbidden: requires permission %s", requiredPerm)
}
