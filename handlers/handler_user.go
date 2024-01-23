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

func HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	apiConfig := r.Context().Value("api_config").(*db.APIConfig)

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err))
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

func HandlerGetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	helpers.RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

func HandlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	apiConfig := r.Context().Value("api_config").(*db.APIConfig)

	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting posts for user: %v", err))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.DatabasePostsToPosts(posts))

}
