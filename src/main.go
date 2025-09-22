package main

import (
	"log"
	"os"
	"strconv"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Print version information
	log.Printf("Simple-Sync v%s (build: %s)", Version, BuildTime)
	log.Printf("Starting application...")

	// Load environment configuration
	log.Printf("Loading environment configuration...")
	envConfig := models.NewEnvironmentConfiguration()
	if err := envConfig.LoadFromEnv(os.Getenv); err != nil {
		log.Fatal("Environment configuration error:", err)
	}

	if err := envConfig.Validate(); err != nil {
		log.Fatal("Environment validation error:", err)
	}
	log.Printf("Environment loaded: PORT=%d, ENV=%s", envConfig.Port, envConfig.Environment)

	// Initialize storage
	store := storage.NewMemoryStorage()

	// Initialize handlers
	h := handlers.NewHandlers(store, envConfig.JWT_SECRET, Version)

	// Setup Gin router
	router := gin.Default()

	// Configure trusted proxies (disable for security in development)
	router.SetTrustedProxies([]string{})

	// Register routes
	v1 := router.Group("/api/v1")

	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Auth routes (no middleware)
	v1.POST("/auth/token", h.PostAuthToken)

	// Health check route (no middleware)
	v1.GET("/health", h.GetHealth)

	// Use port from environment configuration
	port := envConfig.Port

	// Start server
	log.Printf("Starting server on port %d", port)
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
