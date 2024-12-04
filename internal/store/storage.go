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

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostById(context.Context, *Post, int64) error
		Update(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, id int64) ([]UserFeed, error)
		Delete(ctx context.Context, id int64) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetUser(context.Context, *User, int64) error
		AddFollower(ctx context.Context, userID, followerID int64) error
		DeleteFollower(ctx context.Context, userID, followerID int64) error
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
