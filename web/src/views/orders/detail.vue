<template>
    <div class="order-detail-container">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>订单详情</span>
            <el-button @click="goBack">返回</el-button>
          </div>
        </template>
        
        <el-descriptions :column="isMobile ? 1 : 2" border v-loading="loading">
          <el-descriptions-item label="系统订单号">{{ order.id }}</el-descriptions-item>
          <el-descriptions-item label="交易所订单号">{{ order.order_id || '-' }}</el-descriptions-item>

          <el-descriptions-item label="策略ID">{{ order.strategy_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="交易所ID">{{ order.exchange_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="合约代码">{{ order.symbol || '-' }}</el-descriptions-item>
          <el-descriptions-item label="合约代码类型">
            <el-tag type="info">{{ getContractTypeName(order.contract_code) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="合约类型">
            <el-tag type="info">{{ order.contract_type === 'spot' ? '现货' : '合约' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="订单类型">
            <el-tag type="info">{{ order.order_type === 'limit' ? '限价单' : '市价单' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="方向">
            <el-tag :type="order.action === 'buy' ? 'success' : 'danger'">
              {{ order.action === 'buy' ? '买入' : '卖出' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="持仓方向">
            <el-tag type="warning" v-if="order.position_side">
              {{ order.position_side === 'open' ? '开仓' : '平仓' }}
            </el-tag>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="价格">{{ order.price || '-' }}</el-descriptions-item>
          <el-descriptions-item label="数量">{{ order.amount || '-' }}</el-descriptions-item>
          <el-descriptions-item label="成交价格">{{ order.filled_price || '-' }}</el-descriptions-item>
          <el-descriptions-item label="成交数量">{{ order.filled_amount || '-' }}</el-descriptions-item>
          <el-descriptions-item label="手续费">{{ order.fee ? `${order.fee} ${order.fee_currency}` : '-' }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(order.status)">
              {{ getStatusText(order.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatTime(order.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatTime(order.updated_at) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted, onUnmounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { getOrderDetail } from '@/api/order'
  import { formatTime } from '@/utils/format'
  
  const route = useRoute()
  const router = useRouter()
  const loading = ref(false)
  const order = ref({})
  const isMobile = ref(false)

  // 检测当前设备是否为移动端
  const checkIsMobile = () => {
    isMobile.value = window.innerWidth < 768
  }
  
  const fetchData = async () => {
    const id = route.params.id
    if (!id) return
    
    loading.value = true
    try {
      const res = await getOrderDetail(id)
      order.value = res || {}
    } catch (error) {
      console.error('获取订单详情失败:', error)
    } finally {
      loading.value = false
    }
  }
  
  const getStatusType = (status) => {
    const map = {
      'pending': 'warning',
      'completed': 'success',
      'canceled': 'info',
      'failed': 'danger'
    }
    return map[status] || 'info'
  }
  
  const getStatusText = (status) => {
    const map = {
      'pending': '处理中',
      'completed': '已完成',
      'canceled': '已取消',
      'failed': '失败'
    }
    return map[status] || '未知'
  }

  // 获取合约类型名称
  const getContractTypeName = (contractCode) => {
    if (contractCode === undefined || contractCode === null) {
      return '-'
    }
    const typeMap = {
      1: '大A股票',
      2: '商品期货',
      3: 'ETF金融指数',
      4: '虚拟货币'
    }
    return typeMap[contractCode] || contractCode
  }
  
  const goBack = () => {
    router.back()
  }


  
  onMounted(() => {
    checkIsMobile()
    window.addEventListener('resize', checkIsMobile)
    fetchData()
  })

  onUnmounted(() => {
    window.removeEventListener('resize', checkIsMobile)
  })
  </script>
  
  <style scoped>
  .order-detail-container {
    padding: 20px;
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  /* 移动端响应式样式 */
  @media screen and (max-width: 768px) {
    .order-detail-container {
      padding: 10px;
    }
    
    :deep(.el-descriptions__body) {
      background-color: transparent;
    }
    
    :deep(.el-descriptions__label) {
      width: 100px;
      padding: 8px;
    }
    
    :deep(.el-descriptions__content) {
      padding: 8px;
    }
  }
  </style>