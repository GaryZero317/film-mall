import { wxLogin, login, register, getUserInfo } from '../../api/user'

Page({
  data: {
    isRegister: false, // 是否为注册状态
    name: '', // 姓名
    mobile: '', // 手机号
    gender: 1, // 性别，1-男，2-女
    password: '', // 密码
    confirmPassword: '', // 确认密码
    loading: false, // 登录/注册按钮loading
    wxLoading: false, // 微信登录按钮loading
    redirect: '' // 登录成功后的重定向页面
  },

  onLoad(options) {
    // 保存重定向页面
    if (options.redirect) {
      this.setData({ redirect: options.redirect })
    }

    // 页面加载时检查是否有token
    const token = wx.getStorageSync('token')
    if (token) {
      this.navigateBack()
    }
  },

  // 姓名输入
  onNameInput(e) {
    this.setData({ name: e.detail.value })
  },

  // 手机号输入
  onMobileInput(e) {
    this.setData({ mobile: e.detail.value })
  },

  // 性别选择
  onGenderChange(e) {
    this.setData({ gender: parseInt(e.detail.value) })
  },

  // 密码输入
  onPasswordInput(e) {
    this.setData({ password: e.detail.value })
  },

  // 确认密码输入
  onConfirmPasswordInput(e) {
    this.setData({ confirmPassword: e.detail.value })
  },

  // 切换到注册
  switchToRegister() {
    this.setData({
      isRegister: true,
      name: '',
      mobile: '',
      gender: 1,
      password: '',
      confirmPassword: '',
      loading: false
    })
  },

  // 切换到登录
  switchToLogin() {
    this.setData({
      isRegister: false,
      name: '',
      mobile: '',
      gender: 1,
      password: '',
      confirmPassword: '',
      loading: false
    })
  },

  // 登录
  async handleLogin() {
    const { mobile, password } = this.data
    
    if (!mobile || !password) {
      wx.showToast({
        title: '请填写完整信息',
        icon: 'none'
      })
      return
    }

    // 验证手机号格式
    if (!/^1[3-9]\d{9}$/.test(mobile)) {
      wx.showToast({
        title: '手机号格式不正确',
        icon: 'none'
      })
      return
    }

    try {
      this.setData({ loading: true })
      const res = await login({
        mobile,
        password
      })
      
      console.log('[登录] 登录响应:', res)
      
      // 处理登录响应
      if (res && res.data && res.data.accessToken) {
        // 保存token和过期时间
        wx.setStorageSync('token', res.data.accessToken)
        wx.setStorageSync('tokenExpire', res.data.accessExpire)
        
        // 更新全局数据
        getApp().globalData.token = res.data.accessToken
        
        // 获取用户信息并保存
        try {
          const userInfoRes = await getUserInfo()
          if (userInfoRes && userInfoRes.code === 200 && userInfoRes.data) {
            wx.setStorageSync('userInfo', userInfoRes.data)
            console.log('[登录] 已保存用户信息:', userInfoRes.data)
          } else {
            console.error('[登录] 获取用户信息失败:', userInfoRes)
          }
        } catch (error) {
          console.error('[登录] 获取用户信息出错:', error)
        }
        
        wx.showToast({
          title: '登录成功',
          icon: 'success',
          duration: 1500
        })
        
        // 延迟跳转，让用户看到成功提示
        setTimeout(() => {
          // 如果有重定向页面，则跳转到重定向页面
          if (this.data.redirect) {
            wx.redirectTo({
              url: decodeURIComponent(this.data.redirect)
            })
          } else {
            // 否则跳转到首页
            wx.reLaunch({
              url: '/pages/index/index'
            })
          }
        }, 1500)
      } else {
        let errorMsg = '登录失败'
        if (res.code !== 200) {
          switch(res.code) {
            case 30001:
              errorMsg = '用户不存在'
              break
            case 30002:
              errorMsg = '密码错误'
              break
            case 30005:
              errorMsg = '登录已过期'
              break
            case 30006:
              errorMsg = '无效的登录凭证'
              break
            default:
              errorMsg = res.msg || '登录失败，请重试'
          }
        }
        wx.showToast({
          title: errorMsg,
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[登录] 登录失败:', error)
      wx.showToast({
        title: error.message || '登录失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 注册
  async handleRegister() {
    const { name, mobile, gender, password, confirmPassword } = this.data
    
    if (!name || !mobile || !password || !confirmPassword) {
      wx.showToast({
        title: '请填写完整信息',
        icon: 'none'
      })
      return
    }

    // 验证手机号格式
    if (!/^1[3-9]\d{9}$/.test(mobile)) {
      wx.showToast({
        title: '手机号格式不正确',
        icon: 'none'
      })
      return
    }

    if (password !== confirmPassword) {
      wx.showToast({
        title: '两次密码不一致',
        icon: 'none'
      })
      return
    }

    try {
      this.setData({ loading: true })
      const res = await register({
        name,
        mobile,
        gender,
        password
      })
      
      console.log('[注册] 注册响应:', res)
      
      if (res.code === 200 && res.data) {
        wx.showToast({
          title: '注册成功',
          icon: 'success',
          duration: 1500
        })
        
        // 注册成功后切换到登录页，并自动填充手机号和密码
        setTimeout(() => {
          this.setData({
            isRegister: false,
            mobile: mobile,
            password: password,
            loading: false
          })
          // 自动触发登录
          this.handleLogin()
        }, 1500)
      } else {
        let errorMsg = '注册失败'
        switch(res.code) {
          case 30003:
            errorMsg = '手机号已注册'
            break
          case 10001:
            errorMsg = '请填写正确的信息'
            break
          default:
            errorMsg = res.msg || '注册失败，请重试'
        }
        wx.showToast({
          title: errorMsg,
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[注册] 注册失败:', error)
      wx.showToast({
        title: error.message || '注册失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 微信登录
  async handleWxLogin() {
    if (this.data.wxLoading) return
    
    try {
      this.setData({ wxLoading: true })
      
      // 获取微信登录code
      const { code } = await wx.login()
      
      // 调用后端登录接口
      const res = await wxLogin(code)
      
      if (res && res.data && res.data.accessToken) {
        // 保存token和过期时间
        wx.setStorageSync('token', res.data.accessToken)
        wx.setStorageSync('tokenExpire', res.data.accessExpire)
        
        // 更新全局数据
        getApp().globalData.token = res.data.accessToken
        
        // 获取用户信息并保存
        try {
          const userInfoRes = await getUserInfo()
          if (userInfoRes && userInfoRes.code === 200 && userInfoRes.data) {
            wx.setStorageSync('userInfo', userInfoRes.data)
            console.log('[微信登录] 已保存用户信息:', userInfoRes.data)
          } else {
            console.error('[微信登录] 获取用户信息失败:', userInfoRes)
          }
        } catch (error) {
          console.error('[微信登录] 获取用户信息出错:', error)
        }
        
        wx.showToast({
          title: '登录成功',
          icon: 'success',
          duration: 1500
        })
        
        setTimeout(() => {
          // 如果有重定向页面，则跳转到重定向页面
          if (this.data.redirect) {
            wx.redirectTo({
              url: decodeURIComponent(this.data.redirect)
            })
          } else {
            // 否则跳转到首页
            wx.reLaunch({
              url: '/pages/index/index'
            })
          }
        }, 1500)
      } else {
        throw new Error(res.msg || '登录失败')
      }
    } catch (error) {
      console.error('[微信登录] 登录失败:', error)
      wx.showToast({
        title: error.message || '登录失败',
        icon: 'none'
      })
    } finally {
      this.setData({ wxLoading: false })
    }
  },

  // 返回上一页或重定向页面
  navigateBack() {
    const { redirect } = this.data
    if (redirect) {
      wx.redirectTo({
        url: decodeURIComponent(redirect)
      })
    } else {
      // 如果没有重定向页面，直接跳转到首页
      wx.reLaunch({
        url: '/pages/index/index'
      })
    }
  }
}) 