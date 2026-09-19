package net_translator_test

import (
	"testing"

	"VSRT-Lang/internal/net_translator"
)

func TestNewTranslator(t *testing.T) {
	translator := net_translator.NewTranslator("localhost:8081")
	if translator == nil {
		t.Fatalf("Failed to create translator")
	}
}

func newTranslator(t *testing.T) *net_translator.Translator {
	translator := net_translator.NewTranslator("localhost:8081")
	if translator == nil {
		t.Fatalf("Failed to create translator")
	}
	return translator
}

func TestTranslate(t *testing.T) {
	translator := newTranslator(t)
	_, err := translator.Translate("Hello, world!")
	if err != nil {
		t.Fatalf("Failed to translate: %v", err)
	}
}

func TestEmptyTranslation(t *testing.T) {
	translator := newTranslator(t)
	_, err := translator.Translate("")
	if err != nil {
		t.Fatalf("Failed to translate: %v", err)
	}
}

func TestBulkTranslation(t *testing.T) {
	translator := newTranslator(t)
	translations, err := translator.TranslateBulk([]string{"Hello, world!", "Hello, world!", "Hello, world!"})
	if err != nil {
		t.Fatalf("Failed to translate: %v", err)
	}
	if len(translations) != 3 {
		t.Fatalf("Expected 3 translations, got %d", len(translations))
	}
}

func TestEmptyBulkTranslation(t *testing.T) {
	translator := newTranslator(t)
	translations, err := translator.TranslateBulk([]string{})
	if err != nil {
		t.Fatalf("Failed to translate: %v", err)
	}
	if len(translations) != 0 {
		t.Fatalf("Expected 0 translations, got %d", len(translations))
	}
}

func TestBulkTranslationEmptyPhrase(t *testing.T) {
	translator := newTranslator(t)
	translations, err := translator.TranslateBulk([]string{"", "test", "result"})
	if err != nil {
		t.Fatalf("Failed to translate: %v", err)
	}
	if len(translations) != 3 {
		t.Fatalf("Expected 3 translations, got %d", len(translations))
	}
}
