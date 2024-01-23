package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devldm/go-server-rss/db"
	"github.com/devldm/go-server-rss/helpers"
	"github.com/devldm/go-server-rss/models"
	"github.com/go-chi/chi/v5"

	"github.com/devldm/go-server-rss/internal/database"
	"github.com/google/uuid"
)

func HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	apiConfig := r.Context().Value("api_config").(*db.APIConfig)

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err))
		return
	}

	feedFollow, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed_follow: %v", err))
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	apiConfig := r.Context().Value("api_config").(*db.APIConfig)

	feedFollows, err := apiConfig.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting feed follows: %v", err))
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

func HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	apiConfig := r.Context().Value("api_config").(*db.APIConfig)

	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	fmt.Println(feedFollowIDStr)
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing feed follow id: %v", err))
	}

	err = apiConfig.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting feed follow: %v", err))
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil)
}
