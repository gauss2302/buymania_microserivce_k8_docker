package auth

import "github.com/gauss2302/microtest/auth-service/internal/entity"

type Usecase interface {
	Login(request *entity.LoginRequest) (*entity.TokenResponse, error)
	Register(request *entity.RegisterRequest) (*entity.TokenResponse, error)
	ValidateToken(token string) (*entity.TokenDetails, error)
	Logout(token string) error
}
