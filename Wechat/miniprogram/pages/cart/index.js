// pages/cart/index.js
import { getCartList, updateCartItem, removeFromCart, updateCartItemStatus } from '../../api/cart'
import { loginGuard } from '../../utils/auth'

Page(loginGuard({
  data: {
    cartItems: [],
    loading: false,
    totalPrice: 0,
    selectedCount: 0,
    allSelected: false,
    isEmpty: true
  },

  onLoad() {
    this.loadCartList()
  },

  onShow() {
    this.loadCartList()
  },

  // 加载购物车列表
  async loadCartList() {
    try {
      this.setData({ loading: true })
      const res = await getCartList()
      console.log('[购物车] 获取购物车列表响应:', res)

      if (res && (res.code === 0 || res.code === 200)) {
        const cartList = res.data?.list || []
        console.log('[购物车] 原始购物车数据:', cartList)
        
        // 处理商品数据
        const formattedList = cartList.map(item => ({
            id: item.id,
            productId: item.productId,
            productName: item.productName,
            price: parseFloat(item.price || 0).toFixed(2),
            quantity: parseInt(item.quantity || 1),
            selected: Boolean(item.selected),
            stock: parseInt(item.stock || 999),
            // 处理图片URL，确保是完整的URL
            productImage: item.productImage ? 
              (item.productImage.startsWith('http') ? 
                item.productImage : 
                `http://localhost:8001${item.productImage}`
              ) : 
              'http://localhost:8001/uploads/placeholder.png'
        }))

        console.log('[购物车] 格式化后的购物车列表:', formattedList)
        
        // 更新本地存储
        wx.setStorageSync('cartList', formattedList)
        
        // 更新页面数据和计算总价
        this.setData({
          cartItems: formattedList,
          isEmpty: formattedList.length === 0
        })

        // 计算总价
        this.updateTotalPrice()
      } else {
        console.error('[购物车] 获取购物车列表失败:', res)
        throw new Error(res?.msg || '获取购物车列表失败')
      }
    } catch (error) {
      console.error('[购物车] 加载购物车列表失败:', error)
      wx.showToast({
        title: error.message || '加载购物车失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 更新总价和选中状态
  updateTotalPrice() {
    const { cartItems } = this.data
    console.log('[购物车] 开始计算总价，商品列表:', cartItems)
    
    const selectedItems = cartItems.filter(item => item.selected)
    const totalPrice = selectedItems.reduce((sum, item) => {
      return sum + (item.price * item.quantity)
    }, 0)
    
    const selectedCount = selectedItems.length
    const allSelected = selectedCount === cartItems.length && cartItems.length > 0
    
    console.log('[购物车] 计算结果:', {
      总价: totalPrice,
      选中数量: selectedCount,
      全选状态: allSelected
    })
    
    this.setData({
      totalPrice: totalPrice.toFixed(2),
      selectedCount,
      allSelected
    })
  },

  // 更新商品数量
  async onQuantityChange(e) {
    const { id, type } = e.currentTarget.dataset
    const item = this.data.cartItems.find(item => item.id === id)
    if (!item) return

    let newQuantity = item.quantity
    if (type === 'minus' && newQuantity > 1) {
      newQuantity--
    } else if (type === 'plus' && newQuantity < item.stock) {
      newQuantity++
    } else {
      return
    }

    try {
      console.log('[购物车] 更新商品数量:', { id, quantity: newQuantity })
      const res = await updateCartItem({
        id: parseInt(id),
        quantity: newQuantity
      })
      
      if (res && (res.code === 0 || res.code === 200)) {
        const cartItems = this.data.cartItems.map(item => {
          if (item.id === id) {
            return { ...item, quantity: newQuantity }
          }
          return item
        })
        
        this.setData({ cartItems })
        this.updateTotalPrice()
      } else {
        throw new Error(res?.msg || '更新数量失败')
      }
    } catch (error) {
      console.error('[购物车] 更新商品数量失败:', error)
      wx.showToast({
        title: error.message || '更新数量失败',
        icon: 'none'
      })
    }
  },

  // 删除商品
  async onDelete(e) {
    const { id } = e.currentTarget.dataset
    try {
      const res = await removeFromCart(id)
      if (res && (res.code === 0 || res.code === 200)) {
        const cartItems = this.data.cartItems.filter(item => item.id !== id)
        this.setData({ cartItems })
        this.updateTotalPrice()
        wx.showToast({
          title: '删除成功',
          icon: 'success'
        })
      } else {
        throw new Error(res?.msg || '删除失败')
      }
    } catch (error) {
      console.error('[购物车] 删除商品失败:', error)
      wx.showToast({
        title: error.message || '删除失败',
        icon: 'none'
      })
    }
  },

  // 选择商品
  async onItemSelect(e) {
    const { id } = e.currentTarget.dataset
    const cartItems = this.data.cartItems.map(item => {
      if (item.id === id) {
        return { ...item, selected: !item.selected }
      }
      return item
    })
    
    this.setData({ cartItems })
    this.updateTotalPrice()
    
    try {
      const item = cartItems.find(item => item.id === id)
      const res = await updateCartItemStatus({
        id: parseInt(id),
        selected: item.selected
      })
      
      if (res && (res.code !== 0 && res.code !== 200)) {
        throw new Error(res?.msg || '更新状态失败')
      }
    } catch (error) {
      console.error('[购物车] 更新商品选中状态失败:', error)
      wx.showToast({
        title: error.message || '更新状态失败',
        icon: 'none'
      })
    }
  },

  // 全选/取消全选
  onSelectAll() {
    const { allSelected, cartItems } = this.data
    const newCartItems = cartItems.map(item => ({
      ...item,
      selected: !allSelected
    }))
    
    this.setData({
      cartItems: newCartItems,
      allSelected: !allSelected
    })
    this.updateTotalPrice()
  },

  // 去结算
  onCheckout() {
    const selectedItems = this.data.cartItems.filter(item => item.selected)
    console.log('[购物车] 结算 - 选中的商品:', selectedItems)
    
    if (selectedItems.length === 0) {
      console.warn('[购物车] 结算 - 未选择商品')
      wx.showToast({
        title: '请选择商品',
        icon: 'none'
      })
      return
    }
    
    // 保存选中的商品到本地存储
    const checkoutItems = selectedItems.map(item => ({
        product_id: item.productId,
        name: item.productName,
        price: parseFloat(item.price),  // 转换为数字
        quantity: item.quantity,
      product_image: item.productImage  // 使用product_image作为字段名
    }))
    
    console.log('[购物车] 结算 - 准备结算的商品:', checkoutItems)
    wx.setStorageSync('selectedCartItems', checkoutItems)
    
    wx.navigateTo({
      url: '/pages/order/confirm/index?from=cart'
    })
  }
}))