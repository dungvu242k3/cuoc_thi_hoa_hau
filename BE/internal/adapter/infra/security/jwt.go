package security

import (
	"cuoc_thi_hoa_hau/internal/core/domain"
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

func (p *JWTProvider) Generate(user *domain.User, duration time.Duration) (*domain.AuthClaims, string, error) {
	claims := jwt.MapClaims{
		domain.ClaimKeyUserID: user.ID,
		domain.ClaimKeyRole:   user.RoleID,
		domain.ClaimKeyExp:    time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(p.secret))
	if err != nil {
		return nil, "", err
	}

	return &domain.AuthClaims{
		UserID: user.ID,
		Role:   user.RoleID,
	}, tokenStr, nil
}

func (p *JWTProvider) Validate(tokenString string) (*domain.AuthClaims, error) {
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
		userID, _ := claims[domain.ClaimKeyUserID].(string)
		role, _ := claims[domain.ClaimKeyRole].(string)

		return &domain.AuthClaims{
			UserID: userID,
			Role:   role,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
