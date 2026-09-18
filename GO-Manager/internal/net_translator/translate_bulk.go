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
	body, err := json.Marshal(phrases)
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

	translations := make([][]string, 0, len(response))
	for _, item := range response {
		translations = append(translations, item.Translations)
	}

	return translations, nil
}
