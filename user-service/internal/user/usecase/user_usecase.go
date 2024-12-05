package usecase

import (
	"fmt"
	"github.com/gauss2302/microtest/user-service/internal/entity"
	"github.com/gauss2302/microtest/user-service/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) CreateUser(req *entity.CreateUserRequest) (*entity.User, error) {
	// Check if user with this email exists
	if _, err := u.repo.GetUserByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	return u.repo.CreateUser(user)
}

func (u *userUsecase) GetUserByID(id int) (*entity.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userUsecase) GetUserByEmail(email string) (*entity.User, error) {
	return u.repo.GetUserByEmail(email)
}

func (u *userUsecase) UpdateUser(id int, req *entity.UpdateUserRequest) (*entity.User, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Email != nil {
		// Check if new email is already taken
		if *req.Email != user.Email {
			if _, err := u.repo.GetUserByEmail(*req.Email); err == nil {
				return nil, fmt.Errorf("email %s is already taken", *req.Email)
			}
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("error hashing password: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	return u.repo.UpdateUser(user)
}

func (u *userUsecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}

func (u *userUsecase) ListUsers(limit, offset int) ([]*entity.User, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	return u.repo.ListUsers(limit, offset)
}

func (u *userUsecase) VerifyCredentials(email, password string) (*entity.User, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return user, nil
}
