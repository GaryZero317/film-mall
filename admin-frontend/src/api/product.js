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
    data: {
      page: params.page || 1,
      pageSize: params.pageSize || 10
    }
  })
} 