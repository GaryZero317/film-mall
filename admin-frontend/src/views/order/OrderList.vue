<template>
  <div class="order-list-container">
    <el-card class="table-card" shadow="never">
      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-select v-model="status" placeholder="订单状态" clearable @change="handleStatusChange">
          <el-option label="全部" :value="-1" />
          <el-option label="待支付" :value="0" />
          <el-option label="已支付" :value="1" />
          <el-option label="已取消" :value="2" />
          <el-option label="已完成" :value="3" />
        </el-select>
      </div>

      <!-- 订单列表 -->
      <el-table 
        :data="orderList" 
        style="width: 100%" 
        v-loading="loading">
        <el-table-column prop="id" label="订单ID" width="80" align="center" />
        <el-table-column prop="oid" label="订单号" width="160" align="center" />
        <el-table-column prop="uid" label="用户ID" width="80" align="center" />
        <el-table-column label="商品信息" min-width="200" align="left">
          <template #default="scope">
            <div v-for="item in scope.row.items" :key="item.id" class="order-item">
              <el-image 
                :src="item.product_image" 
                :preview-src-list="[item.product_image]"
                style="width: 50px; height: 50px; object-fit: cover;"
              />
              <div class="item-info">
                <div class="product-name">{{ item.product_name }}</div>
                <div class="item-price">
                  <span>¥{{ (item.price / 100).toFixed(2) }}</span>
                  <span class="quantity">x{{ item.quantity }}</span>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="订单金额" width="120" align="center">
          <template #default="scope">
            <div>
              <div>商品：¥{{ (scope.row.total_price / 100).toFixed(2) }}</div>
              <div>运费：¥{{ (scope.row.shipping_fee / 100).toFixed(2) }}</div>
              <div class="total-amount">总计：¥{{ ((scope.row.total_price + scope.row.shipping_fee) / 100).toFixed(2) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getOrderStatusType(scope.row.status)" size="small">
              {{ getOrderStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" min-width="160" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.create_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" align="center">
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

    <!-- 订单详情弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="订单详情"
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      destroy-on-close
    >
      <div v-loading="detailLoading">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="订单号" label-class-name="label-cell" content-class-name="content-cell">
            {{ orderDetail.oid || '--' }}
          </el-descriptions-item>
          <el-descriptions-item label="用户ID" label-class-name="label-cell" content-class-name="content-cell">
            {{ orderDetail.uid || '--' }}
          </el-descriptions-item>
          <el-descriptions-item label="商品信息" label-class-name="label-cell" content-class-name="content-cell">
            <div v-for="item in orderDetail.items" :key="item.id" class="order-item">
              <el-image 
                :src="item.product_image" 
                :preview-src-list="[item.product_image]"
                style="width: 50px; height: 50px; object-fit: cover;"
              />
              <div class="item-info">
                <div class="product-name">{{ item.product_name }}</div>
                <div class="item-price">
                  <span>¥{{ (item.price / 100).toFixed(2) }}</span>
                  <span class="quantity">x{{ item.quantity }}</span>
                </div>
              </div>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="订单金额" label-class-name="label-cell" content-class-name="content-cell">
            <div>商品总额：¥{{ (orderDetail.total_price / 100).toFixed(2) }}</div>
            <div>运费：¥{{ (orderDetail.shipping_fee / 100).toFixed(2) }}</div>
            <div class="total-amount">订单总额：¥{{ ((orderDetail.total_price + orderDetail.shipping_fee) / 100).toFixed(2) }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="订单状态" label-class-name="label-cell" content-class-name="content-cell">
            <el-tag :type="getOrderStatusType(orderDetail.status)" size="small">
              {{ getOrderStatusText(orderDetail.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="备注" label-class-name="label-cell" content-class-name="content-cell">
            {{ orderDetail.remark || '--' }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" label-class-name="label-cell" content-class-name="content-cell">
            {{ formatTime(orderDetail.create_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间" label-class-name="label-cell" content-class-name="content-cell">
            {{ formatTime(orderDetail.update_time) }}
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
import { getAdminOrderList, getOrderDetail } from '../../api/order'

const loading = ref(false)
const orderList = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const status = ref(-1)  // 添加状态筛选

// 详情弹窗相关
const dialogVisible = ref(false)
const detailLoading = ref(false)
const orderDetail = ref({})

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

// 获取订单列表
const fetchOrderList = async () => {
  loading.value = true
  try {
    const res = await getAdminOrderList({
      page: page.value,
      pageSize: pageSize.value,
      status: status.value
    })
    if (res && res.code === 0) {
      orderList.value = res.data?.list || []
      total.value = res.data?.total || 0
    } else {
      throw new Error(res?.msg || '获取订单列表失败')
    }
  } catch (error) {
    console.error('获取订单列表失败:', error)
    ElMessage.error(error.message || '获取订单列表失败')
  } finally {
    loading.value = false
  }
}

// 查看订单详情
const handleDetail = async (row) => {
  dialogVisible.value = true
  detailLoading.value = true
  try {
    const res = await getOrderDetail(row.id)
    console.log('API响应数据:', res)
    
    if (res && res.code === 0) {
      orderDetail.value = res.data || {}
    } else {
      throw new Error(res?.msg || '获取订单详情失败')
    }
    console.log('处理后的订单详情数据:', orderDetail.value)
  } catch (error) {
    console.error('获取订单详情失败:', error)
    ElMessage.error(error.message || '获取订单详情失败')
    // 发生错误时使用表格行数据
    orderDetail.value = row
  } finally {
    detailLoading.value = false
  }
}

// 订单状态
const getOrderStatusText = (status) => {
  if (status === null || status === undefined) return '未知状态'
  const statusMap = {
    0: '待支付',
    1: '已支付',
    2: '已取消',
    3: '已完成'
  }
  return statusMap[status] || '未知状态'
}

const getOrderStatusType = (status) => {
  if (status === null || status === undefined) return 'info'
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
  fetchOrderList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchOrderList()
}

// 状态变化处理
const handleStatusChange = () => {
  page.value = 1
  fetchOrderList()
}

onMounted(() => {
  fetchOrderList()
})
</script>

<style scoped>
.order-list-container {
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

.order-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #ebeef5;
}

.order-item:last-child {
  border-bottom: none;
}

.item-info {
  margin-left: 12px;
  flex: 1;
}

.product-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}

.item-price {
  font-size: 13px;
  color: #606266;
}

.quantity {
  margin-left: 8px;
  color: #909399;
}

.total-amount {
  font-weight: bold;
  color: #f56c6c;
  margin-top: 4px;
}

.search-bar {
  padding: 16px 20px;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  gap: 16px;
}
</style> 