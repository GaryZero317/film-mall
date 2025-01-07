const app = getApp()

const request = (options) => {
  return new Promise((resolve, reject) => {
    const token = app.globalData.token
    wx.request({
      url: `${app.globalData.baseUrl}${options.url}`,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else if (res.statusCode === 401) {
          // token过期，清除本地存储并跳转到登录页
          wx.removeStorageSync('token')
          wx.navigateTo({
            url: '/pages/login/index'
          })
          reject(res)
        } else {
          wx.showToast({
            title: res.data.message || '请求失败',
            icon: 'none'
          })
          reject(res)
        }
      },
      fail: (err) => {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
        reject(err)
      }
    })
  })
}

export default request 