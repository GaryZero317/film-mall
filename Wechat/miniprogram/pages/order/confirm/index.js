// pages/order/confirm/index.js
import { createOrder } from '../../../api/order'
import { getAddressList } from '../../../api/address'
import { getProductDetail } from '../../../api/product'
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
    console.log('[订单确认] 页面加载, 参数:', options)
    
    // 从购物车结算进入
    const { from } = options
    if (from === 'cart') {
      console.log('[订单确认] 从购物车进入')
      this.loadCartProductsAndSetOrderItems()
    } 
    // 从商品详情页直接购买进入
    else {
      const { productId, quantity } = options
      console.log('[订单确认] 从商品详情页进入, 参数:', { productId, quantity })
      if (productId && quantity) {
        this.loadProductAndSetOrderItems(productId, quantity)
      } else {
        console.error('[订单确认] 商品参数无效:', options)
      }
    }
    this.loadAddress()
  },

  // 加载商品详情并设置订单商品
  async loadProductAndSetOrderItems(productId, quantity) {
    console.log('[订单确认] 开始加载商品详情:', { productId, quantity })
    try {
      const res = await getProductDetail(productId)
      console.log('[订单确认] 商品详情响应:', res)
      
      if (res && res.code === 0 && res.data) {
        const product = res.data
        console.log('[订单确认] 商品详情数据:', {
          id: product.id,
          name: product.name,
          price: product.price,
          sellPrice: product.sellPrice,
          selling_price: product.selling_price,
          原始数据: product
        })

        // 确保价格是有效数字
        let price = 0
        // 尝试从不同可能的价格字段中获取价格
        const possiblePriceFields = ['price', 'sellPrice', 'salePrice', 'retailPrice', 'marketPrice', 'amount']
        for (const field of possiblePriceFields) {
          if (product[field] !== undefined && product[field] !== null) {
            const rawPrice = product[field]
            console.log(`[订单确认] 检查价格字段 ${field}:`, {
              值: rawPrice,
              类型: typeof rawPrice
            })
            
            let parsedPrice
            if (field === 'amount') {
              // amount字段单位是分，需要转换为元
              parsedPrice = typeof rawPrice === 'number' ? 
                rawPrice / 100 : 
                parseFloat(String(rawPrice)) / 100
            } else {
              parsedPrice = typeof rawPrice === 'number' ? 
                rawPrice : 
                parseFloat(String(rawPrice))
            }
            
            if (!isNaN(parsedPrice) && parsedPrice > 0) {
              price = parsedPrice
              console.log(`[订单确认] 使用价格字段 ${field}:`, price)
              break
            }
          }
        }

        if (price === 0) {
          console.error('[订单确认] 无法获取有效价格，商品数据:', product)
          wx.showToast({
            title: '商品价格数据无效',
            icon: 'none'
          })
          return
        }

        console.log('[订单确认] 最终确定的价格:', price)

        // 处理图片路径
        let coverImage = '/assets/images/default.png'
        const possibleImageFields = ['coverImage', 'mainImage', 'image', 'cover_image', 'main_image', 'img', 'imgUrl', 'imageUrl', 'cover']
        for (const field of possibleImageFields) {
          if (product[field]) {
            const imagePath = product[field]
            console.log(`[订单确认] 检查图片字段 ${field}:`, imagePath)
            coverImage = imagePath.startsWith('http') 
              ? imagePath 
              : `http://localhost:8001${imagePath}`
            console.log(`[订单确认] 使用图片路径:`, coverImage)
            break
          }
        }

        const orderItem = {
          id: product.id,
          product_id: parseInt(productId),
          name: product.name || product.productName || product.product_name || '未知商品',
          price: price,
          quantity: parseInt(quantity),
          cover_image: coverImage
        }
        console.log('[订单确认] 构造的订单商品:', orderItem)
        this.setOrderItems([orderItem])
      } else {
        console.error('[订单确认] 获取商品详情失败:', {
          响应数据: res,
          错误原因: !res ? '响应为空' : !res.code ? '无响应码' : '无商品数据'
        })
        wx.showToast({
          title: '获取商品信息失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[订单确认] 加载商品详情失败:', {
        错误信息: error.message,
        错误堆栈: error.stack
      })
      wx.showToast({
        title: '获取商品信息失败',
        icon: 'none'
      })
    }
  },

  // 加载购物车商品详情并设置订单商品
  async loadCartProductsAndSetOrderItems() {
    console.log('[订单确认] 开始加载购物车商品')
    try {
      const cartProducts = wx.getStorageSync('selectedCartItems') || []
      console.log('[订单确认] 购物车商品列表:', cartProducts)
      
      if (!cartProducts.length) {
        wx.showToast({
          title: '请先选择商品',
          icon: 'none'
        })
        setTimeout(() => {
          wx.navigateBack()
        }, 1500)
        return
      }

      // 获取每个商品的详细信息
      const orderItemsPromises = cartProducts.map(item => this.getProductDetail(item))
      const orderItems = await Promise.all(orderItemsPromises)
      
      console.log('[订单确认] 订单商品列表:', orderItems)
      
      // 计算总价
      const totalPrice = orderItems.reduce((total, item) => {
        return total + (item.price * item.quantity)
      }, 0)
      
      this.setData({
        orderItems,
        totalPrice: totalPrice.toFixed(2)
      })
    } catch (error) {
      console.error('[订单确认] 加载商品失败:', error)
      wx.showToast({
        title: '加载商品失败',
        icon: 'none'
      })
    }
  },

  // 加载收货地址
  async loadAddress() {
    console.log('[订单确认] 开始加载地址列表')
    try {
      const res = await getAddressList()
      console.log('[订单确认] 获取地址列表响应:', res)
      
      if (res && res.data) {
        // 获取默认地址
        const defaultAddress = res.data.list.find(item => item.isDefault) || res.data.list[0]
        console.log('[订单确认] 选择的默认地址:', defaultAddress)
        this.setData({ address: defaultAddress })
      } else {
        console.error('[订单确认] 获取地址列表失败, 响应数据无效:', res)
      }
    } catch (error) {
      console.error('[订单确认] 加载地址失败, 错误:', error)
    }
  },

  // 设置订单商品
  setOrderItems(items) {
    console.log('[订单确认] 开始设置订单商品, 原始数据:', items)
    
    if (!Array.isArray(items) || items.length === 0) {
      console.error('[订单确认] 商品数据无效:', items)
      wx.showToast({
        title: '商品数据无效',
        icon: 'none'
      })
      return
    }

    const orderItems = items.map(item => {
      // 确保价格是有效数字，如果是字符串则转换为数字
      const price = typeof item.price === 'string' ? parseFloat(item.price) : item.price
      const quantity = parseInt(item.quantity || 1)
      
      if (isNaN(price) || price <= 0) {
        console.error('[订单确认] 商品价格无效:', item)
        wx.showToast({
          title: '商品价格数据无效',
          icon: 'none'
        })
        return null
      }

      const processedItem = {
        id: item.id,  // 原始ID
        product_id: parseInt(item.product_id || item.id),  // 确保product_id是数字类型
        name: item.name || '未知商品',
        price: price,
        quantity: quantity,
        cover_image: item.cover_image || '/assets/images/default.png'
      }
      console.log(`[订单确认] 处理商品 ${processedItem.name}:`, {
        原始数据: item,
        处理后数据: processedItem,
        product_id类型: typeof processedItem.product_id,
        product_id值: processedItem.product_id,
        价格类型: typeof processedItem.price,
        价格值: processedItem.price
      })
      return processedItem
    }).filter(item => item !== null)  // 过滤掉无效的商品

    if (orderItems.length === 0) {
      console.error('[订单确认] 没有有效的商品数据')
      wx.showToast({
        title: '商品数据无效',
        icon: 'none'
      })
      return
    }

    const totalPrice = orderItems.reduce((sum, item) => {
      const itemTotal = item.price * item.quantity
      console.log(`[订单确认] 计算商品 ${item.name} 小计: ${itemTotal} = ${item.price} × ${item.quantity}`)
      return sum + itemTotal
    }, 0)
    
    const totalCount = orderItems.reduce((sum, item) => {
      console.log(`[订单确认] 计算商品 ${item.name} 数量: ${item.quantity}`)
      return sum + item.quantity
    }, 0)

    // 计算运费：3件及以上免运费，否则7元运费
    const shippingFee = totalCount >= 3 ? 0 : 7
    console.log('[订单确认] 计算运费:', {
      商品总数: totalCount,
      是否免运费: totalCount >= 3,
      运费金额: shippingFee
    })

    // 计算订单总价（商品总价 + 运费）
    const finalTotalPrice = totalPrice + shippingFee

    console.log('[订单确认] 订单数据汇总:', {
      商品列表: orderItems,
      商品总价: totalPrice,
      运费: shippingFee,
      订单总价: finalTotalPrice,
      总数量: totalCount
    })

    this.setData({
      orderItems,
      totalPrice: finalTotalPrice.toFixed(2),
      totalCount,
      shippingFee: shippingFee.toFixed(2)
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
      
      // 获取用户ID
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo || !userInfo.id) {
        throw new Error('用户未登录')
      }

      // 构造订单数据
      const orderData = {
        uid: userInfo.id,
        address_id: address.id,
        pid: orderItems[0].product_id,  // 使用product_id作为pid
        quantity: orderItems[0].quantity,
        amount: Math.round(orderItems[0].price * 100),  // 转换为分为单位
        remark,
        total_price: Math.round(parseFloat(this.data.totalPrice) * 100),  // 转换为分为单位
        shipping_fee: Math.round(parseFloat(this.data.shippingFee) * 100),  // 转换为分为单位
        status: 1  // 1: 待支付
      }

      console.log('[订单确认] 提交订单数据:', {
        ...orderData,
        pid_type: typeof orderItems[0].product_id,
        original_id: orderItems[0].id,
        amount_yuan: orderItems[0].price,
        total_price_yuan: parseFloat(this.data.totalPrice),
        shipping_fee_yuan: parseFloat(this.data.shippingFee),
        status_desc: '待支付'
      })

      // 确保有pid和amount
      if (!orderData.pid) {
        throw new Error('商品ID无效')
      }
      if (!orderData.amount || orderData.amount <= 0) {
        throw new Error('商品金额无效')
      }

      const res = await createOrder(orderData)

      if (!res || !res.data || !res.data.id) {
        throw new Error('创建订单失败：返回数据无效')
      }

      // 下单成功后清除购物车中已购买的商品
      const cartItems = wx.getStorageSync('cartItems') || []
      const newCartItems = cartItems.filter(item => 
        !orderItems.find(orderItem => orderItem.id === item.product_id)
      )
      wx.setStorageSync('cartItems', newCartItems)

      // 确保订单ID是数字类型
      const orderId = parseInt(res.data.id)
      if (isNaN(orderId)) {
        throw new Error('创建订单失败：无效的订单ID')
      }

      console.log('[订单确认] 创建订单成功:', {
        订单ID: orderId,
        原始数据: res.data
      })

      // 跳转到支付页面
      wx.navigateTo({
        url: `/pages/order/payment/index?orderId=${orderId}`
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
  },

  async getProductDetail(cartItem) {
    console.log('[订单确认] 开始获取商品详情:', cartItem)
    try {
      // 使用product_id而不是productId
      const productId = cartItem.product_id || cartItem.productId
      if (!productId) {
        console.error('[订单确认] 商品ID无效:', cartItem)
        throw new Error('无效的商品ID')
      }

      const res = await getProductDetail(productId)
      console.log('[订单确认] 商品详情获取结果:', res)
      
      // 如果商品存在且有数据，使用API返回的数据
      if (res && res.code === 0 && res.data) {
        return {
          id: res.data.id,
          product_id: productId,
          name: res.data.name || res.data.productName || cartItem.name || cartItem.productName || '未知商品',
          price: res.data.price || cartItem.price,
          quantity: cartItem.quantity,
          cover_image: res.data.mainImage || res.data.coverImage || cartItem.cover_image || '/assets/images/default.png'
        }
      }
      
      // 如果商品不存在，使用购物车数据
      if (res && res.notFound) {
        console.log('[订单确认] 商品不存在，使用购物车数据:', cartItem)
        return {
          id: cartItem.id || productId,
          product_id: productId,
          name: cartItem.name || cartItem.productName || '未知商品',
          price: cartItem.price,
          quantity: cartItem.quantity,
          cover_image: cartItem.cover_image || cartItem.mainImage || '/assets/images/default.png'
        }
      }

      // 其他错误情况
      throw new Error(res.msg || '获取商品详情失败')
    } catch (error) {
      console.error(`[订单确认] 商品 ${cartItem.product_id || cartItem.productId} 详情获取出错:`, error)
      
      // 如果是网络错误或其他错误，尝试使用购物车数据
      if (cartItem.price && cartItem.quantity) {
        console.log('[订单确认] 使用购物车数据作为回退:', cartItem)
        return {
          id: cartItem.id || productId,
          product_id: cartItem.product_id || cartItem.productId,
          name: cartItem.name || cartItem.productName || '未知商品',
          price: cartItem.price,
          quantity: cartItem.quantity,
          cover_image: cartItem.cover_image || cartItem.mainImage || '/assets/images/default.png'
        }
      }

      wx.showToast({
        title: `商品信息获取失败`,
        icon: 'none'
      })
      throw error
    }
  }
}))