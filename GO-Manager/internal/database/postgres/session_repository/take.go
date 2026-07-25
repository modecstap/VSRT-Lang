package session_repository

import (
	"VSRT-Lang/internal/session"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lib/pq"
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

	rows, err := r.db.Query(`
		SELECT
			phrase,
			base_form,
			translations,
			synonyms,
			antonyms,
			contexts
		FROM records
		WHERE session_id = $1
		ORDER BY id
	`, sessionID)
	if err != nil {
		return s, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec session.Record

		var (
			contextsJSON []byte
			dbContexts   []dbContext
		)

		err = rows.Scan(
			&rec.Phrase,
			&rec.BaseForm,
			pq.Array(&rec.Translations),
			pq.Array(&rec.Synonyms),
			pq.Array(&rec.Antonyms),
			&contextsJSON,
		)
		if err != nil {
			return s, err
		}

		if len(contextsJSON) != 0 {
			if err := json.Unmarshal(contextsJSON, &dbContexts); err != nil {
				return s, err
			}
			rec.Contexts = fromDBContexts(dbContexts)
		}

		s.Records = append(s.Records, rec)
	}

	if err := rows.Err(); err != nil {
		return s, err
	}

	return s, nil
}
