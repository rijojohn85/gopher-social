package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrorNotFound = errors.New("record not found")
var QueryTimeOut = time.Second * 5
var ErrUserAlreadyFollows = errors.New("user already follows")
var ErrDuplicateEmail = errors.New("duplicate email")
var ErrDuplicateUsername = errors.New("duplicate username")
var ErrInvitationExpired = errors.New("invitation expired")
var ErrInvalidToken = errors.New("invalid token")

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostById(context.Context, *Post, int64) error
		Update(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, id int64, fq PaginatedFeedQuery) ([]UserFeed, error)
		Delete(ctx context.Context, id int64) error
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetUser(context.Context, *User, int64) error
		AddFollower(ctx context.Context, userID, followerID int64) error
		DeleteFollower(ctx context.Context, userID, followerID int64) error
		CreateAndInvite(
			ctx context.Context,
			user *User,
			token string,
			exp time.Duration,
		) error
		Activate(ctx context.Context, token string) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}

func withTx(ctx context.Context, db *sql.DB, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := f(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
