package models

import "time"

type QualityLabel string

const (
	HighQuality QualityLabel = "HighQuality"
	LowQuality  QualityLabel = "LowQuality"
)

type ImageMetadata struct {
	ImageID      string       `json:"image_id" dynamodbav:"image_id"`
	FileName     string       `json:"file_name" dynamodbav:"file_name"`
	Bucket       string       `json:"bucket" dynamodbav:"bucket"`
	QualityScore float64      `json:"quality_score" dynamodbav:"quality_score"`
	QualityLabel QualityLabel `json:"quality_label" dynamodbav:"quality_label"`
	Timestamp    time.Time    `json:"timestamp" dynamodbav:"timestamp"`
	URL          string       `json:"url,omitempty" dynamodbav:"-"`
}

type UploadResponse struct {
	ImageID      string       `json:"image_id"`
	QualityScore float64      `json:"quality_score"`
	QualityLabel QualityLabel `json:"quality_label"`
	Bucket       string       `json:"bucket"`
}
