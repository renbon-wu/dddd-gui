package api

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var targets = make(map[string]map[string]interface{})
var targetGroups = make(map[string][]string)

const (
	MaxFileSize     = 10 * 1024 * 1024
	MaxTargetsCount = 10000
)

var (
	domainRegex  = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}(:[0-9]+)?$`)
	ipPortRegex  = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}:\d+$`)
)

// ImportTargets imports targets from a file (TXT, CSV)
func ImportTargets(filePath string) (string, error) {
	// 验证文件路径，防止目录遍历攻击
	cleanPath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanPath, ".") && !filepath.IsAbs(cleanPath) {
		cleanPath = "./" + cleanPath
	}

	// 检查文件大小
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %v", err)
	}

	if fileInfo.Size() > MaxFileSize {
		return "", fmt.Errorf("file too large, maximum size is %d bytes", MaxFileSize)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var importedTargets []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// For CSV, split by commas
			if strings.Contains(line, ",") {
				csvTargets := strings.Split(line, ",")
				for _, csvTarget := range csvTargets {
					csvTarget = strings.TrimSpace(csvTarget)
					if csvTarget != "" {
						// 验证目标格式（简单验证）
						if isValidTarget(csvTarget) {
							importedTargets = append(importedTargets, csvTarget)
						}
					}
				}
			} else {
				// 验证目标格式（简单验证）
				if isValidTarget(line) {
					importedTargets = append(importedTargets, line)
				}
			}
		}
	}

	// 检查目标数量限制
	if len(importedTargets) > MaxTargetsCount {
		return "", fmt.Errorf("too many targets, maximum is %d", MaxTargetsCount)
	}

	// Add targets to the map
	for _, target := range importedTargets {
		targets[target] = map[string]interface{}{
			"target": target,
			"status": "未扫描",
			"time":   time.Now().Format(time.RFC3339),
		}
	}

	return fmt.Sprintf("Imported %d targets successfully", len(importedTargets)), nil
}

// isValidTarget 验证目标格式是否有效
func isValidTarget(target string) bool {
	if target == "" || len(target) > 255 {
		return false
	}

	if net.ParseIP(target) != nil {
		return true
	}

	if _, _, err := net.ParseCIDR(target); err == nil {
		return true
	}

	if strings.Contains(target, "-") {
		parts := strings.Split(target, "-")
		if len(parts) == 2 {
			ip1 := net.ParseIP(strings.TrimSpace(parts[0]))
			ip2 := net.ParseIP(strings.TrimSpace(parts[1]))
			if ip1 != nil && ip2 != nil {
				return true
			}
		}
	}

	if strings.Contains(target, "://") {
		u, err := url.Parse(target)
		if err == nil && u.Host != "" {
			return true
		}
	}

	if domainRegex.MatchString(target) {
		return true
	}

	if ipPortRegex.MatchString(target) {
		parts := strings.Split(target, ":")
		if len(parts) == 2 {
			ip := net.ParseIP(parts[0])
			if ip != nil {
				return true
			}
		}
	}

	return false
}

// GetTargets returns the list of targets
func GetTargets() (map[string]map[string]interface{}, error) {
	return targets, nil
}

// BatchDeleteTargets deletes multiple targets
func BatchDeleteTargets(targetList []string) (string, error) {
	deletedCount := 0
	for _, target := range targetList {
		if _, exists := targets[target]; exists {
			delete(targets, target)
			deletedCount++
		}
	}
	return fmt.Sprintf("Deleted %d targets successfully", deletedCount), nil
}

// BatchScanTargets scans multiple targets
func BatchScanTargets(targetList []string) (string, error) {
	if len(targetList) == 0 {
		return "", fmt.Errorf("no targets to scan")
	}

	// Run the scan
	result, err := RunScan(targetList, map[string]interface{}{})
	if err != nil {
		return "", err
	}

	// Update target status
	for _, target := range targetList {
		if targetInfo, exists := targets[target]; exists {
			targetInfo["status"] = "已扫描"
			targetInfo["scanTime"] = time.Now().Format(time.RFC3339)
			targets[target] = targetInfo
		}
	}

	return result, nil
}

// CreateTargetGroup creates a target group
func CreateTargetGroup(groupName string, targetList []string) (string, error) {
	targetGroups[groupName] = targetList
	return fmt.Sprintf("Created group %s with %d targets", groupName, len(targetList)), nil
}

// GetTargetGroups returns the list of target groups
func GetTargetGroups() (map[string][]string, error) {
	return targetGroups, nil
}

// AddTargetToGroup adds a target to a group
func AddTargetToGroup(groupName string, target string) (string, error) {
	if _, exists := targetGroups[groupName]; !exists {
		targetGroups[groupName] = []string{}
	}
	targetGroups[groupName] = append(targetGroups[groupName], target)
	return fmt.Sprintf("Added target %s to group %s", target, groupName), nil
}

// EditTarget edits a target's address
func EditTarget(oldTarget string, newTarget string) (string, error) {
	if _, exists := targets[oldTarget]; !exists {
		return "", fmt.Errorf("target not found")
	}
	if _, exists := targets[newTarget]; exists {
		return "", fmt.Errorf("target already exists")
	}

	targets[newTarget] = targets[oldTarget]
	delete(targets, oldTarget)

	for groupName, groupTargets := range targetGroups {
		for i, t := range groupTargets {
			if t == oldTarget {
				groupTargets[i] = newTarget
				targetGroups[groupName] = groupTargets
				break
			}
		}
	}

	return fmt.Sprintf("Target updated from '%s' to '%s'", oldTarget, newTarget), nil
}

// RemoveTargetFromGroup removes a target from a group
func RemoveTargetFromGroup(groupName string, target string) (string, error) {
	if _, exists := targetGroups[groupName]; !exists {
		return "", fmt.Errorf("group not found")
	}

	var newTargets []string
	for _, t := range targetGroups[groupName] {
		if t != target {
			newTargets = append(newTargets, t)
		}
	}
	targetGroups[groupName] = newTargets

	return fmt.Sprintf("Removed target %s from group %s", target, groupName), nil
}

// ImportTargetsFromText imports targets from text content
func ImportTargetsFromText(content string) (string, error) {
	lines := strings.Split(content, "\n")
	var importedTargets []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if strings.Contains(line, ",") {
				csvTargets := strings.Split(line, ",")
				for _, csvTarget := range csvTargets {
					csvTarget = strings.TrimSpace(csvTarget)
					if csvTarget != "" && isValidTarget(csvTarget) {
						importedTargets = append(importedTargets, csvTarget)
					}
				}
			} else {
				if isValidTarget(line) {
					importedTargets = append(importedTargets, line)
				}
			}
		}
	}

	if len(importedTargets) > MaxTargetsCount {
		return "", fmt.Errorf("too many targets, maximum is %d", MaxTargetsCount)
	}

	for _, target := range importedTargets {
		targets[target] = map[string]interface{}{
			"target": target,
			"status": "未扫描",
			"time":   time.Now().Format(time.RFC3339),
		}
	}

	return fmt.Sprintf("Imported %d targets successfully", len(importedTargets)), nil
}