import request from '../utils/request'

// 获取购物车列表
export const getCartList = () => {
  return request({
    url: '/cart/list',
    method: 'GET'
  })
}

// 添加商品到购物车
export const addToCart = (data) => {
  return request({
    url: '/cart/add',
    method: 'POST',
    data
  })
}

// 更新购物车商品数量
export const updateCartItem = (data) => {
  return request({
    url: '/cart/update',
    method: 'PUT',
    data
  })
}

// 删除购物车商品
export const removeFromCart = (id) => {
  return request({
    url: `/cart/remove/${id}`,
    method: 'DELETE'
  })
}

// 清空购物车
export const clearCart = () => {
  return request({
    url: '/cart/clear',
    method: 'DELETE'
  })
}

// 购物车商品选中状态
export const updateCartItemStatus = (data) => {
  return request({
    url: '/cart/status',
    method: 'PUT',
    data
  })
} 