package cloudinary

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

// CloudinaryService handles Cloudinary operations for image uploads
type CloudinaryService struct {
	cld    *cloudinary.Cloudinary
	folder string
}

// NewCloudinaryService creates a new CloudinaryService
func NewCloudinaryService(cloudName, apiKey, apiSecret, folder string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	if folder == "" {
		folder = "carzone/cars"
	}

	return &CloudinaryService{
		cld:    cld,
		folder: folder,
	}, nil
}

// UploadBase64Image uploads a base64 image to Cloudinary and returns the secure URL
func (s *CloudinaryService) UploadBase64Image(ctx context.Context, base64Data, fileName string) (string, error) {
	// Clean base64 data - remove data:image/xxx;base64, prefix if present
	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}

	// Generate unique public ID for the image
	publicID := fmt.Sprintf("%s_%d", uuid.New().String(), time.Now().Unix())

	// Create base64 data URI for Cloudinary
	// Cloudinary accepts data URIs in the format: data:image/png;base64,<base64_data>
	dataURI := fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(imageData))

	// Upload to Cloudinary using data URI
	uploadResult, err := s.cld.Upload.Upload(ctx, dataURI, uploader.UploadParams{
		PublicID:     publicID,
		Folder:       s.folder,
		ResourceType: "image",
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	// Return the secure URL
	return uploadResult.SecureURL, nil
}

// DeleteImage deletes an image from Cloudinary using its URL
func (s *CloudinaryService) DeleteImage(ctx context.Context, imageURL string) error {
	// Extract public ID from Cloudinary URL
	// URL format: https://res.cloudinary.com/{cloud_name}/image/upload/v{version}/{folder}/{public_id}.{format}
	publicID := extractPublicIDFromURL(imageURL, s.folder)
	if publicID == "" {
		return fmt.Errorf("invalid Cloudinary URL format: %s", imageURL)
	}

	// Delete from Cloudinary
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})

	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
	}

	return nil
}

// extractPublicIDFromURL extracts the public ID from a Cloudinary URL
// Example URL: https://res.cloudinary.com/demo/image/upload/v1234567890/carzone/cars/abc-123.jpg
// Returns: carzone/cars/abc-123
func extractPublicIDFromURL(url, folder string) string {
	// Find the position of "/upload/"
	uploadIndex := strings.Index(url, "/upload/")
	if uploadIndex == -1 {
		return ""
	}

	// Get everything after "/upload/v{version}/"
	afterUpload := url[uploadIndex+8:] // +8 to skip "/upload/"

	// Skip the version number (e.g., "v1234567890/")
	slashIndex := strings.Index(afterUpload, "/")
	if slashIndex == -1 {
		return ""
	}
	afterVersion := afterUpload[slashIndex+1:]

	// Remove file extension
	lastDot := strings.LastIndex(afterVersion, ".")
	if lastDot != -1 {
		afterVersion = afterVersion[:lastDot]
	}

	return afterVersion
}

// IsCloudinaryURL checks if a URL is a Cloudinary URL
func IsCloudinaryURL(url string) bool {
	return strings.Contains(url, "res.cloudinary.com")
}
