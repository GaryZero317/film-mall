<template>
  <div class="faq-list-container">
    <el-card class="table-card" shadow="never">
      <div class="header-actions">
        <div class="left">
          <el-select v-model="category" placeholder="问题分类" clearable @change="handleCategoryChange">
            <el-option label="全部" :value="-1" />
            <el-option label="商品相关" :value="1" />
            <el-option label="订单相关" :value="2" />
            <el-option label="账户相关" :value="3" />
            <el-option label="其他问题" :value="4" />
          </el-select>
        </div>
        <div class="right">
          <el-button type="primary" @click="handleAdd">添加FAQ</el-button>
        </div>
      </div>

      <!-- FAQ列表 -->
      <el-table 
        :data="faqList" 
        style="width: 100%" 
        v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="question" label="问题" min-width="200" align="left" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getCategoryTag(scope.row.category)" size="small">
              {{ getCategoryText(scope.row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="100" align="center" />
        <el-table-column label="操作" width="180" align="center">
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
          background
        />
      </div>
    </el-card>

    <!-- FAQ表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑FAQ' : '添加FAQ'"
      width="700px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form :model="faqForm" :rules="faqRules" ref="faqFormRef" label-width="80px">
        <el-form-item label="问题" prop="question">
          <el-input v-model="faqForm.question" placeholder="请输入问题"></el-input>
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="faqForm.category" placeholder="请选择分类">
            <el-option label="商品相关" :value="1" />
            <el-option label="订单相关" :value="2" />
            <el-option label="账户相关" :value="3" />
            <el-option label="其他问题" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="faqForm.type" placeholder="请选择类型">
            <el-option label="普通" :value="1" />
            <el-option label="重要" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="faqForm.priority" :min="0" :max="100" />
        </el-form-item>
        <el-form-item label="答案" prop="answer">
          <el-input
            v-model="faqForm.answer"
            type="textarea"
            :rows="8"
            placeholder="请输入答案"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确 定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getFaqList, addFaq, updateFaq, deleteFaq } from '../../api/customer-service'
import { adminService } from '../../api/request'  // 导入adminService
import { callDirectApi } from '../../api/index'  // 导入直接调用API的方法

// FAQ列表数据
const faqList = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const category = ref(-1)

// 表单数据
const dialogVisible = ref(false)
const isEdit = ref(false)
const faqForm = reactive({
  id: null,
  question: '',
  answer: '',
  category: 1,
  type: 1,     // 添加type字段，默认为普通类型
  priority: 0
})
const faqRules = {
  question: [{ required: true, message: '请输入问题', trigger: 'blur' }],
  answer: [{ required: true, message: '请输入答案', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}
const faqFormRef = ref(null)
const submitting = ref(false)

// 获取FAQ列表
const fetchFaqList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (category.value !== -1) {
      params.category = category.value
    }
    
    console.log('发送FAQ请求参数:', params)
    
    // 使用直接调用API的方法
    const response = await callDirectApi('/api/user/service/faq/list', 'post', params)
    console.log('完整API响应:', response)
    
    // 检查响应格式，适当处理
    if (response && response.list) {
      faqList.value = response.list
      total.value = response.total || 0
    } else {
      console.error('响应格式不符合预期:', response)
      ElMessage.warning('获取FAQ列表成功，但数据格式有问题')
      faqList.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取FAQ列表失败', error)
    // 增加更详细的错误信息
    if (error.response) {
      console.error('错误状态码:', error.response.status)
      console.error('错误详情:', error.response.data)
    }
    ElMessage.error('获取FAQ列表失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 处理添加
const handleAdd = () => {
  isEdit.value = false
  faqForm.id = null
  faqForm.question = ''
  faqForm.answer = ''
  faqForm.category = 1
  faqForm.type = 1
  faqForm.priority = 0
  dialogVisible.value = true
}

// 处理编辑
const handleEdit = (row) => {
  isEdit.value = true
  faqForm.id = row.id
  faqForm.question = row.question
  faqForm.answer = row.answer
  faqForm.category = row.category
  faqForm.type = row.type
  faqForm.priority = row.priority
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!faqFormRef.value) return
  
  await faqFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      // 确保字段是正确的数据类型并映射字段名称
      const formData = {
        ...faqForm,
        category: Number(faqForm.category),
        type: Number(faqForm.type),
        priority: Number(faqForm.priority),
        sort: Number(faqForm.priority),  // 添加sort字段，值同priority
        order: Number(faqForm.priority)  // 添加order字段，以防API也检查这个
      }
      
      if (isEdit.value) {
        console.log('更新FAQ数据:', formData)
        await updateFaq(faqForm.id, formData)
        ElMessage.success('更新成功')
      } else {
        console.log('添加FAQ数据:', formData)
        await addFaq(formData)
        ElMessage.success('添加成功')
      }
      dialogVisible.value = false
      fetchFaqList()
    } catch (error) {
      console.error(isEdit.value ? '更新失败' : '添加失败', error)
      if (error.response) {
        console.error('错误状态码:', error.response.status)
        console.error('错误详情:', error.response.data)
      }
      ElMessage.error(isEdit.value ? '更新失败' : '添加失败: ' + (error.message || '未知错误'))
    } finally {
      submitting.value = false
    }
  })
}

// 处理删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除此FAQ吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteFaq(row.id)
      ElMessage.success('删除成功')
      fetchFaqList()
    } catch (error) {
      console.error('删除失败', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

// 页面大小变化
const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  fetchFaqList()
}

// 页码变化
const handleCurrentChange = (val) => {
  page.value = val
  fetchFaqList()
}

// 分类变化
const handleCategoryChange = () => {
  page.value = 1
  fetchFaqList()
}

// 获取分类文本
const getCategoryText = (category) => {
  const categoryMap = {
    1: '商品相关',
    2: '订单相关',
    3: '账户相关',
    4: '其他问题'
  }
  return categoryMap[category] || '未知分类'
}

// 获取分类标签类型
const getCategoryTag = (category) => {
  const categoryMap = {
    1: 'success',
    2: 'warning',
    3: 'info',
    4: 'danger'
  }
  return categoryMap[category] || ''
}

onMounted(() => {
  fetchFaqList()
})
</script>

<style scoped>
.faq-list-container {
  padding: 16px;
}

.table-card {
  margin-bottom: 16px;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style> 