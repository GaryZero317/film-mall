import { getProductList } from '../../api/product'
import { searchProducts } from '../../api/product'

const categories = [
  {
    id: 1,
    name: '彩色胶卷',
    icon: '/assets/images/categories/color.png',
    keywords: ['彩色']
  },
  {
    id: 2,
    name: '黑白胶卷',
    icon: '/assets/images/categories/bw.png',
    keywords: ['黑白']
  },
  {
    id: 3,
    name: '中画幅胶卷',
    icon: '/assets/images/categories/medium.png',
    keywords: ['120']
  },
  {
    id: 4,
    name: '胶卷冲洗',
    icon: '/assets/images/categories/development.png',
    keywords: ['冲洗']
  },
  {
    id: 5,
    name: '拍立得',
    icon: '/assets/images/categories/instant.png',
    keywords: ['拍立得']
  },
  {
    id: 6,
    name: '电影胶片',
    icon: '/assets/images/categories/cinema.png',
    keywords: ['电影卷']
  }
]

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
      // 使用分类的关键词搜索商品
      const keyword = currentCategory.keywords[0]
      console.log('使用关键词搜索商品:', keyword)
      
      const res = await searchProducts({ 
        keyword,
        page: 1,
        pageSize: 50
      })
      console.log('商品列表响应:', res)
      
      if (res.code === 0) {
        const products = res.data.list || []
        // 处理图片URL：mainImage已经包含了完整的URL路径
        const processedProducts = products.map(item => ({
          ...item,
          imageUrl: item.mainImage || '/assets/images/default.png'
        }))
        
        this.setData({ 
          products: processedProducts
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