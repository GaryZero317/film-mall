App({
  globalData: {
    userInfo: null,
    baseUrl: 'http://localhost:8001/api', // 修改端口为8001
    token: ''
  },
  onLaunch() {
    // 获取本地存储的token
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
    }
    
    // 初始化云开发环境
    if (wx.cloud) {
      wx.cloud.init({
        env: 'your-env-id', // 替换为您的云开发环境ID
        traceUser: true
      })
    }
  }
}) 