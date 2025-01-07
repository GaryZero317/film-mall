// pages/product/detail/index.js
import { getProductDetail } from '../../../api/product'
import { addToCart } from '../../../api/cart'

Page({
  data: {
    product: {
      id: 1,
      name: 'Kodak Gold 200 柯达金200胶卷',
      price: 65,
      stock: 100,
      cover_image: '/assets/images/kodak-gold-200.jpg',
      brand: 'Kodak',
      iso: 200,
      exposures: 36,
      expiry_date: '2025-12',
      description: '柯达金系列是一款平价的日常用胶，色彩明亮，颗粒适中，非常适合新手入门。',
      detail: `
        <div style="font-size: 28rpx; color: #333; line-height: 1.6;">
          <p>【产品特点】</p>
          <ul>
            <li>色彩鲜艳，偏暖色调</li>
            <li>颗粒感适中，画面清晰</li>
            <li>宽容度高，适合新手使用</li>
            <li>性价比极高的日常用胶</li>
          </ul>
          <p>【适用场景】</p>
          <ul>
            <li>日常生活记录</li>
            <li>旅行摄影</li>
            <li>人像摄影</li>
            <li>街头摄影</li>
          </ul>
          <p>【使用建议】</p>
          <ul>
            <li>阳光充足时效果最佳</li>
            <li>可以适当过曝1/3档提升色彩饱和度</li>
            <li>室内拍摄建议使用闪光灯</li>
          </ul>
          <p>【产品规格】</p>
          <ul>
            <li>规格：35mm</li>
            <li>感光度：ISO 200</li>
            <li>曝光张数：36张</li>
            <li>产地：美国</li>
          </ul>
        </div>
      `
    },
    loading: false,
    quantity: 1
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
      this.setData({ 
        product: res.data || this.data.product,
        loading: false
      })
      // 设置页面标题
      wx.setNavigationBarTitle({
        title: this.data.product.name
      })
    } catch (error) {
      console.error('加载商品详情失败:', error)
      this.setData({ loading: false })
    }
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
  }
})