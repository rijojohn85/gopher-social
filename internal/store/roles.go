package store

import (
	"context"
	"database/sql"
)

type Role struct {
	RoleName string
	RoleID   int64
}

type RoleStore struct {
	db *sql.DB
}

func (r *RoleStore) GetIDByName(ctx context.Context, role_string string) (*Role, error) {
	role := &Role{}

	query := `SELECT id, name from roles where name=$1`

	err := r.db.QueryRowContext(ctx, query, role_string).Scan(&role.RoleID, &role.RoleName)
	if err != nil {
		return nil, err
	}
	return role, nil
}
