Page({
  data: {
    favorites: []
  },

  onLoad() {
    this.getFavorites();
  },

  onShow() {
    this.getFavorites();
  },

  getFavorites() {
    wx.showLoading({
      title: '加载中...',
    });

    wx.cloud.callFunction({
      name: 'user',
      data: {
        action: 'getFavorites'
      }
    })
    .then(res => {
      const { data } = res.result;
      this.setData({
        favorites: data || []
      });
    })
    .catch(err => {
      console.error('获取收藏列表失败：', err);
      wx.showToast({
        title: '获取收藏列表失败',
        icon: 'none'
      });
    })
    .finally(() => {
      wx.hideLoading();
    });
  },

  goToProduct(e) {
    const { id } = e.currentTarget.dataset;
    wx.navigateTo({
      url: `/pages/product/detail/index?id=${id}`
    });
  },

  removeFavorite(e) {
    const { id } = e.currentTarget.dataset;
    
    wx.showModal({
      title: '提示',
      content: '确定要取消收藏该商品吗？',
      success: (res) => {
        if (res.confirm) {
          wx.showLoading({
            title: '处理中...',
          });

          wx.cloud.callFunction({
            name: 'user',
            data: {
              action: 'removeFavorite',
              productId: id
            }
          })
          .then(() => {
            wx.showToast({
              title: '已取消收藏',
              icon: 'success'
            });
            this.getFavorites();
          })
          .catch(err => {
            console.error('取消收藏失败：', err);
            wx.showToast({
              title: '取消收藏失败',
              icon: 'none'
            });
          })
          .finally(() => {
            wx.hideLoading();
          });
        }
      }
    });
  }
}); 