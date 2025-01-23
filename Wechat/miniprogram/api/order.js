import request from '../utils/request'

// 创建订单
export const createOrder = (data) => {
  return request({
    url: 'http://localhost:8002/api/order/create',
    method: 'POST',
    data
  })
}

// 获取订单列表
export const getOrderList = (params = {}) => {
  return request({
    url: 'http://localhost:8002/api/order/list',
    method: 'POST',
    data: {
      page: params.page || 1,
      pageSize: params.pageSize || 10,
      status: params.status
    }
  })
}

// 获取订单详情
export const getOrderDetail = (id) => {
  return request({
    url: 'http://localhost:8002/api/order/detail',
    method: 'POST',
    data: { id }
  })
}

// 取消订单
export const cancelOrder = (id) => {
  return request({
    url: 'http://localhost:8002/api/order/cancel',
    method: 'POST',
    data: { id }
  })
}

// 确认收货
export const confirmOrder = (id) => {
  return request({
    url: 'http://localhost:8002/api/order/confirm',
    method: 'POST',
    data: { id }
  })
}

// 删除订单
export const deleteOrder = (id) => {
  return request({
    url: 'http://localhost:8002/api/order/delete',
    method: 'POST',
    data: { id }
  })
}

// 获取订单统计
export const getOrderCount = () => {
  return request({
    url: 'http://localhost:8002/api/order/count',
    method: 'GET'
  })
} 