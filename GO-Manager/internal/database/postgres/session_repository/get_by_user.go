package session_repository

import (
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
	"database/sql"
	"errors"
)

func (r *Repository) FindByUser(userId user.UserId) ([]session.Session, error) {
	var sessions []session.Session

	rows, err := r.db.Query(`
		SELECT
			id,
			user_id,
			name
		FROM sessions
		WHERE user_id = $1
	`, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		var s session.Session

		err = rows.Scan(
			&s.ID,
			&s.User,
			&s.Name,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			return nil, err
		}

		s, err = r.applyRecords(&s)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, s)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return sessions, nil
}

