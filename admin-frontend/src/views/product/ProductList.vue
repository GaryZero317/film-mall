<template>
  <div class="product-list-container">
    <div class="action-bar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>添加商品
      </el-button>
    </div>

    <!-- 商品列表 -->
    <el-table :data="productList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="商品名称" width="200" />
      <el-table-column prop="desc" label="商品描述" show-overflow-tooltip />
      <el-table-column prop="stock" label="库存" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.stock > 10 ? 'success' : 'danger'">
            {{ scope.row.stock }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="价格" width="120">
        <template #default="scope">
          ¥{{ (scope.row.amount / 100).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
            {{ scope.row.status === 1 ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
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

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'add' ? '添加商品' : '编辑商品'"
      v-model="dialogVisible"
      width="600px">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="商品描述" prop="desc">
          <el-input
            v-model="form.desc"
            type="textarea"
            rows="3"
            placeholder="请输入商品描述" />
        </el-form-item>
        <el-form-item label="库存数量" prop="stock">
          <el-input-number
            v-model="form.stock"
            :min="0"
            :max="9999"
            controls-position="right" />
        </el-form-item>
        <el-form-item label="商品价格" prop="amount">
          <el-input-number
            v-model="form.amount"
            :min="0"
            :precision="2"
            :step="0.01"
            controls-position="right">
            <template #prefix>¥</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="商品状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">上架</el-radio>
            <el-radio :label="0">下架</el-radio>
          </el-radio-group>
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
import { createProduct, updateProduct, removeProduct, getProductDetail } from '../../api/product'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add')
const productList = ref([])

// 表单相关
const formRef = ref(null)
const form = ref({
  name: '',
  desc: '',
  stock: 0,
  amount: 0,
  status: 1
})

const rules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  desc: [
    { required: true, message: '请输入商品描述', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入库存数量', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入商品价格', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择商品状态', trigger: 'change' }
  ]
}

// 获取商品列表
const fetchProductList = async () => {
  loading.value = true
  try {
    // 这里需要后端提供获取商品列表的接口
    // 临时使用商品详情接口模拟
    const res = await getProductDetail({ id: 1 })
    productList.value = [res]
  } catch (error) {
    console.error('获取商品列表失败:', error)
    ElMessage.error('获取商品���表失败')
  } finally {
    loading.value = false
  }
}

// 添加商品
const handleAdd = () => {
  dialogType.value = 'add'
  form.value = {
    name: '',
    desc: '',
    stock: 0,
    amount: 0,
    status: 1
  }
  dialogVisible.value = true
}

// 编辑商品
const handleEdit = (row) => {
  dialogType.value = 'edit'
  form.value = {
    id: row.id,
    name: row.name,
    desc: row.desc,
    stock: row.stock,
    amount: row.amount / 100, // 转换为元
    status: row.status
  }
  dialogVisible.value = true
}

// 删除商品
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该商品吗？', '提示', {
      type: 'warning'
    })
    
    await removeProduct({ id: row.id })
    ElMessage.success('删除成功')
    fetchProductList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除商品失败:', error)
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
    
    const submitData = {
      ...form.value,
      amount: Math.round(form.value.amount * 100) // 转换为分
    }
    
    if (dialogType.value === 'add') {
      await createProduct(submitData)
      ElMessage.success('添加成功')
    } else {
      await updateProduct(submitData)
      ElMessage.success('更新成功')
    }
    
    dialogVisible.value = false
    fetchProductList()
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchProductList()
})
</script>

<style scoped>
.product-list-container {
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

:deep(.el-input-number .el-input__wrapper) {
  width: 200px;
}
</style> 