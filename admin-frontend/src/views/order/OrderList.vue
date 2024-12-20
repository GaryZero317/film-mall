<template>
  <div class="order-list-container">
    <el-table :data="orderList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="订单ID" width="120" />
      <el-table-column prop="uid" label="用户ID" width="120" />
      <el-table-column prop="pid" label="商品ID" width="120" />
      <el-table-column prop="amount" label="金额">
        <template #default="scope">
          ¥{{ (scope.row.amount / 100).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="getOrderStatusType(scope.row.status)">
            {{ getOrderStatusText(scope.row.status) }}
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
import { getOrderList } from '../../api/order'

const loading = ref(false)
const orderList = ref([])

// 获取订单列表
const fetchOrderList = async () => {
  loading.value = true
  try {
    const res = await getOrderList({ uid: 0 })
    orderList.value = res.data || []
  } catch (error) {
    console.error('获取订单列表失败:', error)
    ElMessage.error('获取订单列表失败')
  } finally {
    loading.value = false
  }
}

// 订单状态
const getOrderStatusText = (status) => {
  const statusMap = {
    0: '待支付',
    1: '已支付',
    2: '已发货',
    3: '已完成',
    4: '已取消'
  }
  return statusMap[status] || '未知状态'
}

const getOrderStatusType = (status) => {
  const typeMap = {
    0: 'warning',
    1: 'success',
    2: 'primary',
    3: 'success',
    4: 'info'
  }
  return typeMap[status] || 'info'
}

// 查看订单详情
const handleDetail = (row) => {
  ElMessage.info('订单详情功能开发中')
}

onMounted(() => {
  fetchOrderList()
})
</script>

<style scoped>
.order-list-container {
  padding: 20px;
}
</style> 