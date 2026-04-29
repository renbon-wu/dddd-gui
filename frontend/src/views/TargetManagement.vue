<template>
  <div class="target-management">
    <div class="content-section">
      <h2>目标管理</h2>
      
      <!-- 操作按钮 -->
      <div class="target-actions" style="margin-bottom: 20px; display: flex; align-items: center; gap: 10px; flex-wrap: wrap;">
        <el-button type="primary" @click="importTargets">导入目标</el-button>
        <el-button @click="batchScan">批量扫描</el-button>
        <el-button type="danger" @click="batchDelete">批量删除</el-button>
        <el-button @click="createGroup">创建分组</el-button>
        <SearchBar @search="handleSearch" />
      </div>

      <!-- 分组管理 -->
      <div class="group-section">
        <h3>目标分组</h3>
        <div class="group-cards">
          <el-card v-for="(groupTargets, groupName) in groups" :key="groupName" class="group-card">
            <template #header>
              <div class="group-header">
                <span>{{ groupName }} ({{ groupTargets.length }}个目标)</span>
                <div>
                  <el-button size="small" @click="scanGroup(groupName)">扫描</el-button>
                  <el-button size="small" type="danger" @click="deleteGroup(groupName)">删除</el-button>
                </div>
              </div>
            </template>
            <div class="group-targets">
              <el-tag v-for="target in groupTargets" :key="target" class="target-tag">
                {{ target.length > 20 ? target.substring(0, 20) + '...' : target }}
                <el-button size="small" @click="removeFromGroup(groupName, target)">移除</el-button>
              </el-tag>
            </div>
            <div v-if="groupTargets.length === 0" class="empty-group">
              分组为空
            </div>
          </el-card>
        </div>
        <div v-if="Object.keys(groups).length === 0" class="no-group">
          暂无分组，点击上方"创建分组"按钮创建
        </div>
      </div>

      <!-- 目标列表 -->
      <div class="target-list-section">
        <h3>目标列表</h3>
        <el-table 
          :data="paginatedTargets" 
          style="width: 100%" 
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55"></el-table-column>
          <el-table-column prop="id" label="ID" width="80"></el-table-column>
          <el-table-column prop="target" label="目标"></el-table-column>
          <el-table-column prop="status" label="状态" width="100"></el-table-column>
          <el-table-column prop="time" label="添加时间" width="200"></el-table-column>
          <el-table-column prop="scanTime" label="扫描时间" width="200"></el-table-column>
          <el-table-column prop="action" label="操作">
            <template #default="scope">
              <el-button size="small" @click="editTarget(scope.row)">编辑</el-button>
              <el-button size="small" @click="scanTarget(scope.row.target)">扫描</el-button>
              <el-button size="small" type="danger" @click="deleteTarget(scope.row.target)">删除</el-button>
              <el-button size="small" @click="addToGroup(scope.row.target)">加入分组</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-if="filteredTargets.length === 0" class="no-target">
          <div class="empty-state">
            <div class="empty-icon">📋</div>
            <p class="empty-text">暂无目标</p>
          </div>
        </div>
        <el-pagination
          v-if="filteredTargets.length > 0"
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="filteredTargets.length"
          layout="total, sizes, prev, pager, next, jumper"
          style="margin-top: 20px; justify-content: center; display: flex;"
        />
      </div>
    </div>

    <!-- 导入目标对话框 -->
    <el-dialog v-model="importDialogVisible" title="导入目标" width="50%">
      <el-form>
        <el-form-item label="导入方式">
          <el-radio-group v-model="importMode">
            <el-radio label="file">文件导入</el-radio>
            <el-radio label="manual">手动输入</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="importMode === 'file'" label="选择文件">
          <el-upload
            class="upload-demo"
            action=""
            :auto-upload="false"
            :on-change="handleFileChange"
            accept=".txt,.csv"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持 TXT（一行一个目标）或 CSV（逗号分隔）格式
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item v-if="importMode === 'manual'" label="目标列表">
          <el-input
            v-model="manualTargets"
            type="textarea"
            placeholder="输入目标，一行一个"
            rows="5"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmImport">导入</el-button>
      </template>
    </el-dialog>

    <!-- 创建分组对话框 -->
    <el-dialog v-model="createGroupVisible" title="创建分组" width="40%">
      <el-form :model="groupForm">
        <el-form-item label="分组名称">
          <el-input v-model="groupForm.name" placeholder="请输入分组名称"></el-input>
        </el-form-item>
        <el-form-item label="选择目标">
          <el-select v-model="groupForm.targets" multiple placeholder="请选择目标">
            <el-option
              v-for="target in targets"
              :key="target.target"
              :label="target.target"
              :value="target.target"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createGroupVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreateGroup">创建</el-button>
      </template>
    </el-dialog>

    <!-- 编辑目标对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑目标" width="40%">
      <el-form :model="editForm">
        <el-form-item label="目标地址">
          <el-input v-model="editForm.target" placeholder="请输入目标地址"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEdit">保存</el-button>
      </template>
    </el-dialog>

    <!-- 添加到分组对话框 -->
    <el-dialog v-model="addToGroupVisible" title="添加到分组" width="40%">
      <el-form :model="addToGroupForm">
        <el-form-item label="选择分组">
          <el-select v-model="addToGroupForm.groupName" placeholder="请选择分组">
            <el-option
              v-for="groupName in Object.keys(groups)"
              :key="groupName"
              :label="groupName"
              :value="groupName"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addToGroupVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddToGroup">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SearchBar from '../components/SearchBar.vue'

// 目标列表
const targets = ref([
  { id: 1, target: '192.168.1.1', status: '已扫描' },
  { id: 2, target: '192.168.1.2', status: '未扫描' }
])

// 目标管理相关
const selectedTargets = ref([])
const groups = ref({})
const importDialogVisible = ref(false)
const createGroupVisible = ref(false)
const addToGroupVisible = ref(false)
const importMode = ref('file')
const manualTargets = ref('')
const selectedFile = ref(null)
const groupForm = ref({ name: '', targets: [] })
const addToGroupForm = ref({ groupName: '', target: '' })
const editDialogVisible = ref(false)
const editForm = ref({ oldTarget: '', target: '' })
const searchText = ref('')
const searchField = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)

// 过滤目标
const filteredTargets = computed(() => {
  if (!searchText.value) {
    return targets.value
  }
  const searchLower = searchText.value.toLowerCase()
  return targets.value.filter(t => {
    if (searchField.value === 'all') {
      return (
        t.target.toLowerCase().includes(searchLower) ||
        t.status.toLowerCase().includes(searchLower)
      )
    } else if (searchField.value === 'name') {
      return t.target.toLowerCase().includes(searchLower)
    } else if (searchField.value === 'desc') {
      return t.status.toLowerCase().includes(searchLower)
    }
    return true
  })
})

// 分页目标
const paginatedTargets = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTargets.value.slice(start, end)
})

const handleSearch = (searchData) => {
  searchText.value = searchData.text
  searchField.value = searchData.field
}

// 加载目标和分组
const loadTargets = async () => {
  try {
    const targetData = await window.go.main.App.GetTargets()
    const targetList = []
    let id = 1
    for (const [target, data] of Object.entries(targetData)) {
      targetList.push({
        id: id,
        target: target,
        status: data.status,
        time: data.time,
        scanTime: data.scanTime
      })
      id++
    }
    targets.value = targetList
  } catch (e) {
    console.error('Failed to load targets:', e)
    ElMessage.error('加载目标列表失败')
  }

  try {
    const groupData = await window.go.main.App.GetTargetGroups()
    groups.value = groupData
  } catch (e) {
    console.error('Failed to load groups:', e)
    ElMessage.error('加载目标分组失败')
  }
}

// 处理文件选择
const handleFileChange = (file) => {
  selectedFile.value = file
}

// 导入目标
const importTargets = () => {
  importDialogVisible.value = true
  importMode.value = 'file'
  manualTargets.value = ''
  selectedFile.value = null
}

// 确认导入
const confirmImport = async () => {
  if (importMode.value === 'file' && !selectedFile.value) {
    ElMessage.warning('请选择文件')
    return
  }

  if (importMode.value === 'manual' && !manualTargets.value) {
    ElMessage.warning('请输入目标')
    return
  }

  try {
    let result
    if (importMode.value === 'file') {
      result = await window.go.main.App.ImportTargets(selectedFile.value.raw.path)
    } else {
      result = await window.go.main.App.ImportTargetsFromText(manualTargets.value)
    }
    ElMessage.success(result)
    await loadTargets()
    importDialogVisible.value = false
    manualTargets.value = ''
    selectedFile.value = null
  } catch (error) {
    ElMessage.error('导入目标失败: ' + error.message)
  }
}

// 处理选择变化
const handleSelectionChange = (val) => {
  selectedTargets.value = val.map(item => item.target)
}

// 批量扫描
const batchScan = async () => {
  if (selectedTargets.value.length === 0) {
    ElMessage.warning('请选择目标')
    return
  }

  try {
    const result = await window.go.main.App.BatchScanTargets(selectedTargets.value)
    ElMessage.success(result)
    await loadTargets()
  } catch (error) {
    ElMessage.error('批量扫描失败: ' + error.message)
  }
}

// 批量删除
const batchDelete = async () => {
  if (selectedTargets.value.length === 0) {
    ElMessage.warning('请选择目标')
    return
  }

  try {
    const result = await window.go.main.App.BatchDeleteTargets(selectedTargets.value)
    ElMessage.success(result)
    await loadTargets()
  } catch (error) {
    ElMessage.error('批量删除失败: ' + error.message)
  }
}

// 扫描单个目标
const scanTarget = async (target) => {
  try {
    const result = await window.go.main.App.BatchScanTargets([target])
    ElMessage.success(result)
    await loadTargets()
  } catch (error) {
    ElMessage.error('扫描失败: ' + error.message)
  }
}

// 删除单个目标
const deleteTarget = async (target) => {
  try {
    const result = await window.go.main.App.BatchDeleteTargets([target])
    ElMessage.success(result)
    await loadTargets()
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

// 创建分组
const createGroup = () => {
  createGroupVisible.value = true
  groupForm.value = { name: '', targets: [] }
}

// 确认创建分组
const confirmCreateGroup = async () => {
  if (!groupForm.value.name) {
    ElMessage.warning('请输入分组名称')
    return
  }
  if (groupForm.value.targets.length === 0) {
    ElMessage.warning('请选择目标')
    return
  }

  try {
    const result = await window.go.main.App.CreateTargetGroup(groupForm.value.name, groupForm.value.targets)
    ElMessage.success(result)
    await loadTargets()
    createGroupVisible.value = false
  } catch (error) {
    ElMessage.error('创建分组失败: ' + error.message)
  }
}

// 添加到分组
const addToGroup = (target) => {
  addToGroupVisible.value = true
  addToGroupForm.value = { groupName: '', target: target }
}

// 确认添加到分组
const confirmAddToGroup = async () => {
  if (!addToGroupForm.value.groupName) {
    ElMessage.warning('请选择分组')
    return
  }

  try {
    const result = await window.go.main.App.AddTargetToGroup(addToGroupForm.value.groupName, addToGroupForm.value.target)
    ElMessage.success(result)
    await loadTargets()
    addToGroupVisible.value = false
  } catch (error) {
    ElMessage.error('添加到分组失败: ' + error.message)
  }
}

// 从分组移除
const removeFromGroup = async (groupName, target) => {
  try {
    const result = await window.go.main.App.RemoveTargetFromGroup(groupName, target)
    ElMessage.success(result)
    await loadTargets()
  } catch (error) {
    ElMessage.error('移除失败: ' + error.message)
  }
}

// 扫描分组
const scanGroup = async (groupName) => {
  const groupTargets = groups.value[groupName]
  if (!groupTargets || groupTargets.length === 0) {
    ElMessage.warning('分组内没有目标')
    return
  }

  try {
    const result = await window.go.main.App.BatchScanTargets(groupTargets)
    ElMessage.success(result)
    await loadTargets()
  } catch (error) {
    ElMessage.error('扫描分组失败: ' + error.message)
  }
}

// 删除分组
const deleteGroup = async (groupName) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分组 "${groupName}" 吗？此操作不可撤销。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    delete groups.value[groupName]
    ElMessage.success('分组已删除')
  } catch (error) {
    ElMessage.info('已取消删除')
  }
}

// 编辑目标
const editTarget = (target) => {
  editForm.value = {
    oldTarget: target,
    target: target
  }
  editDialogVisible.value = true
}

// 确认编辑
const confirmEdit = async () => {
  if (!editForm.value.target.trim()) {
    ElMessage.warning('请输入目标地址')
    return
  }
  if (editForm.value.oldTarget === editForm.value.target) {
    editDialogVisible.value = false
    return
  }

  try {
    const result = await window.go.main.App.EditTarget(editForm.value.oldTarget, editForm.value.target)
    ElMessage.success(result)
    await loadTargets()
    editDialogVisible.value = false
  } catch (error) {
    ElMessage.error('编辑失败: ' + error.message)
  }
}

// 页面加载时获取数据
onMounted(() => {
  loadTargets()
})
</script>

<style scoped>
.target-management {
  padding: 0;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.group-header > div {
  display: flex;
  gap: 8px;
}

.group-section {
  margin-bottom: 30px;
}

.group-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.group-card {
  min-width: 300px;
}

.group-targets {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.target-tag {
  max-width: calc(100% - 60px);
}

.empty-group,
.no-group,
.no-target {
  text-align: center;
  padding: 20px;
  color: #909399;
}

.no-group,
.no-target {
  padding: 40px;
}

.target-list-section {
  margin-top: 20px;
}
</style>