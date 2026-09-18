package net_translator

import (
	"VSRT-Lang/internal/session"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type contextResponse struct {
	Original     string   `json:"original"`
	Translations []string `json:"translations"`
}

func (n Translator) TakeContexts(phrase string) ([]session.Context, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		n.Backend+"/context?target="+url.QueryEscape(phrase),
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
		return nil, fmt.Errorf("context request failed with status %s", resp.Status)
	}

	var response []contextResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contexts := make([]session.Context, 0, len(response))
	for _, item := range response {
		translation := ""
		if len(item.Translations) > 0 {
			translation = item.Translations[0]
		}

		contexts = append(contexts, session.Context{
			Phrase:      item.Original,
			Translation: translation,
		})
	}

	return contexts, nil
}
