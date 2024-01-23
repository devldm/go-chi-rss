package handlers

import (
	"net/http"

	"github.com/devldm/go-server-rss/helpers"
)

func HandlerError(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}
