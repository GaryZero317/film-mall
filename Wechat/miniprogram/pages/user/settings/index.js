// pages/user/settings/index.js
const app = getApp()
import request from '../../../utils/request'

Page({

  /**
   * 页面的初始数据
   */
  data: {
    phoneNumber: '',
    orderNotification: true,
    promotionNotification: true,
    cacheSize: '0KB'
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    this.initData()
  },

  initData: function () {
    // 获取用户信息
    request({
      url: 'http://localhost:8000/api/user/userinfo',
      method: 'POST'
    }).then(res => {
      if (res.code === 200 && res.data) {
        this.setData({
          phoneNumber: res.data.mobile || ''
        })
      }
    }).catch(err => {
      console.error('获取用户信息失败:', err)
    })

    // 获取缓存大小
    wx.getStorageInfo({
      success: res => {
        const sizeInKB = (res.currentSize).toFixed(2)
        this.setData({
          cacheSize: sizeInKB + 'KB'
        })
      }
    })
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

  },

  // 修改密码
  onChangePassword: function () {
    wx.navigateTo({
      url: '/pages/user/settings/change-password/index'
    })
  },

  // 绑定手机
  onBindPhone: function () {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 订单通知开关
  onOrderNotificationChange: function (e) {
    this.setData({
      orderNotification: e.detail.value
    })
  },

  // 优惠活动通知开关
  onPromotionNotificationChange: function (e) {
    this.setData({
      promotionNotification: e.detail.value
    })
  },

  // 隐私政策
  onPrivacyPolicy: function () {
    wx.navigateTo({
      url: '/pages/webview/index?type=privacy'
    })
  },

  // 用户协议
  onUserAgreement: function () {
    wx.navigateTo({
      url: '/pages/webview/index?type=agreement'
    })
  },

  // 清除缓存
  clearCache: function () {
    wx.showModal({
      title: '提示',
      content: '确定要清除缓存吗？',
      success: res => {
        if (res.confirm) {
          wx.clearStorage({
            success: () => {
              this.setData({
                cacheSize: '0KB'
              })
              wx.showToast({
                title: '清除成功',
                icon: 'success'
              })
            }
          })
        }
      }
    })
  },

  // 关于我们
  onAboutUs: function () {
    wx.navigateTo({
      url: '/pages/webview/index?type=about'
    })
  },

  // 退出登录
  onLogout: function () {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: res => {
        if (res.confirm) {
          // 清除本地存储的用户信息
          wx.removeStorageSync('token')
          wx.removeStorageSync('userInfo')
          
          // 跳转到登录页
          wx.reLaunch({
            url: '/pages/login/index'
          })
        }
      }
    })
  }
})