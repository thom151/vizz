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
		http.Redirect(w, r, "api/login", http.StatusFound)
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
		currPage = "6"
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
		respondWithError(w, http.StatusInternalServerError, "error paginating epub "+err.Error())
		return
	}
	data.ID = bookID
	fmt.Println(thread.ThreadID)
	err = GenerateMessagesAndImages(c, thread.ThreadID, cfg.assistant, &data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating messages and images "+err.Error())
		return
	}

	fmt.Println("EXECUTING TEMPLATE AGAIN")
	tmp, err := template.ParseFiles("./static/story.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing tmp")
		return
	}
	err = tmp.Execute(w, data)

}

func GenerateMessagesAndImages(c *openai.Client, threadID, assistantID string, data *PageData) error {

	messagesChan := make(chan string)
	imagesChan := make(chan string)
	errChan := make(chan error, 2)

	go generateMessages(messagesChan, errChan, c, threadID, assistantID, data)

	go func() {
		for prompt := range messagesChan {
			url, err := genImageBase64(c, prompt)
			if err != nil {
				errChan <- err
				return
			}
			imagesChan <- url
		}

		close(imagesChan)
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case image, ok := <-imagesChan:
			if !ok {
				return nil
			}

			fmt.Println("\n\n url: ", image, "\n\n")
			data.Images = append(data.Images, image)
		}
	}

}

func generateMessages(messagesChan chan<- string, errChan chan<- error, c *openai.Client, threadID, assistantID string, data *PageData) {
	defer close(messagesChan)
	_, err := genMessage(c, threadID, stripHTMLTags(string(data.PageContent)))
	if err != nil {
		errChan <- err
		return
	}

	run, err := getRun(c, threadID, assistantID)
	if err != nil {
		errChan <- err
	}

	messages, err := getResponse(c, threadID, run.ID)
	if err != nil {
		errChan <- err
	}

	for _, message := range messages {
		messagesChan <- message
	}

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

	fmt.Println(content.String())

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
	imgRe := regexp.MustCompile(`(?i)<img[^>]*>`)
	// Regular expression to match <a> tags and their content
	linkRe := regexp.MustCompile(`(?i)<a[^>]*?>|</a>`)

	// Remove <img> tags
	cleanedContent := imgRe.ReplaceAllString(htmlContent, "")
	// Remove <a> tags and their content
	cleanedContent = linkRe.ReplaceAllString(cleanedContent, "")

	return cleanedContent
}

func stripHTMLTags(input string) string {
	re := regexp.MustCompile(`(?i)<[^>]*>`)
	return re.ReplaceAllString(input, "")
}
