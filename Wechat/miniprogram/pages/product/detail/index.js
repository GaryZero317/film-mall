// pages/product/detail/index.js
import { getProductDetail, getProductImages } from '../../../api/product'
import { addToCart } from '../../../api/cart'

Page({
  data: {
    product: null,
    images: [],
    loading: false,
    quantity: 1,
    currentImageIndex: 0
  },

  onLoad(options) {
    const { id } = options
    if (id) {
      this.loadProductDetail(id)
      this.loadProductImages(id)
    }
  },

  // 加载商品详情
  async loadProductDetail(id) {
    try {
      this.setData({ loading: true })
      const res = await getProductDetail(id)
      this.setData({ 
        product: res,
        loading: false
      })
      // 设置页面标题
      wx.setNavigationBarTitle({
        title: this.data.product.name
      })
    } catch (error) {
      console.error('加载商品详情失败:', error)
      wx.showToast({
        title: '加载商品详情失败',
        icon: 'none'
      })
      this.setData({ loading: false })
    }
  },

  // 加载商品图片
  async loadProductImages(productId) {
    try {
      const res = await getProductImages(productId)
      this.setData({ 
        images: res || []
      })
    } catch (error) {
      console.error('加载商品图片失败:', error)
    }
  },

  // 轮播图切换事件
  onSwiperChange(e) {
    this.setData({
      currentImageIndex: e.detail.current
    })
  },

  // 数量减少
  onQuantityMinus() {
    if (this.data.quantity > 1) {
      this.setData({
        quantity: this.data.quantity - 1
      })
    }
  },

  // 数量增加
  onQuantityPlus() {
    const { product, quantity } = this.data
    if (quantity < product.stock) {
      this.setData({
        quantity: quantity + 1
      })
    } else {
      wx.showToast({
        title: '超出库存数量',
        icon: 'none'
      })
    }
  },

  // 加入购物车
  async addToCart() {
    const { product, quantity } = this.data
    if (!product) return
    
    try {
      await addToCart({
        product_id: product.id,
        quantity: quantity
      })
      
      wx.showToast({
        title: '已加入购物车',
        icon: 'success'
      })
    } catch (error) {
      console.error('加入购物车失败:', error)
      wx.showToast({
        title: '加入购物车失败',
        icon: 'none'
      })
    }
  },

  // 立即购买
  buyNow() {
    const { product, quantity } = this.data
    if (!product) return
    
    wx.navigateTo({
      url: `/pages/order/confirm/index?productId=${product.id}&quantity=${quantity}`
    })
  },

  // 预览图片
  previewImage(e) {
    const { current } = e.currentTarget.dataset
    const { product, images } = this.data
    const urls = []
    
    // 添加主图到预览列表
    if (product.mainImage) {
      urls.push(product.mainImage)
    }
    
    // 添加其他图片到预览列表
    if (images && images.length > 0) {
      urls.push(...images.map(img => img.url))
    }
    
    wx.previewImage({
      current,
      urls
    })
  }
})