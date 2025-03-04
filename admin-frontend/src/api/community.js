import { communityService } from './request'

/**
 * 获取作品列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export function getWorkList(params) {
  return communityService({
    url: '/api/community/admin/work/list',
    method: 'get',
    params
  })
}

/**
 * 获取作品详情
 * @param {number} id 作品ID
 * @returns {Promise}
 */
export function getWorkDetail(id) {
  console.log('发起获取作品详情请求，ID:', id)
  return communityService({
    url: `/api/community/admin/work/${id}`,
    method: 'get'
  }).then(response => {
    console.log('作品详情原始响应:', response)
    return response
  }).catch(error => {
    console.error('作品详情请求失败:', error)
    return Promise.reject(error)
  })
}

/**
 * 更新作品
 * @param {number} id 作品ID
 * @param {Object} data 更新数据
 * @returns {Promise}
 */
export function updateWork(id, data) {
  return communityService({
    url: `/api/community/admin/work/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除作品
 * @param {number} id 作品ID
 * @returns {Promise}
 */
export function deleteWork(id) {
  return communityService({
    url: `/api/community/admin/work/${id}`,
    method: 'delete'
  })
}

/**
 * 删除评论
 * @param {number} id 评论ID
 * @returns {Promise}
 */
export function deleteComment(id) {
  return communityService({
    url: `/api/community/admin/comment/${id}`,
    method: 'delete'
  })
} 