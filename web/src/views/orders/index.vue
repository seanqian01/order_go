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
          <el-table-column prop="symbol" label="交易对" width="120" />
          <el-table-column label="方向" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.action === 'buy' ? 'success' : 'danger'">
                {{ scope.row.action === 'buy' ? '买入' : '卖出' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="价格" width="120" />
          <el-table-column prop="amount" label="数量" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="180">
            <template #default="scope">
              {{ formatTime(scope.row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="scope">
              <el-button size="small" @click="viewDetail(scope.row)">
                详情
              </el-button>
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
  .pagination-container {
    margin-top: 20px;
    text-align: right;
  }
  </style>