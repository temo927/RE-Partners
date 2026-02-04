package http

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests: browsers send these before actual requests
		// to check CORS permissions. We respond immediately without calling next handler.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimit applies rate limiting middleware: 100 requests per minute per IP.
// Uses token bucket algorithm to prevent abuse while allowing legitimate traffic.
func RateLimit(next http.Handler) http.Handler {
	return httprate.Limit(
		100,
		1*time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP),
	)(next)
}
