package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"visioncloud/config"
	"visioncloud/handlers"
	"visioncloud/services"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize services
	s3Handler, err := services.NewS3Handler(cfg.S3Bucket)
	if err != nil {
		log.Fatalf("Failed to initialize S3: %v", err)
	}

	upscalerService := services.NewUpscalerService("python", cfg.ModelPath)

	// Initialize handlers
	handlers.InitHandlers(upscalerService, s3Handler)

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/upscale", handlers.UpscaleImage).Methods("POST")

	// Start server
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down server...")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
