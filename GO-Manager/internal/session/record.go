package session

type Context struct {
	Phrase      string
	Translation string
}

type Record struct {
	Phrase       string
	Translations []string
	Synonyms     []string
	Antonyms     []string
	BaseForm     string
	Contexts     []Context
}

func NewRecord(
	translator Translator,
	phrase string,
	mainContext string,
) *Record {

	toTranslate := []string{phrase, mainContext}
	translations := translator.TranslateBulk(toTranslate)
	phraseTranslations := translations[0]
	contextTranslations := translations[1]

	var contexts []Context
	contexts = append(contexts, Context{mainContext, contextTranslations[0]})
	contexts = append(contexts, translator.TakeContexts(phrase)[:]...)

	return &Record{
		Phrase:       phrase,
		Translations: phraseTranslations,
		BaseForm:     translator.TakeBaseForm(phrase),
		Synonyms:     translator.TakeSynonyms(phrase),
		Antonyms:     translator.TakeAntonyms(phrase),
		Contexts:     contexts,
	}
}
