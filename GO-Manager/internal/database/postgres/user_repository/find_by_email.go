package user_repository

import (
	"VSRT-Lang/internal/auth"
	"database/sql"
	"errors"
)

func (r *Repository) FindByEmail(email string) (*auth.User, error) {
	user := &auth.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, created_at
         FROM users
         WHERE email = $1`,
		email,
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