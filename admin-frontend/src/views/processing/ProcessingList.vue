<template>
  <div class="processing-list-container">
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="action-bar">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>添加冲洗订单
          </el-button>
        </div>
      </template>

      <!-- 冲洗订单列表 -->
      <el-table 
        :data="processingList" 
        style="width: 100%" 
        v-loading="loading">
        <el-table-column prop="id" label="订单ID" width="80" align="center" />
        <el-table-column prop="customerName" label="客户姓名" min-width="120" />
        <el-table-column prop="phone" label="联系电话" width="150" />
        <el-table-column prop="filmType" label="胶卷类型" width="150">
          <template #default="scope">
            <el-tag>{{ scope.row.filmType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="serviceType" label="冲洗类型" width="150">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.serviceType === '哈苏X5'">{{ scope.row.serviceType }}</el-tag>
            <el-tag type="warning" v-else-if="scope.row.serviceType === '富士SP3000'">{{ scope.row.serviceType }}</el-tag>
            <el-tag v-else>{{ scope.row.serviceType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="处理状态" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ formatStatus(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="200">
          <template #default="scope">
            <el-progress :percentage="scope.row.progress" :status="getProgressStatus(scope.row.progress)" />
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" min-width="180" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="estimatedFinishTime" label="预计完成时间" min-width="180" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.estimatedFinishTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" align="center" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button type="success" size="small" @click="handleUpdateProgress(scope.row)">
              更新进度
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'add' ? '添加冲洗订单' : '编辑冲洗订单'"
      v-model="dialogVisible"
      width="600px">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px">
        <el-form-item label="客户姓名" prop="customerName">
          <el-input v-model="form.customerName" placeholder="请输入客户姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="胶卷类型" prop="filmType">
          <el-select v-model="form.filmType" placeholder="请选择胶卷类型">
            <el-option label="135胶卷" value="135胶卷" />
            <el-option label="120胶卷" value="120胶卷" />
            <el-option label="黑白胶卷" value="黑白胶卷" />
            <el-option label="彩色胶卷" value="彩色胶卷" />
          </el-select>
        </el-form-item>
        <el-form-item label="冲洗类型" prop="serviceType">
          <el-select v-model="form.serviceType" placeholder="请选择冲洗类型">
            <el-option label="哈苏X5" value="哈苏X5" />
            <el-option label="富士SP3000" value="富士SP3000" />
          </el-select>
        </el-form-item>
        <el-form-item label="处理状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择处理状态">
            <el-option :value="1" label="待处理" />
            <el-option :value="2" label="处理中" />
            <el-option :value="3" label="已完成" />
            <el-option :value="4" label="已取消" />
          </el-select>
        </el-form-item>
        <el-form-item label="进度" prop="progress">
          <el-slider v-model="form.progress" :step="10" show-stops />
        </el-form-item>
        <el-form-item label="预计完成时间" prop="estimatedFinishTime">
          <el-date-picker
            v-model="form.estimatedFinishTime"
            type="datetime"
            placeholder="选择预计完成时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledDate"
          />
        </el-form-item>
        <el-form-item label="备注信息" prop="remark">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 更新进度对话框 -->
    <el-dialog
      title="更新冲洗进度"
      v-model="progressDialogVisible"
      width="500px">
      <el-form
        ref="progressFormRef"
        :model="progressForm"
        label-width="120px">
        <el-form-item label="当前进度" prop="progress">
          <el-slider v-model="progressForm.progress" :step="10" show-stops />
        </el-form-item>
        <el-form-item label="处理状态" prop="status">
          <el-select v-model="progressForm.status" placeholder="请选择处理状态">
            <el-option :value="1" label="待处理" />
            <el-option :value="2" label="处理中" />
            <el-option :value="3" label="已完成" />
            <el-option :value="4" label="已取消" />
          </el-select>
        </el-form-item>
        <el-form-item label="进度备注" prop="progressRemark">
          <el-input
            v-model="progressForm.progressRemark"
            type="textarea"
            :rows="3"
            placeholder="请输入进度相关备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="progressDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitProgressUpdate" :loading="submittingProgress">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

// 模拟数据 - 实际项目中应该从API获取
const processingList = ref([
  {
    id: 1,
    customerName: '张三',
    phone: '13800138000',
    filmType: '135胶卷',
    serviceType: '哈苏X5',
    status: 2,
    progress: 40,
    createTime: '2023-02-20 10:30:00',
    estimatedFinishTime: '2023-02-23 18:00:00',
    remark: '需要扫描高清图片'
  },
  {
    id: 2,
    customerName: '李四',
    phone: '13900139000',
    filmType: '120胶卷',
    serviceType: '富士SP3000',
    status: 1,
    progress: 0,
    createTime: '2023-02-21 14:20:00',
    estimatedFinishTime: '2023-02-24 18:00:00',
    remark: ''
  },
  {
    id: 3,
    customerName: '王五',
    phone: '13700137000',
    filmType: '黑白胶卷',
    serviceType: '哈苏X5',
    status: 3,
    progress: 100,
    createTime: '2023-02-19 09:15:00',
    estimatedFinishTime: '2023-02-22 12:00:00',
    remark: '需要冲洗多份'
  }
])

const loading = ref(false)
const submitting = ref(false)
const submittingProgress = ref(false)
const dialogVisible = ref(false)
const progressDialogVisible = ref(false)
const dialogType = ref('add')
const total = ref(processingList.value.length)
const page = ref(1)
const pageSize = ref(10)

// 表单相关
const formRef = ref(null)
const progressFormRef = ref(null)

const form = ref({
  customerName: '',
  phone: '',
  filmType: '',
  serviceType: '',
  status: 1,
  progress: 0,
  estimatedFinishTime: '',
  remark: ''
})

const progressForm = ref({
  id: null,
  progress: 0,
  status: 1,
  progressRemark: ''
})

const rules = {
  customerName: [
    { required: true, message: '请输入客户姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  filmType: [
    { required: true, message: '请选择胶卷类型', trigger: 'change' }
  ],
  serviceType: [
    { required: true, message: '请选择冲洗类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择处理状态', trigger: 'change' }
  ],
  estimatedFinishTime: [
    { required: true, message: '请选择预计完成时间', trigger: 'change' }
  ]
}

// 禁用今天之前的日期
const disabledDate = (time) => {
  return time.getTime() < Date.now() - 8.64e7 // 不能选择今天之前的日期
}

// 时间格式化函数
const formatTime = (time) => {
  if (!time) return '--'
  try {
    const date = new Date(time)
    if (isNaN(date.getTime())) {
      return '--'
    }
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch (error) {
    console.error('时间格式化错误:', error)
    return '--'
  }
}

// 处理状态格式化
const formatStatus = (status) => {
  const statusMap = {
    1: '待处理',
    2: '处理中',
    3: '已完成',
    4: '已取消'
  }
  return statusMap[status] || '未知状态'
}

// 获取状态对应的类型
const getStatusType = (status) => {
  const typeMap = {
    1: 'info', // 待处理显示为蓝色
    2: 'warning', // 处理中显示为橙色
    3: 'success', // 已完成显示为绿色
    4: 'danger' // 已取消显示为红色
  }
  return typeMap[status] || 'info'
}

// 获取进度状态
const getProgressStatus = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 60) return 'warning'
  if (progress > 0) return ''
  return 'exception'
}

// 添加冲洗订单
const handleAdd = () => {
  dialogType.value = 'add'
  form.value = {
    customerName: '',
    phone: '',
    filmType: '',
    serviceType: '',
    status: 1,
    progress: 0,
    estimatedFinishTime: '',
    remark: ''
  }
  dialogVisible.value = true
}

// 编辑冲洗订单
const handleEdit = (row) => {
  dialogType.value = 'edit'
  form.value = { ...row }
  dialogVisible.value = true
}

// 更新冲洗进度
const handleUpdateProgress = (row) => {
  progressForm.value = {
    id: row.id,
    progress: row.progress,
    status: row.status,
    progressRemark: ''
  }
  progressDialogVisible.value = true
}

// 删除冲洗订单
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该冲洗订单吗？', '提示', {
      type: 'warning'
    })
    
    // 这里应该调用API删除订单
    // 模拟删除
    const index = processingList.value.findIndex(item => item.id === row.id)
    if (index !== -1) {
      processingList.value.splice(index, 1)
      total.value -= 1
      ElMessage.success('删除成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除冲洗订单失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // 这里应该调用API保存订单
    // 模拟保存
    setTimeout(() => {
      if (dialogType.value === 'add') {
        const newId = processingList.value.length > 0 
          ? Math.max(...processingList.value.map(item => item.id)) + 1 
          : 1
        const newOrder = {
          ...form.value,
          id: newId,
          createTime: formatTime(new Date())
        }
        processingList.value.push(newOrder)
        total.value += 1
        ElMessage.success('添加成功')
      } else {
        const index = processingList.value.findIndex(item => item.id === form.value.id)
        if (index !== -1) {
          processingList.value[index] = { ...form.value }
          ElMessage.success('更新成功')
        }
      }
      
      dialogVisible.value = false
      submitting.value = false
    }, 1000)
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('表单验证失败，请检查表单')
    submitting.value = false
  }
}

// 提交进度更新
const submitProgressUpdate = async () => {
  submittingProgress.value = true
  
  // 这里应该调用API更新进度
  // 模拟更新
  setTimeout(() => {
    const index = processingList.value.findIndex(item => item.id === progressForm.value.id)
    if (index !== -1) {
      processingList.value[index].progress = progressForm.value.progress
      processingList.value[index].status = progressForm.value.status
      
      // 如果进度为100%，自动设置状态为已完成
      if (progressForm.value.progress === 100 && progressForm.value.status !== 3) {
        processingList.value[index].status = 3
      }
      
      ElMessage.success('进度更新成功')
      progressDialogVisible.value = false
    }
    submittingProgress.value = false
  }, 1000)
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  // 这里应该重新加载数据
}

const handleCurrentChange = (val) => {
  page.value = val
  // 这里应该重新加载数据
}

// 模拟获取冲洗订单列表
const fetchProcessingList = () => {
  loading.value = true
  // 这里应该调用API获取列表
  
  // 模拟API延迟
  setTimeout(() => {
    // 数据已经在上面初始化了
    loading.value = false
  }, 500)
}

onMounted(() => {
  fetchProcessingList()
})
</script>

<style scoped>
.processing-list-container {
  padding: 20px;
  background-color: #f5f7fa;
  box-sizing: border-box;
  min-height: calc(100vh - 60px);
}

.table-card {
  margin: 0;
}

.table-card :deep(.el-card__header) {
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
}

.action-bar {
  display: flex;
  justify-content: flex-start;
  margin: 0;
}

/* 表格样式优化 */
:deep(.el-table) {
  --el-table-border-color: transparent;
  --el-table-header-bg-color: #f5f7fa;
}

:deep(.el-table th.el-table__cell) {
  font-weight: 600;
  color: #606266;
  background-color: #f5f7fa;
  height: 45px;
  border-bottom: 1px solid #ebeef5;
}

:deep(.el-table td.el-table__cell) {
  height: 45px;
  padding: 6px 0;
  border-bottom: 1px solid #ebeef5;
}

:deep(.el-table::before) {
  display: none;
}

:deep(.el-table__inner-wrapper::before) {
  display: none;
}

.pagination-container {
  padding: 15px 20px;
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #ebeef5;
}

:deep(.el-button) {
  padding: 8px 16px;
}

:deep(.el-button--small) {
  padding: 5px 12px;
}

:deep(.el-dialog__body) {
  padding: 20px 40px;
}

/* 按钮间距 */
:deep(.el-button + .el-button) {
  margin-left: 8px;
}

/* 卡片内容区域padding */
:deep(.el-card__body) {
  padding: 0;
}

/* 设置表格容器宽度 */
.table-card {
  width: calc(100vw - 280px);  /* 考虑左侧菜单宽度 */
  margin: 0 auto;
}
</style> 