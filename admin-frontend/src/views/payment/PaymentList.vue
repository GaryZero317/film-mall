<template>
  <div class="payment-list-container">
    <el-table :data="paymentList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="支付ID" width="120" />
      <el-table-column prop="uid" label="用户ID" width="120" />
      <el-table-column prop="oid" label="订单ID" width="120" />
      <el-table-column prop="amount" label="金额">
        <template #default="scope">
          ¥{{ (scope.row.amount / 100).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="source" label="支付方式" width="120">
        <template #default="scope">
          <el-tag>{{ getPaymentSource(scope.row.source) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="getPaymentStatusType(scope.row.status)">
            {{ getPaymentStatusText(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleDetail(scope.row)">
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPaymentDetail } from '../../api/payment'

const loading = ref(false)
const paymentList = ref([])

// 获取支付列表
const fetchPaymentList = async () => {
  loading.value = true
  try {
    // 这里需要后端提供获取支付列表的接口
    // 临时使用支付详情接口模拟
    const res = await getPaymentDetail({ id: 1 })
    paymentList.value = [res]
  } catch (error) {
    console.error('获取支付列表失败:', error)
    ElMessage.error('获取支付列表失败')
  } finally {
    loading.value = false
  }
}

// 支付方式
const getPaymentSource = (source) => {
  const sourceMap = {
    1: '支付宝',
    2: '微信支付',
    3: '银行卡'
  }
  return sourceMap[source] || '其他'
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

// 查看支付详情
const handleDetail = (row) => {
  ElMessage.info('支付详情功能开发中')
}

onMounted(() => {
  fetchPaymentList()
})
</script>

<style scoped>
.payment-list-container {
  padding: 20px;
}
</style> 