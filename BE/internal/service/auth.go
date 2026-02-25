/*
Package service là Tầng Nghiệp vụ (Business Logic), trung tâm của ứng dụng.
Tác dụng:
  - Nơi chứa mọi quy tắc, tính toán cốt lõi. (Vd: Nạp tiền phải trừ phí, đăng ký phải băm mật khẩu).
  - Gọi DAO (Repository) để lưu hoặc lấy data, nhưng Service KHÔNG ĐƯỢC PHÉP biết Database đang dùng là
    Mongo, MySQL hay file text. Service chỉ giao tiếp với DAO thông qua Interface trong package `types`.
  - Được gọi từ Handler (REST) hoặc Resolver (GraphQL).
*/
package service

import (
	"context"
	"errors"
	"time"

	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	repo            types.UserRepository
	passwordEncoder types.PasswordEncoder
	tokenProvider   types.TokenProvider
}

func NewAuthService(repo types.UserRepository, passwordEncoder types.PasswordEncoder, tokenProvider types.TokenProvider) types.AuthService {
	return &AuthService{
		repo:            repo,
		passwordEncoder: passwordEncoder,
		tokenProvider:   tokenProvider,
	}
}

// Register là logic lõi để đăng ký tài khoản.
// Nó thực hiện băm mật khẩu, gọi DB qua repo.Create, và rốt cuộc là sinh JWT token.
func (s *AuthService) Register(ctx context.Context, username, password, role string) (*model.AuthClaims, string, error) {
	hashed, err := s.passwordEncoder.Hash(password)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		Username: username,
		Password: hashed,
		RoleID:   role,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, "", errors.New("username already exists")
		}
		return nil, "", err
	}

	return s.tokenProvider.Generate(user, 72*time.Hour)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*model.AuthClaims, string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := s.passwordEncoder.Compare(user.Password, password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	return s.tokenProvider.Generate(user, 72*time.Hour)
}
