package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type UserIncoming struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	email := UserIncoming{}
	err := decoder.Decode(&email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode email")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	myUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, myUser)
}
