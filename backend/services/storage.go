package services

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	FolderGoodQuality    = "good_quality"
	FolderUpscaled       = "upscaled"
	FolderCouldntUpscale = "couldn't_upscale"
	FolderProcessing     = "processing"
)

// StorageService handles AWS S3 operations
type StorageService struct {
	client     *s3.Client
	bucket     string
	uploader   *manager.Uploader
	downloader *manager.Downloader
}

// NewStorageService creates a new storage service
func NewStorageService(cfg aws.Config, bucket string) *StorageService {
	client := s3.NewFromConfig(cfg)
	return &StorageService{
		client:     client,
		bucket:     bucket,
		uploader:   manager.NewUploader(client),
		downloader: manager.NewDownloader(client),
	}
}

// UploadImage uploads an image to the specified S3 folder
func (ss *StorageService) UploadImage(ctx context.Context, folder, objectKey string, data []byte) (string, error) {
	fullKey := fmt.Sprintf("%s/%s", folder, objectKey)

	result, err := ss.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(fullKey),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	return result.Location, nil
}

// DownloadImage downloads an image from S3
func (ss *StorageService) DownloadImage(ctx context.Context, folder, objectKey string) ([]byte, error) {
	fullKey := fmt.Sprintf("%s/%s", folder, objectKey)

	buf := manager.NewWriteAtBuffer([]byte{})
	_, err := ss.downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(fullKey),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download image from S3: %w", err)
	}

	return buf.Bytes(), nil
}

// GetImageURL generates the S3 URL for an image
func (ss *StorageService) GetImageURL(folder, objectKey string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s", ss.bucket, folder, objectKey)
}

// ListImages lists all images in a folder
func (ss *StorageService) ListImages(ctx context.Context, folder string) ([]string, error) {
	prefix := fmt.Sprintf("%s/", folder)
	paginator := s3.NewListObjectsV2Paginator(ss.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(ss.bucket),
		Prefix: aws.String(prefix),
	})

	var images []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			images = append(images, *obj.Key)
		}
	}

	return images, nil
}

// DeleteImage deletes an image from S3
func (ss *StorageService) DeleteImage(ctx context.Context, folder, objectKey string) error {
	fullKey := fmt.Sprintf("%s/%s", folder, objectKey)

	_, err := ss.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(fullKey),
	})

	if err != nil {
		return fmt.Errorf("failed to delete image from S3: %w", err)
	}

	return nil
}
