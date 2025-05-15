<template>
  <div class="strategy-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>策略管理</span>
          <div>
            <el-button type="primary" @click="openDialog('create')">新增策略</el-button>
            <el-button @click="fetchData">刷新</el-button>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="strategyList"
        border
        style="width: 100%"
        class="responsive-table"
      >
        <el-table-column label="序号" width="80">
          <template #default="scope">
            {{ (pagination.currentPage - 1) * pagination.pageSize + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="策略名称" width="150" />
        <el-table-column prop="code" label="策略代码" width="150" />
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status ? 'success' : 'danger'">
              {{ scope.row.status ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <div class="operation-buttons">
              <el-button size="small" @click="openDialog('edit', scope.row)">
                编辑
              </el-button>
              <el-button 
                size="small" 
                type="danger" 
                @click="handleDelete(scope.row)"
              >
                删除
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
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 新增/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'create' ? '新增策略' : '编辑策略'"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="策略名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入策略名称" />
        </el-form-item>
        <el-form-item label="策略代码" prop="code">
          <el-input v-model="form.code" placeholder="请输入策略代码" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            确认
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatTime } from '@/utils/format'
import axios from 'axios'

const loading = ref(false)
const submitting = ref(false)
const strategyList = ref([])
const total = ref(0)
const isMobile = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('create') // 'create' 或 'edit'
const currentId = ref(null)

const pagination = ref({
  currentPage: 1,
  pageSize: 10
})

const form = ref({
  name: '',
  code: '',
  status: true
})

const formRef = ref(null)
const rules = {
  name: [
    { required: true, message: '请输入策略名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入策略代码', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ]
}

// 检测当前设备是否为移动端
const checkIsMobile = () => {
  isMobile.value = window.innerWidth < 768
}

const fetchData = async () => {
  loading.value = true
  try {
    console.log('开始获取策略列表...')
    const res = await axios.get('/api/strategies', {
      params: {
        page: pagination.value.currentPage,
        limit: pagination.value.pageSize
      }
    })
    console.log('获取策略列表响应:', res)
    if (res.data && res.data.items) {
      strategyList.value = res.data.items
      total.value = res.data.total || 0
      console.log('策略列表数据:', strategyList.value)
    } else {
      console.error('响应数据格式不正确:', res.data)
      strategyList.value = []
      total.value = 0
      ElMessage.warning('获取策略列表数据格式不正确')
    }
  } catch (error) {
    console.error('获取策略列表失败:', error)
    console.error('错误详情:', error.response ? error.response.data : '无响应数据')
    ElMessage.error(`获取策略列表失败: ${error.message || '未知错误'}`)
    strategyList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const openDialog = (type, row) => {
  dialogType.value = type
  dialogVisible.value = true
  
  if (type === 'create') {
    form.value = {
      name: '',
      code: '',
      status: true
    }
    currentId.value = null
  } else {
    form.value = {
      name: row.name,
      code: row.code,
      status: row.status
    }
    currentId.value = row.id
  }
  
  // 在下一个事件循环中重置表单校验结果
  setTimeout(() => {
    formRef.value?.clearValidate()
  }, 0)
}

const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      if (dialogType.value === 'create') {
        await axios.post('/api/strategies', form.value)
        ElMessage.success('策略创建成功')
      } else {
        await axios.put(`/api/strategies/${currentId.value}`, form.value)
        ElMessage.success('策略更新成功')
      }
      
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      console.error('提交策略数据失败:', error)
      ElMessage.error(error.response?.data?.error || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除策略 "${row.name}" 吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await axios.delete(`/api/strategies/${row.id}`)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error('删除策略失败:', error)
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }).catch(() => {
    // 用户取消删除，不做任何操作
  })
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
.strategy-container {
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

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}

/* 操作按钮样式 */
.operation-buttons {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

/* 移动端响应式样式 */
@media screen and (max-width: 768px) {
  .strategy-container {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .responsive-table {
    width: 100%;
    overflow-x: auto;
  }
}
</style>