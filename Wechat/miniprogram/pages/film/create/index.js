import { getFilmPriceList, createFilmOrder } from '../../../api/film'
import { getAddressList } from '../../../api/address'
import { loginGuard } from '../../../utils/auth'

// 使用登录守卫包装Page
Page(loginGuard({
  /**
   * 页面的初始数据
   */
  data: {
    loading: true,
    submitting: false,
    priceList: {
      filmTypes: [],
      filmBrands: [],
      filmSizes: []
    },
    selectedAddress: null,
    totalAmount: 0,
    shippingFee: 15, // 运费默认15元
    form: {
      items: [
        {
          typeIndex: null,
          brandIndex: null,
          sizeIndex: null,
          quantity: 1,
          price: 0,
          remark: ''
        }
      ],
      returnFilm: false,
      remark: '',
      totalPrice: 0
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.loadPriceList()
  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {
    // 检查是否刚从地址选择页面返回
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 1]
    if (currentPage.data.selectedAddress) {
      // 如果已经有选择的地址，不再重新加载
      return
    }
    
    this.loadDefaultAddress()
  },

  /**
   * 加载价格列表
   */
  async loadPriceList() {
    this.setData({
      loading: true
    })

    try {
      const res = await getFilmPriceList()
      if (res && res.code === 0 && res.data) {
        this.setData({
          'priceList.filmTypes': res.data.types || [],
          'priceList.filmBrands': res.data.brands || [],
          'priceList.filmSizes': res.data.sizes || [],
          loading: false
        })
      } else {
        wx.showToast({
          title: '加载价格信息失败',
          icon: 'none'
        })
        this.setData({ loading: false })
        // 使用默认数据
        this.setDefaultPriceList()
      }
    } catch (error) {
      console.error('加载价格列表失败:', error)
      wx.showToast({
        title: '加载价格信息失败，将使用默认数据',
        icon: 'none'
      })
      this.setData({ loading: false })
      // 使用默认数据
      this.setDefaultPriceList()
    }
  },

  /**
   * 设置默认价格列表数据
   */
  setDefaultPriceList() {
    this.setData({
      'priceList.filmTypes': [
        { id: 1, name: '黑白胶片' },
        { id: 2, name: '彩色胶片' },
        { id: 3, name: '正片' }
      ],
      'priceList.filmBrands': [
        { id: 1, name: 'Kodak' },
        { id: 2, name: 'Fuji' },
        { id: 3, name: 'Ilford' },
        { id: 4, name: 'Lomography' }
      ],
      'priceList.filmSizes': [
        { id: 1, name: '135 (35mm)' },
        { id: 2, name: '120 (中画幅)' }
      ]
    })
  },

  /**
   * 加载默认地址
   */
  async loadDefaultAddress() {
    try {
      const res = await getAddressList()
      console.log('[胶片冲洗] 获取地址列表响应:', res)
      
      if (res && res.code === 200 && res.data && res.data.list) {
        // 获取默认地址
        const defaultAddress = res.data.list.find(item => item.isDefault) || res.data.list[0]
        console.log('[胶片冲洗] 选择的默认地址:', defaultAddress)
        this.setData({ selectedAddress: defaultAddress })
      } else {
        console.error('[胶片冲洗] 获取地址列表失败, 响应数据无效:', res)
      }
    } catch (error) {
      console.error('加载地址失败:', error)
      // 不显示错误提示，因为用户可以手动选择地址
    }
  },

  /**
   * 胶片类型变更
   */
  onFilmTypeChange(e) {
    const index = e.currentTarget.dataset.index
    const typeIndex = e.detail.value
    const key = `form.items[${index}].typeIndex`
    
    this.setData({
      [key]: typeIndex
    })
    
    this.updateItemPrice(index)
  },

  /**
   * 胶片品牌变更
   */
  onFilmBrandChange(e) {
    const index = e.currentTarget.dataset.index
    const brandIndex = e.detail.value
    const key = `form.items[${index}].brandIndex`
    
    this.setData({
      [key]: brandIndex
    })
    
    this.updateItemPrice(index)
  },

  /**
   * 胶片尺寸变更
   */
  onFilmSizeChange(e) {
    const index = e.currentTarget.dataset.index
    const sizeIndex = e.detail.value
    const key = `form.items[${index}].sizeIndex`
    
    this.setData({
      [key]: sizeIndex
    })
    
    this.updateItemPrice(index)
  },

  /**
   * 数量变更
   */
  onQuantityChange(e) {
    const index = e.currentTarget.dataset.index
    let quantity = parseInt(e.detail.value)
    
    if (isNaN(quantity) || quantity < 1) {
      quantity = 1
      wx.showToast({
        title: '数量不能小于1',
        icon: 'none'
      })
    }
    
    if (quantity > 10) {
      quantity = 10
      wx.showToast({
        title: '单次最多可冲洗10卷胶片',
        icon: 'none'
      })
    }
    
    const key = `form.items[${index}].quantity`
    
    this.setData({
      [key]: quantity
    })
    
    this.calculateTotal()
    this.updateTotalAmount()
  },

  /**
   * 胶片备注变更
   */
  onItemRemarkChange(e) {
    const index = e.currentTarget.dataset.index
    const remark = e.detail.value
    const key = `form.items[${index}].remark`
    
    this.setData({
      [key]: remark
    })
  },

  /**
   * 订单备注变更
   */
  onRemarkChange(e) {
    this.setData({
      'form.remark': e.detail.value
    })
  },

  /**
   * 切换回寄底片
   */
  toggleReturnFilm(e) {
    this.setData({
      'form.returnFilm': e.detail.value
    })
    
    this.updateTotalAmount()
  },

  /**
   * 添加胶片
   */
  addFilmItem() {
    const items = this.data.form.items
    
    if (items.length >= 5) {
      wx.showToast({
        title: '最多可添加5卷胶片',
        icon: 'none'
      })
      return
    }
    
    items.push({
      typeIndex: null,
      brandIndex: null,
      sizeIndex: null,
      quantity: 1,
      price: 0,
      remark: ''
    })
    
    this.setData({
      'form.items': items
    })
    
    this.calculateTotal()
    this.updateTotalAmount()
  },

  /**
   * 移除胶片
   */
  removeFilmItem(e) {
    const index = e.currentTarget.dataset.index
    const items = this.data.form.items
    
    if (items.length <= 1) {
      wx.showToast({
        title: '至少需要保留一卷胶片',
        icon: 'none'
      })
      return
    }
    
    items.splice(index, 1)
    
    this.setData({
      'form.items': items
    })
    
    this.calculateTotal()
    this.updateTotalAmount()
  },

  /**
   * 更新物品价格
   */
  updateItemPrice(index) {
    const item = this.data.form.items[index]
    const { typeIndex, brandIndex, sizeIndex } = item
    
    // 检查是否已经选择了所有必要的选项
    if (typeIndex === null || brandIndex === null || sizeIndex === null) {
      return
    }
    
    const filmType = this.data.priceList.filmTypes[typeIndex]
    const filmBrand = this.data.priceList.filmBrands[brandIndex]
    const filmSize = this.data.priceList.filmSizes[sizeIndex]
    
    if (!filmType || !filmBrand || !filmSize) {
      return
    }
    
    // 查找价格
    let price = 0
    
    // 实际项目中应该根据API返回的价格数据计算
    // 这里简单设置一个默认价格
    price = 35 // 默认价格35元
    
    // 设置价格
    const priceKey = `form.items[${index}].price`
    
    this.setData({
      [priceKey]: price
    })
    
    this.calculateTotal()
    this.updateTotalAmount()
  },

  /**
   * 计算总价
   */
  calculateTotal() {
    const items = this.data.form.items
    let total = 0
    
    for (const item of items) {
      total += item.price * item.quantity
    }
    
    this.setData({
      'form.totalPrice': total
    })
    
    return total
  },

  /**
   * 更新总金额（包含运费）
   */
  updateTotalAmount() {
    const totalPrice = this.data.form.totalPrice
    const shippingFee = this.data.form.returnFilm ? this.data.shippingFee : 0
    
    this.setData({
      totalAmount: totalPrice + shippingFee
    })
  },

  /**
   * 选择地址
   */
  onSelectAddress() {
    wx.navigateTo({
      url: '/pages/address/list/index?select=true'
    })
  },

  /**
   * 提交订单
   */
  async submitOrder() {
    // 表单验证
    if (!this.validateForm()) {
      return
    }
    
    this.setData({
      submitting: true
    })
    
    try {
      const formData = this.prepareOrderData()
      
      const res = await createFilmOrder(formData)
      console.log('[提交胶片订单] 响应数据:', res) // 添加日志查看响应数据结构
      
      // 添加更详细的响应数据结构日志
      console.log('[提交胶片订单] 响应状态码:', res?.code)
      console.log('[提交胶片订单] 响应消息:', res?.message || res?.msg)
      console.log('[提交胶片订单] 响应数据详情:', res?.data)
      if (res?.data) {
        console.log('[提交胶片订单] 数据字段列表:', Object.keys(res.data))
      }
      
      if (res && res.code === 0 && res.data) {
        wx.showToast({
          title: '订单提交成功',
          icon: 'success'
        })
        
        // 获取订单ID，兼容多种可能的字段名
        const orderId = res.data.orderId || res.data.id || res.data.order_id || res.data.orderid
        console.log('[提交胶片订单] 获取到的订单ID:', orderId)
        
        if (!orderId) {
          console.error('[提交胶片订单] 未获取到有效的订单ID:', res.data)
          wx.showToast({
            title: '创建订单成功，但未获取到订单ID',
            icon: 'none'
          })
          this.setData({ submitting: false })
          return
        }
        
        // 延迟跳转到支付页面
        setTimeout(() => {
          wx.redirectTo({
            url: `/pages/order/payment/index?orderId=${orderId}&type=film`
          })
        }, 1500)
      } else {
        wx.showToast({
          title: res.message || '订单提交失败',
          icon: 'none'
        })
        this.setData({ submitting: false })
      }
    } catch (error) {
      console.error('提交订单失败:', error)
      wx.showToast({
        title: '订单提交失败，请重试',
        icon: 'none'
      })
      this.setData({ submitting: false })
    }
  },

  /**
   * 准备订单数据
   */
  prepareOrderData() {
    console.log('[准备订单数据] 开始准备提交信息:', this.data)
    
    // 获取表单数据
    const formData = this.data
    console.log('[准备订单数据] 原始表单数据:', formData)
    
    // 用户ID - 从userInfo对象中获取
    const userInfo = wx.getStorageSync('userInfo')
    const uid = userInfo ? userInfo.id : 0
    console.log('[准备订单数据] 用户信息:', userInfo)
    console.log('[准备订单数据] 用户ID:', uid)
    
    // 准备订单项数据
    const items = formData.form.items.map((item, index) => {
      // 获取类型、品牌、尺寸信息
      const filmType = this.data.priceList.filmTypes[item.typeIndex]
      const filmBrand = this.data.priceList.filmBrands[item.brandIndex]
      const filmSize = this.data.priceList.filmSizes[item.sizeIndex]
      
      console.log(`[准备订单数据] 处理第${index+1}个订单项:`, {
        索引: index,
        胶片类型: filmType,
        胶片品牌: filmBrand,
        胶片尺寸: filmSize,
        数量: item.quantity,
        价格: item.price
      })
      
      // 确保所有必要字段都存在
      if (!filmType || !filmBrand || !filmSize) {
        console.error(`[准备订单数据] 第${index+1}个订单项缺少必要数据:`, {
          胶片类型: !!filmType,
          胶片品牌: !!filmBrand,
          胶片尺寸: !!filmSize
        })
        
        // 提供默认值
        const defaultType = filmType || { id: 0, name: '胶片' }
        const defaultBrand = filmBrand || { id: 0, name: '标准' }
        const defaultSize = filmSize || { id: 0, name: '标准尺寸' }
        
        return {
          film_type: defaultType.name,
          film_brand: defaultBrand.name,
          size: defaultSize.name,
          quantity: item.quantity,
          price: Math.round(item.price * 100), // 转换为分
          amount: Math.round(item.price * item.quantity * 100), // 转换为分
          remark: item.remark
        }
      }
      
      return {
        film_type: filmType.name,
        film_brand: filmBrand.name,
        size: filmSize.name,
        quantity: item.quantity,
        price: Math.round(item.price * 100), // 转换为分
        amount: Math.round(item.price * item.quantity * 100), // 转换为分
        remark: item.remark
      }
    })
    
    console.log('[准备订单数据] 订单项处理完成:', items)
    
    // 计算总价（以分为单位）
    const totalPrice = items.reduce((sum, item) => sum + item.amount, 0)
    
    // 准备完整的订单数据
    const orderData = {
      uid: uid,
      address_id: formData.selectedAddress?.id || 0,
      return_film: !!formData.form.returnFilm,
      shipping_fee: Math.round(formData.shippingFee * 100), // 转换为分
      total_price: totalPrice,
      total_amount: totalPrice + Math.round(formData.shippingFee * 100), // 总价加运费
      remark: formData.form.remark || '',
      status: 0, // 添加状态字段，0表示待付款
      items: items
    }
    
    console.log('[准备订单数据] 最终订单数据:', orderData)
    return orderData
  },

  /**
   * 验证表单
   */
  validateForm() {
    // 验证地址
    if (!this.data.selectedAddress) {
      wx.showToast({
        title: '请选择收货地址',
        icon: 'none'
      })
      return false
    }
    
    // 验证每个胶片项
    const items = this.data.form.items
    for (let i = 0; i < items.length; i++) {
      const item = items[i]
      
      // 验证胶片类型
      if (item.typeIndex === null) {
        wx.showToast({
          title: `胶片${i + 1}未选择胶片类型`,
          icon: 'none'
        })
        return false
      }
      
      // 验证胶片品牌
      if (item.brandIndex === null) {
        wx.showToast({
          title: `胶片${i + 1}未选择胶片品牌`,
          icon: 'none'
        })
        return false
      }
      
      // 验证胶片尺寸
      if (item.sizeIndex === null) {
        wx.showToast({
          title: `胶片${i + 1}未选择胶片尺寸`,
          icon: 'none'
        })
        return false
      }
      
      // 验证胶片数量
      if (!item.quantity || item.quantity < 1) {
        wx.showToast({
          title: `胶片${i + 1}数量不能小于1`,
          icon: 'none'
        })
        return false
      }
    }
    
    return true
  },

  /**
   * 返回上一页
   */
  onBack() {
    wx.navigateBack()
  }
})) 