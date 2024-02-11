package models

import (
	"time"

	"gorm.io/gorm"
)

// Gorm Configuration Model
type Configurations struct {
	gorm.Model
	MonitoredDirectory string        // Directory to be monitored by the watcher
	TimeInterval       time.Duration // Timeinterval of tasks run
	MagicString        string        // Magic String which is searched in the files of the directories
}
