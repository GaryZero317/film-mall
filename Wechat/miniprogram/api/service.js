import request from '../utils/request'

// 获取常见问题列表
export const getFaqList = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/service/faq/list',
    method: 'POST',
    data
  })
}

// 创建客服问题
export const createQuestion = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/service/create',
    method: 'POST',
    data
  })
}

// 获取用户的问题列表
export const getQuestionList = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/service/list',
    method: 'POST',
    data
  })
}

// 获取问题详情
export const getQuestionDetail = (id) => {
  return request({
    url: 'http://localhost:8000/api/user/service/detail',
    method: 'POST',
    data: { id }
  })
}

// 获取聊天会话，如果不存在则创建
export const getChatSession = () => {
  return request({
    url: 'http://localhost:8000/api/user/chat/session',
    method: 'GET'
  })
}

// 获取聊天历史记录
export const getChatHistory = (sessionId, page = 1, pageSize = 20) => {
  return request({
    url: 'http://localhost:8000/api/user/chat/history',
    method: 'GET',
    data: {
      sessionId,
      page,
      pageSize
    }
  })
}

// 发送聊天消息
export const sendChatMessage = (sessionId, content, type = 1) => {
  return request({
    url: 'http://localhost:8000/api/user/chat/send',
    method: 'POST',
    data: {
      sessionId,
      content,
      type // 1: 文本, 2: 图片, 3: 语音
    }
  })
} 