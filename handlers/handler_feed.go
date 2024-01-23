package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devldm/go-server-rss/db"
	"github.com/devldm/go-server-rss/helpers"
	"github.com/devldm/go-server-rss/internal/database"
	"github.com/devldm/go-server-rss/models"

	"github.com/google/uuid"
)

func HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	apiConfig := r.Context().Value("api_config").(*db.APIConfig)

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err))
		return
	}

	feed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		UserID:    user.ID,
		Url:       params.URL,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.DatabaseFeedToFeed(feed))
}

func HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	apiConfig := r.Context().Value("api_config").(*db.APIConfig)

	feeds, err := apiConfig.DB.GetFeeds(r.Context())
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting feeds: %v", err))
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.DatabaseFeedsToFeeds(feeds))
}
