package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/evandrorm89/httpserver/internal/auth"
)

func (cfg *apiConfig) getRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

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

	rt, err := cfg.dbQueries.GetToken(r.Context(), token)
	if err != nil {
		log.Printf("Error fetching refresh token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rt.ExpiresAt.Before(time.Now()) {
		log.Printf("Refresh token expired: %s", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if rt.RevokedAt.Valid {
		log.Printf("Refresh token revoked: %s", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newToken, err := auth.MakeJWT(rt.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error creating new JWT: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody := response{
		Token: newToken,
	}

	resp, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
