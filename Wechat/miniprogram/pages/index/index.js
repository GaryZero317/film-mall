import { getProductList, searchProducts, getProductImages } from '../../api/product'

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
  async loadProducts(refresh = false) {
    if (this.data.loading || (!refresh && !this.data.hasMore)) return

    try {
      this.setData({ loading: true })
      
      // 构建请求参数
      const params = {
        page: refresh ? 1 : this.data.page,
        pageSize: this.data.pageSize
      }

      // 如果有选中的分类
      if (this.data.selectedCategory !== 0) {
        params.categoryId = this.data.selectedCategory
      }

      // 如果有搜索关键词
      if (this.data.searchValue) {
        params.keyword = this.data.searchValue
      }

      console.log('请求参数:', params)
      const res = await getProductList(params)
      console.log('商品列表响应:', res)

      if (res && res.code === 0 && res.data) {
        const { list, total } = res.data
        
        // 获取每个商品的图片信息
        const productsWithImages = await Promise.all(
          list.map(async (product) => {
            try {
              const imageRes = await getProductImages(product.id)
              console.log('商品图片响应:', imageRes)
              if (imageRes && imageRes.code === 0 && imageRes.data) {
                const mainImage = imageRes.data.find(img => img.isMain)
                console.log('主图信息:', mainImage)
                const imageUrl = mainImage ? `http://localhost:8001${mainImage.imageUrl}` : (imageRes.data[0] ? `http://localhost:8001${imageRes.data[0].imageUrl}` : '')
                console.log('最终图片URL:', imageUrl)
                return {
                  ...product,
                  mainImage: imageUrl || 'http://localhost:8001/uploads/placeholder.png'
                }
              }
              return {
                ...product,
                mainImage: 'http://localhost:8001/uploads/placeholder.png'
              }
            } catch (error) {
              console.error('获取商品图片失败:', error)
              return {
                ...product,
                mainImage: 'http://localhost:8001/uploads/placeholder.png'
              }
            }
          })
        )

        const products = refresh ? productsWithImages : [...this.data.products, ...productsWithImages]
        const page = refresh ? 2 : this.data.page + 1
        const hasMore = products.length < total

        console.log('处理后的数据:', {
          productsLength: products.length,
          page,
          hasMore
        })

        this.setData({
          products,
          page,
          hasMore,
          loading: false
        })
      } else {
        console.error('返回数据格式错误:', res)
        wx.showToast({
          title: res?.msg || '加载失败',
          icon: 'none'
        })
        this.setData({ loading: false })
      }
    } catch (error) {
      console.error('加载商品列表失败:', error)
      wx.showToast({
        title: error.message || '加载失败',
        icon: 'none'
      })
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
    await this.loadProducts(true)
    wx.stopPullDownRefresh()
  },

  // 上拉加载更多
  onReachBottom() {
    this.loadProducts()
  },

  // 跳转到商品详情页
  goToDetail(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    })
  }
}) 