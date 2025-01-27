// pages/cart/index.js
import { getCartList, updateCartItem, removeFromCart, updateCartItemStatus } from '../../api/cart'
import { loginGuard } from '../../utils/auth'

Page(loginGuard({
  data: {
    cartItems: [],
    loading: false,
    totalPrice: 0,
    selectedCount: 0,
    allSelected: false
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
      console.log('购物车列表响应:', res)
      
      if (res && res.code === 0) {
        // 处理图片URL
        const cartItems = (res.data.list || []).map(item => ({
          ...item,
          productImage: item.productImage ? `http://localhost:8001${item.productImage}` : '/assets/images/default.png'
        }))
        console.log('处理后的购物车商品列表:', cartItems)
        
        this.setData({ 
          cartItems,
          selectedAll: false,
          totalPrice: 0
        })
        this.calculateTotal()
      } else {
        console.error('获取购物车列表失败:', res)
        wx.showToast({
          title: res?.msg || '获取购物车列表失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载购物车失败:', error)
      wx.showToast({
        title: '加载购物车失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
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
      console.log('更新购物车数量:', { id, quantity: newQuantity })
      const res = await updateCartItem({
        id: parseInt(id),
        quantity: newQuantity
      })
      console.log('更新购物车数量响应:', res)
      
      if (res.code === 0) {
        const cartItems = this.data.cartItems.map(item => {
          if (item.id === id) {
            return { ...item, quantity: newQuantity }
          }
          return item
        })
        
        this.setData({ cartItems })
        this.calculateTotal()
      } else {
        wx.showToast({
          title: res.msg || '更新数量失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('更新数量失败:', error)
      wx.showToast({
        title: '更新数量失败',
        icon: 'none'
      })
    }
  },

  // 删除商品
  async onDelete(e) {
    const { id } = e.currentTarget.dataset
    try {
      await removeFromCart(id)
      const cartItems = this.data.cartItems.filter(item => item.id !== id)
      this.setData({ cartItems })
      this.calculateTotal()
      wx.showToast({
        title: '删除成功',
        icon: 'success'
      })
    } catch (error) {
      console.error('删除失败:', error)
      wx.showToast({
        title: '删除失败',
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
    this.calculateTotal()
    
    try {
      await updateCartItemStatus({
        id,
        selected: cartItems.find(item => item.id === id).selected
      })
    } catch (error) {
      console.error('更新选中状态失败:', error)
      wx.showToast({
        title: '更新状态失败',
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
    this.calculateTotal()
  },

  // 计算总价和选中数量
  calculateTotal() {
    const { cartItems } = this.data
    const selectedItems = cartItems.filter(item => item.selected)
    const totalPrice = selectedItems.reduce((sum, item) => sum + item.price * item.quantity, 0)
    const selectedCount = selectedItems.length
    const allSelected = selectedCount === cartItems.length && cartItems.length > 0
    
    this.setData({
      totalPrice,
      selectedCount,
      allSelected
    })
  },

  // 去结算
  onCheckout() {
    const selectedItems = this.data.cartItems.filter(item => item.selected)
    if (selectedItems.length === 0) {
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
      price: parseFloat(item.price),
      quantity: parseInt(item.quantity),
      cover_image: item.productImage
    }))
    console.log('[购物车] 准备结算的商品:', checkoutItems)
    wx.setStorageSync('selectedCartItems', checkoutItems)
    
    wx.navigateTo({
      url: '/pages/order/confirm/index?from=cart'
    })
  }
}))