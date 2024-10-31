package redis_cache

import (
	"auth/internal/auth/domain/models"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type UserCache struct {
	Redis
}

func NewUserCache(redis *Redis) *UserCache {
	return &UserCache{
		Redis: *redis,
	}
}

func (r *UserCache) Get(ctx context.Context, username string) (models.User, error) {
	key := createUserKey(username)

	bytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	if err := json.Unmarshal(bytes, &user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserCache) Set(ctx context.Context, username string, user models.User, ttl time.Duration) error {
	key := createUserKey(username)

	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, ttl).Err()
}

func (r *UserCache) Delete(ctx context.Context, username string) error {
	key := createUserKey(username)

	return r.client.Del(ctx, key).Err()
}

func createUserKey(username string) string {
	return fmt.Sprintf("user:%s", username)
}
