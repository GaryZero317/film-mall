// pages/cart/index.js
import { getCartList, updateCartItem, removeFromCart, updateCartItemStatus } from '../../api/cart'

Page({
  data: {
    cartItems: [
      {
        id: 1,
        product_id: 1,
        name: 'Kodak Gold 200 柯达金200胶卷',
        price: 65,
        stock: 100,
        quantity: 2,
        selected: false,
        cover_image: '/assets/images/kodak-gold-200.jpg'
      },
      {
        id: 2,
        product_id: 2,
        name: 'Fujifilm C200 富士C200胶卷',
        price: 55,
        stock: 80,
        quantity: 1,
        selected: false,
        cover_image: '/assets/images/fuji-c200.jpg'
      }
    ],
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
      this.setData({ 
        cartItems: res.data || this.data.cartItems,
        loading: false
      })
      this.calculateTotal()
    } catch (error) {
      console.error('加载购物车失败:', error)
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
      await updateCartItem({
        id,
        quantity: newQuantity
      })
      
      const cartItems = this.data.cartItems.map(item => {
        if (item.id === id) {
          return { ...item, quantity: newQuantity }
        }
        return item
      })
      
      this.setData({ cartItems })
      this.calculateTotal()
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
    
    wx.navigateTo({
      url: '/pages/order/confirm/index'
    })
  }
})