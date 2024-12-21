import { paymentService } from './request'

// 创建支付
export function createPayment(data) {
  return paymentService({
    url: '/api/pay/create',
    method: 'post',
    data
  })
}

// 获取支付详情
export function getPaymentDetail(data) {
  return paymentService({
    url: '/api/pay/detail',
    method: 'post',
    data
  })
}

// 支付回调
export function paymentCallback(data) {
  return paymentService({
    url: '/api/pay/callback',
    method: 'post',
    data
  })
}

// 管理员获取支付列表
export function getAdminPaymentList(params) {
  return paymentService({
    url: '/api/admin/payment/list',
    method: 'post',
    data: {
      page: params.page || 1,
      pageSize: params.pageSize || 10
    }
  })
} 