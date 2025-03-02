<template>
  <div class="question-list-container">
    <el-card class="table-card" shadow="never">
      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-select v-model="status" placeholder="问题状态" clearable @change="handleStatusChange">
          <el-option label="全部" :value="-1" />
          <el-option label="待处理" :value="1" />
          <el-option label="已回复" :value="2" />
          <el-option label="已关闭" :value="3" />
        </el-select>
        <el-select v-model="type" placeholder="问题类型" clearable @change="handleTypeChange">
          <el-option label="全部" :value="-1" />
          <el-option label="商品相关" :value="1" />
          <el-option label="订单相关" :value="2" />
          <el-option label="账户相关" :value="3" />
          <el-option label="其他问题" :value="4" />
        </el-select>
      </div>

      <!-- 问题列表 -->
      <el-table 
        :data="questionList" 
        style="width: 100%" 
        v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="title" label="问题标题" min-width="200" align="left" show-overflow-tooltip />
        <el-table-column prop="user_id" label="用户ID" width="100" align="center" />
        <el-table-column prop="type" label="问题类型" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getQuestionTypeTag(scope.row.type)" size="small">
              {{ getQuestionTypeText(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="提交时间" width="170" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.create_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleDetail(scope.row)">
              {{ scope.row.status === 1 ? '回复' : '查看' }}
            </el-button>
            <el-button 
              v-if="scope.row.status !== 3" 
              type="danger" 
              size="small" 
              @click="handleClose(scope.row)"
            >
              关闭
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

    <!-- 问题详情弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="currentQuestion.status === 1 ? '回复问题' : '问题详情'"
      width="700px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <div v-loading="detailLoading" class="question-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="问题ID">
            {{ currentQuestion.id }}
          </el-descriptions-item>
          <el-descriptions-item label="问题标题">
            {{ currentQuestion.title }}
          </el-descriptions-item>
          <el-descriptions-item label="问题类型">
            <el-tag :type="getQuestionTypeTag(currentQuestion.type)" size="small">
              {{ getQuestionTypeText(currentQuestion.type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="问题状态">
            <el-tag :type="getStatusType(currentQuestion.status)" size="small">
              {{ getStatusText(currentQuestion.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">
            {{ formatTime(currentQuestion.create_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="用户ID">
            {{ currentQuestion.user_id }}
          </el-descriptions-item>
          <el-descriptions-item label="联系方式">
            {{ currentQuestion.contact_way || '--' }}
          </el-descriptions-item>
          <el-descriptions-item label="问题内容">
            <div class="question-content">{{ currentQuestion.content }}</div>
          </el-descriptions-item>
          <el-descriptions-item v-if="currentQuestion.status > 1" label="回复内容">
            <div class="reply-content">{{ currentQuestion.reply }}</div>
          </el-descriptions-item>
          <el-descriptions-item v-if="currentQuestion.status > 1" label="回复时间">
            {{ formatTime(currentQuestion.reply_time) }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 回复表单 -->
        <div v-if="currentQuestion.status === 1" class="reply-form">
          <h3>问题回复</h3>
          <el-form :model="replyForm" :rules="replyRules" ref="replyFormRef">
            <el-form-item prop="reply">
              <el-input
                v-model="replyForm.reply"
                type="textarea"
                :rows="5"
                placeholder="请输入回复内容"
              ></el-input>
            </el-form-item>
          </el-form>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button v-if="currentQuestion.status === 1" type="primary" @click="handleSubmitReply" :loading="submitting">
            提交回复
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 关闭问题确认 -->
    <el-dialog
      v-model="closeDialogVisible"
      title="确认关闭"
      width="400px"
    >
      <p>确定要关闭该问题吗？关闭后将不能再回复。</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeDialogVisible = false">取 消</el-button>
          <el-button type="danger" @click="confirmClose" :loading="closing">
            确认关闭
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  getQuestionList, 
  getQuestionListSafe,
  getQuestionDetail,
  replyQuestion 
} from '@/api/customer-service'
import { adminService } from '@/api/request'

// 问题列表数据
const questionList = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const status = ref(-1)
const type = ref(-1)

// 详情数据
const dialogVisible = ref(false)
const detailLoading = ref(false)
const currentQuestion = ref({})
const replyForm = reactive({
  reply: ''
})
const replyRules = {
  reply: [{ required: true, message: '请输入回复内容', trigger: 'blur' }]
}
const replyFormRef = ref(null)
const submitting = ref(false)

// 关闭问题
const closeDialogVisible = ref(false)
const questionToClose = ref(null)
const closing = ref(false)

// 获取问题列表
const fetchQuestionList = async () => {
  loading.value = true
  try {
    // 确保所有参数都是数字类型，并设置合理的默认值
    const params = {
      page: Number(page.value) || 1,
      page_size: Number(pageSize.value) || 10
    }
    
    // 仅当状态不为-1时添加状态筛选
    if (status.value !== -1 && status.value !== undefined && status.value !== null) {
      params.status = Number(status.value)
    }
    
    // 仅当类型不为-1时添加类型筛选
    if (type.value !== -1 && type.value !== undefined && type.value !== null) {
      params.type = Number(type.value)
    }
    
    // 添加调试日志，查看实际发送的参数
    console.log('发送参数:', JSON.stringify(params))
    
    const res = await getQuestionListSafe(params)
    console.log('API响应:', res)
    
    questionList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取问题列表失败', error)
    // 添加更详细的错误信息
    if (error.response) {
      console.error('错误状态码:', error.response.status)
      console.error('错误详情:', error.response.data)
      
      // 检测时间字段错误，尝试使用备用方案
      if (error.response.data && error.response.data.includes && 
          error.response.data.includes('time') && 
          error.response.data.includes('nil')) {
        
        ElMessage.warning('数据格式问题，正在尝试替代方案获取数据...')
        
        try {
          // 重新构建参数
          const fallbackParams = {
            page: Number(page.value) || 1,
            page_size: Number(pageSize.value) || 10
          }
          
          // 添加筛选条件
          if (status.value !== -1 && status.value !== undefined && status.value !== null) {
            fallbackParams.status = Number(status.value)
          }
          
          if (type.value !== -1 && type.value !== undefined && type.value !== null) {
            fallbackParams.type = Number(type.value)
          }
          
          // 添加忽略空时间字段标记
          fallbackParams.ignore_null_time = true
          
          // 使用原始方法再次尝试
          const res = await getQuestionList(fallbackParams)
          if (res && res.data) {
            questionList.value = res.data.list || []
            total.value = res.data.total || 0
            ElMessage.success('成功获取部分数据')
            return
          }
        } catch (fallbackError) {
          console.error('替代方案也失败了', fallbackError)
          // 继续显示原始错误
        }
      }
    }
    
    ElMessage.error('获取问题列表失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 获取问题详情
const fetchQuestionDetail = async (id) => {
  detailLoading.value = true
  try {
    console.log('获取问题详情，ID:', id)
    const res = await adminService({
      url: '/api/admin/service/detail',
      method: 'post',
      data: { id }
    })
    console.log('问题详情响应:', res)
    
    // 直接使用响应对象，不需要通过data属性访问
    currentQuestion.value = res
    replyForm.reply = ''
  } catch (error) {
    console.error('获取问题详情失败', error)
    // 增加更详细的错误信息
    if (error.response) {
      console.error('错误状态码:', error.response.status)
      console.error('错误详情:', error.response.data)
    }
    ElMessage.error('获取问题详情失败: ' + (error.message || '未知错误'))
  } finally {
    detailLoading.value = false
  }
}

// 显示问题详情
const handleDetail = (row) => {
  dialogVisible.value = true
  fetchQuestionDetail(row.id)
}

// 提交回复
const handleSubmitReply = async () => {
  if (!replyForm.reply.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  
  submitting.value = true
  try {
    await replyQuestion(currentQuestion.value.id, replyForm)
    ElMessage.success('回复成功')
    dialogVisible.value = false
    fetchQuestionList()
  } catch (error) {
    console.error('回复失败', error)
    ElMessage.error('回复失败')
  } finally {
    submitting.value = false
  }
}

// 处理关闭问题
const handleClose = (row) => {
  questionToClose.value = row
  closeDialogVisible.value = true
}

// 确认关闭问题
const confirmClose = async () => {
  if (!questionToClose.value) return
  
  closing.value = true
  try {
    await replyQuestion(questionToClose.value.id, { 
      status: 3,
      reply: '该问题已关闭'
    })
    ElMessage.success('问题已关闭')
    closeDialogVisible.value = false
    fetchQuestionList()
  } catch (error) {
    console.error('关闭问题失败', error)
    ElMessage.error('关闭问题失败')
  } finally {
    closing.value = false
  }
}

// 页面大小变化
const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  fetchQuestionList()
}

// 页码变化
const handleCurrentChange = (val) => {
  page.value = val
  fetchQuestionList()
}

// 状态变化
const handleStatusChange = () => {
  page.value = 1
  fetchQuestionList()
}

// 类型变化
const handleTypeChange = () => {
  page.value = 1
  fetchQuestionList()
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString()
}

// 获取问题类型文本
const getQuestionTypeText = (type) => {
  const typeMap = {
    1: '商品相关',
    2: '订单相关',
    3: '账户相关',
    4: '其他问题'
  }
  return typeMap[type] || '未知类型'
}

// 获取问题类型标签
const getQuestionTypeTag = (type) => {
  const typeMap = {
    1: 'success',
    2: 'warning',
    3: 'info',
    4: 'danger'
  }
  // 确保返回有效的标签类型，默认为 info
  return typeMap[type] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    1: '待处理',
    2: '已回复',
    3: '已关闭'
  }
  return statusMap[status] || '未知状态'
}

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    1: 'warning',
    2: 'success',
    3: 'info'
  }
  // 确保返回有效的标签类型，默认为 info
  return statusMap[status] || 'info'
}

onMounted(() => {
  fetchQuestionList()
})
</script>

<style scoped>
.question-list-container {
  padding: 16px;
}

.table-card {
  margin-bottom: 16px;
}

.search-bar {
  display: flex;
  margin-bottom: 16px;
  gap: 10px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.question-detail {
  padding: 10px;
}

.question-content, .reply-content {
  white-space: pre-wrap;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 4px;
  max-height: 200px;
  overflow-y: auto;
}

.reply-form {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}
</style> 