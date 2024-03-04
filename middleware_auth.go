package main

import (
	"net/http"

	"github.com/ric-ram/go-blog-aggregator/internal/auth"
	"github.com/ric-ram/go-blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, user)
	})
}
