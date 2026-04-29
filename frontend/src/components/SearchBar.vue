<template>
  <el-input
    v-model="searchText"
    placeholder="搜索..."
    prefix-icon="Search"
    @input="debouncedSearch"
    @clear="handleClear"
    clearable
    :style="{ width: '100%', maxWidth: '400px' }"
  >
    <template #append>
      <el-select v-model="searchField" placeholder="字段" style="width: 120px;" @change="handleSearch">
        <el-option label="全部" value="all"></el-option>
        <el-option label="名称" value="name"></el-option>
        <el-option label="描述" value="desc"></el-option>
      </el-select>
    </template>
  </el-input>
</template>

<script setup>
import { ref, watch, onUnmounted } from 'vue'

const searchText = ref('')
const searchField = ref('all')
let debounceTimer = null

const emit = defineEmits(['search'])

const handleSearch = () => {
  emit('search', {
    text: searchText.value,
    field: searchField.value
  })
}

const debouncedSearch = () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  debounceTimer = setTimeout(() => {
    handleSearch()
  }, 300)
}

const handleClear = () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  emit('search', {
    text: '',
    field: searchField.value
  })
}

watch(searchField, () => {
  handleSearch()
})

onUnmounted(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})
</script>

<style scoped>
/* 苹果风格样式 */
.el-input__wrapper {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.el-input__wrapper:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.el-select__wrapper {
  border-radius: 0 12px 12px 0;
}
</style>