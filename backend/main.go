package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	appconfig "visioncloud/config"
	"visioncloud/handlers"
	"visioncloud/services"
)

func main() {
	// Load configuration
	cfg := appconfig.LoadConfig()

	// Initialize AWS configuration
	ctx := context.Background()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
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

	// Initialize handlers
	imageHandler := handlers.NewImageHandler(orchestrator)

	// Set up router with CORS
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/images/upload", imageHandler.UploadImage)
	mux.HandleFunc("/api/images/", imageHandler.GetImage)
	mux.HandleFunc("/api/images/list/", imageHandler.ListProcessed)
	mux.HandleFunc("/api/health", imageHandler.HealthCheck)

	// Wrap with CORS middleware
	handler := withCORS(mux)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting VisionCloud server on port %s", cfg.Port)
		log.Printf("Health check available at http://localhost:%s/api/health", cfg.Port)
		log.Printf("Upload endpoint available at http://localhost:%s/api/images/upload", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("VisionCloud server exited")
}

// withCORS adds CORS headers to all responses
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
