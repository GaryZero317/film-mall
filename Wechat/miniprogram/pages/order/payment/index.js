import { getOrderDetail, payOrder, queryPayStatus } from '../../../api/order'
import { getProductDetail } from '../../../api/product'

Page({
  data: {
    order: null,
    loading: false
  },

  onLoad(options) {
    const { orderId } = options
    if (orderId) {
      // 确保orderId是数字类型
      const id = parseInt(orderId)
      if (isNaN(id)) {
        console.error('订单ID无效:', orderId)
        wx.showToast({
          title: '订单ID无效',
          icon: 'none'
        })
        return
      }
      console.log('[订单支付] 开始加载订单详情:', { 原始ID: orderId, 转换后ID: id })
      this.loadOrderDetail(id)
    }
  },

  // 加载订单详情
  async loadOrderDetail(orderId) {
    try {
      this.setData({ loading: true })
      console.log('[订单支付] 请求订单详情:', { 订单ID: orderId })
      const res = await getOrderDetail(orderId)
      console.log('[订单支付] 获取订单详情成功:', res)
      
      if (!res || !res.data) {
        throw new Error('订单数据无效')
      }

      // 处理订单数据，将金额从分转换为元
      const orderData = res.data
      // 转换订单总价
      orderData.total_price = (orderData.amount / 100).toFixed(2)
      // 转换订单金额
      orderData.amount = (orderData.amount / 100).toFixed(2)
      // 确保数量字段存在
      orderData.quantity = orderData.quantity || 1
      // 设置总数量
      orderData.total_count = orderData.quantity
      // 格式化创建时间
      if (orderData.create_time) {
        const date = new Date(orderData.create_time)
        orderData.create_time = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
      }
      // 确保订单编号存在
      orderData.order_no = orderData.order_no || orderData.id

      // 获取商品详情
      try {
        console.log('[订单支付] 开始获取商品详情:', { 
          商品ID: orderData.pid,
          订单数量: orderData.quantity,
          订单金额: orderData.amount,
          订单总价: orderData.total_price,
          原始数据: orderData
        })
        const productRes = await getProductDetail(orderData.pid)
        console.log('[订单支付] 获取商品详情成功:', productRes)
        
        if (productRes && productRes.data) {
          const product = productRes.data
          console.log('[订单支付] 商品原始数据:', product)

          // 处理商品图片路径
          let coverImage = '/assets/images/default.png'
          const possibleImageFields = ['cover_image', 'main_image', 'image', 'img', 'imgUrl', 'coverImage', 'mainImage', 'imageUrl']
          
          for (const field of possibleImageFields) {
            if (product[field]) {
              const imagePath = product[field]
              console.log(`[订单支付] 尝试图片字段 ${field}:`, imagePath)
              coverImage = imagePath.startsWith('http') ? 
                imagePath : 
                `http://localhost:8001${imagePath}`
              break
            }
          }
          
          // 计算单价（总金额除以数量）
          const amount = parseFloat(orderData.amount)
          const quantity = parseInt(orderData.quantity)
          const unitPrice = (amount / quantity).toFixed(2)
          
          console.log('[订单支付] 价格计算:', {
            订单金额: amount,
            订单数量: quantity,
            计算单价: unitPrice
          })
          
          // 构造商品信息
          orderData.items = [{
            id: orderData.pid,
            name: product.name || '未知商品',
            price: unitPrice,
            quantity: quantity,
            cover_image: coverImage,
            total_price: orderData.amount
          }]

          console.log('[订单支付] 商品信息:', {
            商品名称: product.name,
            商品图片: coverImage,
            单价: unitPrice,
            数量: quantity,
            总价: orderData.amount,
            原始数据: product
          })
        }
      } catch (error) {
        console.error('[订单支付] 获取商品详情失败:', error)
        // 使用默认商品信息
        orderData.items = [{
          id: orderData.pid,
          name: '未知商品',
          price: (orderData.amount / orderData.quantity).toFixed(2),
          quantity: orderData.quantity,
          cover_image: '/assets/images/default.png'
        }]
      }

      console.log('[订单支付] 处理后的订单数据:', orderData)

      this.setData({ 
        order: orderData,
        loading: false
      })
    } catch (error) {
      console.error('[订单支付] 加载订单详情失败:', error)
      this.setData({ loading: false })
      wx.showToast({
        title: error.message || '加载订单失败',
        icon: 'none'
      })
    }
  },

  // 确认支付
  async onPayment() {
    const { order } = this.data
    if (!order) return

    try {
      this.setData({ loading: true })
      
      // 创建支付记录
      const payData = {
        oid: order.id,
        uid: wx.getStorageSync('userInfo').id,
        amount: Math.round(parseFloat(order.amount) * 100) // 转换为分
      }
      console.log('[订单支付] 创建支付记录:', payData)
      const res = await payOrder(payData)
      console.log('[订单支付] 创建支付响应:', res)

      // 模拟支付成功
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
    } catch (error) {
      console.error('[订单支付] 发起支付失败:', error)
      wx.showToast({
        title: error.message || '支付失败',
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