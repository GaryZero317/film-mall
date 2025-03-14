// pages/order/list/index.js
import { getOrderList, cancelOrder, confirmOrder, updateOrderStatus } from '../../../api/order'
import { getProductDetail, getProductImages } from '../../../api/product'

const app = getApp()

// 定义订单支付超时时间（毫秒），与后端保持一致，15分钟
const ORDER_PAYMENT_TIMEOUT = 15 * 60 * 1000;

Page({
  data: {
    orders: [],
    loading: false,
    currentTab: 0,
    tabs: [
      { id: 0, name: '全部', status: null },  // 全部订单不需要状态值
      { id: 1, name: '待付款', status: 0 },   // 0 表示待付款
      { id: 2, name: '待发货', status: 1 },   // 1 表示待发货
      { id: 3, name: '待收货', status: 2 },   // 2 表示待收货
      { id: 4, name: '已完成', status: 3 }    // 3 表示已完成
    ],
    isFirstLoad: true,  // 添加标记，用于区分首次加载和后续显示
    page: 1,           // 当前页码
    pageSize: 10,      // 每页数量
    hasMore: true,     // 是否还有更多数据
    countdownTimers: {} // 存储订单倒计时定时器
  },

  // 格式化日期字符串，确保iOS兼容性
  formatDateString(dateStr) {
    if (!dateStr) return ''
    // 将 "yyyy-MM-dd HH:mm:ss" 转换为 "yyyy/MM/dd HH:mm:ss"
    return dateStr.replace(/-/g, '/')
  },

  onLoad() {
    console.log('[订单列表] 页面加载')
    this.loadOrders()
  },

  onShow() {
    console.log('[订单列表] 页面显示', {
      isFirstLoad: this.data.isFirstLoad,
      currentTab: this.data.currentTab,
      page: this.data.page,
      hasMore: this.data.hasMore
    })
    // 如果不是首次加载，则刷新订单列表
    if (!this.data.isFirstLoad) {
      this.loadOrders()
    }
    // 重置首次加载标记
    this.setData({ isFirstLoad: false })
  },

  onHide() {
    // 页面隐藏时清除所有倒计时
    this.clearAllCountdowns()
  },

  onUnload() {
    // 页面卸载时清除所有倒计时
    this.clearAllCountdowns()
  },

  // 清除所有倒计时
  clearAllCountdowns() {
    const { countdownTimers } = this.data
    Object.keys(countdownTimers).forEach(orderId => {
      clearInterval(countdownTimers[orderId])
    })
    this.setData({
      countdownTimers: {}
    })
  },

  // 切换标签
  onTabChange(e) {
    const index = parseInt(e.currentTarget.dataset.index)
    if (isNaN(index) || index < 0 || index >= this.data.tabs.length) {
      console.error('[订单列表] 无效的标签索引:', index)
      return
    }
    
    console.log('[订单列表] 切换标签:', { 
      原标签: this.data.currentTab,
      新标签: index,
      标签名称: this.data.tabs[index].name,
      当前页码: this.data.page,
      是否有更多: this.data.hasMore
    })

    // 切换标签前清除所有倒计时
    this.clearAllCountdowns()

    // 切换标签时重置分页
    this.setData({
      currentTab: index
    }, () => {
      this.resetPage()
      this.loadOrders()
    })
  },

  // 加载订单列表
  async loadOrders() {
    console.log('[订单列表] 开始加载订单')
    if (this.data.loading) {
      return
    }

    this.setData({ loading: true })

    // 加载前清除已有的倒计时
    this.clearAllCountdowns()

    try {
      // 从全局获取用户信息
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo || !userInfo.id) {
        throw new Error('请先登录')
      }

      // 获取当前标签对应的状态
      const currentStatus = this.data.tabs[this.data.currentTab].status

      // 构造请求参数
      const params = {
        uid: userInfo.id,
        status: currentStatus === null ? -1 : currentStatus,  // 全部订单时传递-1
        page: this.data.page,
        pageSize: this.data.pageSize
      }

      console.log('[订单列表] 请求参数:', params)
      const res = await getOrderList(params)
      console.log('[订单列表] 返回结果:', res)

      if (res.code === 0 && res.data) {
        // 记录原始数据结构
        console.log('[订单列表] 原始数据结构:', {
          dataKeys: Object.keys(res.data),
          isList: Array.isArray(res.data.list),
          listLength: res.data.list ? res.data.list.length : 0
        })
        
        // 使用正确的字段 list 而不是 data
        const orderItems = res.data.list || []
        console.log('[订单列表] 订单数据项:', orderItems)
        
        // 处理订单列表数据
        const orders = orderItems.map(order => {
          // 格式化价格，后端返回的金额单位是分，除以100转换为元
          const amount = ((order.total_price || 0) / 100).toFixed(2)
          // 获取状态文本
          const status_text = this.getStatusText(order.status)
          // 格式化日期
          const create_time = order.create_time
            ? this.formatDateString(order.create_time)
            : ''
            
          // 处理订单中的商品项价格
          const items = Array.isArray(order.items) ? order.items.map(item => {
            // 记录原始商品项数据
            console.log('[订单列表] 商品项原始数据:', item)
              
            // 确保价格从分转换为元，支持多种可能的字段名
            let originalPrice = item.price
            // 根据价格的值判断是否需要转换
            // 如果价格大于1000，可能是分单位，需要转换
            if (originalPrice && parseInt(originalPrice) > 1000) {
              originalPrice = (parseInt(originalPrice) / 100).toFixed(2)
            } else if (originalPrice) {
              // 否则可能已经是元单位，保留两位小数
              originalPrice = parseFloat(originalPrice).toFixed(2)
            } else {
              originalPrice = '0.00'
            }
            
            const itemPrice = originalPrice
            
            const itemAmount = item.amount !== undefined ? 
               (parseInt(item.amount) / 100).toFixed(2) : 
               '0.00'
            
            // 确保商品名称字段统一 - 尝试多种可能的字段名
            const productName = 
              item.product_name || 
              item.name || 
              item.productName || 
              '未知商品'
            
            // 确保商品图片字段统一
            const productImage = 
              item.product_image || 
              item.productImage || 
              item.cover_image || 
              '/assets/images/default.png'
            
            return {
              ...item,
              price: itemPrice,
              amount: itemAmount,
              product_name: productName,
              product_image: productImage
            }
          }) : []

          return {
            ...order,
            amount,
            status_text,
            create_time,
            items: items
          }
        })

        // 详细日志输出处理后的第一个订单信息
        if (orders.length > 0) {
          console.log('[订单列表] 第一个订单处理后数据:', {
            订单ID: orders[0].id,
            订单号: orders[0].oid,
            订单状态: orders[0].status,
            状态文本: orders[0].status_text,
            总金额: orders[0].amount,
            创建时间: orders[0].create_time,
            商品数量: orders[0].items ? orders[0].items.length : 0
          })
          
          // 输出商品项信息
          if (orders[0].items && orders[0].items.length > 0) {
            console.log('[订单列表] 第一个订单的商品项:', orders[0].items.map(item => ({
              商品ID: item.id || item.pid,
              商品名称: item.product_name,
              原始价格: item.price,
              商品数量: item.quantity,
              商品总价: item.amount
            })))
          }
        }

        // 获取总数，使用正确的字段
        const total = res.data.total || 0
        const hasMore = orders.length > 0 && this.data.page * this.data.pageSize < total

        console.log('[订单列表] 处理后的订单数据:', {
          orders,
          ordersLength: orders.length,
          currentPage: this.data.page,
          pageSize: this.data.pageSize,
          total,
          hasMore
        })

        // 更新状态
        this.setData({
          orders: this.data.page === 1 ? orders : [...this.data.orders, ...orders],
          hasMore
        })

        // 设置待支付订单的倒计时
        this.setupOrderCountdowns()
      } else {
        console.error('[订单列表] 获取订单失败:', res.msg || '未知错误')
        wx.showToast({
          title: '获取订单失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[订单列表] 获取订单列表异常:', error.message || error)
      wx.showToast({
        title: error.message || '获取订单失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 为待支付订单设置倒计时
  setupOrderCountdowns() {
    const { orders, countdownTimers } = this.data
    
    // 遍历订单，查找待支付状态的订单
    orders.forEach(order => {
      if (order.status === 0) {
        // 清除可能已存在的该订单定时器
        if (countdownTimers[order.id]) {
          clearInterval(countdownTimers[order.id])
        }
        
        // 启动新的倒计时
        this.startOrderCountdown(order)
      }
    })
  },

  // 为单个订单启动倒计时
  startOrderCountdown(order) {
    // 获取订单创建时间
    const createTime = new Date(this.formatDateString(order.create_time)).getTime()
    
    // 计算过期时间
    const expireTime = createTime + ORDER_PAYMENT_TIMEOUT
    
    // 设置倒计时定时器
    const timerId = setInterval(() => {
      // 计算剩余时间
      const now = new Date().getTime()
      const remainTime = expireTime - now
      
      if (remainTime <= 0) {
        // 倒计时结束，清除计时器
        clearInterval(timerId)
        
        // 更新订单状态
        this.updateOrderCountdownFinished(order.id)
        return
      }
      
      // 计算剩余分钟和秒数
      const minutes = Math.floor(remainTime / (60 * 1000))
      const seconds = Math.floor((remainTime % (60 * 1000)) / 1000)
      
      // 格式化倒计时文本
      const countdownText = `${minutes < 10 ? '0' + minutes : minutes}:${seconds < 10 ? '0' + seconds : seconds}`
      
      // 更新订单的倒计时显示
      this.updateOrderCountdown(order.id, countdownText)
    }, 1000)
    
    // 保存定时器ID
    const newTimers = { ...this.data.countdownTimers, [order.id]: timerId }
    this.setData({ countdownTimers: newTimers })
  },

  // 更新订单的倒计时显示
  updateOrderCountdown(orderId, countdownText) {
    const { orders } = this.data
    const index = orders.findIndex(order => order.id === orderId)
    
    if (index !== -1) {
      const newOrders = [...orders]
      newOrders[index].countdown = countdownText
      
      this.setData({
        orders: newOrders
      })
    }
  },

  // 处理订单倒计时结束
  updateOrderCountdownFinished(orderId) {
    const { orders, countdownTimers } = this.data
    const index = orders.findIndex(order => order.id === orderId)
    
    if (index !== -1) {
      // 清除定时器
      clearInterval(countdownTimers[orderId])
      
      // 从倒计时器列表中移除
      const newTimers = { ...countdownTimers }
      delete newTimers[orderId]
      
      // 修改订单状态为已取消（如果需要的话）
      // 注意：实际状态会由后端定时任务处理，这里只是为了提升用户体验，立即更新UI
      const newOrders = [...orders]
      if (newOrders[index].status === 0) {
        newOrders[index].status = 4 // 使用4表示已取消状态
        newOrders[index].status_text = this.getStatusText(4)
        newOrders[index].countdown = null
      }
      
      this.setData({
        orders: newOrders,
        countdownTimers: newTimers
      })
      
      // 用户体验：提示用户订单已超时
      wx.showToast({
        title: '订单支付超时已自动取消',
        icon: 'none'
      })
    }
  },

  getStatusText(status) {
    switch (parseInt(status)) {
      case 0: return '待付款'
      case 1: return '待发货'
      case 2: return '待收货'
      case 3: return '已完成'
      case 4: return '已取消'
      default: return '未知状态'
    }
  },

  // 跳转到订单详情
  goToDetail(e) {
    const { id } = e.currentTarget.dataset
    if (!id) {
      console.error('[订单列表] 无效的订单ID')
      return
    }
    wx.navigateTo({
      url: `/pages/order/detail/index?orderId=${id}`
    })
  },

  // 支付订单
  async onPayOrder(e) {
    const { id } = e.currentTarget.dataset
    if (!id) {
      console.error('[订单列表] 无效的订单ID')
      return
    }

    try {
      console.log('[订单列表] 更新订单状态:', { orderId: id })
      const result = await updateOrderStatus(id)
      console.log('[订单列表] 更新订单状态响应:', result)

      if (result && result.code === 0) {
        // 显示支付成功提示
        wx.showToast({
          title: '支付成功',
          icon: 'success',
          duration: 2000
        })

        // 延迟1.5秒后刷新订单列表
        setTimeout(() => {
          this.loadOrders()
        }, 1500)
      } else {
        throw new Error(result.msg || '支付失败')
      }
    } catch (error) {
      console.error('[订单列表] 支付失败:', error)
      wx.showToast({
        title: error.message || '支付失败',
        icon: 'none'
      })
    }
  },

  // 取消订单
  async onCancelOrder(e) {
    const { id } = e.currentTarget.dataset
    if (!id) {
      console.error('[订单列表] 无效的订单ID')
      return
    }
    
    try {
      const res = await wx.showModal({
        title: '提示',
        content: '确定要取消该订单吗？',
        confirmText: '确定',
        cancelText: '取消'
      })

      if (res.confirm) {
        console.log('[订单列表] 取消订单:', { orderId: id })
        await this.cancelOrder(id)
      }
    } catch (error) {
      console.error('[订单列表] 取消订单失败:', error)
      wx.showToast({
        title: error.message || '取消失败',
        icon: 'none'
      })
    }
  },

  // 确认收货
  async onConfirmOrder(e) {
    const { id } = e.currentTarget.dataset
    if (!id) {
      console.error('[订单列表] 无效的订单ID')
      return
    }
    
    try {
      const res = await wx.showModal({
        title: '提示',
        content: '确认已收到商品？',
        confirmText: '确认',
        cancelText: '取消'
      })

      if (res.confirm) {
        console.log('[订单列表] 确认收货:', { orderId: id })
        await this.confirmOrder(id)
      }
    } catch (error) {
      console.error('[订单列表] 确认收货失败:', error)
      wx.showToast({
        title: error.message || '确认失败',
        icon: 'none'
      })
    }
  },

  // 查看物流
  onViewLogistics(e) {
    const orderId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/order/logistics/index?id=${orderId}`
    })
  },

  // 再次购买
  onBuyAgain(e) {
    const orderId = e.currentTarget.dataset.id
    // 跳转到确认订单页面，并传入订单ID
    wx.navigateTo({
      url: `/pages/order/confirm/index?order_id=${orderId}&type=rebuy`
    })
  },

  // 下拉刷新
  onPullDownRefresh() {
    console.log('[订单列表] 触发下拉刷新')
    this.setData({
      page: 1,
      hasMore: true
    }, () => {
      this.loadOrders().then(() => {
        wx.stopPullDownRefresh()
      })
    })
  },

  // 重置分页
  resetPage() {
    console.log('[订单列表] 重置分页', {
      原页码: this.data.page,
      原数据量: this.data.orders.length,
      原加载状态: this.data.hasMore
    })
    
    this.setData({
      page: 1,
      hasMore: true,
      orders: []
    })
  },

  // 触底加载更多
  onReachBottom() {
    console.log('[订单列表] 触发触底加载', {
      当前页: this.data.page,
      是否加载中: this.data.loading,
      是否有更多: this.data.hasMore
    })

    // 如果没有更多数据，直接返回
    if (!this.data.hasMore) {
      console.log('[订单列表] 没有更多数据了')
      wx.showToast({
        title: '没有更多订单了',
        icon: 'none',
        duration: 1500
      })
      return
    }

    // 如果正在加载中，跳过重复请求
    if (this.data.loading) {
      console.log('[订单列表] 正在加载中，跳过重复请求')
      return
    }

    const nextPage = this.data.page + 1
    console.log('[订单列表] 加载下一页:', nextPage)
    
    this.setData({
      page: nextPage
    }, () => {
      this.loadOrders()
    })
  },

  async cancelOrder(orderId) {
    try {
      const res = await cancelOrder(orderId)
      console.log('[订单列表] 取消订单响应:', res)

      if (res.code === 0) {
        wx.showToast({
          title: '订单已取消',
          icon: 'success'
        })
        // 刷新订单列表
        this.loadOrders()
      } else {
        let errorMsg = '取消订单失败'
        switch(res.code) {
          case 20001:
            errorMsg = '订单不存在'
            break
          case 20003:
            errorMsg = '订单状态无效'
            break
          case 20004:
            errorMsg = '订单更新失败'
            break
          default:
            errorMsg = res.msg || '取消订单失败'
        }
        wx.showToast({
          title: errorMsg,
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[订单列表] 取消订单失败:', error)
      wx.showToast({
        title: error.message || '取消订单失败',
        icon: 'none'
      })
    }
  },

  async confirmOrder(orderId) {
    try {
      const res = await confirmOrder(orderId)
      console.log('[订单列表] 确认收货响应:', res)

      if (res.code === 0) {
        wx.showToast({
          title: '已确认收货',
          icon: 'success'
        })
        // 刷新订单列表
        this.loadOrders()
      } else {
        let errorMsg = '确认收货失败'
        switch(res.code) {
          case 20001:
            errorMsg = '订单不存在'
            break
          case 20003:
            errorMsg = '订单状态无效'
            break
          case 20004:
            errorMsg = '订单更新失败'
            break
          default:
            errorMsg = res.msg || '确认收货失败'
        }
        wx.showToast({
          title: errorMsg,
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[订单列表] 确认收货失败:', error)
      wx.showToast({
        title: error.message || '确认收货失败',
        icon: 'none'
      })
    }
  }
})