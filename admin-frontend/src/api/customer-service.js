import request from '../utils/request'
import { adminService } from './request'

// 获取客服问题列表 - 原始实现
export function getQuestionList(params) {
  return request({
    url: '/api/admin/service/list',
    method: 'post',
    data: params
  })
}

// 获取客服问题列表 - 使用专用服务实例
export function getQuestionListSafe(params) {
  return adminService({
    url: '/api/admin/service/list',
    method: 'post',
    data: params
  })
}

// 获取问题详情
export function getQuestionDetail(id) {
  return request({
    url: '/api/admin/service/detail',
    method: 'post',
    data: {
      id: id
    }
  })
}

// 回复客服问题 - 使用服务更新API
export function replyQuestion(id, data) {
  return request({
    url: '/api/admin/service/update',
    method: 'post',
    data: {
      id: id,
      reply: data.reply,
      status: 1  // 已回复状态
    }
  })
}

// 获取FAQ列表
export function getFaqList(params) {
  return adminService({
    url: '/api/user/service/faq/list',
    method: 'post',
    data: params
  })
}

// 添加FAQ - 注意：后端可能没有此API
export function addFaq(data) {
  return request({
    url: '/api/admin/service/faq/add',
    method: 'post',
    data
  })
}

// 更新FAQ - 注意：后端可能没有此API
export function updateFaq(id, data) {
  return request({
    url: '/api/admin/service/faq/update',
    method: 'post',
    data: {
      id,
      ...data
    }
  })
}

// 删除FAQ - 注意：后端可能没有此API
export function deleteFaq(id) {
  return request({
    url: '/api/admin/service/faq/delete',
    method: 'post',
    data: {
      id
    }
  })
}

// 获取聊天会话列表 (保持管理员API路径)
export function getChatSessions(params) {
  return request({
    url: '/api/admin/chat/sessions',
    method: 'get',
    params
  })
}

// 获取聊天记录 (保持管理员API路径)
export function getChatMessages(userId, params) {
  return request({
    url: `/api/admin/chat/history`,
    method: 'get',
    params: {
      userId,
      ...params
    }
  })
}

// 发送聊天消息 (保持管理员API路径)
export function sendChatMessage(userId, data) {
  return request({
    url: `/api/admin/chat/send`,
    method: 'post',
    data: {
      userId,
      ...data
    }
  })
} 