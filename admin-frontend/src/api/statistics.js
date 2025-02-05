import request from '@/utils/request'

// 获取热门商品统计
export function getHotProducts(data) {
  return request({
    url: '/api/statistics/hot-products',
    method: 'post',
    data
  })
}

// 获取商品类别统计
export function getCategoryStats(data) {
  return request({
    url: '/api/statistics/category-stats',
    method: 'post',
    data
  })
}

// 获取用户行为统计
export function getUserBehavior(data) {
  return request({
    url: '/api/statistics/user-behavior',
    method: 'post',
    data
  })
}

// 获取用户活跃度统计
export function getUserActivity(data) {
  return request({
    url: '/api/statistics/user-activity',
    method: 'post',
    data
  })
}

// 导出统计数据
export function exportStatistics(data) {
  return request({
    url: '/api/statistics/export',
    method: 'post',
    data,
    responseType: 'blob'
  })
} 