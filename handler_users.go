package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/vizz/internal/auth"
	"github.com/thom151/vizz/internal/database"
)

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "/static/signup.html")
		return

	case http.MethodPost:

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

		log.Printf("Decoded incoming user data: %v\n", userInc)

		fmt.Printf("%s, %v, %v", userInc.Email, hashed, userID)
		userParams := database.CreateUserParams{
			Email:          userInc.Email,
			HashedPassword: hashed,
			ID:             uuid.New().String(),
			UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
			CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		}

		user, err := cfg.db.CreateUser(r.Context(), userParams)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create user"+err.Error())
			return
		}

		_ = User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
		log.Println("User created successfully, redirecting to /app")
		http.Redirect(w, r, "/app", http.StatusFound)

		return
	}
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "/static/login.html")
		return

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Could not parse form data", http.StatusBadRequest)
			return
		}

		type UserLoginInc struct {
			Email           string `json:"email"`
			Password        string `json:"password"`
			ExpirationInSec int    `json:"expires_in_seconds"`
		}

		type response struct {
			User
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}

		decoder := json.NewDecoder(r.Body)
		var userLoginInc UserLoginInc
		err = decoder.Decode(&userLoginInc)
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

		userUUID, err := uuid.Parse(user.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error parsing id in login")
			return
		}
		expTime := time.Hour
		accToken, err := auth.MakeJWT(userUUID, cfg.secret, expTime)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed creating acc token")
			return
		}
		fmt.Println("Created user: %v\n", accToken)

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed creating ref token")
			return
		}

		refreshParams := database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour).Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
			RevokedAt: sql.NullString{},
		}

		_, err = cfg.db.CreateRefreshToken(r.Context(), refreshParams)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating ref token (db)")
			return
		}

		_ = User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		//fmt.Println("Created user: %v\n", myUser)

		http.Redirect(w, r, "/app", http.StatusSeeOther)
		return

		//respondWithJSON(w, http.StatusOK, response{
		//	User:         myUser,
		//	Token:        accToken,
		//	RefreshToken: refreshToken,
		//})

	}
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
		ID:             userID.String(),
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
