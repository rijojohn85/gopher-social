package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/lib/pq"
)

type User struct {
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	ID        int64    `json:"id"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(input string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.hash = hash
	p.text = &input
	return nil
}

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
  INSERT INTO users(username, password, email) VALUES($1, $2, $3) RETURNING id, created_at
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := tx.QueryRowContext(
		ctx, query,
		user.Username,
		user.Password.hash,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
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

func (s *UserStore) CreateAndInvite(
	ctx context.Context,
	user *User,
	token string,
	exp time.Duration,
) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		//create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		//Create The user invite
		err := s.createUserInvite(ctx, tx, user.ID, token, exp)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) createUserInvite(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	token string,
	exp time.Duration,
) error {

	query := `
insert into user_invitations(token, user_id, expiry) values($1, $2, $3);
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}
func (s *UserStore) Activate(ctx context.Context, token string) error {
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	query := `
SELECT user_id, expiry from user_invitations where token = $1;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, hashToken)
	var userID int64
	var expiryString string
	err := row.Scan(&userID, &expiryString)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrInvalidToken
		default:
			return err
		}
	}
	expiry, err := time.Parse(time.RFC3339, expiryString)
	if err != nil {
		return err
	}
	if time.Now().After(expiry) {
		return ErrInvitationExpired
	}
	query = `
update users set is_active=true where id = $1;
`
	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorNotFound
	}
	query = `
delete from user_invitations where token = $1;
`
	s.db.ExecContext(ctx, query, hashToken)
	return nil
}
