import { filmService } from './request'

// 获取胶片冲洗订单列表
export function getFilmOrderList(params) {
  return filmService({
    url: '/api/film/admin/order/list',
    method: 'get',
    params
  })
}

// 获取胶片冲洗订单详情
export function getFilmOrderDetail(id) {
  return filmService({
    url: `/api/film/admin/order/${id}`,
    method: 'get'
  })
}

// 更新胶片冲洗订单
export function updateFilmOrder(id, data) {
  return filmService({
    url: `/api/film/admin/order/${id}`,
    method: 'put',
    data
  })
}

// 删除胶片冲洗订单
export function deleteFilmOrder(id) {
  return filmService({
    url: `/api/film/admin/order/${id}`,
    method: 'delete'
  })
} 