package endpoints

import (
	"dirwatcher/db"
	"dirwatcher/models"
	"dirwatcher/services"
	"dirwatcher/structures"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
*

		This method will Create new Configurations
		@Accept json
	        @Produce json
		@Success 200
		@Router /configurations/:id [POST]
*/
func CreateConfiguration(c *gin.Context) {
	// Parse request body and create a new Configuration
	configsDto := structures.Configuration{}
	if err := c.BindJSON(&configsDto); err != nil {
		log.Printf("[CreateConfiguration] Error while binding the json :%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while creating configuration"})
		return
	}

	configsDto.TimeInterval = configsDto.TimeInterval * time.Second
	configModel := models.Configurations{MonitoredDirectory: configsDto.MonitoredDirectory, MagicString: configsDto.MagicString, TimeInterval: configsDto.TimeInterval}
	// Save it to data store
	if err := db.DB.Create(&configModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while creating configuration"})
		return
	}

	responseDto := structures.ConfigurationResponse{ConfigurationId: configModel.ID, MonitoredDirectory: configsDto.MonitoredDirectory, TimeInterval: configModel.TimeInterval, MagicString: configsDto.MagicString}

	c.JSON(http.StatusOK, gin.H{"message": "Configutaion saved succesfully", "configurations": responseDto})
}

/*
*

		This method will fetch all the Configurations
	        @Produce json
		@Success 200
		@Router /configurations [GET]
*/
func GetAllConfigurations(c *gin.Context) {
	// Retrieve all configurations
	responseDto, err := services.GetAllConfigurations()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching configuration"})
		return
	}

	// Return them as JSON
	c.JSON(http.StatusOK, gin.H{"configurations": responseDto})
}

/*
*

		This method will Fetch Configurations by passing Id as Path Parameter
		@PathParam id
	        @Produce json
		@Success 200
		@Router /configurations/:id [GET]
*/
func GetConfigurationByID(c *gin.Context) {
	// Get configuration by ID
	var config models.Configurations
	id := c.Params.ByName("id")
	if err := db.DB.Where("id = ?", id).First(&config).Error; err != nil {
		log.Printf("[GetConfigurationByID] Error while getting configs by Id %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("configuration for Id: %v not found", id)})
		return
	}

	// Return it as JSON
	c.JSON(http.StatusOK, gin.H{"configuration": structures.Configuration{MonitoredDirectory: config.MonitoredDirectory, TimeInterval: config.TimeInterval, MagicString: config.MagicString}})
}

/*
*

		This method will Update Configurations by passing ConfigurationId along with Configuration Details in Request Body.
		@Accept json
	        @Produce json
		@Success 200
		@Router /configurations/ [PUT]
*/
func UpdateConfiguration(c *gin.Context) {

	configsDto := structures.ConfigurationResponse{}
	if err := c.BindJSON(&configsDto); err != nil {
		log.Printf("[UpdateConfiguration] Error while binding the json :%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Update Configurations failed :%v", err)})
		return
	}

	configsDto.TimeInterval = configsDto.TimeInterval * time.Second
	configModel := models.Configurations{MonitoredDirectory: configsDto.MonitoredDirectory, MagicString: configsDto.MagicString, TimeInterval: configsDto.TimeInterval}
	configModel.ID = configsDto.ConfigurationId
	// Save the updated configuration
	if err := db.DB.Save(&configModel).Error; err != nil {
		log.Printf("[UpdateConfiguration] Error while updating configuration to DB :%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while updating configuration"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

/*
*

		This method will Delete Configurations by passing Id as Path Parameter
		@PathParam id
	        @Produce json
		@Success 200
		@Router /configurations/:id [DELETE]
*/
func DeleteConfiguration(c *gin.Context) {
	id := c.Params.ByName("id")
	var conf models.Configurations
	// Delete configuration by ID
	// Remove it from data store
	if err := db.DB.Where("id = ?", id).Delete(&conf).Error; err != nil {
		log.Printf("[UpdateConfiguration] Error while deleting configuration to DB :%v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while deleting configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}
