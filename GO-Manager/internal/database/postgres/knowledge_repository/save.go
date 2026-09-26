package knowledge_repository

import (
	"VSRT-Lang/internal/card"
	"log/slog"
)

// Save implements [card.Repository].
func (r *Repository) Save(k *card.Knowledge) (card.KnowledgeId, error) {
	_, err := r.db.Exec(
		`INSERT INTO record_knowledge (session_id, phrase_key, ease_factor, repetitions, interval_days, due_at)
         VALUES ($1, $2, $3, $4, $5, $6)
         ON CONFLICT (session_id, phrase_key) DO UPDATE SET
             ease_factor   = excluded.ease_factor,
             repetitions   = excluded.repetitions,
             interval_days = excluded.interval_days,
             due_at        = excluded.due_at`,
		k.ID.SessionID,
		k.ID.PhraseKey,
		k.EaseFactor,
		k.Repetitions,
		k.IntervalDays,
		k.DueAt,
	)
	if err != nil {
		return card.KnowledgeId{}, err
	}

	slog.Info("record knowledge saved",
		"session_id", k.ID.SessionID,
		"phrase_key", k.ID.PhraseKey,
		"ease_factor", k.EaseFactor,
		"repetitions", k.Repetitions,
		"interval_days", k.IntervalDays,
		"due_at", k.DueAt,
	)
	return k.ID, nil
}
