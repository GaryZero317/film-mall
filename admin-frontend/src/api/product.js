import { productService } from './request'

// 创建商品
export function createProduct(data) {
  return productService({
    url: '/api/product/create',
    method: 'post',
    data
  })
}

// 更新商品
export function updateProduct(data) {
  return productService({
    url: '/api/product/update',
    method: 'post',
    data
  })
}

// 删除商品
export function removeProduct(data) {
  return productService({
    url: '/api/product/remove',
    method: 'post',
    data
  })
}

// 获取商品详情
export function getProductDetail(data) {
  return productService({
    url: '/api/product/detail',
    method: 'post',
    data
  })
}

// 管理员获取商品列表
export function getAdminProductList(params) {
  return productService({
    url: '/api/admin/product/list',
    method: 'post',
    data: params
  })
}

// 设置商品主图
export const setMainImage = (data) => {
  console.log('设置主图请求数据:', data) // 添加日志
  return productService.post('/api/product/images/setMain', {
    productId: Number(data.productId), // 确保是数字类型
    imageUrl: data.imageUrl
  })
}

// 上传商品图片
export function uploadImage(data) {
  return productService({
    url: '/api/product/upload',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
} 