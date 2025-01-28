import { orderService } from './request'

// 创建订单
export function createOrder(data) {
  return orderService({
    url: '/api/order',
    method: 'post',
    data
  })
}

// 更新订单
export function updateOrder(id, data) {
  return orderService({
    url: `/api/order/${id}`,
    method: 'put',
    data
  })
}

// 删除订单
export function deleteOrder(id) {
  return orderService({
    url: `/api/order/${id}`,
    method: 'delete'
  })
}

// 获取订单详情
export function getOrderDetail(id) {
  return orderService({
    url: `/api/order/${id}`,
    method: 'get'
  })
}

// 获取订单列表（用户端）
export function getOrderList(params) {
  return orderService({
    url: '/api/order/list',
    method: 'get',
    params: {
      page: params.page || 1,
      page_size: params.pageSize || 10,
      status: params.status,
      uid: params.uid
    }
  })
}

// 获取订单列表（管理端）
export function getAdminOrderList(params) {
  return orderService({
    url: '/api/order/list',
    method: 'get',
    params: {
      page: params.page || 1,
      page_size: params.pageSize || 10,
      status: params.status || -1,  // -1表示查询所有状态
      uid: 0  // 管理员查看所有订单，uid设为0
    }
  })
}

// 更新订单状态
export function updateOrderStatus(id, status) {
  return orderService({
    url: `/api/order/${id}`,
    method: 'put',
    data: {
      status
    }
  })
} 