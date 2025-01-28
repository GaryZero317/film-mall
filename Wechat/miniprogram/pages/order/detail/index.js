import { getOrderDetail } from '../../../api/order'
import { getAddressDetail } from '../../../api/address'

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
    actualPrice: '0.00'
  },

  onLoad(options) {
    if (options.id) {
      this.setData({
        orderId: options.id
      })
      this.getOrderDetail()
    }
  },

  async getOrderDetail() {
    try {
      const res = await getOrderDetail(this.data.orderId)
      if (res.code === 0 && res.data) {
        const { data } = res
        console.log('订单详情数据:', data)
        
        // 处理商品列表数据
        const items = Array.isArray(data.items) ? data.items.map(item => ({
          ...item,
          price: Number(item.price || 0),
          quantity: Number(item.quantity || item.num || 0)
        })) : []

        // 计算商品总价
        const totalPrice = calculateTotalPrice(items)

        // 处理展示用的商品列表
        const goodsList = items.map(item => ({
          id: item.id,
          product_image: item.product_image || '',
          product_name: item.product_name || '',
          price: formatPrice(item.price),
          quantity: item.quantity
        }))

        const orderData = {
          orderNo: data.oid || '',
          createTime: data.create_time || '',
          orderStatus: typeof data.status === 'number' ? data.status : 0,
          goods: goodsList,
          totalPrice: formatPrice(totalPrice),
          freight: formatPrice(data.shipping_fee || 0),
          actualPrice: formatPrice(data.actual_price || totalPrice + (data.shipping_fee || 0))
        }

        this.setData(orderData)

        // 获取地址信息
        if (data.address_id) {
          this.getAddressInfo(data.address_id)
        }
      } else {
        wx.showToast({
          title: res.msg || '获取订单详情失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('获取订单详情失败:', error)
      wx.showToast({
        title: '获取订单详情失败',
        icon: 'none'
      })
    }
  },

  async getAddressInfo(addressId) {
    try {
      const res = await getAddressDetail(addressId)
      console.log('地址API响应:', res)
      if (res?.data?.address) {
        const address = formatAddress(res.data.address)
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