import { ElMessage } from 'element-plus'

export async function safeApiCall(apiFunc, ...args) {
  try {
    const result = await apiFunc(...args)
    return { success: true, data: result }
  } catch (error) {
    const message = error?.message || String(error)
    ElMessage.error(`操作失败: ${message}`)
    return { success: false, error }
  }
}

export async function withLoading(loadingRef, apiFunc, ...args) {
  if (loadingRef && typeof loadingRef === 'object') {
    loadingRef.value = true
  }
  try {
    const result = await safeApiCall(apiFunc, ...args)
    return result
  } finally {
    if (loadingRef && typeof loadingRef === 'object') {
      loadingRef.value = false
    }
  }
}

export function handleApiError(error, defaultMessage = '操作失败') {
  const message = error?.message || defaultMessage
  ElMessage.error(message)
  console.error('API Error:', error)
}

export function showSuccess(message) {
  ElMessage.success(message)
}

export function showError(message) {
  ElMessage.error(message)
}

export function showWarning(message) {
  ElMessage.warning(message)
}

export function showInfo(message) {
  ElMessage.info(message)
}
