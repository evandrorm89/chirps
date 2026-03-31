package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/evandrorm89/httpserver/internal/auth"
	"github.com/evandrorm89/httpserver/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expiresInSeconds := 3600

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Second*time.Duration(expiresInSeconds))
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
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

	refreshToken := auth.MakeRefreshToken()

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), refreshTokenParams)

	respBody := response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
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
