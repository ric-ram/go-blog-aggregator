package main

import (
	"net/http"
	"strconv"

	"github.com/ric-ram/go-blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	defaultLimit := 10
	var postLimit int
	queryLimit := r.URL.Query().Get("limit")
	if queryLimit != "" {
		var err error
		postLimit, err = strconv.Atoi(queryLimit)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error parsing query parameter")
			return
		}
	} else {
		postLimit = defaultLimit
	}

	dbPostList, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts from database")
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostListToPostList(dbPostList))
}
