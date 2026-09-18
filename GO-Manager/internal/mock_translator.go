package internal

import "VSRT-Lang/internal/session"

type MockTranslator struct {
}

func (m MockTranslator) TakeBaseForm(word string) (string, error) {
	return "base", nil
}

func (m MockTranslator) Translate(phrase string) ([]string, error) {
	result := []string{"тест", "испытание"}
	return result, nil
}

func (m MockTranslator) TranslateBulk(phrases []string) ([][]string, error) {
	result := [][]string{
		{"тест", "испытание"},
		{"второй", "следующий"},
	}
	return result, nil
}

func (m MockTranslator) TakeContexts(phrase string) ([]session.Context, error) {
	result := []session.Context{
		{
			Phrase:      "Test Context",
			Translation: "тестовый контекст",
		},
	}
	return result, nil
}

func (m MockTranslator) TakeSynonyms(word string) ([]string, error) {
	result := []string{"синоним", "синоним2"}
	return result, nil
}

func (m MockTranslator) TakeAntonyms(word string) ([]string, error) {
	result := []string{"антоним", "антоним2"}
	return result, nil
}
