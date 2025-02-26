import { getOrderDetail, payOrder, payCallback, updateOrderStatus } from '../../../api/order'
import { getProductDetail } from '../../../api/product'
import { getFilmOrderDetail, updateFilmOrderStatus } from '../../../api/film'

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
    },
    isFilmOrder: false // 标记是否为胶片订单
  },

  onLoad(options) {
    console.log('[订单支付] 页面加载参数:', options)
    const { orderId, payId, type } = options
    
    // 设置订单类型
    const isFilmOrder = type === 'film'
    this.setData({ isFilmOrder })
    
    if (!orderId) {
      console.error('[订单支付] 订单ID未提供:', options)
      wx.showToast({
        title: '订单ID未提供',
        icon: 'none'
      })
      
      // 添加延迟返回
      setTimeout(() => {
        wx.navigateBack()
      }, 2000)
      return
    }
    
    // 确保orderId是数字类型
    const id = parseInt(orderId)
    if (isNaN(id)) {
      console.error('[订单支付] 订单ID无效:', orderId)
      wx.showToast({
        title: '订单ID无效',
        icon: 'none'
      })
      
      // 添加延迟返回
      setTimeout(() => {
        wx.navigateBack()
      }, 2000)
      return
    }
    
    console.log('[订单支付] 开始加载订单详情:', { 
      原始ID: orderId, 
      转换后ID: id, 
      支付ID: payId, 
      类型: type,
      是否胶片订单: isFilmOrder
    })
    
    // 保存支付ID到data中
    this.setData({ payId: payId })
    
    // 根据类型加载不同的订单详情
    if (isFilmOrder) {
      this.loadFilmOrderDetail(id)
    } else {
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

  // 添加胶片订单详情加载方法
  async loadFilmOrderDetail(orderId) {
    try {
      this.setData({ loading: true })
      console.log('[胶片订单支付] 请求订单详情:', { 订单ID: orderId })
      const res = await getFilmOrderDetail(orderId)
      console.log('[胶片订单支付] 获取订单详情成功，原始数据:', JSON.stringify(res))
      console.log('[胶片订单支付] 响应数据结构:', {
        状态码: res?.code,
        消息: res?.message || res?.msg,
        数据字段: res?.data ? Object.keys(res.data) : '无数据'
      })
      
      if (!res) {
        throw new Error('订单详情请求失败')
      }
      
      if (res.code !== 0) {
        throw new Error(res.message || res.msg || `请求错误(${res.code})`)
      }
      
      if (!res.data) {
        throw new Error('未返回订单数据')
      }

      // 处理订单数据
      const orderData = res.data
      console.log('[胶片订单支付] 处理前的订单数据:', JSON.stringify(orderData))
      
      // 转换价格（分转元）
      // 确保字段存在且是数字，否则使用默认值
      const totalPrice = orderData.total_price || orderData.totalPrice || 0
      const shippingFee = orderData.shipping_fee || orderData.shippingFee || 0
      
      orderData.total_price = (parseInt(totalPrice) / 100).toFixed(2)
      orderData.shipping_fee = (parseInt(shippingFee) / 100).toFixed(2)
      
      // 检查是否回寄底片
      const isReturnFilm = Boolean(orderData.return_film)
      console.log('[胶片订单支付] 回寄底片状态:', isReturnFilm)
      orderData.return_film = isReturnFilm
      
      // 如果不回寄底片，则不收取运费
      const effectiveShippingFee = isReturnFilm ? parseInt(shippingFee) : 0
      console.log('[胶片订单支付] 实际计算的运费:', effectiveShippingFee)
      
      // 计算应付金额（只有回寄底片才加运费）
      const totalAmount = parseInt(totalPrice) + effectiveShippingFee
      orderData.amount = (totalAmount / 100).toFixed(2)
      
      // 如果不回寄底片，显示运费为0
      if (!isReturnFilm) {
        orderData.shipping_fee = "0.00"
      }
      
      // 处理订单项，确保字段存在
      if (orderData.items && orderData.items.length > 0) {
        console.log('[胶片订单支付] 处理订单项，原始数据:', JSON.stringify(orderData.items))
        
        orderData.items = orderData.items.map(item => {
          console.log('[胶片订单支付] 处理订单项:', item)
          
          // 获取价格，支持不同的字段名
          const price = item.price || item.unit_price || 0
          const amount = item.amount || (price * (item.quantity || 1)) || 0
          
          // 构建商品名称，确保各字段存在
          let typeName = item.film_type_name || ''
          let brandName = item.film_brand_name || ''
          const sizeName = item.size || item.film_size_name || '未知尺寸'
          
          // 如果缺少类型或品牌，尝试从item.name获取（如果存在）
          if ((!typeName || !brandName) && item.name) {
            const nameParts = item.name.split(' ')
            if (nameParts.length >= 2) {
              if (!typeName) typeName = nameParts[0]
              if (!brandName) brandName = nameParts[1]
            }
          }
          
          // 如果仍然缺少类型或品牌，设置默认值
          if (!typeName) typeName = '胶片'
          if (!brandName) brandName = '标准'
          
          // 构建完整的商品名称
          const fullName = `${typeName} ${brandName} ${sizeName}`
          
          return {
            ...item,
            price: (parseInt(price) / 100).toFixed(2),
            amount: (parseInt(amount) / 100).toFixed(2),
            name: fullName,
            product_image: '/assets/images/film-icon.png'
          }
        })
        
        console.log('[胶片订单支付] 处理后的订单项:', orderData.items)
      } else {
        // 如果没有订单项，创建一个默认项
        orderData.items = [{
          id: 0,
          name: '胶片冲洗服务',
          price: orderData.total_price,
          quantity: 1,
          amount: orderData.total_price,
          product_image: '/assets/images/film-icon.png'
        }]
      }

      // 设置订单号，支持多种可能的字段名
      orderData.order_no = orderData.order_no || orderData.foid || orderData.id || orderId
      
      // 设置订单ID，确保存在
      orderData.id = orderData.id || orderId
      
      // 设置订单状态，支持不同的字段名
      const status = parseInt(orderData.status || 0)
      orderData.status = status
      orderData.pay_status = status === 0 ? 0 : 1
      
      // 添加状态文本
      orderData.statusText = this.data.orderStatus[status] || '未知状态'
      orderData.payStatusText = this.data.payStatus[orderData.pay_status] || '未知支付状态'

      console.log('[胶片订单支付] 处理后的订单数据:', orderData)

      this.setData({ 
        order: orderData,
        loading: false
      })

      // 只有当支付状态不为待支付时才自动跳转
      if (orderData.status !== 0) {
        setTimeout(() => {
          wx.navigateBack()
        }, 3000)
      }
    } catch (error) {
      console.error('[胶片订单支付] 加载订单详情失败:', error)
      this.setData({ loading: false })
      wx.showToast({
        title: error.message || '加载订单失败',
        icon: 'none'
      })
      
      // 添加延迟返回
      setTimeout(() => {
        wx.navigateBack()
      }, 2000)
    }
  },

  // 确认支付
  onPayment() {
    const { order, isFilmOrder } = this.data
    if (!order) {
      wx.showToast({
        title: '订单信息不存在',
        icon: 'none'
      })
      return
    }

    console.log(`[订单支付] 开始支付订单，ID: ${order.id}, 订单号: ${order.order_no || order.foid}, 当前状态: ${order.status}, 状态描述: ${order.status_desc}`)
    console.log(`[订单支付] 回寄底片状态: ${order.return_film}`)
    
    // 模拟支付逻辑
    wx.showLoading({
      title: '支付处理中...',
      mask: true
    })
    
    // 胶片订单和商品订单使用不同的处理逻辑
    if (isFilmOrder) {
      // 胶片订单支付处理 - 使用updateFilmOrderStatus更新订单状态
      const { updateFilmOrderStatus } = require('../../../api/film')
      
      updateFilmOrderStatus(order.id)
        .then(res => {
          wx.hideLoading()
          console.log('[订单支付] 胶片订单状态更新响应:', res)
          
          if (res && res.code === 0) {
            // 更新成功
            wx.showToast({
              title: '支付成功，订单已更新为冲洗处理中',
              icon: 'none',
              duration: 2000
            })
            
            // 更新本地状态显示
            const updatedOrder = {...order, status: 1, status_desc: '冲洗处理中'}
            this.setData({
              order: updatedOrder,
              'paySuccess': true
            })
            
            // 发送支付结果事件
            const eventChannel = this.getOpenerEventChannel()
            if (eventChannel && eventChannel.emit) {
              eventChannel.emit('paymentResult', { success: true })
            }
          } else {
            // 更新失败
            console.error('[订单支付] 订单状态更新失败:', res)
            wx.showToast({
              title: res.msg || '支付失败',
              icon: 'none'
            })
          }
          
          // 延迟返回首页
          setTimeout(() => {
            wx.switchTab({
              url: '/pages/index/index'
            })
          }, 2000)
        })
        .catch(err => {
          wx.hideLoading()
          console.error('[订单支付] 支付异常:', err)
          wx.showToast({
            title: '支付过程中发生错误',
            icon: 'none'
          })
        })
    } else {
      // 普通商品订单 - 保持原有逻辑
      updateOrderStatus(order.id)
        .then(res => {
          wx.hideLoading()
          if (res.code === 0) {
            // 执行支付回调
            payCallback(this.data.payId)
              .then(callbackRes => {
                wx.showToast({
                  title: '支付成功',
                  icon: 'success'
                })
                this.setData({
                  'paySuccess': true
                })
                
                // 发送支付结果事件
                const eventChannel = this.getOpenerEventChannel()
                if (eventChannel && eventChannel.emit) {
                  eventChannel.emit('paymentResult', { success: true })
                }
                
                // 延迟返回首页
                setTimeout(() => {
                  wx.switchTab({
                    url: '/pages/index/index'
                  })
                }, 1500)
              })
              .catch(callbackErr => {
                console.error('[订单支付] 支付回调异常:', callbackErr)
                wx.showToast({
                  title: '支付异常',
                  icon: 'none'
                })
              })
          } else {
            console.error('[订单支付] 订单状态更新失败:', res)
            wx.showToast({
              title: res.msg || '订单支付失败',
              icon: 'none'
            })
          }
        })
        .catch(err => {
          wx.hideLoading()
          console.error('[订单支付] 支付异常:', err)
          wx.showToast({
            title: '支付过程中发生错误',
            icon: 'none'
          })
        })
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