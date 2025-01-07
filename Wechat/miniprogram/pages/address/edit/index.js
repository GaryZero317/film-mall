const app = getApp()
const areaData = require('../../../utils/area-data')

Page({
  data: {
    address: {
      name: '',
      mobile: '',
      province: '',
      city: '',
      district: '',
      detail: '',
      is_default: false
    },
    region: [],
    regionValue: [0, 0, 0],
    provinces: [],
    cities: [],
    districts: [],
    showRegionPicker: false
  },

  onLoad(options) {
    // 初始化省市区数据
    this.initRegionData()
    
    // 如果有id参数，说明是编辑模式
    if (options.id) {
      this.loadAddress(options.id)
    }
  },

  // 初始化省市区数据
  initRegionData() {
    const provinces = areaData.provinces
    const cities = provinces[0].children
    const districts = cities[0].children

    this.setData({
      provinces,
      cities,
      districts
    })
  },

  // 加载地址详情
  async loadAddress(id) {
    try {
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/api/addresses/${id}`,
        method: 'GET',
        header: {
          'Authorization': `Bearer ${wx.getStorageSync('token')}`
        }
      })

      if (res.statusCode === 200) {
        const address = res.data
        this.setData({ 
          address,
          region: [address.province, address.city, address.district]
        })
        // 设置省市区选择器的值
        this.setRegionValue()
      }
    } catch (error) {
      console.error('加载地址详情失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    }
  },

  // 设置省市区选择器的值
  setRegionValue() {
    const { provinces, address } = this.data
    let provinceIndex = provinces.findIndex(item => item.name === address.province)
    if (provinceIndex === -1) provinceIndex = 0

    const cities = provinces[provinceIndex].children
    let cityIndex = cities.findIndex(item => item.name === address.city)
    if (cityIndex === -1) cityIndex = 0

    const districts = cities[cityIndex].children
    let districtIndex = districts.findIndex(item => item.name === address.district)
    if (districtIndex === -1) districtIndex = 0

    this.setData({
      cities,
      districts,
      regionValue: [provinceIndex, cityIndex, districtIndex]
    })
  },

  // 显示地区选择器
  onShowRegionPicker() {
    this.setData({ showRegionPicker: true })
  },

  // 隐藏地区选择器
  onHideRegionPicker() {
    this.setData({ showRegionPicker: false })
  },

  // 地区选择器值改变
  onRegionChange(e) {
    const value = e.detail.value
    const [provinceIndex, cityIndex] = value
    
    // 获取新的城市和区县数据
    const cities = this.data.provinces[provinceIndex].children
    const districts = cities[cityIndex].children

    this.setData({
      cities,
      districts,
      regionValue: value
    })
  },

  // 确认选择地区
  onConfirmRegion() {
    const { provinces, cities, districts, regionValue } = this.data
    const [provinceIndex, cityIndex, districtIndex] = regionValue
    
    const region = [
      provinces[provinceIndex].name,
      cities[cityIndex].name,
      districts[districtIndex].name
    ]

    this.setData({
      region,
      'address.province': region[0],
      'address.city': region[1],
      'address.district': region[2],
      showRegionPicker: false
    })
  },

  // 提交表单
  async onSubmit(e) {
    const formData = e.detail.value
    const { address, region } = this.data

    // 表单验证
    if (!formData.name) {
      wx.showToast({
        title: '请输入收货人姓名',
        icon: 'none'
      })
      return
    }

    if (!formData.mobile) {
      wx.showToast({
        title: '请输入手机号码',
        icon: 'none'
      })
      return
    }

    if (!/^1\d{10}$/.test(formData.mobile)) {
      wx.showToast({
        title: '手机号码格式不正确',
        icon: 'none'
      })
      return
    }

    if (region.length === 0) {
      wx.showToast({
        title: '请选择所在地区',
        icon: 'none'
      })
      return
    }

    if (!formData.detail) {
      wx.showToast({
        title: '请输入详细地址',
        icon: 'none'
      })
      return
    }

    // 构造请求数据
    const data = {
      name: formData.name,
      mobile: formData.mobile,
      province: region[0],
      city: region[1],
      district: region[2],
      detail: formData.detail,
      is_default: formData.is_default
    }

    try {
      const url = address.id 
        ? `${app.globalData.baseUrl}/api/addresses/${address.id}`
        : `${app.globalData.baseUrl}/api/addresses`
      const method = address.id ? 'PUT' : 'POST'

      const res = await wx.request({
        url,
        method,
        data,
        header: {
          'Authorization': `Bearer ${wx.getStorageSync('token')}`
        }
      })

      if (res.statusCode === 200) {
        wx.showToast({
          title: '保存成功',
          icon: 'success'
        })
        // 返回上一页
        wx.navigateBack()
      } else {
        wx.showToast({
          title: '保存失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('保存地址失败:', error)
      wx.showToast({
        title: '保存失败',
        icon: 'none'
      })
    }
  }
}) 