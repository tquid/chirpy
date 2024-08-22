package main

import (
	"log"
	"net/http"
)

func healthz(w http.ResponseWriter, _ *http.Request) {
  w.Header().Add("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(200)
  w.Write([]byte("OK"))
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
  mux.HandleFunc("/healthz", healthz)
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
