package test

import (
	"reflect"
	"testing"

	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

func TestNewRecord(t *testing.T) {
	translator := stubTranslator{}
	record := session.NewRecord(translator, "hello", "world")

	if record.Phrase != "hello" {
		t.Fatalf("expected phrase hello, got %q", record.Phrase)
	}

	if !reflect.DeepEqual(record.Translations, translator.TranslateBulk([]string{"hello"})[0]) {
		t.Fatalf(
			"expected translations %v, got %v",
			translator.TranslateBulk([]string{"hello"}), record.Translations,
		)
	}

	if record.BaseForm != "base-form:hello" {
		t.Fatalf(
			"expected base form base-form:hello, got %q",
			record.BaseForm,
		)
	}

	if !reflect.DeepEqual(record.Synonyms, translator.TakeSynonyms("hello")) {
		t.Fatalf(
			"expected synonyms %v, got %v",
			translator.TakeSynonyms("hello"), record.Synonyms,
		)
	}

	if !reflect.DeepEqual(record.Antonyms, translator.TakeAntonyms("hello")) {
		t.Fatalf(
			"expected antonyms %v, got %v",
			translator.TakeAntonyms("hello"), record.Antonyms,
		)
	}

	if len(record.Contexts) != len(translator.TakeContexts("hello"))+1 {
		t.Fatalf(
			"expected %d contexts, got %d",
			len(translator.TakeContexts("hello"))+1, len(record.Contexts),
		)
	}

	if record.Contexts[0].Phrase != "world" {
		t.Fatalf("expected first context phrase 'world', got %q", record.Contexts[0].Phrase)
	}

	for i, context := range translator.TakeContexts("hello") {
		if record.Contexts[i+1].Phrase != context.Phrase ||
			record.Contexts[i+1].Translation != context.Translation {
			t.Fatalf(
				"unexpected context at index %d: %+v",
				i, record.Contexts[i+1],
			)
		}
	}
}

func TestSessionSaveAndGetRecords(t *testing.T) {
	translator := stubTranslator{}
	session := session.NewSession(user.UserId("user-123"), "test-session", translator)

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

func TestCountWhenSaveOneRecord(t *testing.T) {
	translator := stubTranslator{}
	session := session.NewSession(user.UserId("user-123"), "test-session", translator)
	_ = session.SaveRecord("world", "hello", translator)
	
	records := session.GetRecords()
	if records[0].Count != 1 {
		t.Fatalf("expected record count to be 1, got %d", records[0].Count)
	}
}

func TestSessionSaveRecordWithExistingPhrase(t *testing.T) {
	translator := stubTranslator{}
	session := session.NewSession(user.UserId("user-123"), "test-session", translator)
	_ = session.SaveRecord("world", "hello", translator)
	_ = session.SaveRecord("world", "hello", translator)

	records := session.GetRecords()

	if len(records) != 1 {
		t.Fatalf(
			"expected GetRecords to return 1 record, got %d",
			len(records),
		)
	}

	if records[0].Count != 2 {
		t.Fatalf(
			"expected record count to be 2, got %d",
			records[0].Count,
		)
	}
}
