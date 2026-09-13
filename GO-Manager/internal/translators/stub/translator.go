package stub

import "VSRT-Lang/internal/session"

type Translator struct{}

func (Translator) Translate(phrase string) []string {
	return []string{"translated:" + phrase}
}

func (Translator) TranslateBulk(phrases []string) [][]string {
	return [][]string{
		{"phrase-translation-1", "phrase-translation-2"},
		{"context-translation"},
	}
}

func (Translator) TakeContexts(phrase string) []session.Context {
	return []session.Context{{Phrase: "additional context", Translation: "additional-translation"}}
}

func (Translator) TakeSynonyms(word string) []string {
	return []string{"synonym1", "synonym2"}
}

func (Translator) TakeAntonyms(word string) []string {
	return []string{"antonym1", "antonym2"}
}

func (Translator) TakeBaseForm(word string) string {
	return "base-form:" + word
}
