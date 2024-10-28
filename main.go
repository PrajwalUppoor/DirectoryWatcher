package main

import (
	"context"
	"dirwatcher/db"
	"dirwatcher/endpoints"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Get the Authorization header
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing or invalid"})
// 			c.Abort()
// 			return
// 		}

// 		// Extract the token
// 		token := strings.TrimPrefix(authHeader, "Bearer ")

// 		// Here, you can add logic to validate the token
// 		if token != "your_secret_token" { // Replace with your token validation logic
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

func main() {
	r := gin.Default()
	r.Use(RateLimiter(1, 3))
	// r.Use(AuthMiddleware())

	// API endpoints
	// task endpoints
	// Starts a task with the given configurationId
	r.POST("/task/start", endpoints.StartDirectoryWatchTask)
	// Stops a task with the given configurationId
	r.POST("/task/stop", endpoints.StopDirectoryWatchTask)
	// Get all task details
	r.GET("/task/details", endpoints.GetTaskDetails)

	// Create a new configuration
	r.POST("/configurations", endpoints.CreateConfiguration)
	// Get all configurations
	r.GET("/configurations", endpoints.GetAllConfigurations)
	// Get a single configuration by ID
	r.GET("/configurations/:id", endpoints.GetConfigurationByID)
	// Update a configuration
	r.PUT("/configurations", endpoints.UpdateConfiguration)
	// Delete a configuration by ID
	r.DELETE("/configurations/:id", endpoints.DeleteConfiguration)

	// User profile routes
	r.POST("/profiles", endpoints.CreateProfile)
	r.GET("/profiles/:id", endpoints.GetProfile)
	r.PUT("/profiles/:id", endpoints.UpdateProfile)

	// Itinerary routes
	r.POST("/itineraries", endpoints.CreateItinerary)
	r.GET("/itineraries", endpoints.GetItineraries)
	r.GET("/itineraries/:id", endpoints.GetItinerary)
	r.PUT("/itineraries/:id", endpoints.UpdateItinerary)
	r.DELETE("/itineraries/:id", endpoints.DeleteItinerary)

	// Booking history routes
	r.POST("/booking-history", endpoints.CreateBookingHistory)
	r.GET("/booking-history/:id", endpoints.GetBookingHistory)

	// WebSocket route
	r.GET("/ws/chat", endpoints.WebSocketHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Something went wrong and got recoverd: %v", r)
		}
	}()

	if err := db.Connect(); err != nil {
		log.Printf("Error while connecting to Database: %v", err)
	}

	// Graceful shutdown used to exit the application
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Println("Server started on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server listen error: %v", err)
	}

}
