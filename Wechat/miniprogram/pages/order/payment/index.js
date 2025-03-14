import { getOrderDetail, payOrder, payCallback, updateOrderStatus } from '../../../api/order'
import { getProductDetail } from '../../../api/product'
import { getFilmOrderDetail, updateFilmOrderStatus } from '../../../api/film'

Page({
  data: {
    orderId: '',
    order: null,
    payId: null,
    isFilmOrder: false,
    paySuccess: false,
    isLoading: true,
    error: null,
    paymentMethods: [
      { id: 1, name: '微信支付', icon: 'wechat', selected: true },
      { id: 2, name: '支付宝', icon: 'alipay', selected: false },
      { id: 3, name: '银行卡', icon: 'creditcard', selected: false }
    ],
    paymentInProgress: false,
    orderStatus: {
      0: '待支付',
      1: '已支付',
      2: '待收货',
      9: '已取消'
    },
    payStatus: {
      0: '未支付',
      1: '已支付',
      2: '支付失败',
      9: '已取消'
    },
    countdown: '', // 倒计时显示
    countdownTimer: null, // 倒计时定时器ID
    orderExpireTime: null, // 订单过期时间
    placeholderHeight: '160rpx' // 占位符高度
  },

  onLoad(options) {
    console.log('[订单支付] 页面加载参数:', options)
    const { orderId, payId, type } = options
    
    // 设置默认的倒计时时间（从当前时间开始15分钟倒计时）
    this.setDefaultCountdown()
    
    // 立即启动倒计时（不依赖订单数据）
    const timer = setInterval(() => {
      this.updateCountdown()
    }, 1000)
    this.setData({ countdownTimer: timer })
    
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

  onShow() {
    // 如果有倒计时，重新启动
    if (this.data.orderExpireTime) {
      this.startCountdown()
    }
  },

  onHide() {
    // 页面隐藏时清除倒计时
    this.clearCountdown()
  },

  onUnload() {
    // 页面卸载时清除倒计时
    this.clearCountdown()
  },

  // 启动倒计时
  startCountdown() {
    // 清除可能存在的定时器
    this.clearCountdown()
    
    try {
      // 获取订单创建时间，如果无法获取则使用当前时间
      let createTime = new Date()
      const order = this.data.order
      
      if (order && order.create_time) {
        const timeStr = order.create_time
        // 尝试直接解析日期
        const date = new Date(timeStr)
        if (!isNaN(date.getTime())) {
          createTime = date
        } else {
          // 尝试替换格式再解析
          const altDate = new Date(timeStr.replace(/-/g, '/'))
          if (!isNaN(altDate.getTime())) {
            createTime = altDate
          }
          // 如果仍然无效，使用当前时间（已赋值为默认值）
        }
      }
      
      // 计算过期时间（创建时间 + 15分钟）
      const expireTime = new Date(createTime.getTime() + 15 * 60 * 1000)
      
      // 更新数据
      this.setData({ orderExpireTime: expireTime })
      
      // 立即更新一次倒计时
      this.updateCountdown()
      
      // 设置定时器，每秒更新倒计时
      const timer = setInterval(() => {
        this.updateCountdown()
      }, 1000)
      
      this.setData({ countdownTimer: timer })
      
    } catch (error) {
      console.error('[订单支付] 初始化倒计时发生错误:', error)
      // 出错时设置默认倒计时
      this.setDefaultCountdown()
    }
  },
  
  // 设置默认倒计时（15分钟）
  setDefaultCountdown() {
    // 计算15分钟后的时间
    const expireTime = new Date(new Date().getTime() + 15 * 60 * 1000)
    
    // 设置数据
    this.setData({ 
      orderExpireTime: expireTime,
      countdown: '15:00'
    })
    
    // 调整占位区域高度，确保页面内容不被遮挡
    const query = wx.createSelectorQuery()
    query.select('.countdown-container').boundingClientRect()
    query.exec(res => {
      if (res && res[0]) {
        const height = res[0].height
        // 设置占位符高度，增加一些额外空间
        this.setData({
          placeholderHeight: height + 20 + 'px'
        })
      } else {
        // 默认高度
        this.setData({
          placeholderHeight: '160rpx'
        })
      }
    })
  },
  
  // 更新倒计时显示
  updateCountdown() {
    try {
      const now = new Date()
      const expireTime = this.data.orderExpireTime
      
      if (!expireTime) {
        this.setDefaultCountdown()
        return
      }
      
      // 计算剩余时间（毫秒）
      let remainTime = expireTime.getTime() - now.getTime()
      
      if (remainTime <= 0) {
        // 订单已过期
        this.clearCountdown()
        this.setData({ countdown: '00:00' })
        
        // 更新订单状态为已取消（状态码9）
        const order = this.data.order
        if (order && order.status === 0) {
          // 先更新UI显示
          this.setData({
            'order.status': 9,
            'order.statusText': '已取消',
            'order.pay_status': 9,
            'order.payStatusText': '已取消'
          })
          
          // 提示用户订单已过期
          wx.showToast({
            title: '订单支付超时已自动取消',
            icon: 'none',
            duration: 2000
          })
          
          // 调用API更新数据库中的订单状态
          if (this.data.isFilmOrder) {
            // 胶片订单取消
            const { updateFilmOrderStatus } = require('../../../api/film')
            updateFilmOrderStatus(order.id, 9)
              .then(res => {
                console.log('[订单支付] 胶片订单已自动取消:', res)
              })
              .catch(err => {
                console.error('[订单支付] 胶片订单取消失败:', err)
              })
          } else {
            // 普通商品订单取消
            const { updateOrderStatus } = require('../../../api/order')
            updateOrderStatus(order.id, 9)
              .then(res => {
                console.log('[订单支付] 订单已自动取消:', res)
              })
              .catch(err => {
                console.error('[订单支付] 订单取消失败:', err)
              })
            
            // 如果存在支付记录，更新支付状态为已取消
            if (this.data.payId) {
              const { payCallback } = require('../../../api/order')
              payCallback({
                id: parseInt(this.data.payId),
                uid: parseInt(order.uid) || 0,
                oid: parseInt(order.id),
                amount: parseInt(parseFloat(order.amount || 0) * 100),
                source: 0,
                status: 9  // 设置为已取消状态
              })
                .then(res => {
                  console.log('[订单支付] 支付记录已更新为已取消:', res)
                })
                .catch(err => {
                  console.error('[订单支付] 支付记录更新失败:', err)
                })
            }
          }
        }
        
        return
      }
      
      // 计算分钟和秒数
      const minutes = Math.floor(remainTime / (60 * 1000))
      remainTime %= (60 * 1000)
      const seconds = Math.floor(remainTime / 1000)
      
      // 格式化为 MM:SS
      const countdown = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
      
      this.setData({ countdown })
    } catch (error) {
      console.error('[订单支付] 更新倒计时出错:', error)
      this.setData({ countdown: '15:00' })
    }
  },
  
  // 清除倒计时
  clearCountdown() {
    if (this.data.countdownTimer) {
      clearInterval(this.data.countdownTimer)
      this.setData({ countdownTimer: null })
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
      console.log('[订单支付] 原始订单数据:', JSON.stringify(orderData))
      
      // 检查原始订单中的商品项
      if (orderData.items && orderData.items.length > 0) {
        console.log('[订单支付] 原始订单中的商品项:', orderData.items.map(item => ({
          id: item.id,
          pid: item.pid,
          name: item.name,
          product_name: item.product_name
        })))
      }
      
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

      // 构建商品信息数组
      let items = []
      
      // 先检查订单中是否已经包含商品信息
      if (orderData.items && orderData.items.length > 0) {
        // 如果订单中已经包含商品信息，直接使用
        console.log('[订单支付] 使用订单中现有的商品信息')
        items = orderData.items.map(item => {
          const processedItem = {
            id: item.id || 0,
            pid: item.pid || 0,
            name: item.product_name || item.name || '未知商品',
            product_name: item.product_name || '',
            price: (parseInt(item.price) / 100).toFixed(2),
            quantity: parseInt(item.quantity) || 1,
            total_price: (parseInt(item.total_price || item.amount) / 100).toFixed(2),
            cover_image: item.product_image || item.cover_image || '/assets/images/film-icon.png',
            product_image: item.product_image || item.cover_image || '/assets/images/film-icon.png'
          }
          
          console.log('[订单支付] 处理后的商品项:', processedItem)
          return processedItem
        })
      } 
      // 获取商品详情（如果存在商品ID且没有商品列表）
      else if (orderData.pid) {
        try {
          console.log('[订单支付] 开始获取商品详情:', { 商品ID: orderData.pid })
          
          const productRes = await getProductDetail(orderData.pid)
          console.log('[订单支付] 获取商品详情成功:', productRes)
          
          if (productRes && productRes.data) {
            const product = productRes.data
            
            // 处理商品图片路径
            let coverImage = '/assets/images/film-icon.png'
            if (product.product_image) {
              coverImage = product.product_image.startsWith('http') 
                ? product.product_image 
                : `http://localhost:8001${product.product_image}`
            }
            
            // 计算单价（总金额除以数量）
            const unitPrice = (parseFloat(orderData.amount) / orderData.quantity).toFixed(2)
            
            // 构造商品信息
            const productItem = {
              id: orderData.pid,
              name: product.name || '未知商品',
              product_name: product.name || '',
              price: unitPrice,
              quantity: orderData.quantity,
              cover_image: coverImage,
              product_image: coverImage,
              total_price: orderData.amount
            }
            
            console.log('[订单支付] 从商品详情创建的商品项:', productItem)
            items.push(productItem)
          }
        } catch (error) {
          console.error('[订单支付] 获取商品详情失败:', error)
          // 使用默认商品信息
          items.push({
            id: orderData.pid || 0,
            name: '未知商品',
            price: (parseFloat(orderData.amount) / orderData.quantity).toFixed(2),
            quantity: orderData.quantity,
            cover_image: '/assets/images/film-icon.png',
            product_image: '/assets/images/film-icon.png',
            total_price: orderData.amount
          })
        }
      } else {
        // 没有商品信息时使用默认值
        console.log('[订单支付] 使用默认商品信息')
        items.push({
          id: 0,
          name: '未知商品',
          price: orderData.amount,
          quantity: orderData.quantity || 1,
          cover_image: '/assets/images/film-icon.png',
          product_image: '/assets/images/film-icon.png',
          total_price: orderData.amount
        })
      }
      
      // 确保items数组存在并且有至少一个元素
      if (!items || items.length === 0) {
        items = [{
          id: 0,
          name: '未知商品',
          price: orderData.amount,
          quantity: 1,
          cover_image: '/assets/images/film-icon.png',
          product_image: '/assets/images/film-icon.png',
          total_price: orderData.amount
        }]
      }
      
      // 将处理后的items赋值给orderData
      orderData.items = items
      
      console.log('[订单支付] 处理后的商品信息:', items)
      
      // 最后再次检查确保每个商品都有name字段
      if (orderData.items && orderData.items.length > 0) {
        orderData.items = orderData.items.map(item => {
          if (!item.name) {
            if (item.product_name) {
              item.name = item.product_name;
              console.log(`[订单支付] 为商品ID ${item.id} 设置名称:`, item.name);
            } else {
              item.name = '未知商品';
              console.log(`[订单支付] 为商品ID ${item.id} 设置默认名称:`, item.name);
            }
          }
          return item;
        });
      }
      
      console.log('[订单支付] 最终处理后的订单数据:', JSON.stringify(orderData.items))

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

      // 启动倒计时
      this.startCountdown()
      
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
      
      // 构建商品信息数组
      let items = []
      
      // 处理订单项，确保字段存在
      if (orderData.items && orderData.items.length > 0) {
        console.log('[胶片订单支付] 处理订单项，原始数据:', JSON.stringify(orderData.items))
        
        items = orderData.items.map(item => {
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
        
        console.log('[胶片订单支付] 处理后的订单项:', items)
      } else {
        // 如果没有订单项，创建一个默认项
        items.push({
          id: 0,
          name: '未知商品',
          price: orderData.total_price,
          quantity: 1,
          amount: orderData.total_price,
          product_image: '/assets/images/film-icon.png'
        })
      }
      
      // 确保items数组存在并且有至少一个元素
      if (!items || items.length === 0) {
        items = [{
          id: 0,
          name: '未知商品',
          price: orderData.total_price,
          quantity: 1,
          amount: orderData.total_price, 
          product_image: '/assets/images/film-icon.png'
        }]
      }
      
      // 将处理后的items赋值给orderData
      orderData.items = items

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

      // 启动倒计时
      this.startCountdown()
      
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

  // 支付处理
  onPayment() {
    // 检查是否已经在处理支付
    if (this.data.paymentInProgress) {
      console.log('[订单支付] 支付已在进行中，避免重复提交')
      return
    }
    
    const { order, isFilmOrder } = this.data
    if (!order) {
      wx.showToast({
        title: '订单数据不存在',
        icon: 'none'
      })
      return
    }
    
    // 设置支付处理中状态
    this.setData({
      paymentInProgress: true
    })
    
    // 检查订单状态
    if (order.status !== 0) {
      wx.showToast({
        title: '订单状态已变更，无法支付',
        icon: 'none'
      })
      
      // 刷新订单状态
      setTimeout(() => {
        // 重置支付状态
        this.setData({
          paymentInProgress: false
        })
        
        // 刷新订单数据
        if (isFilmOrder) {
          this.loadFilmOrderDetail(order.id)
        } else {
          this.loadOrderDetail(order.id)
        }
      }, 1500)
      return
    }
    
    console.log(`[订单支付] 开始支付处理，订单ID: ${order.id}`)
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
              'paySuccess': true,
              'paymentInProgress': false
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
            
            // 重置支付状态
            this.setData({
              paymentInProgress: false
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
          
          // 重置支付状态
          this.setData({
            paymentInProgress: false
          })
        })
    } else {
      // 普通商品订单 - 保持原有逻辑
      // 首先检查订单状态是否为待支付状态
      if (order.status !== 0) {
        wx.hideLoading()
        wx.showToast({
          title: '订单状态已更改，无法支付',
          icon: 'none'
        })
        
        // 刷新订单详情
        this.loadOrderDetail(order.id)
        return
      }

      // 先执行支付回调，确认支付成功后再更新订单状态
      payCallback({
        id: parseInt(this.data.payId),  // 确保ID是整数
        uid: parseInt(order.uid),       // 确保用户ID是整数
        oid: parseInt(order.id),        // 确保订单ID是整数
        amount: parseInt(parseFloat(order.amount) * 100),  // 转换元为分，确保是整数
        source: 0,  // 0表示微信支付
        status: 1   // 1表示已支付
      })
        .then(callbackRes => {
          console.log('[订单支付] 支付回调成功:', callbackRes)
          
          // 支付回调成功后，更新订单状态
          return updateOrderStatus(order.id)
            .then(res => {
              wx.hideLoading()
              if (res.code === 0) {
                wx.showToast({
                  title: '支付成功',
                  icon: 'success'
                })
                this.setData({
                  'paySuccess': true,
                  'paymentInProgress': false
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
              } else {
                console.error('[订单支付] 更新订单状态失败:', res)
                wx.showToast({
                  title: res.msg || '支付状态更新失败',
                  icon: 'none'
                })
                
                // 重置支付状态
                this.setData({
                  paymentInProgress: false
                })
              }
            })
            .catch(updateErr => {
              wx.hideLoading()
              console.error('[订单支付] 更新订单状态异常:', updateErr)
              wx.showToast({
                title: '支付成功但状态更新失败',
                icon: 'none'
              })
              
              // 重置支付状态
              this.setData({
                paymentInProgress: false
              })
              
              // 刷新订单详情
              this.loadOrderDetail(order.id)
            })
        })
        .catch(callbackErr => {
          wx.hideLoading()
          console.error('[订单支付] 支付回调异常:', callbackErr)
          // 处理支付回调失败
          wx.showToast({
            title: '支付处理失败',
            icon: 'none'
          })
          
          // 重置支付状态
          this.setData({
            paymentInProgress: false
          })
          
          // 刷新订单详情
          this.loadOrderDetail(order.id)
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