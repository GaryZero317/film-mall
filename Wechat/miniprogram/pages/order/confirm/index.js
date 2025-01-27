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
      const cartItems = wx.getStorageSync('selectedCartItems') || []
      console.log('[订单确认] 从购物车进入, 商品数据:', cartItems)
      this.setOrderItems(cartItems)
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
      const price = parseFloat(item.price || 0)
      const quantity = parseInt(item.quantity || 1)
      const processedItem = {
        id: item.product_id,
        name: item.name || '未知商品',
        price: price,
        quantity: quantity,
        cover_image: item.cover_image || '/assets/images/default.png'
      }
      console.log(`[订单确认] 处理商品 ${processedItem.name}:`, {
        原始数据: item,
        处理后数据: processedItem
      })
      return processedItem
    })

    const totalPrice = orderItems.reduce((sum, item) => {
      const itemTotal = item.price * item.quantity
      console.log(`[订单确认] 计算商品 ${item.name} 小计: ${itemTotal} = ${item.price} × ${item.quantity}`)
      return sum + itemTotal
    }, 0)
    
    const totalCount = orderItems.reduce((sum, item) => {
      console.log(`[订单确认] 计算商品 ${item.name} 数量: ${item.quantity}`)
      return sum + item.quantity
    }, 0)

    console.log('[订单确认] 订单数据汇总:', {
      商品列表: orderItems,
      总价: totalPrice,
      总数量: totalCount
    })

    this.setData({
      orderItems,
      totalPrice: totalPrice.toFixed(2),
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