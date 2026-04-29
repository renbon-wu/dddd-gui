package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dddd/common"
	"dddd/lib/ddfinger"
	"dddd/structs"

	"gopkg.in/yaml.v3"
)

func getFingerYamlPath() string {
	fingerPath := structs.GlobalConfig.FingerConfigFilePath
	if fingerPath == "" {
		fingerPath = "config/finger.yaml"
	}
	return fingerPath
}

func readFingerYaml() (map[string][]string, error) {
	fingerPath := getFingerYamlPath()
	result := make(map[string][]string)
	
	if _, err := os.Stat(fingerPath); os.IsNotExist(err) {
		return result, nil
	}
	
	data, err := os.ReadFile(fingerPath)
	if err != nil {
		return nil, err
	}
	
	var fps map[string]interface{}
	err = yaml.Unmarshal(data, &fps)
	if err != nil {
		return nil, err
	}
	
	for productName, rulesInterface := range fps {
		var rules []string
		rulesList, ok := rulesInterface.([]interface{})
		if !ok {
			continue
		}
		for _, ruleInterface := range rulesList {
			rule, ok := ruleInterface.(string)
			if ok {
				rules = append(rules, rule)
			}
		}
		result[productName] = rules
	}
	
	return result, nil
}

func writeFingerYaml(data map[string][]string) error {
	fingerPath := getFingerYamlPath()
	
	dir := filepath.Dir(fingerPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	
	return os.WriteFile(fingerPath, yamlData, 0644)
}

func UploadFingerprint(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read fingerprint file: %v", err)
	}

	fps := make(map[string]interface{})
	err = yaml.Unmarshal(data, &fps)
	if err != nil {
		return "", fmt.Errorf("invalid YAML format: %v", err)
	}

	for productName, rulesInterface := range fps {
		rulesList, ok := rulesInterface.([]interface{})
		if !ok {
			return "", fmt.Errorf("invalid format for fingerprint: %s, expected array of rules", productName)
		}
		for _, ruleInterface := range rulesList {
			rule, ok := ruleInterface.(string)
			if !ok {
				return "", fmt.Errorf("invalid rule format for fingerprint: %s, expected string", productName)
			}
			ruleData := ddfinger.ParseRule(rule)
			if len(ruleData) == 0 {
				return "", fmt.Errorf("invalid rule syntax for fingerprint: %s", productName)
			}
		}
	}

	configPath := structs.GlobalConfig.FingerConfigFilePath
	if configPath == "" {
		configPath = "config/finger.yaml"
	}
	if !strings.Contains(configPath, "/") {
		configPath = "config/" + configPath
	}
	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	err = os.WriteFile(configPath, data, 0666)
	if err != nil {
		return "", fmt.Errorf("failed to write fingerprint file: %v", err)
	}

	structs.FingerprintDB = nil
	common.ParseFingerDB()

	return fmt.Sprintf("Fingerprint uploaded successfully, %d fingerprints loaded", len(structs.FingerprintDB)), nil
}

func ValidateFingerprintByName(name string) (string, error) {
	details, err := GetFingerprintDetails(name)
	if err != nil {
		return "", err
	}

	rules, ok := details["rules"].([]string)
	if !ok || len(rules) == 0 {
		return "", fmt.Errorf("fingerprint has no valid rules")
	}

	for _, rule := range rules {
		ruleData := ddfinger.ParseRule(rule)
		if len(ruleData) == 0 {
			return "", fmt.Errorf("invalid rule syntax in fingerprint: %s", rule)
		}
	}

	return fmt.Sprintf("Fingerprint '%s' format validation passed", name), nil
}

func ValidateFingerprint(rule string) (string, error) {
	ruleData := ddfinger.ParseRule(rule)
	if len(ruleData) == 0 {
		return "", fmt.Errorf("invalid fingerprint rule syntax")
	}

	return "Fingerprint format validation passed", nil
}

func TestFingerprint(nameOrRule string, target string) (string, error) {
	if target == "" {
		_, err := GetFingerprintDetails(nameOrRule)
		if err == nil {
			return ValidateFingerprintByName(nameOrRule)
		}
		return ValidateFingerprint(nameOrRule)
	}

	details, err := GetFingerprintDetails(nameOrRule)
	if err == nil {
		rules, ok := details["rules"].([]string)
		if ok && len(rules) > 0 {
			for _, rule := range rules {
				ruleData := ddfinger.ParseRule(rule)
				if len(ruleData) == 0 {
					return "", fmt.Errorf("invalid rule syntax in fingerprint")
				}
			}
		}
	} else {
		ruleData := ddfinger.ParseRule(nameOrRule)
		if len(ruleData) == 0 {
			return "", fmt.Errorf("invalid fingerprint rule syntax")
		}
	}

	return fmt.Sprintf("Fingerprint test would run on %s. This functionality requires actual HTTP request implementation", target), nil
}

func GetFingerprints() ([]string, error) {
	var fingerNames []string
	nameMap := make(map[string]bool)
	for _, finger := range structs.FingerprintDB {
		if !nameMap[finger.ProductName] {
			fingerNames = append(fingerNames, finger.ProductName)
			nameMap[finger.ProductName] = true
		}
	}
	return fingerNames, nil
}

func GetFingerprintDetails(name string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	fingerData, err := readFingerYaml()
	if err != nil {
		return nil, err
	}
	
	if rules, ok := fingerData[name]; ok {
		result["name"] = name
		result["rules"] = rules
	} else {
		var rules []string
		for _, finger := range structs.FingerprintDB {
			if finger.ProductName == name {
				rules = append(rules, finger.AllString)
			}
		}
		result["name"] = name
		result["rules"] = rules
	}
	
	return result, nil
}

func GetAllFingerprintDetails() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	
	fingerData, err := readFingerYaml()
	if err != nil {
		return nil, err
	}
	
	nameMap := make(map[string]bool)
	for name := range fingerData {
		nameMap[name] = true
	}
	for _, finger := range structs.FingerprintDB {
		nameMap[finger.ProductName] = true
	}
	
	for name := range nameMap {
		details, _ := GetFingerprintDetails(name)
		result = append(result, details)
	}
	
	return result, nil
}

func AddFingerprint(name string, rules []string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("fingerprint name cannot be empty")
	}
	
	for _, rule := range rules {
		ruleData := ddfinger.ParseRule(rule)
		if len(ruleData) == 0 {
			return "", fmt.Errorf("invalid rule syntax: %s", rule)
		}
	}
	
	fingerData, err := readFingerYaml()
	if err != nil {
		return "", err
	}
	
	fingerData[name] = rules
	
	err = writeFingerYaml(fingerData)
	if err != nil {
		return "", err
	}
	
	structs.FingerprintDB = nil
	common.ParseFingerDB()
	
	return fmt.Sprintf("Fingerprint '%s' added successfully", name), nil
}

func EditFingerprint(oldName string, newName string, rules []string) (string, error) {
	if oldName == "" {
		return "", fmt.Errorf("old fingerprint name cannot be empty")
	}
	
	for _, rule := range rules {
		ruleData := ddfinger.ParseRule(rule)
		if len(ruleData) == 0 {
			return "", fmt.Errorf("invalid rule syntax: %s", rule)
		}
	}
	
	fingerData, err := readFingerYaml()
	if err != nil {
		return "", err
	}
	
	if _, ok := fingerData[oldName]; ok {
		delete(fingerData, oldName)
	}
	
	if newName != "" {
		fingerData[newName] = rules
	}
	
	err = writeFingerYaml(fingerData)
	if err != nil {
		return "", err
	}
	
	structs.FingerprintDB = nil
	common.ParseFingerDB()
	
	return fmt.Sprintf("Fingerprint updated successfully"), nil
}

func DeleteFingerprint(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("fingerprint name cannot be empty")
	}
	
	fingerData, err := readFingerYaml()
	if err != nil {
		return "", err
	}
	
	delete(fingerData, name)
	
	err = writeFingerYaml(fingerData)
	if err != nil {
		return "", err
	}
	
	structs.FingerprintDB = nil
	common.ParseFingerDB()
	
	return fmt.Sprintf("Fingerprint '%s' deleted successfully", name), nil
}