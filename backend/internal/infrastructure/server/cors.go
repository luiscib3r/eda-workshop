package server

import (
	"backend/internal/infrastructure/config"
	"net/http"
)

func cors(cfg config.CorsConfig) func(http.Handler) http.Handler {
	origin := "*"
	if len(cfg.AllowedOrigins) > 0 {
		origin = cfg.AllowedOrigins[0]
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == "OPTIONS" {
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
