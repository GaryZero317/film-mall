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
    if (this.data.loading) {
      console.log('[订单列表] 正在加载中，跳过重复请求', {
        currentTab: this.data.currentTab,
        page: this.data.page,
        hasMore: this.data.hasMore
      })
      return
    }

    try {
      console.log('[订单列表] 开始加载订单', {
        currentTab: this.data.currentTab,
        page: this.data.page,
        pageSize: this.data.pageSize,
        hasMore: this.data.hasMore
      })
      
      this.setData({ loading: true })
      
      // 获取用户信息
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo || !userInfo.id) {
        console.error('[订单列表] 用户未登录')
        throw new Error('请先登录')
      }

      // 构造请求参数
      const params = {
        uid: parseInt(userInfo.id),
        page: this.data.page,
        pageSize: this.data.pageSize
      }

      // 根据当前标签添加状态过滤
      if (this.data.currentTab > 0) {
        // 将界面状态映射到后端状态
        const statusMap = {
          1: 0, // 待付款
          2: 1, // 待发货
          3: 2, // 待收货
          4: 3  // 已完成
        }
        params.status = statusMap[this.data.currentTab]
      }

      console.log('[订单列表] 请求参数:', params)
      const res = await getOrderList(params)
      console.log('[订单列表] 响应数据:', res)

      if (!res || res.code !== 0 || !res.data) {
        console.error('[订单列表] 响应数据异常:', res)
        throw new Error(res?.msg || '加载订单失败')
      }

      // 处理订单数据
      const orders = Array.isArray(res.data) ? res.data : []
      console.log('[订单列表] 原始订单数据:', orders)

      // 获取所有订单中的商品详情
      const processOrders = orders.map(async (order) => {
        try {
          console.log('[订单列表] 获取商品详情:', { 订单ID: order.id, 商品ID: order.pid })
          const [productRes, imageRes] = await Promise.all([
            getProductDetail(order.pid),
            getProductImages(order.pid)
          ])
          
          if (productRes && productRes.code === 0 && productRes.data) {
            const product = productRes.data
            
            // 处理商品图片
            let imageUrl = ''
            if (imageRes && imageRes.code === 0 && imageRes.data) {
              const mainImage = imageRes.data.find(img => img.isMain)
              imageUrl = mainImage ? mainImage.imageUrl : (imageRes.data[0]?.imageUrl || '')
            }
            
            console.log('[订单列表] 商品详情:', { 
              订单ID: order.id,
              商品ID: product.id,
              商品名称: product.name,
              商品图片: imageUrl,
              商品价格: product.price
            })
            
            // 使用商品详情更新订单数据
            order.product_name = product.name
            order.cover_image = imageUrl
          } else {
            console.warn('[订单列表] 获取商品详情失败:', productRes)
          }
        } catch (error) {
          console.error('[订单列表] 获取商品详情异常:', error)
        }

        // 转换金额（分到元）
        order.amount = (order.amount / 100).toFixed(2)
        
        // 格式化时间
        if (order.create_time) {
          const date = new Date(order.create_time)
          order.create_time = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
        }
        
        // 设置状态文本
        order.status_text = this.getStatusText(order.status)
        
        // 构造商品列表
        order.items = [{
          id: order.pid,
          name: order.product_name || '未知商品',
          cover_image: order.cover_image || '/assets/images/default.png',
          price: order.amount,
          quantity: order.quantity || 1
        }]

        // 处理商品图片
        order.items.forEach(item => {
          console.log('[订单列表] 处理商品:', {
            商品ID: item.id,
            商品名称: item.name,
            原始图片: item.cover_image,
            价格: item.price,
            数量: item.quantity
          })

          if (!item.id) {
            console.warn('[订单列表] 商品ID无效:', item)
          }

          if (item.cover_image && !item.cover_image.startsWith('http')) {
            item.cover_image = `http://localhost:8001${item.cover_image}`
          }
        })

        return order
      })

      // 等待所有订单处理完成
      const processedOrders = await Promise.all(processOrders)
      console.log('[订单列表] 处理后的订单数据:', processedOrders)

      // 更新订单列表和分页状态
      const newOrders = this.data.page === 1 ? processedOrders : [...this.data.orders, ...processedOrders]
      const hasMore = processedOrders.length >= this.data.pageSize

      console.log('[订单列表] 更新状态:', {
        是否第一页: this.data.page === 1,
        原订单数量: this.data.orders.length,
        新增订单数量: processedOrders.length,
        更新后总数量: newOrders.length,
        是否有更多: hasMore,
        订单状态统计: processedOrders.reduce((acc, order) => {
          acc[order.status_text] = (acc[order.status_text] || 0) + 1
          return acc
        }, {})
      })

      this.setData({
        orders: newOrders,
        loading: false,
        hasMore
      })
    } catch (error) {
      console.error('[订单列表] 加载订单列表失败:', error)
      this.setData({ 
        orders: [],
        loading: false,
        hasMore: false
      })
      wx.showToast({
        title: error.message || '加载订单失败',
        icon: 'none'
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
        const result = await cancelOrder(id)
        console.log('[订单列表] 取消订单响应:', result)

        if (result && result.code === 0) {
          wx.showToast({
            title: '订单已取消',
            icon: 'success'
          })
          // 刷新订单列表
          this.loadOrders()
        } else {
          throw new Error(result.msg || '取消订单失败')
        }
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
        const result = await confirmOrder(id)
        console.log('[订单列表] 确认收货响应:', result)

        if (result && result.code === 0) {
          wx.showToast({
            title: '确认成功',
            icon: 'success'
          })
          // 刷新订单列表
          this.loadOrders()
        } else {
          throw new Error(result.msg || '确认收货失败')
        }
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
    this.loadOrders().then(() => {
      wx.stopPullDownRefresh()
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
    console.log('[订单列表] 触底加载', {
      当前页码: this.data.page,
      是否加载中: this.data.loading,
      是否有更多: this.data.hasMore,
      当前数据量: this.data.orders.length
    })

    if (this.data.hasMore && !this.data.loading) {
      const nextPage = this.data.page + 1
      console.log('[订单列表] 加载下一页:', nextPage)
      
      this.setData({
        page: nextPage
      }, () => {
        this.loadOrders()
      })
    } else {
      console.log('[订单列表] 不满足加载条件，跳过加载')
    }
  }
})