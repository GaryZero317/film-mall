// pages/order/list/index.js
import { getOrderList, cancelOrder, confirmOrder, updateOrderStatus } from '../../../api/order'
import { getProductDetail, getProductImages } from '../../../api/product'

const app = getApp()

Page({
  data: {
    orders: [],
    loading: false,
    currentTab: 0,
    tabs: [
      { id: 0, name: '全部' },
      { id: 1, name: '待付款' },
      { id: 2, name: '待发货' },
      { id: 3, name: '待收货' },
      { id: 4, name: '已完成' }
    ],
    isFirstLoad: true,  // 添加标记，用于区分首次加载和后续显示
    page: 1,           // 当前页码
    pageSize: 10,      // 每页数量
    hasMore: true      // 是否还有更多数据
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

    try {
      // 从全局获取用户信息
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo || !userInfo.id) {
        throw new Error('请先登录')
      }

      const res = await getOrderList({
        uid: userInfo.id,
        status: this.data.currentTab,
        page: this.data.page,
        pageSize: this.data.pageSize
      })

      console.log('[订单列表] 获取订单列表响应:', res)

      if (res.code === 0) {
        const orders = res.data?.list || []
        // 处理订单数据
        const formattedOrders = orders.map(order => ({
          id: order.id,
          oid: order.oid,
          status: order.status,
          status_text: this.getStatusText(order.status),
          amount: (order.total_price / 100).toFixed(2),
          create_time: this.formatDateString(order.create_time),
          items: order.items.map(item => ({
            id: item.id,
            product_name: item.product_name,
            product_image: item.product_image ? 
              (item.product_image.startsWith('http') ? 
                item.product_image : 
                `http://localhost:8001${item.product_image}`
              ) : 
              'http://localhost:8001/uploads/placeholder.png',
            price: (item.price / 100).toFixed(2),
            quantity: item.quantity
          }))
        }))
        
        // 按创建时间排序，最新的在前面
        formattedOrders.sort((a, b) => {
          return new Date(b.create_time) - new Date(a.create_time)
        })
        
        this.setData({
          orders: this.data.page === 1 ? formattedOrders : [...this.data.orders, ...formattedOrders],
          hasMore: orders.length === this.data.pageSize
        })
      } else {
        let errorMsg = '获取订单列表失败'
        switch(res.code) {
          case 10001:
            errorMsg = '参数错误'
            break
          case 10002:
            errorMsg = '服务器内部错误'
            break
          default:
            errorMsg = res.msg || '获取订单列表失败'
        }
        wx.showToast({
          title: errorMsg,
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[订单列表] 加载订单失败:', error)
      
      // 如果是未登录错误，跳转到登录页面
      if (error.message === '请先登录') {
        wx.showToast({
          title: '请先登录',
          icon: 'none',
          duration: 2000,
          complete: () => {
            setTimeout(() => {
              wx.navigateTo({
                url: '/pages/login/index'
              })
            }, 1000)
          }
        })
      } else {
        wx.showToast({
          title: error.message || '加载订单失败',
          icon: 'none'
        })
      }
    } finally {
      this.setData({ 
        loading: false,
        refreshing: false
      })
    }
  },

  // 获取状态文本
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