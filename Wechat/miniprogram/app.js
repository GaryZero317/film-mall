App({
  globalData: {
    userInfo: null,
    baseUrl: 'http://localhost:8888/api', // 后端接口基础地址
    token: ''
  },
  onLaunch() {
    // 获取本地存储的token
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
    }
  }
}) 