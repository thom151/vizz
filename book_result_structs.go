package main

type BookResults struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Authors []struct {
			Name      string `json:"name"`
			BirthYear int    `json:"birth_year"`
			DeathYear int    `json:"death_year"`
		} `json:"authors"`
		Translators []any    `json:"translators"`
		Subjects    []string `json:"subjects"`
		Bookshelves []string `json:"bookshelves"`
		Languages   []string `json:"languages"`
		Copyright   bool     `json:"copyright"`
		MediaType   string   `json:"media_type"`
		Formats     struct {
			TextHTML                    string `json:"text/html"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		DownloadCount int `json:"download_count"`
		Formats0      struct {
			TextHTML                    string `json:"text/html"`
			TextHTMLCharsetUtf8         string `json:"text/html; charset=utf-8"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetIso88591    string `json:"text/plain; charset=iso-8859-1"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		Formats1 struct {
			TextHTML                    string `json:"text/html"`
			TextHTMLCharsetIso88591     string `json:"text/html; charset=iso-8859-1"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			TextPlainCharsetIso88591    string `json:"text/plain; charset=iso-8859-1"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		Formats2 struct {
			TextHTML                    string `json:"text/html"`
			TextHTMLCharsetUtf8         string `json:"text/html; charset=utf-8"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetUtf8        string `json:"text/plain; charset=utf-8"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		Formats3 struct {
			ImageJpeg         string `json:"image/jpeg"`
			TextPlain         string `json:"text/plain"`
			TextHTML          string `json:"text/html"`
			AudioOgg          string `json:"audio/ogg"`
			AudioMp4          string `json:"audio/mp4"`
			AudioMpeg         string `json:"audio/mpeg"`
			ApplicationRdfXML string `json:"application/rdf+xml"`
		} `json:"formats,omitempty"`
		Formats4 struct {
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			TextHTML                    string `json:"text/html"`
			TextHTMLCharsetUsASCII      string `json:"text/html; charset=us-ascii"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetIso88591    string `json:"text/plain; charset=iso-8859-1"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		Formats5 struct {
			TextHTML                    string `json:"text/html"`
			TextHTMLCharsetUtf8         string `json:"text/html; charset=utf-8"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			TextPlainCharsetUtf8        string `json:"text/plain; charset=utf-8"`
			TextPlainCharsetIso88591    string `json:"text/plain; charset=iso-8859-1"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		Formats6 struct {
			TextHTML                    string `json:"text/html"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetIso88591    string `json:"text/plain; charset=iso-8859-1"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
		} `json:"formats,omitempty"`
		Formats7 struct {
			TextHTML                    string `json:"text/html"`
			TextHTMLCharsetUtf8         string `json:"text/html; charset=utf-8"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			TextPlainCharsetUtf8        string `json:"text/plain; charset=utf-8"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		Formats8 struct {
			TextHTML                    string `json:"text/html"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			TextPlainCharsetIso88591    string `json:"text/plain; charset=iso-8859-1"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		Formats9 struct {
			TextHTML                    string `json:"text/html"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetIso88591    string `json:"text/plain; charset=iso-8859-1"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
		} `json:"formats,omitempty"`
	} `json:"results"`
}
