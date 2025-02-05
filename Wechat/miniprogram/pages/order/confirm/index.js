// pages/order/confirm/index.js
import { createOrder, payOrder } from '../../../api/order'
import { getAddressList } from '../../../api/address'
import { getProductDetail } from '../../../api/product'
import { loginGuard } from '../../../utils/auth'
import { removeFromCart } from '../../../api/cart'

// 添加Array.includes的polyfill
if (!Array.prototype.includes) {
  Array.prototype.includes = function(searchElement, fromIndex) {
    if (this == null) {
      throw new TypeError('"this" is null or undefined')
    }
    const o = Object(this)
    const len = o.length >>> 0
    if (len === 0) return false
    const n = fromIndex | 0
    let k = Math.max(n >= 0 ? n : len + n, 0)
    while (k < len) {
      if (o[k] === searchElement) return true
      k++
    }
    return false
  }
}

Page(loginGuard({
  data: {
    address: null,
    orderItems: [],
    totalPrice: 0,
    totalCount: 0,
    remark: '',
    loading: false,
    fromCart: false,  // 添加来源标记
    shippingFee: 0,   // 添加运费字段
  },

  onLoad(options) {
    console.log('[订单确认] 页面加载, 参数:', options)
    
    // 从购物车结算进入
    const { from } = options
    if (from === 'cart') {
      console.log('[订单确认] 从购物车进入')
      this.setData({ fromCart: true })
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

  // 从购物车加载商品
  loadCartProductsAndSetOrderItems() {
    try {
      // 获取已选中的购物车商品
      const selectedItems = wx.getStorageSync('selectedCartItems') || []
      console.log('[订单确认] 已选中的购物车商品:', selectedItems)
      
      if (!selectedItems.length) {
        console.error('[订单确认] 没有选中的购物车商品')
        return
      }
      
      // 计算总数量和商品总价
      const totalCount = selectedItems.reduce((total, item) => total + item.quantity, 0)
      const productTotal = selectedItems.reduce((total, item) => total + item.price * item.quantity, 0)
      
      // 计算运费
      const shippingFee = totalCount >= 3 ? 0 : 7
      const totalPrice = productTotal + shippingFee
      
      // 设置订单商品
      this.setData({
        orderItems: selectedItems.map(item => ({
          product_id: item.product_id,
          name: item.name,
          price: item.price,
          quantity: item.quantity,
          cover_image: item.product_image  // 使用从购物车传递的完整图片路径
        })),
        totalPrice: totalPrice.toFixed(2),
        totalCount,
        shippingFee: shippingFee.toFixed(2)
      })
      
      console.log('[订单确认] 设置订单商品成功:', this.data.orderItems)
    } catch (error) {
      console.error('[订单确认] 加载购物车商品失败:', error)
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

      return {
        ...item,
        price: price,
        quantity: quantity
      }
    }).filter(Boolean)

    // 计算总数量和商品总价
    const totalCount = orderItems.reduce((total, item) => total + item.quantity, 0)
    const productTotal = orderItems.reduce((total, item) => total + item.price * item.quantity, 0)
    
    // 计算运费
    const shippingFee = totalCount >= 3 ? 0 : 7
    const totalPrice = productTotal + shippingFee

    this.setData({
      orderItems,
      totalCount,
      totalPrice: totalPrice.toFixed(2),
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

  // 清除已购买的购物车商品
  async clearCartItems() {
    try {
      console.log('[订单确认] 开始清除购物车商品')
      
      // 获取已选中的购物车商品ID
      const selectedItems = wx.getStorageSync('selectedCartItems') || []
      console.log('[订单确认] 需要清除的商品:', selectedItems)
      
      if (selectedItems.length === 0) {
        console.log('[订单确认] 没有需要清除的商品')
        return
      }

      // 获取当前购物车数据，用于获取购物车商品的ID
      const cartList = wx.getStorageSync('cartList') || []
      console.log('[订单确认] 当前购物车数据:', cartList)

      // 找到需要删除的购物车商品ID
      const cartItemsToDelete = cartList.filter(cartItem => 
        selectedItems.some(selected => selected.product_id === cartItem.productId)
      )
      console.log('[订单确认] 需要删除的购物车商品:', cartItemsToDelete)

      // 逐个删除选中的购物车商品
      for (const item of cartItemsToDelete) {
        try {
          const res = await removeFromCart(item.id)
          console.log(`[订单确认] 删除购物车商品 ${item.id} 响应:`, res)
          
          if (!res || (res.code !== 0 && res.code !== 200)) {
            console.error(`[订单确认] 删除购物车商品 ${item.id} 失败:`, res)
          }
        } catch (err) {
          console.error(`[订单确认] 删除购物车商品 ${item.id} 出错:`, err)
        }
      }

      // 更新本地购物车数据
      const newCartList = cartList.filter(item => 
        !selectedItems.some(selected => selected.product_id === item.productId)
      )
      console.log('[订单确认] 更新后的购物车数据:', newCartList)
      
      // 更新购物车数据
      wx.setStorageSync('cartList', newCartList)
      
      // 清除已选商品缓存
      wx.removeStorageSync('selectedCartItems')
      
      // 尝试更新购物车页面
      const pages = getCurrentPages()
      const cartPage = pages.find(p => p.route === 'pages/cart/index')
      if (cartPage && typeof cartPage.loadCartList === 'function') {
        console.log('[订单确认] 刷新购物车页面')
        cartPage.loadCartList()
      }
      
      console.log('[订单确认] 清除购物车商品完成')
    } catch (error) {
      console.error('[订单确认] 清除购物车商品失败:', error)
      // 不显示错误提示，因为这不应该影响用户体验
      console.warn('[订单确认] 清除购物车失败，但不影响订单提交:', error.message)
    }
  },

  // 提交订单
  async submitOrder() {
    console.log('[订单确认] 开始提交订单')
    if (this.data.loading) {
      console.log('[订单确认] 订单正在提交中，请勿重复点击')
      return
    }

    if (!this.data.address) {
      wx.showToast({
        title: '请选择收货地址',
        icon: 'none'
      })
      return
    }

    if (!this.data.orderItems.length) {
      wx.showToast({
        title: '订单商品数据无效',
        icon: 'none'
      })
      return
    }

    this.setData({ loading: true })
    
    try {
      // 获取用户信息
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo || !userInfo.id) {
        throw new Error('用户信息无效，请重新登录')
      }

      // 计算总价（转换为分）
      const totalPriceInCents = Math.round(parseFloat(this.data.totalPrice) * 100)
      
      const orderData = {
        uid: userInfo.id,
        address_id: this.data.address.id,
        total_price: totalPriceInCents,
        shipping_fee: 700,
        remark: this.data.remark || '',
        status: 0, // 新订单状态为待支付
        items: this.data.orderItems.map(item => {
          const priceInCents = Math.round(item.price * 100)
          return {
            pid: item.product_id,
            product_name: item.name,
            product_image: item.cover_image,
            price: priceInCents,
            quantity: item.quantity,
            amount: priceInCents * item.quantity
          }
        })
      }

      console.log('[订单确认] 提交订单数据:', orderData)
      
      const res = await createOrder(orderData)
      console.log('[订单确认] 创建订单响应:', res)

      if (!res || res.code !== 0 || !res.data || !res.data.id) {
        throw new Error(res?.msg || '创建订单失败')
      }

      // 创建支付记录（状态为未支付）
      const payData = {
        oid: res.data.id,
        uid: userInfo.id,
        amount: totalPriceInCents
      }
      console.log('[订单确认] 创建支付记录:', payData)
      const payRes = await payOrder(payData)
      console.log('[订单确认] 创建支付响应:', payRes)

      if (!payRes || (payRes.code !== 0 && payRes.code !== 200 && payRes.msg !== 'success')) {
        throw new Error(payRes?.msg || '创建支付记录失败')
      }

      // 保存支付记录ID到本地存储
      const payId = payRes.data.id
      wx.setStorageSync('currentPayId', payId)

      // 清除购物车中已购买的商品
      if (this.data.fromCart) {
        await this.clearCartItems()
      }

      // 直接跳转到支付页面
      wx.redirectTo({
        url: `/pages/order/payment/index?orderId=${res.data.id}&payId=${payId}`
      })

    } catch (error) {
      console.error('[订单确认] 提交订单失败:', error)
      wx.showToast({
        title: error.message || '创建订单失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 获取商品详情
  async getProductDetail(cartItem) {
    try {
      const res = await getProductDetail(cartItem.product_id)
      console.log('[订单确认] 商品详情响应:', res)
      
      if (res && (res.code === 0 || res.code === 200)) {
        const productData = res.data || res
        return {
          ...cartItem,
          stock: productData.stock,
          amount: productData.amount,
          image: productData.cover_image || cartItem.cover_image
        }
      } else {
        console.error('[订单确认] 商品详情获取失败:', res)
        // 如果获取失败，返回购物车中的数据
        return cartItem
      }
    } catch (error) {
      console.error('[订单确认] 商品', cartItem.product_id, '详情获取出错:', error)
      // 如果出错，返回购物车中的数据
      console.log('[订单确认] 使用购物车数据作为回退:', cartItem)
      return cartItem
    }
  }
}))