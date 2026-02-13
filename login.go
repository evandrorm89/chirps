package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/evandrorm89/httpserver/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password hash: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !c {
		log.Printf("Password mismatch for user %s", params.Email)
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	respBody := response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
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
