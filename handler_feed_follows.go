package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ric-ram/go-blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	id := uuid.New()
	createdAt := time.Now().UTC()
	updatedAt := createdAt

	feedFollow, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) handlerGetFeedFollowsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowsList, err := cfg.DB.GetAllFeedFollowsFromUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting feed follows for user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowListToFeedFollowList(feedFollowsList))
}

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	paramFeedFollowsID := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(paramFeedFollowsID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Feed Follow ID")
		return
	}

	err = cfg.DB.DeleteFeedFollows(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting Feed Follow")
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
