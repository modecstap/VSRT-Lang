package session

import (
	"reflect"
	"testing"

	"VSRT-Lang/internal/user"
)

type stubTranslator struct{}

func (stubTranslator) Translate(phrase string) []string {
	return []string{"translated:" + phrase}
}

func (stubTranslator) TranslateBulk(phrases []string) [][]string {
	return [][]string{
		{"phrase-translation-1", "phrase-translation-2"},
		{"context-translation"},
	}
}

func (stubTranslator) TakeContexts(phrase string) []Context {
	return []Context{{Phrase: "additional context", Translation: "additional-translation"}}
}

func (stubTranslator) TakeSynonyms(word string) []string {
	return []string{"synonym1", "synonym2"}
}

func (stubTranslator) TakeAntonyms(word string) []string {
	return []string{"antonym1", "antonym2"}
}

func (stubTranslator) TakeBaseForm(word string) string {
	return "base-form:" + word
}

func TestNewRecord(t *testing.T) {
	translator := stubTranslator{}
	record := NewRecord(translator, "hello", "world")

	if record.Phrase != "hello" {
		t.Fatalf("expected phrase hello, got %q", record.Phrase)
	}

	expectedTranslations := []string{"phrase-translation-1", "phrase-translation-2"}
	if !reflect.DeepEqual(record.Translations, expectedTranslations) {
		t.Fatalf("expected translations %v, got %v", expectedTranslations, record.Translations)
	}

	if record.BaseForm != "base-form:hello" {
		t.Fatalf("expected base form base-form:hello, got %q", record.BaseForm)
	}

	if !reflect.DeepEqual(record.Synonyms, []string{"synonym1", "synonym2"}) {
		t.Fatalf("expected synonyms [synonym1 synonym2], got %v", record.Synonyms)
	}

	if !reflect.DeepEqual(record.Antonyms, []string{"antonym1", "antonym2"}) {
		t.Fatalf("expected antonyms [antonym1 antonym2], got %v", record.Antonyms)
	}

	if len(record.Contexts) != 2 {
		t.Fatalf("expected 2 contexts, got %d", len(record.Contexts))
	}

	if record.Contexts[0].Phrase != "world" || record.Contexts[0].Translation != "context-translation" {
		t.Fatalf("unexpected primary context: %+v", record.Contexts[0])
	}

	if record.Contexts[1].Phrase != "additional context" || record.Contexts[1].Translation != "additional-translation" {
		t.Fatalf("unexpected additional context: %+v", record.Contexts[1])
	}
}

func TestSessionSaveAndGetRecords(t *testing.T) {
	translator := stubTranslator{}
	session := NewSession(user.UserId("user-123"), "test-session", translator)

	record := session.SaveRecord("world", "hello", translator)

	if record.Phrase != "hello" {
		t.Fatalf("expected returned record phrase hello, got %q", record.Phrase)
	}

	gotRecords := session.GetRecords()
	if len(gotRecords) != 1 {
		t.Fatalf("expected GetRecords to return 1 record, got %d", len(gotRecords))
	}

	if !reflect.DeepEqual(gotRecords[0], record) {
		t.Fatalf("expected GetRecords result to match record, got %v", gotRecords[0])
	}
}
