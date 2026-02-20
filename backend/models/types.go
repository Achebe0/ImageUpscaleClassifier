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
