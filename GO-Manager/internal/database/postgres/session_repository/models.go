package session_repository

import (
	"VSRT-Lang/internal/session"
	"encoding/json"

	"github.com/lib/pq"
)

type dbContext struct {
	Phrase      string `json:"phrase"`
	Translation string `json:"translation"`
}

func toDBContexts(src []session.Context) []dbContext {
	dst := make([]dbContext, len(src))
	for i, c := range src {
		dst[i] = dbContext(c)
	}
	return dst
}

func fromDBContexts(src []dbContext) []session.Context {
	dst := make([]session.Context, len(src))
	for i, c := range src {
		dst[i] = session.Context(c)
	}
	return dst
}

func (r *Repository) applyRecords(s *session.Session) (session.Session, error) {
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
	`, s.ID)
	if err != nil {
		return *s, err
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
			return *s, err
		}

		if len(contextsJSON) != 0 {
			if err := json.Unmarshal(contextsJSON, &dbContexts); err != nil {
				return *s, err
			}
			rec.Contexts = fromDBContexts(dbContexts)
		}

		s.Records = append(s.Records, rec)
	}

	if err := rows.Err(); err != nil {
		return *s, err
	}

	return *s, nil
}
