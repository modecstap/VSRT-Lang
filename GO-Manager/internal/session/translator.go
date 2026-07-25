package session

type Translator interface {
	Translate(phrase string) (translations []string)
	TranslateBulk(phrases []string) (translations [][]string)
	TakeContexts(phrase string) (contexts []Context)
	TakeSynonyms(word string) (synonyms []string)
	TakeAntonyms(word string) (antonyms []string)
	TakeBaseForm(word string) (baseForm string)
}
