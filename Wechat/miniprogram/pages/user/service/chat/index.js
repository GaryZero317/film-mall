// 导入API函数
const service = require('../../../../api/service');
const { getChatSession, getChatHistory, sendChatMessage } = service;

// 页面定义
Page({
  data: {
    // 基础数据
    messages: [],
    inputValue: '',
    loading: true,
    sendingMessage: false,
    sessionId: '',
    page: 1,
    pageSize: 20,
    hasMore: true,
    loadingMore: false
  },

  // 页面加载
  onLoad: function(options) {
    console.log('聊天页面加载成功');
    this.initChatSession();
  },
  
  // 初始化聊天会话
  initChatSession: function() {
    this.setData({ loading: true });
    
    // 获取或创建聊天会话
    getChatSession().then(res => {
      console.log('获取聊天会话结果:', res);
      if (res.code === 0 && res.data) {
        const sessionId = res.data.sessionId;
        
        this.setData({
          sessionId: sessionId,
          loading: false
        });
        
        // 加载聊天历史记录
        this.loadChatHistory();
      } else {
        wx.showToast({
          title: res.msg || '初始化聊天失败',
          icon: 'none'
        });
        this.setData({ loading: false });
      }
    }).catch(err => {
      console.error('初始化聊天会话失败', err);
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      this.setData({ loading: false });
    });
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
    
    const token = wx.getStorageSync('token');
    if (!token) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    if (!this.data.sessionId) {
      wx.showToast({
        title: '聊天会话未初始化',
        icon: 'none'
      });
      return;
    }
    
    // 创建临时用户消息(乐观UI更新)
    const tempMessage = {
      id: `temp-${Date.now()}`,
      content: content,
      senderType: 1, // 用户消息
      createTime: Math.floor(Date.now() / 1000),
      isTemp: true // 标记为临时消息
    };
    
    // 更新UI，清空输入框
    this.setData({
      messages: [...this.data.messages, tempMessage],
      inputValue: '',
      sendingMessage: true
    });
    
    // 滚动到底部
    this.scrollToBottom();
    
    // 调用API发送消息
    sendChatMessage(this.data.sessionId, content, 1).then(res => {
      console.log('消息发送成功', res);
      if (res.code === 0) {
        // 用真实消息替换临时消息
        const realMessage = {
          id: res.data.id || Date.now(),
          content: content,
          senderType: 1,
          createTime: res.data.createTime || Math.floor(Date.now() / 1000)
        };
        
        // 替换临时消息
        const updatedMessages = this.data.messages.map(msg => 
          msg.id === tempMessage.id ? realMessage : msg
        );
        
        this.setData({ messages: updatedMessages });
        
        // 如果有自动回复消息，直接添加到列表中
        if (res.data.serviceReply) {
          const serviceMessage = {
            id: res.data.serviceReply.id,
            content: res.data.serviceReply.content,
            senderType: 2, // 客服消息
            createTime: res.data.serviceReply.createTime
          };
          
          // 延迟显示客服回复，模拟真实场景
          setTimeout(() => {
            this.setData({
              messages: [...this.data.messages, serviceMessage]
            });
            this.scrollToBottom();
          }, 1000);
        } else {
          // 短暂延迟后刷新消息列表以获取可能的自动回复
          setTimeout(() => {
            this.refreshMessages();
          }, 1000);
        }
      } else {
        // 标记为发送失败
        const failedMessages = this.data.messages.map(msg => {
          if (msg.id === tempMessage.id) {
            return { ...msg, isTemp: false, isFailed: true };
          }
          return msg;
        });
        
        this.setData({ messages: failedMessages });
        
        wx.showToast({
          title: res.msg || '发送失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      console.error('消息发送请求失败', err);
      // 标记为发送失败
      const failedMessages = this.data.messages.map(msg => {
        if (msg.id === tempMessage.id) {
          return { ...msg, isTemp: false, isFailed: true };
        }
        return msg;
      });
      
      this.setData({ messages: failedMessages });
      
      wx.showToast({
        title: '网络错误，请稍后重试',
        icon: 'none'
      });
    }).finally(() => {
      this.setData({ sendingMessage: false });
    });
  },
  
  // 刷新消息列表
  refreshMessages: function() {
    // 如果没有会话ID，不执行操作
    if (!this.data.sessionId) return;
    
    // 获取最新消息
    getChatHistory(this.data.sessionId, 1, this.data.pageSize).then(res => {
      if (res.code === 0 && res.data && res.data.list) {
        // 映射消息格式
        const historyMessages = res.data.list.map(item => ({
          id: item.id,
          content: item.content,
          senderType: item.senderType, // 1-用户消息，2-客服消息
          createTime: item.createTime
        }));
        
        // 按时间顺序排序历史消息
        historyMessages.sort((a, b) => a.createTime - b.createTime);
        
        this.setData({
          messages: historyMessages,
          hasMore: historyMessages.length >= this.data.pageSize
        });
        
        // 滚动到底部
        this.scrollToBottom();
      }
    }).catch(err => {
      console.error('刷新消息失败', err);
    });
  },
  
  // 重发失败消息
  resendMessage: function(e) {
    const index = e.currentTarget.dataset.index;
    const message = this.data.messages[index];
    
    // 从消息列表中移除失败消息
    const updatedMessages = [...this.data.messages];
    updatedMessages.splice(index, 1);
    
    this.setData({
      messages: updatedMessages,
      inputValue: message.content
    });
    
    // 调用发送方法
    this.sendMessage();
  },
  
  // 加载聊天历史记录
  loadChatHistory: function() {
    // 如果没有会话ID，不执行操作
    if (!this.data.sessionId) return;
    
    this.setData({ loadingMore: true });
    
    getChatHistory(this.data.sessionId, this.data.page, this.data.pageSize).then(res => {
      console.log('获取聊天历史成功', res);
      if (res.code === 0 && res.data) {
        // 映射消息格式
        const historyMessages = res.data.list.map(item => ({
          id: item.id,
          content: item.content,
          senderType: item.senderType, // 1-用户消息，2-客服消息
          createTime: item.createTime
        }));
        
        // 按时间顺序排序历史消息
        historyMessages.sort((a, b) => a.createTime - b.createTime);
        
        if (this.data.page === 1) {
          // 如果是第一页，直接设置消息列表
          this.setData({
            messages: historyMessages,
            page: this.data.page + 1,
            hasMore: historyMessages.length >= this.data.pageSize
          });
          
          // 如果没有历史消息，显示欢迎消息
          if (historyMessages.length === 0) {
            this.setData({
              messages: [{
                id: Date.now(),
                content: '您好，我是客服助手，请问有什么可以帮您？',
                senderType: 2, // 客服消息
                createTime: Math.floor(Date.now() / 1000)
              }]
            });
          }
          
          // 滚动到底部
          this.scrollToBottom();
        } else {
          // 如果是加载更多，则将历史消息添加到现有消息前面
          this.setData({
            messages: [...historyMessages, ...this.data.messages],
            page: this.data.page + 1,
            hasMore: historyMessages.length >= this.data.pageSize
          });
        }
      } else {
        wx.showToast({
          title: res.msg || '获取历史记录失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      console.error('获取聊天历史失败', err);
      wx.showToast({
        title: '网络错误，请稍后重试',
        icon: 'none'
      });
    }).finally(() => {
      this.setData({ loadingMore: false });
    });
  },
  
  // 加载更多历史消息
  loadMoreHistory: function() {
    if (this.data.loadingMore || !this.data.hasMore) return;
    console.log('加载更多历史消息');
    this.loadChatHistory();
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