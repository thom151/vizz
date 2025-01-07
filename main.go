package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thom151/vizz/internal/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	bookCache      map[string]string
	openai_key     string
	assistant      string
}

func main() {
	production := true
	if !production {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Cannot load env" + err.Error())
		}

	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Db url not set")
	} else {
		log.Printf("Using DB_URL: %s", dbURL)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("cannot read port")
	}

	open_api := os.Getenv("OPEN_API")
	if open_api == "" {
		log.Fatal("cannot read openapi")
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal("Cannot open db" + err.Error())
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform not set")
	}

	dbQueries := database.New(db)

	if err != nil {
		log.Fatal("error parsing tmp in main")
	}

	ass := os.Getenv("ASSISTANT")
	if ass == "" {
		log.Fatal("assistant id not set")
	}
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         os.Getenv("SECRET"),
		bookCache:      make(map[string]string),
		openai_key:     open_api,
		assistant:      ass,
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

	_ = http.FileServer(http.Dir("./static/index.html"))

	mux.HandleFunc("/app/", apiCfg.handlerIndex)
	mux.HandleFunc("/api/healthz", hanlderReadiness)
	mux.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("/api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("/api/search", apiCfg.handlerSearchBooks)
	mux.HandleFunc("/api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUserUpdate)
	mux.HandleFunc("/story", apiCfg.handlerStory)
	mux.HandleFunc("/upload", apiCfg.handlerCreateBook)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func hanlderReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing ok")
		return
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Reset is only allowed in dev environment."))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error writing error ahaha")
			return
		}
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error resetting hits")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Hits reset to 0 and database reset to initial state."))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing reset")
		return
	}
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
	_, err := w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing hits")
		return
	}
}

func (cfg *apiConfig) handlerIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
