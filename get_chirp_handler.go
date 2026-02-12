package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}

	id := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error parsing chirp ID: %v", err)
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp: %v", err)
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	respBody := response{
		Chirp: Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		},
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
