<template>
  <div class="poc-management">
    <div class="content-section">
      <h2>POC管理</h2>
      
      <div class="poc-actions" style="margin-bottom: 20px; display: flex; align-items: center; gap: 10px; flex-wrap: wrap;">
        <el-button type="primary" @click="uploadPOC">上传POC</el-button>
        <el-button @click="testPOC">测试POC</el-button>
        <el-button @click="showWorkflowDialog">工作流管理</el-button>
        <el-select v-model="activeCategory" placeholder="选择分类" style="width: 150px;">
          <el-option label="全部" value="all"></el-option>
          <el-option v-for="category in categories" :key="category" :label="category" :value="category"></el-option>
        </el-select>
        <SearchBar @search="handleSearch" />
      </div>

      <el-table :data="paginatedPOCs" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="POC名称"></el-table-column>
        <el-table-column prop="action" label="操作">
          <template #default="scope">
            <el-button size="small" @click="testPOC(scope.row.name)">测试</el-button>
            <el-button size="small" type="danger" @click="deletePOC(scope.row.name)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="filteredPOCs.length === 0" class="empty-state-wrapper">
        <div class="empty-state">
          <div class="empty-icon">⚠️</div>
          <p class="empty-text">暂无POC数据</p>
        </div>
      </div>
      <el-pagination
        v-if="filteredPOCs.length > 0"
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="filteredPOCs.length"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: center; display: flex;"
      />
    </div>

    <el-dialog v-model="workflowDialogVisible" title="工作流管理" width="60%">
      <el-form :model="workflowForm" label-width="120px">
        <el-form-item label="选择指纹">
          <el-select v-model="selectedFingerprint" placeholder="请选择指纹" @change="loadWorkflowDetails">
            <el-option
              v-for="fingerprint in fingerprints"
              :key="fingerprint.id"
              :label="fingerprint.name"
              :value="fingerprint.name"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="执行类型">
          <el-checkbox-group v-model="selectedTypes">
            <el-checkbox label="root">根目录执行</el-checkbox>
            <el-checkbox label="dir">目录执行</el-checkbox>
            <el-checkbox label="base">基础路径执行</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="选择POC">
          <el-select v-model="selectedPocs" multiple placeholder="请选择POC">
            <el-option
              v-for="poc in pocs"
              :key="poc.id"
              :label="poc.name"
              :value="poc.name"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="当前配置">
          <div v-if="currentWorkflow">
            <div><strong>执行类型:</strong> {{ currentWorkflow.type.join(', ') }}</div>
            <div><strong>POC列表:</strong> {{ currentWorkflow.pocs.join(', ') }}</div>
          </div>
          <div v-else>暂无配置</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="workflowDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="deleteWorkflow" v-if="selectedFingerprint">删除工作流</el-button>
        <el-button type="primary" @click="saveWorkflow">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="uploadDialogVisible" title="上传POC" width="50%">
      <el-form>
        <el-form-item label="选择文件">
          <el-upload
            class="upload-demo"
            action=""
            :auto-upload="false"
            :on-change="handleFileChange"
            accept=".yaml,.yml">
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持 YAML 格式的 POC 文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmUpload">上传</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="testDialogVisible" title="测试POC" width="50%">
      <el-form :model="testForm" label-width="100px">
        <el-form-item label="测试模式">
          <el-radio-group v-model="testForm.mode">
            <el-radio label="format">格式检测</el-radio>
            <el-radio label="real">真实检测</el-radio>
          </el-radio-group>
        </el-form-item>
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
        <el-form-item v-if="testForm.mode === 'real'" label="测试目标">
          <el-input v-model="testForm.target" placeholder="输入测试目标URL"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="runTest">测试</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SearchBar from '../components/SearchBar.vue'

const pocs = ref([])
const fingerprints = ref([])
const categories = ref(['RCE', 'XSS', 'SQLi', 'XXE'])
const workflows = ref({})
const workflowDialogVisible = ref(false)
const uploadDialogVisible = ref(false)
const testDialogVisible = ref(false)
const workflowForm = ref({})
const testForm = ref({ mode: 'format', poc: '', target: '' })
const selectedFile = ref(null)
const activeCategory = ref('all')
const selectedFingerprint = ref('')
const selectedPocs = ref([])
const selectedTypes = ref(['root'])
const currentWorkflow = ref(null)
const searchText = ref('')
const searchField = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)

const filteredPOCs = computed(() => {
  let result = pocs.value
  // 按分类过滤
  if (activeCategory.value !== 'all') {
    result = result.filter(p => p.category === activeCategory.value)
  }
  // 按搜索条件过滤
  if (searchText.value) {
    const searchLower = searchText.value.toLowerCase()
    result = result.filter(p => {
      if (searchField.value === 'all') {
        return (
          p.name.toLowerCase().includes(searchLower) ||
          p.type.toLowerCase().includes(searchLower)
        )
      } else if (searchField.value === 'name') {
        return p.name.toLowerCase().includes(searchLower)
      } else if (searchField.value === 'desc') {
        return p.type.toLowerCase().includes(searchLower)
      }
      return true
    })
  }
  return result
})

const paginatedPOCs = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredPOCs.value.slice(start, end)
})

const handleSearch = (searchData) => {
  searchText.value = searchData.text
  searchField.value = searchData.field
}

const loadData = async () => {
  try {
    const pocData = await window.go.main.App.GetPOCs()
    pocs.value = pocData.map((name, index) => ({
      id: index + 1,
      name: name,
      category: 'RCE',
      type: 'Nuclei',
      source: 'GitHub',
      time: new Date().toLocaleString()
    }))
  } catch (e) {
    console.error('Failed to load POCs:', e)
  }

  try {
    const fingerData = await window.go.main.App.GetFingerprints()
    fingerprints.value = fingerData.map((name, index) => ({
      id: index + 1,
      name: name
    }))
  } catch (e) {
    console.error('Failed to load fingerprints:', e)
  }

  try {
    const workflowData = await window.go.main.App.GetWorkflows()
    workflows.value = workflowData
  } catch (e) {
    console.error('Failed to load workflows:', e)
  }
}

const loadWorkflowDetails = async () => {
  if (!selectedFingerprint.value) {
    currentWorkflow.value = null
    selectedTypes.value = ['root']
    selectedPocs.value = []
    return
  }

  try {
    const details = await window.go.main.App.GetWorkflowDetails(selectedFingerprint.value)
    currentWorkflow.value = details
    selectedTypes.value = details.type || ['root']
    selectedPocs.value = details.pocs || []
  } catch (e) {
    currentWorkflow.value = null
    selectedTypes.value = ['root']
    selectedPocs.value = []
  }
}

const handleFileChange = (file) => {
  selectedFile.value = file
}

const uploadPOC = () => {
  uploadDialogVisible.value = true
  selectedFile.value = null
}

const confirmUpload = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请选择文件')
    return
  }

  try {
    const result = await window.go.main.App.UploadPOC(selectedFile.value.raw.path)
    ElMessage.success(result)
    await loadData()
    uploadDialogVisible.value = false
  } catch (error) {
    ElMessage.error('上传POC失败: ' + error)
  }
}

const testPOC = (name) => {
  testForm.value = { mode: 'format', poc: name || '', target: '' }
  testDialogVisible.value = true
}

const runTest = async () => {
  if (!testForm.value.poc) {
    ElMessage.warning('请选择POC')
    return
  }
  if (testForm.value.mode === 'real' && !testForm.value.target) {
    ElMessage.warning('请输入测试目标')
    return
  }

  try {
    let result
    if (testForm.value.mode === 'format') {
      result = await window.go.main.App.ValidatePOC(testForm.value.poc)
    } else {
      result = await window.go.main.App.TestPOC(testForm.value.poc, testForm.value.target)
    }
    ElMessage.success(result)
    testDialogVisible.value = false
  } catch (error) {
    ElMessage.error('测试失败: ' + error)
  }
}

const deletePOC = async (name) => {
  try {
    await ElMessageBox.confirm('确定要删除该POC吗?', '确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const result = await window.go.main.App.DeletePOC(name)
    ElMessage.success(result)
    await loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除POC失败: ' + error)
    }
  }
}

const showWorkflowDialog = () => {
  selectedFingerprint.value = ''
  selectedPocs.value = []
  selectedTypes.value = ['root']
  currentWorkflow.value = null
  workflowDialogVisible.value = true
}

const saveWorkflow = async () => {
  if (!selectedFingerprint.value) {
    ElMessage.warning('请选择指纹')
    return
  }
  if (selectedPocs.value.length === 0) {
    ElMessage.warning('请选择至少一个POC')
    return
  }
  if (selectedTypes.value.length === 0) {
    ElMessage.warning('请选择至少一个执行类型')
    return
  }

  try {
    const result = await window.go.main.App.UpdateWorkflow(
      selectedFingerprint.value,
      selectedPocs.value,
      selectedTypes.value
    )
    ElMessage.success(result)
    await loadData()
    await loadWorkflowDetails()
  } catch (error) {
    ElMessage.error('更新工作流失败: ' + error)
  }
}

const deleteWorkflow = async () => {
  if (!selectedFingerprint.value) {
    ElMessage.warning('请选择指纹')
    return
  }

  try {
    await ElMessageBox.confirm('确定要删除该工作流吗?', '确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const result = await window.go.main.App.DeleteWorkflow(selectedFingerprint.value)
    ElMessage.success(result)
    selectedFingerprint.value = ''
    selectedPocs.value = []
    selectedTypes.value = ['root']
    currentWorkflow.value = null
    await loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除工作流失败: ' + error)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.poc-management {
  padding: 0;
}

.poc-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}
</style>