package cache

import (
	"context"

	"github.com/rijojohn85/social/internal/store"
)

func NewMockCache() Storage {
	return Storage{
		Users: &MockUserCache{},
	}
}

type MockUserCache struct{}

func (m *MockUserCache) Get(context.Context, *store.User, int64) error {
	return nil
}

func (m *MockUserCache) Set(context.Context, *store.User) error {
	return nil
}
