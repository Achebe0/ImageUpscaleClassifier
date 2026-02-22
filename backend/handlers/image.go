package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"visioncloud/services"
)

// ImageHandler handles image-related HTTP requests
type ImageHandler struct {
	orchestrator *services.PipelineOrchestrator
}

// NewImageHandler creates a new image handler
func NewImageHandler(orchestrator *services.PipelineOrchestrator) *ImageHandler {
	return &ImageHandler{
		orchestrator: orchestrator,
	}
}

// UploadImageRequest represents the upload request
type UploadImageRequest struct {
	Quality float64 `json:"quality_threshold"` // Optional, defaults to service default
}

// UploadImageResponse represents the upload response
type UploadImageResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Result  *services.ProcessingResult `json:"result,omitempty"`
	Error   string                     `json:"error,omitempty"`
}

// BatchUploadResponse represents batch upload response
type BatchUploadResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Results []*services.ProcessingResult `json:"results,omitempty"`
	Error   string                       `json:"error,omitempty"`
}

// UploadImage handles single image upload and processing
// POST /api/images/upload
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse form data with max 1GB file size
	if err := r.ParseMultipartForm(1024 * 1024 * 1024); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadImageResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse form: %v", err),
		})
		return
	}

	// Get file from form
	file, header, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadImageResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get image file: %v", err),
		})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadImageResponse{
			Success: false,
			Error:   "Invalid image format. Supported: jpg, jpeg, png, webp",
		})
		return
	}

	// Read file content
	imageData, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UploadImageResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to read image: %v", err),
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	// Process image through pipeline
	result := h.orchestrator.ProcessImage(ctx, imageData, header.Filename)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UploadImageResponse{
		Success: result.Status == "success",
		Message: fmt.Sprintf("Image processed: %s", result.Status),
		Result:  result,
	})
}

// GetImage retrieves a processed image information
// GET /api/images/{folder}/{filename}
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// In a real implementation, you'd extract folder and filename from URL params
	// This is a placeholder response

	response := map[string]interface{}{
		"success": true,
		"message": "Image retrieval not yet implemented",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ListProcessed lists all processed images in a folder
// GET /api/images/list/{folder}
func (h *ImageHandler) ListProcessed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// In a real implementation, you'd extract folder from URL params
	// This is a placeholder response

	response := map[string]interface{}{
		"success": true,
		"message": "Image listing not yet implemented",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthCheck returns the health status
// GET /api/health
func (h *ImageHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
