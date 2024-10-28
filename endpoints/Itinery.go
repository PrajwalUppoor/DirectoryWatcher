package endpoints

import (
	"dirwatcher/db"
	"dirwatcher/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateItinerary creates a new itinerary
func CreateItinerary(c *gin.Context) {
	var itinerary models.Itinerary
	if err := c.ShouldBindJSON(&itinerary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(&itinerary)
	c.JSON(http.StatusCreated, itinerary)
}

// GetItinerary retrieves an itinerary by ID
func GetItinerary(c *gin.Context) {
	var itinerary models.Itinerary
	if err := db.DB.First(&itinerary, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Itinerary not found"})
		return
	}
	c.JSON(http.StatusOK, itinerary)
}

// UpdateItinerary updates an existing itinerary
func UpdateItinerary(c *gin.Context) {
	var updatedItinerary models.Itinerary
	if err := c.ShouldBindJSON(&updatedItinerary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var itinerary models.Itinerary
	if err := db.DB.First(&itinerary, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Itinerary not found"})
		return
	}

	// Update fields
	itinerary.TripName = updatedItinerary.TripName
	itinerary.StartDate = updatedItinerary.StartDate
	itinerary.EndDate = updatedItinerary.EndDate
	itinerary.Destination = updatedItinerary.Destination
	itinerary.TravelMode = updatedItinerary.TravelMode
	itinerary.Accommodation = updatedItinerary.Accommodation
	itinerary.TotalCost = updatedItinerary.TotalCost
	itinerary.Activities = updatedItinerary.Activities
	itinerary.Preferences = updatedItinerary.Preferences
	itinerary.Notes = updatedItinerary.Notes
	itinerary.Status = updatedItinerary.Status
	itinerary.Demographics = updatedItinerary.Demographics

	// Save updated itinerary
	if err := db.DB.Save(&itinerary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update itinerary"})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// DeleteItinerary deletes an existing itinerary
func DeleteItinerary(c *gin.Context) {
	if err := db.DB.Delete(&models.Itinerary{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete itinerary"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetItineraries retrieves all itineraries
func GetItineraries(c *gin.Context) {
	var itineraries []models.Itinerary
	if err := db.DB.Find(&itineraries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch itineraries"})
		return
	}

	c.JSON(http.StatusOK, itineraries)
}
