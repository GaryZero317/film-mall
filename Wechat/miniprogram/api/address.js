const app = getApp()

// 获取地址列表
export function getAddressList() {
  return wx.request({
    url: `${app.globalData.baseUrl}/api/addresses`,
    method: 'GET',
    header: {
      'Authorization': `Bearer ${wx.getStorageSync('token')}`
    }
  })
}

// 获取地址详情
export function getAddressDetail(id) {
  return wx.request({
    url: `${app.globalData.baseUrl}/api/addresses/${id}`,
    method: 'GET',
    header: {
      'Authorization': `Bearer ${wx.getStorageSync('token')}`
    }
  })
}

// 新增地址
export function createAddress(data) {
  return wx.request({
    url: `${app.globalData.baseUrl}/api/addresses`,
    method: 'POST',
    data,
    header: {
      'Authorization': `Bearer ${wx.getStorageSync('token')}`
    }
  })
}

// 更新地址
export function updateAddress(id, data) {
  return wx.request({
    url: `${app.globalData.baseUrl}/api/addresses/${id}`,
    method: 'PUT',
    data,
    header: {
      'Authorization': `Bearer ${wx.getStorageSync('token')}`
    }
  })
}

// 删除地址
export function deleteAddress(id) {
  return wx.request({
    url: `${app.globalData.baseUrl}/api/addresses/${id}`,
    method: 'DELETE',
    header: {
      'Authorization': `Bearer ${wx.getStorageSync('token')}`
    }
  })
}

// 设置默认地址
export function setDefaultAddress(id) {
  return wx.request({
    url: `${app.globalData.baseUrl}/api/addresses/${id}/default`,
    method: 'POST',
    header: {
      'Authorization': `Bearer ${wx.getStorageSync('token')}`
    }
  })
} 