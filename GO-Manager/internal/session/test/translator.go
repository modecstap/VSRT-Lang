package test

import "VSRT-Lang/internal/session"

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

func (stubTranslator) TakeContexts(phrase string) []session.Context {
	return []session.Context{{Phrase: "additional context", Translation: "additional-translation"}}
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
