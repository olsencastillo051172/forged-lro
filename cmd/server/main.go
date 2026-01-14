package main

import (
	"log"
	"net/http"
	"os"
)

const version = "dev" // luego lo amarramos a tags/ldflags

func main() {
	// Puerto con fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Router mínimo (sin frameworks)
	mux := http.NewServeMux()

	// RegisterRoutes (mínimo)
	registerRoutes(mux)

	addr := "0.0.0.0:" + port
	log.Printf("FORGED-LRO server starting on %s", addr)

	// IMPORTANTE: NO goroutine, y el error no se ignora
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(version))
	})
}

