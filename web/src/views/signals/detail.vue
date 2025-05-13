<template>
    <div class="signal-detail-container">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>信号详情</span>
            <el-button @click="goBack">返回</el-button>
          </div>
        </template>
        
        <el-descriptions :column="2" border v-loading="loading">
          <el-descriptions-item label="信号ID">{{ signal.id }}</el-descriptions-item>
          <el-descriptions-item label="交易对">{{ signal.symbol }}</el-descriptions-item>
          <el-descriptions-item label="方向">
            <el-tag :type="signal.action === 'buy' ? 'success' : 'danger'">
              {{ signal.action === 'buy' ? '买入' : '卖出' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="价格">{{ signal.price }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatTime(signal.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatTime(signal.updated_at) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { getSignalDetail } from '@/api/signal'
  import { formatTime } from '@/utils/format'
  
  const route = useRoute()
  const router = useRouter()
  const loading = ref(false)
  const signal = ref({})
  
  const fetchData = async () => {
    const id = route.params.id
    if (!id) return
    
    loading.value = true
    try {
      const res = await getSignalDetail(id)
      signal.value = res || {}
    } catch (error) {
      console.error('获取信号详情失败:', error)
    } finally {
      loading.value = false
    }
  }
  
  const goBack = () => {
    router.back()
  }
  
  onMounted(() => {
    fetchData()
  })
  </script>
  
  <style scoped>
  .signal-detail-container {
    padding: 20px;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  </style>