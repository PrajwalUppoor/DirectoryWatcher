package endpoints

import (
	"dirwatcher/db"
	"dirwatcher/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBookingHistory(c *gin.Context) {
	var booking models.BookingHistory
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(&booking)
	c.JSON(http.StatusCreated, booking)
}

func GetBookingHistory(c *gin.Context) {
	var booking models.BookingHistory
	if err := db.DB.First(&booking, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, booking)
}
