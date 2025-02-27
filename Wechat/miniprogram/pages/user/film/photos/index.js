import { getFilmOrderDetail } from '../../../../api/film'

Page({
  data: {
    orderId: null,
    orderDetail: null,
    photos: [],
    loading: false,
    baseUrl: '' // 基础URL
  },

  onLoad(options) {
    if (options.id) {
      // 获取基础URL
      const baseUrl = wx.getStorageSync('baseUrl') || 'http://localhost:8007'
      this.setData({
        orderId: options.id,
        baseUrl: baseUrl
      })
      
      this.loadOrderDetail()
    } else {
      wx.showToast({
        title: '订单ID不存在',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    }
  },

  // 加载订单详情
  async loadOrderDetail() {
    if (!this.data.orderId || this.data.loading) return
    
    try {
      this.setData({ loading: true })
      
      const res = await getFilmOrderDetail(this.data.orderId)
      
      if (res.code === 0 && res.data) {
        // 提取照片列表
        let photos = res.data.photos || []
        
        // 处理照片URL
        photos = photos.map(photo => {
          // 如果URL不是以http开头，则添加基础URL
          if (photo.url && !photo.url.startsWith('http')) {
            // 确保URL格式正确
            let url = photo.url
            if (!url.startsWith('/')) {
              url = '/' + url
            }
            photo.url = this.data.baseUrl + url
          }
          console.log('照片URL:', photo.url)
          return photo
        })
        
        this.setData({
          orderDetail: res.data,
          photos: photos
        })
        
        if (photos.length === 0) {
          wx.showToast({
            title: '该订单暂无照片',
            icon: 'none'
          })
        }
      } else {
        wx.showToast({
          title: res.msg || '获取订单详情失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('获取胶片冲洗订单详情失败:', error)
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },
  
  // 预览照片
  previewPhoto(e) {
    const { url } = e.currentTarget.dataset
    const { photos } = this.data
    
    const urls = photos.map(item => item.url)
    
    wx.previewImage({
      current: url,
      urls: urls
    })
  },
  
  // 图片加载出错
  onImageError(e) {
    const { index } = e.currentTarget.dataset
    console.error(`图片 ${index+1} 加载失败:`, this.data.photos[index].url)
    
    // 可以在这里更新图片URL或显示错误占位图
    let photos = [...this.data.photos]
    // 如果需要可以更新错误图片的URL
    // photos[index].url = '/assets/images/photo-error.png'
    
    // this.setData({ photos })
  },
  
  // 返回上一页
  onBack() {
    wx.navigateBack()
  }
}) 