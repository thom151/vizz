package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/taylorskalyo/goreader/epub"
	"github.com/thom151/vizz/internal/auth"
	"github.com/thom151/vizz/internal/database"
)

func (cfg *apiConfig) handlerCreateBook(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header, r.Cookies())
	if err != nil {
		return
	}

	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		http.Redirect(w, r, "/api/login", http.StatusFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "./static/upload.html")
		return
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Could not parse form")
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "err retrieving file "+err.Error())
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error reading file")
			return
		}

		tempFileName := fmt.Sprintf("./temp/%s", handler.Filename)
		err = os.WriteFile(tempFileName, fileBytes, 0644)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error saving temp file")
			return
		}

		bookEntryParams, err := createBookParams(tempFileName)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		bookEntryParams.EpubPath = tempFileName

		book, err := cfg.db.CreateBookEntry(r.Context(), bookEntryParams)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, "err creating book")
			return
		}

		respondWithJSON(w, http.StatusOK, book)

	}
}

func createBookParams(file string) (database.CreateBookEntryParams, error) {
	epubReader, err := epub.OpenReader(file)
	if err != nil {
		return database.CreateBookEntryParams{}, err
	}
	defer epubReader.Close()

	book := epubReader.Rootfiles[0]

	toNullString := func(s string) sql.NullString {
		if s == "" {
			return sql.NullString{Valid: false}
		}
		return sql.NullString{String: s, Valid: true}
	}

	bookEntryParams := database.CreateBookEntryParams{
		Title:       book.Title,
		Author:      toNullString(book.Creator),
		Description: toNullString(book.Description),
		EpubPath:    "",
	}

	return bookEntryParams, nil
}
