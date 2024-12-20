<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <!-- 数据概览卡片 -->
      <el-col :span="6">
        <el-card shadow="hover" class="data-card">
          <template #header>
            <div class="card-header">
              <span>商品总数</span>
              <el-icon><Goods /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ statistics.productCount }}</div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="data-card">
          <template #header>
            <div class="card-header">
              <span>订单总数</span>
              <el-icon><List /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ statistics.orderCount }}</div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="data-card">
          <template #header>
            <div class="card-header">
              <span>支付总额</span>
              <el-icon><Money /></el-icon>
            </div>
          </template>
          <div class="card-value">¥{{ (statistics.totalAmount / 100).toFixed(2) }}</div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="data-card">
          <template #header>
            <div class="card-header">
              <span>用户数量</span>
              <el-icon><User /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ statistics.userCount }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近订单列表 -->
    <el-card class="recent-orders" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>最近订单</span>
        </div>
      </template>
      
      <el-table :data="recentOrders" style="width: 100%">
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
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Goods, List, Money, User } from '@element-plus/icons-vue'
import { getOrderList } from '../api/order'

const statistics = ref({
  productCount: 0,
  orderCount: 0,
  totalAmount: 0,
  userCount: 0
})

const recentOrders = ref([])
const loading = ref(false)

// 获取最近订单
const fetchRecentOrders = async () => {
  loading.value = true
  try {
    // 初始化默认数据
    recentOrders.value = []
    statistics.value = {
      productCount: 0,
      orderCount: 0,
      totalAmount: 0,
      userCount: 0
    }

    const res = await getOrderList({ uid: 0 })
    if (res && res.data) {
      recentOrders.value = Array.isArray(res.data) ? res.data.slice(0, 5) : []
      
      // 统计数据
      statistics.value = {
        productCount: 0, // 这里需要对接实际的统计接口
        orderCount: Array.isArray(res.data) ? res.data.length : 0,
        totalAmount: Array.isArray(res.data) ? res.data.reduce((sum, order) => sum + (order.amount || 0), 0) : 0,
        userCount: 0 // 这里需要对接实际的统计接口
      }
    }
  } catch (error) {
    console.error('获取订单列表失败:', error)
    ElMessage({
      message: '获取订单列表失败，将显示默认数据',
      type: 'warning'
    })
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

onMounted(() => {
  fetchRecentOrders()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.data-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
}

.data-card .card-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  text-align: center;
  margin-top: 10px;
}

.recent-orders .card-header {
  font-size: 16px;
  font-weight: bold;
}
</style> 