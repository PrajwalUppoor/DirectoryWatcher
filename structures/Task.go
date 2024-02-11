package structures

import (
	"time"
)

type Status int

const (
	Success Status = iota
	Failed
	InProgress
)

// TaskResult represents the result of a task run.
type Task struct {
	TaskId                 uint32
	StartTime              time.Time
	EndTime                time.Time
	TotalRuntime           time.Duration
	FilesAdded             []string
	FilesDeleted           []string
	Configuration          Configuration
	MagicStringOccurrences int
	Status                 Status // "success", "failed", or "in progress"
}

// Task Request Dto
type TaskDto struct {
	ConfigurationId int `json:"configurationId"`
}
