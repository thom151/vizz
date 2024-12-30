package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thom151/vizz/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {

	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Db url not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot open db")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform not set")
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	handler := http.FileServer(http.Dir("./static/"))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", handler)))
	mux.HandleFunc("GET /api/healthz", hanlderReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func hanlderReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileserverHits.Store(0)
	cfg.db.Reset(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}
