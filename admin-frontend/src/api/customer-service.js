import request from '../utils/request'

// 获取客服问题列表
export function getQuestionList(params) {
  return request({
    url: '/customer-service/questions',
    method: 'get',
    params
  })
}

// 获取问题详情
export function getQuestionDetail(id) {
  return request({
    url: `/customer-service/questions/${id}`,
    method: 'get'
  })
}

// 回复客服问题
export function replyQuestion(id, data) {
  return request({
    url: `/customer-service/questions/${id}/reply`,
    method: 'post',
    data
  })
}

// 获取FAQ列表
export function getFaqList(params) {
  return request({
    url: '/customer-service/faq',
    method: 'get',
    params
  })
}

// 添加FAQ
export function addFaq(data) {
  return request({
    url: '/customer-service/faq',
    method: 'post',
    data
  })
}

// 更新FAQ
export function updateFaq(id, data) {
  return request({
    url: `/customer-service/faq/${id}`,
    method: 'put',
    data
  })
}

// 删除FAQ
export function deleteFaq(id) {
  return request({
    url: `/customer-service/faq/${id}`,
    method: 'delete'
  })
}

// 获取聊天会话列表
export function getChatSessions(params) {
  return request({
    url: '/customer-service/chat/sessions',
    method: 'get',
    params
  })
}

// 获取聊天记录
export function getChatMessages(sessionId, params) {
  return request({
    url: `/customer-service/chat/sessions/${sessionId}/messages`,
    method: 'get',
    params
  })
}

// 发送聊天消息
export function sendChatMessage(sessionId, data) {
  return request({
    url: `/customer-service/chat/sessions/${sessionId}/messages`,
    method: 'post',
    data
  })
} 