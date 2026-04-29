package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dddd/common"
	"dddd/structs"

	"gopkg.in/yaml.v3"
)

const (
	MaxPOCFileSize  = 5 * 1024 * 1024
	MaxPOCNamesCount = 100
)

func getWorkflowYamlPath() string {
	workflowPath := structs.GlobalConfig.WorkflowYamlPath
	if workflowPath == "" {
		workflowPath = "config/workflow.yaml"
	}
	return workflowPath
}

func readWorkflowYaml() (map[string]map[string]interface{}, error) {
	workflowPath := getWorkflowYamlPath()
	result := make(map[string]map[string]interface{})
	
	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		return result, nil
	}
	
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return nil, err
	}
	
	var fps map[string]interface{}
	err = yaml.Unmarshal(data, &fps)
	if err != nil {
		return nil, err
	}
	
	for productName, workflowInterface := range fps {
		workflowMap, ok := workflowInterface.(map[string]interface{})
		if !ok {
			continue
		}
		
		workflow := make(map[string]interface{})
		workflow["name"] = productName
		
		if types, ok := workflowMap["type"]; ok {
			var typeList []string
			if typeSlice, ok := types.([]interface{}); ok {
				for _, t := range typeSlice {
					if typeStr, ok := t.(string); ok {
						typeList = append(typeList, typeStr)
					}
				}
			}
			workflow["type"] = typeList
		}
		
		if pocs, ok := workflowMap["pocs"]; ok {
			var pocList []string
			if pocSlice, ok := pocs.([]interface{}); ok {
				for _, p := range pocSlice {
					if pocStr, ok := p.(string); ok {
						pocList = append(pocList, pocStr)
					}
				}
			}
			workflow["pocs"] = pocList
		}
		
		result[productName] = workflow
	}
	
	return result, nil
}

func writeWorkflowYaml(data map[string]map[string]interface{}) error {
	workflowPath := getWorkflowYamlPath()
	
	dir := filepath.Dir(workflowPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	
	yamlData := make(map[string]interface{})
	for name, workflow := range data {
		workflowEntry := make(map[string]interface{})
		if types, ok := workflow["type"]; ok {
			workflowEntry["type"] = types
		}
		if pocs, ok := workflow["pocs"]; ok {
			workflowEntry["pocs"] = pocs
		}
		yamlData[name] = workflowEntry
	}
	
	dataBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		return err
	}
	
	return os.WriteFile(workflowPath, dataBytes, 0644)
}

func UploadPOC(filePath string) (string, error) {
	cleanPath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanPath, ".") && !filepath.IsAbs(cleanPath) {
		cleanPath = "./" + cleanPath
	}

	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat POC file: %v", err)
	}

	if fileInfo.Size() > MaxPOCFileSize {
		return "", fmt.Errorf("POC file too large, maximum size is %d bytes", MaxPOCFileSize)
	}

	ext := strings.ToLower(filepath.Ext(cleanPath))
	if ext != ".yaml" && ext != ".yml" {
		return "", fmt.Errorf("POC file must have .yaml or .yml extension")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to read POC file: %v", err)
	}

	var poc map[string]interface{}
	err = yaml.Unmarshal(data, &poc)
	if err != nil {
		return "", fmt.Errorf("invalid YAML format: %v", err)
	}

	if _, ok := poc["id"]; !ok {
		return "", fmt.Errorf("POC must have an 'id' field")
	}
	if _, ok := poc["info"]; !ok {
		return "", fmt.Errorf("POC must have an 'info' field")
	}

	pocDir := structs.GlobalConfig.NucleiTemplate
	if pocDir == "" {
		pocDir = "config/pocs"
	}
	if _, err := os.Stat(pocDir); os.IsNotExist(err) {
		os.MkdirAll(pocDir, 0755)
	}
	fileName := filepath.Base(cleanPath)
	destPath := filepath.Join(pocDir, fileName)
	err = os.WriteFile(destPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write POC file: %v", err)
	}

	return fmt.Sprintf("POC uploaded successfully to %s", destPath), nil
}

func ValidatePOC(pocName string) (string, error) {
	if pocName == "" {
		return "", fmt.Errorf("POC name cannot be empty")
	}

	pocDir := structs.GlobalConfig.NucleiTemplate
	if pocDir == "" {
		pocDir = "config/pocs"
	}

	var pocPath string
	if filepath.IsAbs(pocName) || strings.HasPrefix(pocName, ".") {
		pocPath = pocName
	} else {
		potentialPath := filepath.Join(pocDir, pocName)
		if strings.HasSuffix(pocName, ".yaml") || strings.HasSuffix(pocName, ".yml") {
			pocPath = potentialPath
		} else {
			for _, ext := range []string{".yaml", ".yml"} {
				potentialPathWithExt := potentialPath + ext
				if _, err := os.Stat(potentialPathWithExt); err == nil {
					pocPath = potentialPathWithExt
					break
				}
			}
		}
	}

	if pocPath == "" {
		return "", fmt.Errorf("could not find POC: %s", pocName)
	}

	if _, err := os.Stat(pocPath); os.IsNotExist(err) {
		return "", fmt.Errorf("POC file not found: %s", pocPath)
	}

	data, err := os.ReadFile(pocPath)
	if err != nil {
		return "", fmt.Errorf("failed to read POC file: %v", err)
	}

	var poc map[string]interface{}
	err = yaml.Unmarshal(data, &poc)
	if err != nil {
		return "", fmt.Errorf("invalid YAML format for POC: %v", err)
	}

	if _, ok := poc["id"]; !ok {
		return "", fmt.Errorf("POC missing required 'id' field")
	}
	if _, ok := poc["info"]; !ok {
		return "", fmt.Errorf("POC missing required 'info' field")
	}

	return fmt.Sprintf("POC %s format validation successful", pocName), nil
}

func TestPOC(pocName string, target string) (string, error) {
	if target == "" {
		return ValidatePOC(pocName)
	}

	validationResult, err := ValidatePOC(pocName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s. Would run on target %s using Nuclei", validationResult, target), nil
}

func UpdateWorkflow(fingerprintName string, pocs []string, types []string) (string, error) {
	if fingerprintName == "" {
		return "", fmt.Errorf("fingerprint name cannot be empty")
	}

	if len(fingerprintName) > 100 {
		return "", fmt.Errorf("fingerprint name too long, maximum length is 100 characters")
	}

	if len(pocs) > MaxPOCNamesCount {
		return "", fmt.Errorf("too many POCs, maximum is %d", MaxPOCNamesCount)
	}

	for _, pocName := range pocs {
		if pocName == "" {
			return "", fmt.Errorf("POC name cannot be empty")
		}
		if len(pocName) > 255 {
			return "", fmt.Errorf("POC name too long, maximum length is 255 characters")
		}
	}
	
	if len(types) == 0 {
		types = []string{"root"}
	}

	workflowData, err := readWorkflowYaml()
	if err != nil {
		return "", err
	}

	workflowData[fingerprintName] = map[string]interface{}{
		"type": types,
		"pocs": pocs,
	}

	err = writeWorkflowYaml(workflowData)
	if err != nil {
		return "", err
	}

	structs.WorkFlowDB = nil
	common.ReadWorkFlowDB()

	return fmt.Sprintf("Workflow updated successfully for %s, %d POCs added", fingerprintName, len(pocs)), nil
}

func GetWorkflowDetails(name string) (map[string]interface{}, error) {
	workflowData, err := readWorkflowYaml()
	if err != nil {
		return nil, err
	}
	
	if workflow, ok := workflowData[name]; ok {
		return workflow, nil
	}
	
	if wf, ok := structs.WorkFlowDB[name]; ok {
		result := make(map[string]interface{})
		result["name"] = name
		
		var types []string
		if wf.RootType {
			types = append(types, "root")
		}
		if wf.DirType {
			types = append(types, "dir")
		}
		if wf.BaseType {
			types = append(types, "base")
		}
		if len(types) == 0 {
			types = []string{"root"}
		}
		result["type"] = types
		result["pocs"] = wf.PocsName
		return result, nil
	}
	
	return nil, fmt.Errorf("workflow not found")
}

func GetAllWorkflowDetails() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	
	workflowData, err := readWorkflowYaml()
	if err != nil {
		return nil, err
	}
	
	nameMap := make(map[string]bool)
	for name := range workflowData {
		nameMap[name] = true
	}
	for name := range structs.WorkFlowDB {
		nameMap[name] = true
	}
	
	for name := range nameMap {
		details, _ := GetWorkflowDetails(name)
		if details != nil {
			result = append(result, details)
		}
	}
	
	return result, nil
}

func DeleteWorkflow(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("workflow name cannot be empty")
	}
	
	workflowData, err := readWorkflowYaml()
	if err != nil {
		return "", err
	}
	
	delete(workflowData, name)
	
	err = writeWorkflowYaml(workflowData)
	if err != nil {
		return "", err
	}
	
	structs.WorkFlowDB = nil
	common.ReadWorkFlowDB()
	
	return fmt.Sprintf("Workflow '%s' deleted successfully", name), nil
}

func GetPOCs() ([]string, error) {
	var pocNames []string
	pocDir := structs.GlobalConfig.NucleiTemplate
	if pocDir == "" {
		pocDir = "config/pocs"
	}

	if _, err := os.Stat(pocDir); os.IsNotExist(err) {
		return pocNames, nil
	}

	err := filepath.Walk(pocDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			pocNames = append(pocNames, filepath.Base(path))
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan POC directory: %v", err)
	}

	return pocNames, nil
}

func DeletePOC(name string) (string, error) {
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return "", fmt.Errorf("invalid POC name: path traversal not allowed")
	}

	pocDir := structs.GlobalConfig.NucleiTemplate
	if pocDir == "" {
		pocDir = "config/pocs"
	}

	pocPath := filepath.Join(pocDir, name)
	if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
		for _, ext := range []string{".yaml", ".yml"} {
			if _, err := os.Stat(pocPath + ext); err == nil {
				pocPath += ext
				break
			}
		}
	}

	absPocDir, err := filepath.Abs(pocDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve POC directory: %v", err)
	}
	absPocPath, err := filepath.Abs(pocPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve POC path: %v", err)
	}

	if !strings.HasPrefix(absPocPath, absPocDir) {
		return "", fmt.Errorf("invalid POC path: must be within POC directory")
	}

	if _, err := os.Stat(pocPath); os.IsNotExist(err) {
		return "", fmt.Errorf("POC file not found: %s", name)
	}

	err = os.Remove(pocPath)
	if err != nil {
		return "", fmt.Errorf("failed to delete POC file: %v", err)
	}

	return fmt.Sprintf("POC %s deleted successfully", name), nil
}

func GetWorkflows() (map[string][]string, error) {
	result := make(map[string][]string)
	for name, wf := range structs.WorkFlowDB {
		result[name] = wf.PocsName
	}
	return result, nil
}
