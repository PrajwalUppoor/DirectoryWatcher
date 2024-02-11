package services

import (
	"bufio"
	"dirwatcher/db"
	"dirwatcher/models"
	"dirwatcher/structures"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	Tasks      = make(map[uint32]structures.Task) // Map to store the tasks in memory
	TasksMutex sync.RWMutex                       // Mutex for concurrent map access

)

/*
*

	This method will create a connection object to postgres Sql
	@param config Configuration to be used for the perticular task
	@param stopChan Channel to stop that perticular task
	@return error
*/
func DirectoryWatcher(config structures.Configuration, stopChan chan struct{}) {

	task := structures.Task{
		StartTime:     time.Now(),
		Configuration: config,
		Status:        structures.InProgress,
	}
	log.Printf("[DirectoryWatcher] task with taskId:%v \n", task)
	modelConfigurations := models.Configurations{MonitoredDirectory: task.Configuration.MonitoredDirectory, TimeInterval: task.Configuration.TimeInterval, MagicString: task.Configuration.MagicString}
	modelTask := models.Task{StartTime: task.StartTime, Configurations: modelConfigurations, Status: models.Status(task.Status)}

	if err := db.DB.Create(&modelTask).Error; err != nil {
		log.Printf("[DirectoryWatcher] Error occured while creating task : %v", err)
		return
	}
	taskId := uint32(modelTask.ID)
	task.TaskId = taskId

	fileWithCount := make(map[string]int)
	processDirectory(config, fileWithCount)
	task.MagicStringOccurrences = sumOfValuesOfMap(fileWithCount)

	log.Printf("[DirectoryWatcher]filecount is %v", fileWithCount)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("[DirectoryWatcher] Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(config.MonitoredDirectory)
	if err != nil {
		log.Println("[DirectoryWatcher] Error adding directory to watcher:", err)
		return
	}
	log.Printf("[DirectoryWatcher] Watching directory: %s\n", config.MonitoredDirectory)
	ticker := time.NewTicker(task.Configuration.TimeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			TasksMutex.RLock()
			val, ok := Tasks[taskId]
			TasksMutex.RUnlock()
			if !ok || (ok && !reflect.DeepEqual(val, task)) {
				TasksMutex.Lock()
				Tasks[taskId] = task
				TasksMutex.Unlock()
				log.Printf("[DirectoryWatcher] Saving it to local map :%v", task)
			}
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				log.Printf("[DirectoryWatcher]  File %s has added \n", event.Name)
				task.FilesAdded = append(task.FilesAdded, event.Name)
				fileWithCount[event.Name] = 0
			}
			if event.Has(fsnotify.Remove) {
				log.Printf("[DirectoryWatcher]  File %s has deleted \n", event.Name)
				task.FilesDeleted = append(task.FilesDeleted, event.Name)
				task.MagicStringOccurrences -= fileWithCount[event.Name]
				delete(fileWithCount, event.Name)
			}
			if event.Has(fsnotify.Write) {
				fileName := event.Name
				log.Printf("[DirectoryWatcher]  File %s has increased magic strings \n", fileName)
				fileWithCount[fileName] = countMagicString(fileName, config)
				task.MagicStringOccurrences = sumOfValuesOfMap(fileWithCount)
			}

		case <-stopChan:
			log.Printf("[DirectoryWatcher] Stopped watching directory: %s\n", config.MonitoredDirectory)
			TasksMutex.RLock()
			task = Tasks[taskId]
			TasksMutex.RUnlock()
			task.EndTime = time.Now()
			task.TotalRuntime = task.EndTime.Sub(task.StartTime)
			task.Status = structures.Success
			log.Printf("[DirectoryWatcher] Task is completed with : %v", task)

			modelTaskValue := models.Task{StartTime: task.StartTime, EndTime: task.EndTime, TotalRuntime: task.TotalRuntime, FilesAdded: strings.Join(task.FilesAdded, ","), FilesDeleted: strings.Join(task.FilesDeleted, ","), MagicStringOccurrences: task.MagicStringOccurrences, Configurations: modelConfigurations, Status: models.Status(task.Status), ConfigurationsId: int(modelConfigurations.ID)}
			if err := db.DB.Save(&modelTaskValue).Error; err != nil {
				log.Printf("[DirectoryWatcher] Error occured while saving task details : %v", err)
			}
			return
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Printf("[DirectoryWatcher] Stopped watching directory %s due to error: %v \n", config.MonitoredDirectory, err)
				task.EndTime = time.Now()
				task.TotalRuntime = task.EndTime.Sub(task.StartTime)
				task.Status = structures.Failed
				log.Printf("[DirectoryWatcher] Task has failed with: %v", task)
				TasksMutex.RLock()
				log.Printf("[DirectoryWatcher] Task map values:%v", Tasks)
				TasksMutex.RUnlock()
				modelTaskValue := models.Task{StartTime: task.StartTime, EndTime: task.EndTime, TotalRuntime: task.TotalRuntime, FilesAdded: strings.Join(task.FilesAdded, ","), FilesDeleted: strings.Join(task.FilesDeleted, ","), MagicStringOccurrences: task.MagicStringOccurrences, Configurations: modelConfigurations, Status: models.Status(task.Status), ConfigurationsId: int(modelConfigurations.ID)}
				if err := db.DB.Save(&modelTaskValue).Error; err != nil {
					log.Printf("[DirectoryWatcher] Error occured while saving task details : %v", err)
					return
				}
			}
		}

	}

}

/*
*

	This method is used to count magic string occurrences in a file
	@param fileName name of the file where the magic strings are counted
	@param config Configuration to be used for the perticular task
	@return count of the magic string occurence
*/
func countMagicString(filename string, config structures.Configuration) int {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("[countMagicString] Error opening file:", err)

	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	// Iterate through each line
	for scanner.Scan() {
		line := scanner.Text()
		// Read file content and count occurrences
		count += strings.Count(line, config.MagicString)
	}

	if err := scanner.Err(); err != nil {
		log.Println("[countMagicString] Error reading file:", err)
	}
	// Return the count
	return count
}

/*
*

	This method is used to walk through the files of the directory and update the map with filename along with corresponding magicstring count
	@param fileWithCount Map of file with corresponding count of magic strings
	@param config Configuration to be used for the perticular task
	@return void
*/
func processDirectory(config structures.Configuration, fileWithCount map[string]int) {
	existingFiles := make(map[string]bool)
	var mu sync.Mutex // Mutex for concurrent map access
	// Walk through the directory
	err := filepath.Walk(config.MonitoredDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file exists in the map
		mu.Lock()
		_, exists := existingFiles[path]
		mu.Unlock()

		if !exists {
			// New file added
			fileWithCount[path] = 0
			mu.Lock()
			existingFiles[path] = true
			mu.Unlock()
		}

		// Count magic string occurrences
		count := countMagicString(path, config)
		fileWithCount[path] = count
		log.Printf("[processDirectory] File: %s, Magic String Count: %d\n", path, count)
		return nil
	})

	if err != nil {
		log.Println("[processDirectory] Error walking directory:", err)
	}
}

/*
*

	This method is used to get the sum of values of the given map of string with integer values
	@param valueMap
	@return sum of the values of the map
*/
func sumOfValuesOfMap(valueMap map[string]int) int {
	sum := 0
	for _, val := range valueMap {
		sum += val
	}
	return sum
}
