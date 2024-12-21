<template>
  <div class="payment-list-container">
    <el-card class="table-card" shadow="never">
      <!-- 支付列表 -->
      <el-table 
        :data="paymentList" 
        style="width: 100%" 
        v-loading="loading">
        <el-table-column prop="id" label="支付ID" width="80" align="center" />
        <el-table-column prop="oid" label="订单ID" width="80" align="center" />
        <el-table-column prop="amount" label="金额" width="100" align="center">
          <template #default="scope">
            ¥{{ (scope.row.amount / 100).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getPaymentStatusType(scope.row.status)" size="small">
              {{ getPaymentStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" min-width="160" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="updateTime" label="更新时间" min-width="160" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.updateTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleDetail(scope.row)">
              详情
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
          @update:current-page="page = $event"
          @update:page-size="pageSize = $event"
          background
        />
      </div>
    </el-card>

    <!-- 支付详情弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="支付详情"
      width="500px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      destroy-on-close
    >
      <div v-loading="detailLoading">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="支付ID" label-class-name="label-cell" content-class-name="content-cell">
            {{ paymentDetail.id || '--' }}
          </el-descriptions-item>
          <el-descriptions-item label="订单ID" label-class-name="label-cell" content-class-name="content-cell">
            {{ paymentDetail.orderId || '--' }}
          </el-descriptions-item>
          <el-descriptions-item label="支付金额" label-class-name="label-cell" content-class-name="content-cell">
            {{ paymentDetail.amount ? `¥${(paymentDetail.amount / 100).toFixed(2)}` : '--' }}
          </el-descriptions-item>
          <el-descriptions-item label="支付状态" label-class-name="label-cell" content-class-name="content-cell">
            <el-tag :type="getPaymentStatusType(paymentDetail.status)" size="small">
              {{ getPaymentStatusText(paymentDetail.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" label-class-name="label-cell" content-class-name="content-cell">
            {{ formatTime(paymentDetail.createTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间" label-class-name="label-cell" content-class-name="content-cell">
            {{ formatTime(paymentDetail.updateTime) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminPaymentList, getPaymentDetail } from '../../api/payment'

const loading = ref(false)
const paymentList = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 详情弹窗相关
const dialogVisible = ref(false)
const detailLoading = ref(false)
const paymentDetail = ref({})

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

// 获取支付列表
const fetchPaymentList = async () => {
  loading.value = true
  try {
    const res = await getAdminPaymentList({
      page: page.value,
      pageSize: pageSize.value
    })
    paymentList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取支付列表失败:', error)
    ElMessage.error('获取支付列表失败')
  } finally {
    loading.value = false
  }
}

// 支付状态
const getPaymentStatusText = (status) => {
  if (status === null || status === undefined) return '未知状态'
  const statusMap = {
    0: '待支付',
    1: '支付成功',
    2: '支付失败',
    3: '已退款'
  }
  return statusMap[status] || '未知状态'
}

const getPaymentStatusType = (status) => {
  if (status === null || status === undefined) return 'info'
  const typeMap = {
    0: 'warning',
    1: 'success',
    2: 'danger',
    3: 'info'
  }
  return typeMap[status] || 'info'
}

// 查看支付详情
const handleDetail = async (row) => {
  dialogVisible.value = true
  detailLoading.value = true
  try {
    const res = await getPaymentDetail({ id: row.id })
    console.log('API响应数据:', res) // 添加调试日志
    
    // 使用API返回的数据，将oid映射为orderId
    if (res && res.code === 0) {
      paymentDetail.value = {
        id: res.id || row.id || '--',
        orderId: res.oid || row.oid || '--',
        amount: res.amount || row.amount || 0,
        status: res.status ?? row.status,
        createTime: res.createTime || row.createTime || '',
        updateTime: res.updateTime || row.updateTime || ''
      }
    } else {
      // 如果API返回错误，使用表格行数据
      paymentDetail.value = {
        id: row.id || '--',
        orderId: row.oid || '--',
        amount: row.amount || 0,
        status: row.status,
        createTime: row.createTime || '',
        updateTime: row.updateTime || ''
      }
    }
    console.log('处理后的支付详情数据:', paymentDetail.value)
  } catch (error) {
    console.error('获取支付详情失败:', error)
    ElMessage.error('获取支付详情失败')
    // 发生错误时使用表格行数据
    paymentDetail.value = {
      id: row.id || '--',
      orderId: row.oid || '--',
      amount: row.amount || 0,
      status: row.status,
      createTime: row.createTime || '',
      updateTime: row.updateTime || ''
    }
  } finally {
    detailLoading.value = false
  }
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  fetchPaymentList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchPaymentList()
}

onMounted(() => {
  fetchPaymentList()
})
</script>

<style scoped>
.payment-list-container {
  padding: 20px;
  background-color: #f5f7fa;
  box-sizing: border-box;
  min-height: calc(100vh - 60px);
}

.table-card {
  margin: 0;
  width: calc(100vw - 280px);  /* 考虑左侧菜单宽度 */
}

.table-card :deep(.el-card__body) {
  padding: 0;
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

:deep(.el-tag--small) {
  height: 20px;
  padding: 0 6px;
  font-size: 12px;
}

/* 详情弹窗样式 */
:deep(.el-descriptions__label) {
  width: 120px;
  justify-content: flex-end;
  padding: 12px 16px;
}

:deep(.el-descriptions__content) {
  padding: 12px 16px;
}

:deep(.el-dialog__body) {
  padding: 20px;
}

:deep(.el-dialog__footer) {
  padding: 0 20px 20px;
  border-top: none;
}

.label-cell {
  background-color: #f5f7fa;
  color: #606266;
  font-weight: 600;
}

.content-cell {
  color: #303133;
}
</style> 