const app = getApp();
const { uploadWork, updateWork, getWorkDetail } = require('../../../../api/community');

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
    filmTypes: ['135', '120', 'APS-C', '4x5', '其他'],
    filmBrands: ['柯达', '富士', '乐凯', '爱尔福', '其他'],
  },

  onLoad(options) {
    if (options.id) {
      this.setData({
        id: options.id,
        isEdit: true,
        loading: true
      });
      this.loadWorkDetail(options.id);
    }
  },

  async loadWorkDetail(id) {
    try {
      const work = await getWorkDetail(id);
      if (work) {
        this.setData({
          title: work.title || '',
          description: work.description || '',
          filmType: work.filmType || '',
          filmBrand: work.filmBrand || '',
          cameraModel: work.cameraModel || '',
          lensInfo: work.lensInfo || '',
          exifData: work.exifData || '',
          images: work.images || [],
          loading: false
        });
      }
    } catch (error) {
      console.error('加载作品详情失败', error);
      wx.showToast({
        title: '加载作品详情失败',
        icon: 'none'
      });
      this.setData({ loading: false });
    }
  },

  onTitleInput(e) {
    this.setData({ title: e.detail.value });
  },

  onDescriptionInput(e) {
    this.setData({ description: e.detail.value });
  },

  onFilmTypeChange(e) {
    this.setData({ filmType: this.data.filmTypes[e.detail.value] });
  },

  onFilmBrandChange(e) {
    this.setData({ filmBrand: this.data.filmBrands[e.detail.value] });
  },

  onCameraModelInput(e) {
    this.setData({ cameraModel: e.detail.value });
  },

  onLensInfoInput(e) {
    this.setData({ lensInfo: e.detail.value });
  },

  onExifDataInput(e) {
    this.setData({ exifData: e.detail.value });
  },

  onChooseImage() {
    const { images, maxImageCount } = this.data;
    const remainCount = maxImageCount - images.length;
    
    if (remainCount <= 0) {
      wx.showToast({
        title: `最多只能上传${maxImageCount}张图片`,
        icon: 'none'
      });
      return;
    }

    wx.chooseImage({
      count: remainCount,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        // 合并临时文件路径
        const tempFilePaths = res.tempFilePaths;
        const tempFiles = res.tempFiles;
        
        this.setData({
          tempImages: [...this.data.tempImages, ...tempFilePaths]
        });
      }
    });
  },

  onPreviewImage(e) {
    const { current } = e.currentTarget.dataset;
    const { images, tempImages } = this.data;
    const allImages = [...images, ...tempImages];
    
    wx.previewImage({
      current,
      urls: allImages
    });
  },

  onRemoveImage(e) {
    const { index, type } = e.currentTarget.dataset;
    
    if (type === 'temp') {
      const tempImages = [...this.data.tempImages];
      tempImages.splice(index, 1);
      this.setData({ tempImages });
    } else {
      const images = [...this.data.images];
      images.splice(index, 1);
      this.setData({ images });
    }
  },

  validateForm() {
    const { title, description, images, tempImages } = this.data;
    
    if (!title.trim()) {
      wx.showToast({
        title: '请输入作品标题',
        icon: 'none'
      });
      return false;
    }
    
    if (!description.trim()) {
      wx.showToast({
        title: '请输入作品描述',
        icon: 'none'
      });
      return false;
    }
    
    if (images.length === 0 && tempImages.length === 0) {
      wx.showToast({
        title: '请至少上传一张图片',
        icon: 'none'
      });
      return false;
    }
    
    return true;
  },

  async uploadImages() {
    const { tempImages } = this.data;
    const uploadedImages = [];
    
    if (tempImages.length === 0) {
      return [];
    }
    
    for (let i = 0; i < tempImages.length; i++) {
      const filePath = tempImages[i];
      try {
        wx.showLoading({
          title: `上传图片 ${i + 1}/${tempImages.length}`,
          mask: true
        });
        
        const result = await new Promise((resolve, reject) => {
          wx.uploadFile({
            url: `${app.globalData.baseUrl}/api/upload/image`,
            filePath,
            name: 'file',
            header: {
              'Authorization': `Bearer ${wx.getStorageSync('token')}`
            },
            success: (res) => {
              if (res.statusCode === 200) {
                const data = JSON.parse(res.data);
                resolve(data.url);
              } else {
                reject(new Error('上传失败'));
              }
            },
            fail: reject
          });
        });
        
        uploadedImages.push(result);
      } catch (error) {
        console.error('上传图片失败', error);
        wx.hideLoading();
        wx.showToast({
          title: '上传图片失败',
          icon: 'none'
        });
        return null;
      }
    }
    
    wx.hideLoading();
    return uploadedImages;
  },

  async onSubmit() {
    if (!this.validateForm()) {
      return;
    }
    
    this.setData({ loading: true });
    
    try {
      // 上传新图片
      const uploadedImages = await this.uploadImages();
      if (uploadedImages === null) {
        this.setData({ loading: false });
        return;
      }
      
      // 合并已有图片和新上传的图片
      const allImages = [...this.data.images, ...uploadedImages];
      
      const workData = {
        title: this.data.title,
        description: this.data.description,
        filmType: this.data.filmType,
        filmBrand: this.data.filmBrand,
        cameraModel: this.data.cameraModel,
        lensInfo: this.data.lensInfo,
        exifData: this.data.exifData,
        images: allImages
      };
      
      if (this.data.isEdit) {
        // 更新作品
        await updateWork(this.data.id, workData);
        wx.showToast({
          title: '更新成功',
          icon: 'success'
        });
      } else {
        // 创建新作品
        await uploadWork(workData);
        wx.showToast({
          title: '发布成功',
          icon: 'success'
        });
      }
      
      // 返回上一页
      setTimeout(() => {
        wx.navigateBack();
      }, 1500);
    } catch (error) {
      console.error('提交失败', error);
      wx.showToast({
        title: '提交失败，请重试',
        icon: 'none'
      });
      this.setData({ loading: false });
    }
  },

  onCancel() {
    wx.navigateBack();
  }
}); 