# Bug修复计划

## 概述

根据全面代码检测报告，发现了多个需要修复的问题。本计划将按照优先级顺序修复这些问题。

---

## 修复优先级

### 高优先级 - 立即修复

1. **编译错误修复**
2. **功能逻辑错误修复**

### 中优先级 - 短期优化

1. **用户体验优化**
2. **代码质量提升**

---

## 修复计划详情

### 第一阶段：编译错误修复

| 序号 | 文件 | 问题 | 修复方案 |
|------|------|------|---------|
| 1 | `internal/api/config.go` | `SaveConfigToFile()`未定义 | 修改为`SaveConfig()` |
| 2 | `internal/api/result.go` | 缺少`strings`包导入 | 添加`"strings"`导入 |
| 3 | `frontend/src/views/ScanTask.vue` | 未导入`ElMessageBox` | 添加导入 |
| 4 | `frontend/src/views/TargetManagement.vue` | `deleteGroup`变量引用错误 | 修改为`groups.value` |
| 5 | `frontend/src/views/FingerprintManagement.vue` | `EditFingerprint`参数错误 | 修正参数个数 |

### 第二阶段：功能逻辑修复

| 序号 | 文件 | 问题 | 修复方案 |
|------|------|------|---------|
| 6 | `cmd/dddd-gui/main.go` | `InitTaskCleanup()`未调用 | 在startup中调用 |
| 7 | `internal/api/storage.go` | `AutoSaveTasks()`未启动 | 在startup中调用 |
| 8 | `internal/api/main.go` | `LoadDictionaries()`已调用 | 确认状态 |

### 第三阶段：用户体验优化

| 序号 | 文件 | 问题 | 修复方案 |
|------|------|------|---------|
| 9 | `frontend/src/views/TargetManagement.vue` | 搜索框布局杂乱 | 重新组织布局 |
| 10 | `frontend/src/views/ScanTask.vue` | 缺少删除确认对话框 | 添加确认对话框 |

---

## 执行步骤

### 步骤1：修复后端编译错误

1. 修改 `config.go` 中的 `SaveConfigToFile()` 为 `SaveConfig()`
2. 在 `result.go` 中添加 `strings` 包导入
3. 在 `main.go` 中添加 `InitTaskCleanup()` 调用

### 步骤2：修复前端编译错误

1. 在 `ScanTask.vue` 中添加 `ElMessageBox` 导入
2. 修复 `TargetManagement.vue` 中的变量引用错误
3. 修复 `FingerprintManagement.vue` 中的API调用参数

### 步骤3：添加缺失的API绑定

1. 在 `main.go` 中绑定 `GetDictionaries` API

---

## 风险评估

| 风险 | 等级 | 缓解措施 |
|------|------|---------|
| 编译失败 | 低 | 修复前备份代码 |
| 数据丢失 | 低 | 测试环境验证 |
| API兼容性 | 中 | 前后端同步更新 |

---

## 验证方法

1. **编译验证**: `go build ./...`
2. **前端构建**: `npm run build`
3. **功能测试**: 手动测试各功能模块

---

## 预期成果

修复完成后，项目应能够正常编译、构建和运行，所有基础功能可正常使用。