package service

import (
	"context"
	"errors"
	"time"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	repo            port.UserRepository
	passwordEncoder port.PasswordEncoder
	tokenProvider   port.TokenProvider
}

func NewAuthService(repo port.UserRepository, passwordEncoder port.PasswordEncoder, tokenProvider port.TokenProvider) port.AuthService {
	return &AuthService{
		repo:            repo,
		passwordEncoder: passwordEncoder,
		tokenProvider:   tokenProvider,
	}
}

func (s *AuthService) Register(ctx context.Context, username, password, role string) (*domain.AuthClaims, string, error) {
	// 1. Hash Password
	hashed, err := s.passwordEncoder.Hash(password)
	if err != nil {
		return nil, "", err
	}

	// 2. Create User
	user := &domain.User{
		Username: username,
		Password: hashed,
		RoleID:   role,
	}
	// 3. Create User (Atomic Insert)
	if err := s.repo.Create(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, "", errors.New("username already exists")
		}
		return nil, "", err
	}

	// 4. Generate Token
	return s.tokenProvider.Generate(user, 72*time.Hour)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*domain.AuthClaims, string, error) {
	// 1. Find User
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// 2. Check Password
	if err := s.passwordEncoder.Compare(user.Password, password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// 3. Generate Token
	return s.tokenProvider.Generate(user, 72*time.Hour)
}
