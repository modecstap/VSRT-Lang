package session_repository

import (
	"VSRT-Lang/internal/session"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
)

func (r *Repository) Save(s *session.Session) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = r.saveSession(s, tx)
	if err != nil {
		return err
	}

	err = r.updateRecords(s, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) updateRecords(s *session.Session, tx *sql.Tx) error {
	_, err := tx.Exec(
		`DELETE FROM records WHERE session_id=$1`,
		s.ID,
	)

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT INTO records(
            session_id,
            phrase,
            base_form,
            translations,
            synonyms,
            antonyms,
            contexts
        )
        VALUES($1,$2,$3,$4,$5,$6,$7)
    `)

	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	for _, rec := range s.Records {

		ctxJSON, err := json.Marshal(toDBContexts(rec.Contexts))
		if err != nil {
			return err
		}

		_, err = stmt.Exec(
			s.ID,
			rec.Phrase,
			rec.BaseForm,
			pq.Array(rec.Translations),
			pq.Array(rec.Synonyms),
			pq.Array(rec.Antonyms),
			ctxJSON,
		)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) saveSession(s *session.Session, tx *sql.Tx) error {
	var exists bool

	err := tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM sessions
			WHERE id = $1
		)
	`, s.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = tx.Exec(`
			UPDATE sessions
			SET user_id = $1,
			    name = $2
			WHERE id = $3
		`, s.User, s.Name, s.ID)

		return err
	}

	err = tx.QueryRow(`
		INSERT INTO sessions (user_id, name)
		VALUES ($1, $2)
		RETURNING id
	`, s.User, s.Name).Scan(&s.ID)

	return err
}
