package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
)

type AWSConfig struct {
	S3Client     *aws.Config
	ApprovedBucket string
	RejectedBucket string
	DynamoTable    string
	Region         string
}

func LoadConfig() (*aws.Config, error) {
	_ = godotenv.Load() // Load .env file if present

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	return &cfg, nil
}

func GetApprovedBucket() string {
	return os.Getenv("S3_APPROVED_BUCKET")
}

func GetRejectedBucket() string {
	return os.Getenv("S3_REJECTED_BUCKET")
}

func GetDynamoTable() string {
	return os.Getenv("DYNAMODB_TABLE")
}
