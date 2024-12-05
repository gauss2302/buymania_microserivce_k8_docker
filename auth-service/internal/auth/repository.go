package auth

import (
	"github.com/gauss2302/microtest/auth-service/internal/entity"
	"time"
)

type Repository interface {
	StoreToken(userID int, token string, expiration time.Duration) error
	//GetToken(token string) (*entity)
	GetToken(token string) (*entity.TokenDetails, error)
	DeleteToken(token string) error
}
