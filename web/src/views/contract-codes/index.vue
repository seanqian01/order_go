<template>
  <div class="contract-code-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>交易对列表</span>
          <div>
            <el-button type="primary" @click="addContractCode">新增交易对</el-button>
            <el-button type="default" @click="fetchData">刷新</el-button>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="contractCodeList"
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
        <el-table-column prop="code" label="交易对名称" width="120" />
        <el-table-column prop="min_amount" label="最小交易量" width="120" />
        <el-table-column prop="amount_precision" label="数量精度" width="100" />
        <el-table-column prop="price_precision" label="价格精度" width="100" />
        <el-table-column prop="max_position_ratio" label="最大仓位比例" width="120">
          <template #default="scope">
            {{ scope.row.max_position_ratio }}%
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status ? 'success' : 'danger'">
              {{ scope.row.status ? '启用' : '禁用' }}
            </el-tag>
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
              <el-button size="small" type="primary" @click="editContractCode(scope.row)">
                编辑
              </el-button>
              <el-button size="small" type="danger" @click="confirmDelete(scope.row)">
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
          layout="total, sizes, prev, pager, next"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 编辑/新增交易对对话框 -->
    <el-dialog 
      :title="dialogTitle" 
      v-model="dialogVisible" 
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form 
        :model="contractCodeForm" 
        :rules="rules" 
        ref="contractCodeFormRef" 
        label-width="120px"
        label-position="right"
      >
        <el-form-item label="合约代码" prop="symbol">
          <el-input v-model="contractCodeForm.symbol" placeholder="例如: BTCUSDT" />
        </el-form-item>
        <el-form-item label="交易对名称" prop="code">
          <el-input v-model="contractCodeForm.code" placeholder="例如: BTC/USDT" />
        </el-form-item>
        <el-form-item label="最小交易量" prop="min_amount">
          <el-input-number 
            v-model="contractCodeForm.min_amount" 
            :precision="5" 
            :step="0.001" 
            :min="0.00001"
          />
        </el-form-item>
        <el-form-item label="数量精度" prop="amount_precision">
          <el-input-number 
            v-model="contractCodeForm.amount_precision" 
            :precision="0" 
            :step="1" 
            :min="0" 
            :max="10"
          />
        </el-form-item>
        <el-form-item label="价格精度" prop="price_precision">
          <el-input-number 
            v-model="contractCodeForm.price_precision" 
            :precision="0" 
            :step="1" 
            :min="0" 
            :max="10"
          />
        </el-form-item>
        <el-form-item label="最大仓位比例" prop="max_position_ratio">
          <el-input-number 
            v-model="contractCodeForm.max_position_ratio" 
            :precision="2" 
            :step="1" 
            :min="0" 
            :max="100"
          />
          <span class="ratio-unit">%</span>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch 
            v-model="contractCodeForm.status" 
            active-text="启用" 
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
// 移除未使用的useRouter导入
// import { useRouter } from 'vue-router'
import request from '@/api/request'
import { formatTime } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'

// 移除未使用的router变量
// const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const contractCodeList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogType = ref('add') // 'add' 或 'edit'
const contractCodeFormRef = ref(null)

// 分页数据
const pagination = reactive({
  currentPage: 1,
  pageSize: 10
})

// 表单数据
const contractCodeForm = reactive({
  id: null,
  symbol: '',
  code: '',
  min_amount: 0.001,
  amount_precision: 3,
  price_precision: 5,
  max_position_ratio: 10,
  status: true
})

// 表单验证规则
const rules = {
  symbol: [
    { required: true, message: '请输入合约代码', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入交易对名称', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  min_amount: [
    { required: true, message: '请输入最小交易量', trigger: 'blur' }
  ],
  amount_precision: [
    { required: true, message: '请输入数量精度', trigger: 'blur' }
  ],
  price_precision: [
    { required: true, message: '请输入价格精度', trigger: 'blur' }
  ],
  max_position_ratio: [
    { required: true, message: '请输入最大仓位比例', trigger: 'blur' }
  ]
}

// 计算属性：对话框标题
const dialogTitle = computed(() => {
  return dialogType.value === 'add' ? '新增交易对' : '编辑交易对'
})

// 获取交易对列表数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/contract-codes', {
      params: {
        page: pagination.currentPage,
        limit: pagination.pageSize
      }
    })
    // 修复数据访问错误，因为request.js中的响应拦截器已经提取了response.data
    contractCodeList.value = res.items
    total.value = res.total
  } catch (error) {
    ElMessage.error('获取交易对列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 打开新增交易对对话框
const addContractCode = () => {
  dialogType.value = 'add'
  resetForm()
  dialogVisible.value = true
}

// 打开编辑交易对对话框
const editContractCode = (row) => {
  dialogType.value = 'edit'
  resetForm()
  Object.assign(contractCodeForm, row)
  dialogVisible.value = true
}

// 确认删除交易对
const confirmDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除交易对 ${row.symbol} 吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    deleteContractCode(row.id)
  }).catch(() => {
    // 用户取消删除
  })
}

// 删除交易对
const deleteContractCode = async (id) => {
  try {
    await request.delete(`/api/contract-codes/${id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

// 提交表单
const submitForm = async () => {
  if (!contractCodeFormRef.value) return
  
  await contractCodeFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (dialogType.value === 'add') {
          await request.post('/api/contract-codes', contractCodeForm)
          ElMessage.success('新增交易对成功')
        } else {
          await request.put(`/api/contract-codes/${contractCodeForm.id}`, contractCodeForm)
          ElMessage.success('更新交易对成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        ElMessage.error('操作失败: ' + error.message)
      } finally {
        submitting.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (contractCodeFormRef.value) {
    contractCodeFormRef.value.resetFields()
  }
  
  // 重置为默认值
  Object.assign(contractCodeForm, {
    id: null,
    symbol: '',
    code: '',
    min_amount: 0.001,
    amount_precision: 3,
    price_precision: 5,
    max_position_ratio: 10,
    status: true
  })
}

// 分页大小变化
const handleSizeChange = (val) => {
  pagination.pageSize = val
  fetchData()
}

// 当前页变化
const handleCurrentChange = (val) => {
  pagination.currentPage = val
  fetchData()
}

// 页面加载时获取数据
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.contract-code-container {
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

.ratio-unit {
  margin-left: 5px;
}

/* 操作按钮样式 */
.operation-buttons {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}

/* 移动端响应式样式 */
@media screen and (max-width: 768px) {
  .contract-code-container {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .card-header > div {
    margin-top: 10px;
    display: flex;
    flex-wrap: wrap;
  }
  
  .card-header > div > button {
    margin-bottom: 5px;
  }
  
  .pagination-container {
    justify-content: center;
  }
  
  /* 表格在移动端上的优化 */
  .responsive-table {
    font-size: 12px;
  }
  
  /* 对话框在移动端上的优化 */
  :deep(.el-dialog) {
    width: 95% !important;
  }
}
</style>
