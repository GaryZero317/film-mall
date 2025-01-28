import request from '../utils/request'

// 获取购物车列表
export const getCartList = () => {
  return request({
    url: '/api/cart/list',
    method: 'GET'
  })
}

// 添加商品到购物车
export const addToCart = (data) => {
  return request({
    url: '/api/cart/add',
    method: 'POST',
    data
  })
}

// 更新购物车商品数量
export const updateCartItem = (data) => {
  return request({
    url: '/api/cart/quantity',
    method: 'PUT',
    data
  })
}

// 删除购物车商品
export const removeFromCart = (id) => {
  return request({
    url: `/api/cart/${id}`,
    method: 'DELETE'
  })
}

// 清空购物车
export const clearCart = () => {
  return request({
    url: '/api/cart/clear',
    method: 'DELETE'
  })
}

// 购物车商品选中状态
export const updateCartItemStatus = (data) => {
  return request({
    url: '/api/cart/selected',
    method: 'PUT',
    data
  })
}

// 清除购物车商品
export function clearCartItems(productIds) {
  return request({
    url: '/api/cart/clear',
    method: 'DELETE',
    data: {
      product_ids: productIds
    }
  })
} 