package user_repository

import (
	"VSRT-Lang/internal/user"
)

func (r *Repository) Save(u *user.User) error {
	var avatar any
	var media any
	if len(u.Avatar.Bytes) > 0 {
		avatar = u.Avatar.Bytes
		media = u.Avatar.MediaType
	}

	result, err := r.db.Exec(
		`UPDATE users
         SET username = $1, email = $2, password = $3, avatar = $4, avatar_media_type = $5
         WHERE id = $6`,
		u.Username,
		u.Email,
		u.Password,
		avatar,
		media,
		u.ID,
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
