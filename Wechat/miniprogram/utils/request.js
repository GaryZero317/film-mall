const app = getApp()

const request = (options) => {
  const { url, method = 'GET', data, noAuth = false } = options

  // 只有在需要认证且有token的情况下才添加token
  const token = noAuth ? '' : wx.getStorageSync('token')
  
  return new Promise((resolve, reject) => {
    console.log('发起请求:', {
      url,
      method,
      data,
      noAuth
    })

    const headers = {
      'Content-Type': 'application/json'
    }

    // 只有在需要认证时才添加token
    if (!noAuth && token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    wx.request({
      url: url.startsWith('http') ? url : `${app.globalData.baseUrl}${url}`,
      method,
      data,
      header: headers,
      success: (res) => {
        console.log('请求成功:', res)
        if (res.statusCode === 200) {
          // 处理不同的返回格式
          if (res.data && (res.data.code !== undefined)) {
            // 标准格式：{ code: 0, msg: '', data: {} }
            if (res.data.code === 0) {
              resolve(res.data)
            } else {
              const errorMsg = res.data.msg || '请求失败'
              console.error('业务错误:', errorMsg)
              reject(new Error(errorMsg))
            }
          } else if (typeof res.data === 'string' && res.data.includes('field')) {
            // 处理验证错误
            console.error('验证错误:', res.data)
            reject(new Error(res.data))
          } else {
            // 直接返回数据格式
            resolve({
              code: 0,
              msg: 'success',
              data: res.data
            })
          }
        } else if (res.statusCode === 401 && !noAuth) {
          // 只有需要认证的请求才处理401错误
          wx.removeStorageSync('token')
          wx.showToast({
            title: '请先登录',
            icon: 'none',
            duration: 2000,
            complete: () => {
              setTimeout(() => {
                wx.redirectTo({
                  url: '/pages/login/index'
                })
              }, 1000)
            }
          })
          reject(new Error('未登录或登录已过期'))
        } else {
          let errorMsg = '请求失败'
          if (res.data) {
            errorMsg = typeof res.data === 'string' ? res.data : (res.data.msg || '请求失败')
          }
          console.error('HTTP错误:', errorMsg)
          reject(new Error(errorMsg))
        }
      },
      fail: (error) => {
        console.error('请求错误:', error)
        reject(new Error('网络错误，请检查网络连接'))
      }
    })
  })
}

export default request 