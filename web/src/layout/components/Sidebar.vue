<template>
    <div class="sidebar-container">
      <div class="logo-container">
        <router-link to="/" class="logo-link">
          <h1 class="logo-title">IBS量化交易管理系统</h1>
        </router-link>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><odometer /></el-icon>
          <span>控制台</span>
        </el-menu-item>
        
        <el-sub-menu index="/signals">
          <template #title>
            <el-icon><data-line /></el-icon>
            <span>交易信号</span>
          </template>
          <el-menu-item index="/signals/list">信号列表</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/orders">
          <template #title>
            <el-icon><shopping-cart /></el-icon>
            <span>交易订单</span>
          </template>
          <el-menu-item index="/orders/list">订单列表</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </div>
  </template>
  
  <script setup>
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  
  const route = useRoute()
  const activeMenu = computed(() => {
    return route.path
  })
  </script>
  
  <style scoped>
  .logo-container {
    height: 50px;
    padding: 10px 0;
    text-align: center;
  }
  .logo-title {
    color: #fff;
    font-size: 18px;
    margin: 0;
  }
  .logo-link {
    text-decoration: none;
  }
  </style>