package endpoints

import (
	"dirwatcher/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SendNotification(c *gin.Context) {
	var request struct {
		UserID  int    `json:"user_id"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := services.NewNotificationService().SendNotification(request.UserID, request.Message)
	c.JSON(http.StatusOK, notification)
}

func GetUserNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	userId, err := strconv.Atoi(userID)
	if err != nil {
		notifications := services.NewNotificationService().GetNotifications(userId)
		c.JSON(http.StatusOK, notifications)
	}
	c.JSON(http.StatusOK, nil)

}
