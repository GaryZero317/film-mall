import { categories } from '../../utils/categories'
import { searchProducts } from '../../api/product'

const app = getApp()

Page({
  data: {
    categories: categories,
    currentCategory: categories[0],
    products: [],
    loading: false,
    hasMore: true,
    page: 1,
    pageSize: 10
  },

  onLoad() {
    this.loadProducts()
  },

  // 加载分类列表
  async loadCategories() {
    try {
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/categories`,
        method: 'GET'
      })

      if (res.statusCode === 200) {
        const categories = res.data
        this.setData({ 
          categories,
          currentCategory: categories[0]
        })
        // 加载第一个分类的商品
        this.loadProducts()
      }
    } catch (error) {
      console.error('加载分类列表失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    }
  },

  // 切换分类
  onCategoryClick(e) {
    const category = e.currentTarget.dataset.category
    this.setData({ 
      currentCategory: category,
      products: []
    })
    this.loadProducts()
  },

  // 加载商品
  async loadProducts() {
    const { currentCategory } = this.data
    if (!currentCategory) return

    try {
      this.setData({ loading: true })
      // 使用分类关键词搜索商品
      const keywords = currentCategory.keywords.join('|')
      const res = await searchProducts({ keywords })
      console.log('分类商品搜索结果:', res)
      
      if (res.code === 0) {
        this.setData({ 
          products: res.data.list || []
        })
      } else {
        wx.showToast({
          title: res.msg || '加载商品失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载商品失败:', error)
      wx.showToast({
        title: '加载商品失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 点击商品
  onProductClick(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    })
  },

  // 下拉刷新
  async onPullDownRefresh() {
    this.setData({
      products: [],
      page: 1,
      hasMore: true
    })
    await this.loadProducts()
    wx.stopPullDownRefresh()
  },

  // 上拉加载更多
  onLoadMore() {
    this.loadProducts()
  }
}) 