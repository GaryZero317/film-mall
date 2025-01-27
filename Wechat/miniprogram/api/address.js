const app = getApp()
const request = require('../utils/request').default

module.exports = {
  // 获取地址列表
  getAddressList() {
    return request({
      url: '/api/address/list',
      method: 'GET'
    })
  },

  // 新增地址
  createAddress(data) {
    return request({
      url: '/api/address/add',
      method: 'POST',
      data
    })
  },

  // 更新地址
  updateAddress(id, data) {
    return request({
      url: '/api/address/update',
      method: 'POST',
      data: {
        id,
        ...data
      }
    })
  },

  // 删除地址
  deleteAddress(id) {
    return request({
      url: '/api/address/delete',
      method: 'POST',
      data: {
        id
      }
    })
  },

  // 设置默认地址
  setDefaultAddress(id) {
    return request({
      url: '/api/address/setDefault',
      method: 'POST',
      data: {
        id
      }
    })
  }
} 