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
	User      User      `json:"user"`
}

type UserFeed struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, id int64, fq PaginatedFeedQuery) ([]UserFeed, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
			f.user_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	rows, err := s.db.QueryContext(
		ctx,
		query,
		id,
		fq.Limit,
		fq.Offset,
		//fq.Search,
		//pq.Array(fq.Tags),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feed []UserFeed
	for rows.Next() {
		var post UserFeed
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreateAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		); err != nil {
			return nil, err
		}
		feed = append(feed, post)
	}
	return feed, nil
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
