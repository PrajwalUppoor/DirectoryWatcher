package services

import (
	"dirwatcher/db"
	"dirwatcher/models"
	"dirwatcher/structures"
	"log"
)

/*
*

	This method is used to fetch all the configuartions from the Database
	@return []structures.ConfigurationResponse
	@return error
*/
func GetAllConfigurations() ([]structures.ConfigurationResponse, error) {

	configs := []models.Configurations{}

	if err := db.DB.Find(&configs).Error; err != nil {
		log.Printf("[GetAllConfigurations] Could not fetch all the configurations:%v", err)
		return nil, err
	}

	configurations := []structures.ConfigurationResponse{}

	for _, conf := range configs {
		configurations = append(configurations, structures.ConfigurationResponse{ConfigurationId: conf.ID, MagicString: conf.MagicString, MonitoredDirectory: conf.MonitoredDirectory, TimeInterval: conf.TimeInterval})
	}
	return configurations, nil

}
