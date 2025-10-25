package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	store := storage.NewStorage()

	// Initialize handlers
	h, err := handlers.NewHandlers(store, Version)
	if err != nil {
		log.Fatal("Failed to initialize handlers:", err)
	}

	// Setup Gin router
	// Set Gin mode from environment configuration (production => Release mode)
	if envConfig.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	var router *gin.Engine
	if envConfig.IsProduction() {
		// Production: use New() to control middleware and logging
		router = gin.New()
		// Recovery middleware to avoid panics crashing the process
		router.Use(gin.Recovery())
		// Minimal, structured-ish logger for production
		router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: func(param gin.LogFormatterParams) string {
				// timestamp, client ip, method, path, status, latency, error (if any)
				return fmt.Sprintf("%s - %s \"%s %s\" %d %s %s\n",
					param.TimeStamp.Format(time.RFC3339),
					param.ClientIP,
					param.Method,
					param.Path,
					param.StatusCode,
					param.Latency,
					param.ErrorMessage,
				)
			},
		}))
	} else {
		// Development: defaults with logger + recovery
		router = gin.Default()
	}

	// Configure trusted proxies (disable for security in development)
	router.SetTrustedProxies([]string{})
	// Return 405 for known path but unsupported method
	router.HandleMethodNotAllowed = true

	// Register routes
	v1 := router.Group("/api/v1")

	// Health check route (no middleware)
	v1.GET("/health", h.GetHealth)
	// Public setup route (no middleware)
	v1.POST("/user/exchangeToken", h.PostSetupExchangeToken)

	// Protected group: all routes here require X-API-Key
	auth := v1.Group("")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)
	auth.POST("/acl", h.PostAcl)

	// Auth routes (with middleware for permission checks)
	auth.POST("/user/resetKey", h.PostUserResetKey)
	auth.POST("/user/generateToken", h.PostUserGenerateToken)

	// Use port from environment configuration
	port := envConfig.Port

	// Start server with graceful shutdown
	addr := ":" + strconv.Itoa(port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in background
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Printf("Server exiting")
}
