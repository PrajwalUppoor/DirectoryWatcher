package structures

import (
	"time"
)

// Configuration holds the task configuration.
type Configuration struct {
	MonitoredDirectory string        `json:"fileDirectory"`
	TimeInterval       time.Duration `json:"timeInterval"`
	MagicString        string        `json:"magicString"`
}

// Configuratiun Response Dto
type ConfigurationResponse struct {
	ConfigurationId    uint          `json:"configurationId"`
	MonitoredDirectory string        `json:"fileDirectory"`
	TimeInterval       time.Duration `json:"timeInterval"`
	MagicString        string        `json:"magicString"`
}
