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
)

func main() {
	r := gin.Default()

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
