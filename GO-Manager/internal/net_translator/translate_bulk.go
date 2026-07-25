package net_translator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type translateBulkItem struct {
	Original     string   `json:"original"`
	Translations []string `json:"translations"`
}

func (n Translator) TranslateBulk(phrases []string) [][]string {
	body, err := json.Marshal(phrases)
	if err != nil {
		return nil
	}

	request, err := http.NewRequest(
		http.MethodPost,
		n.Backend+"/translate-bulk",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := n.Client.Do(request)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var response []translateBulkItem
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil
	}

	translations := make([][]string, 0, len(response))
	for _, item := range response {
		translations = append(translations, item.Translations)
	}

	return translations
}
