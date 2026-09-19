package test

import (
	"reflect"
	"testing"

	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
)

func TestNewRecord(t *testing.T) {
	translator := stub.Translator{}
	record, err := session.NewRecord(
		session.NewRecordCommand{
			Translator:  translator,
			Phrase:      "hello",
			MainContext: "world",
		},
	)
	if err != nil {
		t.Fatalf("NewRecord returned unexpected error: %v", err)
	}

	translations, err := translator.TranslateBulk([]string{"hello", "world"})
	if err != nil {
		t.Fatalf("translator.TranslateBulk returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(record.Translations, translations[0]) {
		t.Fatalf(
			"expected translations %v, got %v",
			translations[0], record.Translations,
		)
	}

	if record.BaseForm != "base-form:hello" {
		t.Fatalf(
			"expected base form base-form:hello, got %q",
			record.BaseForm,
		)
	}

	synonyms, err := translator.TakeSynonyms("hello")
	if err != nil {
		t.Fatalf("translator.TakeSynonyms returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(record.Synonyms, synonyms) {
		t.Fatalf(
			"expected synonyms %v, got %v",
			synonyms, record.Synonyms,
		)
	}

	antonyms, err := translator.TakeAntonyms("hello")
	if err != nil {
		t.Fatalf("translator.TakeAntonyms returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(record.Antonyms, antonyms) {
		t.Fatalf(
			"expected antonyms %v, got %v",
			antonyms, record.Antonyms,
		)
	}

	contexts, err := translator.TakeContexts("hello")
	if err != nil {
		t.Fatalf("translator.TakeContexts returned unexpected error: %v", err)
	}
	if len(record.Contexts) != len(contexts)+1 {
		t.Fatalf(
			"expected %d contexts, got %d",
			len(contexts)+1, len(record.Contexts),
		)
	}

	if record.Contexts[0].Phrase != "world" {
		t.Fatalf("expected first context phrase 'world', got %q", record.Contexts[0].Phrase)
	}

	for i, context := range contexts {
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
	translator := stub.Translator{}
	session := session.NewSession(user.UserId("user-123"), "test-session")

	record, err := session.SaveRecord("world", "hello", translator)
	if err != nil {
		t.Fatalf("SaveRecord returned unexpected error: %v", err)
	}

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
	translator := stub.Translator{}
	session := session.NewSession(user.UserId("user-123"), "test-session")
	_, err := session.SaveRecord("world", "hello", translator)
	if err != nil {
		t.Fatalf("SaveRecord returned unexpected error: %v", err)
	}

	records := session.GetRecords()
	if records[0].Count != 1 {
		t.Fatalf("expected record count to be 1, got %d", records[0].Count)
	}
}

func TestSessionSaveRecordWithExistingPhrase(t *testing.T) {
	translator := stub.Translator{}
	session := session.NewSession(user.UserId("user-123"), "test-session")
	_, err := session.SaveRecord("world", "hello", translator)
	if err != nil {
		t.Fatalf("SaveRecord returned unexpected error: %v", err)
	}
	_, err = session.SaveRecord("world", "hello", translator)
	if err != nil {
		t.Fatalf("SaveRecord returned unexpected error on second save: %v", err)
	}

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
