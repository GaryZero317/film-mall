import { getFilmOrderList } from '../../../api/film'

Page({
  data: {
    filmOrders: [],
    loading: false,
    page: 1,
    pageSize: 10,
    total: 0,
    hasMore: true,
    currentStatus: -1, // -1 表示全部状态
    tabs: [
      { id: -1, name: '全部' },
      { id: 0, name: '待付款' },
      { id: 1, name: '冲洗中' },
      { id: 2, name: '待收货' },
      { id: 3, name: '已完成' }
    ]
  },

  onLoad(options) {
    // 如果有传入的状态，则设置当前状态
    if (options.status) {
      this.setData({
        currentStatus: parseInt(options.status)
      })
    }
    this.loadFilmOrders()
  },

  onPullDownRefresh() {
    this.refreshOrders()
  },

  onReachBottom() {
    if (this.data.hasMore) {
      this.loadMoreOrders()
    }
  },

  // 切换状态标签
  onTabChange(e) {
    const status = e.currentTarget.dataset.id
    this.setData({
      currentStatus: status,
      filmOrders: [],
      page: 1,
      hasMore: true
    })
    this.loadFilmOrders()
  },

  // 加载订单数据
  async loadFilmOrders() {
    if (this.data.loading) return
    
    try {
      this.setData({ loading: true })
      
      const params = {
        page: this.data.page,
        page_size: this.data.pageSize
      }
      
      // 添加状态过滤
      if (this.data.currentStatus !== -1) {
        params.status = this.data.currentStatus
      }
      
      const res = await getFilmOrderList(params)
      
      if (res.code === 0 && res.data) {
        const { list, total } = res.data
        
        this.setData({
          filmOrders: [...this.data.filmOrders, ...list],
          total,
          hasMore: this.data.filmOrders.length + list.length < total
        })
      }
    } catch (error) {
      console.error('获取胶片冲洗订单列表失败:', error)
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
      wx.stopPullDownRefresh()
    }
  },

  // 刷新订单
  refreshOrders() {
    this.setData({
      filmOrders: [],
      page: 1,
      hasMore: true
    })
    this.loadFilmOrders()
  },
  
  // 加载更多订单
  loadMoreOrders() {
    this.setData({
      page: this.data.page + 1
    })
    this.loadFilmOrders()
  },
  
  // 查看订单详情
  viewOrderDetail(e) {
    const orderId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/user/film/detail/index?id=${orderId}`
    })
  },
  
  // 返回上一页
  onBack() {
    wx.navigateBack()
  }
}) 