// internal/user/usecase/interface.go
package usecase

import "github.com/gauss2302/microtest/user-service/internal/entity"

type UserUsecase interface {
	CreateUser(req *entity.CreateUserRequest) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(id int, req *entity.UpdateUserRequest) (*entity.User, error)
	DeleteUser(id int) error
	ListUsers(limit, offset int) ([]*entity.User, error)
	VerifyCredentials(email, password string) (*entity.User, error)
}
