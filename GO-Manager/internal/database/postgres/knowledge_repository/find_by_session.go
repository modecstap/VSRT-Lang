package knowledge_repository

import (
	"VSRT-Lang/internal/card"
	"log/slog"
)

// FindBySession implements [card.Repository].
func (r *Repository) FindBySession(sessionID int64) ([]card.Knowledge, error) {
	rows, err := r.db.Query(`
		SELECT
			phrase_key,
			ease_factor,
			repetitions,
			interval_days,
			due_at
		FROM record_knowledge
		WHERE session_id = $1
		ORDER BY phrase_key
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	known := make([]card.Knowledge, 0)
	for rows.Next() {
		k := card.Knowledge{ID: card.KnowledgeId{SessionID: sessionID}}
		if err := rows.Scan(
			&k.ID.PhraseKey,
			&k.EaseFactor,
			&k.Repetitions,
			&k.IntervalDays,
			&k.DueAt,
		); err != nil {
			return nil, err
		}
		k.DueAt = k.DueAt.UTC()
		known = append(known, k)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	slog.Info("record knowledge found by session",
		"session_id", sessionID,
		"count", len(known),
	)
	return known, nil
}
