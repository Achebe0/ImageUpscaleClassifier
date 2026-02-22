package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// ProcessingResult contains the result of image processing
type ProcessingResult struct {
	OriginalKey   string    `json:"original_key"`
	Status        string    `json:"status"` // success, skipped, error
	Folder        string    `json:"folder"`
	S3URL         string    `json:"s3_url,omitempty"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	ProcessedAt   time.Time `json:"processed_at"`
	QualityScore  float64   `json:"quality_score"`
	UpscaleScale  int       `json:"upscale_scale,omitempty"`
}

// PipelineOrchestrator orchestrates the image upscaling pipeline
type PipelineOrchestrator struct {
	qualityService *QualityService
	storageService *StorageService
	upscaleScript  string
	tempDir        string
	upscaleScale   int
}

// NewPipelineOrchestrator creates a new pipeline orchestrator
func NewPipelineOrchestrator(
	qualityService *QualityService,
	storageService *StorageService,
	upscaleScript string,
	upscaleScale int,
) *PipelineOrchestrator {
	if upscaleScale <= 0 {
		upscaleScale = 2 // default 2x upscaling
	}

	tempDir := os.TempDir()

	return &PipelineOrchestrator{
		qualityService: qualityService,
		storageService: storageService,
		upscaleScript:  upscaleScript,
		tempDir:        filepath.Join(tempDir, "visioncloud"),
		upscaleScale:   upscaleScale,
	}
}

// ProcessImage processes a single image through the pipeline
func (po *PipelineOrchestrator) ProcessImage(ctx context.Context, imageData []byte, objectKey string) *ProcessingResult {
	result := &ProcessingResult{
		OriginalKey: objectKey,
		ProcessedAt: time.Now(),
		Status:      "error",
	}

	// Step 1: Assess image quality
	assessment, err := po.qualityService.AssessQuality(imageData)
	if err != nil {
		result.Status = "error"
		result.Folder = FolderCouldntUpscale
		result.ErrorMessage = fmt.Sprintf("Quality assessment failed: %v", err)
		po.storageService.UploadImage(ctx, result.Folder, objectKey, imageData)
		return result
	}

	result.QualityScore = assessment.QualityScore

	// Step 2: Check if image is already good quality
	if po.qualityService.IsGoodQuality(assessment) {
		result.Status = "success"
		result.Folder = FolderGoodQuality
		url, err := po.storageService.UploadImage(ctx, result.Folder, objectKey, imageData)
		if err != nil {
			result.Status = "error"
			result.ErrorMessage = fmt.Sprintf("Failed to upload good quality image: %v", err)
			return result
		}
		result.S3URL = url
		return result
	}

	// Step 3: Image needs upscaling - attempt upscale
	result.UpscaleScale = po.upscaleScale
	upscaledData, err := po.upscaleImage(ctx, imageData)
	if err != nil {
		result.Status = "error"
		result.Folder = FolderCouldntUpscale
		result.ErrorMessage = fmt.Sprintf("Upscaling failed: %v", err)

		// Still upload the original image to couldn't_upscale folder
		po.storageService.UploadImage(ctx, result.Folder, objectKey, imageData)
		return result
	}

	// Step 4: Upload upscaled image
	result.Status = "success"
	result.Folder = FolderUpscaled
	url, err := po.storageService.UploadImage(ctx, result.Folder, objectKey, upscaledData)
	if err != nil {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("Failed to upload upscaled image: %v", err)
		return result
	}

	result.S3URL = url
	return result
}

// upscaleImage calls the Python upscaling script via subprocess
func (po *PipelineOrchestrator) upscaleImage(ctx context.Context, imageData []byte) ([]byte, error) {
	// Create temporary directory if it doesn't exist
	if err := os.MkdirAll(po.tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Generate unique filenames for input/output
	timestamp := time.Now().UnixNano()
	inputPath := filepath.Join(po.tempDir, fmt.Sprintf("input_%d.png", timestamp))
	outputPath := filepath.Join(po.tempDir, fmt.Sprintf("output_%d.png", timestamp))

	// Save image to temp file
	if err := os.WriteFile(inputPath, imageData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp image: %w", err)
	}
	defer os.Remove(inputPath)
	defer os.Remove(outputPath)

	// Call Python upscaling script directly
	cmd := exec.CommandContext(ctx, "python", po.upscaleScript,
		"--input", inputPath,
		"--output", outputPath,
		"--scale", strconv.Itoa(po.upscaleScale),
	)

	// Capture stderr for debugging
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command with timeout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("upscaling script failed: %w, stderr: %s", err, stderr.String())
	}

	// Check if output file was created
	if _, err := os.Stat(outputPath); err != nil {
		return nil, fmt.Errorf("upscaled image not created: %w", err)
	}

	// Read the upscaled image
	upscaledData, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read upscaled image: %w", err)
	}

	return upscaledData, nil
}

// ProcessImageBatch processes multiple images (useful for batch operations)
func (po *PipelineOrchestrator) ProcessImageBatch(ctx context.Context, images map[string][]byte) []*ProcessingResult {
	results := make([]*ProcessingResult, 0, len(images))

	for objectKey, imageData := range images {
		result := po.ProcessImage(ctx, imageData, objectKey)
		results = append(results, result)
	}

	return results
}
