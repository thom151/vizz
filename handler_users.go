package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/vizz/internal/auth"
	"github.com/thom151/vizz/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type UserIncoming struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userInc := UserIncoming{}
	err := decoder.Decode(&userInc)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode email")
		return
	}

	hashed, err := auth.HashPassword(userInc.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing pass")
		return
	}
	userID := uuid.New()

	fmt.Printf("%s, %v, %v", userInc.Email, hashed, userID)
	userParams := database.CreateUserParams{
		Email:          userInc.Email,
		HashedPassword: hashed,
		ID:             userID,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user"+err.Error())
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type UserLoginInc struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ExpirationInSec int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	var userLoginInc UserLoginInc
	err := decoder.Decode(&userLoginInc)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Login failed decoding")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), userLoginInc.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Email doesn't exist")
		return
	}

	err = auth.CheckPasswordHash(userLoginInc.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email/pass")
		return
	}

	expTime := time.Hour
	accToken, err := auth.MakeJWT(user.ID, cfg.secret, expTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed creating acc token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed creating ref token")
		return
	}

	refreshParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating ref token (db)")
		return
	}

	myUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, response{
		User:         myUser,
		Token:        accToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No bearer")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid jwt")
		return
	}

	decoder := json.NewDecoder(r.Body)
	uParams := params{}
	err = decoder.Decode(&uParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding params")
		return
	}

	hashed, err := auth.HashPassword(uParams.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	updatedParams := database.UpdateUserParams{
		ID:             userID,
		Email:          uParams.Email,
		HashedPassword: hashed,
	}

	user, err := cfg.db.UpdateUser(r.Context(), updatedParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
