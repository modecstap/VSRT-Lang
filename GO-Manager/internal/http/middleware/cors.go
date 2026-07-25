package middleware

import "net/http"

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	AllowCredentials bool
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	origins := make(map[string]struct{})
	for _, origin := range cfg.AllowedOrigins {
		origins[origin] = struct{}{}
	}

	methods := join(cfg.AllowedMethods)
	headers := join(cfg.AllowedHeaders)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" {
				if _, ok := origins["*"]; ok {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if _, ok := origins[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
				}

				if methods != "" {
					w.Header().Set("Access-Control-Allow-Methods", methods)
				}

				if headers != "" {
					w.Header().Set("Access-Control-Allow-Headers", headers)
				}

				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			// Preflight request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func join(values []string) string {
	if len(values) == 0 {
		return ""
	}

	result := values[0]
	for _, v := range values[1:] {
		result += ", " + v
	}

	return result
}