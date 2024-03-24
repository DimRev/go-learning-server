package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	const port = "8080"
	const filepathRoot = "./public"

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	// * Router Handlers
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handlerStatic(filepathRoot)))
	mux.HandleFunc("/healthz", handlerReadiness)

	log.Printf("Serving on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileServerHits++
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerStatic(filepathRoot string) http.Handler {
	return http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
}
