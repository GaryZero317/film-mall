import request from '../utils/request'

// 创建订单
export const createOrder = (data) => {
  return request({
    url: '/order/create',
    method: 'POST',
    data
  })
}

// 获取订单列表
export const getOrderList = (status) => {
  return request({
    url: '/order/list',
    method: 'GET',
    data: { status }
  })
}

// 获取订单详情
export const getOrderDetail = (id) => {
  return request({
    url: `/order/detail/${id}`,
    method: 'GET'
  })
}

// 取消订单
export const cancelOrder = (id) => {
  return request({
    url: `/order/cancel/${id}`,
    method: 'PUT'
  })
}

// 确认收货
export const confirmOrder = (id) => {
  return request({
    url: `/order/confirm/${id}`,
    method: 'PUT'
  })
}

// 删除订单
export const deleteOrder = (id) => {
  return request({
    url: `/order/delete/${id}`,
    method: 'DELETE'
  })
}

// 支付订单
export const payOrder = (id) => {
  return request({
    url: `/order/pay/${id}`,
    method: 'POST'
  })
} 