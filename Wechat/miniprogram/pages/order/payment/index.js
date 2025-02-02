import { getOrderDetail, payOrder, payCallback, updateOrderStatus } from '../../../api/order'
import { getProductDetail } from '../../../api/product'

Page({
  data: {
    order: null,
    loading: false,
    orderStatus: {
      0: '待支付',
      1: '已支付',
      2: '已取消'
    },
    payStatus: {
      0: '未支付',
      1: '已支付',
      2: '支付失败'
    }
  },

  onLoad(options) {
    const { orderId, payId } = options
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
      console.log('[订单支付] 开始加载订单详情:', { 原始ID: orderId, 转换后ID: id, 支付ID: payId })
      // 保存支付ID到data中
      this.setData({ payId: payId })
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
      
      if (!res || res.code !== 0 || !res.data) {
        throw new Error('订单数据无效')
      }

      // 处理订单数据，将金额从分转换为元
      const orderData = res.data
      // 转换订单总价和金额（确保是数字类型）
      const totalPriceInCents = parseInt(orderData.total_price) || 0
      const amountInCents = parseInt(orderData.amount) || totalPriceInCents
      
      orderData.total_price = (totalPriceInCents / 100).toFixed(2)
      orderData.amount = (amountInCents / 100).toFixed(2)
      
      // 确保数量字段存在
      orderData.quantity = parseInt(orderData.quantity) || 1
      // 设置总数量
      orderData.total_count = orderData.quantity
      
      // 格式化创建时间
      if (orderData.create_time) {
        try {
          let dateStr = orderData.create_time
          if (typeof dateStr === 'string') {
            dateStr = dateStr.replace(' ', 'T')
          }
          const date = new Date(dateStr)
          
          if (isNaN(date.getTime())) {
            console.error('[订单支付] 无效的日期格式:', orderData.create_time)
            orderData.create_time = '日期格式无效'
          } else {
            orderData.create_time = `${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
          }
        } catch (error) {
          console.error('[订单支付] 日期格式化失败:', error)
          orderData.create_time = '日期格式化失败'
        }
      }
      
      // 确保订单编号存在
      orderData.order_no = orderData.order_no || orderData.oid || orderData.id

      // 获取商品详情（如果存在商品ID）
      if (orderData.pid) {
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
            
            // 处理商品图片路径
            let coverImage = '/assets/images/default.png'
            if (product.product_image) {
              coverImage = product.product_image.startsWith('http') 
                ? product.product_image 
                : `http://localhost:8001${product.product_image}`
            }
            
            // 计算单价（总金额除以数量）
            const unitPrice = (parseFloat(orderData.amount) / orderData.quantity).toFixed(2)
            
            // 构造商品信息
            orderData.items = [{
              id: orderData.pid,
              name: product.name || '未知商品',
              price: unitPrice,
              quantity: orderData.quantity,
              cover_image: coverImage,
              product_image: coverImage,
              total_price: orderData.amount
            }]
          }
        } catch (error) {
          console.error('[订单支付] 获取商品详情失败:', error)
          // 使用默认商品信息
          orderData.items = [{
            id: orderData.pid,
            name: '未知商品',
            price: (parseFloat(orderData.amount) / orderData.quantity).toFixed(2),
            quantity: orderData.quantity,
            cover_image: '/assets/images/default.png',
            product_image: '/assets/images/default.png',
            total_price: orderData.amount
          }]
        }
      } else if (orderData.items && orderData.items.length > 0) {
        // 如果订单中已经包含商品信息，直接使用
        orderData.items = orderData.items.map(item => ({
          ...item,
          price: (parseInt(item.price) / 100).toFixed(2),
          total_price: (parseInt(item.amount) / 100).toFixed(2),
          cover_image: item.product_image || item.cover_image || '/assets/images/default.png',
          product_image: item.product_image || item.cover_image || '/assets/images/default.png'
        }))
      } else {
        // 没有商品信息时使用默认值
        orderData.items = [{
          id: 0,
          name: '未知商品',
          price: orderData.amount,
          quantity: orderData.quantity,
          cover_image: '/assets/images/default.png',
          product_image: '/assets/images/default.png',
          total_price: orderData.amount
        }]
      }

      // 确保订单状态为数字类型
      orderData.status = parseInt(orderData.status) || 0
      orderData.pay_status = parseInt(orderData.pay_status) || 0
      
      // 添加订单状态显示
      orderData.statusText = this.data.orderStatus[orderData.status] || '未知状态'
      orderData.payStatusText = this.data.payStatus[orderData.pay_status] || '未知支付状态'

      console.log('[订单支付] 处理后的订单数据:', orderData)

      this.setData({ 
        order: orderData,
        loading: false
      })

      // 只有当支付状态为已支付时才自动跳转
      if (orderData.pay_status === 1) {
        setTimeout(() => {
          wx.redirectTo({
            url: '/pages/order/list/index'
          })
        }, 3000)
      }
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
    const { order, payId } = this.data
    if (!order) return

    try {
      this.setData({ loading: true })
      
      // 调用支付回调，更新支付状态
      const callbackData = {
        id: payId,  // 使用保存的支付记录ID
        uid: wx.getStorageSync('userInfo').id,
        oid: order.id,
        amount: Math.round(parseFloat(order.amount) * 100),
        source: 0,  // 支付来源：0表示小程序
        status: 1  // 支付状态：1表示已支付
      }
      console.log('[订单支付] 发起支付回调:', callbackData)
      const callbackRes = await payCallback(callbackData)
      console.log('[订单支付] 支付回调响应:', callbackRes)

      // 更新订单状态为已支付
      console.log('[订单支付] 更新订单状态:', { orderId: order.id })
      const updateRes = await updateOrderStatus(order.id)
      console.log('[订单支付] 更新订单状态响应:', updateRes)

      if (!updateRes || updateRes.code !== 0) {
        throw new Error(updateRes?.msg || '更新订单状态失败')
      }

      // 重新加载订单详情以获取最新状态
      await this.loadOrderDetail(order.id)

      // 显示支付成功提示
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