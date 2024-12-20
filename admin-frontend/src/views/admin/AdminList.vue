<template>
  <div class="admin-list-container">
    <div class="action-bar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>添加管理员
      </el-button>
    </div>

    <!-- 管理员列表 -->
    <el-table :data="adminList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" width="180" />
      <el-table-column prop="level" label="级别">
        <template #default="scope">
          <el-tag :type="scope.row.level === 1 ? 'danger' : 'info'">
            {{ scope.row.level === 1 ? '超级管理员' : '普通管理员' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">
            编辑
          </el-button>
          <el-button 
            type="danger" 
            size="small" 
            @click="handleDelete(scope.row)"
            :disabled="scope.row.level === 1">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'add' ? '添加管理员' : '编辑管理员'"
      v-model="dialogVisible"
      width="500px">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="dialogType === 'edit'" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
        <el-form-item label="级别" prop="level">
          <el-select v-model="form.level" placeholder="请选择级别">
            <el-option label="普通管理员" :value="0" />
            <el-option label="超级管理员" :value="1" />
          </el-select>
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
import { createAdmin, updateAdmin, deleteAdmin, getAdminInfo } from '../../api/admin'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add') // 'add' 或 'edit'
const adminList = ref([])

// 表单相关
const formRef = ref(null)
const form = ref({
  username: '',
  password: '',
  level: 0
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  level: [
    { required: true, message: '请选择级别', trigger: 'change' }
  ]
}

// 获取管理员列表
const fetchAdminList = async () => {
  loading.value = true
  try {
    // 这里需要后端提供获取管理员列表的接口
    const res = await getAdminInfo()
    adminList.value = [res] // 临时方案：只显示当前管理员
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
    password: '',
    level: 0
  }
  dialogVisible.value = true
}

// 编辑管理员
const handleEdit = (row) => {
  dialogType.value = 'edit'
  form.value = {
    id: row.id,
    username: row.username,
    password: '',
    level: row.level
  }
  dialogVisible.value = true
}

// 删除管理员
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该管理员吗？', '提示', {
      type: 'warning'
    })
    
    await deleteAdmin({ id: row.id })
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

onMounted(() => {
  fetchAdminList()
})
</script>

<style scoped>
.admin-list-container {
  padding: 20px;
}

.action-bar {
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 