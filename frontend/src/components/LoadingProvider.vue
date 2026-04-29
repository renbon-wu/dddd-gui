<template>
  <div class="loading-provider">
    <slot></slot>
    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-content">
        <div class="loading-spinner"></div>
        <p class="loading-text">{{ loadingText }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, provide, inject } from 'vue'

const isLoading = ref(false)
const loadingText = ref('加载中...')

const showLoading = (text = '加载中...') => {
  loadingText.value = text
  isLoading.value = true
}

const hideLoading = () => {
  isLoading.value = false
}

provide('loading', {
  showLoading,
  hideLoading
})
</script>

<style scoped>
.loading-provider {
  position: relative;
  width: 100%;
  height: 100%;
}

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.loading-content {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  padding: 32px 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(0, 122, 255, 0.1);
  border-top: 4px solid #007aff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.loading-text {
  margin: 0;
  color: #1d1d1f;
  font-size: 15px;
  font-weight: 500;
  letter-spacing: -0.3px;
}
</style>
