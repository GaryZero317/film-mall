// 检查是否已登录
export const checkLogin = () => {
  const token = wx.getStorageSync('token')
  const tokenExpire = wx.getStorageSync('tokenExpire')
  
  // 检查token是否存在
  if (!token) {
    navigateToLogin()
    return false
  }

  // 检查token是否过期
  if (tokenExpire && Date.now() / 1000 >= tokenExpire) {
    wx.removeStorageSync('token')
    wx.removeStorageSync('tokenExpire')
    navigateToLogin()
    return false
  }

  return true
}

// 页面登录守卫
export const loginGuard = (pageObj) => {
  // 保存原始的onLoad函数
  const originalOnLoad = pageObj.onLoad
  const originalOnShow = pageObj.onShow

  // 重写onLoad函数
  pageObj.onLoad = function(options) {
    // 检查登录状态
    const isLoggedIn = checkLogin()
    
    // 如果已登录，则执行原始的onLoad函数
    if (isLoggedIn && originalOnLoad) {
      originalOnLoad.call(this, options)
    }
  }

  // 重写onShow函数
  pageObj.onShow = function() {
    // 检查登录状态
    const isLoggedIn = checkLogin()
    
    // 如果已登录，则执行原始的onShow函数
    if (isLoggedIn && originalOnShow) {
      originalOnShow.call(this)
    }
  }

  return pageObj
}

// 跳转到登录页
const navigateToLogin = () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const url = currentPage ? currentPage.route : null
  
  wx.showToast({
    title: '请先登录',
    icon: 'none',
    duration: 1500,
    complete: () => {
      setTimeout(() => {
        // 将当前页面路径作为参数传递给登录页
        wx.navigateTo({
          url: `/pages/login/index?redirect=${url || ''}`
        })
      }, 1500)
    }
  })
} 