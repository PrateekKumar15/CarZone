package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/service/s3"
	"github.com/gorilla/mux"
)

// ImageUploadMiddleware handles simple image uploads to S3
func ImageUploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only process POST and PUT requests
		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Try to parse as CarRequest
		var carRequest models.CarRequest
		if err := json.Unmarshal(body, &carRequest); err != nil {
			// If it's not a valid CarRequest, just pass it through
			r.Body = io.NopCloser(bytes.NewReader(body))
			next.ServeHTTP(w, r)
			return
		}

		// Handle image uploads and cleanup
		if r.Method == "PUT" {
			// For updates, cleanup old images if needed
			vars := mux.Vars(r)
			carID := vars["id"]
			if carID != "" {
				cleanupOldImages(r.Context(), carID, carRequest.Images)
			}
		}

		// If there are images and they look like base64, upload them
		if len(carRequest.Images) > 0 {
			s3Service, err := s3.NewS3Service(
				getEnv("AWS_S3_BUCKET_NAME", "carzone-images"),
				getEnv("AWS_REGION", "us-east-1"),
			)
			if err == nil {
				for i, img := range carRequest.Images {
					if !isURL(img) { // If not already a URL, try to upload
						if url, err := s3Service.UploadBase64Image(r.Context(), img, "car_image.jpg"); err == nil {
							carRequest.Images[i] = url
						}
					}
				}
			}
		}

		// Put the (possibly modified) request back
		newBody, _ := json.Marshal(carRequest)
		r.Body = io.NopCloser(bytes.NewReader(newBody))

		next.ServeHTTP(w, r)
	})
}

// cleanupOldImages removes old images from S3 when updating a car
func cleanupOldImages(ctx context.Context, carID string, newImages []string) {
	// Get old car images from database
	oldImages := GetCarImages(ctx, carID)
	if len(oldImages) == 0 {
		return
	}

	s3Service, err := s3.NewS3Service(
		getEnv("AWS_S3_BUCKET_NAME", "carzone-images"),
		getEnv("AWS_REGION", "us-east-1"),
	)
	if err != nil {
		return
	}

	// Find images that are being removed (exist in old but not in new)
	for _, oldImage := range oldImages {
		if isS3URL(oldImage) && !contains(newImages, oldImage) {
			s3Service.DeleteImage(ctx, oldImage)
		}
	}
}

// DeleteCarImages removes all images for a deleted car
func DeleteCarImages(ctx context.Context, imageURLs []string) {
	if len(imageURLs) == 0 {
		return
	}

	s3Service, err := s3.NewS3Service(
		getEnv("AWS_S3_BUCKET_NAME", "carzone-images"),
		getEnv("AWS_REGION", "us-east-1"),
	)
	if err != nil {
		return
	}

	for _, imageURL := range imageURLs {
		if isS3URL(imageURL) {
			s3Service.DeleteImage(ctx, imageURL)
		}
	}
}

// GetCarImages is a placeholder function for getting existing car images
// This should be implemented when you update the car store
func GetCarImages(ctx context.Context, carID string) []string {
	// TODO: Implement this function to get car images from database
	// This will be used by cleanupOldImages function

	// Placeholder implementation - returns empty slice
	return []string{}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func isURL(str string) bool {
	return len(str) > 4 && (str[:4] == "http" || str[:5] == "https")
}

func isS3URL(str string) bool {
	return strings.Contains(str, "s3.amazonaws.com") || strings.Contains(str, ".s3.")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
