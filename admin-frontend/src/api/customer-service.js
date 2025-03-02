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
      status: data.status || 2  // 修改为已回复状态(2)，同时保持对手动关闭(3)的支持
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
  console.log('调用getChatSessions API，参数:', params)
  // 使用adminService而不是基本request，以便获得特定错误处理
  return adminService({
    url: '/api/admin/chat/sessions',
    method: 'get',
    params,
    timeout: 5000 // 设置较短的超时
  }).catch(error => {
    console.error('getChatSessions API调用失败:', error)
    // 构造一个默认响应，避免前端崩溃
    return {
      code: 0,
      msg: 'API暂时不可用，返回模拟数据',
      list: []
    }
  })
}

// 获取聊天记录 (保持管理员API路径)
export function getChatMessages(userId, params) {
  console.log('调用getChatMessages API，参数:', { userId, ...params })
  return adminService({
    url: `/api/admin/chat/history`,
    method: 'get',
    params: {
      userId,
      ...params
    },
    timeout: 5000 // 设置较短的超时
  }).catch(error => {
    console.error('getChatMessages API调用失败:', error)
    // 构造一个默认响应，避免前端崩溃
    return {
      code: 0,
      msg: 'API暂时不可用，返回模拟数据',
      data: {
        total: 0,
        list: []
      }
    }
  })
}

// 发送聊天消息 (保持管理员API路径)
export function sendChatMessage(userId, data) {
  console.log('调用sendChatMessage API，参数:', { userId, ...data })
  return adminService({
    url: `/api/admin/chat/send`,
    method: 'post',
    data: {
      userId,
      ...data
    },
    timeout: 5000 // 设置较短的超时
  }).catch(error => {
    console.error('sendChatMessage API调用失败:', error)
    // 返回错误信息并且警告用户消息可能未发送成功
    throw new Error('消息发送失败，请稍后重试')
  })
} 