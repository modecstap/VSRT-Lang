package user_repository

import (
	"VSRT-Lang/internal/user"
	"database/sql"
	"errors"
)

func (r *Repository) FindByID(id string) (*user.User, error) {
	found := &user.User{}
	var media sql.NullString
	err := r.db.QueryRow(
		`SELECT id, username, email, password, created_at, avatar, avatar_media_type
         FROM users
         WHERE id = $1`,
		id,
	).Scan(
		&found.ID, &found.Username, &found.Email, &found.Password, &found.CreatedAt,
		&found.Avatar.Bytes, &media,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
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
