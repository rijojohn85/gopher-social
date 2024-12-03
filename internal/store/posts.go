package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	CreateAt  string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Tags      []string  `json:"tags"`
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(
	ctx context.Context,
	post *Post,
) error {
	query := `
  INSERT INTO posts(content, title, user_id, tags)
  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreateAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
UPDATE posts
SET content = $1, title = $2, updated_at = NOW(), tags = $3, version=version + 1
WHERE id = $4 and version=$5
RETURNING version
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		pq.Array(post.Tags),
		post.ID,
		post.Version,
	).Scan(&post.Version)
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

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `
Delete FROM posts
WHERE id = $1;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	result, err := s.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrorNotFound
	} else if count != 1 {
		return errors.New("update: rows affected != 1")
	}
	return nil
}

func (s *PostStore) GetPostById(ctx context.Context, post *Post, id int64) error {
	query := `
  SELECT id, user_id, title, content, created_at, updated_at, tags, version
  FROM posts
  where id = $1
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreateAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(
			err,
			sql.ErrNoRows,
		):
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}
