import request from '../utils/request'

// 获取商品列表
export const getProductList = (params = {}) => {
  return request({
    url: '/api/product/list',
    method: 'GET',
    data: {
      page: params.page || 1,
      pageSize: params.pageSize || 10,
      keyword: params.keyword || ''
    },
    noAuth: true
  })
}

// 获取商品详情
export const getProductDetail = (id) => {
  // 参数验证
  const productId = parseInt(id)
  if (!productId || isNaN(productId) || productId <= 0) {
    return Promise.reject(new Error('无效的商品ID'))
  }

  return request({
    url: '/api/product/detail',
    method: 'POST',
    data: {
      "id": productId
    },
    header: {
      'content-type': 'application/json'
    },
    noAuth: true
  }).then(res => {
    // 如果是正常的成功响应
    if ((res.code === 0 || res.code === 200) && res.data) {
      return {
        code: 0,
        data: res.data
      }
    } else if (res.code === 200 && !res.data) {
      // 如果是直接返回的数据对象
      return {
        code: 0,
        data: res
      }
    }
    
    // 如果是商品不存在的响应
    if (res.notFound) {
      return {
        code: 404,
        msg: '商品不存在'
      }
    }
    
    // 其他错误情况
    throw new Error(res.msg || '获取商品详情失败')
  }).catch(error => {
    console.error('获取商品详情请求失败:', error.message)
    throw error
  })
}

// 创建商品
export const createProduct = (data) => {
  return request({
    url: '/api/product/create',
    method: 'POST',
    data
  })
}

// 更新商品
export const updateProduct = (data) => {
  return request({
    url: '/api/product/update',
    method: 'POST',
    data
  })
}

// 删除商品
export const removeProduct = (id) => {
  return request({
    url: '/api/product/remove',
    method: 'POST',
    data: { id }
  })
}

// 上传图片
export const uploadImage = (filePath) => {
  return new Promise((resolve, reject) => {
    wx.uploadFile({
      url: 'http://localhost:8001/api/product/upload',  // 上传文件需要完整URL
      filePath: filePath,
      name: 'file',
      success: (res) => {
        try {
          const data = JSON.parse(res.data)
          resolve(data)
        } catch (error) {
          reject(error)
        }
      },
      fail: reject
    })
  })
}

// 添加商品图片
export const addProductImages = (productId, imageUrls) => {
  return request({
    url: '/api/product/images/add',
    method: 'POST',
    data: { productId, imageUrls }
  })
}

// 设置商品主图
export const setMainImage = (productId, imageUrl) => {
  return request({
    url: '/api/product/images/setMain',
    method: 'POST',
    data: { productId, imageUrl }
  })
}

// 获取商品分类
export const getCategories = () => {
  return request({
    url: '/api/product/v1/category/list',
    method: 'POST'
  })
}

// 获取轮播图
export const getBanners = () => {
  return request({
    url: '/api/banner/v1/list',
    method: 'POST'
  })
}

// 搜索商品
export const searchProducts = (params = {}) => {
  return request({
    url: '/api/product/search',
    method: 'GET',
    data: {
      keyword: params.keyword,
      page: params.page || 1,
      pageSize: params.pageSize || 10
    },
    noAuth: true
  })
}

// 获取商品库存
export const getProductStock = (id) => {
  return request({
    url: '/api/product/stock',
    method: 'POST',
    data: { id }
  })
}

// 获取商品图片
export function getProductImages(productId) {
  return request({
    url: '/api/product/images/list',
    method: 'POST',
    data: {
      productId: Number(productId)
    },
    noAuth: true
  })
} 