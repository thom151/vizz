package main

type BookResults struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		ID      int        `json:"id"`
		Title   string     `json:"title"`
		Authors []struct { /*...*/
		} `json:"authors"`
		Translators   []any             `json:"translators"`
		Subjects      []string          `json:"subjects"`
		Bookshelves   []string          `json:"bookshelves"`
		Languages     []string          `json:"languages"`
		Copyright     bool              `json:"copyright"`
		MediaType     string            `json:"media_type"`
		Formats       map[string]string `json:"formats,omitempty"`
		DownloadCount int               `json:"download_count"`
	} `json:"results"`
}
