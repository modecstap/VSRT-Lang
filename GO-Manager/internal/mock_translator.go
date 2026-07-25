package internal

import "VSRT-Lang/internal/session"

type MockTranslator struct {
}

func (m MockTranslator) TakeBaseForm(word string) (baseForm string) {
	return "base"
}

func (m MockTranslator) Translate(phrase string) []string {
	result := []string{"тест", "испытание"}
	return result
}

func (m MockTranslator) TranslateBulk(phrases []string) [][]string {
	result := [][]string{
		{"тест", "испытание"},
		{"второй", "следующий"},
	}
	return result
}

func (m MockTranslator) TakeContexts(phrase string) []session.Context {
	result := []session.Context{
		{
			Phrase:      "Test Context",
			Translation: "тестовый контекст",
		},
	}
	return result
}

func (m MockTranslator) TakeSynonyms(word string) []string {
	result := []string{"синоним", "синоним2"}
	return result
}

func (m MockTranslator) TakeAntonyms(word string) []string {
	result := []string{"антоним", "антоним2"}
	return result
}
