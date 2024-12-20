import { orderService } from './request'

// 创建订单
export function createOrder(data) {
  return orderService({
    url: '/api/admin/order/create',
    method: 'post',
    data
  })
}

// 更新订单
export function updateOrder(data) {
  return orderService({
    url: '/api/admin/order/update',
    method: 'post',
    data
  })
}

// 删除订单
export function removeOrder(data) {
  return orderService({
    url: '/api/admin/order/remove',
    method: 'post',
    data
  })
}

// 获取订单详情
export function getOrderDetail(data) {
  return orderService({
    url: '/api/admin/order/detail',
    method: 'post',
    data
  })
}

// 获取订单列表
export function getOrderList(params) {
  return orderService({
    url: '/api/admin/order/list',
    method: 'post',
    data: {
      uid: params.uid || 0,
      page: params.page || 1,
      pageSize: params.pageSize || 10
    }
  })
}

// 更新订单状态
export function updateOrderStatus(data) {
  return orderService({
    url: '/api/admin/order/status',
    method: 'post',
    data
  })
} 