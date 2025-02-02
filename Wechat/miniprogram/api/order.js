import request from '../utils/request'

// 创建订单
export const createOrder = (data) => {
  return request({
    url: '/api/order',
    method: 'POST',
    data
  })
}

// 获取订单列表
export const getOrderList = (params = {}) => {
  return request({
    url: '/api/order/list',
    method: 'GET',
    params: {
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
    url: `/api/order/${id}`,
    method: 'GET'
  })
}

// 获取地址详情
export const getAddressDetail = (id) => {
  return request({
    url: `/api/address/${id}`,
    method: 'GET'
  })
}

// 取消订单
export const cancelOrder = (id) => {
  return request({
    url: '/api/order/cancel',
    method: 'POST',
    data: { id }
  })
}

// 确认收货
export const confirmOrder = (id) => {
  return request({
    url: '/api/order/confirm',
    method: 'POST',
    data: { id }
  })
}

// 删除订单
export const deleteOrder = (id) => {
  return request({
    url: '/api/order/delete',
    method: 'POST',
    data: { id }
  })
}

// 获取订单统计
export const getOrderCount = () => {
  return request({
    url: '/api/order/count',
    method: 'GET'
  })
}

// 查询支付状态
export const queryPayStatus = (payId) => {
  return request({
    url: '/api/pay/detail',
    method: 'POST',
    data: { id: payId }
  })
}

// 查询订单支付记录
export function queryPayByOrderId(orderId) {
  return request({
    url: '/api/pay/query/order',
    method: 'POST',
    data: {
      orderId
    }
  })
}

// 查询支付参数
export const getPayParams = (payId) => {
  return request({
    url: '/api/pay/params',
    method: 'POST',
    data: { id: payId }
  })
}

// 支付订单
export const payOrder = (data) => {
  return request({
    baseUrl: 'pay',
    url: '/api/pay/create',
    method: 'POST',
    data: {
      oid: data.oid,
      uid: data.uid,
      amount: data.amount
    }
  })
}

// 支付回调
export const payCallback = (data) => {
  return request({
    baseUrl: 'pay',
    url: '/api/pay/callback',
    method: 'POST',
    data
  })
}

// 查询订单支付状态
export const getOrderPayStatus = (orderId) => {
  return request({
    url: '/api/pay/detail',
    method: 'POST',
    data: { orderId }
  })
}

// 更新订单状态为已支付
export function updateOrderStatus(orderId) {
  return request({
    url: '/api/order/update',
    method: 'POST',
    data: {
      id: orderId,
      status: 1,  // 1表示已支付，待发货
      status_desc: '待发货'
    }
  })
} 