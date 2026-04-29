package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"dddd/structs"
)

const (
	storageDir  = ".dddd"
	tasksFile   = "tasks.json"
	resultsFile = "results.json"
	templatesFile = "templates.json"
)

var (
	storageMutex  sync.Mutex
	stopAutoSave  chan struct{}
	autoSaveDone  chan struct{}
)

func getStoragePath(filename string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", storageDir, filename)
	}
	return filepath.Join(homeDir, storageDir, filename)
}

func initStorageDir() error {
	storagePath := getStoragePath("")
	return os.MkdirAll(storagePath, 0755)
}

func saveToFile(filename string, data interface{}) error {
	storageMutex.Lock()
	defer storageMutex.Unlock()

	if err := initStorageDir(); err != nil {
		return err
	}

	filePath := getStoragePath(filename)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

func loadFromFile(filename string, data interface{}) error {
	storageMutex.Lock()
	defer storageMutex.Unlock()

	filePath := getStoragePath(filename)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	return json.Unmarshal(jsonData, data)
}

func SaveTasks() error {
	type taskData struct {
		Status       map[string]string              `json:"status"`
		Results      map[string]map[string]interface{} `json:"results"`
		Templates    map[string]map[string]interface{} `json:"templates"`
		CreationTime map[string]string              `json:"creationTime"`
	}

	creationTimeStr := make(map[string]string)
	for taskID, t := range taskCreationTime {
		creationTimeStr[taskID] = t.Format(time.RFC3339)
	}

	data := taskData{
		Status:       taskStatus,
		Results:      taskResults,
		Templates:    taskTemplates,
		CreationTime: creationTimeStr,
	}

	return saveToFile(tasksFile, data)
}

func LoadTasks() error {
	type taskData struct {
		Status       map[string]string              `json:"status"`
		Results      map[string]map[string]interface{} `json:"results"`
		Templates    map[string]map[string]interface{} `json:"templates"`
		CreationTime map[string]string              `json:"creationTime"`
	}

	var data taskData
	if err := loadFromFile(tasksFile, &data); err != nil {
		return err
	}

	if data.Status != nil {
		taskStatus = data.Status
	}
	if data.Results != nil {
		taskResults = data.Results
	}
	if data.Templates != nil {
		taskTemplates = data.Templates
	}
	if data.CreationTime != nil {
		for taskID, tStr := range data.CreationTime {
			t, err := time.Parse(time.RFC3339, tStr)
			if err == nil {
				taskCreationTime[taskID] = t
			}
		}
	}

	return nil
}

func SaveConfigToFile() error {
	return saveToFile("config.json", structs.GlobalConfig)
}

func LoadConfig() error {
	return loadFromFile("config.json", &structs.GlobalConfig)
}

func init() {
	initStorageDir()
	LoadTasks()
	LoadConfig()
}

func AutoSaveTasks() {
	stopAutoSave = make(chan struct{})
	autoSaveDone = make(chan struct{})
	
	go func() {
		defer close(autoSaveDone)
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				SaveTasks()
			case <-stopAutoSave:
				return
			}
		}
	}()
}

func StopAutoSave() {
	if stopAutoSave != nil {
		close(stopAutoSave)
		<-autoSaveDone
	}
}