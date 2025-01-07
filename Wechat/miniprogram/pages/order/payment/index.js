import { getOrderDetail, payOrder } from '../../../api/order'

Page({
  data: {
    order: null,
    loading: false,
    paymentMethods: [
      { id: 'wxpay', name: '微信支付', icon: '/assets/images/wxpay.png' },
      { id: 'alipay', name: '支付宝支付', icon: '/assets/images/alipay.png' },
      { id: 'balance', name: '余额支付', icon: '/assets/images/balance.png' }
    ],
    selectedPayment: 'wxpay'
  },

  onLoad(options) {
    const { orderId } = options
    if (orderId) {
      this.loadOrderDetail(orderId)
    }
  },

  // 加载订单详情
  async loadOrderDetail(orderId) {
    try {
      this.setData({ loading: true })
      const res = await getOrderDetail(orderId)
      this.setData({ 
        order: res.data,
        loading: false
      })
    } catch (error) {
      console.error('加载订单详情失败:', error)
      this.setData({ loading: false })
      wx.showToast({
        title: '加载订单失败',
        icon: 'none'
      })
    }
  },

  // 选择支付方式
  onSelectPayment(e) {
    const { method } = e.currentTarget.dataset
    this.setData({
      selectedPayment: method
    })
  },

  // 确认支付
  async onPayment() {
    const { order, selectedPayment } = this.data
    if (!order) return

    try {
      this.setData({ loading: true })
      const res = await payOrder(order.id)
      
      // 根据选择的支付方式调用不同的支付接口
      if (selectedPayment === 'wxpay') {
        const payParams = res.data
        wx.requestPayment({
          timeStamp: payParams.timeStamp,
          nonceStr: payParams.nonceStr,
          package: payParams.package,
          signType: payParams.signType,
          paySign: payParams.paySign,
          success: () => {
            wx.showToast({
              title: '支付成功',
              icon: 'success'
            })
            // 支付成功后跳转到订单列表
            setTimeout(() => {
              wx.redirectTo({
                url: '/pages/order/list/index'
              })
            }, 1500)
          },
          fail: () => {
            wx.showToast({
              title: '支付失败',
              icon: 'none'
            })
          }
        })
      } else {
        // 其他支付方式的处理逻辑
        wx.showToast({
          title: '暂不支持该支付方式',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('支付失败:', error)
      wx.showToast({
        title: '支付失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 取消支付
  onCancel() {
    wx.showModal({
      title: '提示',
      content: '确定要取消支付吗？',
      success: (res) => {
        if (res.confirm) {
          wx.navigateBack()
        }
      }
    })
  }
}) 