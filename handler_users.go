package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ric-ram/go-blog-aggregator/internal/auth"
	"github.com/ric-ram/go-blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	id := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerGetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No authorization token provided")
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), headerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user from token")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
