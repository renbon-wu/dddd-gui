package api

import (
	"fmt"
	"os"
	"path/filepath"

	"dddd/structs"

	"gopkg.in/yaml.v3"
)

const configFilePath = "config/config.yaml"

// ConfigData 配置数据结构
type ConfigData struct {
	APIKeys      map[string]string `yaml:"api_keys"`
	UIData       map[string]string `yaml:"ui"`
	Dictionaries map[string]string `yaml:"dictionaries"`
}

// GetConfig returns the current configuration
func GetConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"proxy":          structs.GlobalConfig.HTTPProxy,
		"webThreads":     structs.GlobalConfig.WebThreads,
		"webTimeout":     structs.GlobalConfig.WebTimeout,
		"fingerConfig":   structs.GlobalConfig.FingerConfigFilePath,
		"workflowYaml":   structs.GlobalConfig.WorkflowYamlPath,
		"nucleiTemplate": structs.GlobalConfig.NucleiTemplate,
		"subdomainDict":  structs.GlobalConfig.SubdomainWordListFile,
		"Hunter":         structs.GlobalConfig.Hunter,
		"Fofa":           structs.GlobalConfig.Fofa,
		"Quake":          structs.GlobalConfig.Quake,
	}, nil
}

// SaveConfig saves the configuration
func SaveConfig(config map[string]interface{}) (string, error) {
	if proxy, ok := config["proxy"].(string); ok {
		structs.GlobalConfig.HTTPProxy = proxy
	}
	if webThreads, ok := config["webThreads"].(float64); ok {
		structs.GlobalConfig.WebThreads = int(webThreads)
	}
	if webTimeout, ok := config["webTimeout"].(float64); ok {
		structs.GlobalConfig.WebTimeout = int(webTimeout)
	}
	if fingerConfig, ok := config["fingerConfig"].(string); ok {
		structs.GlobalConfig.FingerConfigFilePath = fingerConfig
	}
	if workflowYaml, ok := config["workflowYaml"].(string); ok {
		structs.GlobalConfig.WorkflowYamlPath = workflowYaml
	}
	if nucleiTemplate, ok := config["nucleiTemplate"].(string); ok {
		structs.GlobalConfig.NucleiTemplate = nucleiTemplate
	}
	if subdomainDict, ok := config["subdomainDict"].(string); ok {
		structs.GlobalConfig.SubdomainWordListFile = subdomainDict
	}
	if Hunter, ok := config["Hunter"].(bool); ok {
		structs.GlobalConfig.Hunter = Hunter
	}
	if Fofa, ok := config["Fofa"].(bool); ok {
		structs.GlobalConfig.Fofa = Fofa
	}
	if Quake, ok := config["Quake"].(bool); ok {
		structs.GlobalConfig.Quake = Quake
	}

	err := saveConfigToFile()
	if err != nil {
		return "", err
	}

	return "Configuration saved successfully", nil
}

// GetAPIs returns the API configurations
func GetAPIs() (map[string]string, error) {
	config, err := loadConfigFromFile()
	if err != nil {
		return map[string]string{
			"Hunter": "",
			"Fofa":   "",
			"Quake":  "",
		}, nil
	}
	return config.APIKeys, nil
}

// SaveAPI saves an API configuration
func SaveAPI(apiName string, apiKey string) (string, error) {
	config, err := loadConfigFromFile()
	if err != nil {
		return "", err
	}

	if config.APIKeys == nil {
		config.APIKeys = make(map[string]string)
	}
	config.APIKeys[apiName] = apiKey

	err = saveConfigData(config)
	if err != nil {
		return "", err
	}

	return apiName + " API key saved successfully", nil
}

// SaveApiConfig saves the API configuration
func SaveApiConfig(apiConfig map[string]interface{}) (string, error) {
	config, err := loadConfigFromFile()
	if err != nil {
		return "", err
	}

	if config.APIKeys == nil {
		config.APIKeys = make(map[string]string)
	}

	if fofaKey, ok := apiConfig["fofaKey"].(string); ok {
		config.APIKeys["Fofa"] = fofaKey
	}
	if shodanKey, ok := apiConfig["shodanKey"].(string); ok {
		config.APIKeys["Shodan"] = shodanKey
	}
	if censysID, ok := apiConfig["censysID"].(string); ok {
		config.APIKeys["CensysID"] = censysID
	}
	if censysSecret, ok := apiConfig["censysSecret"].(string); ok {
		config.APIKeys["CensysSecret"] = censysSecret
	}

	err = saveConfigData(config)
	if err != nil {
		return "", err
	}

	return "API configuration saved successfully", nil
}

// SaveUiConfig saves the UI configuration
func SaveUiConfig(uiConfig map[string]interface{}) (string, error) {
	config, err := loadConfigFromFile()
	if err != nil {
		return "", err
	}

	if config.UIData == nil {
		config.UIData = make(map[string]string)
	}

	if theme, ok := uiConfig["theme"].(string); ok {
		config.UIData["theme"] = theme
	}
	if language, ok := uiConfig["language"].(string); ok {
		config.UIData["language"] = language
	}
	if fontSize, ok := uiConfig["fontSize"].(string); ok {
		config.UIData["fontSize"] = fontSize
	}

	err = saveConfigData(config)
	if err != nil {
		return "", err
	}

	return "UI configuration saved successfully", nil
}

// GetDictionaries returns the dictionary paths
func GetDictionaries() (map[string]string, error) {
	config, err := loadConfigFromFile()
	if err != nil {
		return map[string]string{
			"subdomain": structs.GlobalConfig.SubdomainWordListFile,
		}, nil
	}
	return config.Dictionaries, nil
}

// SaveDictionary saves a dictionary path
func SaveDictionary(dictType string, dictPath string) (string, error) {
	structs.GlobalConfig.SubdomainWordListFile = dictPath

	config, err := loadConfigFromFile()
	if err != nil {
		return "", err
	}

	if config.Dictionaries == nil {
		config.Dictionaries = make(map[string]string)
	}
	config.Dictionaries[dictType] = dictPath

	err = saveConfigData(config)
	if err != nil {
		return "", err
	}

	return dictType + " dictionary path saved successfully", nil
}

// AddDictionary adds a new dictionary
func AddDictionary(name string, path string) (string, error) {
	if name == "" || path == "" {
		return "", fmt.Errorf("dictionary name and path cannot be empty")
	}

	config, err := loadConfigFromFile()
	if err != nil {
		return "", err
	}

	if config.Dictionaries == nil {
		config.Dictionaries = make(map[string]string)
	}
	config.Dictionaries[name] = path

	err = saveConfigData(config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Dictionary '%s' added successfully", name), nil
}

// DeleteDictionary deletes a dictionary
func DeleteDictionary(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("dictionary name cannot be empty")
	}

	config, err := loadConfigFromFile()
	if err != nil {
		return "", err
	}

	if config.Dictionaries != nil {
		delete(config.Dictionaries, name)
	}

	err = saveConfigData(config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Dictionary '%s' deleted successfully", name), nil
}

// loadConfigFromFile loads configuration from file
func loadConfigFromFile() (*ConfigData, error) {
	config := &ConfigData{
		APIKeys:      make(map[string]string),
		UIData:       make(map[string]string),
		Dictionaries: make(map[string]string),
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// saveConfigData saves configuration data to file
func saveConfigData(config *ConfigData) error {
	dir := filepath.Dir(configFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, data, 0644)
}

// saveConfigToFile saves global config to file
func saveConfigToFile() error {
	config, err := loadConfigFromFile()
	if err != nil {
		return err
	}

	config.Dictionaries["subdomain"] = structs.GlobalConfig.SubdomainWordListFile

	return saveConfigData(config)
}

// LoadDictionaries loads dictionary configurations from file to global config
func LoadDictionaries() {
	config, err := loadConfigFromFile()
	if err != nil {
		return
	}

	if config.Dictionaries != nil {
		if path, ok := config.Dictionaries["subdomain"]; ok && path != "" {
			structs.GlobalConfig.SubdomainWordListFile = path
		}
	}

	if config.APIKeys != nil {
		if _, ok := config.APIKeys["Hunter"]; ok {
			structs.GlobalConfig.Hunter = true
		}
		if _, ok := config.APIKeys["Fofa"]; ok {
			structs.GlobalConfig.Fofa = true
		}
		if _, ok := config.APIKeys["Quake"]; ok {
			structs.GlobalConfig.Quake = true
		}
	}
}