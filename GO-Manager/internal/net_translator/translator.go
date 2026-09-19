package net_translator

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Translator struct {
	Client  http.Client
	Backend string
}

func NewTranslator(host string) *Translator {
	backend := host
	if !strings.HasPrefix(backend, "http://") && !strings.HasPrefix(backend, "https://") {
		backend = fmt.Sprintf("http://%s", backend)
	}
	backend = strings.TrimRight(backend, "/")
	if !strings.HasSuffix(backend, "/api/translator") {
		backend += "/api/translator"
	}

	return &Translator{
		Client: http.Client{
			Timeout: 30 * time.Second,
		},
		Backend: backend,
	}
}
