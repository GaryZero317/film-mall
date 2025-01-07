import request from '../utils/request'

// 获取商品列表
export const getProductList = (params) => {
  return request({
    url: '/film/list',
    method: 'GET',
    data: params
  })
}

// 获取商品详情
export const getProductDetail = (id) => {
  return request({
    url: `/film/detail/${id}`,
    method: 'GET'
  })
}

// 获取商品分类
export const getCategories = () => {
  return request({
    url: '/film/categories',
    method: 'GET'
  })
}

// 获取轮播图
export const getBanners = () => {
  return request({
    url: '/banner/list',
    method: 'GET'
  })
}

// 搜索商品
export const searchProducts = (keyword) => {
  return request({
    url: '/film/search',
    method: 'GET',
    data: { keyword }
  })
}

// 获取商品库存
export const getProductStock = (id) => {
  return request({
    url: `/film/stock/${id}`,
    method: 'GET'
  })
} 