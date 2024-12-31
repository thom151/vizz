package main

import (
	"net/http"
	"time"

	"github.com/thom151/vizz/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bearer missing")
		return
	}

	token, err := cfg.db.GetToken(r.Context(), refToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Token doesn't exist")
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusBadRequest, "Token expired")
		return
	}

	newJWT, err := auth.MakeJWT(token.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Creating jwt failed")
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusBadRequest, "Token has been revoked")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newJWT,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	token, err := cfg.db.GetToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Token does not exist")
		return
	}

	err = cfg.db.RevokeToken(r.Context(), token.Token)
	if err != nil {
		respondWithError(w, 500, "Token not revoked")
	}

	respondWithJSON(w, 204, nil)
}
