package handlers

import (
	"net/http"

	"github.com/devldm/go-server-rss/helpers"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}
