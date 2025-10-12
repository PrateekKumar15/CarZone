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
	"github.com/PrateekKumar15/CarZone/service/cloudinary"
	"github.com/gorilla/mux"
)

// ImageUploadMiddleware handles image uploads to Cloudinary
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
			println("📸 Processing", len(carRequest.Images), "images...")
			cloudinaryService, err := cloudinary.NewCloudinaryService(
				getEnv("CLOUDINARY_CLOUD_NAME", ""),
				getEnv("CLOUDINARY_API_KEY", ""),
				getEnv("CLOUDINARY_API_SECRET", ""),
				getEnv("CLOUDINARY_FOLDER", "carzone/cars"),
			)
			if err != nil {
				println("❌ Failed to initialize Cloudinary service:", err.Error())
			} else {
				println("✅ Cloudinary service initialized successfully")
				for i, img := range carRequest.Images {
					if !isURL(img) { // If not already a URL, try to upload
						println("📤 Uploading image", i+1, "- Size:", len(img), "bytes")
						if url, err := cloudinaryService.UploadBase64Image(r.Context(), img, "car_image.jpg"); err == nil {
							println("✅ Image", i+1, "uploaded successfully:", url)
							carRequest.Images[i] = url
						} else {
							println("❌ Failed to upload image", i+1, ":", err.Error())
						}
					} else {
						println("⏭️  Image", i+1, "is already a URL, skipping")
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

// cleanupOldImages removes old images from Cloudinary when updating a car
func cleanupOldImages(ctx context.Context, carID string, newImages []string) {
	// Get old car images from database
	oldImages := GetCarImages(ctx, carID)
	if len(oldImages) == 0 {
		return
	}

	cloudinaryService, err := cloudinary.NewCloudinaryService(
		getEnv("CLOUDINARY_CLOUD_NAME", ""),
		getEnv("CLOUDINARY_API_KEY", ""),
		getEnv("CLOUDINARY_API_SECRET", ""),
		getEnv("CLOUDINARY_FOLDER", "carzone/cars"),
	)
	if err != nil {
		return
	}

	// Find images that are being removed (exist in old but not in new)
	for _, oldImage := range oldImages {
		if cloudinary.IsCloudinaryURL(oldImage) && !contains(newImages, oldImage) {
			cloudinaryService.DeleteImage(ctx, oldImage)
		}
	}
}

// DeleteCarImages removes all images for a deleted car
func DeleteCarImages(ctx context.Context, imageURLs []string) {
	if len(imageURLs) == 0 {
		return
	}

	cloudinaryService, err := cloudinary.NewCloudinaryService(
		getEnv("CLOUDINARY_CLOUD_NAME", ""),
		getEnv("CLOUDINARY_API_KEY", ""),
		getEnv("CLOUDINARY_API_SECRET", ""),
		getEnv("CLOUDINARY_FOLDER", "carzone/cars"),
	)
	if err != nil {
		return
	}

	for _, imageURL := range imageURLs {
		if cloudinary.IsCloudinaryURL(imageURL) {
			cloudinaryService.DeleteImage(ctx, imageURL)
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
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func isCloudinaryURL(str string) bool {
	return strings.Contains(str, "res.cloudinary.com")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
