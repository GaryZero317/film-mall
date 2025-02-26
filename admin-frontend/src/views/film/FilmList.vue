<template>
  <div class="container">
    <div class="page-header">
      <h2>胶片冲洗订单管理</h2>
    </div>
    
    <div class="toolbar">
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="用户ID">
          <el-input v-model="queryParams.uid" placeholder="用户ID" clearable></el-input>
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">搜索</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table v-loading="loading" :data="orderList" border style="width: 100%">
      <el-table-column prop="id" label="订单ID" width="80" />
      <el-table-column prop="foid" label="订单号" width="180" />
      <el-table-column prop="uid" label="用户ID" width="80" />
      <el-table-column label="总价" width="120">
        <template #default="scope">
          {{ (scope.row.total_price / 100).toFixed(2) }} 元
        </template>
      </el-table-column>
      <el-table-column label="运费" width="120">
        <template #default="scope">
          {{ (scope.row.shipping_fee / 100).toFixed(2) }} 元
        </template>
      </el-table-column>
      <el-table-column label="回寄底片" width="80">
        <template #default="scope">
          <el-tag :type="scope.row.return_film ? 'success' : 'info'">
            {{ scope.row.return_film ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="订单状态" width="120">
        <template #default="scope">
          <el-tag :type="getStatusTagType(scope.row.status)">
            {{ scope.row.status_desc }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="create_time" label="创建时间" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button size="small" type="primary" @click="viewOrder(scope.row)">查看</el-button>
          <el-button size="small" type="warning" @click="editOrder(scope.row)">编辑</el-button>
          <el-button v-if="scope.row.status === 3" size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        :current-page="queryParams.page"
        :page-size="queryParams.page_size"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 编辑订单对话框 -->
    <el-dialog v-model="dialogVisible" title="编辑胶片冲洗订单" width="500px">
      <el-form ref="formRef" :model="form" label-width="120px">
        <el-form-item label="订单号">
          <span>{{ form.foid }}</span>
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="form.status">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="回寄底片">
          <el-switch v-model="form.return_film" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getFilmOrderList, updateFilmOrder, deleteFilmOrder } from '@/api/film'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const orderList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const formRef = ref(null)

// 查询参数
const queryParams = reactive({
  page: 1,
  page_size: 10,
  uid: '',
  status: ''
})

// 表单数据
const form = reactive({
  id: null,
  foid: '',
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

// 获取订单列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getFilmOrderList(queryParams)
    if (res.code === 0) {
      orderList.value = res.data.list
      total.value = res.data.total
    } else {
      ElMessage.error(res.msg || '获取订单列表失败')
    }
  } catch (error) {
    console.error('获取订单列表出错:', error)
    ElMessage.error('获取订单列表出错')
  } finally {
    loading.value = false
  }
}

// 查询按钮点击
const handleQuery = () => {
  queryParams.page = 1
  getList()
}

// 重置查询
const resetQuery = () => {
  queryParams.uid = ''
  queryParams.status = ''
  handleQuery()
}

// 分页大小变更
const handleSizeChange = (newSize) => {
  queryParams.page_size = newSize
  getList()
}

// 当前页变更
const handleCurrentChange = (newPage) => {
  queryParams.page = newPage
  getList()
}

// 查看订单详情
const viewOrder = (row) => {
  router.push(`/film/detail/${row.id}`)
}

// 编辑订单
const editOrder = (row) => {
  form.id = row.id
  form.foid = row.foid
  form.status = row.status
  form.return_film = row.return_film
  form.remark = row.remark
  dialogVisible.value = true
}

// 提交表单
const submitForm = async () => {
  try {
    const res = await updateFilmOrder(form.id, {
      status: form.status,
      return_film: form.return_film,
      remark: form.remark
    })
    
    if (res.code === 0) {
      ElMessage.success('更新订单成功')
      dialogVisible.value = false
      getList()
    } else {
      ElMessage.error(res.msg || '更新订单失败')
    }
  } catch (error) {
    console.error('更新订单出错:', error)
    ElMessage.error('更新订单出错')
  }
}

// 删除订单
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该订单吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteFilmOrder(row.id)
      if (res.code === 0) {
        ElMessage.success('删除成功')
        getList()
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
  getList()
})
</script>

<style scoped>
.container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
}

.toolbar {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style> 