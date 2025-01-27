import request from '../utils/request'

// 用户登录
export const login = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/login',
    method: 'POST',
    data,
    noAuth: true
  })
}

// 微信登录
export const wxLogin = (code) => {
  return request({
    url: 'http://localhost:8000/api/user/wx-login',
    method: 'POST',
    data: { code },
    noAuth: true
  })
}

// 用户注册
export const register = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/register',
    method: 'POST',
    data,
    noAuth: true
  })
}

// 获取用户信息
export const getUserInfo = () => {
  return request({
    url: 'http://localhost:8000/api/user/userinfo',
    method: 'POST'
  })
}

// 更新用户信息
export const updateUserInfo = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/update',
    method: 'POST',
    data
  })
} 