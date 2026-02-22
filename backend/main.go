//go:build !ec2

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gorilla/mux"

	appConfig "visioncloud/config"
	"visioncloud/handlers"
	"visioncloud/services"
)

func main() {
	// Load configuration
	cfg := appConfig.LoadConfig()
	log.Printf("Starting VisionCloud server on port %s", cfg.Port)

	// Initialize AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.AWSRegion))
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Initialize services
	qualityService := services.NewQualityService(cfg.QualityThreshold)
	storageService := services.NewStorageService(awsCfg, cfg.S3Bucket)
	orchestrator := services.NewPipelineOrchestrator(
		qualityService,
		storageService,
		cfg.UpscaleScript,
		cfg.UpscaleScale,
	)

	log.Printf("Quality threshold: %.2f", cfg.QualityThreshold)
	log.Printf("Upscale script: %s", cfg.UpscaleScript)
	log.Printf("Upscale scale: %dx", cfg.UpscaleScale)
	log.Printf("S3 bucket: %s", cfg.S3Bucket)

	// Initialize handlers
	imageHandler := handlers.NewImageHandler(orchestrator)

	// Setup routes
	router := mux.NewRouter()

	// CORS middleware
	router.Use(corsMiddleware)

	// Health check
	router.HandleFunc("/api/health", imageHandler.HealthCheck).Methods("GET")

	// Image endpoints
	router.HandleFunc("/api/images/upload", imageHandler.UploadImage).Methods("POST")
	router.HandleFunc("/api/images/list/{folder}", imageHandler.ListProcessed).Methods("GET")
	router.HandleFunc("/api/images/{folder}/{filename}", imageHandler.GetImage).Methods("GET")

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
