import { getProductList, getCategories } from '../../api/product'

Page({
  data: {
    banners: [
      {
        id: 1,
        image_url: '/assets/images/banner1.jpg',
        product_id: 1,
        title: 'Kodak 柯达 Gold 200 胶卷'
      },
      {
        id: 2,
        image_url: '/assets/images/banner2.jpg',
        product_id: 2,
        title: 'Fujifilm 富士 C200 胶卷'
      }
    ],
    categories: [
      { id: 1, name: '彩色负片' },
      { id: 2, name: '黑白负片' },
      { id: 3, name: '正片' },
      { id: 4, name: '一次性相机' },
      { id: 5, name: '中画幅胶卷' },
      { id: 6, name: '拍立得' }
    ],
    products: [
      {
        id: 1,
        name: 'Kodak Gold 200 柯达金200胶卷',
        price: 65,
        stock: 100,
        cover_image: '/assets/images/kodak-gold-200.jpg',
        brand: 'Kodak',
        iso: 200,
        exposures: 36,
        description: '柯达金系列是一款平价的日常用胶，色彩明亮，颗粒适中，非常适合新手入门。'
      },
      {
        id: 2,
        name: 'Fujifilm C200 富士C200胶卷',
        price: 55,
        stock: 80,
        cover_image: '/assets/images/fuji-c200.jpg',
        brand: 'Fujifilm',
        iso: 200,
        exposures: 36,
        description: '富士C200是一款性价比极高的彩色负片，色彩还原自然，绿色和蓝色的表现尤为出色。'
      },
      {
        id: 3,
        name: 'Ilford HP5 Plus 400 黑白胶卷',
        price: 75,
        stock: 50,
        cover_image: '/assets/images/ilford-hp5.jpg',
        brand: 'Ilford',
        iso: 400,
        exposures: 36,
        description: '经典的黑白胶卷，颗粒感细腻，层次丰富，是黑白摄影爱好者的首选。'
      }
    ],
    loading: false,
    selectedCategory: null,
    searchValue: ''
  },

  onLoad() {
    this.loadCategories()
    this.loadProducts()
  },

  // 加载分类
  async loadCategories() {
    try {
      const res = await getCategories()
      this.setData({ categories: res.data || this.data.categories })
    } catch (error) {
      console.error('加载分类失败:', error)
    }
  },

  // 加载商品列表
  async loadProducts(params = {}) {
    try {
      this.setData({ loading: true })
      const res = await getProductList(params)
      this.setData({ 
        products: res.data || this.data.products,
        loading: false
      })
    } catch (error) {
      console.error('加载商品列表失败:', error)
      this.setData({ loading: false })
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
    const params = { keyword: searchValue }
    if (selectedCategory) {
      params.category_id = selectedCategory
    }
    this.loadProducts(params)
  },

  // 切换分类
  async switchCategory(e) {
    const categoryId = e.currentTarget.dataset.id
    this.setData({ selectedCategory: categoryId })
    const params = { category_id: categoryId }
    if (this.data.searchValue) {
      params.keyword = this.data.searchValue
    }
    await this.loadProducts(params)
  },

  // 轮播图点击
  onBannerTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    })
  },

  // 跳转到商品详情页
  goToDetail(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    })
  }
}) 