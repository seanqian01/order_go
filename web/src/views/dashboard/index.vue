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
              <h2 class="clickable-count big-number" @click="goToSignals">{{ stats.signalCount || 0 }}</h2>
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
              <h2 class="clickable-count big-number" @click="goToOrders">{{ stats.orderCount || 0 }}</h2>
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
              <div class="account-value-container">
                <h2 class="big-number">{{ stats.accountValue || '0.00' }} USDT</h2>
                <el-tooltip :content="refreshTooltip" placement="top">
                  <el-button 
                    class="refresh-btn" 
                    :icon="RefreshRight" 
                    circle 
                    size="small" 
                    @click="refreshAccountValue"
                    :loading="refreshing"
                    :disabled="cooldownActive"
                    :type="cooldownActive ? 'info' : 'primary'"
                  />
                </el-tooltip>
              </div>
              <div v-if="lastRefreshTime" class="refresh-time">
                上次更新时间: {{ formatTime(lastRefreshTime) }}
                <span v-if="cooldownActive">(还需{{ cooldownTimeLeft }})</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted, computed, onUnmounted } from 'vue'
  import { useRouter } from 'vue-router'
  import request from '@/api/request'
  import { RefreshRight } from '@element-plus/icons-vue'
  import { formatTime } from '@/utils/format'
  import { ElMessage } from 'element-plus'
  
  const router = useRouter()
  
  const stats = ref({
    signalCount: 0,
    orderCount: 0,
    accountValue: '0.00'
  })
  const refreshing = ref(false)
  const lastRefreshTime = ref(null)
  const cooldownActive = ref(false)
  const cooldownSeconds = ref(0)
  let cooldownTimer = null
  
  const goToSignals = () => {
    router.push('/signals/list')
  }

  const goToOrders = () => {
    router.push('/orders/list')
  }

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
  
  const startCooldown = () => {
    cooldownActive.value = true
    cooldownSeconds.value = 5 * 60 // 5分钟冷却时间
    
    if (cooldownTimer) {
      clearInterval(cooldownTimer)
    }
    
    cooldownTimer = setInterval(() => {
      if (cooldownSeconds.value <= 0) {
        cooldownActive.value = false
        clearInterval(cooldownTimer)
        cooldownTimer = null
      } else {
        cooldownSeconds.value--
      }
    }, 1000)
  }
  
  const cooldownTimeLeft = computed(() => {
    const minutes = Math.floor(cooldownSeconds.value / 60)
    const seconds = cooldownSeconds.value % 60
    return `${minutes}:${seconds.toString().padStart(2, '0')}`
  })
  
  const refreshTooltip = computed(() => {
    return cooldownActive.value 
      ? `刷新冷却中 (${cooldownTimeLeft.value})` 
      : '刷新账户总值'
  })
  
  const refreshAccountValue = async () => {
    if (refreshing.value || cooldownActive.value) return
    
    refreshing.value = true
    try {
      const res = await request({
        url: '/api/refresh-account',
        method: 'post'
      })
      stats.value.accountValue = res.accountValue
      lastRefreshTime.value = new Date()
      startCooldown()
      ElMessage.success('账户总值已刷新')
    } catch (error) {
      console.error('刷新账户总值失败:', error)
      ElMessage.error('刷新账户总值失败')
    } finally {
      refreshing.value = false
    }
  }
  
  onMounted(() => {
    fetchStats()
  })
  
  onUnmounted(() => {
    if (cooldownTimer) {
      clearInterval(cooldownTimer)
    }
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
  .account-value-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}
.refresh-btn {
  font-size: 14px;
}
.refresh-time {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  text-align: center;
}
.clickable-count {
  cursor: pointer;
  color: #409EFF;
  transition: color 0.3s;
}
.clickable-count:hover {
  color: #66b1ff;
  text-decoration: underline;
}
.big-number {
  font-size: 36px;
  font-weight: bold;
}
</style>