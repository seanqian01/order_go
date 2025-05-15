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
          <el-descriptions-item label="编号">{{ signal.id }}</el-descriptions-item>
          <el-descriptions-item label="合约代码">{{ signal.symbol }}</el-descriptions-item>
          <el-descriptions-item label="代码简称">{{ signal.scode }}</el-descriptions-item>
          <el-descriptions-item label="合约类型">
            <el-tag type="info">
              {{ getContractTypeName(signal.contractType) || '未知类型' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="方向">
            <el-tag :type="signal.action === 'buy' ? 'success' : 'danger'">
              {{ signal.action === 'buy' ? '买入' : '卖出' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="价格">{{ signal.price }}</el-descriptions-item>
          <el-descriptions-item label="提醒标题">{{ signal.alert_title }}</el-descriptions-item>
          <el-descriptions-item label="策略ID">{{ signal.strategy_id }}</el-descriptions-item>
          <el-descriptions-item label="时间周期">{{ signal.time_circle }}</el-descriptions-item>
          <el-descriptions-item label="处理状态">
            <el-tag :type="getProcessStatusType(signal.process_status)">
              {{ getProcessStatusName(signal.process_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="处理原因" :span="2">{{ signal.process_reason || '无' }}</el-descriptions-item>
          <el-descriptions-item label="信号接收时间">{{ formatTime(signal.created_at) }}</el-descriptions-item>
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

  // 获取合约类型名称
  const getContractTypeName = (contractType) => {
    if (contractType === undefined || contractType === null) {
      return null
    }
    const typeMap = {
      1: '大A股票',
      2: '商品期货',
      3: 'ETF金融指数',
      4: '虚拟货币'
    }
    return typeMap[contractType] || '未知合约类型'
  }
  
  // 获取处理状态名称
  const getProcessStatusName = (status) => {
    switch (status) {
      case 'pending':
        return '未处理'
      case 'invalid':
        return '信号无效'
      case 'valid_no_order':
        return '有效未下单'
      case 'processed':
        return '已下单'
      default:
        return '未知状态'
    }
  }

  // 获取处理状态标签类型
  const getProcessStatusType = (status) => {
    switch (status) {
      case 'pending':
        return 'info'
      case 'invalid':
        return 'danger'
      case 'valid_no_order':
        return 'warning'
      case 'processed':
        return 'success'
      default:
        return 'info'
    }
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