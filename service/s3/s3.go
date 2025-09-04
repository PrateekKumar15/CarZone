package s3

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Service handles simple AWS S3 operations
type S3Service struct {
	client     *s3.Client
	bucketName string
}

// NewS3Service creates a new simple S3Service
func NewS3Service(bucketName, region string) (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return &S3Service{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}, nil
}

// UploadBase64Image uploads a base64 image and returns the URL
func (s *S3Service) UploadBase64Image(ctx context.Context, base64Data, fileName string) (string, error) {
	// Clean base64 data
	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Generate unique filename
	uniqueName := fmt.Sprintf("cars/%s-%d.jpg", uuid.New().String(), time.Now().Unix())

	// Upload to S3
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(uniqueName),
		Body:   bytes.NewReader(imageData),
	})
	if err != nil {
		return "", err
	}

	// Return URL
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, uniqueName), nil
}

// DeleteImage deletes an image from S3 using its URL
func (s *S3Service) DeleteImage(ctx context.Context, imageURL string) error {
	// Extract key from URL
	// URL format: https://bucket-name.s3.amazonaws.com/key
	parts := strings.Split(imageURL, "/")
	if len(parts) < 4 {
		return fmt.Errorf("invalid S3 URL format")
	}

	key := strings.Join(parts[3:], "/")

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	return err
}
