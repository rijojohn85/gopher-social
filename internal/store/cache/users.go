package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rijojohn85/social/internal/store"
)

type UserStore struct {
	rdb *redis.Client
}

func (r *UserStore) Get(ctx context.Context, user *store.User, userID int64) error {
	cacheKey := fmt.Sprintf("user-%d", userID)
	data, err := r.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return err
	}
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf(
		"user-%d",
		user.ID,
	)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.rdb.SetEX(ctx, cacheKey, json, time.Minute*10).Err()
}
