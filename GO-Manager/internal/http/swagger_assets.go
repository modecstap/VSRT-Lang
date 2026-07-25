package router

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed docs/*
var embeddedDocs embed.FS

// RegisterSwagger registers handlers that serve embedded OpenAPI docs and Swagger UI.
func RegisterSwagger(mux *http.ServeMux) {
	// serve openapi.yaml
	docBytes, err := embeddedDocs.ReadFile("docs/openapi.yaml")
	if err == nil {
		mux.HandleFunc("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-yaml")
			_, _ = w.Write(docBytes)
		})
	} else {
		mux.HandleFunc("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}

	// serve swagger-ui static files from embedded FS
	sub, err := fs.Sub(embeddedDocs, "docs/swagger-ui")
	if err == nil {
		mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(sub))))
	} else {
		mux.Handle("/swagger/", http.NotFoundHandler())
	}
}
