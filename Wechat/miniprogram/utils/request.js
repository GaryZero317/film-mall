const app = getApp()

const request = (options) => {
  return new Promise((resolve, reject) => {
    const token = wx.getStorageSync('token')
    wx.request({
      url: options.url.startsWith('http') ? options.url : `${app.globalData.baseUrl}${options.url}`,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        if (res.statusCode === 200) {
          // 处理不同的返回格式
          if (res.data && (res.data.code !== undefined)) {
            // 标准格式：{ code: 0, msg: '', data: {} }
            if (res.data.code === 0) {
              resolve(res.data.data)
            } else {
              reject(new Error(res.data.msg || '请求失败'))
            }
          } else {
            // 直接返回数据格式
            resolve(res.data)
          }
        } else if (res.statusCode === 401) {
          // token过期或未登录，跳转到登录页
          wx.removeStorageSync('token')
          wx.showToast({
            title: '请先登录',
            icon: 'none',
            duration: 2000,
            complete: () => {
              setTimeout(() => {
                wx.navigateTo({
                  url: '/pages/login/index'
                })
              }, 1000)
            }
          })
          reject(res)
        } else {
          wx.showToast({
            title: (res.data && res.data.msg) || '请求失败',
            icon: 'none'
          })
          reject(res)
        }
      },
      fail: (err) => {
        console.error('请求失败:', err)
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