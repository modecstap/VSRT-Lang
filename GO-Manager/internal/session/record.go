package session

import "fmt"


type Context struct {
	Phrase     string `json:"phrase"`
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

type NewRecordCommand struct {
	Translator  Translator
	Phrase      string
	MainContext string
}

func NewRecord(cmd NewRecordCommand) (*Record, error) {
	if err := validateCommand(cmd); err != nil {
		return nil, err
	}

	translations, err := translateRecord(cmd)
	if err != nil {
		return nil, err
	}

	recordData, err := collectRecordData(cmd)
	if err != nil {
		return nil, err
	}

	return buildRecord(cmd, translations, recordData), nil
}

type recordTranslations struct {
	Phrase  []string
	Context string
}

type recordData struct {
	Contexts  []Context
	BaseForm  string
	Synonyms  []string
	Antonyms  []string
}

func validateCommand(cmd NewRecordCommand) error {
	if cmd.Translator == nil {
		return fmt.Errorf("translator is required")
	}

	if cmd.Phrase == "" {
		return fmt.Errorf("phrase is required")
	}

	return nil
}

func translateRecord(
	cmd NewRecordCommand,
) (recordTranslations, error) {
	phrases := []string{cmd.Phrase}
	if cmd.MainContext != "" {
		phrases = append(phrases, cmd.MainContext)
	}

	translations, err := cmd.Translator.TranslateBulk(phrases)
	if err != nil {
		return recordTranslations{}, err
	}

	if len(translations) < len(phrases) {
		return recordTranslations{}, fmt.Errorf(
			"translator returned %d translation groups for %d phrases",
			len(translations),
			len(phrases),
		)
	}

	result := recordTranslations{
		Phrase: translations[0],
	}

	if cmd.MainContext == "" {
		return result, nil
	}

	if len(translations[1]) == 0 {
		return recordTranslations{}, fmt.Errorf(
			"translator returned no translation for context %q",
			cmd.MainContext,
		)
	}

	result.Context = translations[1][0]
	return result, nil
}

func collectRecordData(cmd NewRecordCommand) (recordData, error) {
	contexts, err := cmd.Translator.TakeContexts(cmd.Phrase)
	if err != nil {
		return recordData{}, err
	}

	baseForm, err := cmd.Translator.TakeBaseForm(cmd.Phrase)
	if err != nil {
		return recordData{}, err
	}

	synonyms, err := cmd.Translator.TakeSynonyms(cmd.Phrase)
	if err != nil {
		return recordData{}, err
	}

	antonyms, err := cmd.Translator.TakeAntonyms(cmd.Phrase)
	if err != nil {
		return recordData{}, err
	}

	return recordData{
		Contexts: contexts,
		BaseForm: baseForm,
		Synonyms: synonyms,
		Antonyms: antonyms,
	}, nil
}

func buildRecord(
	cmd NewRecordCommand,
	translations recordTranslations,
	data recordData,
) *Record {
	contexts := make([]Context, 0, len(data.Contexts)+1)
	if cmd.MainContext != "" {
		contexts = append(contexts, Context{
			Phrase:      cmd.MainContext,
			Translation: translations.Context,
		})
	}
	contexts = append(contexts, data.Contexts...)

	return &Record{
		Phrase:       cmd.Phrase,
		Translations: translations.Phrase,
		Synonyms:     data.Synonyms,
		Antonyms:     data.Antonyms,
		BaseForm:     data.BaseForm,
		Contexts:     contexts,
		Count:        1,
	}
}
