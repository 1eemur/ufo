package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"ufo/internal/config"
	"ufo/internal/handlers"
	"ufo/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Ensure data directory exists
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Setup router
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware)

	// Create upload handler
	uploadHandler := handlers.NewUploadHandler(cfg)

	// Routes
	r.HandleFunc("/", uploadHandler.HomePage).Methods("GET")
	r.HandleFunc("/upload", uploadHandler.UploadFile).Methods("POST")
	r.HandleFunc("/files", uploadHandler.ListFiles).Methods("GET")

	// Static files
	staticDir := filepath.Join("web", "static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	log.Printf("UFO File Upload Server starting on port %s", cfg.Port)
	log.Printf("Upload directory: %s", cfg.DataDir)
	log.Printf("Max file size: %d MB", cfg.MaxFileSize/(1024*1024))

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
