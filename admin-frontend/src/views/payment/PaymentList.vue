<template>
  <div class="payment-list-container">
    <el-table :data="paymentList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="支付ID" width="120" />
      <el-table-column prop="uid" label="用户ID" width="120" />
      <el-table-column prop="oid" label="订单ID" width="120" />
      <el-table-column prop="amount" label="金额" width="120">
        <template #default="scope">
          ¥{{ (scope.row.amount).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="source" label="支付方式" width="120">
        <template #default="scope">
          {{ getPaymentSourceText(scope.row.source) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="120">
        <template #default="scope">
          <el-tag :type="getPaymentStatusType(scope.row.status)">
            {{ getPaymentStatusText(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="180" />
      <el-table-column prop="updateTime" label="更新时间" width="180" />
    </el-table>

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
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminPaymentList } from '../../api/payment'

const loading = ref(false)
const paymentList = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

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

// 支付方式
const getPaymentSourceText = (source) => {
  const sourceMap = {
    0: '支付宝',
    1: '微信支付',
    2: '银行卡'
  }
  return sourceMap[source] || '未知方式'
}

// 支付状态
const getPaymentStatusText = (status) => {
  const statusMap = {
    0: '待支付',
    1: '支付成功',
    2: '支付失败',
    3: '已退款'
  }
  return statusMap[status] || '未知状态'
}

const getPaymentStatusType = (status) => {
  const typeMap = {
    0: 'warning',
    1: 'success',
    2: 'danger',
    3: 'info'
  }
  return typeMap[status] || 'info'
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
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 