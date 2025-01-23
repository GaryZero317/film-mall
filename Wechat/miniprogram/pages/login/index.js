import { wxLogin, login, register } from '../../api/user'

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

  // 账号密码登录
  async handleLogin() {
    const { mobile, password } = this.data
    if (!mobile || !password) {
      wx.showToast({
        title: '请输入手机号和密码',
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
      
      console.log('登录响应数据:', res) // 打印响应数据，方便调试
      
      // 检查返回的数据结构
      if (res && res.accessToken) {
        // 保存token
        wx.setStorageSync('token', res.accessToken)
        if (res.accessExpire) {
          wx.setStorageSync('tokenExpire', res.accessExpire)
        }
        getApp().globalData.token = res.accessToken
        
        wx.showToast({
          title: '登录成功',
          icon: 'success',
          duration: 1500
        })
        
        setTimeout(() => {
          this.navigateBack()
        }, 1500)
      } else {
        throw new Error('登录失败，请检查账号密码')
      }
    } catch (error) {
      console.error('登录失败，详细错误:', error)
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
      
      console.log('注册响应数据:', res)
      
      // 后端直接返回用户数据，包含 id、name、gender、mobile
      if (res && res.id) {
        wx.showToast({
          title: '注册成功',
          icon: 'success',
          duration: 1500
        })
        
        // 注册成功后切换到登录页，并自动填充手机号
        setTimeout(() => {
          this.setData({
            isRegister: false,
            mobile: res.mobile, // 使用返回的手机号
            password: '',
            loading: false
          })
        }, 1500)
      } else {
        throw new Error('注册失败，请重试')
      }
    } catch (error) {
      console.error('注册失败，详细错误:', error)
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
      
      if (res.code === 0 && res.data) {
        const { accessToken, accessExpire } = res.data
        // 保存token和过期时间
        wx.setStorageSync('token', accessToken)
        wx.setStorageSync('tokenExpire', accessExpire)
        getApp().globalData.token = accessToken
        
        wx.showToast({
          title: '登录成功',
          icon: 'success',
          duration: 1500
        })
        
        setTimeout(() => {
          this.navigateBack()
        }, 1500)
      } else {
        throw new Error(res.msg || '登录失败')
      }
    } catch (error) {
      console.error('登录失败:', error)
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
        url: `/${redirect}`
      })
    } else {
      const pages = getCurrentPages()
      if (pages.length > 1) {
        wx.navigateBack()
      } else {
        wx.switchTab({
          url: '/pages/index/index'
        })
      }
    }
  }
}) 