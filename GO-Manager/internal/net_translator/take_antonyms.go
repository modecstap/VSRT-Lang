package net_translator

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func (n Translator) TakeAntonyms(word string) []string {

	request, err := http.NewRequest(
		http.MethodGet,
		n.Backend+"/antonyms?target="+url.QueryEscape(word),
		nil,
	)
	if err != nil {
		return nil
	}

	resp, err := n.Client.Do(request)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var response []string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil
	}

	return response
}
