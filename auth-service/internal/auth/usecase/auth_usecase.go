package usecase

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gauss2302/microtest/auth-service/internal/auth"
	"github.com/gauss2302/microtest/auth-service/internal/entity"
)

type AuthUsecase struct {
	repo            auth.Repository
	tokenExpiration time.Duration
	userServiceURL  string
	httpClient      *http.Client
}

func NewAuthUsecase(repo auth.Repository, tokenExpiration time.Duration, userServiceURL string) *AuthUsecase {
	return &AuthUsecase{
		repo:            repo,
		tokenExpiration: tokenExpiration,
		userServiceURL:  userServiceURL,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (u *AuthUsecase) Login(request *entity.LoginRequest) (*entity.TokenResponse, error) {
	// Verify credentials with user service
	userID, err := u.verifyCredentials(request.Email, request.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	// Generate token
	token, err := u.generateToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Store token in memcached
	err = u.repo.StoreToken(userID, token, u.tokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return &entity.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(u.tokenExpiration.Seconds()),
		CreatedAt:   time.Now(),
	}, nil
}

func (u *AuthUsecase) Register(request *entity.RegisterRequest) (*entity.TokenResponse, error) {
	// Call user service to create user
	userID, err := u.createUser(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	token, err := u.generateToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Store token
	err = u.repo.StoreToken(userID, token, u.tokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return &entity.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(u.tokenExpiration.Seconds()),
		CreatedAt:   time.Now(),
	}, nil
}

func (u *AuthUsecase) Logout(token string) error {
	return u.repo.DeleteToken(token)
}

func (u *AuthUsecase) ValidateToken(token string) (*entity.TokenDetails, error) {
	return u.repo.GetToken(token)
}

// Helper methods

func (u *AuthUsecase) verifyCredentials(email, password string) (int, error) {
	reqBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return 0, err
	}

	resp, err := u.httpClient.Post(
		u.userServiceURL+"/users/verify",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid credentials")
	}

	var user struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u *AuthUsecase) createUser(req *entity.RegisterRequest) (int, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	resp, err := u.httpClient.Post(
		u.userServiceURL+"/users",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to create user")
	}

	var user struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u *AuthUsecase) generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
