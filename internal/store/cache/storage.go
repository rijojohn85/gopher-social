package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/rijojohn85/social/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, *store.User, int64) error
		Set(context.Context, *store.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
