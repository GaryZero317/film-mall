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
      uid: params.uid,
      status: params.status || 0,
      page: params.page || 1,
      page_size: params.pageSize || 10
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

// 查询支付状态
export const queryPayStatus = (payId) => {
  return request({
    url: 'http://localhost:8003/api/pay/detail',
    method: 'POST',
    data: { id: payId }
  })
}

// 查询订单支付记录
export function queryPayByOrderId(orderId) {
  return request({
    url: 'http://localhost:8003/api/pay/query/order',
    method: 'POST',
    data: {
      orderId
    }
  })
}

// 查询支付参数
export const getPayParams = (payId) => {
  return request({
    url: 'http://localhost:8003/api/pay/params',
    method: 'POST',
    data: { id: payId }
  })
}

// 支付订单
export const payOrder = (data) => {
  return request({
    url: 'http://localhost:8003/api/pay/create',
    method: 'POST',
    data: {
      oid: data.oid,
      uid: data.uid,
      amount: data.amount
    }
  })
}

// 查询订单支付状态
export const getOrderPayStatus = (orderId) => {
  return request({
    url: 'http://localhost:8003/api/pay/detail',
    method: 'POST',
    data: { orderId }
  })
} 