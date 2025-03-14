import { uploadWork, updateWork, getWorkDetail, uploadWorkImage, deleteWorkImage } from '../../../../api/community'

Page({
  data: {
    id: null,
    title: '',
    description: '',
    filmType: '',
    filmBrand: '',
    cameraModel: '',
    lensInfo: '',
    exifData: '',
    images: [],
    tempImages: [],
    loading: false,
    isEdit: false,
    maxImageCount: 9,
    filmTypes: ['彩色胶卷', '黑白胶卷', '中画幅胶卷', '拍立得', '电影胶片'],
    filmBrands: [
      '柯达 Kodak',
      '富士 Fujifilm',
      '乐凯 Lucky',
      '伊尔福 Ilford',
      '宝丽来 Polaroid',
      'Cinestill',
      '其他'
    ],
    cameraBrands: [
      '佳能 Canon',
      '尼康 Nikon',
      '索尼 Sony',
      '富士 Fujifilm',
      '徕卡 Leica',
      '哈苏 Hasselblad',
      '宾得 Pentax',
      '奥林巴斯 Olympus',
      '松下 Panasonic',
      '禄来 Rollei',
      '勃朗尼卡 Bronica',
      '玛米亚 Mamiya',
      '宝丽来 Polaroid',
      '其他'
    ],
    lensTypes: [
      '广角定焦 14mm',
      '广角定焦 20mm',
      '广角定焦 24mm',
      '广角定焦 28mm',
      '广角定焦 35mm',
      '标准定焦 50mm',
      '人像定焦 85mm',
      '人像定焦 105mm',
      '远摄定焦 135mm',
      '远摄定焦 200mm',
      '广角变焦 16-35mm',
      '标准变焦 24-70mm',
      '远摄变焦 70-200mm',
      '微距镜头',
      '移轴镜头',
      '其他'
    ],
    exifOptions: [
      'jpg',
      'png',
      'gif',
      'raw',
      '其他'
    ]
  },

  onLoad(options) {
    if (options.id) {
      this.setData({
        id: parseInt(options.id),
        isEdit: true
      })
      this.loadWorkDetail()
    }
  },

  async loadWorkDetail() {
    if (!this.data.id) return
    this.setData({ loading: true })
    
    try {
      const res = await getWorkDetail(this.data.id)
      if (res.code === 0 || res.code === 200) {
        const work = res.data.work || res.data
        this.setData({
          title: work.title || '',
          description: work.description || '',
          filmType: work.film_type || '',
          filmBrand: work.film_brand || '',
          cameraModel: work.camera || '',
          lensInfo: work.lens || '',
          exifData: work.exif_info || '',
          images: work.images || []
        })
      } else {
        wx.showToast({
          title: res.msg || '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载作品详情失败:', error)
      wx.showToast({
        title: error.message || '加载失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  onTitleInput(e) {
    this.setData({ title: e.detail.value })
  },

  onDescriptionInput(e) {
    this.setData({ description: e.detail.value })
  },

  onFilmTypeChange(e) {
    this.setData({ filmType: this.data.filmTypes[e.detail.value] })
  },

  onFilmBrandChange(e) {
    this.setData({ filmBrand: this.data.filmBrands[e.detail.value] })
  },

  onCameraModelChange(e) {
    this.setData({ cameraModel: this.data.cameraBrands[e.detail.value] })
  },

  onLensInfoChange(e) {
    this.setData({ lensInfo: this.data.lensTypes[e.detail.value] })
  },

  onExifDataChange(e) {
    this.setData({ exifData: this.data.exifOptions[e.detail.value] })
  },

  onChooseImage() {
    const { images, tempImages, maxImageCount } = this.data
    const remainCount = maxImageCount - images.length - tempImages.length
    
    if (remainCount <= 0) {
      wx.showToast({
        title: `最多只能上传${maxImageCount}张图片`,
        icon: 'none'
      })
      return
    }

    wx.chooseImage({
      count: remainCount,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        this.setData({
          tempImages: [...tempImages, ...res.tempFilePaths]
        })
      }
    })
  },

  onPreviewImage(e) {
    const { current } = e.currentTarget.dataset
    const { images, tempImages } = this.data
    const allImages = [...images, ...tempImages]
    
    wx.previewImage({
      current,
      urls: allImages
    })
  },

  onRemoveImage(e) {
    const { index, type } = e.currentTarget.dataset
    
    if (type === 'temp') {
      const tempImages = [...this.data.tempImages]
      tempImages.splice(index, 1)
      this.setData({ tempImages })
    } else {
      const images = [...this.data.images]
      const removedImage = images[index]
      images.splice(index, 1)
      this.setData({ images })

      // 如果是编辑模式，删除已上传的图片
      if (this.data.isEdit && removedImage) {
        this.deleteImage(removedImage)
      }
    }
  },

  async deleteImage(imageUrl) {
    try {
      // 从URL中提取图片ID
      const matches = imageUrl.match(/\/([^\/]+)$/)
      if (matches && matches[1]) {
        const imageId = matches[1]
        await deleteWorkImage(imageId)
      }
    } catch (error) {
      console.error('删除图片失败:', error)
    }
  },

  validateForm() {
    const { title, description, images, tempImages } = this.data
    
    if (!title.trim()) {
      wx.showToast({
        title: '请输入作品标题',
        icon: 'none'
      })
      return false
    }
    
    if (!description.trim()) {
      wx.showToast({
        title: '请输入作品描述',
        icon: 'none'
      })
      return false
    }
    
    if (images.length === 0 && tempImages.length === 0) {
      wx.showToast({
        title: '请至少上传一张图片',
        icon: 'none'
      })
      return false
    }
    
    return true
  },

  async onSubmit() {
    if (!this.validateForm()) {
      return
    }
    
    this.setData({ loading: true })
    
    try {
      // 获取用户信息
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo) {
        throw new Error('获取用户信息失败，请重新登录')
      }

      // 使用显示名称作为作者名
      const authorName = userInfo.name || userInfo.nickName || userInfo.mobile
      if (!authorName) {
        throw new Error('获取用户名失败，请重新登录')
      }
      
      // 先创建作品
      const workData = {
        title: this.data.title,
        description: this.data.description,
        film_type: this.data.filmType,
        film_brand: this.data.filmBrand,
        camera: this.data.cameraModel,
        lens: this.data.lensInfo,
        exif_info: this.data.exifData,
        status: 1,
        cover_url: '',
        images: [],
        name: authorName
      }
      
      let workRes
      if (this.data.isEdit) {
        // 更新作品
        workRes = await updateWork(this.data.id, workData)
      } else {
        // 创建新作品
        workRes = await uploadWork(workData)
      }

      if (workRes.code !== 0 && workRes.code !== 200) {
        throw new Error(workRes.msg || '操作失败')
      }

      const workId = this.data.isEdit ? this.data.id : workRes.data.id

      // 上传新图片
      const { tempImages } = this.data
      const uploadedImages = []
      let coverUrl = ''  // 用于存储第一张图片的URL作为封面
      
      if (tempImages.length > 0) {
        for (let i = 0; i < tempImages.length; i++) {
          const filePath = tempImages[i]
          try {
            wx.showLoading({
              title: `上传图片 ${i + 1}/${tempImages.length}`,
              mask: true
            })
            
            const res = await uploadWorkImage(filePath, workId)
            if (res.code === 0 || res.code === 200) {
              uploadedImages.push(res.data.url)
              // 将第一张上传的图片设为封面
              if (i === 0) {
                coverUrl = res.data.url
              }
            } else {
              throw new Error(res.msg || '上传失败')
            }
          } catch (error) {
            console.error('上传图片失败:', error)
            wx.hideLoading()
            wx.showToast({
              title: error.message || '上传图片失败',
              icon: 'none'
            })
            return
          }
        }
        
        wx.hideLoading()
      }
      
      // 更新作品的图片列表和封面
      if (uploadedImages.length > 0) {
        const allImages = [...this.data.images, ...uploadedImages]
        await updateWork(workId, { 
          ...workData, 
          images: allImages,
          cover_url: coverUrl || this.data.images[0]?.url || ''  // 使用第一张图片作为封面
        })
      }

      wx.showToast({
        title: this.data.isEdit ? '更新成功' : '发布成功',
        icon: 'success'
      })
      
      // 返回上一页并刷新列表
      setTimeout(() => {
        const pages = getCurrentPages()
        const prevPage = pages[pages.length - 2]
        if (prevPage) {
          prevPage.onRefresh()
        }
        wx.navigateBack()
      }, 1500)
    } catch (error) {
      console.error('提交失败:', error)
      wx.showToast({
        title: error.message || '提交失败，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  onCancel() {
    wx.navigateBack()
  }
}) 