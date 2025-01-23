// pages/user/index.js
import { loginGuard } from '../../utils/auth'
import { getUserInfo } from '../../api/user'

const app = getApp()

Page(loginGuard({
  data: {
    isLogin: false,
    userInfo: null,
    orderCount: {
      unpaid: 0,
      undelivered: 0,
      delivered: 0
    },
    loading: false
  },

  onLoad() {
    this.loadUserInfo()
  },

  onShow() {
    // 每次显示页面时检查登录状态
    this.checkLoginStatus()
    if (this.data.isLogin) {
      this.loadUserInfo()
      this.loadOrderCount()
    }
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    this.setData({ isLogin: !!token })
  },

  // 加载用户信息
  async loadUserInfo() {
    try {
      this.setData({ loading: true })
      const res = await getUserInfo()
      if (res.code === 0) {
        this.setData({ userInfo: res.data })
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
    } finally {
      this.setData({ loading: false })
    }
  },

  // 加载订单数量
  async loadOrderCount() {
    try {
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/orders/count`,
        method: 'GET',
        header: {
          'Authorization': `Bearer ${wx.getStorageSync('token')}`
        }
      })

      if (res.statusCode === 200) {
        this.setData({ orderCount: res.data })
      }
    } catch (error) {
      console.error('加载订单数量失败:', error)
    }
  },

  // 登录
  async onLogin() {
    try {
      // 获取用户信息
      const userInfoRes = await wx.getUserProfile({
        desc: '用于完善会员资料'
      })

      // 获取登录code
      const loginRes = await wx.login()
      
      // 发送登录请求
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/login`,
        method: 'POST',
        data: {
          code: loginRes.code,
          user_info: userInfoRes.userInfo
        }
      })

      if (res.statusCode === 200) {
        // 保存token
        wx.setStorageSync('token', res.data.token)
        this.setData({ 
          isLogin: true,
          userInfo: res.data.user
        })
        // 加载订单数量
        this.loadOrderCount()
      } else {
        wx.showToast({
          title: '登录失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('登录失败:', error)
      wx.showToast({
        title: '登录失败',
        icon: 'none'
      })
    }
  },

  // 退出登录
  async onLogout() {
    const res = await wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      confirmText: '退出',
      cancelText: '取消'
    })

    if (res.confirm) {
      // 清除token和用户信息
      wx.removeStorageSync('token')
      this.setData({
        isLogin: false,
        userInfo: null,
        orderCount: {
          unpaid: 0,
          undelivered: 0,
          delivered: 0
        }
      })
    }
  },

  // 查看全部订单
  onViewAllOrders() {
    if (!this.data.isLogin) {
      this.showLoginTip()
      return
    }
    wx.navigateTo({
      url: '/pages/order/list/index'
    })
  },

  // 查看指定状态的订单
  onViewOrders(e) {
    if (!this.data.isLogin) {
      this.showLoginTip()
      return
    }
    const type = e.currentTarget.dataset.type
    wx.navigateTo({
      url: `/pages/order/list/index?tab=${type}`
    })
  },

  // 查看收货地址
  onViewAddress() {
    if (!this.data.isLogin) {
      this.showLoginTip()
      return
    }
    wx.navigateTo({
      url: '/pages/address/list/index'
    })
  },

  // 查看收藏
  onViewFavorites() {
    if (!this.data.isLogin) {
      this.showLoginTip()
      return
    }
    wx.navigateTo({
      url: '/pages/user/favorites/index'
    })
  },

  // 查看优惠券
  onViewCoupons() {
    if (!this.data.isLogin) {
      this.showLoginTip()
      return
    }
    wx.navigateTo({
      url: '/pages/user/coupons/index'
    })
  },

  // 查看设置
  onViewSettings() {
    wx.navigateTo({
      url: '/pages/user/settings/index'
    })
  },

  // 显示登录提示
  showLoginTip() {
    wx.showToast({
      title: '请先登录',
      icon: 'none'
    })
  }
}))