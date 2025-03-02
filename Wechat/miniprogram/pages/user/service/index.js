import { getFaqList, createQuestion } from '../../../api/service';

Page({
  data: {
    faqOpen: [],
    faqList: [],
    loading: false,
    questionType: 1, // 默认问题类型
    questionContent: '', // 问题内容
    showSubmitForm: false, // 是否显示提交表单
  },

  onLoad: function (options) {
    // 页面加载时的初始化工作
    this.loadFaqList();
  },
  
  // 加载常见问题列表
  loadFaqList: function() {
    this.setData({ loading: true });
    
    getFaqList({ page: 1, pageSize: 10 })
      .then(res => {
        console.log('获取FAQ返回数据:', res);
        if ((res.code === 0 || res.code === 200) && res.data) {
          // 兼容两种可能的数据格式
          let faqList = [];
          
          if (res.data.list) {
            // 格式: {data: {list: [...]}}
            faqList = res.data.list;
          } else if (Array.isArray(res.data)) {
            // 格式: {data: [...]}
            faqList = res.data;
          } else {
            // 其他可能的格式，将整个data作为列表
            faqList = [res.data];
          }
          
          // 初始化展开状态数组
          const faqOpen = new Array(faqList.length).fill(false);
          
          this.setData({
            faqList,
            faqOpen,
            loading: false
          });
          
          console.log('FAQ列表加载成功', faqList);
        } else {
          console.error('FAQ数据格式不正确', res);
          this.setData({ 
            loading: false,
            faqList: [] 
          });
        }
      })
      .catch(err => {
        console.error('加载FAQ失败', err);
        this.setData({ 
          loading: false,
          faqList: [] 
        });
        wx.showToast({
          title: '加载失败，请重试',
          icon: 'none'
        });
      });
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
  },
  
  // 显示问题提交表单
  showSubmitForm: function() {
    this.setData({ showSubmitForm: true });
  },
  
  // 隐藏问题提交表单
  hideSubmitForm: function() {
    this.setData({ showSubmitForm: false });
  },
  
  // 问题类型选择
  onTypeChange: function(e) {
    this.setData({
      questionType: parseInt(e.detail.value)
    });
  },
  
  // 问题内容输入
  onContentInput: function(e) {
    this.setData({
      questionContent: e.detail.value
    });
  },
  
  // 提交问题
  submitQuestion: function() {
    const { questionType, questionContent } = this.data;
    
    if (!questionContent.trim()) {
      wx.showToast({
        title: '请输入问题内容',
        icon: 'none'
      });
      return;
    }
    
    const data = {
      type: questionType,
      content: questionContent.trim()
    };
    
    wx.showLoading({ title: '提交中...' });
    
    createQuestion(data)
      .then(res => {
        wx.hideLoading();
        if (res.code === 0) {
          wx.showToast({
            title: '提交成功',
            icon: 'success'
          });
          // 重置表单
          this.setData({
            questionContent: '',
            showSubmitForm: false
          });
        } else {
          throw new Error(res.msg || '提交失败');
        }
      })
      .catch(err => {
        wx.hideLoading();
        console.error('提交问题失败', err);
        wx.showToast({
          title: err.message || '提交失败，请重试',
          icon: 'none'
        });
      });
  },
  
  // 跳转到聊天页面
  goToChat: function() {
    console.log('跳转到聊天页面');
    wx.navigateTo({
      url: './chat/index',
      success: function() {
        console.log('跳转成功');
      },
      fail: function(err) {
        console.error('跳转到聊天页面失败', err);
        // 尝试另一种路径
        wx.navigateTo({
          url: '/pages/user/service/chat/index',
          success: function() {
            console.log('使用绝对路径跳转成功');
          },
          fail: function(secondErr) {
            console.error('两种跳转方式均失败', secondErr);
            wx.showToast({
              title: '跳转失败，请重试',
              icon: 'none'
            });
          }
        });
      }
    });
  }
}); 