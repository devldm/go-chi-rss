package middleware

import (
	"context"
	"net/http"

	"github.com/devldm/go-server-rss/db"
)

// func ConfigMiddleware(next http.Handler, config *db.APIConfig) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.WithValue(r.Context(), "api_config", config)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func ConfigMiddleware(config *db.APIConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "api_config", config)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
