// 导入API函数，并确保路径正确
// 如果api导入有问题，我们可以先提供一个模拟接口
const apiService = {
  getChatSession: () => Promise.resolve({ code: 0, data: { sessionId: 'mock-session-id' } }),
  getChatHistory: () => Promise.resolve({ code: 0, data: { list: [] } }),
  sendChatMessage: () => Promise.resolve({ code: 0, data: { messageId: Date.now() } })
};

// 尝试正确导入API，如果失败则使用模拟接口
let getChatSession, getChatHistory, sendChatMessage;
try {
  const api = require('../../../../api/service');
  getChatSession = api.getChatSession || apiService.getChatSession;
  getChatHistory = api.getChatHistory || apiService.getChatHistory;
  sendChatMessage = api.sendChatMessage || apiService.sendChatMessage;
  console.log('API导入成功');
} catch (err) {
  console.error('API导入失败，使用模拟接口', err);
  getChatSession = apiService.getChatSession;
  getChatHistory = apiService.getChatHistory;
  sendChatMessage = apiService.sendChatMessage;
}

// 页面定义
Page({
  data: {
    // 基础数据
    sessionId: 'default-session-id',
    messages: [{
      id: Date.now(),
      content: '您好，我是客服助手，请问有什么可以帮您？',
      senderType: 2, // 客服消息
      createTime: Math.floor(Date.now() / 1000)
    }],
    inputValue: '',
    loading: false,
    sendingMessage: false,
    page: 1,
    pageSize: 20,
    hasMore: true,
    loadingMore: false
  },

  // 页面加载
  onLoad: function(options) {
    console.log('聊天页面加载成功');
    this.scrollToBottom();
  },

  // 返回上一页
  navigateBack: function() {
    wx.navigateBack({ delta: 1 });
  },

  // 输入框内容变化
  onInputChange: function(e) {
    this.setData({
      inputValue: e.detail.value
    });
  },

  // 发送消息
  sendMessage: function() {
    const content = this.data.inputValue.trim();
    if (!content) return;
    
    // 创建用户消息
    const userMessage = {
      id: Date.now(),
      content: content,
      senderType: 1, // 用户消息
      createTime: Math.floor(Date.now() / 1000)
    };
    
    // 更新UI
    this.setData({
      messages: [...this.data.messages, userMessage],
      inputValue: ''
    });
    
    // 滚动到底部
    this.scrollToBottom();
    
    // 模拟客服回复
    setTimeout(() => {
      const replyMessage = {
        id: Date.now(),
        content: '感谢您的咨询，我们会尽快处理您的问题。',
        senderType: 2, // 客服消息
        createTime: Math.floor(Date.now() / 1000)
      };
      
      this.setData({
        messages: [...this.data.messages, replyMessage]
      });
      
      // 滚动到底部
      this.scrollToBottom();
    }, 1000);
  },
  
  // 加载更多历史消息
  loadMoreHistory: function() {
    console.log('加载更多历史消息');
    // 这里可以调用API加载更多消息
  },
  
  // 滚动到底部
  scrollToBottom: function() {
    setTimeout(() => {
      wx.createSelectorQuery()
        .select('#message-container')
        .node()
        .exec(res => {
          if (res[0] && res[0].node) {
            const scrollView = res[0].node;
            scrollView.scrollIntoView({
              selector: '.message-item:last-child',
              block: 'end'
            });
          }
        });
    }, 100);
  }
}) 