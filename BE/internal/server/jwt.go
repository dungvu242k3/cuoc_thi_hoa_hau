/*
Package server chứa các công cụ hỗ trợ cho việc vận hành HTTP và danh tính.
Tác dụng của jwt.go:
- Ẩn giấu thư viện "github.com/golang-jwt/jwt/v5" đi.
- Cung cấp các hàm tạo Token (khi Login) và kiểm tra Token (để Middleware chặn mồm).
*/
package server

import (
	"cuoc_thi_hoa_hau/internal/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret string
}

func NewJWTProvider(secret string) *JWTProvider {
	return &JWTProvider{
		secret: secret,
	}
}

// Generate tạo ra một JWT String chứa (UserID, Role) và thời hạn (Exp).
// Chuỗi này được ký bằng thuật toán HS256 với khóa mật (secret) cấu hình trong YAML.
func (p *JWTProvider) Generate(user *model.User, duration time.Duration) (*model.AuthClaims, string, error) {
	claims := jwt.MapClaims{
		model.ClaimKeyUserID: user.ID,
		model.ClaimKeyRole:   user.RoleID,
		model.ClaimKeyExp:    time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(p.secret))
	if err != nil {
		return nil, "", err
	}

	return &model.AuthClaims{
		UserID: user.ID,
		Role:   user.RoleID,
	}, tokenStr, nil
}

func (p *JWTProvider) Validate(tokenString string) (*model.AuthClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, _ := claims[model.ClaimKeyUserID].(string)
		role, _ := claims[model.ClaimKeyRole].(string)

		return &model.AuthClaims{
			UserID: userID,
			Role:   role,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
