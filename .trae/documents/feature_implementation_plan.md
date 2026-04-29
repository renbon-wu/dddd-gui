# 功能实现方案规划

## 一、批量操作优化

### 1.1 需求分析

当前批量操作功能较为基础，需要优化以下方面：
- 支持批量选择后统一操作
- 提供批量导入/导出功能

### 1.2 实现方案

#### 1.2.1 批量导入优化

**新增导入格式支持**：
- TXT格式（一行一个目标）
- CSV格式（支持多列数据）
- JSON格式（结构化数据）

#### 1.2.2 批量导出功能

**导出格式**：
- TXT格式
- JSON格式（包含完整信息）

### 1.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `frontend/src/views/TargetManagement.vue` | 添加批量导入/导出 |
| `internal/api/target.go` | 添加批量导入导出API |

---

## 二、搜索过滤功能

### 2.1 需求分析

在各列表页面添加搜索框，支持快速过滤：
- 指纹管理页面
- POC管理页面
- 目标管理页面
- 结果分析页面

### 2.2 实现方案

#### 2.2.1 统一搜索组件

**创建通用搜索组件**：`components/SearchBar.vue`

```vue
<template>
  <el-input
    v-model="searchText"
    placeholder="搜索..."
    prefix-icon="Search"
    @input="handleSearch"
    clearable
  >
    <template #append>
      <el-select v-model="searchField" placeholder="字段">
        <el-option label="全部" value="all"></el-option>
        <el-option label="名称" value="name"></el-option>
        <el-option label="描述" value="desc"></el-option>
      </el-select>
    </template>
  </el-input>
</template>
```

#### 2.2.2 在各页面集成

**指纹管理页面**：
- 按指纹名称搜索
- 按规则内容搜索

**POC管理页面**：
- 按POC名称搜索
- 按标签搜索

**目标管理页面**：
- 按目标地址搜索
- 按分组搜索

**结果分析页面**：
- 按目标搜索
- 按漏洞名称搜索
- 按级别搜索

### 2.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `frontend/src/components/SearchBar.vue` | 新建通用搜索组件 |
| `frontend/src/views/FingerprintManagement.vue` | 集成搜索功能 |
| `frontend/src/views/PocManagement.vue` | 集成搜索功能 |
| `frontend/src/views/TargetManagement.vue` | 集成搜索功能 |
| `frontend/src/views/ResultAnalysis.vue` | 集成搜索功能 |

---

## 三、智能目标识别

### 3.1 需求分析

自动识别输入的目标类型：
- IP地址：自动识别为单个目标
- CIDR网段：自动识别为网段扫描
- 域名：自动识别为网站扫描
- 支持识别Hunter/Fofa查询语句

### 3.2 实现方案

#### 3.2.1 目标类型检测算法

**正则表达式规则**：

```javascript
// IP地址检测
const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/

// CIDR网段检测
const cidrRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/([0-9]|[1-2][0-9]|3[0-2])$/

// 域名检测
const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$/

// URL检测
const urlRegex = /^https?:\/\/[^\s]+$/

// Hunter查询检测
const hunterQueryRegex = /^ip="|^icp\.name=|^icp\.domain=/

// Fofa查询检测
const fofaQueryRegex = /^domain="|^ip="|^host=/
```

#### 3.2.2 识别逻辑

```javascript
const detectTargetType = (target) => {
  if (cidrRegex.test(target)) {
    return { type: 'cidr', value: target, label: 'CIDR网段' }
  } else if (ipRegex.test(target)) {
    return { type: 'ip', value: target, label: 'IP地址' }
  } else if (urlRegex.test(target)) {
    return { type: 'url', value: target, label: 'URL地址' }
  } else if (domainRegex.test(target)) {
    return { type: 'domain', value: target, label: '域名' }
  } else if (hunterQueryRegex.test(target)) {
    return { type: 'hunter', value: target, label: 'Hunter查询' }
  } else if (fofaQueryRegex.test(target)) {
    return { type: 'fofa', value: target, label: 'Fofa查询' }
  }
  return { type: 'unknown', value: target, label: '未知类型' }
}
```

#### 3.2.3 UI反馈

在输入框旁显示识别结果：

```vue
<el-input v-model="target" @blur="detectType">
</el-input>
<el-tag :type="getTagType(detectedType)" size="small">
  {{ detectedType.label }}
</el-tag>
```

### 3.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `frontend/src/utils/targetDetector.js` | 新建目标类型检测工具 |
| `frontend/src/views/ScanTask.vue` | 集成智能识别功能 |
| `frontend/src/views/TargetManagement.vue` | 集成智能识别功能 |

---

## 四、一键扫描模板

### 4.1 需求分析

在首页添加常用扫描模板的快捷入口：
- 内网扫描（默认端口）
- 外网扫描（全端口）
- 仅指纹识别
- 完整漏洞扫描
- 子域名枚举
- Hunter/Fofa查询

### 4.2 实现方案

#### 4.2.1 模板配置结构

```javascript
const scanTemplates = [
  {
    id: 'internal',
    name: '内网扫描',
    description: '扫描内网资产，使用默认端口',
    icon: 'Network',
    config: {
      ports: 'Top1000',
      scanType: 'tcp',
      enablePoc: true,
      enableFingerprint: true
    }
  },
  {
    id: 'external',
    name: '外网扫描',
    description: '扫描外网资产，使用全端口',
    icon: 'Globe',
    config: {
      ports: '1-65535',
      scanType: 'syn',
      enablePoc: true,
      enableFingerprint: true
    }
  },
  {
    id: 'fingerprint-only',
    name: '仅指纹识别',
    description: '仅识别资产指纹，不进行漏洞扫描',
    icon: 'Fingerprint',
    config: {
      ports: 'Top1000',
      scanType: 'tcp',
      enablePoc: false,
      enableFingerprint: true
    }
  },
  {
    id: 'full-scan',
    name: '完整漏洞扫描',
    description: '全面扫描，包括指纹识别和漏洞检测',
    icon: 'ShieldAlert',
    config: {
      ports: 'Top2000',
      scanType: 'tcp',
      enablePoc: true,
      enableFingerprint: true,
      deepScan: true
    }
  },
  {
    id: 'subdomain',
    name: '子域名枚举',
    description: '枚举子域名并扫描',
    icon: 'Link',
    config: {
      ports: 'Top100',
      scanType: 'tcp',
      enablePoc: true,
      enableFingerprint: true,
      subdomain: true
    }
  }
]
```

#### 4.2.2 首页快捷入口

```vue
<div class="template-grid">
  <el-card
    v-for="template in scanTemplates"
    :key="template.id"
    class="template-card"
    @click="applyTemplate(template)"
  >
    <div class="template-icon">
      <el-icon><component :is="template.icon" /></el-icon>
    </div>
    <h3>{{ template.name }}</h3>
    <p>{{ template.description }}</p>
  </el-card>
</div>
```

#### 4.2.3 模板应用逻辑

```javascript
const applyTemplate = (template) => {
  scanForm.value = {
    ports: template.config.ports,
    scanType: template.config.scanType,
    enablePoc: template.config.enablePoc,
    enableFingerprint: template.config.enableFingerprint,
    deepScan: template.config.deepScan || false,
    subdomain: template.config.subdomain || false
  }
  ElMessage.success(`已应用「${template.name}」模板`)
}
```

### 4.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `frontend/src/views/ScanTask.vue` | 添加模板选择区域 |
| `frontend/src/utils/scanTemplates.js` | 新建扫描模板配置 |

---

## 五、指纹测试（格式检测+真实测试分开）

### 5.1 需求分析

指纹测试功能需要分开两种模式：
1. **格式检测**：仅验证指纹规则语法是否正确
2. **真实测试**：发送真实HTTP请求，验证指纹规则是否匹配

### 5.2 实现方案

#### 5.2.1 后端API实现

**修改文件**：`internal/api/fingerprint.go`

```go
// 仅验证指纹规则格式
func ValidateFingerprint(rule string) (string, error) {
    ruleData := ddfinger.ParseRule(rule)
    if len(ruleData) == 0 {
        return "", fmt.Errorf("无效的指纹规则语法")
    }
    return "指纹规则语法验证通过", nil
}

// 真实HTTP测试指纹规则
func TestFingerprint(rule string, target string) (string, error) {
    // 1. 先验证规则格式
    ruleData := ddfinger.ParseRule(rule)
    if len(ruleData) == 0 {
        return "", fmt.Errorf("无效的指纹规则语法")
    }
    
    // 2. 发送HTTP请求
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(target)
    if err != nil {
        return "", fmt.Errorf("无法连接目标: %v", err)
    }
    defer resp.Body.Close()
    
    // 3. 读取响应
    body, _ := io.ReadAll(resp.Body)
    headers := resp.Header
    
    // 4. 匹配指纹规则
    matches := ddfinger.Match(ruleData, string(body), headers, resp.StatusCode)
    
    result := map[string]interface{}{
        "success": matches,
        "statusCode": resp.StatusCode,
        "headers": headers,
        "bodyPreview": string(body[:min(500, len(body))]) + "...",
    }
    
    if matches {
        return fmt.Sprintf("指纹匹配成功: %v", result), nil
    }
    return fmt.Sprintf("指纹未匹配: %v", result), nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

#### 5.2.2 前端界面优化

```vue
<el-dialog v-model="testDialogVisible" title="测试指纹" width="60%">
  <el-form :model="testForm">
    <el-form-item label="指纹规则">
      <el-input v-model="testForm.rule" type="textarea" :rows="4"></el-input>
    </el-form-item>
    
    <el-form-item>
      <el-radio-group v-model="testMode" label="测试模式">
        <el-radio label="format">仅格式检测</el-radio>
        <el-radio label="real">真实HTTP请求</el-radio>
      </el-radio-group>
    </el-form-item>
    
    <el-form-item v-if="testMode === 'real'" label="测试目标">
      <el-input v-model="testForm.target" placeholder="https://example.com"></el-input>
    </el-form-item>
    
    <el-form-item>
      <div v-if="testResult" :class="['result-box', testResult.success ? 'success' : 'error']">
        <h4>{{ testResult.message }}</h4>
        <div v-if="testResult.details">
          <el-collapse>
            <el-collapse-item title="查看详细信息">
              <pre>{{ testResult.details }}</pre>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </el-form-item>
  </el-form>
  <template #footer>
    <el-button @click="testDialogVisible = false">取消</el-button>
    <el-button type="primary" @click="runTest" :loading="testing">测试</el-button>
  </template>
</el-dialog>
```

#### 5.2.3 前端测试逻辑

```javascript
const testMode = ref('format')

const runTest = async () => {
  if (!testForm.value.rule) {
    ElMessage.warning('请输入指纹规则')
    return
  }
  
  if (testMode.value === 'real' && !testForm.value.target) {
    ElMessage.warning('请输入测试目标')
    return
  }
  
  try {
    testing.value = true
    if (testMode.value === 'format') {
      result = await window.go.main.App.ValidateFingerprint(testForm.value.rule)
      testResult.value = {
        success: true,
        message: result
      }
    } else {
      result = await window.go.main.App.TestFingerprint(testForm.value.rule, testForm.value.target)
      testResult.value = {
        success: true,
        message: '测试完成',
        details: result
      }
    }
    ElMessage.success('测试完成')
  } catch (err) {
    testResult.value = {
      success: false,
      message: '测试失败',
      details: err
    }
    ElMessage.error('测试失败')
  } finally {
    testing.value = false
  }
}
```

### 5.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `internal/api/fingerprint.go` | 添加ValidateFingerprint和改进TestFingerprint |
| `frontend/src/views/FingerprintManagement.vue` | 优化测试界面，支持两种模式 |
| `cmd/dddd-gui/main.go` | 注册ValidateFingerprint API |

---

## 六、POC测试（格式检测+真实测试分开）

### 6.1 需求分析

POC测试功能需要分开两种模式：
1. **格式检测**：仅验证POC文件格式是否符合Nuclei规范
2. **真实测试**：调用Nuclei执行真实POC测试

### 6.2 实现方案

#### 6.2.1 后端API实现

**修改文件**：`internal/api/poc.go`

```go
// 仅验证POC文件格式
func ValidatePOC(pocName string) (string, error) {
    pocPath := findPOC(pocName)
    if pocPath == "" {
        return "", fmt.Errorf("POC未找到: %s", pocName)
    }
    
    // 读取并解析POC文件验证格式
    data, err := os.ReadFile(pocPath)
    if err != nil {
        return "", fmt.Errorf("读取POC文件失败: %v", err)
    }
    
    // 尝试解析为YAML验证格式
    var pocMap map[string]interface{}
    if err := yaml.Unmarshal(data, &pocMap); err != nil {
        return "", fmt.Errorf("POC文件格式错误: %v", err)
    }
    
    // 检查必要字段
    requiredFields := []string{"id", "name", "author", "severity"}
    for _, field := range requiredFields {
        if _, ok := pocMap[field]; !ok {
            return "", fmt.Errorf("POC缺少必要字段: %s", field)
        }
    }
    
    return fmt.Sprintf("POC格式验证通过: %s (severity: %v)", pocMap["name"], pocMap["severity"]), nil
}

// 真实调用Nuclei测试POC
func TestPOC(pocName string, target string) (string, error) {
    // 1. 先验证POC格式
    pocPath := findPOC(pocName)
    if pocPath == "" {
        return "", fmt.Errorf("POC未找到: %s", pocName)
    }
    
    // 2. 构建nuclei命令
    cmd := exec.Command("nuclei", "-u", target, "-t", pocPath, "-silent")
    
    // 3. 执行命令
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("Nuclei执行失败: %v", err)
    }
    
    resultStr := string(output)
    if strings.Contains(resultStr, "[+]") {
        return fmt.Sprintf("POC匹配成功:\n%s", resultStr), nil
    }
    
    return fmt.Sprintf("POC未匹配\n输出:\n%s", resultStr), nil
}
```

#### 6.2.2 前端界面优化

```vue
<el-dialog v-model="testDialogVisible" title="测试POC" width="60%">
  <el-form :model="testForm">
    <el-form-item label="POC名称">
      <el-select v-model="testForm.poc" placeholder="选择POC">
        <el-option 
          v-for="poc in pocs" 
          :key="poc.id" 
          :label="poc.name" 
          :value="poc.name"
        ></el-option>
      </el-select>
    </el-form-item>
    
    <el-form-item>
      <el-radio-group v-model="testMode" label="测试模式">
        <el-radio label="format">仅格式检测</el-radio>
        <el-radio label="real">调用Nuclei</el-radio>
      </el-radio-group>
    </el-form-item>
    
    <el-form-item v-if="testMode === 'real'" label="测试目标">
      <el-input v-model="testForm.target" placeholder="https://example.com"></el-input>
    </el-form-item>
    
    <el-form-item>
      <div v-if="testResult" class="result-box">
        <h4>{{ testResult.message }}</h4>
        <div v-if="testResult.details" class="mt-3">
          <el-collapse>
            <el-collapse-item title="查看输出详情">
              <pre class="bg-gray-800 text-white p-3 rounded text-sm">{{ testResult.details }}</pre>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </el-form-item>
  </el-form>
  <template #footer>
    <el-button @click="testDialogVisible = false">取消</el-button>
    <el-button type="primary" @click="runTest" :loading="testing">测试</el-button>
  </template>
</el-dialog>
```

### 6.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `internal/api/poc.go` | 添加ValidatePOC和改进TestPOC |
| `frontend/src/views/PocManagement.vue` | 优化测试界面，支持两种模式 |
| `cmd/dddd-gui/main.go` | 注册ValidatePOC API |

---

## 七、扫描进度显示

### 7.1 需求分析

添加扫描进度实时显示功能。

### 7.2 实现方案

#### 7.2.1 WebSocket实时通信

**后端实现**：

```go
// 扫描进度结构体
type ScanProgress struct {
    TaskID    string  `json:"taskId"`
    Total     int     `json:"total"`
    Completed int     `json:"completed"`
    Progress  float64 `json:"progress"`
    Status    string  `json:"status"`
    Message   string  `json:"message"`
}

// 使用通道传递进度
func RunScanWithProgress(taskID string, targets []string, options map[string]interface{}) (<-chan ScanProgress, error) {
    progressChan := make(chan ScanProgress)
    
    go func() {
        total := len(targets)
        for i, target := range targets {
            // 扫描目标
            scanTarget(target, options)
            
            // 发送进度
            progressChan <- ScanProgress{
                TaskID:    taskID,
                Total:     total,
                Completed: i + 1,
                Progress:  float64(i+1) / float64(total) * 100,
                Status:    "running",
                Message:   fmt.Sprintf("扫描中: %s", target),
            }
        }
        
        // 完成
        progressChan <- ScanProgress{
            TaskID:    taskID,
            Total:     total,
            Completed: total,
            Progress:  100,
            Status:    "completed",
            Message:   "扫描完成",
        }
        
        close(progressChan)
    }()
    
    return progressChan, nil
}
```

#### 7.2.2 前端进度展示

```vue
<el-dialog v-model="scanProgressVisible" title="扫描进度" width="60%">
  <div class="progress-container">
    <el-progress :percentage="progress" :status="progressStatus" :stroke-width="20"></el-progress>
    <div class="progress-info">
      <span>{{ progressMessage }}</span>
      <span class="float-right">{{ completed }} / {{ total }}</span>
    </div>
    <div class="progress-log">
      <div v-for="(log, index) in scanLogs" :key="index" class="log-item">
        <el-tag type="info" size="small" effect="dark">{{ index + 1 }}</el-tag>
        <span class="ml-2">{{ log }}</span>
      </div>
    </div>
  </div>
  <template #footer>
    <el-button @click="stopScan" type="danger">停止扫描</el-button>
    <el-button @click="closeProgress" type="primary">关闭</el-button>
  </template>
</el-dialog>
```

### 7.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `internal/api/task.go` | 添加进度通道支持 |
| `frontend/src/views/ScanTask.vue` | 添加进度对话框 |

---

## 八、结果分析图表展示

### 8.1 需求分析

完善可视化图表展示功能。

### 8.2 实现方案

#### 8.2.1 图表类型

1. **漏洞级别分布** - 饼图/环形图
2. **指纹类型分布** - 柱状图
3. **扫描趋势** - 折线图
4. **端口开放分布** - 热力图/饼图

#### 8.2.2 使用Chart.js实现

```javascript
import { Chart, ArcElement, Tooltip, Legend, PieController, BarController, LineController, CategoryScale, LinearScale, PointElement, LineElement, BarElement } from 'chart.js'

Chart.register(ArcElement, Tooltip, Legend, PieController, BarController, LineController, CategoryScale, LinearScale, PointElement, LineElement, BarElement)

// 漏洞级别分布饼图
const levelChart = new Chart(levelChartRef.value, {
  type: 'pie',
  data: {
    labels: ['高危', '中危', '低危', '信息'],
    datasets: [{
      data: [highCount, mediumCount, lowCount, infoCount],
      backgroundColor: ['#ff3b30', '#ff9500', '#007aff', '#909399']
    }]
  },
  options: {
    responsive: true,
    plugins: {
      legend: {
        position: 'bottom',
      },
    },
  },
})

// 指纹类型分布柱状图
const fingerprintChart = new Chart(fingerprintChartRef.value, {
  type: 'bar',
  data: {
    labels: fingerprintNames,
    datasets: [{
      label: '数量',
      data: fingerprintCounts,
      backgroundColor: '#007aff'
    }]
  },
  options: {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
    },
  },
})

// 扫描趋势折线图
const trendChart = new Chart(trendChartRef.value, {
  type: 'line',
  data: {
    labels: scanDates,
    datasets: [{
      label: '发现资产',
      data: assetCounts,
      borderColor: '#007aff',
      tension: 0.1
    }]
  },
  options: {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
    },
  },
})
```

#### 8.2.3 前端界面

```vue
<div class="charts-grid">
  <el-card class="chart-card">
    <template #header>漏洞级别分布</template>
    <canvas ref="levelChart"></canvas>
  </el-card>
  <el-card class="chart-card">
    <template #header>指纹类型分布</template>
    <canvas ref="fingerprintChart"></canvas>
  </el-card>
  <el-card class="chart-card">
    <template #header>扫描趋势</template>
    <canvas ref="trendChart"></canvas>
  </el-card>
  <el-card class="chart-card">
    <template #header>端口开放分布</template>
    <canvas ref="portChart"></canvas>
  </el-card>
</div>
```

### 8.3 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `frontend/src/views/ResultAnalysis.vue` | 添加图表展示 |
| `frontend/src/utils/chartHelper.js` | 新建图表工具函数 |
| `frontend/package.json` | 添加Chart.js依赖 |

---

## 九、dddd功能GUI缺失检查补充

### 9.1 需求分析

根据`/workspace/details.md`检查dddd功能的GUI覆盖情况，添加缺失的功能界面。

### 9.2 功能检查列表

| dddd功能 | GUI覆盖情况 | 状态 |
|---------|------------|-----|
| 端口扫描（TCP/SYN） | 扫描任务页面 | ✅ 已有 |
| 主机发现（ICMP/TCP） | 需要添加设置项 | ⚠️ 补充 |
| 协议识别（Nmap） | 需要添加设置项 | ⚠️ 补充 |
| 子域名枚举（爆破+被动） | 需要添加设置项 | ⚠️ 补充 |
| 主动指纹探测 | 需要添加设置项 | ⚠️ 补充 |
| Hunter查询 | 需要添加API配置和查询页面 | ⚠️ 补充 |
| Fofa查询 | 需要添加API配置和查询页面 | ⚠️ 补充 |
| Quake查询 | 需要添加API配置和查询页面 | ⚠️ 补充 |
| 弱口令爆破（FTP/SSH/SMB等） | 需要添加新页面 | ⚠️ 补充 |
| 审计日志 | 需要添加日志查看页面 | ⚠️ 补充 |
| 结果导出（TXT/JSON/HTML） | 结果分析页面 | ✅ 已有 |
| 配置文件管理 | 系统设置页面 | ✅ 已有 |

### 9.3 新增功能实现方案

#### 9.3.1 扫描高级设置页面

在扫描任务页面添加高级设置折叠面板：

```vue
<el-collapse>
  <el-collapse-item title="高级设置">
    <el-form :model="advancedConfig">
      <el-form-item label="主机发现">
        <el-checkbox v-model="advancedConfig.enableIcmp">启用ICMP</el-checkbox>
        <el-checkbox v-model="advancedConfig.enableTcpPing" class="ml-4">启用TCP Ping</el-checkbox>
      </el-form-item>
      <el-form-item label="协议识别">
        <el-switch v-model="advancedConfig.enableProtocol"></el-switch>
        <span class="ml-2 text-gray-500">使用Nmap识别协议</span>
      </el-form-item>
      <el-form-item label="主动指纹探测">
        <el-switch v-model="advancedConfig.enableActiveFingerprint"></el-switch>
      </el-form-item>
      <el-form-item label="子域名枚举">
        <el-checkbox v-model="advancedConfig.enableSubdomain">启用子域名枚举</el-checkbox>
        <el-checkbox v-model="advancedConfig.enableSubdomainBrute" class="ml-4">启用爆破</el-checkbox>
        <el-checkbox v-model="advancedConfig.enableSubfinder" class="ml-4">启用Subfinder</el-checkbox>
      </el-form-item>
      <el-form-item label="审计日志">
        <el-switch v-model="advancedConfig.enableAuditLog"></el-switch>
        <span class="ml-2 text-gray-500">记录详细扫描行为</span>
      </el-form-item>
    </el-form>
  </el-collapse-item>
</el-collapse>
```

#### 9.3.2 网络空间搜索引擎查询页面

新增`frontend/src/views/SpaceSearch.vue`：

```vue
<template>
  <div class="space-search">
    <el-card class="content-section mb-4">
      <template #header>网络空间资产搜索</template>
      <el-radio-group v-model="searchType" style="margin-bottom: 16px;">
        <el-radio label="hunter">Hunter</el-radio>
        <el-radio label="fofa">Fofa</el-radio>
        <el-radio label="quake">Quake</el-radio>
      </el-radio-group>
      <el-form :model="searchForm">
        <el-form-item label="查询语句">
          <el-input v-model="searchForm.query" type="textarea" placeholder="输入查询语句"></el-input>
        </el-form-item>
        <el-form-item v-if="searchType === 'hunter'" label="最大页数">
          <el-input-number v-model="searchForm.hunterPageCount" :min="1" :max="100"></el-input-number>
        </el-form-item>
        <el-form-item v-if="searchType === 'fofa'" label="最大数量">
          <el-input-number v-model="searchForm.fofaMaxCount" :min="1" :max="10000"></el-input-number>
        </el-form-item>
        <el-form-item v-if="searchType === 'quake'" label="最大数量">
          <el-input-number v-model="searchForm.quakeMaxCount" :min="1" :max="10000"></el-input-number>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="searchAndScan">查询并扫描</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card class="content-section">
      <template #header>查询结果</template>
      <el-table :data="searchResults" style="width: 100%">
        <el-table-column prop="ip" label="IP"></el-table-column>
        <el-table-column prop="port" label="端口" width="100"></el-table-column>
        <el-table-column prop="domain" label="域名"></el-table-column>
        <el-table-column prop="title" label="标题" show-overflow-tooltip></el-table-column>
        <el-table-column prop="action" label="操作">
          <template #default="scope">
            <el-button size="small" @click="addToTargets(scope.row)">添加目标</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
```

#### 9.3.3 弱口令爆破页面

新增`frontend/src/views/BruteForce.vue`：

```vue
<template>
  <div class="brute-force">
    <el-card class="content-section mb-4">
      <template #header>弱口令爆破</template>
      <el-form :model="bruteConfig">
        <el-form-item label="目标">
          <el-input v-model="bruteConfig.targets" type="textarea" placeholder="192.168.0.1/24, 一行一个"></el-input>
        </el-form-item>
        <el-form-item label="服务类型">
          <el-checkbox-group v-model="bruteConfig.services">
            <el-checkbox label="ssh">SSH</el-checkbox>
            <el-checkbox label="ftp">FTP</el-checkbox>
            <el-checkbox label="smb">SMB</el-checkbox>
            <el-checkbox label="mysql">MySQL</el-checkbox>
            <el-checkbox label="mssql">MSSQL</el-checkbox>
            <el-checkbox label="redis">Redis</el-checkbox>
            <el-checkbox label="rdp">RDP</el-checkbox>
            <el-checkbox label="telnet">Telnet</el-checkbox>
            <el-checkbox label="postgresql">PostgreSQL</el-checkbox>
            <el-checkbox label="oracle">Oracle</el-checkbox>
            <el-checkbox label="mongodb">MongoDB</el-checkbox>
            <el-checkbox label="shiro">Shiro Key</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="凭证">
          <el-radio-group v-model="bruteConfig.credentialMode">
            <el-radio label="manual">手动输入</el-radio>
            <el-radio label="file">文件导入</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="bruteConfig.credentialMode === 'manual'" label="凭证内容">
          <el-input v-model="bruteConfig.credentials" type="textarea" placeholder="格式: username : password, 一行一个"></el-input>
        </el-form-item>
        <el-form-item v-if="bruteConfig.credentialMode === 'file'" label="凭证文件">
          <el-upload :auto-upload="false" :on-change="handleFileChange" accept=".txt">
            <el-button type="primary">选择文件</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="startBrute">开始爆破</el-button>
          <el-button @click="stopBrute" type="danger">停止</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card class="content-section">
      <template #header>爆破结果</template>
      <el-table :data="bruteResults" style="width: 100%">
        <el-table-column prop="target" label="目标"></el-table-column>
        <el-table-column prop="service" label="服务" width="100"></el-table-column>
        <el-table-column prop="username" label="用户名" width="120"></el-table-column>
        <el-table-column prop="password" label="密码" width="120"></el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag type="success" size="small" v-if="scope.row.success">成功</el-tag>
            <el-tag type="info" size="small" v-else>失败</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
```

#### 9.3.4 审计日志页面

新增`frontend/src/views/AuditLog.vue`：

```vue
<template>
  <div class="audit-log">
    <el-card class="content-section mb-4">
      <template #header>审计日志</template>
      <el-form :inline="true">
        <el-form-item label="任务ID">
          <el-input v-model="taskIdFilter" placeholder="输入任务ID"></el-input>
        </el-form-item>
        <el-form-item label="日志级别">
          <el-select v-model="levelFilter" placeholder="全部">
            <el-option label="全部" value=""></el-option>
            <el-option label="INFO" value="INFO"></el-option>
            <el-option label="WARN" value="WARN"></el-option>
            <el-option label="ERROR" value="ERROR"></el-option>
            <el-option label="REQUEST" value="REQUEST"></el-option>
            <el-option label="RESPONSE" value="RESPONSE"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="refreshLogs">刷新</el-button>
          <el-button @click="clearLogs">清空</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card class="content-section">
      <template #header>日志列表</template>
      <div class="log-container">
        <div v-for="(log, index) in filteredLogs" :key="index" class="log-item">
          <div class="log-header">
            <el-tag :type="getLogType(log.level)" size="small">{{ log.level }}</el-tag>
            <span class="log-time ml-2 text-gray-500">{{ log.time }}</span>
            <span class="log-task-id ml-2 text-gray-500">[{{ log.taskId }}]</span>
          </div>
          <div class="log-message mt-1">{{ log.message }}</div>
          <div v-if="log.details" class="log-details">
            <el-collapse>
              <el-collapse-item title="查看详情">
                <pre class="bg-gray-800 text-white p-3 rounded text-sm">{{ log.details }}</pre>
              </el-collapse-item>
            </el-collapse>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>
```

### 9.4 后端API补充

需要在对应API文件中补充功能：
- 网络空间搜索引擎API（Hunter/Fofa/Quake查询）
- 弱口令爆破API
- 审计日志查询API
- 完整扫描参数支持（传递dddd的所有命令行参数）

### 9.5 代码文件

| 文件 | 修改内容 |
|-----|---------|
| `frontend/src/views/SpaceSearch.vue` | 新建网络空间搜索页面 |
| `frontend/src/views/BruteForce.vue` | 新建弱口令爆破页面 |
| `frontend/src/views/AuditLog.vue` | 新建审计日志页面 |
| `frontend/src/views/ScanTask.vue` | 添加高级设置面板 |
| `frontend/src/App.vue` | 添加新的导航菜单项 |
| `internal/api/space.go` | 新建网络空间搜索API |
| `internal/api/brute.go` | 新建弱口令爆破API |
| `internal/api/audit.go` | 新建审计日志API |
| `internal/api/task.go` | 补充完整扫描参数支持 |
| `cmd/dddd-gui/main.go` | 注册新API函数 |

---

## 十、实施计划

### 10.1 第一阶段（基础优化）

| 任务 | 时间 | 负责人 |
|-----|------|-------|
| 搜索过滤功能 | 1天 | 前端开发 |
| 智能目标识别（支持Hunter/Fofa） | 1天 | 前端开发 |
| 一键扫描模板（增加子域名枚举等） | 1天 | 前端开发 |

### 10.2 第二阶段（核心功能）

| 任务 | 时间 | 负责人 |
|-----|------|-------|
| 指纹测试（格式检测+真实测试分开） | 2天 | 前后端协作 |
| POC测试（格式检测+真实测试分开） | 2天 | 前后端协作 |
| 扫描进度显示 | 2天 | 前后端协作 |

### 10.3 第三阶段（可视化）

| 任务 | 时间 | 负责人 |
|-----|------|-------|
| 结果分析图表展示 | 2天 | 前端开发 |
| 批量操作优化 | 1天 | 前端开发 |

### 10.4 第四阶段（功能补充）

| 任务 | 时间 | 负责人 |
|-----|------|-------|
| 网络空间搜索引擎页面（Hunter/Fofa/Quake） | 3天 | 前后端协作 |
| 弱口令爆破页面 | 3天 | 前后端协作 |
| 审计日志页面 | 2天 | 前后端协作 |
| 扫描高级设置补充 | 2天 | 前端开发 |
| 完整测试和优化 | 2天 | 测试人员 |

---

## 十一、依赖与资源

### 11.1 前端依赖

| 依赖 | 版本 | 用途 |
|-----|------|-----|
| chart.js | ^4.0.0 | 图表展示 |
| @vueuse/core | ^10.0.0 | 工具函数 |

### 11.2 后端依赖

| 依赖 | 用途 |
|-----|-----|
| gorilla/websocket | WebSocket通信（可选） |

### 11.3 外部工具

| 工具 | 用途 |
|-----|-----|
| nuclei | POC测试执行 |
| nmap | 协议识别 |
| subfinder | 子域名枚举 |
| masscan | SYN扫描 |

---

## 十二、最终导航菜单结构

```
├── 扫描任务
│   ├── 新建扫描
│   ├── 扫描模板
│   └── 历史任务
├── 目标管理
│   ├── 目标列表
│   ├── 目标分组
│   └── 批量导入
├── 指纹管理
│   ├── 指纹列表
│   ├── 添加指纹
│   └── 测试指纹（格式+真实）
├── POC管理
│   ├── POC列表
│   ├── 上传POC
│   ├── 工作流管理
│   └── 测试POC（格式+真实）
├── 网络空间搜索
│   ├── Hunter查询
│   ├── Fofa查询
│   └── Quake查询
├── 弱口令爆破
│   ├── 新建爆破
│   └── 爆破结果
├── 结果分析
│   ├── 结果概览
│   ├── 图表展示
│   └── 导出报告
├── 审计日志
│   ├── 日志列表
│   └── 日志搜索
└── 系统设置
    ├── 全局配置
    ├── API配置
    ├── 字典管理
    └── 界面设置
```