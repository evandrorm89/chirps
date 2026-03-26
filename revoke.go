package main

import (
	"log"
	"net/http"

	"github.com/evandrorm89/httpserver/internal/auth"
)

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token == "" {
		log.Printf("No token provided in Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.dbQueries.UpdateToken(r.Context(), token)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
