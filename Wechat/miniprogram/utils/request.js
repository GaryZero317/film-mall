const app = getApp()

const request = (options) => {
  const { url, method = 'GET', data } = options

  const token = wx.getStorageSync('token')
  
  return new Promise((resolve, reject) => {
    console.log('发起请求:', {
      url,
      method,
      data
    })

    wx.request({
      url: url.startsWith('http') ? url : `${app.globalData.baseUrl}${url}`,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      },
      success: (res) => {
        console.log('请求成功:', res)
        if (res.statusCode === 200) {
          // 处理不同的返回格式
          if (res.data && (res.data.code !== undefined)) {
            // 标准格式：{ code: 0, msg: '', data: {} }
            if (res.data.code === 0) {
              resolve(res.data)  // 返回完整的响应数据，包括 code, msg, data
            } else {
              reject(res.data)
            }
          } else {
            // 直接返回数据格式
            resolve({
              code: 0,
              msg: 'success',
              data: res.data
            })
          }
        } else if (res.statusCode === 401) {
          // token过期或无效，需要重新登录
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
          console.error('请求失败:', res)
          wx.showToast({
            title: (res.data && res.data.msg) || '请求失败',
            icon: 'none'
          })
          reject(res)
        }
      },
      fail: (error) => {
        console.error('请求错误:', error)
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
        reject(error)
      }
    })
  })
}

export default request 