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

// 获取收货地址列表
export const getAddressList = () => {
  return request({
    url: 'http://localhost:8000/api/address/list',
    method: 'GET'
  })
}

// 添加收货地址
export const addAddress = (data) => {
  return request({
    url: 'http://localhost:8000/api/address/add',
    method: 'POST',
    data
  })
}

// 更新收货地址
export const updateAddress = (data) => {
  return request({
    url: 'http://localhost:8000/api/address/update',
    method: 'POST',
    data
  })
}

// 删除收货地址
export const deleteAddress = (id) => {
  return request({
    url: 'http://localhost:8000/api/address/delete',
    method: 'POST',
    data: { id }
  })
}

// 设置默认收货地址
export const setDefaultAddress = (id) => {
  return request({
    url: 'http://localhost:8000/api/address/set-default',
    method: 'POST',
    data: { id }
  })
} 