package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/vizz/internal/auth"
)

func (cfg *apiConfig) handlerSearchBooks(w http.ResponseWriter, r *http.Request) {
	type response struct {
		BookResults BookResults
		UserID      uuid.UUID
	}

	token, err := auth.GetBearerToken(r.Header, r.Cookies())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find acc token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid jwt")
		return
	}

	book := r.URL.Query().Get("book")

	fmt.Printf("Book: %s\n", string(book))

	url := "https://gutendex.com/books?search=" + book

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Creating req failed")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Client Request do failed")
		return
	}

	defer resp.Body.Close()

	var bookResults BookResults
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&bookResults)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding  results failed")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		BookResults: bookResults,
		UserID:      userID,
	})
}
