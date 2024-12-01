package store

import (
	"context"
	"database/sql"
	"errors"
)

var ErrorNotFound = errors.New("record not found")

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostById(context.Context, *Post, int64) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
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
