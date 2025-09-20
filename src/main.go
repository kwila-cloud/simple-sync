package main

import (
	"log"
	"os"

	"simple-sync/src/handlers"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize storage
	store := storage.NewMemoryStorage()

	// Initialize handlers
	h := handlers.NewHandlers(store)

	// Setup Gin router
	router := gin.Default()

	// Register routes
	router.GET("/events", h.GetEvents)
	router.POST("/events", h.PostEvents)

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