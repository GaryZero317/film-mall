const app = getApp()
const areaData = require('../../../utils/area-data')
const { getAddressList, createAddress, updateAddress } = require('../../../api/address')

console.log('[地址编辑] baseUrl:', app.globalData.baseUrl)

Page({
  data: {
    address: {
      name: '',
      phone: '',
      province: '',
      city: '',
      district: '',
      detailAddr: '',
      isDefault: false
    },
    region: [],
    regionValue: [0, 0, 0],
    provinces: [],
    cities: [],
    districts: [],
    showRegionPicker: false,
    addressId: null
  },

  onLoad(options) {
    console.log('[地址编辑] onLoad options:', options)
    // 初始化省市区数据
    this.initRegionData()
    
    // 如果有id参数，说明是编辑模式
    if (options.id) {
      this.setData({ addressId: parseInt(options.id) })
      this.loadAddress(parseInt(options.id))
    }
  },

  // 初始化省市区数据
  initRegionData() {
    console.log('[地址编辑] 初始化省市区数据')
    const provinces = areaData.provinces || []
    if (provinces.length === 0) {
      console.error('[地址编辑] 省份数据为空')
      return
    }

    const cities = provinces[0].children || []
    if (cities.length === 0) {
      console.error('[地址编辑] 城市数据为空')
      return
    }

    const districts = cities[0].children || []
    if (districts.length === 0) {
      console.error('[地址编辑] 区县数据为空')
      return
    }

    console.log('[地址编辑] 初始化地区数据:', {
      provinces: provinces.map(p => p.name),
      cities: cities.map(c => c.name),
      districts: districts.map(d => d.name)
    })

    this.setData({
      provinces,
      cities,
      districts
    })
  },

  // 加载地址详情
  async loadAddress(id) {
    console.log('[地址编辑] 开始加载地址详情, id:', id)
    try {
      const res = await getAddressList()
      console.log('[地址编辑] 获取地址列表响应:', res)

      if (res && res.data) {
        const address = res.data.list.find(item => item.id === id)
        console.log('[地址编辑] 找到的地址详情:', address)
        if (address) {
          this.setData({ 
            address: {
              name: address.name || '',
              phone: address.phone || '',
              province: address.province || '',
              city: address.city || '',
              district: address.district || '',
              detailAddr: address.detailAddr || '',
              isDefault: address.isDefault || false
            },
            region: [
              address.province || '',
              address.city || '',
              address.district || ''
            ]
          })
          // 设置省市区选择器的值
          this.setRegionValue()
        } else {
          console.error('[地址编辑] 未找到对应的地址')
          wx.showToast({
            title: '地址不存在',
            icon: 'none'
          })
          setTimeout(() => {
            wx.navigateBack()
          }, 1500)
        }
      } else {
        console.error('[地址编辑] 获取地址列表失败:', res)
        wx.showToast({
          title: '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('[地址编辑] 加载地址详情失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    }
  },

  // 设置省市区选择器的值
  setRegionValue() {
    console.log('[地址编辑] 设置省市区选择器的值')
    const { provinces, address } = this.data
    if (!provinces || provinces.length === 0) {
      console.error('[地址编辑] 省份数据为空')
      return
    }

    // 查找省份索引
    let provinceIndex = provinces.findIndex(item => item.name === address.province)
    if (provinceIndex === -1) {
      console.log('[地址编辑] 未找到对应的省份，使用默认值')
      provinceIndex = 0
    }

    // 获取城市列表
    const cities = provinces[provinceIndex].children || []
    if (cities.length === 0) {
      console.error('[地址编辑] 城市数据为空')
      return
    }

    // 查找城市索引
    let cityIndex = cities.findIndex(item => item.name === address.city)
    if (cityIndex === -1) {
      console.log('[地址编辑] 未找到对应的城市，使用默认值')
      cityIndex = 0
    }

    // 获取区县列表
    const districts = cities[cityIndex].children || []
    if (districts.length === 0) {
      console.error('[地址编辑] 区县数据为空')
      return
    }

    // 查找区县索引
    let districtIndex = districts.findIndex(item => item.name === address.district)
    if (districtIndex === -1) {
      console.log('[地址编辑] 未找到对应的区县，使用默认值')
      districtIndex = 0
    }

    console.log('[地址编辑] 省市区索引:', { 
      provinceIndex, 
      cityIndex, 
      districtIndex,
      province: provinces[provinceIndex].name,
      city: cities[cityIndex].name,
      district: districts[districtIndex].name
    })

    this.setData({
      cities,
      districts,
      regionValue: [provinceIndex, cityIndex, districtIndex]
    })
  },

  // 显示地区选择器
  onShowRegionPicker() {
    console.log('[地址编辑] 显示地区选择器')
    this.setData({ showRegionPicker: true })
  },

  // 隐藏地区选择器
  onHideRegionPicker() {
    console.log('[地址编辑] 隐藏地区选择器')
    this.setData({ showRegionPicker: false })
  },

  // 地区选择器值改变
  onRegionChange(e) {
    console.log('[地址编辑] 地区选择器值改变:', e.detail.value)
    const value = e.detail.value
    const [provinceIndex, cityIndex] = value
    
    // 获取新的城市和区县数据
    const { provinces } = this.data
    if (!provinces || !provinces[provinceIndex]) {
      console.error('[地址编辑] 省份数据不存在:', { provinces, provinceIndex })
      return
    }

    const cities = provinces[provinceIndex].children || []
    if (!cities[cityIndex]) {
      console.error('[地址编辑] 城市数据不存在:', { cities, cityIndex })
      return
    }

    const districts = cities[cityIndex].children || []

    console.log('[地址编辑] 地区数据:', {
      provinces: provinces[provinceIndex].name,
      cities: cities[cityIndex].name,
      districts: districts.map(d => d.name)
    })

    this.setData({
      cities,
      districts,
      regionValue: value
    })
  },

  // 确认选择地区
  onConfirmRegion() {
    console.log('[地址编辑] 确认选择地区')
    const { provinces, cities, districts, regionValue } = this.data
    const [provinceIndex, cityIndex, districtIndex] = regionValue
    
    const region = [
      provinces[provinceIndex].name,
      cities[cityIndex].name,
      districts[districtIndex].name
    ]

    console.log('[地址编辑] 选择的地区:', region)

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
    console.log('[地址编辑] 提交表单数据:', e.detail.value)
    const formData = e.detail.value
    const { region, addressId } = this.data

    // 表单验证
    if (!formData.name) {
      wx.showToast({
        title: '请输入收货人姓名',
        icon: 'none'
      })
      return
    }

    if (!formData.phone) {
      wx.showToast({
        title: '请输入手机号码',
        icon: 'none'
      })
      return
    }

    if (!/^1\d{10}$/.test(formData.phone)) {
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

    if (!formData.detailAddr) {
      wx.showToast({
        title: '请输入详细地址',
        icon: 'none'
      })
      return
    }

    // 构造请求数据
    const data = {
      name: formData.name,
      phone: formData.phone,
      province: region[0],
      city: region[1],
      district: region[2],
      detailAddr: formData.detailAddr,
      isDefault: formData.isDefault || false
    }

    console.log('[地址编辑] 准备提交的数据:', data)

    try {
      let res
      if (addressId) {
        console.log('[地址编辑] 更新地址')
        res = await updateAddress(addressId, data)
        console.log('[地址编辑] 更新地址响应:', res)
        // 更新地址返回 code
        if (res && res.data && (res.data.code === 0 || res.data.code === '0')) {
          wx.showToast({
            title: '保存成功',
            icon: 'success'
          })
          setTimeout(() => {
            wx.navigateBack()
          }, 1500)
        } else {
          console.error('[地址编辑] 更新地址失败:', res)
          wx.showToast({
            title: res.data?.message || '保存失败',
            icon: 'none'
          })
        }
      } else {
        console.log('[地址编辑] 新增地址')
        res = await createAddress(data)
        console.log('[地址编辑] 新增地址响应:', res)
        // 新增地址返回 id
        if (res && res.data && res.data.id) {
          wx.showToast({
            title: '保存成功',
            icon: 'success'
          })
          setTimeout(() => {
            wx.navigateBack()
          }, 1500)
        } else {
          console.error('[地址编辑] 新增地址失败:', res)
          wx.showToast({
            title: '保存失败',
            icon: 'none'
          })
        }
      }
    } catch (error) {
      console.error('[地址编辑] 保存地址失败:', error)
      wx.showToast({
        title: '保存失败',
        icon: 'none'
      })
    }
  }
}) 