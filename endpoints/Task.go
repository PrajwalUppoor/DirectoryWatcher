package endpoints

import (
	"dirwatcher/db"
	"dirwatcher/models"
	"dirwatcher/services"
	"dirwatcher/structures"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	watchers = make(map[uint]chan struct{}) // Map to store directory watchers

)

/*
*

		This method will Start Task for given ConfigurationId
		@Accept json
	        @Produce json
		@Success 200
		@Router /task/start [POST]
*/
func StartDirectoryWatchTask(c *gin.Context) {

	stopChan := make(chan struct{})
	configsDto := structures.TaskDto{}
	if err := c.BindJSON(&configsDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Bad Request : %v", configsDto),
		})
		return
	}
	configModel := models.Configurations{}
	if err := db.DB.Where("id = ?", configsDto.ConfigurationId).First(&configModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Configuration not found for ID %v: %s", configsDto.ConfigurationId, err),
		})
		return
	}

	configs := structures.Configuration{MonitoredDirectory: configModel.MonitoredDirectory, TimeInterval: configModel.TimeInterval, MagicString: configModel.MagicString}

	go services.DirectoryWatcher(configs, stopChan)

	// Store the stop channel in the map
	watchers[configModel.ID] = stopChan

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Started watching directory: %s", configs.MonitoredDirectory),
	})
}

/*
*

		This method will Stop Task for given ConfigurationId
		@Accept json
	        @Produce json
		@Success 200
		@Router /task/stop [POST]
*/
func StopDirectoryWatchTask(c *gin.Context) {

	configsDto := structures.TaskDto{}
	if err := c.BindJSON(&configsDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Error while Binding the json : %s", err),
		})
		return
	}
	configModel := models.Configurations{}
	if err := db.DB.Where("id = ?", configsDto.ConfigurationId).First(&configModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Configuration not found for ID %v: %s", configsDto.ConfigurationId, err),
		})
		return
	}

	configs := structures.Configuration{MonitoredDirectory: configModel.MonitoredDirectory, TimeInterval: configModel.TimeInterval, MagicString: configModel.MagicString}

	if stopChan, ok := watchers[configModel.ID]; ok {
		close(stopChan)
		delete(watchers, configModel.ID)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Stopped watching directory: %s", configs.MonitoredDirectory),
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Watcher for Directory %s not found in the list of tasks which are started", configs.MonitoredDirectory),
		})
	}
}

/*
*

		This method will fetch all task details
		@Accept json
	        @Produce json
		@Success 200
		@Router /task/details [POST]
*/
func GetTaskDetails(c *gin.Context) {

	taskModels := []models.Task{}

	if err := db.DB.Preload("Configurations").Find(&taskModels).Error; err != nil {
		log.Printf("[GetTaskDetails] Could not fetch all the tasks:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while getting task"})
		return
	}

	responseDto := []structures.Task{}

	for _, task := range taskModels {
		services.TasksMutex.RLock()
		if val, ok := services.Tasks[uint32(task.ID)]; ok {
			responseDto = append(responseDto, val)
		} else {
			taskConfig := structures.Configuration{MonitoredDirectory: task.Configurations.MonitoredDirectory, TimeInterval: task.Configurations.TimeInterval, MagicString: task.Configurations.MagicString}
			responseDto = append(responseDto, structures.Task{TaskId: uint32(task.ID), StartTime: task.StartTime, EndTime: task.EndTime, TotalRuntime: task.TotalRuntime, FilesAdded: strings.Split(task.FilesAdded, ","), FilesDeleted: strings.Split(task.FilesDeleted, ","), Configuration: taskConfig, MagicStringOccurrences: task.MagicStringOccurrences, Status: structures.Status(task.Status)})
		}
		services.TasksMutex.RUnlock()
	}
	c.JSON(http.StatusOK, gin.H{"tasks": responseDto})
}
