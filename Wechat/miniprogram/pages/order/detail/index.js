import { getOrderDetail } from '../../../api/order'
import { getAddressDetail } from '../../../api/address'

// 订单支付超时时间（毫秒），与后端保持一致，15分钟
const ORDER_PAYMENT_TIMEOUT = 15 * 60 * 1000;

// 格式化金额
const formatPrice = (price) => {
  if (typeof price !== 'number') return '0.00'
  return (price / 100).toFixed(2)
}

// 计算商品总价
const calculateTotalPrice = (items) => {
  if (!Array.isArray(items)) return 0
  return items.reduce((total, item) => {
    const price = item.price || 0
    const quantity = item.quantity || 0
    return total + (price * quantity)
  }, 0)
}

// 处理地址数据
const formatAddress = (addressData) => {
  if (!addressData) return null
  return {
    name: addressData.name || '',
    phone: addressData.phone || '',
    province: addressData.province || '',
    city: addressData.city || '',
    district: addressData.district || '',
    address: addressData.detailAddr || ''
  }
}

// 格式化日期字符串，确保iOS兼容性
const formatDateString = (dateStr) => {
  if (!dateStr) return ''
  // 将 "yyyy-MM-dd HH:mm:ss" 转换为 "yyyy/MM/dd HH:mm:ss"
  return dateStr.replace(/-/g, '/')
}

Page({
  data: {
    orderId: '',
    orderNo: '',
    createTime: '',
    orderStatus: 0,
    address: null,
    goods: [],
    totalPrice: '0.00',
    freight: '0.00',
    actualPrice: '0.00',
    countdown: null,  // 倒计时显示
    countdownTimer: null  // 倒计时定时器ID
  },

  onLoad(options) {
    console.log('订单详情页面参数:', options)
    
    let orderId = ''
    if (options.id) {
      orderId = options.id
      console.log('从options.id获取订单ID:', orderId)
    } else if (options.orderId) {
      orderId = options.orderId
      console.log('从options.orderId获取订单ID:', orderId)
    } else {
      console.error('没有提供订单ID参数!')
      wx.showToast({
        title: '订单ID不存在',
        icon: 'none'
      })
      // 延迟返回
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
      return
    }
    
    this.setData({
      orderId: orderId
    })
    
    // 确保orderId有效后再请求详情
    if (orderId) {
      console.log('开始请求订单详情，ID:', orderId)
      this.getOrderDetail()
    }
  },

  onHide() {
    // 页面隐藏时清除倒计时
    this.clearCountdownTimer()
  },

  onUnload() {
    // 页面卸载时清除倒计时
    this.clearCountdownTimer()
  },

  // 清除倒计时定时器
  clearCountdownTimer() {
    if (this.data.countdownTimer) {
      clearInterval(this.data.countdownTimer)
      this.setData({
        countdownTimer: null
      })
    }
  },

  async getOrderDetail() {
    try {
      console.log('开始获取订单详情，订单ID:', this.data.orderId)
      const res = await getOrderDetail(this.data.orderId)
      console.log('获取订单详情API响应:', res)
      
      if (res.code === 0 && res.data) {
        const { data } = res
        console.log('订单详情数据:', data)
        
        // 处理商品列表数据
        const items = Array.isArray(data.items) ? data.items.map(item => ({
          ...item,
          price: Number(item.price || 0),
          quantity: Number(item.quantity || item.num || 0)
        })) : []

        console.log('处理后的商品列表:', items)

        // 计算商品总价
        const totalPrice = calculateTotalPrice(items)
        
        // 订单信息设置
        this.setData({
          orderNo: data.oid || '',
          createTime: data.create_time || '',
          orderStatus: data.status || 0,
          goods: items,
          totalPrice: formatPrice(totalPrice),
          freight: formatPrice(data.shipping_fee || 0),
          actualPrice: formatPrice(data.total_price || 0)
        })
        
        // 如果有地址ID，获取地址信息
        if (data.address_id) {
          this.getAddressInfo(data.address_id)
        }
        
        // 如果是待支付状态，启动倒计时
        if (data.status === 0) {
          this.startOrderCountdown(data.create_time)
        }
      } else {
        console.error('获取订单详情失败:', res.msg || '未知错误')
        wx.showToast({
          title: '获取订单详情失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('获取订单详情异常:', error)
      wx.showToast({
        title: '获取订单详情失败',
        icon: 'none'
      })
    }
  },
  
  // 启动订单倒计时
  startOrderCountdown(createTimeStr) {
    // 清除可能已存在的倒计时
    this.clearCountdownTimer()
    
    if (!createTimeStr) {
      console.error('订单创建时间为空，无法启动倒计时')
      return
    }
    
    // 获取订单创建时间
    const createTime = new Date(formatDateString(createTimeStr)).getTime()
    
    // 计算过期时间
    const expireTime = createTime + ORDER_PAYMENT_TIMEOUT
    
    // 设置倒计时定时器
    const timerId = setInterval(() => {
      // 计算剩余时间
      const now = new Date().getTime()
      const remainTime = expireTime - now
      
      if (remainTime <= 0) {
        // 倒计时结束，清除计时器
        this.clearCountdownTimer()
        
        // 更新订单状态为已取消（UI显示）
        this.setData({
          orderStatus: 4, // 4表示已取消
          countdown: null
        })
        
        // 提示用户
        wx.showToast({
          title: '订单支付超时已自动取消',
          icon: 'none'
        })
        
        return
      }
      
      // 计算剩余分钟和秒数
      const minutes = Math.floor(remainTime / (60 * 1000))
      const seconds = Math.floor((remainTime % (60 * 1000)) / 1000)
      
      // 格式化倒计时文本
      const countdownText = `${minutes < 10 ? '0' + minutes : minutes}:${seconds < 10 ? '0' + seconds : seconds}`
      
      // 更新倒计时显示
      this.setData({
        countdown: countdownText
      })
    }, 1000)
    
    // 保存定时器ID
    this.setData({
      countdownTimer: timerId
    })
  },

  async getAddressInfo(addressId) {
    try {
      const res = await getAddressDetail(addressId)
      console.log('地址API响应:', res)
      if (res && res.address) {
        const address = formatAddress(res.address)
        console.log('格式化后的地址信息:', address)
        this.setData({ 
          address: {
            name: address.name,
            phone: address.phone,
            province: address.province,
            city: address.city,
            district: address.district,
            address: address.address
          }
        })
      } else {
        console.error('地址数据格式不正确:', res)
      }
    } catch (error) {
      console.error('获取地址信息失败:', error)
    }
  }
}) 