package net_translator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type translateRequest struct {
	Phrase string `json:"Phrase"`
}

type translateResponse struct {
	Phrase       string   `json:"original"`
	Translations []string `json:"translations"`
}

func (n Translator) Translate(phrase string) []string {
	body, err := json.Marshal(translateRequest{Phrase: phrase})
	if err != nil {
		return nil
	}

	request, err := http.NewRequest(
		http.MethodPost,
		n.Backend+"/translate",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil
	}

	resp, err := n.Client.Do(request)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var response translateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil
	}

	return response.Translations
}
