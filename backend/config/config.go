package config

import (
	"os"
)

type Config struct {
	Port      string
	S3Bucket  string
	ModelPath string
	AWSRegion string
}

func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		S3Bucket:  getEnv("S3_BUCKET", "visioncloud-bucket"),
		ModelPath: getEnv("MODEL_PATH", "./models/upscaler.pth"),
		AWSRegion: getEnv("AWS_REGION", "us-east-1"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
