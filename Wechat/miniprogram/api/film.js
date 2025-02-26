import request from '../utils/request'

/**
 * 获取用户胶片冲洗订单列表
 * @param {Object} params - 查询参数 {page, page_size, status}
 * @returns {Promise}
 */
export function getFilmOrderList(params = {}) {
  return request({
    url: '/api/film/user/order/list',
    method: 'get',
    data: params
  })
}

/**
 * 获取胶片冲洗订单详情
 * @param {Number} id - 订单ID
 * @returns {Promise}
 */
export function getFilmOrderDetail(id) {
  return request({
    url: `/api/film/user/order/${id}`,
    method: 'get'
  })
}

/**
 * 创建胶片冲洗订单
 * @param {Object} data - 订单数据
 * @returns {Promise}
 */
export function createFilmOrder(data) {
  return request({
    url: '/api/film/order',
    method: 'post',
    data
  })
}

/**
 * 更新胶片冲洗订单
 * @param {Number} id - 订单ID
 * @param {Object} data - 更新数据
 * @returns {Promise}
 */
export function updateFilmOrder(id, data) {
  return request({
    url: `/api/film/user/order/${id}`,
    method: 'put',
    data
  })
}

/**
 * 获取胶片冲洗服务价格列表
 * @returns {Promise}
 */
export function getFilmPriceList() {
  return request({
    url: '/api/film/prices',
    method: 'get'
  })
}

/**
 * 更新胶片订单状态（模拟支付成功后）
 * @param {Number} id - 订单ID
 * @returns {Promise}
 */
export function updateFilmOrderStatus(id) {
  return request({
    url: `/api/film/user/order/${id}`,
    method: 'put',
    data: {
      status: 1 // 更新为"冲洗处理中"状态
    }
  })
}

/**
 * 模拟支付成功后的状态更新（使用特殊API绕过普通用户权限限制）
 * 注意：此函数仅用于演示/测试环境
 * @param {Number} id - 订单ID
 * @returns {Promise}
 */
export function simulatePaymentSuccess(id) {
  return request({
    url: `/api/film/payment/simulate/${id}`,
    method: 'post'
  })
} 