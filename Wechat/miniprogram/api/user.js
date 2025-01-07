import request from '../utils/request'

// 用户登录
export const login = (data) => {
  return request({
    url: '/user/login',
    method: 'POST',
    data
  })
}

// 获取用户信息
export const getUserInfo = () => {
  return request({
    url: '/user/info',
    method: 'GET'
  })
}

// 更新用户信息
export const updateUserInfo = (data) => {
  return request({
    url: '/user/update',
    method: 'PUT',
    data
  })
}

// 获取收货地址列表
export const getAddressList = () => {
  return request({
    url: '/address/list',
    method: 'GET'
  })
}

// 添加收货地址
export const addAddress = (data) => {
  return request({
    url: '/address/add',
    method: 'POST',
    data
  })
}

// 更新收货地址
export const updateAddress = (data) => {
  return request({
    url: '/address/update',
    method: 'PUT',
    data
  })
}

// 删除收货地址
export const deleteAddress = (id) => {
  return request({
    url: `/address/delete/${id}`,
    method: 'DELETE'
  })
}

// 设置默认收货地址
export const setDefaultAddress = (id) => {
  return request({
    url: `/address/default/${id}`,
    method: 'PUT'
  })
} 