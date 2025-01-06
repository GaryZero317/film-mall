import { adminService } from './request'

// 管理员登录
export function login(data) {
  return adminService({
    url: '/api/admin/login',
    method: 'post',
    data
  })
}

// 获取管理员列表
export function getAdminList(params) {
  return adminService({
    url: '/api/admin/list',
    method: 'post',
    data: {
      page: params.page || 1,
      pageSize: params.pageSize || 10
    }
  })
}

// 创建管理员
export function createAdmin(data) {
  return adminService({
    url: '/api/admin/create',
    method: 'post',
    data
  })
}

// 更新管理员
export function updateAdmin(data) {
  return adminService({
    url: '/api/admin/update',
    method: 'post',
    data
  })
}

// 删除管理员
export function removeAdmin(data) {
  return adminService({
    url: '/api/admin/delete',
    method: 'post',
    data
  })
} 