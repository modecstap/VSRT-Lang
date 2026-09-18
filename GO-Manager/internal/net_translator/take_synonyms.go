package net_translator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (n Translator) TakeSynonyms(word string) ([]string, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		n.Backend+"/synonyms?target="+url.QueryEscape(word),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := n.Client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("synonyms request failed with status %s", resp.Status)
	}

	var response []string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
