import { getFilmOrderDetail } from '../../../../api/film'

Page({
  data: {
    orderId: null,
    orderDetail: null,
    loading: false,
    statusMap: {
      0: { text: '待付款', color: '#ff9900' },
      1: { text: '冲洗处理中', color: '#1890ff' },
      2: { text: '待收货', color: '#52c41a' },
      3: { text: '已完成', color: '#666666' }
    }
  },

  onLoad(options) {
    if (options.id) {
      this.setData({
        orderId: options.id
      })
      this.loadOrderDetail()
    } else {
      wx.showToast({
        title: '订单ID不存在',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    }
  },

  // 加载订单详情
  async loadOrderDetail() {
    if (!this.data.orderId || this.data.loading) return
    
    try {
      this.setData({ loading: true })
      
      const res = await getFilmOrderDetail(this.data.orderId)
      
      if (res.code === 0 && res.data) {
        this.setData({
          orderDetail: res.data
        })
      } else {
        wx.showToast({
          title: res.msg || '获取订单详情失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('获取胶片冲洗订单详情失败:', error)
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },
  
  // 复制订单号
  copyOrderNo() {
    if (this.data.orderDetail && this.data.orderDetail.foid) {
      wx.setClipboardData({
        data: this.data.orderDetail.foid,
        success: () => {
          wx.showToast({
            title: '订单号已复制',
            icon: 'success'
          })
        }
      })
    }
  },
  
  // 返回上一页
  onBack() {
    wx.navigateBack()
  },
  
  // 支付订单
  payOrder() {
    if (!this.data.orderDetail) return
    
    // 调用支付接口
    wx.showToast({
      title: '暂未开通支付功能',
      icon: 'none'
    })
  },
  
  // 查看物流
  viewLogistics() {
    wx.showToast({
      title: '暂未开通物流查询',
      icon: 'none'
    })
  },
  
  // 确认收货
  confirmReceive() {
    if (!this.data.orderDetail) return
    
    wx.showModal({
      title: '确认收货',
      content: '确定已收到商品？',
      success: (res) => {
        if (res.confirm) {
          wx.showToast({
            title: '确认收货成功',
            icon: 'success'
          })
          // 重新加载订单详情
          this.loadOrderDetail()
        }
      }
    })
  }
}) 