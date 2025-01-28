// pages/user/index.js
import { getUserInfo } from '../../api/user'

const app = getApp()

Page({
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
    this.checkLoginStatus()
  },

  onShow() {
    this.checkLoginStatus()
    if (this.data.isLogin) {
      this.loadUserInfo()
    }
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    console.log('当前token:', token)
    this.setData({ 
      isLogin: !!token,
      userInfo: token ? wx.getStorageSync('userInfo') : null
    })
  },

  // 加载用户信息
  async loadUserInfo() {
    if (!this.data.isLogin) {
      return
    }

    try {
      this.setData({ loading: true })
      const res = await getUserInfo()
      console.log('获取用户信息结果:', res)
      if (res.code === 200 && res.data) {
        this.setData({ 
          userInfo: res.data,
          isLogin: true
        })
        // 保存用户信息到本地存储
        wx.setStorageSync('userInfo', res.data)
      } else {
        // 如果获取用户信息失败，可能是token失效
        this.setData({ 
          isLogin: false,
          userInfo: null 
        })
        wx.removeStorageSync('token')
        wx.removeStorageSync('userInfo')
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 获取用户信息失败时，清除登录状态
      this.setData({ 
        isLogin: false,
        userInfo: null 
      })
      wx.removeStorageSync('token')
      wx.removeStorageSync('userInfo')
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
  onLogin() {
    wx.navigateTo({
      url: '/pages/login/index'
    })
  },

  // 退出登录
  onLogout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          // 清除token和用户信息
          wx.removeStorageSync('token')
          wx.removeStorageSync('userInfo')
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
      }
    })
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
})