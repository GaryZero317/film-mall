import { getProductList, searchProducts, getProductImages } from '../../api/product'

Page({
  data: {
    banners: [
      {
        id: 1,
        imageUrl: '/assets/images/banner1.jpg',
        link: '/pages/product/detail/index?id=1'
      },
      {
        id: 2,
        imageUrl: '/assets/images/banner2.jpg',
        link: '/pages/product/detail/index?id=2'
      },
      {
        id: 3,
        imageUrl: '/assets/images/banner3.jpg',
        link: '/pages/product/detail/index?id=3'
      }
    ],
    products: [],
    loading: false,
    searchValue: '',
    page: 1,
    pageSize: 10,
    hasMore: true
  },

  onLoad() {
    this.loadProducts(true)
  },

  // 加载商品列表
  async loadProducts(refresh = false) {
    if (this.data.loading || (!refresh && !this.data.hasMore)) return

    try {
      this.setData({ loading: true })
      
      // 构建请求参数
      const params = {
        page: refresh ? 1 : this.data.page,
        pageSize: this.data.pageSize,
        keyword: this.data.searchValue
      }

      console.log('开始加载商品列表，参数:', params)
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
                console.log('商品ID:', product.id, '主图信息:', mainImage)
                const imageUrl = mainImage ? `http://localhost:8001${mainImage.imageUrl}` : (imageRes.data[0] ? `http://localhost:8001${imageRes.data[0].imageUrl}` : '')
                console.log('最终图片URL:', imageUrl)
                return {
                  ...product,
                  mainImage: imageUrl || '/assets/images/default.png'
                }
              }
              return {
                ...product,
                mainImage: '/assets/images/default.png'
              }
            } catch (error) {
              console.error('获取商品图片失败:', error)
              return {
                ...product,
                mainImage: '/assets/images/default.png'
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
          hasMore,
          products
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
  },

  // 点击banner
  onBannerTap(e) {
    const item = e.currentTarget.dataset.item
    if (item.type === 'product') {
      wx.navigateTo({
        url: `/pages/product/detail/index?id=${item.productId}`
      })
    }
  },

  // 跳转到搜索页面
  goToSearch() {
    wx.navigateTo({
      url: '/pages/search/index'
    })
  },

  // 胶片冲洗服务
  navigateToFilmCreate() {
    // 判断是否登录
    const token = wx.getStorageSync('token')
    if (!token) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateTo({
          url: '/pages/login/index'
        })
      }, 1500)
      return
    }
    
    wx.navigateTo({
      url: '/pages/film/create/index'
    })
  }
}) 