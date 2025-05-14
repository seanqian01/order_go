<template>
    <div class="signal-container">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>交易信号列表</span>
            <el-button type="primary" @click="fetchData">刷新</el-button>
          </div>
        </template>
        
        <el-table
          v-loading="loading"
          :data="signalList"
          border
          style="width: 100%"
          class="responsive-table"
        >
          <el-table-column label="序号" width="80">
            <template #default="scope">
              {{ (pagination.currentPage - 1) * pagination.pageSize + scope.$index + 1 }}
            </template>
          </el-table-column>
          <el-table-column prop="symbol" label="合约代码" width="120" />
          <el-table-column label="方向" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.action === 'buy' ? 'success' : 'danger'">
                {{ scope.row.action === 'buy' ? '买入' : '卖出' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="价格" width="120" />
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
            :layout="isMobile ? 'total, prev, next, jumper' : 'total, sizes, prev, pager, next'"
            :total="total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted, onUnmounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { getSignalList } from '@/api/signal'
  import { formatTime } from '@/utils/format'
  
  const router = useRouter()
  const loading = ref(false)
  const signalList = ref([])
  const total = ref(0)
  const isMobile = ref(false)
  const pagination = ref({
    currentPage: 1,
    pageSize: 10
  })
  
  // 检测当前设备是否为移动端
  const checkIsMobile = () => {
    isMobile.value = window.innerWidth < 768
  }
  
  const fetchData = async () => {
    loading.value = true
    try {
      const res = await getSignalList({
        page: pagination.value.currentPage,
        limit: pagination.value.pageSize
      })
      signalList.value = res.items || []
      total.value = res.total || 0
    } catch (error) {
      console.error('获取信号列表失败:', error)
    } finally {
      loading.value = false
    }
  }
  
  const viewDetail = (row) => {
    router.push(`/signals/detail/${row.id}`)
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
    checkIsMobile()
    window.addEventListener('resize', checkIsMobile)
    fetchData()
  })

  onUnmounted(() => {
    window.removeEventListener('resize', checkIsMobile)
  })
  </script>
  
  <style scoped>
  .signal-container {
    padding: 20px;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  /* 移动端响应式样式 */
  @media screen and (max-width: 768px) {
    .signal-container {
      padding: 10px;
    }
    
    .card-header {
      flex-direction: column;
      align-items: flex-start;
    }
    
    .card-header > span {
      margin-bottom: 10px;
    }
    
    .pagination-container {
      justify-content: center;
    }
    
    .responsive-table {
      font-size: 12px;
    }
  }
  </style>