package session

type Translator interface {
	Translate(phrase string) ([]string, error)
	TranslateBulk(phrases []string) ([][]string, error)
	TakeContexts(phrase string) ([]Context, error)
	TakeSynonyms(word string) ([]string, error)
	TakeAntonyms(word string) ([]string, error)
	TakeBaseForm(word string) (string, error)
}
