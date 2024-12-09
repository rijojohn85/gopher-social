package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(context.Context, *sql.Tx, *User) error {
	return nil
}

func (m *MockUserStore) GetUser(context.Context, *User, int64) error {
	return nil
}

func (m *MockUserStore) AddFollower(ctx context.Context, userID, followerID int64) error {
	return nil
}

func (m *MockUserStore) DeleteFollower(ctx context.Context, userID, followerID int64) error {
	return nil
}

func (m *MockUserStore) CreateAndInvite(
	ctx context.Context,
	user *User,
	token string,
	exp time.Duration,
) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (m *MockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}
