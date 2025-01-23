import request from '../utils/request'

// 创建支付订单
export const createPayment = (orderId) => {
  return request({
    url: 'http://localhost:8003/api/pay/create',
    method: 'POST',
    data: { orderId }
  })
}

// 查询支付状态
export const queryPayment = (orderId) => {
  return request({
    url: 'http://localhost:8003/api/pay/query',
    method: 'POST',
    data: { orderId }
  })
}

// 取消支付
export const cancelPayment = (orderId) => {
  return request({
    url: 'http://localhost:8003/api/pay/cancel',
    method: 'POST',
    data: { orderId }
  })
}

// 申请退款
export const refund = (orderId, reason) => {
  return request({
    url: 'http://localhost:8003/api/pay/refund',
    method: 'POST',
    data: { 
      orderId,
      reason
    }
  })
}

// 查询退款状态
export const queryRefund = (orderId) => {
  return request({
    url: 'http://localhost:8003/api/pay/refund/query',
    method: 'POST',
    data: { orderId }
  })
} 