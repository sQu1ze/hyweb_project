package main

import (
	"context"
	"fmt"
	_ "hyweb-api/docs/swagger"
	"hyweb-api/internal/config"
	"hyweb-api/internal/repository"
	"hyweb-api/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Gin-Gonic API
// @version 1.0
// @description This is a sample server for a Gin-Gonic API.
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database connection
	repository.InitDB()
	defer func() {
		if repository.DB != nil {
			err := repository.DB.Close()
			if err != nil {
				log.Fatalf("Error closing database: %v", err)
			}
			log.Println("Database connection closed")
		}
	}()

	// Setup Gin router
	r := router.SetupRouter()

	// Start HTTP server
	addr := fmt.Sprintf(":%s", config.AppConfig.AppPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Goroutine to start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	log.Printf("Server listening on %s", addr)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline to wait for current requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
