package main

import (
	"log"
	"net/http"

	"github.com/evandrorm89/httpserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing chirp ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if chirp.UserID != userID {
		log.Printf("User %v is not authorized to delete chirp %v", userID, chirpID)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
