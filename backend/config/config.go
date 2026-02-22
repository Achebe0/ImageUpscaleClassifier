package config

import (
	"os"
	"strconv"
	"log"
	"path/filepath"
)

type Config struct {
	Port             string
	S3Bucket         string
	ModelPath        string
	AWSRegion        string
	QualityThreshold float64
	UpscaleScript    string
	UpscaleScale     int
}

func init() {
	// Try to load .env file from current directory or parent directories
	loadEnvFile()
}

func loadEnvFile() {
	envPath := ".env"
	
	// Try current directory first
	if _, err := os.Stat(envPath); err != nil {
		// Try parent directory
		envPath = filepath.Join("..", ".env")
		if _, err := os.Stat(envPath); err == nil {
			loadDotEnv(envPath)
			return
		}
		// Try two levels up
		envPath = filepath.Join("..", "..", ".env")
		if _, err := os.Stat(envPath); err == nil {
			loadDotEnv(envPath)
			return
		}
		// .env file not found, that's okay
		return
	}
	
	loadDotEnv(envPath)
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Warning: Could not open .env file at %s: %v", path, err)
		return
	}
	defer file.Close()

	// Simple .env parser
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n == 0 || err != nil {
			break
		}
		
		lines := string(buf[:n])
		parseDotEnvLines(lines)
	}
}

func parseDotEnvLines(content string) {
	lines := splitLines(content)
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		
		parts := splitKeyValue(line)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			// Only set if not already set
			if _, exists := os.LookupEnv(key); !exists {
				os.Setenv(key, value)
			}
		}
	}
}

func splitLines(content string) []string {
	var lines []string
	var current string
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if len(current) > 0 {
		lines = append(lines, current)
	}
	return lines
}

func splitKeyValue(line string) []string {
	for i, char := range line {
		if char == '=' {
			return []string{line[:i], line[i+1:]}
		}
	}
	return []string{}
}

func LoadConfig() *Config {
	qualityThreshold := 0.5
	if val := os.Getenv("QUALITY_THRESHOLD"); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			qualityThreshold = f
		}
	}

	upscaleScale := 2
	if val := os.Getenv("UPSCALE_SCALE"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			upscaleScale = i
		}
	}

	return &Config{
		Port:             getEnv("PORT", "8080"),
		S3Bucket:         getEnv("S3_BUCKET", "visioncloud-bucket"),
		ModelPath:        getEnv("MODEL_PATH", "./models/upscaler.pth"),
		AWSRegion:        getEnv("AWS_REGION", "us-east-1"),
		QualityThreshold: qualityThreshold,
		UpscaleScript:    getEnv("UPSCALE_SCRIPT", "../python/upscaler/upscale.py"),
		UpscaleScale:     upscaleScale,
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
