const app = getApp()
const { getAddressList, setDefaultAddress, deleteAddress } = require('../../../api/address')

console.log('[地址列表] baseUrl:', app.globalData.baseUrl)

Page({
  data: {
    addresses: [],
    isSelectMode: false // 是否是选择地址模式
  },

  onLoad(options) {
    console.log('[地址列表] onLoad options:', options)
    // 如果是从确认订单页面跳转来的，设置为选择模式
    if (options.select) {
      this.setData({ isSelectMode: true })
    }
  },

  onShow() {
    console.log('[地址列表] onShow')
    this.loadAddresses()
  },

  // 加载地址列表
  async loadAddresses() {
    console.log('[地址列表] 开始加载地址列表')
    try {
      const token = wx.getStorageSync('token')
      console.log('[地址列表] token:', token)
      
      const res = await getAddressList()
      console.log('[地址列表] 获取地址列表响应:', res)

      if (res && res.data) {
        console.log('[地址列表] 地址列表数据:', res.data)
        this.setData({ addresses: res.data.list || [] })
      } else {
        console.error('[地址列表] 获取地址列表失败:', res)
        wx.showToast({
          title: '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[地址列表] 加载地址列表失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    }
  },

  // 选择地址
  onSelectAddress(e) {
    console.log('[地址列表] 选择地址:', e.currentTarget.dataset)
    if (!this.data.isSelectMode) return
    
    const id = e.currentTarget.dataset.id
    const address = this.data.addresses.find(item => item.id === id)
    console.log('[地址列表] 选中的地址:', address)
    
    // 返回到上一页并传递选中的地址
    const pages = getCurrentPages()
    const prevPage = pages[pages.length - 2]
    prevPage.setData({ selectedAddress: address })
    wx.navigateBack()
  },

  // 设为默认地址
  async onSetDefault(e) {
    const id = e.currentTarget.dataset.id
    console.log('[地址列表] 设置默认地址, id:', id)
    
    try {
      const res = await setDefaultAddress(id)
      console.log('[地址列表] 设置默认地址响应:', res)

      // 判断是否成功：支持code为0或200的情况
      if (res.code === 0 || res.code === 200) {
        console.log('[地址列表] 设置默认地址成功:', res)
        wx.showToast({
          title: '设置成功',
          icon: 'success'
        })
        this.loadAddresses()
      } else {
        console.error('[地址列表] 设置默认地址失败:', res)
        wx.showToast({
          title: res?.message || res?.msg || '设置失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[地址列表] 设置默认地址失败:', error)
      wx.showToast({
        title: '设置失败',
        icon: 'none'
      })
    }
  },

  // 编辑地址
  onEdit(e) {
    const id = e.currentTarget.dataset.id
    console.log('[地址列表] 编辑地址, id:', id)
    wx.navigateTo({
      url: `/pages/address/edit/index?id=${id}`
    })
  },

  // 删除地址
  async onDelete(e) {
    const id = e.currentTarget.dataset.id
    console.log('[地址列表] 准备删除地址, id:', id)
    
    const res = await wx.showModal({
      title: '提示',
      content: '确定要删除该地址吗？',
      confirmText: '删除',
      cancelText: '取消'
    })

    if (res.confirm) {
      try {
        console.log('[地址列表] 开始删除地址')
        const result = await deleteAddress(id)
        console.log('[地址列表] 删除地址响应:', result)

        // 判断删除是否成功：支持code为0或200的情况
        if (result && (result.code === 0 || result.code === 200)) {
          console.log('[地址列表] 删除地址成功:', result)
          wx.showToast({
            title: '删除成功',
            icon: 'success'
          })
          this.loadAddresses()
        } else {
          console.error('[地址列表] 删除地址失败:', result)
          wx.showToast({
            title: result?.msg || '删除失败',
            icon: 'none'
          })
        }
      } catch (error) {
        console.error('[地址列表] 删除地址失败:', error)
        wx.showToast({
          title: '删除失败',
          icon: 'none'
        })
      }
    }
  },

  // 新增地址
  onAdd() {
    console.log('[地址列表] 跳转到新增地址页面')
    wx.navigateTo({
      url: '/pages/address/edit/index'
    })
  }
}) 