<template>
    <div class="app-wrapper" :class="{ 'mobile': isMobile, 'sidebar-open': sidebarOpen }">
      <div class="mobile-mask" v-if="isMobile && sidebarOpen" @click="closeSidebar"></div>
      <sidebar class="sidebar-container" />
      <div class="main-container">
        <navbar @toggle-sidebar="toggleSidebar" />
        <app-main />
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted, onUnmounted } from 'vue'
  import Sidebar from './components/Sidebar.vue'
  import Navbar from './components/Navbar.vue'
  import AppMain from './components/AppMain.vue'
  
  // 响应式布局相关状态
  const isMobile = ref(false)
  const sidebarOpen = ref(true)
  
  // 检测当前设备是否为移动端
  const checkIsMobile = () => {
    isMobile.value = window.innerWidth < 768
    // 在移动端默认关闭侧边栏
    if (isMobile.value) {
      sidebarOpen.value = false
    } else {
      sidebarOpen.value = true
    }
  }
  
  // 切换侧边栏显示/隐藏
  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }
  
  // 关闭侧边栏
  const closeSidebar = () => {
    sidebarOpen.value = false
  }
  
  // 监听窗口大小变化
  onMounted(() => {
    checkIsMobile()
    window.addEventListener('resize', checkIsMobile)
  })
  
  onUnmounted(() => {
    window.removeEventListener('resize', checkIsMobile)
  })
  </script>
  
  <style scoped>
  .app-wrapper {
    position: relative;
    height: 100vh;
    width: 100%;
  }
  .sidebar-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 210px;
    height: 100%;
    background-color: #304156;
    overflow-y: auto;
    transition: all 0.3s;
    z-index: 1001;
  }
  .main-container {
    margin-left: 210px;
    position: relative;
    min-height: 100%;
    transition: all 0.3s;
  }
  
  /* 移动端样式 */
  .mobile .sidebar-container {
    transform: translateX(-100%);
    width: 210px;
  }
  .mobile.sidebar-open .sidebar-container {
    transform: translateX(0);
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
  }
  .mobile .main-container {
    margin-left: 0;
  }
  
  /* 遮罩层，点击关闭侧边栏 */
  .mobile-mask {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.3);
    z-index: 1000;
  }
  </style>