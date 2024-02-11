package models

import (
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	Success Status = iota
	Failed
	InProgress
)

// Gorm Task Model
type Task struct {
	gorm.Model
	StartTime              time.Time      //Start Time of the task
	EndTime                time.Time      // End  Time of the task
	TotalRuntime           time.Duration  // Total Runtime of the task
	FilesAdded             string         // Files added to directory
	FilesDeleted           string         // Files deleted from directory
	MagicStringOccurrences int            // Total Number of Occurence of Magic String
	Status                 Status         // Enum for Status of the Task success,failed,inProgress
	ConfigurationsId       int            // ConfigurationId
	Configurations         Configurations `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //Configuration
}
