package session

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrPhraseRequired = errors.New("phrase is required")

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
	Contexts []Context
	BaseForm string
	Synonyms []string
	Antonyms []string
}

func validateCommand(cmd NewRecordCommand) error {
	if cmd.Translator == nil {
		return fmt.Errorf("translator is required")
	}

	if strings.TrimSpace(cmd.Phrase) == "" {
		return ErrPhraseRequired
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
	var (
		wg sync.WaitGroup

		contexts    []Context
		contextsErr error
		baseForm    string
		baseFormErr error
		synonyms    []string
		synonymsErr error
		antonyms    []string
		antonymsErr error
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		contexts, contextsErr = cmd.Translator.TakeContexts(cmd.Phrase)
	}()
	go func() {
		defer wg.Done()
		baseForm, baseFormErr = cmd.Translator.TakeBaseForm(cmd.Phrase)
	}()
	go func() {
		defer wg.Done()
		synonyms, synonymsErr = cmd.Translator.TakeSynonyms(cmd.Phrase)
	}()
	go func() {
		defer wg.Done()
		antonyms, antonymsErr = cmd.Translator.TakeAntonyms(cmd.Phrase)
	}()
	wg.Wait()

	switch {
	case contextsErr != nil:
		return recordData{}, contextsErr
	case baseFormErr != nil:
		return recordData{}, baseFormErr
	case synonymsErr != nil:
		return recordData{}, synonymsErr
	case antonymsErr != nil:
		return recordData{}, antonymsErr
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
		Phrase:       strings.TrimSpace(cmd.Phrase),
		Translations: translations.Phrase,
		Synonyms:     data.Synonyms,
		Antonyms:     data.Antonyms,
		BaseForm:     data.BaseForm,
		Contexts:     contexts,
		Count:        1,
	}
}

func (r *Record) registerSave(context string, translator Translator) error {
	if err := r.ensureUserContext(context, translator); err != nil {
		return err
	}
	r.Count++
	return nil
}

func (r *Record) ensureUserContext(context string, translator Translator) error {
	if context == "" {
		return nil
	}
	for _, existing := range r.Contexts {
		if existing.Phrase == context {
			return nil
		}
	}

	translation, err := translateUserContext(translator, context)
	if err != nil {
		return err
	}

	r.Contexts = append(r.Contexts, Context{
		Phrase:      context,
		Translation: translation,
	})
	return nil
}

func translateUserContext(translator Translator, context string) (string, error) {
	translations, err := translator.Translate(context)
	if err != nil {
		return "", err
	}
	if len(translations) == 0 {
		return "", fmt.Errorf("translator returned no translation for context %q", context)
	}
	return translations[0], nil
}
