package net_translator

import "net/http"

type Translator struct {
	Client  http.Client
	Backend string
}
