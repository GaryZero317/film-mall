import { getProductList, searchProducts } from '../../api/product'

Page({
  data: {
    banners: [],
    categories: [
      { id: 0, name: '全部' },
      { id: 1, name: '彩色胶卷' },
      { id: 2, name: '黑白胶卷' },
      { id: 3, name: '135胶卷' },
      { id: 4, name: '120胶卷' },
      { id: 5, name: '反转片' },
      { id: 6, name: '拍立得' }
    ],
    // 用于布局的分类数组，将分类分成两行
    categoryRows: [
      [
        { id: 1, name: '彩色胶卷' },
        { id: 2, name: '黑白胶卷' },
        { id: 3, name: '135胶卷' }
      ],
      [
        { id: 4, name: '120胶卷' },
        { id: 5, name: '反转片' },
        { id: 6, name: '拍立得' }
      ]
    ],
    products: [],
    loading: false,
    selectedCategory: 0,
    searchValue: '',
    page: 1,
    pageSize: 10,
    hasMore: true
  },

  onLoad() {
    this.loadProducts()
  },

  // 加载商品列表
  async loadProducts(reset = false) {
    if (this.data.loading) return
    
    try {
      this.setData({ loading: true })
      const { page, pageSize, selectedCategory, searchValue, categories } = this.data
      
      // 获取当前分类名称作为搜索关键词
      let keyword = searchValue
      if (selectedCategory !== 0) {
        const category = categories.find(c => c.id === selectedCategory)
        if (category) {
          keyword = keyword ? `${keyword} ${category.name}` : category.name
        }
      }

      let res
      if (keyword) {
        // 使用搜索接口
        res = await searchProducts(keyword)
      } else {
        // 使用商品列表接口
        const params = {
          page: page,
          pageSize: pageSize
        }
        res = await getProductList(params)
      }
      
      if (res && res.list) {
        this.setData({
          products: reset ? res.list : [...this.data.products, ...res.list],
          hasMore: res.list.length === pageSize
        })
      } else {
        throw new Error('获取商品列表失败')
      }
    } catch (error) {
      console.error('加载商品列表失败:', error)
      wx.showToast({
        title: '加载商品列表失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 切换分类
  async switchCategory(e) {
    const categoryId = e.currentTarget.dataset.id
    this.setData({ 
      selectedCategory: categoryId,
      page: 1,
      products: []
    })
    await this.loadProducts(true)
  },

  // 搜索输入
  onSearchInput(e) {
    this.setData({
      searchValue: e.detail.value
    })
  },

  // 执行搜索
  async onSearch() {
    this.setData({
      page: 1,
      products: []
    })
    await this.loadProducts(true)
  },

  // 下拉刷新
  async onPullDownRefresh() {
    this.setData({
      page: 1,
      products: []
    })
    await this.loadProducts(true)
    wx.stopPullDownRefresh()
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadProducts()
    }
  },

  // 跳转到商品详情页
  goToDetail(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    })
  }
}) 