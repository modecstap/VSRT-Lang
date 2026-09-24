package user_repository

import (
	"VSRT-Lang/internal/user"
	"database/sql"
	"errors"
)

func (r *Repository) FindByUsername(username string) (*user.User, error) {
	found := &user.User{}
	var media sql.NullString
	err := r.db.QueryRow(
		`SELECT id, username, email, password, created_at, avatar, avatar_media_type
         FROM users
         WHERE username = $1`,
		username,
	).Scan(
		&found.ID, &found.Username, &found.Email, &found.Password, &found.CreatedAt,
		&found.Avatar.Bytes, &media,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	if found.Avatar.Bytes == nil {
		found.Avatar = user.Avatar{}
	} else {
		found.Avatar.MediaType = media.String
	}
	return found, nil
}
