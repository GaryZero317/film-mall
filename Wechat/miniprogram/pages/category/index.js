const app = getApp()

Page({
  data: {
    categories: [],
    currentCategory: null,
    products: [],
    loading: false,
    hasMore: true,
    page: 1,
    pageSize: 10
  },

  onLoad() {
    this.loadCategories()
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
    const id = e.currentTarget.dataset.id
    const currentCategory = this.data.categories.find(item => item.id === id)
    
    this.setData({
      currentCategory,
      products: [],
      page: 1,
      hasMore: true
    })
    
    this.loadProducts()
  },

  // 加载商品列表
  async loadProducts() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      const { page, pageSize, currentCategory } = this.data
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/products`,
        method: 'GET',
        data: {
          page,
          page_size: pageSize,
          category_id: currentCategory.id
        }
      })

      if (res.statusCode === 200) {
        const { data, total } = res.data
        this.setData({
          products: [...this.data.products, ...data],
          page: page + 1,
          hasMore: this.data.products.length + data.length < total
        })
      } else {
        wx.showToast({
          title: '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载商品列表失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 点击商品跳转到详情页
  onProductClick(e) {
    const id = e.currentTarget.dataset.id
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