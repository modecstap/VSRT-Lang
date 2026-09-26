package knowledge_repository

import (
	"VSRT-Lang/internal/card"
	"database/sql"
	"errors"
	"log/slog"
)

// Take implements [card.Repository].
func (r *Repository) Take(id card.KnowledgeId) (*card.Knowledge, error) {
	k := card.Knowledge{ID: id}

	err := r.db.QueryRow(`
		SELECT
			ease_factor,
			repetitions,
			interval_days,
			due_at
		FROM record_knowledge
		WHERE session_id = $1 AND phrase_key = $2
	`, id.SessionID, id.PhraseKey).Scan(
		&k.EaseFactor,
		&k.Repetitions,
		&k.IntervalDays,
		&k.DueAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	k.DueAt = k.DueAt.UTC()

	slog.Info("record knowledge taken",
		"session_id", k.ID.SessionID,
		"phrase_key", k.ID.PhraseKey,
		"ease_factor", k.EaseFactor,
		"repetitions", k.Repetitions,
		"interval_days", k.IntervalDays,
		"due_at", k.DueAt,
	)
	return &k, nil
}
