package test

import (
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

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

type overlapTranslator struct {
	mu                        sync.Mutex
	entered                   int
	release                   chan struct{}
	ctxPhrase, basePhrase     string
	synPhrase, antPhrase      string
	bulkCalls, translateCalls int
}

func newOverlapTranslator() *overlapTranslator {
	return &overlapTranslator{release: make(chan struct{})}
}

func (f *overlapTranslator) gate(phrase string, seen *string) error {
	f.mu.Lock()
	*seen = phrase
	f.entered++
	if f.entered == 4 {
		close(f.release)
	}
	f.mu.Unlock()
	select {
	case <-f.release:
		return nil
	case <-time.After(time.Second):
		return errors.New("description calls did not overlap")
	}
}

func (f *overlapTranslator) Translate(string) ([]string, error) {
	f.mu.Lock()
	f.translateCalls++
	f.mu.Unlock()
	return []string{"unused"}, nil
}

func (f *overlapTranslator) TranslateBulk(phrases []string) ([][]string, error) {
	f.mu.Lock()
	f.bulkCalls++
	f.mu.Unlock()
	out := [][]string{{"phrase-tr"}}
	if len(phrases) == 2 {
		out = append(out, []string{"context-tr"})
	}
	return out, nil
}

func (f *overlapTranslator) TakeContexts(phrase string) ([]session.Context, error) {
	if err := f.gate(phrase, &f.ctxPhrase); err != nil {
		return nil, err
	}
	return []session.Context{{Phrase: "extra", Translation: "extra-tr"}}, nil
}

func (f *overlapTranslator) TakeBaseForm(phrase string) (string, error) {
	if err := f.gate(phrase, &f.basePhrase); err != nil {
		return "", err
	}
	return "base", nil
}

func (f *overlapTranslator) TakeSynonyms(phrase string) ([]string, error) {
	if err := f.gate(phrase, &f.synPhrase); err != nil {
		return nil, err
	}
	return []string{"syn"}, nil
}

func (f *overlapTranslator) TakeAntonyms(phrase string) ([]string, error) {
	if err := f.gate(phrase, &f.antPhrase); err != nil {
		return nil, err
	}
	return []string{"ant"}, nil
}

type scenarioTranslator struct {
	mu                                       sync.Mutex
	bulkErr, ctxErr, baseErr, synErr, antErr error
	bulkN, trN, ctxN, baseN, synN, antN      int
}

func (f *scenarioTranslator) hit(n *int) {
	f.mu.Lock()
	*n++
	f.mu.Unlock()
}

func (f *scenarioTranslator) Translate(string) ([]string, error) {
	f.hit(&f.trN)
	return []string{"ctx-tr"}, nil
}

func (f *scenarioTranslator) TranslateBulk(phrases []string) ([][]string, error) {
	f.hit(&f.bulkN)
	if f.bulkErr != nil {
		return nil, f.bulkErr
	}
	out := [][]string{{"phrase-tr"}}
	if len(phrases) == 2 {
		out = append(out, []string{"context-tr"})
	}
	return out, nil
}

func (f *scenarioTranslator) TakeContexts(string) ([]session.Context, error) {
	f.hit(&f.ctxN)
	return nil, f.ctxErr
}

func (f *scenarioTranslator) TakeBaseForm(string) (string, error) {
	f.hit(&f.baseN)
	return "", f.baseErr
}

func (f *scenarioTranslator) TakeSynonyms(string) ([]string, error) {
	f.hit(&f.synN)
	return nil, f.synErr
}

func (f *scenarioTranslator) TakeAntonyms(string) ([]string, error) {
	f.hit(&f.antN)
	return nil, f.antErr
}

func (f *scenarioTranslator) descEach(n int) bool {
	return f.ctxN == n && f.baseN == n && f.synN == n && f.antN == n
}

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

func TestDescriptionCallsOverlap(t *testing.T) {
	tr := newOverlapTranslator()
	svc := session.NewService(memory.NewSessionRepository(), tr)
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	mustAddRecord(t, svc, userAnna, id, "  hello  ", "user context")

	if tr.ctxPhrase != "  hello  " || tr.basePhrase != "  hello  " ||
		tr.synPhrase != "  hello  " || tr.antPhrase != "  hello  " {
		t.Fatalf("phrases = %q %q %q %q", tr.ctxPhrase, tr.basePhrase, tr.synPhrase, tr.antPhrase)
	}
	if tr.bulkCalls != 1 || tr.translateCalls != 0 {
		t.Fatalf("bulk = %d, translate = %d", tr.bulkCalls, tr.translateCalls)
	}

	record := mustFindRecord(t, recordsOf(t, svc, userAnna, id), "hello")
	if record.Phrase != "hello" || record.Count != 1 {
		t.Fatalf("phrase = %q, count = %d", record.Phrase, record.Count)
	}
	if len(record.Translations) != 1 || record.Translations[0] != "phrase-tr" {
		t.Fatalf("translations = %v", record.Translations)
	}
	if len(record.Contexts) != 2 ||
		record.Contexts[0] != (session.Context{Phrase: "user context", Translation: "context-tr"}) ||
		record.Contexts[1] != (session.Context{Phrase: "extra", Translation: "extra-tr"}) {
		t.Fatalf("contexts = %+v", record.Contexts)
	}
	if record.BaseForm != "base" || len(record.Synonyms) != 1 || record.Synonyms[0] != "syn" ||
		len(record.Antonyms) != 1 || record.Antonyms[0] != "ant" {
		t.Fatalf("record = %+v", record)
	}
}

func TestDescriptionErrorDropsRecord(t *testing.T) {
	errBulk := errors.New("bulk")
	errCtx := errors.New("contexts")
	errBase := errors.New("base")
	errSyn := errors.New("syn")
	errAnt := errors.New("ant")

	tests := []struct {
		name string
		tr   *scenarioTranslator
		want error
		n    int
	}{
		{"bulk", &scenarioTranslator{bulkErr: errBulk}, errBulk, 0},
		{"contexts", &scenarioTranslator{ctxErr: errCtx, baseErr: errBase}, errCtx, 1},
		{"base", &scenarioTranslator{baseErr: errBase, synErr: errSyn}, errBase, 1},
		{"synonyms", &scenarioTranslator{synErr: errSyn, antErr: errAnt}, errSyn, 1},
		{"antonyms", &scenarioTranslator{antErr: errAnt}, errAnt, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := session.NewService(memory.NewSessionRepository(), tt.tr)
			id := mustCreateSession(t, svc, userAnna, "Урок 1")
			_, err := svc.AddRecord(session.AddRecordCommand{
				UserId:    userAnna,
				SessionId: id,
				Phrase:    "hello",
				Context:   "c",
			})
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
			if got := len(recordsOf(t, svc, userAnna, id)); got != 0 {
				t.Fatalf("got %d records, want 0", got)
			}
			if !tt.tr.descEach(tt.n) {
				t.Fatalf("calls ctx=%d base=%d syn=%d ant=%d, want %d",
					tt.tr.ctxN, tt.tr.baseN, tt.tr.synN, tt.tr.antN, tt.n)
			}
		})
	}
}

func TestRepeatSaveSkipsDescription(t *testing.T) {
	tr := &scenarioTranslator{}
	svc := session.NewService(memory.NewSessionRepository(), tr)
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	mustAddRecord(t, svc, userAnna, id, "hello", "c1")
	if tr.bulkN != 1 || tr.trN != 0 || !tr.descEach(1) {
		t.Fatalf("after first: bulk=%d translate=%d ctx=%d base=%d syn=%d ant=%d",
			tr.bulkN, tr.trN, tr.ctxN, tr.baseN, tr.synN, tr.antN)
	}

	mustAddRecord(t, svc, userAnna, id, "hello", "c2")
	if tr.trN != 1 || tr.bulkN != 1 || !tr.descEach(1) {
		t.Fatalf("after c2: bulk=%d translate=%d ctx=%d base=%d syn=%d ant=%d",
			tr.bulkN, tr.trN, tr.ctxN, tr.baseN, tr.synN, tr.antN)
	}

	mustAddRecord(t, svc, userAnna, id, "hello", "c2")
	if tr.trN != 1 {
		t.Fatalf("translate = %d, want 1", tr.trN)
	}

	record := mustAddRecord(t, svc, userAnna, id, "hello", "")
	if tr.trN != 1 || !tr.descEach(1) || record.Count != 4 {
		t.Fatalf("translate=%d count=%d ctx=%d base=%d syn=%d ant=%d",
			tr.trN, record.Count, tr.ctxN, tr.baseN, tr.synN, tr.antN)
	}
}
