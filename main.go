// Package main is the entry point of the application
package main

// Import necessary packages
import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sahilq312/workly/initializer"
	"github.com/sahilq312/workly/routes"
)

func init() {
	initializer.LoadEnvVariale()
	initializer.ConnectPostgresDatabase()
}

func main() {
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Workly",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is healthy",
		})
	})

	routes.AuthRoutes(r)
	routes.PostRoutes(r)
	routes.CompanyRoutes(r)
	routes.UserRoutes(r)
	routes.JobRoutes(r)
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	// Set up graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	// Log server exit
	log.Println("Server exiting")
}
