package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gauss2302/microtest/auth-service/internal/auth"
	"github.com/gauss2302/microtest/auth-service/internal/entity"
	"time"
)

type AuthUsecase struct {
	repo            auth.Repository
	tokenExpiration time.Duration
	userServiceURL  string
}

func NewAuthUsecase(repo auth.Repository, tokenExpiration time.Duration, userServiceURL string) *AuthUsecase {
	return &AuthUsecase{
		repo:            repo,
		tokenExpiration: tokenExpiration,
		userServiceURL:  userServiceURL,
	}
}
func (u *AuthUsecase) Login(request *entity.LoginRequest) (*entity.TokenResponse, error) {
	// Call user service to validate username and password

	// Generate token

	// Store token in memcached

	// Return token

	return nil, nil
}

func (u *AuthUsecase) Logout(token string) error {
	return u.repo.DeleteToken(token)
}

func (u *AuthUsecase) Register(request *entity.RegisterRequest) (*entity.TokenResponse, error) {
	// Hash password

	// Call user service to create user

	// Generate & store token
	return nil, nil
}

// Helpers

func (u *AuthUsecase) generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil

}

func (u *AuthUsecase) verifyCredentials(email, password string) error {
	return nil
}
