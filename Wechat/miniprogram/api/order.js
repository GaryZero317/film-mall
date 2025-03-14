import request from '../utils/request'

// 创建订单
export const createOrder = (data) => {
  // 确保支付状态初始为未支付
  if (data && !data.hasOwnProperty('pay_status')) {
    data.pay_status = 0 // 0表示未支付
  }
  
  console.log('[订单API] 创建订单，请求参数:', JSON.stringify(data))
  
  return request({
    url: '/api/order',
    method: 'POST',
    data
  }).then(res => {
    console.log('[订单API] 创建订单成功:', res)
    return res
  }).catch(err => {
    console.error('[订单API] 创建订单失败:', err)
    throw err
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
  console.log('调用订单详情API, 订单ID:', id)
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
  // 不再自动设置status，使用调用者提供的值
  console.log('[支付API] 调用支付回调接口，参数:', JSON.stringify(data))
  return request({
    url: '/api/pay/callback',
    method: 'POST',
    data
  }).then(res => {
    console.log('[支付API] 支付回调成功:', res)
    return res
  }).catch(err => {
    console.error('[支付API] 支付回调失败:', err)
    // 仍然将错误抛出，由上层处理
    throw err
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
    url: `/api/order/${orderId}`,
    method: 'PUT',
    data: {
      status: 1  // 1表示已支付，待发货
    }
  })
} 