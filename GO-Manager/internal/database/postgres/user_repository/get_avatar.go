package user_repository

import (
	"VSRT-Lang/internal/user"
	"database/sql"
	"errors"
)

func (r *Repository) GetAvatar(id user.UserId) (user.Avatar, error) {
	var (
		bytes     []byte
		mediaType sql.NullString
	)
	err := r.db.QueryRow(
		`SELECT avatar, avatar_media_type FROM users WHERE id = $1`,
		string(id),
	).Scan(&bytes, &mediaType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.Avatar{}, errors.New("user not found")
		}
		return user.Avatar{}, err
	}
	if bytes == nil {
		return user.Avatar{}, user.ErrNoAvatar
	}
	return user.Avatar{Bytes: bytes, MediaType: mediaType.String}, nil
}
