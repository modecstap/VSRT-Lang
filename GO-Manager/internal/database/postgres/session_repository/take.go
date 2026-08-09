package session_repository

import (
	"VSRT-Lang/internal/session"
	"database/sql"
	"errors"
)

func (r *Repository) Take(sessionID int) (session.Session, error) {
	var s session.Session

	err := r.db.QueryRow(`
		SELECT
			id,
			user_id,
			name
		FROM sessions
		WHERE id = $1
	`, sessionID).Scan(
		&s.ID,
		&s.User,
		&s.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s, err
		}
		return s, err
	}

	s, err = r.applyRecords(&s)

	return s, nil
}
