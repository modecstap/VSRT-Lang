package stub

import "VSRT-Lang/internal/session"

type Translator struct{}

func (Translator) Translate(phrase string) ([]string, error) {
	return []string{"translated:" + phrase}, nil
}

func (Translator) TranslateBulk(phrases []string) ([][]string, error) {
	return [][]string{
		{"phrase-translation-1", "phrase-translation-2"},
		{"context-translation"},
	}, nil
}

func (Translator) TakeContexts(phrase string) ([]session.Context, error) {
	return []session.Context{{Phrase: "additional context", Translation: "additional-translation"}}, nil
}

func (Translator) TakeSynonyms(word string) ([]string, error) {
	return []string{"synonym1", "synonym2"}, nil
}

func (Translator) TakeAntonyms(word string) ([]string, error) {
	return []string{"antonym1", "antonym2"}, nil
}

func (Translator) TakeBaseForm(word string) (string, error) {
	return "base-form:" + word, nil
}
