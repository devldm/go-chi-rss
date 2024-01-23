package middleware

import (
	"fmt"
	"net/http"

	"github.com/devldm/go-server-rss/db"
	"github.com/devldm/go-server-rss/helpers"
	"github.com/devldm/go-server-rss/internal/database"
	"github.com/devldm/go-server-rss/internal/database/auth"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiConfig := r.Context().Value("api_config").(*db.APIConfig)

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			helpers.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		fetchedUser, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, fetchedUser)
	}
}
