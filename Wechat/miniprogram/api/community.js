import request from '../utils/request'

/**
 * 获取社区作品列表
 * @param {Object} params 查询参数 
 */
export function getWorkList(params = {}) {
  return request({
    url: '/api/community/work/list',
    method: 'GET',
    params
  })
}

/**
 * 获取社区作品详情
 * @param {number} id 作品ID
 */
export function getWorkDetail(id) {
  if (!id) {
    return Promise.reject(new Error('作品ID不能为空'))
  }
  return request({
    url: `/api/community/work/${id}`,
    method: 'GET'
  })
}

/**
 * 获取用户个人作品列表
 */
export function getUserWorks() {
  return request({
    url: '/api/community/user/work/list',
    method: 'GET'
  })
}

/**
 * 作品点赞/取消点赞
 * @param {number} id 作品ID
 * @param {boolean} isLike 是否点赞
 */
export function likeWork(id, isLike) {
  return request({
    url: '/api/community/user/like',
    method: 'POST',
    data: { work_id: id, action: isLike ? 1 : 0 }
  })
}

/**
 * 获取作品评论列表
 * @param {number} workId 作品ID
 */
export function getComments(workId) {
  if (!workId) {
    return Promise.reject(new Error('作品ID不能为空'))
  }
  return request({
    url: '/api/community/comment/list',
    method: 'GET',
    params: { work_id: workId }
  })
}

/**
 * 发表评论
 * @param {number} workId 作品ID
 * @param {string} content 评论内容
 */
export function addComment(workId, content) {
  return request({
    url: '/api/community/user/comment',
    method: 'POST',
    data: { work_id: workId, content }
  })
}

/**
 * 删除评论
 * @param {number} commentId 评论ID
 */
export function deleteComment(commentId) {
  return request({
    url: `/api/community/user/comment/${commentId}`,
    method: 'DELETE'
  })
}

/**
 * 上传社区作品
 * @param {Object} data 作品数据
 */
export function uploadWork(data) {
  return request({
    url: '/api/community/user/work',
    method: 'POST',
    data
  })
}

/**
 * 更新社区作品
 * @param {number} id 作品ID
 * @param {Object} data 更新数据
 */
export function updateWork(id, data) {
  return request({
    url: `/api/community/user/work/${id}`,
    method: 'PUT',
    data
  })
}

/**
 * 删除社区作品
 * @param {number} id 作品ID
 */
export function deleteWork(id) {
  return request({
    url: `/api/community/user/work/${id}`,
    method: 'DELETE'
  })
}

/**
 * 上传作品图片
 * @param {string} filePath 图片文件路径
 * @param {number} workId 作品ID
 */
export function uploadWorkImage(filePath, workId) {
  return new Promise((resolve, reject) => {
    wx.uploadFile({
      url: `${getApp().globalData.baseUrl.community}/api/community/user/work/image/upload`,
      filePath,
      name: 'file',
      formData: {
        work_id: workId
      },
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token')}`
      },
      success: (res) => {
        if (res.statusCode === 200) {
          resolve(JSON.parse(res.data))
        } else {
          reject(new Error(res.data))
        }
      },
      fail: reject
    })
  })
}

/**
 * 删除作品图片
 * @param {number} imageId 图片ID
 */
export function deleteWorkImage(imageId) {
  return request({
    url: `/api/community/user/work/image/${imageId}`,
    method: 'DELETE'
  })
} 