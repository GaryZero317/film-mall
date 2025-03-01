Page({
  data: {
    faqOpen: [false, false, false],
  },

  onLoad: function (options) {
    // 页面加载时的初始化工作
  },

  // 切换FAQ问题的展开状态
  toggleFaq: function (e) {
    const index = e.currentTarget.dataset.index;
    const faqOpen = [...this.data.faqOpen];
    faqOpen[index] = !faqOpen[index];
    this.setData({ faqOpen });
  },

  // 拨打客服电话
  callService: function () {
    wx.makePhoneCall({
      phoneNumber: '400-123-4567',
      success: function () {
        console.log('拨打电话成功');
      },
      fail: function (err) {
        console.log('拨打电话失败:', err);
      }
    });
  }
}); 