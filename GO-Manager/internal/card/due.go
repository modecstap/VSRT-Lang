package card

import (
	"slices"
	"strings"
	"time"

	"VSRT-Lang/internal/session"
)

type dueRecord struct {
	record session.Record
	dueAt  time.Time
}

func DueRecords(records []session.Record, knowledge []Knowledge, now time.Time, limit int) []session.Record {
	byKey := make(map[string]Knowledge, len(knowledge))
	for _, k := range knowledge {
		byKey[k.ID.PhraseKey] = k
	}

	due := make([]dueRecord, 0, len(records))
	for _, rec := range records {
		key := session.RecordKey(rec.Phrase)
		k, ok := byKey[key]
		if !ok {
			k = *NewKnowledge(KnowledgeId{PhraseKey: key})
		}
		if k.DueBy(now) {
			due = append(due, dueRecord{record: rec, dueAt: k.DueAt})
		}
	}

	slices.SortFunc(due, func(a, b dueRecord) int {
		if c := a.dueAt.Compare(b.dueAt); c != 0 {
			return c
		}
		return strings.Compare(a.record.Phrase, b.record.Phrase)
	})

	if limit > 0 && len(due) > limit {
		due = due[:limit]
	}

	result := make([]session.Record, 0, len(due))
	for _, d := range due {
		result = append(result, d.record)
	}
	return result
}
