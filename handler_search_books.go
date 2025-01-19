package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/vizz/internal/auth"
	"github.com/thom151/vizz/internal/database"
)

type BookData struct {
	Books []database.Book
}

func (cfg *apiConfig) handlerSearchBooks(w http.ResponseWriter, r *http.Request) {

	type response struct {
		BookFile string
		UserID   uuid.UUID
	}

	token, err := auth.GetBearerToken(r.Header, r.Cookies())
	fmt.Println("Token: ", err)
	log.Printf("Got cookies: %v\n", token)
	if err != nil {
		http.Redirect(w, r, "/api/login", http.StatusFound)
		return
	}

	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		fmt.Println("Not logged in")
		http.Redirect(w, r, "/api/login", http.StatusFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "./static/search.html")
		return

	case http.MethodPost:

		bookQuery := r.FormValue("book")

		var bookParam sql.NullString
		if bookQuery != "" {
			bookParam = sql.NullString{String: bookQuery, Valid: true}
		} else {
			bookParam = sql.NullString{Valid: false}
		}
		books, err := cfg.db.GetBooks(r.Context(), bookParam)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Database query failed")
			return
		}

		tmp, err := template.ParseFiles("./static/results.html")
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error parsing template")
			return
		}

		bookData := BookData{
			Books: books,
		}

		fmt.Printf("Bookdata: %v\n", bookQuery)
		err = tmp.Execute(w, bookData)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error executing data")
			return
		}
	}
}

func apiSearch(book string) (BookResults, error) {
	url := "https://gutendex.com/books?search=" + book

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BookResults{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return BookResults{}, err
	}

	defer resp.Body.Close()

	var bookResults BookResults
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&bookResults)
	if err != nil {
		return BookResults{}, err
	}

	return bookResults, nil
}

func getFileFromBookResults(results BookResults) string {
	file := results.Results[0].Formats["text/plain; charset=us-ascii"]
	return file
}
