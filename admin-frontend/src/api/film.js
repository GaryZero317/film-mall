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

// 上传胶片照片
export function uploadFilmPhoto(data) {
  return filmService({
    url: '/api/film/admin/photo/upload',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 获取胶片照片列表
export function getFilmPhotos(filmOrderId) {
  return filmService({
    url: '/api/film/admin/photo/list',
    method: 'get',
    params: { film_order_id: filmOrderId }
  })
}

// 删除胶片照片
export function deleteFilmPhoto(photoId) {
  return filmService({
    url: `/api/film/admin/photo/${photoId}`,
    method: 'delete'
  })
} 