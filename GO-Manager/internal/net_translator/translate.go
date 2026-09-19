package net_translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type translateRequest struct {
	Phrase string `json:"Phrase"`
}

type translateResponse struct {
	Phrase       string   `json:"original"`
	Translations []string `json:"translations"`
}

func (n Translator) Translate(phrase string) ([]string, error) {
	if phrase == "" {
		return []string{}, nil
	}

	body, err := json.Marshal(translateRequest{Phrase: phrase})
	if err != nil {
		return nil, err
	}

	endpoint := n.Backend + "/translate?phrase=" + url.QueryEscape(phrase)
	request, err := http.NewRequest(
		http.MethodPost,
		endpoint,
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
		return nil, fmt.Errorf("translate request failed with status %s", resp.Status)
	}

	var response translateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Translations, nil
}
