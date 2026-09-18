package session

import "fmt"

type Context struct {
	Phrase      string `json:"phrase"`
	Translation string `json:"translation"`
}

type Record struct {
	Phrase       string    `json:"phrase"`
	Translations []string  `json:"translations"`
	Synonyms     []string  `json:"synonyms"`
	Antonyms     []string  `json:"antonyms"`
	BaseForm     string    `json:"base_form"`
	Contexts     []Context `json:"contexts"`
	Count        int64     `json:"count"`
}

func NewRecord(
	translator Translator,
	phrase string,
	mainContext string,
) (*Record, error) {
	toTranslate := []string{phrase, mainContext}
	translations, err := translator.TranslateBulk(toTranslate)
	if err != nil {
		return nil, err
	}
	if len(translations) < 2 {
		return nil, fmt.Errorf("translator returned %d translation groups for %d phrases", len(translations), len(toTranslate))
	}

	phraseTranslations := translations[0]
	contextTranslations := translations[1]
	if len(contextTranslations) == 0 {
		return nil, fmt.Errorf("translator returned no translation for context %q", mainContext)
	}

	contexts, err := translator.TakeContexts(phrase)
	if err != nil {
		return nil, err
	}

	baseForm, err := translator.TakeBaseForm(phrase)
	if err != nil {
		return nil, err
	}

	synonyms, err := translator.TakeSynonyms(phrase)
	if err != nil {
		return nil, err
	}

	antonyms, err := translator.TakeAntonyms(phrase)
	if err != nil {
		return nil, err
	}

	result := &Record{
		Phrase:       phrase,
		Translations: phraseTranslations,
		BaseForm:     baseForm,
		Synonyms:     synonyms,
		Antonyms:     antonyms,
		Count:        1,
	}

	result.Contexts = append(result.Contexts, Context{Phrase: mainContext, Translation: contextTranslations[0]})
	result.Contexts = append(result.Contexts, contexts...)

	return result, nil
}
