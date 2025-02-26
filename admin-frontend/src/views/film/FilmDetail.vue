<template>
  <div class="container">
    <div class="page-header">
      <h2>胶片冲洗订单详情</h2>
      <el-button @click="$router.push('/film/list')">返回列表</el-button>
    </div>

    <el-card v-loading="loading" class="info-card">
      <template #header>
        <div class="card-header">
          <span>订单基本信息</span>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单ID">{{ orderInfo.id }}</el-descriptions-item>
        <el-descriptions-item label="订单号">{{ orderInfo.foid }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ orderInfo.uid }}</el-descriptions-item>
        <el-descriptions-item label="地址ID">{{ orderInfo.address_id }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusTagType(orderInfo.status)">
            {{ orderInfo.status_desc }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="回寄底片">
          <el-tag :type="orderInfo.return_film ? 'success' : 'info'">
            {{ orderInfo.return_film ? '是' : '否' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="总价">
          {{ orderInfo.total_price ? (orderInfo.total_price / 100).toFixed(2) + ' 元' : '0.00 元' }}
        </el-descriptions-item>
        <el-descriptions-item label="运费">
          {{ orderInfo.shipping_fee ? (orderInfo.shipping_fee / 100).toFixed(2) + ' 元' : '0.00 元' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ orderInfo.create_time }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ orderInfo.update_time }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ orderInfo.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>订单项信息</span>
        </div>
      </template>

      <el-table :data="orderInfo.items || []" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="film_type" label="胶片类型" />
        <el-table-column prop="film_brand" label="胶片品牌" />
        <el-table-column prop="size" label="尺寸规格" />
        <el-table-column prop="quantity" label="数量" width="80" />
        <el-table-column label="单价" width="120">
          <template #default="scope">
            {{ (scope.row.price / 100).toFixed(2) }} 元
          </template>
        </el-table-column>
        <el-table-column label="总价" width="120">
          <template #default="scope">
            {{ (scope.row.amount / 100).toFixed(2) }} 元
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
      </el-table>
    </el-card>

    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>订单操作</span>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-form-item label="订单状态">
          <el-select v-model="form.status">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="回寄底片">
          <el-switch v-model="form.return_film" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleUpdate">更新订单</el-button>
          <el-button v-if="orderInfo.status === 3" type="danger" @click="handleDelete">删除订单</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getFilmOrderDetail, updateFilmOrder, deleteFilmOrder } from '@/api/film'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const orderInfo = ref({})

// 表单数据
const form = reactive({
  status: 0,
  return_film: false,
  remark: ''
})

// 订单状态选项
const statusOptions = [
  { value: 0, label: '待付款' },
  { value: 1, label: '冲洗处理中' },
  { value: 2, label: '待收货' },
  { value: 3, label: '已完成' }
]

// 获取订单详情
const getOrderDetail = async (id) => {
  loading.value = true
  try {
    const res = await getFilmOrderDetail(id)
    if (res.code === 0) {
      orderInfo.value = res.data
      form.status = res.data.status
      form.return_film = res.data.return_film
      form.remark = res.data.remark || ''
    } else {
      ElMessage.error(res.msg || '获取订单详情失败')
    }
  } catch (error) {
    console.error('获取订单详情出错:', error)
    ElMessage.error('获取订单详情出错')
  } finally {
    loading.value = false
  }
}

// 更新订单
const handleUpdate = async () => {
  try {
    const res = await updateFilmOrder(orderInfo.value.id, {
      status: form.status,
      return_film: form.return_film,
      remark: form.remark
    })
    
    if (res.code === 0) {
      ElMessage.success('更新订单成功')
      getOrderDetail(orderInfo.value.id)
    } else {
      ElMessage.error(res.msg || '更新订单失败')
    }
  } catch (error) {
    console.error('更新订单出错:', error)
    ElMessage.error('更新订单出错')
  }
}

// 删除订单
const handleDelete = () => {
  ElMessageBox.confirm('确定要删除该订单吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteFilmOrder(orderInfo.value.id)
      if (res.code === 0) {
        ElMessage.success('删除成功')
        router.push('/film/list')
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error) {
      console.error('删除订单出错:', error)
      ElMessage.error('删除订单出错')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const statusMap = {
    0: 'info',
    1: 'warning',
    2: 'primary',
    3: 'success'
  }
  return statusMap[status] || 'info'
}

onMounted(() => {
  const id = route.params.id
  if (id) {
    getOrderDetail(id)
  }
})
</script>

<style scoped>
.container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.info-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 