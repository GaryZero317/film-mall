import request from '../utils/request'

// 获取商品列表
export const getProductList = (params = {}) => {
  return request({
    url: '/product/list',
    method: 'POST',
    data: {
      page: params.page || 1,
      pageSize: params.pageSize || 10
    }
  })
}

// 获取商品详情
export const getProductDetail = (id) => {
  return request({
    url: '/product/detail',
    method: 'POST',
    data: { id }
  })
}

// 获取商品分类（这个接口似乎后端还没有提供，暂时返回固定数据）
export const getCategories = () => {
  return Promise.resolve({
    data: [
      { id: 1, name: '彩色负片' },
      { id: 2, name: '黑白负片' },
      { id: 3, name: '正片' },
      { id: 4, name: '一次性相机' },
      { id: 5, name: '中画幅胶卷' },
      { id: 6, name: '拍立得' }
    ]
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
    url: '/product/search',
    method: 'GET',
    data: { keyword }
  })
}

// 获取商品库存
export const getProductStock = (id) => {
  return request({
    url: `/product/stock/${id}`,
    method: 'GET'
  })
} 