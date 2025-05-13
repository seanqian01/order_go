<template>
    <div class="dashboard-container">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>信号总数</span>
              </div>
            </template>
            <div class="card-body">
              <h2>{{ stats.signalCount || 0 }}</h2>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>订单总数</span>
              </div>
            </template>
            <div class="card-body">
              <h2>{{ stats.orderCount || 0 }}</h2>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>账户总值</span>
              </div>
            </template>
            <div class="card-body">
              <h2>{{ stats.accountValue || '0.00' }} USDT</h2>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import request from '@/api/request'
  
  const stats = ref({
    signalCount: 0,
    orderCount: 0,
    accountValue: '0.00'
  })
  
  const fetchStats = async () => {
    try {
      const res = await request({
        url: '/api/stats',
        method: 'get'
      })
      stats.value = res
    } catch (error) {
      console.error('获取统计数据失败:', error)
    }
  }
  
  onMounted(() => {
    fetchStats()
  })
  </script>
  
  <style scoped>
  .dashboard-container {
    padding: 20px;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .card-body {
    text-align: center;
    padding: 20px 0;
  }
  </style>