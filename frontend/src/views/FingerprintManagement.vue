<template>
  <div class="fingerprint-management">
    <div class="content-section">
      <h2>指纹管理</h2>
      
      <div class="fingerprint-actions" style="margin-bottom: 20px; display: flex; align-items: center; gap: 10px; flex-wrap: wrap;">
        <el-button type="primary" @click="uploadFingerprint">上传指纹文件</el-button>
        <el-button @click="testFingerprint">测试指纹</el-button>
        <el-button @click="addFingerprint">添加指纹</el-button>
        <el-select v-model="activeCategory" placeholder="选择分类" style="width: 150px;">
          <el-option label="全部" value="all"></el-option>
          <el-option v-for="category in categories" :key="category" :label="category" :value="category"></el-option>
        </el-select>
        <SearchBar @search="handleSearch" />
      </div>

      <el-table :data="paginatedFingerprints" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="指纹名称"></el-table-column>
        <el-table-column prop="rules" label="规则" width="400">
          <template #default="scope">
            <div style="max-height: 60px; overflow: hidden; text-overflow: ellipsis;">
              {{ scope.row.rules.join(', ') }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作">
          <template #default="scope">
            <el-button size="small" @click="editFingerprint(scope.row)">编辑</el-button>
            <el-button size="small" @click="testFingerprintByName(scope.row.name)">测试</el-button>
            <el-button size="small" type="danger" @click="deleteFingerprint(scope.row.name)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="filteredFingerprints.length === 0" class="empty-state-wrapper">
        <div class="empty-state">
          <div class="empty-icon">🔒</div>
          <p class="empty-text">暂无指纹数据</p>
        </div>
      </div>
      <el-pagination
        v-if="filteredFingerprints.length > 0"
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="filteredFingerprints.length"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: center; display: flex;"
      />
    </div>

    <el-dialog v-model="uploadDialogVisible" title="上传指纹文件" width="50%">
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
                支持 YAML 格式的指纹文件
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

    <el-dialog v-model="editDialogVisible" :title="isEdit ? '编辑指纹' : '添加指纹'" width="60%">
      <el-form :model="fingerprintForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="fingerprintForm.name" placeholder="请输入指纹名称"></el-input>
        </el-form-item>
        <el-form-item label="规则">
          <el-input
            v-model="fingerprintForm.rulesText"
            type="textarea"
            placeholder="输入指纹规则，多个规则用逗号分隔"
            rows="10"
          ></el-input>
          <div style="font-size: 12px; color: #909399; margin-top: 5px;">
            示例: body="example.com", title="管理系统", header="Server: Apache"
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveFingerprint">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="testDialogVisible" title="测试指纹" width="50%">
      <el-form :model="testForm" label-width="100px">
        <el-form-item label="测试模式">
          <el-radio-group v-model="testForm.mode">
            <el-radio label="format">格式检测</el-radio>
            <el-radio label="real">真实检测</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="指纹名称">
          <el-select v-model="testForm.fingerprint" placeholder="选择指纹">
            <el-option
              v-for="fingerprint in fingerprints"
              :key="fingerprint.id"
              :label="fingerprint.name"
              :value="fingerprint.name"
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
import { ElMessage } from 'element-plus'
import SearchBar from '../components/SearchBar.vue'

const fingerprints = ref([])
const categories = ref(['全部'])
const activeCategory = ref('all')
const uploadDialogVisible = ref(false)
const editDialogVisible = ref(false)
const testDialogVisible = ref(false)
const isEdit = ref(false)
const originalFingerprintName = ref('')
const originalName = ref('')
const selectedFile = ref(null)
const searchText = ref('')
const searchField = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)

const fingerprintForm = ref({
  name: '',
  rulesText: ''
})

const testForm = ref({
  mode: 'format',
  fingerprint: '',
  target: ''
})

const filteredFingerprints = computed(() => {
  let result = fingerprints.value
  // 按分类过滤
  if (activeCategory.value !== 'all') {
    result = result.filter(f => f.category === activeCategory.value)
  }
  // 按搜索条件过滤
  if (searchText.value) {
    const searchLower = searchText.value.toLowerCase()
    result = result.filter(f => {
      if (searchField.value === 'all') {
        return (
          f.name.toLowerCase().includes(searchLower) ||
          f.rules.join(', ').toLowerCase().includes(searchLower)
        )
      } else if (searchField.value === 'name') {
        return f.name.toLowerCase().includes(searchLower)
      } else if (searchField.value === 'desc') {
        return f.rules.join(', ').toLowerCase().includes(searchLower)
      }
      return true
    })
  }
  return result
})

const paginatedFingerprints = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredFingerprints.value.slice(start, end)
})

const handleSearch = (searchData) => {
  searchText.value = searchData.text
  searchField.value = searchData.field
}

const loadFingerprints = async () => {
  try {
    const fingerData = await window.go.main.App.GetAllFingerprintDetails()
    const fingerList = []
    let id = 1
    if (Array.isArray(fingerData)) {
      for (const item of fingerData) {
        if (item && item.name) {
          const ruleArray = Array.isArray(item.rules) ? item.rules : []
          fingerList.push({
            id: id,
            name: item.name,
            rules: ruleArray,
            category: 'default'
          })
          id++
        }
      }
    }
    fingerprints.value = fingerList
    
    const uniqueCategories = [...new Set(fingerprints.value.map(f => f.category))]
    categories.value = ['全部', ...uniqueCategories]
  } catch (e) {
    console.error('Failed to load fingerprints:', e)
  }
}

const handleFileChange = (file) => {
  selectedFile.value = file
}

const uploadFingerprint = () => {
  uploadDialogVisible.value = true
  selectedFile.value = null
}

const confirmUpload = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请选择文件')
    return
  }

  try {
    const result = await window.go.main.App.UploadFingerprint(selectedFile.value.raw.path)
    ElMessage.success(result)
    await loadFingerprints()
    uploadDialogVisible.value = false
  } catch (error) {
    ElMessage.error('上传失败: ' + error.message)
  }
}

const addFingerprint = () => {
  isEdit.value = false
  fingerprintForm.value = {
    name: '',
    rulesText: ''
  }
  editDialogVisible.value = true
}

const editFingerprint = (fingerprint) => {
  isEdit.value = true
  originalFingerprintName.value = fingerprint.name
  fingerprintForm.value = {
    name: fingerprint.name,
    rulesText: fingerprint.rules.join(', ')
  }
  editDialogVisible.value = true
}

const saveFingerprint = async () => {
  if (!fingerprintForm.value.name) {
    ElMessage.warning('请输入指纹名称')
    return
  }
  if (!fingerprintForm.value.rulesText) {
    ElMessage.warning('请输入指纹规则')
    return
  }

  const rules = fingerprintForm.value.rulesText.split(',').map(r => r.trim()).filter(r => r)

  try {
    let result
    if (isEdit.value) {
      result = await window.go.main.App.EditFingerprint(originalFingerprintName.value, fingerprintForm.value.name, rules)
    } else {
      result = await window.go.main.App.AddFingerprint(fingerprintForm.value.name, rules)
    }
    ElMessage.success(result)
    await loadFingerprints()
    editDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message)
  }
}

const deleteFingerprint = async (name) => {
  try {
    const result = await window.go.main.App.DeleteFingerprint(name)
    ElMessage.success(result)
    await loadFingerprints()
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

const testFingerprint = () => {
  testForm.value = {
    mode: 'format',
    fingerprint: '',
    target: ''
  }
  testDialogVisible.value = true
}

const testFingerprintByName = (name) => {
  testForm.value = {
    mode: 'format',
    fingerprint: name,
    target: ''
  }
  testDialogVisible.value = true
}

const runTest = async () => {
  if (!testForm.value.fingerprint) {
    ElMessage.warning('请选择指纹')
    return
  }
  if (testForm.value.mode === 'real' && !testForm.value.target) {
    ElMessage.warning('请输入测试目标')
    return
  }

  try {
    let result
    if (testForm.value.mode === 'format') {
      result = await window.go.main.App.ValidateFingerprintByName(testForm.value.fingerprint)
    } else {
      result = await window.go.main.App.TestFingerprint(testForm.value.fingerprint, testForm.value.target)
    }
    ElMessage.success(result)
    testDialogVisible.value = false
  } catch (error) {
    ElMessage.error('测试失败: ' + error.message)
  }
}

onMounted(() => {
  loadFingerprints()
})
</script>

<style scoped>
.fingerprint-management {
  padding: 0;
}
</style>