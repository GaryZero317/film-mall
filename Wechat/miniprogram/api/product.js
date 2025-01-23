import request from '../utils/request'

// 获取商品列表
export const getProductList = (data) => {
  return request({
    url: 'http://localhost:8001/api/product/list',
    method: 'GET',
    data: {
      page: data.page || 1,
      pageSize: data.pageSize || 10
    }
  })
}

// 获取商品详情
export const getProductDetail = (id) => {
  return request({
    url: 'http://localhost:8001/api/product/detail',
    method: 'POST',
    data: { id: parseInt(id) }
  })
}

// 创建商品
export const createProduct = (data) => {
  return request({
    url: 'http://localhost:8001/api/product/create',
    method: 'POST',
    data
  })
}

// 更新商品
export const updateProduct = (data) => {
  return request({
    url: 'http://localhost:8001/api/product/update',
    method: 'POST',
    data
  })
}

// 删除商品
export const removeProduct = (id) => {
  return request({
    url: 'http://localhost:8001/api/product/remove',
    method: 'POST',
    data: { id }
  })
}

// 上传图片
export const uploadImage = (filePath) => {
  return new Promise((resolve, reject) => {
    wx.uploadFile({
      url: 'http://localhost:8001/api/product/upload',
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
    url: 'http://localhost:8001/api/product/images/add',
    method: 'POST',
    data: { productId, imageUrls }
  })
}

// 设置商品主图
export const setMainImage = (productId, imageUrl) => {
  return request({
    url: 'http://localhost:8001/api/product/images/setMain',
    method: 'POST',
    data: { productId, imageUrl }
  })
}

// 获取商品分类
export const getCategories = () => {
  return request({
    url: 'http://localhost:8001/api/product/v1/category/list',
    method: 'POST'
  })
}

// 获取轮播图
export const getBanners = () => {
  return request({
    url: 'http://localhost:8001/api/banner/v1/list',
    method: 'POST'
  })
}

// 搜索商品
export const searchProducts = (keyword) => {
  return request({
    url: 'http://localhost:8001/api/product/search',
    method: 'POST',
    data: { keyword }
  })
}

// 获取商品库存
export const getProductStock = (id) => {
  return request({
    url: 'http://localhost:8001/api/product/stock',
    method: 'POST',
    data: { id }
  })
}

// 获取商品图片
export function getProductImages(productId) {
  return request({
    url: '/api/product/images',
    method: 'POST',
    data: {
      productId
    }
  })
} 