package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"

	"dddd/internal/api"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed config/default/*
var defaultConfig embed.FS

func init() {
	configDir := "config"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
		files, _ := defaultConfig.ReadDir("config/default")
		for _, file := range files {
			data, _ := defaultConfig.ReadFile("config/default/" + file.Name())
			os.WriteFile(filepath.Join(configDir, file.Name()), data, 0644)
		}
	}
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "VulScanX",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	api.InitTaskCleanup()
	api.AutoSaveTasks()
	api.LoadDictionaries()
}

func (a *App) RunScan(targets []string, options map[string]interface{}) (string, error) {
	return api.RunScan(targets, options)
}

func (a *App) RunScanWithID(taskID string, targets []string, options map[string]interface{}) (string, error) {
	return api.RunScanWithID(taskID, targets, options)
}

func (a *App) GetTaskStatus(taskID string) (string, error) {
	return api.GetTaskStatus(taskID)
}

func (a *App) GetTasks() (map[string]map[string]interface{}, error) {
	return api.GetTasks()
}

func (a *App) SaveTaskTemplate(name string, config map[string]interface{}) (string, error) {
	return api.SaveTaskTemplate(name, config)
}

func (a *App) LoadTaskTemplate(name string) (map[string]interface{}, error) {
	return api.LoadTaskTemplate(name)
}

func (a *App) GetTaskTemplates() ([]string, error) {
	return api.GetTaskTemplates()
}

func (a *App) DeleteTaskTemplate(name string) (string, error) {
	return api.DeleteTaskTemplate(name)
}

func (a *App) ImportTargets(filePath string) (string, error) {
	return api.ImportTargets(filePath)
}

func (a *App) ImportTargetsFromText(content string) (string, error) {
	return api.ImportTargetsFromText(content)
}

func (a *App) EditTarget(oldTarget string, newTarget string) (string, error) {
	return api.EditTarget(oldTarget, newTarget)
}

func (a *App) GetTargets() (map[string]map[string]interface{}, error) {
	return api.GetTargets()
}

func (a *App) BatchDeleteTargets(targetList []string) (string, error) {
	return api.BatchDeleteTargets(targetList)
}

func (a *App) BatchScanTargets(targetList []string) (string, error) {
	return api.BatchScanTargets(targetList)
}

func (a *App) CreateTargetGroup(groupName string, targetList []string) (string, error) {
	return api.CreateTargetGroup(groupName, targetList)
}

func (a *App) GetTargetGroups() (map[string][]string, error) {
	return api.GetTargetGroups()
}

func (a *App) AddTargetToGroup(groupName string, target string) (string, error) {
	return api.AddTargetToGroup(groupName, target)
}

func (a *App) RemoveTargetFromGroup(groupName string, target string) (string, error) {
	return api.RemoveTargetFromGroup(groupName, target)
}

func (a *App) UploadFingerprint(filePath string) (string, error) {
	return api.UploadFingerprint(filePath)
}

func (a *App) ValidateFingerprint(rule string) (string, error) {
	return api.ValidateFingerprint(rule)
}

func (a *App) ValidateFingerprintByName(name string) (string, error) {
	return api.ValidateFingerprintByName(name)
}

func (a *App) TestFingerprint(rule string, target string) (string, error) {
	return api.TestFingerprint(rule, target)
}

func (a *App) GetFingerprints() ([]string, error) {
	return api.GetFingerprints()
}

func (a *App) DeleteFingerprint(name string) (string, error) {
	return api.DeleteFingerprint(name)
}

func (a *App) UploadPOC(filePath string) (string, error) {
	return api.UploadPOC(filePath)
}

func (a *App) ValidatePOC(pocName string) (string, error) {
	return api.ValidatePOC(pocName)
}

func (a *App) TestPOC(pocName string, target string) (string, error) {
	return api.TestPOC(pocName, target)
}

func (a *App) UpdateWorkflow(fingerprintName string, pocNames []string, types []string) (string, error) {
	return api.UpdateWorkflow(fingerprintName, pocNames, types)
}

func (a *App) GetAllFingerprintDetails() ([]map[string]interface{}, error) {
	return api.GetAllFingerprintDetails()
}

func (a *App) GetFingerprintDetails(name string) (map[string]interface{}, error) {
	return api.GetFingerprintDetails(name)
}

func (a *App) AddFingerprint(name string, rules []string) (string, error) {
	return api.AddFingerprint(name, rules)
}

func (a *App) EditFingerprint(oldName string, newName string, rules []string) (string, error) {
	return api.EditFingerprint(oldName, newName, rules)
}

func (a *App) GetWorkflowDetails(name string) (map[string]interface{}, error) {
	return api.GetWorkflowDetails(name)
}

func (a *App) GetAllWorkflowDetails() ([]map[string]interface{}, error) {
	return api.GetAllWorkflowDetails()
}

func (a *App) DeleteWorkflow(name string) (string, error) {
	return api.DeleteWorkflow(name)
}

func (a *App) GetPOCs() ([]string, error) {
	return api.GetPOCs()
}

func (a *App) DeletePOC(name string) (string, error) {
	return api.DeletePOC(name)
}

func (a *App) GetWorkflows() (map[string][]string, error) {
	return api.GetWorkflows()
}

func (a *App) GetResults(taskID string) (map[string]interface{}, error) {
	return api.GetResults(taskID)
}

func (a *App) GetResultOverview() (map[string]interface{}, error) {
	return api.GetResultOverview()
}

func (a *App) ExportResults(format string, exportRange string) (string, error) {
	return api.ExportResults(format, exportRange)
}

func (a *App) ExportResult(target string, format string) (string, error) {
	return api.ExportResult(target, format)
}

func (a *App) GenerateReport(reportType string, format string, title string) (string, error) {
	return api.GenerateReport(reportType, format, title)
}

func (a *App) GetConfig() (map[string]interface{}, error) {
	return api.GetConfig()
}

func (a *App) SaveConfig(config map[string]interface{}) (string, error) {
	return api.SaveConfig(config)
}

func (a *App) GetAPIs() (map[string]string, error) {
	return api.GetAPIs()
}

func (a *App) SaveAPI(apiName string, apiKey string) (string, error) {
	return api.SaveAPI(apiName, apiKey)
}

func (a *App) GetDictionaries() (map[string]string, error) {
	return api.GetDictionaries()
}

func (a *App) SaveDictionary(dictType string, dictPath string) (string, error) {
	return api.SaveDictionary(dictType, dictPath)
}

func (a *App) SaveApiConfig(apiConfig map[string]interface{}) (string, error) {
	return api.SaveApiConfig(apiConfig)
}

func (a *App) SaveUiConfig(uiConfig map[string]interface{}) (string, error) {
	return api.SaveUiConfig(uiConfig)
}

func (a *App) AddDictionary(name string, path string) (string, error) {
	return api.AddDictionary(name, path)
}

func (a *App) DeleteDictionary(name string) (string, error) {
	return api.DeleteDictionary(name)
}
