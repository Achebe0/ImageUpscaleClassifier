package models

type UpscaleRequest struct {
	Image   []byte `json:"image"`
	Scale   int    `json:"scale"`
	ModelID string `json:"model_id"`
}

type UpscaleResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ProcessingStatus struct {
	Status       string  `json:"status"`        // success, skipped, error
	Folder       string  `json:"folder"`        // good_quality, upscaled, couldn't_upscale
	QualityScore float64 `json:"quality_score"`
	S3URL        string  `json:"s3_url,omitempty"`
	Error        string  `json:"error,omitempty"`
}
