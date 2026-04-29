<template>
  <div class="scan-task">
    <div class="content-section">
      <h2>自定义扫描</h2>
      
      <!-- 任务创建 -->
      <el-form :model="scanForm" label-width="80px">
        <el-form-item label="目标">
          <div style="position: relative;">
            <el-input
              v-model="scanForm.targets"
              type="textarea"
              placeholder="输入目标，一行一个"
              rows="4"
              @blur="detectCurrentTargets"
            ></el-input>
            <div v-if="detectedTargets.length > 0" style="margin-top: 8px;">
              <el-tag 
                v-for="(target, index) in detectedTargets" 
                :key="index"
                :type="getTypeTagType(target.type)"
                size="small"
                style="margin-right: 5px; margin-bottom: 5px;"
              >
                {{ target.label }}: {{ target.value.length > 20 ? target.value.substring(0, 20) + '...' : target.value }}
              </el-tag>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="端口">
          <el-input v-model="scanForm.ports" placeholder="默认 Top1000"></el-input>
        </el-form-item>
        <el-form-item label="扫描类型">
          <el-select v-model="scanForm.scanType" placeholder="选择扫描类型">
            <el-option label="TCP" value="tcp"></el-option>
            <el-option label="SYN" value="syn"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="runScan">开始扫描</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="saveTemplate">保存模板</el-button>
          <el-button @click="loadTemplate">加载模板</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 任务历史 -->
    <div class="content-section">
      <h3>任务历史</h3>
      <div style="display: flex; justify-content: space-between; margin-bottom: 10px;">
        <el-button size="small" @click="loadTasks" :loading="loading">刷新任务</el-button>
        <el-button size="small" @click="clearTasks">清空历史</el-button>
      </div>
      <el-table 
        :data="tasks" 
        style="width: 100%"
        :max-height="500"
        :row-key="row => row.id"
        border
      >
        <el-table-column prop="id" label="任务ID" width="150"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getTagType(scope.row.status)">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="targets" label="目标" min-width="200">
          <template #default="scope">
            <el-tooltip :content="scope.row.targets?.join('\n') || 'N/A'" placement="top">
              <div class="targets-truncate">{{ scope.row.targets?.join(', ') || 'N/A' }}</div>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="time" label="时间" width="200"></el-table-column>
        <el-table-column prop="action" label="操作" width="100">
          <template #default="scope">
            <el-button size="small" @click="viewTask(scope.row.id)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="tasks.length === 0" style="text-align: center; padding: 40px; color: #909399;">
        暂无任务历史
      </div>
    </div>

    <!-- 快捷扫描模板 -->
    <div class="content-section">
      <h3>快捷扫描模板</h3>
      <div class="template-grid">
        <div 
          v-for="template in allTemplates" 
          :key="template.id"
          class="template-card"
          @click="applyTemplate(template)"
        >
          <div class="template-icon">
            <el-icon :size="32">
              <component :is="template.icon" />
            </el-icon>
          </div>
          <h4>{{ template.name }}</h4>
          <p>{{ template.description }}</p>
        </div>
      </div>
    </div>

    <!-- 模板管理 -->
    <div class="content-section">
      <h3>任务模板</h3>
      <el-tag v-for="template in templates" :key="template" style="margin: 5px;">
        {{ template }}
        <el-button size="small" @click="deleteTemplate(template)" style="margin-left: 5px;">删除</el-button>
      </el-tag>
      <div v-if="templates.length === 0" style="text-align: center; padding: 20px; color: #909399;">
        暂无模板
      </div>
    </div>

    <!-- 任务详情对话框 -->
    <el-dialog v-model="taskDetailVisible" title="任务详情" width="60%">
      <div v-if="currentTask">
        <el-descriptions :column="1">
          <el-descriptions-item label="任务ID">{{ currentTask.id }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ currentTask.status }}</el-descriptions-item>
          <el-descriptions-item label="目标">{{ currentTask.targets?.join('\n') || 'N/A' }}</el-descriptions-item>
          <el-descriptions-item label="选项">{{ JSON.stringify(currentTask.options || {}, null, 2) }}</el-descriptions-item>
          <el-descriptions-item label="结果">{{ currentTask.result }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ currentTask.time }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="taskDetailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 保存模板对话框 -->
    <el-dialog v-model="saveTemplateVisible" title="保存模板" width="40%">
      <el-form :model="templateForm">
        <el-form-item label="模板名称">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="saveTemplateVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmSaveTemplate">保存</el-button>
      </template>
    </el-dialog>

    <!-- 加载模板对话框 -->
    <el-dialog v-model="loadTemplateVisible" title="加载模板" width="40%">
      <el-form :model="templateForm">
        <el-form-item label="选择模板">
          <el-select v-model="templateForm.name" placeholder="请选择模板">
            <el-option
              v-for="template in templates"
              :key="template"
              :label="template"
              :value="template"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="loadTemplateVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmLoadTemplate">加载</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, Promotion, Search, Warning, Link } from '@element-plus/icons-vue'
import { detectTargetType, getTypeTagType } from '../utils/targetDetector'
import { getAllTemplates } from '../utils/scanTemplates'

// 所有扫描模板
const allTemplates = ref(getAllTemplates())

// 应用模板
const applyTemplate = (template) => {
  scanForm.value.ports = template.config.ports
  scanForm.value.scanType = template.config.scanType
  ElMessage.success(`已应用模板: ${template.name}`)
}

// 扫描表单
const scanForm = ref({
  targets: '',
  ports: '',
  scanType: 'tcp'
})

// 目标检测相关
const detectedTargets = ref([])

// 检测当前输入的目标
const detectCurrentTargets = () => {
  const targetLines = scanForm.value.targets.split('\n').filter(t => t.trim())
  detectedTargets.value = targetLines.map(t => detectTargetType(t))
}

// 任务管理相关
const tasks = ref([])
const templates = ref([])
const taskDetailVisible = ref(false)
const saveTemplateVisible = ref(false)
const loadTemplateVisible = ref(false)
const currentTask = ref(null)
const templateForm = ref({ name: '' })
const loading = ref(false)

// 加载任务和模板
const loadTasks = async () => {
  loading.value = true
  try {
    const [taskData, templateData] = await Promise.all([
      window.go.main.App.GetTasks(),
      window.go.main.App.GetTaskTemplates()
    ])
    
    const taskList = []
    for (const [id, data] of Object.entries(taskData)) {
      taskList.push({
        id: id,
        status: data.status,
        targets: data.result?.targets || [],
        time: data.result?.time || '',
        options: data.result?.options,
        result: data.result?.result
      })
    }
    taskList.sort((a, b) => new Date(b.time) - new Date(a.time))
    tasks.value = taskList
    
    templates.value = templateData
  } catch (e) {
    console.error('Failed to load data:', e)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 清空任务历史
const clearTasks = async () => {
  try {
    tasks.value = []
    ElMessage.success('任务历史已清空')
  } catch (error) {
    ElMessage.error('清空任务失败: ' + error.message)
  }
}

// 获取标签类型
const getTagType = (status) => {
  switch (status) {
    case 'running':
      return 'warning'
    case 'completed':
      return 'success'
    case 'failed':
      return 'danger'
    default:
      return ''
  }
}

// 运行扫描
const runScan = async () => {
  const targets = scanForm.value.targets.split('\n').filter(t => t.trim())
  if (targets.length === 0) {
    ElMessage.warning('请输入扫描目标')
    return
  }

  try {
    const taskID = 'task_' + Date.now()
    const options = {
      ports: scanForm.value.ports || 'Top1000',
      scanType: scanForm.value.scanType
    }
    const result = await window.go.main.App.RunScanWithID(taskID, targets, options)
    ElMessage.success(result)
    await loadTasks()
  } catch (error) {
    ElMessage.error('扫描失败: ' + error.message)
  }
}

// 查看任务详情
const viewTask = (taskID) => {
  const task = tasks.value.find(t => t.id === taskID)
  if (task) {
    currentTask.value = task
    taskDetailVisible.value = true
  }
}

// 保存模板
const saveTemplate = () => {
  saveTemplateVisible.value = true
}

// 确认保存模板
const confirmSaveTemplate = async () => {
  if (!templateForm.value.name) {
    ElMessage.warning('请输入模板名称')
    return
  }

  try {
    const config = {
      targets: scanForm.value.targets,
      ports: scanForm.value.ports,
      scanType: scanForm.value.scanType
    }
    const result = await window.go.main.App.SaveTaskTemplate(templateForm.value.name, config)
    ElMessage.success(result)
    await loadTasks()
    saveTemplateVisible.value = false
    templateForm.value.name = ''
  } catch (error) {
    ElMessage.error('保存模板失败: ' + error.message)
  }
}

// 加载模板
const loadTemplate = () => {
  loadTemplateVisible.value = true
}

// 确认加载模板
const confirmLoadTemplate = async () => {
  if (!templateForm.value.name) {
    ElMessage.warning('请选择模板')
    return
  }

  try {
    const template = await window.go.main.App.LoadTaskTemplate(templateForm.value.name)
    if (template) {
      scanForm.value.targets = template.targets || ''
      scanForm.value.ports = template.ports || ''
      scanForm.value.scanType = template.scanType || 'tcp'
      ElMessage.success('模板加载成功')
    }
    loadTemplateVisible.value = false
    templateForm.value.name = ''
  } catch (error) {
    ElMessage.error('加载模板失败: ' + error.message)
  }
}

// 删除模板
const deleteTemplate = async (name) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除此模板吗？',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    const result = await window.go.main.App.DeleteTaskTemplate(name)
    ElMessage.success(result)
    await loadTasks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + error)
    }
  }
}

// 重置表单
const resetForm = () => {
  scanForm.value = {
    targets: '',
    ports: '',
    scanType: 'tcp'
  }
}

// 页面加载时获取数据
onMounted(() => {
  loadTasks()
})
</script>

<style scoped>
.scan-task {
  padding: 0;
}

/* 模板网格样式 */
.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.template-card {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: white;
}

.template-card:hover {
  border-color: var(--el-color-primary);
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.template-icon {
  color: var(--el-color-primary);
  margin-bottom: 12px;
}

.template-card h4 {
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
  font-size: 16px;
}

.template-card p {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  line-height: 1.5;
}

/* 目标列表截断样式 */
.targets-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

/* 表格行悬停效果 */
:deep(.el-table__row) {
  transition: all 0.3s ease;
}

:deep(.el-table__row:hover) {
  background-color: rgba(245, 158, 11, 0.05);
}
</style>