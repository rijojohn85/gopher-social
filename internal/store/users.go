package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
}

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
  INSERT INTO users(username, password, email) VALUES($1, $2, $3) RETURNING id, created_at
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx, query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetUser(ctx context.Context, user *User, id int64) error {
	query := `
		SELECT id, username, email, created_at FROM users WHERE id = $1;
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx, query, id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) AddFollower(ctx context.Context, userID, followerID int64) error {
	query := `
INSERT INTO followers(user_id, follower_id) VALUES($1, $2)
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		if prErr, ok := err.(*pq.Error); ok && prErr.Code == "23505" {
			return ErrUserAlreadyFollows
		} else {
			return err
		}
	}
	return nil
}

func (s *UserStore) DeleteFollower(ctx context.Context, userID, followerID int64) error {
	query := `
DELETE FROM followers WHERE user_id = $1 AND follower_id = $2;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}
