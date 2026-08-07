package user_repository

import (
	"VSRT-Lang/internal/user"
)

func (r *Repository) Create(user *user.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, username, email, password, created_at)
         VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
	)
	return err
}
