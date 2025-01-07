package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/taylorskalyo/goreader/epub"
	"github.com/thom151/vizz/internal/auth"
	"github.com/thom151/vizz/internal/database"
)

type PageData struct {
	PageContent template.HTML
	PrevPage    int
	NextPage    int
	IsFirstPage bool
	IsLastPage  bool
	ID          int
	Images      []string
}

func (cfg *apiConfig) handlerStory(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header, r.Cookies())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find acc token in story")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, invalid token")
		return
	}

	bookIDStr := r.URL.Query().Get("id")
	if bookIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing book id")
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	book, err := cfg.db.GetBook(r.Context(), int64(bookID))
	currPage := r.URL.Query().Get("page")
	if currPage == "" {
		currPage = "1"
	}

	threadParams := database.GetThreadParams{
		UserID: userID.String(),
		BookID: book.ID,
	}

	c := openai.NewClient(cfg.openai_key)
	thread, err := cfg.db.GetThread(r.Context(), threadParams)
	if err != nil {

		threadID, err := genThread(c, book.Title)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "cannot generate thread "+err.Error())
			return
		}
		thread, err = cfg.db.CreateThread(r.Context(), database.CreateThreadParams{
			ThreadID: threadID.ID,
			UserID:   userID.String(),
			BookID:   book.ID,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "cannot create threads "+err.Error())
			return
		}

		fmt.Println("Generated thread sir")
	}

	fmt.Printf("Currpage: %v\n", currPage)

	currPageInt, err := strconv.Atoi(currPage)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid page number")
		return
	}

	data, err := paginateEpubContent(book.EpubPath, currPageInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error paginating epub")
		return
	}
	data.ID = bookID

	tmp, err := template.ParseFiles("./static/story.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing tmp")
		return
	}
	err = tmp.Execute(w, data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error executing data")
		return
	}

	_, err = genMessage(c, thread.ThreadID, string(data.PageContent))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to gen msg")
		return
	}

	run, err := getRun(c, thread.ThreadID, cfg.assistant)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get run")
		return
	}

	messages, err := getResponse(c, thread.ThreadID, run.ID)
	data.Images = []string{}
	for i := 0; i < len(messages); i++ {
		url, err := genImageBase64(c, messages[i])
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to gen image "+err.Error())
			return
		}

		data.Images = append(data.Images, url)
	}
	err = tmp.Execute(w, data)
}

func paginateEpubContent(filePath string, currentPage int) (PageData, error) {
	epubReader, err := epub.OpenReader(filePath)
	if err != nil {
		return PageData{}, err
	}

	spineItems := epubReader.Rootfiles[0].Spine.Itemrefs
	totalPages := len(spineItems)

	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	currentItem := spineItems[currentPage-1].Item
	if currentItem == nil {
		return PageData{}, fmt.Errorf("no items")
	}

	reader, err := currentItem.Open()
	if err != nil {
		return PageData{}, err
	}
	defer reader.Close()

	var content strings.Builder
	_, err = io.Copy(&content, reader)
	if err != nil {
		return PageData{}, err
	}

	pageData := PageData{
		PageContent: template.HTML(RemoveUnwantedTags(content.String())),
		PrevPage:    currentPage - 1,
		NextPage:    currentPage + 1,
		IsFirstPage: currentPage == 1,
		IsLastPage:  currentPage == totalPages,
	}

	return pageData, nil
}

func RemoveUnwantedTags(htmlContent string) string {
	// Regular expression to match <img> tags
	imgRe := regexp.MustCompile(`<img[^>]*>`)
	// Regular expression to match <a> tags and their content
	linkRe := regexp.MustCompile(`<a[^>]*>.*?</a>`)

	// Remove <img> tags
	cleanedContent := imgRe.ReplaceAllString(htmlContent, "")
	// Remove <a> tags and their content
	cleanedContent = linkRe.ReplaceAllString(cleanedContent, "")

	return cleanedContent
}
