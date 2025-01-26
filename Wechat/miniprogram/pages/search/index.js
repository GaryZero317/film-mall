import { searchProducts, getProductImages } from '../../api/product'

const MAX_HISTORY = 10
const HISTORY_KEY = 'search_history'

Page({
  data: {
    keyword: '',
    history: [],
    products: [],
    loading: false,
    page: 1,
    pageSize: 10,
    hasMore: true,
    // 热门品牌
    brands: [
      '柯达',
      '富士',
      '伊尔福',
      '宝丽来'
    ],
    // 热门型号
    models: [
      'Portra 400',
      'Gold 200',
      'ColorPlus 200',
      '拍立得',
      'C200',
      'HP5',
      '炮塔',
      '135',
      '120'
    ]
  },

  onLoad() {
    // 加载搜索历史
    const history = wx.getStorageSync(HISTORY_KEY) || []
    this.setData({ history })
  },

  // 输入关键词
  onInput(e) {
    this.setData({
      keyword: e.detail.value
    })
  },

  // 清空输入
  onClear() {
    this.setData({
      keyword: '',
      products: []
    })
  },

  // 取消搜索
  onCancel() {
    wx.navigateBack()
  },

  // 点击历史记录
  onTagTap(e) {
    const { keyword } = e.currentTarget.dataset
    this.setData({ keyword }, () => {
      this.search()
    })
  },

  // 清空历史记录
  clearHistory() {
    wx.showModal({
      title: '提示',
      content: '确定要清空搜索历史吗？',
      success: (res) => {
        if (res.confirm) {
          wx.removeStorageSync(HISTORY_KEY)
          this.setData({ history: [] })
        }
      }
    })
  },

  // 保存搜索历史
  saveHistory(keyword) {
    let history = this.data.history
    // 移除已存在的相同关键词
    history = history.filter(item => item !== keyword)
    // 添加到开头
    history.unshift(keyword)
    // 限制数量
    if (history.length > MAX_HISTORY) {
      history = history.slice(0, MAX_HISTORY)
    }
    // 保存到本地和data
    wx.setStorageSync(HISTORY_KEY, history)
    this.setData({ history })
  },

  // 执行搜索
  async search(loadMore = false) {
    if (this.data.loading || (!loadMore && !this.data.hasMore)) return
    if (!this.data.keyword.trim()) return

    const page = loadMore ? this.data.page : 1
    
    try {
      this.setData({ loading: true })
      
      const params = {
        page,
        pageSize: this.data.pageSize,
        keyword: this.data.keyword
      }

      const res = await searchProducts(params)
      
      if (res && res.code === 0 && res.data) {
        const { list, total } = res.data
        
        // 获取商品图片
        const productsWithImages = await Promise.all(
          list.map(async (product) => {
            try {
              const imageRes = await getProductImages(product.id)
              if (imageRes && imageRes.code === 0 && imageRes.data) {
                const mainImage = imageRes.data.find(img => img.isMain)
                const imageUrl = mainImage ? `http://localhost:8001${mainImage.imageUrl}` : (imageRes.data[0] ? `http://localhost:8001${imageRes.data[0].imageUrl}` : '')
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

        const products = loadMore ? [...this.data.products, ...productsWithImages] : productsWithImages
        const hasMore = products.length < total

        this.setData({
          products,
          page: page + 1,
          hasMore,
          loading: false
        })

        // 保存搜索历史
        if (!loadMore) {
          this.saveHistory(this.data.keyword)
        }
      }
    } catch (error) {
      console.error('搜索失败:', error)
      wx.showToast({
        title: '搜索失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 确认搜索
  onSearch() {
    this.search()
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore) {
      this.search(true)
    }
  },

  // 跳转到商品详情
  goToDetail(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    })
  }
}) 