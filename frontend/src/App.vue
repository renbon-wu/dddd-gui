<template>
  <LoadingProvider>
    <div class="app-container">
      <el-container class="container">
        <!-- 左侧导航菜单 -->
        <el-aside width="200px" class="sidebar">
          <div class="logo">
            <h2>VulScanX</h2>
          </div>
          <el-menu
            :default-active="activeMenu"
            class="menu"
            router
            background-color="rgba(255, 255, 255, 0.05)"
            text-color="#e8e8ed"
            active-text-color="#007aff"
            @select="handleMenuSelect"
          >
            <el-menu-item index="scan">
              <el-icon><Search /></el-icon>
              <span>扫描任务</span>
            </el-menu-item>
            <el-menu-item index="target">
              <el-icon><Location /></el-icon>
              <span>目标管理</span>
            </el-menu-item>
            <el-menu-item index="fingerprint">
              <el-icon><Lock /></el-icon>
              <span>指纹管理</span>
            </el-menu-item>
            <el-menu-item index="poc">
              <el-icon><Warning /></el-icon>
              <span>POC管理</span>
            </el-menu-item>
            <el-menu-item index="result">
              <el-icon><DataAnalysis /></el-icon>
              <span>结果分析</span>
            </el-menu-item>
            <el-menu-item index="setting">
              <el-icon><Setting /></el-icon>
              <span>系统设置</span>
            </el-menu-item>
          </el-menu>
        </el-aside>

        <!-- 主内容区 -->
        <el-container>
          <!-- 顶部操作栏 -->
          <el-header class="header">
            <div class="header-left">
              <h1>dddd 漏洞扫描工具</h1>
            </div>
            <div class="header-right">
              <el-button type="primary" @click="runScan">开始扫描</el-button>
            </div>
          </el-header>

          <!-- 内容区域 -->
          <el-main class="main-content">
            <ScanTask v-if="activeMenu === 'scan'" />
            <TargetManagement v-if="activeMenu === 'target'" />
            <FingerprintManagement v-if="activeMenu === 'fingerprint'" />
            <PocManagement v-if="activeMenu === 'poc'" />
            <ResultAnalysis v-if="activeMenu === 'result'" />
            <SystemSettings v-if="activeMenu === 'setting'" />
          </el-main>
        </el-container>
      </el-container>
    </div>
  </LoadingProvider>
</template>

<script setup>
import { ref, provide, inject } from 'vue'
import { Search, Location, Lock, Warning, DataAnalysis, Setting } from '@element-plus/icons-vue'
import LoadingProvider from './components/LoadingProvider.vue'
import ScanTask from './views/ScanTask.vue'
import TargetManagement from './views/TargetManagement.vue'
import FingerprintManagement from './views/FingerprintManagement.vue'
import PocManagement from './views/PocManagement.vue'
import ResultAnalysis from './views/ResultAnalysis.vue'
import SystemSettings from './views/SystemSettings.vue'

const activeMenu = ref('scan')

const handleMenuSelect = (key) => {
  activeMenu.value = key
}

const runScan = () => {
  activeMenu.value = 'scan'
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  background-color: #f5f5f7;
  color: #1d1d1f;
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app-container {
  height: 100vh;
  overflow: hidden;
  background: linear-gradient(180deg, #f5f5f7 0%, #e8e8ed 100%);
}

.container {
  height: 100%;
}

.sidebar {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-right: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.05);
}

.logo {
  text-align: center;
  padding: 24px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.logo h2 {
  color: #ffffff;
  font-size: 17px;
  margin: 0;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.menu {
  flex: 1;
  border-right: none;
  overflow-y: auto;
  background: transparent;
}

.el-menu-item {
  transition: all 0.2s ease;
  margin: 4px 8px;
  border-radius: 10px;
}

.el-menu-item:hover {
  background: rgba(0, 122, 255, 0.08) !important;
}

.el-menu-item.is-active {
  background: rgba(0, 122, 255, 0.15) !important;
}

.el-menu-item.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: #007aff;
  border-radius: 0 3px 3px 0;
}

.header {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 56px;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.05);
}

.header-left h1 {
  font-size: 17px;
  color: #1d1d1f;
  margin: 0;
  font-weight: 600;
  letter-spacing: -0.5px;
}

.header-right {
  display: flex;
  gap: 12px;
}

.el-button--primary {
  background: linear-gradient(135deg, #007aff 0%, #5856d6 100%) !important;
  border-color: transparent !important;
  border-radius: 10px !important;
  padding: 8px 20px !important;
  font-weight: 500;
  box-shadow: 0 4px 14px rgba(0, 122, 255, 0.3);
  transition: all 0.2s ease;
}

.el-button--primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(0, 122, 255, 0.4);
}

.el-button--primary:active {
  transform: translateY(0);
}

.main-content {
  padding: 20px;
  overflow-y: auto;
  background: transparent;
  height: calc(100vh - 56px);
}

.content-section {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  margin-bottom: 20px;
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.content-section:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
}

.content-section h2 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #1d1d1f;
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.3px;
}

.content-section h3 {
  margin-top: 0;
  margin-bottom: 16px;
  color: #6e6e73;
  font-size: 15px;
  font-weight: 500;
}

@media screen and (max-width: 1024px) {
  .sidebar {
    width: 180px !important;
  }
  
  .header-left h1 {
    font-size: 16px;
  }
  
  .main-content {
    padding: 16px;
  }
  
  .content-section {
    padding: 20px;
  }
}

@media screen and (max-width: 768px) {
  .sidebar {
    width: 160px !important;
  }
  
  .logo h2 {
    font-size: 15px;
  }
  
  .header-left h1 {
    font-size: 15px;
  }
  
  .main-content {
    padding: 12px;
  }
  
  .content-section {
    padding: 16px;
  }
  
  .content-section h2 {
    font-size: 16px;
  }
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.25);
}

.el-table {
  border-radius: 12px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.9);
}

.el-table th {
  background: rgba(0, 0, 0, 0.02) !important;
  font-weight: 500;
  color: #6e6e73;
  font-size: 13px;
}

.el-table td {
  font-size: 14px;
  color: #1d1d1f;
}

.el-card {
  border-radius: 12px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.el-card__header {
  background: rgba(0, 0, 0, 0.02);
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
}

.el-form-item {
  margin-bottom: 16px;
}

.el-input, .el-select, .el-textarea {
  border-radius: 10px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.8);
}

.el-input:focus, .el-select:focus, .el-textarea:focus {
  border-color: #007aff;
  box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
}

.el-tag {
  border-radius: 8px;
  padding: 4px 10px;
  font-size: 12px;
}

.el-tag--danger {
  background: rgba(255, 59, 48, 0.1);
  color: #ff3b30;
  border: none;
}

.el-tag--warning {
  background: rgba(255, 149, 0, 0.1);
  color: #ff9500;
  border: none;
}

.el-tag--success {
  background: rgba(52, 199, 89, 0.1);
  color: #34c759;
  border: none;
}

.el-tag--info {
  background: rgba(0, 122, 255, 0.1);
  color: #007aff;
  border: none;
}

.empty-state-wrapper {
  width: 100%;
  display: flex;
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 60px 40px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-text {
  font-size: 16px;
  color: #8a8a8d;
  margin: 0;
}

.el-dialog {
  border-radius: 20px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
}

.el-dialog__header {
  background: rgba(0, 0, 0, 0.02);
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
  padding: 20px 24px;
}

.el-dialog__title {
  font-size: 16px;
  font-weight: 600;
  color: #1d1d1f;
}

.el-dialog__body {
  padding: 24px;
}

.el-dialog__footer {
  padding: 16px 24px;
  border-top: 1px solid rgba(0, 0, 0, 0.04);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.el-button {
  border-radius: 8px;
  padding: 8px 16px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.el-button--default {
  color: #6e6e73;
  background: rgba(0, 0, 0, 0.04);
  border: none;
}

.el-button--default:hover {
  background: rgba(0, 0, 0, 0.08);
}

.el-button--danger {
  background: rgba(255, 59, 48, 0.1);
  color: #ff3b30;
  border: none;
}

.el-button--danger:hover {
  background: rgba(255, 59, 48, 0.2);
}

.el-descriptions {
  background: rgba(0, 0, 0, 0.02);
  border-radius: 12px;
  padding: 16px;
}

.el-descriptions__label {
  color: #6e6e73;
  font-weight: 500;
}

.el-descriptions__content {
  color: #1d1d1f;
}
</style>