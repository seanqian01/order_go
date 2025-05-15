<template>
    <div class="order-container">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>交易订单列表</span>
            <el-button type="primary" @click="fetchData">刷新</el-button>
          </div>
        </template>
        
        <el-table
          v-loading="loading"
          :data="orderList"
          border
          style="width: 100%"
        >
          <el-table-column label="序号" width="80">
            <template #default="scope">
              {{ (pagination.currentPage - 1) * pagination.pageSize + scope.$index + 1 }}
            </template>
          </el-table-column>
          <el-table-column prop="id" label="系统订单号" width="80" />
          <el-table-column label="交易所订单号" width="150">
            <template #default="scope">
              <el-text truncated>
                {{ scope.row.order_id || '-' }}
              </el-text>
            </template>
          </el-table-column>
          <el-table-column prop="symbol" label="合约代码" width="100" />
          <el-table-column label="交易方向" width="80">
            <template #default="scope">
              <el-tag :type="scope.row.action === 'buy' ? 'success' : 'danger'">
                {{ scope.row.action === 'buy' ? '买入' : '卖出' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="position_side" label="开仓平仓" width="80">
            <template #default="scope">
              <el-tag type="warning" v-if="scope.row.position_side">
                {{ scope.row.position_side === 'open' ? '开仓' : '平仓' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="order_type" label="订单类型" width="80">
            <template #default="scope">
              {{ scope.row.order_type === 'limit' ? '限价单' : '市价单' }}
            </template>
          </el-table-column>
          <el-table-column prop="price" label="价格" width="100">
            <template #default="scope">
              {{ scope.row.price ? formatNumber(scope.row.price) : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="数量" width="100">
            <template #default="scope">
              {{ scope.row.amount ? formatNumber(scope.row.amount) : '-' }}
            </template>
          </el-table-column>
          <el-table-column label="成交价" width="100">
            <template #default="scope">
              <span :class="{ 'highlight-text': scope.row.filled_price > 0 }">
                {{ scope.row.filled_price ? formatNumber(scope.row.filled_price) : '-' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="成交量" width="100">
            <template #default="scope">
              <span :class="{ 'highlight-text': scope.row.filled_amount > 0 }">
                {{ scope.row.filled_amount ? formatNumber(scope.row.filled_amount) : '-' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="手续费（USDT）" width="120">
            <template #default="scope">
              <span :class="{ 'highlight-text': scope.row.fee > 0 }">
                {{ scope.row.fee ? formatNumber(scope.row.fee) : '-' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="80">
            <template #default="scope">
              <el-tag 
                :type="getStatusType(scope.row.status)" 
                :effect="['filled', 'partially_filled'].includes(scope.row.status) ? 'dark' : 'light'"
                :class="{
                  'filled-status': scope.row.status === 'filled',
                  'partially-filled-status': scope.row.status === 'partially_filled'
                }"
              >
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="150">
            <template #default="scope">
              {{ formatTime(scope.row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="更新时间" width="150">
            <template #default="scope">
              {{ formatTime(scope.row.updated_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="scope">
              <div class="operation-buttons">
                <el-button size="small" @click="viewDetail(scope.row)">
                  详情
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <div class="pagination-container">
          <el-pagination
            :current-page="pagination.currentPage"
            :page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            :total="total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { getOrderList } from '@/api/order'
  import { useRouter } from 'vue-router'
  import { formatTime } from '@/utils/format'
  
  const router = useRouter()
  const loading = ref(false)
  const orderList = ref([])
  const total = ref(0)
  const pagination = ref({
    currentPage: 1,
    pageSize: 10
  })
  
  const fetchData = async () => {
    loading.value = true
    try {
      const res = await getOrderList({
        page: pagination.value.currentPage,
        limit: pagination.value.pageSize
      })
      orderList.value = res.items || []
      total.value = res.total || 0
    } catch (error) {
      console.error('获取订单列表失败:', error)
    } finally {
      loading.value = false
    }
  }
  
  const getStatusType = (status) => {
    const map = {
      'pending': 'warning',
      'completed': 'success',
      'filled': 'success',
      'partially_filled': 'warning',
      'canceled': 'info',
      'failed': 'danger'
    }
    return map[status] || 'info'
  }
  
  const getStatusText = (status) => {
    const map = {
      'created': '已创建',
      'pending': '处理中',
      'filled': '已成交',
      'partially_filled': '部分成交',
      'canceled': '已取消',
      'failed': '失败'
    }
    return map[status] || '未知'
  }

  const formatNumber = (num) => {
    if (num === undefined || num === null) return '-'
  
    // 如果是整数，不显示小数部分
    if (Number.isInteger(num)) return num.toString()
  
    // 直接返回原始价格值，不限制小数位数
    // 使用String()转换以保留原始精度
    return String(num)
  }

  
  const viewDetail = (row) => {
    router.push(`/orders/detail/${row.id}`)
  }
  
  const handleSizeChange = (val) => {
    pagination.value.pageSize = val
    fetchData()
  }

  const handleCurrentChange = (val) => {
    pagination.value.currentPage = val
    fetchData()
  }
  
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.order-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.operation-buttons {
  display: flex;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 高亮显示成交信息 */
.highlight-text {
  color: #409EFF;
  font-weight: bold;
}

/* 已成交状态的高亮样式 */
.filled-status {
  background-color: #67C23A !important;
  color: white !important;
  font-weight: bold;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  transform: scale(1.05);
}

/* 部分成交状态的高亮样式 */
.partially-filled-status {
  background-color: #E6A23C !important;
  color: white !important;
  font-weight: bold;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  transform: scale(1.05);
}
</style>