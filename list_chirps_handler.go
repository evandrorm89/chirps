package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) listChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type response []Chirp
	chirps, err := cfg.dbQueries.ListChirps(r.Context())
	if err != nil {
		log.Printf("Error listing chirps: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	re := []Chirp{}

	for _, c := range chirps {
		chirp := Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
		re = append(re, chirp)
	}

	var respBody response

	respBody = re

	resp, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
