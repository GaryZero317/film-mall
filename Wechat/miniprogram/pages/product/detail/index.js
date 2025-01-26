// pages/product/detail/index.js
import { getProductDetail, getProductImages } from '../../../api/product'
import { addToCart } from '../../../api/cart'

Page({
  data: {
    product: null,
    images: [],
    loading: false,
    quantity: 1,
    currentImageIndex: 0,
    imageList: [], // 用于存储所有图片URL
    mainImage: '',
    formattedAmount: '0.00' // 添加格式化后的金额
  },

  onLoad(options) {
    const { id } = options
    if (id) {
      this.loadProductDetail(id)
    }
  },

  // 加载商品详情
  async loadProductDetail(id) {
    try {
      this.setData({ loading: true })
      const res = await getProductDetail(id)
      console.log('商品详情响应:', res)
      
      if (res && res.code === 0 && res.data) {
        // 格式化金额
        const formattedAmount = res.data.amount ? (res.data.amount / 100).toFixed(2) : '0.00'
        
        this.setData({ 
          product: res.data,
          formattedAmount,
          loading: false
        })
        // 设置页面标题
        wx.setNavigationBarTitle({
          title: res.data.name || '商品详情'
        })
        // 加载商品图片
        await this.loadProductImages(id)
      } else {
        throw new Error(res?.msg || '获取商品详情失败')
      }
    } catch (error) {
      console.error('加载商品详情失败:', error)
      wx.showToast({
        title: error.message || '加载商品详情失败',
        icon: 'none'
      })
      this.setData({ loading: false })
    }
  },

  // 加载商品图片
  async loadProductImages(productId) {
    try {
      const res = await getProductImages(productId)
      if (res && res.code === 0 && res.data) {
        const mainImage = res.data.find(img => img.isMain)
        const imageList = res.data.map(img => `http://localhost:8001${img.imageUrl}`)
        
        this.setData({
          images: res.data,
          imageList: imageList,
          mainImage: mainImage ? `http://localhost:8001${mainImage.imageUrl}` : (imageList[0] || 'http://localhost:8001/uploads/placeholder.png')
        })
      }
    } catch (error) {
      console.error('获取商品图片失败:', error)
      this.setData({
        mainImage: 'http://localhost:8001/uploads/placeholder.png'
      })
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
        productId: product.id,
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
    const { imageList } = this.data
    
    if (imageList && imageList.length > 0) {
      wx.previewImage({
        current,
        urls: imageList
      })
    }
  }
})