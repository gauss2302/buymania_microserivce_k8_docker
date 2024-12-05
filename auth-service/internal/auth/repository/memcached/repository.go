package memcached

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gauss2302/microtest/auth-service/internal/entity"
	"time"
)

type AuthRepository struct {
	client *memcache.Client
}

func NewAuthRepository(client *memcache.Client) *AuthRepository {
	return &AuthRepository{
		client: client,
	}
}

func (r *AuthRepository) StoreToken(userID int, token string, expiration time.Duration) error {
	tokenDetails := &entity.TokenDetails{
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiration),
	}

	// Serialization of tokenDetails
	value, err := json.Marshal(tokenDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal token details: %w", err)
	}

	// Store token in memcached
	err = r.client.Set(&memcache.Item{
		Key:        token,
		Value:      value,
		Expiration: int32(expiration.Seconds()),
	})

	if err != nil {
		return fmt.Errorf("failed to store token in memcached: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetToken(token string) (*entity.TokenDetails, error) {
	// Get token from memcached
	item, err := r.client.Get(token)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil, fmt.Errorf("token not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get token from memcached: %w", err)
	}

	// Deserialization of tokenDetails
	var tokenDetails entity.TokenDetails
	if err := json.Unmarshal(item.Value, &tokenDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token details: %w", err)
	}

	// Check if token is expired
	if time.Now().After(tokenDetails.ExpiresAt) {
		return nil, fmt.Errorf("token is expired")
	}

	return &tokenDetails, nil
}

func (r *AuthRepository) DeleteToken(token string) error {
	err := r.client.Delete(token)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return fmt.Errorf("failed to delete token from memcached: %w", err)
	}
	return nil
}
