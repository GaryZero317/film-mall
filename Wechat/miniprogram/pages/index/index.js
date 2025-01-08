import { getProductList, getCategories } from '../../api/product'

Page({
  data: {
    banners: [],
    categories: [],
    products: [],
    loading: false,
    selectedCategory: null,
    searchValue: '',
    page: 1,
    pageSize: 10,
    hasMore: true
  },

  onLoad() {
    this.loadCategories()
    this.loadProducts()
  },

  // 加载分类
  async loadCategories() {
    try {
      const res = await getCategories()
      if (res && res.data) {
        this.setData({ categories: res.data })
      }
    } catch (error) {
      console.error('加载分类失败:', error)
      wx.showToast({
        title: '加载分类失败',
        icon: 'none'
      })
    }
  },

  // 加载商品列表
  async loadProducts(params = {}) {
    if (this.data.loading || (!this.data.hasMore && !params.refresh)) return

    try {
      this.setData({ loading: true })
      
      const requestParams = {
        page: params.refresh ? 1 : this.data.page,
        pageSize: this.data.pageSize,
        ...params
      }

      const res = await getProductList(requestParams)
      
      if (res && res.data) {
        const { list, total } = res.data
        const formattedProducts = list.map(item => ({
          id: item.id,
          name: item.name,
          price: item.amount / 100, // 将分转换为元
          stock: item.stock,
          cover_image: item.mainImage,
          description: item.desc
        }))

        this.setData({
          products: params.refresh ? formattedProducts : [...this.data.products, ...formattedProducts],
          page: params.refresh ? 2 : this.data.page + 1,
          hasMore: this.data.products.length < total
        })
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

  // 下拉刷新
  async onPullDownRefresh() {
    await this.loadProducts({ refresh: true })
    wx.stopPullDownRefresh()
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadProducts()
    }
  },

  // 搜索输入
  onSearchInput(e) {
    this.setData({
      searchValue: e.detail.value
    })
  },

  // 执行搜索
  onSearch() {
    const { searchValue, selectedCategory } = this.data
    const params = {
      refresh: true
    }
    if (searchValue) {
      params.keyword = searchValue
    }
    if (selectedCategory) {
      params.category_id = selectedCategory
    }
    this.loadProducts(params)
  },

  // 切换分类
  async switchCategory(e) {
    const categoryId = e.currentTarget.dataset.id
    this.setData({ 
      selectedCategory: this.data.selectedCategory === categoryId ? null : categoryId 
    })
    
    const params = {
      refresh: true
    }
    if (this.data.selectedCategory) {
      params.category_id = this.data.selectedCategory
    }
    if (this.data.searchValue) {
      params.keyword = this.data.searchValue
    }
    await this.loadProducts(params)
  },

  // 跳转到商品详情页
  goToDetail(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    })
  }
}) 