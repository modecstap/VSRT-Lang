package user_repository

import (
	"VSRT-Lang/internal/user"
	"database/sql"
	"errors"
)

func (r *Repository) FindByID(id string) (*user.User, error) {
	user := &user.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, created_at
         FROM users
         WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
