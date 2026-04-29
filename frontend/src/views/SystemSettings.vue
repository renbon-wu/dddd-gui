<template>
  <div class="system-settings">
    <!-- 全局配置 -->
    <div class="content-section">
      <h2>全局配置</h2>
      <el-form :model="globalConfig" label-width="120px">
        <el-form-item label="HTTP代理">
          <el-input v-model="globalConfig.httpProxy" placeholder="例如: http://127.0.0.1:8080"></el-input>
        </el-form-item>
        <el-form-item label="线程数">
          <el-input-number v-model="globalConfig.threads" :min="1" :max="100" placeholder="线程数"></el-input-number>
        </el-form-item>
        <el-form-item label="超时时间">
          <el-input-number v-model="globalConfig.timeout" :min="1" :max="30" placeholder="超时时间(秒)"></el-input-number>
        </el-form-item>
        <el-form-item label="重试次数">
          <el-input-number v-model="globalConfig.retries" :min="0" :max="5" placeholder="重试次数"></el-input-number>
        </el-form-item>
        <el-form-item label="子域名字典">
          <el-input v-model="globalConfig.subdomainDict" placeholder="子域名字典文件路径"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveGlobalConfig">保存</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- API配置 -->
    <div class="content-section">
      <h2>API配置</h2>
      <el-form :model="apiConfig" label-width="120px">
        <el-form-item label="Fofa API Key">
          <el-input v-model="apiConfig.fofaKey" placeholder="Fofa API Key"></el-input>
        </el-form-item>
        <el-form-item label="Shodan API Key">
          <el-input v-model="apiConfig.shodanKey" placeholder="Shodan API Key"></el-input>
        </el-form-item>
        <el-form-item label="Censys API ID">
          <el-input v-model="apiConfig.censysID" placeholder="Censys API ID"></el-input>
        </el-form-item>
        <el-form-item label="Censys API Secret">
          <el-input v-model="apiConfig.censysSecret" type="password" placeholder="Censys API Secret"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveApiConfig">保存</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 字典管理 -->
    <div class="content-section">
      <h2>字典管理</h2>
      <el-table :data="dictionaries" style="width: 100%">
        <el-table-column prop="name" label="字典名称"></el-table-column>
        <el-table-column prop="path" label="文件路径"></el-table-column>
        <el-table-column prop="size" label="大小"></el-table-column>
        <el-table-column prop="action" label="操作">
          <template #default="scope">
            <el-button size="small" @click="editDictionary(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteDictionary(scope.row.name)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="dictionaries.length === 0" style="text-align: center; padding: 40px; color: #909399;">
        暂无字典
      </div>
      <el-button type="primary" @click="addDictionary" style="margin-top: 10px;">添加字典</el-button>
    </div>

    <!-- 界面设置 -->
    <div class="content-section">
      <h2>界面设置</h2>
      <el-form :model="uiConfig" label-width="120px">
        <el-form-item label="主题">
          <el-select v-model="uiConfig.theme" placeholder="选择主题">
            <el-option label="默认" value="default"></el-option>
            <el-option label="暗色" value="dark"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="语言">
          <el-select v-model="uiConfig.language" placeholder="选择语言">
            <el-option label="中文" value="zh"></el-option>
            <el-option label="English" value="en"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="字体大小">
          <el-select v-model="uiConfig.fontSize" placeholder="选择字体大小">
            <el-option label="小" value="small"></el-option>
            <el-option label="中" value="medium"></el-option>
            <el-option label="大" value="large"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveUiConfig">保存</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 添加字典对话框 -->
    <el-dialog v-model="addDictionaryVisible" title="添加字典" width="40%">
      <el-form :model="dictionaryForm">
        <el-form-item label="字典名称">
          <el-input v-model="dictionaryForm.name" placeholder="请输入字典名称"></el-input>
        </el-form-item>
        <el-form-item label="文件路径">
          <el-input v-model="dictionaryForm.path" placeholder="请输入字典文件路径"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDictionaryVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddDictionary">添加</el-button>
      </template>
    </el-dialog>

    <!-- 编辑字典对话框 -->
    <el-dialog v-model="editDictionaryVisible" title="编辑字典" width="40%">
      <el-form :model="dictionaryForm">
        <el-form-item label="字典名称">
          <el-input v-model="dictionaryForm.name" placeholder="请输入字典名称"></el-input>
        </el-form-item>
        <el-form-item label="文件路径">
          <el-input v-model="dictionaryForm.path" placeholder="请输入字典文件路径"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDictionaryVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEditDictionary">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

// 全局配置
const globalConfig = ref({
  httpProxy: '',
  threads: 10,
  timeout: 5,
  retries: 2,
  subdomainDict: ''
})

// API配置
const apiConfig = ref({
  fofaKey: '',
  shodanKey: '',
  censysID: '',
  censysSecret: ''
})

// 界面设置
const uiConfig = ref({
  theme: 'default',
  language: 'zh',
  fontSize: 'medium'
})

// 字典管理
const dictionaries = ref([
  { name: '子域名字典', path: './dict/subdomain.txt', size: '100KB' },
  { name: '端口字典', path: './dict/port.txt', size: '10KB' }
])

// 对话框状态
const addDictionaryVisible = ref(false)
const editDictionaryVisible = ref(false)
const dictionaryForm = ref({ name: '', path: '' })
const originalDictName = ref('')

// 加载配置
const loadConfig = async () => {
  try {
    const config = await window.go.main.App.GetConfig()
    globalConfig.value = {
      httpProxy: config.httpProxy || '',
      threads: config.webThreads || 10,
      timeout: config.webTimeout || 5,
      retries: config.retries || 2,
      subdomainDict: config.subdomainDict || ''
    }
  } catch (e) {
    console.error('Failed to load config:', e)
  }
}

// 保存全局配置
const saveGlobalConfig = async () => {
  try {
    const config = {
      proxy: globalConfig.value.httpProxy,
      webThreads: globalConfig.value.threads,
      webTimeout: globalConfig.value.timeout,
      retries: globalConfig.value.retries,
      subdomainDict: globalConfig.value.subdomainDict
    }
    const result = await window.go.main.App.SaveConfig(config)
    ElMessage.success(result)
  } catch (error) {
    ElMessage.error('保存配置失败: ' + error.message)
  }
}

// 保存API配置
const saveApiConfig = async () => {
  try {
    const result = await window.go.main.App.SaveApiConfig(apiConfig.value)
    ElMessage.success(result)
  } catch (error) {
    ElMessage.error('保存API配置失败: ' + error.message)
  }
}

// 保存界面配置
const saveUiConfig = async () => {
  try {
    const result = await window.go.main.App.SaveUiConfig(uiConfig.value)
    ElMessage.success(result)
  } catch (error) {
    ElMessage.error('保存界面配置失败: ' + error.message)
  }
}

// 添加字典
const addDictionary = () => {
  addDictionaryVisible.value = true
  dictionaryForm.value = { name: '', path: '' }
}

// 确认添加字典
const confirmAddDictionary = async () => {
  if (!dictionaryForm.value.name || !dictionaryForm.value.path) {
    ElMessage.warning('请填写完整信息')
    return
  }

  try {
    const result = await window.go.main.App.AddDictionary(dictionaryForm.value.name, dictionaryForm.value.path)
    ElMessage.success(result)
    dictionaries.value.push({
      name: dictionaryForm.value.name,
      path: dictionaryForm.value.path,
      size: '未知'
    })
    addDictionaryVisible.value = false
  } catch (error) {
    ElMessage.error('添加字典失败: ' + error.message)
  }
}

// 编辑字典
const editDictionary = (dictionary) => {
  originalDictName.value = dictionary.name
  dictionaryForm.value = {
    name: dictionary.name,
    path: dictionary.path
  }
  editDictionaryVisible.value = true
}

// 确认编辑字典
const confirmEditDictionary = async () => {
  if (!dictionaryForm.value.name || !dictionaryForm.value.path) {
    ElMessage.warning('请填写完整信息')
    return
  }

  try {
    if (originalDictName.value !== dictionaryForm.value.name) {
      await window.go.main.App.DeleteDictionary(originalDictName.value)
    }
    const result = await window.go.main.App.AddDictionary(dictionaryForm.value.name, dictionaryForm.value.path)
    ElMessage.success(result)
    
    const index = dictionaries.value.findIndex(d => d.name === originalDictName.value)
    if (index !== -1) {
      dictionaries.value[index] = {
        name: dictionaryForm.value.name,
        path: dictionaryForm.value.path,
        size: dictionaries.value[index].size
      }
    }
    editDictionaryVisible.value = false
  } catch (error) {
    ElMessage.error('编辑字典失败: ' + error.message)
  }
}

// 删除字典
const deleteDictionary = async (name) => {
  try {
    const result = await window.go.main.App.DeleteDictionary(name)
    ElMessage.success(result)
    dictionaries.value = dictionaries.value.filter(d => d.name !== name)
  } catch (error) {
    ElMessage.error('删除字典失败: ' + error.message)
  }
}

// 页面加载时获取配置
onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.system-settings {
  padding: 0;
}
</style>