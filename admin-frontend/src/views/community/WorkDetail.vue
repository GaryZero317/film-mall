<template>
  <div class="container" v-loading="loading">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="goBack">返回</el-button>
        <h2>作品详情</h2>
      </div>
    </div>

    <!-- 调试信息 -->
    <el-alert
      v-if="debugMode"
      title="调试信息"
      type="info"
      :closable="false"
      show-icon
    >
      <pre>作品ID: {{ workId }}</pre>
      <pre>原始数据: {{ JSON.stringify(rawResponse, null, 2) }}</pre>
      <pre>处理后数据: {{ JSON.stringify(workData, null, 2) }}</pre>
    </el-alert>

    <el-card v-if="workData">
      <template #header>
        <div class="card-header">
          <span>{{ workData.title || '未命名作品' }}</span>
          <el-tag :type="getStatusTagType(workData.status)">{{ getStatusText(workData.status) }}</el-tag>
        </div>
      </template>
      <div class="card-content">
        <p><strong>作品ID:</strong> {{ workData.id || '无' }}</p>
        <p><strong>用户ID:</strong> {{ workData.uid || '无' }}</p>
        <p><strong>用户名:</strong> {{ workData.name || workData.author?.nickname || '未知用户' }}</p>
        <p><strong>创建时间:</strong> {{ workData.create_time || '无' }}</p>
        <p><strong>描述:</strong> {{ workData.description || '无描述' }}</p>
        <p v-if="workData.film_type"><strong>胶片类型:</strong> {{ workData.film_type }}</p>
        <p v-if="workData.film_brand"><strong>胶片品牌:</strong> {{ workData.film_brand }}</p>
        
        <!-- 显示封面图 -->
        <div v-if="workData.cover_url" class="cover-image">
          <p><strong>封面图:</strong></p>
          <el-image 
            :src="getImageUrl(workData.cover_url)" 
            fit="contain"
            :preview-src-list="[getImageUrl(workData.cover_url)]"
            style="max-width: 300px; max-height: 300px">
          </el-image>
        </div>
        
        <!-- 互动数据 -->
        <div class="interaction-data">
          <p><strong>互动数据:</strong></p>
          <p>浏览: {{ workData.view_count || 0 }} | 点赞: {{ workData.like_count || 0 }} | 评论: {{ workData.comment_count || 0 }}</p>
        </div>

        <!-- 评论列表 -->
        <div class="comments-section">
          <h3>评论列表</h3>
          <div v-loading="commentsLoading">
            <div v-if="comments.length > 0" class="comment-list">
              <div v-for="comment in comments" :key="comment.id" class="comment-item">
                <div class="comment-header">
                  <span class="comment-user">{{ comment.user?.nickname || '未知用户' }}</span>
                  <span class="comment-time">{{ comment.create_time }}</span>
                </div>
                <div class="comment-content">{{ comment.content }}</div>
              </div>
            </div>
            <el-empty v-else description="暂无评论"></el-empty>
          </div>
        </div>
      </div>
    </el-card>
    <div v-else class="no-data">
      <el-empty description="暂无数据"></el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getWorkDetail, getCommentList } from '../../api/community'

const route = useRoute()
const router = useRouter()
const workId = ref(parseInt(route.params.id) || 0)
const loading = ref(false)
const workData = ref(null)
const debugMode = ref(false) // 关闭调试模式
const rawResponse = ref(null)
const commentsLoading = ref(false)
const comments = ref([])

// 获取状态标签类型
const getStatusTagType = (status) => {
  if (status === 0) return 'info'
  if (status === 1) return 'success'
  if (status === 2) return 'danger'
  return 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  if (status === 0) return '草稿'
  if (status === 1) return '已发布'
  if (status === 2) return '已删除'
  return '未知'
}

// 获取完整的图片URL
const getImageUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return `http://localhost:8008${url}`
}

// 加载作品详情
const loadWorkDetail = () => {
  loading.value = true
  console.log('正在加载作品详情，ID:', workId.value)
  
  getWorkDetail(workId.value).then(res => {
    console.log('获取到的作品详情数据:', res)
    
    // 增强数据处理逻辑
    if (res.data) {
      if (res.data.work) {
        // 如果是嵌套结构
        console.log('处理嵌套结构的作品数据')
        workData.value = {
          ...res.data.work, // 直接使用work对象的数据
          author: res.data.author // 直接使用data中的author
        }
        console.log('处理后的作品数据:', workData.value)
      } else {
        // 如果是扁平结构
        console.log('处理扁平结构的作品数据')
        workData.value = res.data
      }
    } else if (res.work) {
      // 另一种可能的响应结构
      console.log('直接使用响应中的work数据')
      workData.value = {
        ...res.work,
        author: res.author
      }
    } else {
      // 如果没有符合预期的数据结构
      console.error('未找到符合预期的数据结构:', res)
      ElMessage.warning('获取作品详情数据结构异常')
    }
    rawResponse.value = res
  }).catch(err => {
    console.error('获取作品详情失败', err)
    ElMessage.error('获取作品详情失败')
  }).finally(() => {
    loading.value = false
  })
}

// 加载评论列表
const loadComments = async () => {
  if (!workId.value) return
  
  commentsLoading.value = true
  try {
    const res = await getCommentList(workId.value)
    comments.value = res.data.list || []
  } catch (err) {
    console.error('加载评论失败:', err)
    ElMessage.error('加载评论失败')
  } finally {
    commentsLoading.value = false
  }
}

// 返回列表页
const goBack = () => {
  router.push('/community/works')
}

// 页面加载时获取数据
onMounted(() => {
  loadWorkDetail()
  loadComments() // 加载评论列表
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

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-content {
  line-height: 1.8;
}

.no-data {
  padding: 40px 0;
}

.comments-section {
  margin-top: 20px;
}

.comment-list {
  margin-top: 10px;
}

.comment-item {
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
  color: #666;
  font-size: 14px;
}

.comment-content {
  color: #333;
}

.comment-user {
  font-weight: bold;
}

.comment-time {
  color: #999;
}
</style> 