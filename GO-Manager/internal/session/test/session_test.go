package test

import (
	"errors"
	"strings"
	"testing"

	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/session"
)

func findRecord(records []session.Record, phrase string) (session.Record, bool) {
	for _, record := range records {
		if strings.EqualFold(record.Phrase, phrase) {
			return record, true
		}
	}
	return session.Record{}, false
}

func mustFindRecord(t *testing.T, records []session.Record, phrase string) session.Record {
	t.Helper()
	record, ok := findRecord(records, phrase)
	if !ok {
		t.Fatalf("record %q not found", phrase)
	}
	return record
}

func countContext(record session.Record, phrase string) int {
	n := 0
	for _, context := range record.Contexts {
		if context.Phrase == phrase {
			n++
		}
	}
	return n
}

type failTranslator struct {
	err error
}

func (t failTranslator) Translate(string) ([]string, error) { return nil, t.err }
func (t failTranslator) TranslateBulk([]string) ([][]string, error) {
	return nil, t.err
}
func (t failTranslator) TakeContexts(string) ([]session.Context, error) {
	return nil, t.err
}
func (t failTranslator) TakeSynonyms(string) ([]string, error) { return nil, t.err }
func (t failTranslator) TakeAntonyms(string) ([]string, error) { return nil, t.err }
func (t failTranslator) TakeBaseForm(string) (string, error)   { return "", t.err }

func TestFirstSaveCreatesRecord(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	mustAddRecord(t, svc, userAnna, id, "hello", "hello world")

	records := recordsOf(t, svc, userAnna, id)
	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}

	record := mustFindRecord(t, records, "hello")
	if record.Count != 1 {
		t.Fatalf("count = %d, want 1", record.Count)
	}
	if len(record.Translations) == 0 || record.BaseForm == "" ||
		len(record.Synonyms) == 0 || len(record.Antonyms) == 0 {
		t.Fatalf("record missing translation data: %+v", record)
	}
	if countContext(record, "hello world") == 0 {
		t.Fatalf("missing context %q in %+v", "hello world", record.Contexts)
	}
}

func TestEmptyContextAllowedOnCreate(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	mustAddRecord(t, svc, userAnna, id, "hello", "")

	records := recordsOf(t, svc, userAnna, id)
	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}

	record := mustFindRecord(t, records, "hello")
	if record.Count != 1 {
		t.Fatalf("count = %d, want 1", record.Count)
	}
	for _, context := range record.Contexts {
		if context.Phrase == "" {
			t.Fatal("empty context was stored")
		}
	}
}

func TestRepeatSaveIncrementsCount(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")
	mustAddRecord(t, svc, userAnna, id, "hello", "")

	mustAddRecord(t, svc, userAnna, id, "hello", "")

	records := recordsOf(t, svc, userAnna, id)
	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}
	if records[0].Count != 2 {
		t.Fatalf("count = %d, want 2", records[0].Count)
	}
}

func TestRepeatSaveWithEmptyContextIncrementsCount(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")
	mustAddRecord(t, svc, userAnna, id, "hello", "hello world")

	mustAddRecord(t, svc, userAnna, id, "hello", "")

	record := mustFindRecord(t, recordsOf(t, svc, userAnna, id), "hello")
	if record.Count != 2 {
		t.Fatalf("count = %d, want 2", record.Count)
	}
}

func TestNewContextAddedToSameRecord(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")
	mustAddRecord(t, svc, userAnna, id, "hello", "hello world")

	mustAddRecord(t, svc, userAnna, id, "hello", "say hello")

	records := recordsOf(t, svc, userAnna, id)
	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}

	record := records[0]
	if record.Count != 2 {
		t.Fatalf("count = %d, want 2", record.Count)
	}
	if countContext(record, "hello world") == 0 || countContext(record, "say hello") == 0 {
		t.Fatalf("missing contexts in %+v", record.Contexts)
	}
}

func TestDuplicateContextNotRepeated(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")
	mustAddRecord(t, svc, userAnna, id, "hello", "hello world")

	mustAddRecord(t, svc, userAnna, id, "hello", "hello world")

	records := recordsOf(t, svc, userAnna, id)
	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}

	record := records[0]
	if record.Count != 2 {
		t.Fatalf("count = %d, want 2", record.Count)
	}
	if got := countContext(record, "hello world"); got != 1 {
		t.Fatalf("context count = %d, want 1", got)
	}
}

func TestPhraseCaseDoesNotCreateNewRecord(t *testing.T) {
	for _, phrase := range []string{"hello", "HELLO", "HeLLo"} {
		t.Run(phrase, func(t *testing.T) {
			svc := setupService()
			id := mustCreateSession(t, svc, userAnna, "Урок 1")
			mustAddRecord(t, svc, userAnna, id, "Hello", "")

			mustAddRecord(t, svc, userAnna, id, phrase, "")

			records := recordsOf(t, svc, userAnna, id)
			if len(records) != 1 {
				t.Fatalf("got %d records, want 1", len(records))
			}
			if records[0].Count != 2 {
				t.Fatalf("count = %d, want 2", records[0].Count)
			}
		})
	}
}

func TestEmptyPhraseRejected(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	_, err := svc.AddRecord(session.AddRecordCommand{
		UserId:    userAnna,
		SessionId: id,
		Phrase:    "",
	})
	if !errors.Is(err, session.ErrPhraseRequired) {
		t.Fatalf("error = %v, want %v", err, session.ErrPhraseRequired)
	}

	if got := len(recordsOf(t, svc, userAnna, id)); got != 0 {
		t.Fatalf("got %d records, want 0", got)
	}
}

func TestListSessionRecords(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")
	mustAddRecord(t, svc, userAnna, id, "hello", "")
	mustAddRecord(t, svc, userAnna, id, "world", "")

	records := recordsOf(t, svc, userAnna, id)
	if _, ok := findRecord(records, "hello"); !ok {
		t.Fatal("missing record hello")
	}
	if _, ok := findRecord(records, "world"); !ok {
		t.Fatal("missing record world")
	}
}

func TestDeleteRecord(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")
	mustAddRecord(t, svc, userAnna, id, "hello", "")

	if err := svc.DeleteRecord(userAnna, id, "hello"); err != nil {
		t.Fatalf("DeleteRecord: %v", err)
	}

	if _, ok := findRecord(recordsOf(t, svc, userAnna, id), "hello"); ok {
		t.Fatal("record hello still exists")
	}
}

func TestDeleteMissingRecordSucceeds(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	if err := svc.DeleteRecord(userAnna, id, "hello"); err != nil {
		t.Fatalf("DeleteRecord: %v", err)
	}

	if _, ok := findRecord(recordsOf(t, svc, userAnna, id), "hello"); ok {
		t.Fatal("record hello exists")
	}
}

func TestSaveRecordToForeignSession(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userOther, "Чужой урок")

	_, err := svc.AddRecord(session.AddRecordCommand{
		UserId:    userAnna,
		SessionId: id,
		Phrase:    "hello",
	})
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("error = %v, want %v", err, session.ErrSessionNotFound)
	}

	if _, ok := findRecord(recordsOf(t, svc, userOther, id), "hello"); ok {
		t.Fatal("record appeared in foreign session")
	}
}

func TestGetRecordsOfForeignSession(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userOther, "Чужой урок")
	mustAddRecord(t, svc, userOther, id, "hello", "")

	_, err := svc.GetSession(userAnna, id)
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("error = %v, want %v", err, session.ErrSessionNotFound)
	}
}

func TestDeleteRecordFromForeignSession(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userOther, "Чужой урок")
	mustAddRecord(t, svc, userOther, id, "hello", "")

	err := svc.DeleteRecord(userAnna, id, "hello")
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("error = %v, want %v", err, session.ErrSessionNotFound)
	}

	if _, ok := findRecord(recordsOf(t, svc, userOther, id), "hello"); !ok {
		t.Fatal("foreign record was deleted")
	}
}

func TestSaveRecordWhenTranslatorUnavailable(t *testing.T) {
	svc := session.NewService(memory.NewSessionRepository(), failTranslator{err: errServiceUnavailable})
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	_, err := svc.AddRecord(session.AddRecordCommand{
		UserId:    userAnna,
		SessionId: id,
		Phrase:    "hello",
	})
	if err == nil {
		t.Fatal("expected translator error")
	}

	if got := len(recordsOf(t, svc, userAnna, id)); got != 0 {
		t.Fatalf("got %d records, want 0", got)
	}
}
