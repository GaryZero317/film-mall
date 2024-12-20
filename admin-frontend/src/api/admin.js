import { userService } from './request'

// 管理员登录
export function login(data) {
  return userService({
    url: '/api/admin/login',
    method: 'post',
    data
  })
}

// 创建管理员
export function createAdmin(data) {
  return userService({
    url: '/api/admin/create',
    method: 'post',
    data
  })
}

// 更新管理员
export function updateAdmin(data) {
  return userService({
    url: '/api/admin/update',
    method: 'post',
    data
  })
}

// 删除管理员
export function deleteAdmin(data) {
  return userService({
    url: '/api/admin/delete',
    method: 'post',
    data
  })
}

// 获取管理员信息
export function getAdminInfo() {
  return userService({
    url: '/api/admin/info',
    method: 'post'
  })
} 