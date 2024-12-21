<template>
  <div class="admin-list-container">
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="action-bar">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>添加管理员
          </el-button>
        </div>
      </template>

      <!-- 管理员列表 -->
      <el-table 
        :data="adminList" 
        style="width: 100%" 
        v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="createTime" label="创建时间" min-width="180" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="updateTime" label="更新时间" min-width="180" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.updateTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">
              编辑
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
          @update:current-page="page = $event"
          @update:page-size="pageSize = $event"
          background
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'add' ? '添加管理员' : '编辑管理员'"
      v-model="dialogVisible"
      width="500px"
      destroy-on-close>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        status-icon>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item 
          label="密码" 
          prop="password"
          :rules="dialogType === 'add' ? rules.password : []">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password />
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getAdminList, createAdmin, updateAdmin, removeAdmin } from '../../api/admin'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add')
const adminList = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 表单相关
const formRef = ref(null)
const form = ref({
  username: '',
  password: ''
})

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

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

// 获取管理员列表
const fetchAdminList = async () => {
  loading.value = true
  try {
    const res = await getAdminList({
      page: page.value,
      pageSize: pageSize.value
    })
    adminList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取管理员列表失败:', error)
    ElMessage.error('获取管理员列表失败')
  } finally {
    loading.value = false
  }
}

// 添加管理员
const handleAdd = () => {
  dialogType.value = 'add'
  form.value = {
    username: '',
    password: ''
  }
  dialogVisible.value = true
}

// 编辑管理员
const handleEdit = (row) => {
  dialogType.value = 'edit'
  form.value = {
    id: row.id,
    username: row.username,
    password: ''
  }
  dialogVisible.value = true
}

// 删除管理员
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该管理员吗？', '提示', {
      type: 'warning'
    })
    
    await removeAdmin({ id: row.id })
    ElMessage.success('删除成功')
    fetchAdminList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除管理员失败:', error)
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
    
    if (dialogType.value === 'add') {
      await createAdmin(form.value)
      ElMessage.success('添加成功')
    } else {
      await updateAdmin(form.value)
      ElMessage.success('更新成功')
    }
    
    dialogVisible.value = false
    fetchAdminList()
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
  }
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  fetchAdminList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchAdminList()
}

onMounted(() => {
  fetchAdminList()
})
</script>

<style scoped>
.admin-list-container {
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