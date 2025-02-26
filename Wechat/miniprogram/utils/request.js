const app = getApp()

const request = (options) => {
  const { url, method = 'GET', data, params, noAuth = false } = options

  console.log('app.globalData:', app.globalData)
  console.log('baseUrl配置:', app.globalData?.baseUrl)

  // 获取服务类型和实际路径
  let baseUrl = 'http://localhost:8001' // 默认baseUrl
  
  try {
    // 如果url已经包含完整域名，不使用baseUrl
    if (url.indexOf('localhost:') !== -1) {
      baseUrl = ''
    }
    // 根据API路径选择不同的服务地址
    else if (url.startsWith('/api/cart')) {
      baseUrl = 'http://localhost:8004'
    } else if (url.startsWith('/api/address')) {
      baseUrl = 'http://localhost:8005'
    } else if (url.startsWith('/api/order')) {
      baseUrl = 'http://localhost:8002'
    } else if (url.startsWith('/api/pay')) {
      baseUrl = 'http://localhost:8003'
    } else if (url.startsWith('/api/film')) {
      baseUrl = 'http://localhost:8007'
    }

    console.log('选择的baseUrl:', baseUrl, '请求路径:', url)
  } catch (error) {
    console.error('处理baseUrl时出错:', error)
  }

  // 只有在需要认证且有token的情况下才添加token
  const token = noAuth ? '' : wx.getStorageSync('token')
  
  return new Promise((resolve, reject) => {
    // 处理GET请求的查询参数
    let fullUrl = baseUrl ? `${baseUrl}${url}` : url
    if (method.toUpperCase() === 'GET' && params) {
      const queryString = Object.keys(params)
        .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
        .join('&')
      fullUrl = `${fullUrl}${fullUrl.includes('?') ? '&' : '?'}${queryString}`
    }
    
    console.log('完整请求URL:', fullUrl)

    console.log('发起请求:', {
      baseUrl,
      url: fullUrl,
      method,
      data: method.toUpperCase() === 'GET' ? null : data,
      noAuth,
      headers: options.header
    })

    const headers = {
      'Content-Type': 'application/json'
    }

    // 只有在需要认证时才添加token
    if (!noAuth && token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    wx.request({
      url: fullUrl,
      method,
      data,
      header: headers,
      success: (res) => {
        console.log('[请求工具] 收到响应:', {
          状态码: res.statusCode,
          响应数据: res.data,
          响应类型: typeof res.data
        })
        
        if (res.statusCode === 200) {
          // 处理不同的返回格式
          if (res.data && typeof res.data === 'object') {
            console.log('[请求工具] 处理对象响应:', res.data)

            // 如果是地址详情接口，直接返回数据
            if (url.startsWith('/api/address/') && method === 'GET' && !url.endsWith('/list')) {
              console.log('[请求工具] 地址详情响应:', res.data)
              resolve(res.data)
              return
            }

            // 如果响应中已经包含了code字段，直接返回
            if ('code' in res.data) {
              // 如果是成功响应（code为0或200），直接返回
              if (res.data.code === 200 || res.data.code === 0) {
                console.log('[请求工具] 成功响应:', res.data)
                resolve(res.data)
                return
              }
              
              // 处理错误响应
              let errorMsg = res.data.msg || '请求失败'
              switch(res.data.code) {
                // 基础错误码
                case 10001:
                  errorMsg = '参数错误'
                  break
                case 10002:
                  errorMsg = '服务器内部错误'
                  break
                
                // 用户相关错误码
                case 30001:
                  errorMsg = '用户不存在'
                  break
                case 30002:
                  errorMsg = '密码错误'
                  break
                case 30003:
                  errorMsg = '手机号已注册'
                  break
                case 30004:
                  errorMsg = '验证码错误'
                  break
                case 30005:
                  errorMsg = '登录已过期'
                  break
                case 30006:
                  errorMsg = '无效的登录凭证'
                  break
                
                // 订单相关错误码
                case 20001:
                  errorMsg = '订单不存在'
                  break
                case 20002:
                  errorMsg = '订单创建失败'
                  break
                case 20003:
                  errorMsg = '订单状态无效'
                  break
                case 20004:
                  errorMsg = '订单更新失败'
                  break
                case 20005:
                  errorMsg = '订单金额无效'
                  break
              }
              console.error('[请求工具] 业务错误:', errorMsg)
              reject(new Error(errorMsg))
              return
            }

            // 包装其他响应
            console.log('[请求工具] 包装对象响应:', {
              code: 200,
              msg: 'success',
              data: res.data
            })
            resolve({
              code: 200,
              msg: 'success',
              data: res.data
            })
          } else {
            // 处理字符串响应
            if (typeof res.data === 'string') {
              // 处理验证错误
              if (res.data.indexOf('field') !== -1) {
                console.error('[请求工具] 验证错误:', res.data)
                reject(new Error(res.data))
                return
              }

              // 处理成功字符串
              if (res.data === 'success' || res.data === '成功') {
                const response = {
                  code: 200,
                  msg: '成功',
                  data: null
                }
                console.log('[请求工具] 成功字符串响应:', response)
                resolve(response)
                return
              }

              // 其他字符串响应
              const response = {
                code: 200,
                msg: res.data,
                data: null
              }
              console.log('[请求工具] 其他字符串响应:', response)
              resolve(response)
              return
            }

            // 其他类型响应
            const response = {
              code: 200,
              msg: 'success',
              data: res.data
            }
            console.log('[请求工具] 其他类型响应:', response)
            resolve(response)
            return
          }
        } else if (res.statusCode === 401 && !noAuth) {
          // 处理401错误，未登录或登录过期
          wx.removeStorageSync('token')
          wx.showToast({
            title: '请重新登录',
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
        } else if (res.statusCode === 400) {
          // 处理400错误，通常是请求参数错误或资源不存在
          let errorMsg = '请求参数错误'
          if (res.data) {
            if (typeof res.data === 'string' && res.data.indexOf('no rows') !== -1) {
              // 对于"no rows in result set"错误，返回特定的响应格式
              resolve({
                code: 404,
                msg: '商品不存在',
                data: null,
                notFound: true
              })
              return
            }
            errorMsg = typeof res.data === 'string' ? res.data : (res.data.msg || '请求参数错误')
          }
          console.error('请求错误:', errorMsg, res)
          reject(new Error(errorMsg))
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