const app = getApp()

Page({
  data: {
    addresses: [],
    isSelectMode: false // 是否是选择地址模式
  },

  onLoad(options) {
    // 如果是从确认订单页面跳转来的，设置为选择模式
    if (options.select) {
      this.setData({ isSelectMode: true })
    }
  },

  onShow() {
    this.loadAddresses()
  },

  // 加载地址列表
  async loadAddresses() {
    try {
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/addresses`,
        method: 'GET',
        header: {
          'Authorization': `Bearer ${wx.getStorageSync('token')}`
        }
      })

      if (res.statusCode === 200) {
        this.setData({ addresses: res.data })
      }
    } catch (error) {
      console.error('加载地址列表失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    }
  },

  // 选择地址
  onSelectAddress(e) {
    if (!this.data.isSelectMode) return
    
    const id = e.currentTarget.dataset.id
    const address = this.data.addresses.find(item => item.id === id)
    
    // 返回到上一页并传递选中的地址
    const pages = getCurrentPages()
    const prevPage = pages[pages.length - 2]
    prevPage.setData({ selectedAddress: address })
    wx.navigateBack()
  },

  // 设为默认地址
  async onSetDefault(e) {
    const id = e.currentTarget.dataset.id
    
    try {
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/addresses/${id}/default`,
        method: 'POST',
        header: {
          'Authorization': `Bearer ${wx.getStorageSync('token')}`
        }
      })

      if (res.statusCode === 200) {
        wx.showToast({
          title: '设置成功',
          icon: 'success'
        })
        this.loadAddresses()
      } else {
        wx.showToast({
          title: '设置失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('设置默认地址失败:', error)
      wx.showToast({
        title: '设置失败',
        icon: 'none'
      })
    }
  },

  // 编辑地址
  onEdit(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/address/edit/index?id=${id}`
    })
  },

  // 删除地址
  async onDelete(e) {
    const id = e.currentTarget.dataset.id
    
    const res = await wx.showModal({
      title: '提示',
      content: '确定要删除该地址吗？',
      confirmText: '删除',
      cancelText: '取消'
    })

    if (res.confirm) {
      try {
        const result = await wx.request({
          url: `${app.globalData.baseUrl}/api/addresses/${id}`,
          method: 'DELETE',
          header: {
            'Authorization': `Bearer ${wx.getStorageSync('token')}`
          }
        })

        if (result.statusCode === 200) {
          wx.showToast({
            title: '删除成功',
            icon: 'success'
          })
          this.loadAddresses()
        } else {
          wx.showToast({
            title: '删除失败',
            icon: 'none'
          })
        }
      } catch (error) {
        console.error('删除地址失败:', error)
        wx.showToast({
          title: '删除失败',
          icon: 'none'
        })
      }
    }
  },

  // 新增地址
  onAdd() {
    wx.navigateTo({
      url: '/pages/address/edit/index'
    })
  }
}) 