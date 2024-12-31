package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	decodeError := ErrorResponse{
		Error: msg,
	}
	errDat, _ := json.Marshal(decodeError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(errDat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	w.Write([]byte(resp))

}
