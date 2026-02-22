package services

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

// QualityAssessment holds the assessment result
type QualityAssessment struct {
	QualityScore float64 // 0-1 score
	Width        int
	Height       int
	Format       string
}

// QualityService handles image quality assessment
type QualityService struct {
	QualityThreshold float64 // Below this threshold, image needs upscaling
}

// NewQualityService creates a new quality assessment service
func NewQualityService(threshold float64) *QualityService {
	if threshold < 0 {
		threshold = 0
	}
	if threshold > 1 {
		threshold = 1
	}
	return &QualityService{
		QualityThreshold: threshold,
	}
}

// AssessQuality evaluates image quality based on dimensions
// Uses resolution as a proxy for quality (lower resolution = lower quality)
// You can enhance this with more sophisticated metrics like sharpness, blur detection, etc.
func (qs *QualityService) AssessQuality(imageData []byte) (*QualityAssessment, error) {
	img, format, err := image.DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	assessment := &QualityAssessment{
		Width:  img.Width,
		Height: img.Height,
		Format: format,
	}

	// Calculate quality score based on resolution
	// Normalized to 0-1: assumes max good quality at 4K resolution
	maxResolution := float64(3840 * 2160)
	currentResolution := float64(img.Width * img.Height)
	assessment.QualityScore = currentResolution / maxResolution

	// Cap at 1.0
	if assessment.QualityScore > 1.0 {
		assessment.QualityScore = 1.0
	}

	return assessment, nil
}

// NeedsUpscaling returns true if image quality is below threshold
func (qs *QualityService) NeedsUpscaling(assessment *QualityAssessment) bool {
	return assessment.QualityScore < qs.QualityThreshold
}

// IsGoodQuality returns true if image quality is above threshold
func (qs *QualityService) IsGoodQuality(assessment *QualityAssessment) bool {
	return assessment.QualityScore >= qs.QualityThreshold
}
