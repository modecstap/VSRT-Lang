package net_translator

import (
	"VSRT-Lang/internal/session"
	"bytes"
	"encoding/json"
	"net/http"
)

type contextRequest struct {
	Content string `json:"content"`
}

type contextResponse struct {
	Original     string   `json:"original"`
	Translations []string `json:"translations"`
}

func (n Translator) TakeContexts(phrase string) []session.Context {
	body, err := json.Marshal(contextRequest{
		Content: phrase,
	})
	if err != nil {
		return nil
	}

	request, err := http.NewRequest(
		http.MethodPost,
		n.Backend+"/context",
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

	var response []contextResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil
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

	return contexts
}
