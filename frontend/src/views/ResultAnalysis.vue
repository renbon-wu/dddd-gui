<template>
  <div class="result-analysis">
    <div class="content-section">
      <h2>结果分析</h2>
      
      <!-- 结果概览 -->
      <div class="overview-cards">
        <el-card class="overview-card">
          <div class="overview-item">
            <span class="overview-label">总扫描目标</span>
            <span class="overview-value">{{ overview.totalTargets }}</span>
          </div>
          <div class="overview-item">
            <span class="overview-label">发现指纹</span>
            <span class="overview-value">{{ overview.fingerprintsFound }}</span>
          </div>
          <div class="overview-item">
            <span class="overview-label">发现漏洞</span>
            <span class="overview-value">{{ overview.vulnerabilitiesFound }}</span>
          </div>
          <div class="overview-item">
            <span class="overview-label">扫描时间</span>
            <span class="overview-value">{{ overview.scanTime }}</span>
          </div>
        </el-card>
      </div>
      
      <!-- 结果概览图表 -->
      <div class="overview-charts" style="margin-top: 30px; display: flex; gap: 20px; flex-wrap: wrap;">
        <el-card class="chart-card" style="flex: 1; min-width: 400px;">
          <template #header>
            <div class="card-header">
              <span>漏洞级别分布</span>
            </div>
          </template>
          <div class="chart-container">
            <canvas ref="levelChart"></canvas>
          </div>
        </el-card>
        <el-card class="chart-card" style="flex: 1; min-width: 400px;">
          <template #header>
            <div class="card-header">
              <span>指纹类型分布</span>
            </div>
          </template>
          <div class="chart-container">
            <canvas ref="fingerprintChart"></canvas>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 结果列表 -->
    <div class="content-section">
      <h3>扫描结果列表</h3>
      <el-table :data="paginatedResults" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="target" label="目标"></el-table-column>
        <el-table-column prop="fingerprint" label="指纹" width="150"></el-table-column>
        <el-table-column prop="vulnerability" label="漏洞"></el-table-column>
        <el-table-column prop="level" label="级别" width="100">
          <template #default="scope">
            <el-tag :type="getLevelType(scope.row.level)">{{ scope.row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scanTime" label="扫描时间" width="200"></el-table-column>
        <el-table-column prop="action" label="操作">
          <template #default="scope">
            <el-button size="small" @click="viewResultDetail(scope.row.id)">查看详情</el-button>
            <el-button size="small" @click="exportResult(scope.row.id)">导出</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="results.length === 0" class="empty-state-wrapper">
        <div class="empty-state">
          <div class="empty-icon">📊</div>
          <p class="empty-text">暂无扫描结果</p>
        </div>
      </div>
      <el-pagination
        v-if="results.length > 0"
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="results.length"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: center; display: flex;"
      />
    </div>

    <!-- 导出选项 -->
    <div class="content-section">
      <h3>导出选项</h3>
      <el-form :inline="true" :model="exportForm">
        <el-form-item label="导出格式">
          <el-select v-model="exportForm.format" placeholder="选择格式">
            <el-option label="CSV" value="csv"></el-option>
            <el-option label="JSON" value="json"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="导出范围">
          <el-select v-model="exportForm.range" placeholder="选择范围">
            <el-option label="全部结果" value="all"></el-option>
            <el-option label="仅漏洞" value="vulnerabilities"></el-option>
            <el-option label="仅指纹" value="fingerprints"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="exportAllResults">导出</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 报告生成 -->
    <div class="content-section">
      <h3>报告生成</h3>
      <el-form :model="reportForm">
        <el-form-item label="报告类型">
          <el-select v-model="reportForm.type" placeholder="选择报告类型">
            <el-option label="详细报告" value="detailed"></el-option>
            <el-option label="摘要报告" value="summary"></el-option>
            <el-option label="漏洞报告" value="vulnerability"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="报告格式">
          <el-select v-model="reportForm.format" placeholder="选择报告格式">
            <el-option label="HTML" value="html"></el-option>
            <el-option label="PDF" value="pdf"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="报告标题">
          <el-input v-model="reportForm.title" placeholder="输入报告标题"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="generateReport">生成报告</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 结果详情对话框 -->
    <el-dialog v-model="resultDetailVisible" title="结果详情" width="60%">
      <div v-if="currentResult">
        <el-descriptions :column="1">
          <el-descriptions-item label="目标">{{ currentResult.target }}</el-descriptions-item>
          <el-descriptions-item label="指纹">{{ currentResult.fingerprint }}</el-descriptions-item>
          <el-descriptions-item label="漏洞">{{ currentResult.vulnerability }}</el-descriptions-item>
          <el-descriptions-item label="级别">{{ currentResult.level }}</el-descriptions-item>
          <el-descriptions-item label="详情">{{ currentResult.detail }}</el-descriptions-item>
          <el-descriptions-item label="扫描时间">{{ currentResult.scanTime }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="resultDetailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { ElMessage } from 'element-plus'

const overview = ref({
  totalTargets: 0,
  fingerprintsFound: 0,
  vulnerabilitiesFound: 0,
  scanTime: ''
})

const results = ref([
  { id: 1, target: '192.168.1.1', fingerprint: 'Apache', vulnerability: 'CVE-2021-41773', level: 'high', scanTime: '2023-10-01 12:00:00', detail: 'Apache HTTP Server 2.4.49 路径穿越漏洞' },
  { id: 2, target: '192.168.1.2', fingerprint: 'Nginx', vulnerability: 'None', level: 'info', scanTime: '2023-10-01 12:01:00', detail: '未发现漏洞' },
  { id: 3, target: '192.168.1.3', fingerprint: 'Tomcat', vulnerability: 'CVE-2020-1938', level: 'medium', scanTime: '2023-10-01 12:02:00', detail: 'Tomcat AJP 协议文件读取漏洞' },
  { id: 4, target: '192.168.1.4', fingerprint: 'MySQL', vulnerability: 'CVE-2019-3733', level: 'low', scanTime: '2023-10-01 12:03:00', detail: 'MySQL 权限提升漏洞' }
])

const exportForm = ref({ format: 'csv', range: 'all' })
const reportForm = ref({ type: 'detailed', format: 'html', title: '扫描报告' })
const resultDetailVisible = ref(false)
const currentResult = ref(null)
const currentPage = ref(1)
const pageSize = ref(20)

const paginatedResults = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return results.value.slice(start, end)
})

const levelChart = ref(null)
const fingerprintChart = ref(null)
let levelChartInstance = null
let fingerprintChartInstance = null

const loadResults = async () => {
  try {
    const tasksData = await window.go.main.App.GetTasks()
    const resultList = []
    let id = 1
    for (const [taskID, taskData] of Object.entries(tasksData)) {
      const result = taskData.result
      if (result && result.fingerprints) {
        for (const fp of result.fingerprints) {
          resultList.push({
            id: id++,
            target: fp.url || fp.target || 'Unknown',
            fingerprint: fp.fingerprint || 'Unknown',
            vulnerability: 'None',
            level: 'info',
            detail: fp.title || '',
            scanTime: result.time || ''
          })
        }
      }
    }
    results.value = resultList
  } catch (e) {
    console.error('Failed to load results:', e)
  }

  try {
    const overviewData = await window.go.main.App.GetResultOverview()
    overview.value = overviewData
  } catch (e) {
    console.error('Failed to load overview:', e)
  }

  nextTick(() => {
    initCharts()
  })
}

const initCharts = () => {
  if (typeof window.Chart === 'undefined') {
    console.warn('Chart.js not loaded, skipping chart initialization')
    return
  }

  if (levelChartInstance) {
    levelChartInstance.destroy()
  }
  if (fingerprintChartInstance) {
    fingerprintChartInstance.destroy()
  }

  const levelData = {
    labels: ['高危', '中危', '低危', '信息'],
    datasets: [{
      label: '漏洞数量',
      data: [
        results.value.filter(r => r.level === 'high').length,
        results.value.filter(r => r.level === 'medium').length,
        results.value.filter(r => r.level === 'low').length,
        results.value.filter(r => r.level === 'info').length
      ],
      backgroundColor: ['#ff3b30', '#ff9500', '#007aff', '#909399']
    }]
  }

  const fingerprintData = {
    labels: [...new Set(results.value.map(r => r.fingerprint))],
    datasets: [{
      label: '出现次数',
      data: [...new Set(results.value.map(r => r.fingerprint))].map(fingerprint => 
        results.value.filter(r => r.fingerprint === fingerprint).length
      ),
      backgroundColor: [
        '#ff3b30', '#ff9500', '#007aff', '#34c759', '#909399',
        '#f08080', '#ffa500', '#87ceeb', '#98fb98', '#d3d3d3'
      ]
    }]
  }

  if (levelChart.value) {
    levelChartInstance = new window.Chart(levelChart.value, {
      type: 'pie',
      data: levelData,
      options: {
        responsive: true,
        maintainAspectRatio: false
      }
    })
  }

  if (fingerprintChart.value) {
    fingerprintChartInstance = new window.Chart(fingerprintChart.value, {
      type: 'bar',
      data: fingerprintData,
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              precision: 0
            }
          }
        }
      }
    })
  }
}

const getLevelType = (level) => {
  switch (level) {
    case 'high':
      return 'danger'
    case 'medium':
      return 'warning'
    case 'low':
      return 'info'
    default:
      return ''
  }
}

const viewResultDetail = (id) => {
  const result = results.value.find(r => r.id === id)
  if (result) {
    currentResult.value = result
    resultDetailVisible.value = true
  }
}

const exportResult = async (id) => {
  try {
    const result = results.value.find(r => r.id === id)
    if (result) {
      const resultData = await window.go.main.App.ExportResult(result.target, exportForm.value.format)
      ElMessage.success('结果已导出')
    }
  } catch (error) {
    ElMessage.error('导出失败: ' + error.message)
  }
}

const exportAllResults = async () => {
  try {
    const result = await window.go.main.App.ExportResults(exportForm.value.format, exportForm.value.range)
    ElMessage.success('结果已导出')
  } catch (error) {
    ElMessage.error('导出失败: ' + error.message)
  }
}

const generateReport = async () => {
  try {
    const result = await window.go.main.App.GenerateReport(reportForm.value.type, reportForm.value.format, reportForm.value.title)
    ElMessage.success('报告已生成')
  } catch (error) {
    ElMessage.error('生成报告失败: ' + error.message)
  }
}

onMounted(() => {
  loadResults()
})

onUnmounted(() => {
  if (levelChartInstance) {
    levelChartInstance.destroy()
  }
  if (fingerprintChartInstance) {
    fingerprintChartInstance.destroy()
  }
})
</script>

<style scoped>
.result-analysis {
  padding: 0;
}

.overview-cards {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.overview-card {
  flex: 1;
  min-width: 200px;
}

.overview-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.overview-label {
  font-size: 14px;
  color: #6e6e73;
}

.overview-value {
  font-size: 18px;
  font-weight: bold;
  color: #007aff;
}

.chart-container {
  height: 300px;
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>