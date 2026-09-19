package net_translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type translateBulkItem struct {
	Original     string   `json:"original"`
	Translations []string `json:"translations"`
}

func (n Translator) TranslateBulk(phrases []string) ([][]string, error) {
	if len(phrases) == 0 {
		return [][]string{}, nil
	}

	pending := make([]int, 0, len(phrases))
	toTranslate := make([]string, 0, len(phrases))
	translations := make([][]string, len(phrases))
	for i, phrase := range phrases {
		if phrase == "" {
			translations[i] = []string{}
			continue
		}
		pending = append(pending, i)
		toTranslate = append(toTranslate, phrase)
	}
	if len(toTranslate) == 0 {
		return translations, nil
	}

	body, err := json.Marshal(toTranslate)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		n.Backend+"/translate-bulk",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := n.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("translate-bulk request failed with status %s", resp.Status)
	}

	var response []translateBulkItem
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response) != len(toTranslate) {
		return nil, fmt.Errorf("translate-bulk returned %d items, expected %d", len(response), len(toTranslate))
	}

	for i, item := range response {
		translations[pending[i]] = item.Translations
	}

	return translations, nil
}
