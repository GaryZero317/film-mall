// pages/order/list/index.js
import { getOrderList, cancelOrder, confirmOrder, deleteOrder } from '../../../api/order'

const app = getApp()

Page({
  data: {
    tabs: [
      { id: 'all', name: '全部' },
      { id: 'unpaid', name: '待付款' },
      { id: 'undelivered', name: '待发货' },
      { id: 'delivered', name: '待收货' },
      { id: 'completed', name: '已完成' }
    ],
    activeTab: 'all',
    orders: [],
    loading: false,
    hasMore: true,
    page: 1,
    pageSize: 10
  },

  onLoad(options) {
    // 如果从其他页面传入了tab参数，则切换到对应标签
    if (options.tab) {
      this.setData({ activeTab: options.tab })
    }
    this.loadOrders()
  },

  // 切换标签
  onTabChange(e) {
    const tab = e.currentTarget.dataset.tab
    this.setData({
      activeTab: tab,
      orders: [],
      page: 1,
      hasMore: true
    })
    this.loadOrders()
  },

  // 加载订单列表
  async loadOrders() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      const { page, pageSize, activeTab } = this.data
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/orders`,
        method: 'GET',
        data: {
          page,
          page_size: pageSize,
          status: activeTab === 'all' ? '' : activeTab
        },
        header: {
          'Authorization': `Bearer ${wx.getStorageSync('token')}`
        }
      })

      if (res.statusCode === 200) {
        const { data, total } = res.data
        const formattedOrders = data.map(order => ({
          ...order,
          status_text: this.getStatusText(order.status)
        }))

        this.setData({
          orders: [...this.data.orders, ...formattedOrders],
          page: page + 1,
          hasMore: this.data.orders.length + formattedOrders.length < total
        })
      } else {
        wx.showToast({
          title: '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载订单列表失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 获取订单状态文本
  getStatusText(status) {
    const statusMap = {
      unpaid: '待付款',
      undelivered: '待发货',
      delivered: '待收货',
      completed: '已完成',
      cancelled: '已取消'
    }
    return statusMap[status] || status
  },

  // 查看订单详情
  onViewDetail(e) {
    const orderId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/order/detail/index?id=${orderId}`
    })
  },

  // 支付订单
  async onPayOrder(e) {
    const orderId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/order/payment/index?id=${orderId}`
    })
  },

  // 取消订单
  async onCancelOrder(e) {
    const orderId = e.currentTarget.dataset.id
    
    const res = await wx.showModal({
      title: '提示',
      content: '确定要取消该订单吗？',
      confirmText: '确定',
      cancelText: '取消'
    })

    if (res.confirm) {
      try {
        const result = await wx.request({
          url: `${app.globalData.baseUrl}/api/orders/${orderId}/cancel`,
          method: 'POST',
          header: {
            'Authorization': `Bearer ${wx.getStorageSync('token')}`
          }
        })

        if (result.statusCode === 200) {
          wx.showToast({
            title: '订单已取消',
            icon: 'success'
          })
          // 刷新订单列表
          this.setData({
            orders: [],
            page: 1,
            hasMore: true
          })
          this.loadOrders()
        } else {
          wx.showToast({
            title: '取消失败',
            icon: 'none'
          })
        }
      } catch (error) {
        console.error('取消订单失败:', error)
        wx.showToast({
          title: '取消失败',
          icon: 'none'
        })
      }
    }
  },

  // 确认收货
  async onConfirmOrder(e) {
    const orderId = e.currentTarget.dataset.id
    
    const res = await wx.showModal({
      title: '提示',
      content: '确认已收到商品？',
      confirmText: '确认',
      cancelText: '取消'
    })

    if (res.confirm) {
      try {
        const result = await wx.request({
          url: `${app.globalData.baseUrl}/api/orders/${orderId}/confirm`,
          method: 'POST',
          header: {
            'Authorization': `Bearer ${wx.getStorageSync('token')}`
          }
        })

        if (result.statusCode === 200) {
          wx.showToast({
            title: '确认成功',
            icon: 'success'
          })
          // 刷新订单列表
          this.setData({
            orders: [],
            page: 1,
            hasMore: true
          })
          this.loadOrders()
        } else {
          wx.showToast({
            title: '确认失败',
            icon: 'none'
          })
        }
      } catch (error) {
        console.error('确认收货失败:', error)
        wx.showToast({
          title: '确认失败',
          icon: 'none'
        })
      }
    }
  },

  // 删除订单
  async onDeleteOrder(e) {
    const orderId = e.currentTarget.dataset.id
    
    const res = await wx.showModal({
      title: '提示',
      content: '确定要删除该订单吗？',
      confirmText: '删除',
      cancelText: '取消'
    })

    if (res.confirm) {
      try {
        const result = await wx.request({
          url: `${app.globalData.baseUrl}/api/orders/${orderId}`,
          method: 'DELETE',
          header: {
            'Authorization': `Bearer ${wx.getStorageSync('token')}`
          }
        })

        if (result.statusCode === 200) {
          wx.showToast({
            title: '删除成功',
            icon: 'success'
          })
          // 刷新订单列表
          this.setData({
            orders: [],
            page: 1,
            hasMore: true
          })
          this.loadOrders()
        } else {
          wx.showToast({
            title: '删除失败',
            icon: 'none'
          })
        }
      } catch (error) {
        console.error('删除订单失败:', error)
        wx.showToast({
          title: '删除失败',
          icon: 'none'
        })
      }
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
  async onPullDownRefresh() {
    this.setData({
      orders: [],
      page: 1,
      hasMore: true
    })
    await this.loadOrders()
    wx.stopPullDownRefresh()
  },

  // 上拉加载更多
  onReachBottom() {
    this.loadOrders()
  }
})