package net_translator

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func (n Translator) TakeBaseForm(word string) string {
	request, err := http.NewRequest(
		http.MethodGet,
		n.Backend+"/base-form?target="+url.QueryEscape(word),
		nil,
	)
	if err != nil {
		return ""
	}

	resp, err := n.Client.Do(request)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	var response string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ""
	}

	return response

}
