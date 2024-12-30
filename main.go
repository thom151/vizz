package main

import (
	"log"
	"net/http"
)

func main() {

	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	handler := http.FileServer(http.Dir("./static/"))

	mux.Handle("/app/", http.StripPrefix("/app", handler))
	mux.HandleFunc("/healthz", hanlderReadiness)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func hanlderReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
