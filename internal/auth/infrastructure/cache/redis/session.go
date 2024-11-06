package redis_cache

import (
	"auth/internal/auth/domain/models"
	redisdb "auth/internal/pkg/database/redis"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type SessionCache struct {
	redisdb.Redis
}

func NewSessionCache(redis *redisdb.Redis) *SessionCache {
	return &SessionCache{
		Redis: *redis,
	}
}

func (r *SessionCache) Get(ctx context.Context, id string) (models.Session, error) {
	key := createSessionKey(id)

	bytes, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return models.Session{}, err
	}

	var session models.Session
	if err := json.Unmarshal(bytes, &session); err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (r *SessionCache) Set(ctx context.Context, id string, session models.Session, ttl time.Duration) error {
	key := createSessionKey(id)

	bytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, bytes, ttl).Err()
}

func (r *SessionCache) Delete(ctx context.Context, id string) error {
	key := createSessionKey(id)

	return r.Client.Del(ctx, key).Err()
}

func createSessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}
