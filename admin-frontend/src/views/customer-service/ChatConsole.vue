<template>
  <div class="chat-console-container">
    <div class="chat-layout">
      <!-- 左侧会话列表 -->
      <div class="session-list">
        <div class="list-header">
          <h3>用户会话列表</h3>
          <el-tag type="success" v-if="isConnected">在线</el-tag>
          <el-tag type="danger" v-else>离线</el-tag>
        </div>
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索用户ID"
            clearable
            prefix-icon="Search"
          />
        </div>
        <div class="session-items" v-loading="loadingSessions">
          <div
            v-for="session in filteredSessions"
            :key="session.id"
            class="session-item"
            :class="{ active: currentSession && currentSession.id === session.id }"
            @click="selectSession(session)"
          >
            <el-avatar :size="40">{{ session.user_id }}</el-avatar>
            <div class="session-info">
              <div class="user-id">用户ID: {{ session.user_id }}</div>
              <div class="last-message">{{ session.last_message || '暂无消息' }}</div>
            </div>
            <div class="session-meta">
              <div class="time">{{ formatTime(session.update_time) }}</div>
              <el-badge v-if="session.unread_count > 0" :value="session.unread_count" class="unread-badge" />
            </div>
          </div>
          <div v-if="filteredSessions.length === 0" class="empty-tip">
            暂无会话
          </div>
        </div>
      </div>

      <!-- 右侧聊天区域 -->
      <div class="chat-area">
        <template v-if="currentSession">
          <!-- 聊天头部 -->
          <div class="chat-header">
            <div class="user-info">
              <span class="user-id">用户ID: {{ currentSession.user_id }}</span>
            </div>
            <div class="actions">
              <el-button size="small" type="danger" @click="endSession">结束会话</el-button>
            </div>
          </div>

          <!-- 聊天消息列表 -->
          <div class="chat-messages" ref="messageContainer" v-loading="loadingMessages">
            <div v-if="messageList.length === 0" class="empty-message">
              暂无消息记录
            </div>
            <div v-else>
              <div
                v-for="(message, index) in messageList"
                :key="index"
                class="message-item"
                :class="{ 'admin-message': message.sender_type === 2, 'user-message': message.sender_type === 1 }"
              >
                <div class="message-avatar">
                  <el-avatar :size="36" :icon="message.sender_type === 1 ? 'UserFilled' : 'Service'"></el-avatar>
                </div>
                <div class="message-content">
                  <div class="message-text">{{ message.content }}</div>
                  <div class="message-time">{{ formatTime(message.create_time) }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 消息输入框 -->
          <div class="chat-input">
            <el-input
              v-model="messageInput"
              type="textarea"
              :rows="3"
              placeholder="请输入消息内容..."
              resize="none"
              @keydown.enter.exact.prevent="sendMessage"
            />
            <div class="input-actions">
              <el-button type="primary" :disabled="!messageInput.trim()" @click="sendMessage">
                发送 <el-icon class="el-icon--right"><Position /></el-icon>
              </el-button>
            </div>
          </div>
        </template>

        <!-- 未选择会话时显示 -->
        <div v-else class="no-session">
          <el-empty description="请选择一个会话开始聊天" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getChatSessions, getChatMessages, sendChatMessage } from '../../api/customer-service'
import { Position, Search, UserFilled, Service } from '@element-plus/icons-vue'

// 当前连接状态
const isConnected = ref(false)
let websocket = null

// 会话数据
const sessionList = ref([])
const loadingSessions = ref(false)
const searchKeyword = ref('')
const currentSession = ref(null)

// 消息数据
const messageList = ref([])
const loadingMessages = ref(false)
const messageInput = ref('')
const messageContainer = ref(null)

// 过滤后的会话列表
const filteredSessions = computed(() => {
  if (!searchKeyword.value) return sessionList.value
  const keyword = searchKeyword.value.toLowerCase()
  return sessionList.value.filter(session => 
    session.user_id.toString().includes(keyword)
  )
})

// 获取会话列表
const fetchChatSessions = async () => {
  loadingSessions.value = true
  try {
    const res = await getChatSessions()
    sessionList.value = res.data.list
  } catch (error) {
    console.error('获取会话列表失败', error)
    ElMessage.error('获取会话列表失败')
  } finally {
    loadingSessions.value = false
  }
}

// 获取聊天消息
const fetchChatMessages = async () => {
  if (!currentSession.value) return
  
  loadingMessages.value = true
  try {
    const res = await getChatMessages(currentSession.value.user_id, { page: 1, pageSize: 20 })
    if (res.data && res.data.list) {
      messageList.value = res.data.list.reverse() // 倒序，最新的消息在底部
    } else {
      messageList.value = []
    }
    
    // 更新当前会话的未读数
    if (currentSession.value.unread_count > 0) {
      const sessionIndex = sessionList.value.findIndex(s => s.user_id === currentSession.value.user_id)
      if (sessionIndex > -1) {
        sessionList.value[sessionIndex].unread_count = 0
      }
    }
    
    nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('获取聊天消息失败', error)
    ElMessage.error('获取聊天消息失败')
  } finally {
    loadingMessages.value = false
  }
}

// 选择会话
const selectSession = (session) => {
  currentSession.value = session
  fetchChatMessages()
}

// 发送消息
const sendMessage = async () => {
  if (!messageInput.value.trim() || !currentSession.value) return
  
  try {
    const adminInfo = JSON.parse(localStorage.getItem('adminInfo') || '{}')
    const message = {
      userId: currentSession.value.user_id,
      content: messageInput.value.trim(),
      type: 1 // 文本消息
    }
    
    // 清空输入框
    messageInput.value = ''
    
    // 调用发送接口
    await sendChatMessage(currentSession.value.user_id, {
      content: message.content,
      type: 1
    })
    
    // 本地添加消息，方便显示
    messageList.value.push({
      id: new Date().getTime(),
      userId: currentSession.value.user_id,
      adminId: adminInfo.id,
      direction: 2, // 管理员到用户
      content: message.content,
      readStatus: 0,
      createTime: Math.floor(Date.now() / 1000)
    })
    
    scrollToBottom()
  } catch (error) {
    console.error('发送消息失败', error)
    ElMessage.error('发送消息失败')
  }
}

// 结束会话
const endSession = () => {
  if (!currentSession.value) return
  
  ElMessageBox.confirm('确定要结束此会话吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // TODO: 调用结束会话接口
      const sessionIndex = sessionList.value.findIndex(s => s.id === currentSession.value.id)
      if (sessionIndex > -1) {
        sessionList.value.splice(sessionIndex, 1)
      }
      
      ElMessage.success('会话已结束')
      currentSession.value = null
      messageList.value = []
    } catch (error) {
      console.error('结束会话失败', error)
      ElMessage.error('结束会话失败')
    }
  }).catch(() => {})
}

// 连接WebSocket
const connectWebSocket = () => {
  // 获取token
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('请先登录')
    return
  }
  
  // 获取管理员信息
  const adminInfo = JSON.parse(localStorage.getItem('adminInfo') || '{}')
  const adminId = adminInfo.id
  
  if (!adminId) {
    ElMessage.error('无法获取管理员信息')
    return
  }
  
  // 创建WebSocket连接，使用正确的路径
  const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8000'
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsBaseUrl = baseUrl.replace(/^https?:/, wsProtocol)
  
  websocket = new WebSocket(`${wsBaseUrl}/api/admin/chat/connect?adminId=${adminId}&token=${token}`)
  
  // 连接开启
  websocket.onopen = () => {
    isConnected.value = true
    ElMessage.success('客服系统已连接')
    
    // 发送心跳包
    startHeartbeat()
  }
  
  // 接收消息
  websocket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      
      // 处理消息类型
      if (data.type === 'message') {
        handleNewMessage(data.data)
      } else if (data.type === 'session') {
        handleSessionUpdate(data.data)
      } else if (data.type === 'pong') {
        // 心跳响应，不做处理
      }
    } catch (error) {
      console.error('处理消息失败', error)
    }
  }
  
  // 连接关闭
  websocket.onclose = () => {
    isConnected.value = false
    stopHeartbeat()
    
    // 尝试重连
    setTimeout(() => {
      ElMessage.warning('正在尝试重新连接...')
      connectWebSocket()
    }, 3000)
  }
  
  // 连接错误
  websocket.onerror = (error) => {
    console.error('WebSocket错误', error)
    isConnected.value = false
    stopHeartbeat()
  }
}

// 处理新消息
const handleNewMessage = (message) => {
  // 如果是当前会话的消息，直接添加到消息列表
  if (currentSession.value && message.session_id === currentSession.value.id) {
    messageList.value.push(message)
    scrollToBottom()
  }
  
  // 更新会话列表中的最后一条消息
  const sessionIndex = sessionList.value.findIndex(s => s.id === message.session_id)
  if (sessionIndex > -1) {
    const session = sessionList.value[sessionIndex]
    session.last_message = message.content
    session.update_time = message.create_time
    
    // 如果不是当前会话，增加未读数
    if (!currentSession.value || currentSession.value.id !== message.session_id) {
      session.unread_count = (session.unread_count || 0) + 1
    }
    
    // 将会话移到最前面
    sessionList.value.splice(sessionIndex, 1)
    sessionList.value.unshift(session)
  } else {
    // 如果会话不存在，说明是新会话，刷新会话列表
    fetchChatSessions()
  }
}

// 处理会话更新
const handleSessionUpdate = (session) => {
  const sessionIndex = sessionList.value.findIndex(s => s.id === session.id)
  if (sessionIndex > -1) {
    sessionList.value[sessionIndex] = session
  } else {
    sessionList.value.unshift(session)
  }
}

// 心跳包
let heartbeatTimer = null
const startHeartbeat = () => {
  heartbeatTimer = setInterval(() => {
    if (websocket && isConnected.value) {
      websocket.send(JSON.stringify({ type: 'ping' }))
    }
  }, 30000) // 30秒一次心跳
}

const stopHeartbeat = () => {
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  const date = new Date(timestamp * 1000)
  
  // 今天的消息只显示时间
  const now = new Date()
  const isSameDay = date.getDate() === now.getDate() && 
                   date.getMonth() === now.getMonth() && 
                   date.getFullYear() === now.getFullYear()
  
  if (isSameDay) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else {
    return date.toLocaleString('zh-CN', { 
      month: '2-digit', 
      day: '2-digit',
      hour: '2-digit', 
      minute: '2-digit'
    })
  }
}

// 监听消息列表变化，自动滚动到底部
watch(messageList, () => {
  scrollToBottom()
})

onMounted(() => {
  fetchChatSessions()
  connectWebSocket()
})

onUnmounted(() => {
  // 关闭WebSocket连接
  if (websocket) {
    websocket.close()
  }
  
  // 停止心跳
  stopHeartbeat()
})
</script>

<style scoped>
.chat-console-container {
  padding: 16px;
  height: calc(100vh - 120px);
  min-height: 500px;
}

.chat-layout {
  display: flex;
  background-color: #fff;
  height: 100%;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

/* 左侧会话列表 */
.session-list {
  width: 280px;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
}

.list-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ebeef5;
}

.list-header h3 {
  margin: 0;
  font-size: 16px;
}

.search-box {
  padding: 10px;
  border-bottom: 1px solid #ebeef5;
}

.session-items {
  flex: 1;
  overflow-y: auto;
}

.session-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  transition: background-color 0.3s;
  position: relative;
}

.session-item:hover {
  background-color: #f5f7fa;
}

.session-item.active {
  background-color: #ecf5ff;
}

.session-info {
  flex: 1;
  margin-left: 12px;
  overflow: hidden;
}

.user-id {
  font-weight: bold;
  margin-bottom: 4px;
}

.last-message {
  color: #606266;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.session-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.time {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.unread-badge {
  margin-top: 4px;
}

.empty-tip {
  text-align: center;
  padding: 20px;
  color: #909399;
}

/* 右侧聊天区域 */
.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 16px;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-messages {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.message-item {
  display: flex;
  margin-bottom: 16px;
}

.admin-message {
  flex-direction: row-reverse;
}

.message-avatar {
  margin: 0 8px;
}

.message-content {
  max-width: 60%;
}

.user-message .message-content {
  background-color: #f2f2f2;
  border-radius: 8px;
  padding: 10px;
}

.admin-message .message-content {
  background-color: #ecf5ff;
  border-radius: 8px;
  padding: 10px;
  text-align: right;
}

.message-text {
  word-break: break-word;
  white-space: pre-wrap;
}

.message-time {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.chat-input {
  padding: 16px;
  border-top: 1px solid #ebeef5;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.no-session {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.empty-message {
  text-align: center;
  color: #909399;
  padding: 20px;
}
</style> 