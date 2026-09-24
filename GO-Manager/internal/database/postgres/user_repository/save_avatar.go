package user_repository

import (
	"VSRT-Lang/internal/user"
)

func (r *Repository) SaveAvatar(id user.UserId, avatar user.Avatar) error {
	result, err := r.db.Exec(
		`UPDATE users SET avatar = $1, avatar_media_type = $2 WHERE id = $3`,
		avatar.Bytes,
		avatar.MediaType,
		string(id),
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return user.ErrNotFound
	}
	return nil
}
