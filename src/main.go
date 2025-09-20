package main

import (
	"log"
	"os"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Print version information
	log.Printf("Simple-Sync v%s (build: %s)", Version, BuildTime)

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-in-production" // Default for development
		log.Printf("Warning: Using default JWT secret. Set JWT_SECRET environment variable for production.")
	}

	// Initialize storage
	store := storage.NewMemoryStorage()

	// Initialize handlers
	h := handlers.NewHandlers(store, jwtSecret)

	// Setup Gin router
	router := gin.Default()

	// Configure trusted proxies (disable for security in development)
	router.SetTrustedProxies([]string{})

	// Register routes
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Auth routes (no middleware)
	router.POST("/auth/token", h.PostAuthToken)

	// Get port from environment or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
