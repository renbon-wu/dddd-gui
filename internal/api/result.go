package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// GetResults returns the scan results
func GetResults(taskID string) (map[string]interface{}, error) {
	if taskID == "" {
		// Return all results
		return map[string]interface{}{
			"tasks": taskResults,
		}, nil
	}

	// Return specific task results
	result, exists := taskResults[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	return map[string]interface{}{
		"task":  taskID,
		"result": result,
	}, nil
}

// ExportResults exports results to CSV or JSON format (frontend compatible)
func ExportResults(format string, exportRange string) (string, error) {
	var data interface{}
	var filename string

	// Filter results based on range
	filteredResults := make(map[string]map[string]interface{})
	if exportRange == "all" {
		filteredResults = taskResults
	} else if exportRange == "vulnerabilities" {
		for id, result := range taskResults {
			if res, ok := result["result"].(string); ok && res != "" {
				filteredResults[id] = result
			}
		}
	} else if exportRange == "fingerprints" {
		filteredResults = taskResults
	}

	data = filteredResults
	filename = fmt.Sprintf("dddd-results-all-%s", time.Now().Format("20060102-150405"))

	switch format {
	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %v", err)
		}
		filename += ".json"
		err = os.WriteFile(filename, jsonData, 0666)
		if err != nil {
			return "", fmt.Errorf("failed to write JSON file: %v", err)
		}

	case "csv":
		csvData := "Target,Status,Result,Time\n"
		for _, result := range filteredResults {
			targets, ok := result["targets"].([]string)
			if ok {
				for _, target := range targets {
					csvData += fmt.Sprintf("%s,completed,%s,%s\n", target, result["result"], result["time"])
				}
			}
		}
		filename += ".csv"
		err := os.WriteFile(filename, []byte(csvData), 0666)
		if err != nil {
			return "", fmt.Errorf("failed to write CSV file: %v", err)
		}

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	return fmt.Sprintf("Results exported to %s", filename), nil
}

// GetResultOverview returns an overview of all results
func GetResultOverview() (map[string]interface{}, error) {
	totalTargets := 0
	fingerprintsFound := 0
	vulnerabilitiesFound := 0
	scanTime := ""

	for _, result := range taskResults {
		if targets, ok := result["targets"].([]string); ok {
			totalTargets += len(targets)
		}
		if time, ok := result["time"].(string); ok && time > scanTime {
			scanTime = time
		}
	}

	return map[string]interface{}{
		"totalTargets":        totalTargets,
		"fingerprintsFound":   fingerprintsFound,
		"vulnerabilitiesFound": vulnerabilitiesFound,
		"scanTime":            scanTime,
	}, nil
}

// ExportResult exports a single result to CSV or JSON format
func ExportResult(target string, format string) (string, error) {
	for _, result := range taskResults {
		if targets, ok := result["targets"].([]string); ok {
			for _, t := range targets {
				if t == target {
					var data interface{} = result
					filename := fmt.Sprintf("dddd-result-%s-%s", target, time.Now().Format("20060102-150405"))

					switch format {
					case "json":
						jsonData, err := json.MarshalIndent(data, "", "  ")
						if err != nil {
							return "", fmt.Errorf("failed to marshal JSON: %v", err)
						}
						filename += ".json"
						err = os.WriteFile(filename, jsonData, 0666)
						if err != nil {
							return "", fmt.Errorf("failed to write JSON file: %v", err)
						}

					case "csv":
						csvData := "Target,Status,Result,Time\n"
						csvData += fmt.Sprintf("%s,completed,%s,%s\n", target, result["result"], result["time"])
						filename += ".csv"
						err := os.WriteFile(filename, []byte(csvData), 0666)
						if err != nil {
							return "", fmt.Errorf("failed to write CSV file: %v", err)
						}

					default:
						return "", fmt.Errorf("unsupported format: %s", format)
					}

					return fmt.Sprintf("Result exported to %s", filename), nil
				}
			}
		}
	}

	return "", fmt.Errorf("target not found: %s", target)
}

// GenerateReport generates a scan report (frontend compatible)
func GenerateReport(reportType string, format string, title string) (string, error) {
	var reportData string
	var filename string

	if title == "" {
		title = "dddd 扫描报告"
	}

	reportData = fmt.Sprintf("# %s\n\n", title)
	reportData += fmt.Sprintf("生成时间: %s\n\n", time.Now().Format(time.RFC3339))

	if reportType == "summary" {
		reportData += "## 概览\n\n"
		reportData += fmt.Sprintf("总任务数: %d\n", len(taskResults))
		
		totalTargets := 0
		for _, result := range taskResults {
			if targets, ok := result["targets"].([]string); ok {
				totalTargets += len(targets)
			}
		}
		reportData += fmt.Sprintf("总目标数: %d\n\n", totalTargets)
	} else if reportType == "vulnerability" {
		reportData += "## 漏洞报告\n\n"
		for taskID, result := range taskResults {
			if res, ok := result["result"].(string); ok && res != "" && res != "Scan completed" {
				reportData += fmt.Sprintf("### 任务: %s\n", taskID)
				targets, ok := result["targets"].([]string)
				if ok {
					for _, target := range targets {
						reportData += fmt.Sprintf("- 目标: %s\n", target)
					}
				}
				reportData += fmt.Sprintf("结果: %s\n\n", res)
			}
		}
	} else {
		reportData += "## 详细报告\n\n"
		for taskID, result := range taskResults {
			reportData += fmt.Sprintf("### 任务: %s\n", taskID)
			targets, ok := result["targets"].([]string)
			if ok {
				reportData += fmt.Sprintf("目标数: %d\n", len(targets))
				reportData += "目标列表:\n"
				for _, target := range targets {
					reportData += fmt.Sprintf("- %s\n", target)
				}
			}
			reportData += fmt.Sprintf("结果: %s\n", result["result"])
			reportData += fmt.Sprintf("时间: %s\n\n", result["time"])
		}
	}

	filename = fmt.Sprintf("dddd-report-%s-%s", reportType, time.Now().Format("20060102-150405"))
	if format == "html" {
		filename += ".html"
		reportData = "<!DOCTYPE html><html><head><title>" + title + "</title><style>body{font-family:sans-serif;max-width:800px;margin:0 auto;padding:20px;}h1{color:#1d1d1f;}h2{color:#6e6e73;border-bottom:1px solid #e8e8ed;padding-bottom:10px;}h3{color:#1d1d1f;margin-top:20px;}ul{list-style-type:disc;margin-left:20px;}</style></head><body>" +
			reportData + "</body></html>"
		// Convert markdown-like content to HTML
		reportData = strings.ReplaceAll(reportData, "# ", "<h1>")
		reportData = strings.ReplaceAll(reportData, "#\n", "</h1>\n")
		reportData = strings.ReplaceAll(reportData, "## ", "<h2>")
		reportData = strings.ReplaceAll(reportData, "##\n", "</h2>\n")
		reportData = strings.ReplaceAll(reportData, "### ", "<h3>")
		reportData = strings.ReplaceAll(reportData, "###\n", "</h3>\n")
		reportData = strings.ReplaceAll(reportData, "- ", "<li>")
		reportData = strings.ReplaceAll(reportData, "\n- ", "</li>\n<li>")
	} else {
		filename += ".md"
	}

	err := os.WriteFile(filename, []byte(reportData), 0666)
	if err != nil {
		return "", fmt.Errorf("failed to write report file: %v", err)
	}

	return fmt.Sprintf("Report generated: %s", filename), nil
}