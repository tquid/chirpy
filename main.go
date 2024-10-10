package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/tquid/chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	platform       string
	userStore      UserStore
	chirpStore     ChirpStore
	jwtSecret      string
}

const jwtExpireSeconds = 3600

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("unable to open connection to %s: %v", dbURL, err)
	}
	dbQueries := database.New(db)
	pgUserStore := &PgUserStore{
		db: dbQueries,
	}
	pgChirpStore := &PgChirpStore{
		db: dbQueries,
	}
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		platform:       platform,
		userStore:      pgUserStore,
		chirpStore:     pgChirpStore,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpById)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateEmailPassword)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
