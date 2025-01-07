Page({
  data: {
    orderNo: '',
    createTime: '',
    goods: [],
    totalPrice: 0,
    freight: 0,
    actualPrice: 0,
    address: {},
    orderStatus: ''
  },

  onLoad(options) {
    if (options.id) {
      this.getOrderDetail(options.id);
    }
  },

  getOrderDetail(orderId) {
    wx.showLoading({
      title: '加载中...',
    });

    // 这里需要调用获取订单详情的API
    wx.cloud.callFunction({
      name: 'order',
      data: {
        action: 'getOrderDetail',
        orderId: orderId
      }
    })
    .then(res => {
      const { data } = res.result;
      if (data) {
        this.setData({
          orderNo: data.orderNo,
          createTime: data.createTime,
          goods: data.goods,
          totalPrice: data.totalPrice,
          freight: data.freight,
          actualPrice: data.actualPrice,
          address: data.address,
          orderStatus: data.orderStatus
        });
      }
    })
    .catch(err => {
      console.error('获取订单详情失败：', err);
      wx.showToast({
        title: '获取订单详情失败',
        icon: 'none'
      });
    })
    .finally(() => {
      wx.hideLoading();
    });
  }
}); 