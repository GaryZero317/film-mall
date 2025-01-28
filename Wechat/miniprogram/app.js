App({
  globalData: {
    userInfo: null,
    baseUrl: {
      cart: 'http://localhost:8004',    // 购物车服务
      address: 'http://localhost:8005',   // 地址服务
      order: 'http://localhost:8002',
      pay: 'http://localhost:8003',
      product: 'http://localhost:8001'
    },
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
        env: 'your-env-id', 
        traceUser: true
      })
    }
  }
}) 