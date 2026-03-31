package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) webhookHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if params.Event != "user.upgraded" {
		log.Printf("Received unsupported event type: %s", params.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = cfg.dbQueries.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		log.Printf("Error upgrading user to Chiry Red")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
