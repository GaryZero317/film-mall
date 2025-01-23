// pages/order/confirm/index.js
import { createOrder } from '../../../api/order'
import { getAddressList } from '../../../api/user'
import { loginGuard } from '../../../utils/auth'

Page(loginGuard({
  data: {
    address: null,
    orderItems: [],
    totalPrice: 0,
    totalCount: 0,
    remark: '',
    loading: false
  },

  onLoad(options) {
    // 从购物车结算进入
    const { from } = options
    if (from === 'cart') {
      const cartItems = wx.getStorageSync('selectedCartItems') || []
      this.setOrderItems(cartItems)
    } 
    // 从商品详情页直接购买进入
    else {
      const { productId, quantity } = options
      if (productId && quantity) {
        this.setOrderItems([{
          product_id: parseInt(productId),
          quantity: parseInt(quantity)
        }])
      }
    }
    this.loadAddress()
  },

  // 加载收货地址
  async loadAddress() {
    try {
      const res = await getAddressList()
      // 获取默认地址
      const defaultAddress = res.data.find(item => item.is_default) || res.data[0]
      this.setData({ address: defaultAddress })
    } catch (error) {
      console.error('加载地址失败:', error)
    }
  },

  // 设置订单商品
  setOrderItems(items) {
    const orderItems = items.map(item => ({
      id: item.product_id,
      name: item.name,
      price: item.price,
      quantity: item.quantity,
      cover_image: item.cover_image
    }))

    const totalPrice = orderItems.reduce((sum, item) => sum + item.price * item.quantity, 0)
    const totalCount = orderItems.reduce((sum, item) => sum + item.quantity, 0)

    this.setData({
      orderItems,
      totalPrice,
      totalCount
    })
  },

  // 选择收货地址
  onSelectAddress() {
    wx.navigateTo({
      url: '/pages/address/list/index?from=order'
    })
  },

  // 备注输入
  onRemarkInput(e) {
    this.setData({
      remark: e.detail.value
    })
  },

  // 提交订单
  async submitOrder() {
    const { address, orderItems, remark } = this.data
    if (!address) {
      wx.showToast({
        title: '请选择收货地址',
        icon: 'none'
      })
      return
    }

    if (orderItems.length === 0) {
      wx.showToast({
        title: '订单商品不能为空',
        icon: 'none'
      })
      return
    }

    try {
      this.setData({ loading: true })
      const res = await createOrder({
        address_id: address.id,
        items: orderItems.map(item => ({
          product_id: item.id,
          quantity: item.quantity
        })),
        remark
      })

      // 下单成功后清除购物车中已购买的商品
      const cartItems = wx.getStorageSync('cartItems') || []
      const newCartItems = cartItems.filter(item => 
        !orderItems.find(orderItem => orderItem.id === item.product_id)
      )
      wx.setStorageSync('cartItems', newCartItems)

      // 跳转到支付页面
      wx.navigateTo({
        url: `/pages/order/payment/index?orderId=${res.data.id}`
      })
    } catch (error) {
      console.error('创建订单失败:', error)
      wx.showToast({
        title: '创建订单失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  }
}))