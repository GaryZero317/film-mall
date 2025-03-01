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
import { getQuestionList, getQuestionDetail, replyQuestion } from '../../api/customer-service'

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
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (status.value !== -1) {
      params.status = status.value
    }
    if (type.value !== -1) {
      params.type = type.value
    }
    const res = await getQuestionList(params)
    questionList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取问题列表失败', error)
    ElMessage.error('获取问题列表失败')
  } finally {
    loading.value = false
  }
}

// 获取问题详情
const fetchQuestionDetail = async (id) => {
  detailLoading.value = true
  try {
    const res = await getQuestionDetail(id)
    currentQuestion.value = res.data
    replyForm.reply = ''
  } catch (error) {
    console.error('获取问题详情失败', error)
    ElMessage.error('获取问题详情失败')
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
  return typeMap[type] || ''
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
  return statusMap[status] || ''
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