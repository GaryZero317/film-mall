import { payService } from './request'

// 创建支付
export function createPayment(data) {
  return payService({
    url: '/api/pay/create',
    method: 'post',
    data
  })
}

// 获取支付详情
export function getPaymentDetail(data) {
  return payService({
    url: '/api/pay/detail',
    method: 'post',
    data
  })
}

// 支付回调
export function paymentCallback(data) {
  return payService({
    url: '/api/pay/callback',
    method: 'post',
    data
  })
} 