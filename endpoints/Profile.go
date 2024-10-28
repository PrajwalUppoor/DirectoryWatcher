package endpoints

import (
	"dirwatcher/db"
	"dirwatcher/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler Functions
func CreateProfile(c *gin.Context) {
	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(&profile)
	c.JSON(http.StatusCreated, profile)
}

func GetProfile(c *gin.Context) {
	var profile models.UserProfile
	if err := db.DB.First(&profile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates the profile details based on user preferences
func UpdateProfile(c *gin.Context) {
	var updatedProfile models.UserProfile
	if err := c.ShouldBindJSON(&updatedProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.UserProfile
	if err := db.DB.First(&profile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Update fields as needed
	profile.Name = updatedProfile.Name
	profile.Email = updatedProfile.Email
	profile.Preferences = updatedProfile.Preferences // Assuming Preferences is a struct or a string

	// Save updated profile
	if err := db.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
