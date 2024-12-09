package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNotFound         = errors.New("record not found")
	QueryTimeOut          = time.Second * 5
	ErrUserAlreadyFollows = errors.New("user already follows")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateUsername  = errors.New("duplicate username")
	ErrInvitationExpired  = errors.New("invitation expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrEmailNotConfirmed  = errors.New("email not confirmed")
)

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
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		Delete(ctx context.Context, userID int64) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
	Roles interface {
		GetIDByName(context.Context, string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
		Roles:    &RoleStore{db},
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
