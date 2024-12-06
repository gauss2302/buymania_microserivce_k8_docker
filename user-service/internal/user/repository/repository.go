package repository

import "github.com/gauss2302/microtest/user-service/internal/entity"

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int) error
	ListUsers(limit, offset int) ([]*entity.User, error)
}
