import request from '../utils/request'

// 获取服务器地址
const getBaseUrl = () => {
  // 从缓存获取基础URL，如果没有则使用端口8000
  return wx.getStorageSync('baseUrl') || 'http://localhost:8000';
}

// 获取常见问题列表
const getFaqList = (data) => {
  return request({
    url: '/api/user/service/faq/list',
    method: 'POST',
    data
  }).then(res => {
    // 如果后台返回的是标准HTTP成功状态码，但我们需要的是自定义格式
    // 这里进行转换，确保返回格式一致
    if (res.code === 200 && !res.hasOwnProperty('data')) {
      return {
        code: 0,
        data: res
      };
    }
    return res;
  });
}

// 创建客服问题
const createQuestion = (data) => {
  return request({
    url: '/api/user/service/submit',
    method: 'POST',
    data
  }).then(res => {
    // 确保返回格式一致
    if (res.code === 200 && !res.hasOwnProperty('data')) {
      return {
        code: 0,
        data: res
      };
    }
    return res;
  });
}

// 获取用户的问题列表
const getQuestionList = (data) => {
  return request({
    url: '/api/user/service/list',
    method: 'POST',
    data
  });
}

// 获取问题详情
const getQuestionDetail = (id) => {
  return request({
    url: '/api/user/service/detail',
    method: 'POST',
    data: { id }
  });
}

// 获取或创建聊天会话 (模拟实现)
const getChatSession = () => {
  // 本地缓存的会话ID
  let sessionId = wx.getStorageSync('chatSessionId');
  
  // 如果没有会话ID，创建一个新的
  if (!sessionId) {
    sessionId = 'session_' + Date.now();
    wx.setStorageSync('chatSessionId', sessionId);
  }
  
  // 返回一个成功的Promise
  return Promise.resolve({
    code: 0,
    data: {
      sessionId: sessionId,
      createTime: Math.floor(Date.now() / 1000)
    },
    msg: '会话获取成功'
  });
}

// 获取聊天历史记录 (模拟实现)
const getChatHistory = (sessionId, page = 1, pageSize = 20) => {
  // 尝试从本地缓存获取消息
  const key = `chat_history_${sessionId}`;
  let messages = wx.getStorageSync(key) || [];
  
  // 如果是第一页且没有消息，添加一条欢迎消息
  if (page === 1 && messages.length === 0) {
    const welcomeMessage = {
      id: Date.now(),
      content: '您好，我是客服助手，请问有什么可以帮您？',
      senderType: 2, // 客服消息
      createTime: Math.floor(Date.now() / 1000)
    };
    messages.push(welcomeMessage);
    wx.setStorageSync(key, messages);
  }
  
  // 按页码筛选消息
  const startIdx = (page - 1) * pageSize;
  const endIdx = startIdx + pageSize;
  const pagedMessages = messages.slice(startIdx, endIdx);
  
  // 返回成功的Promise
  return Promise.resolve({
    code: 0,
    data: {
      list: pagedMessages,
      total: messages.length,
      page: page,
      pageSize: pageSize
    },
    msg: '获取聊天记录成功'
  });
}

// 发送聊天消息 (模拟实现)
const sendChatMessage = (sessionId, content, type = 1) => {
  // 获取现有消息
  const key = `chat_history_${sessionId}`;
  let messages = wx.getStorageSync(key) || [];
  
  // 创建用户消息
  const userMessage = {
    id: `user_${Date.now()}`,
    content: content,
    senderType: 1, // 用户消息
    createTime: Math.floor(Date.now() / 1000)
  };
  
  // 创建客服自动回复
  const replyOptions = [
    '您好，我们已收到您的消息，正在处理中。',
    '感谢您的咨询，请问还有其他问题吗？',
    '您的问题我们已记录，稍后会有客服专员与您联系。',
    '我们将尽快解决您的问题，请耐心等待。',
    '您的反馈对我们很重要，谢谢支持！'
  ];
  
  const randomReply = replyOptions[Math.floor(Math.random() * replyOptions.length)];
  
  const serviceMessage = {
    id: `service_${Date.now() + 1}`,
    content: randomReply,
    senderType: 2, // 客服消息
    createTime: Math.floor(Date.now() / 1000) + 1
  };
  
  // 将消息添加到历史记录
  messages.push(userMessage);
  
  // 保存到本地缓存
  wx.setStorageSync(key, messages);
  
  // 返回成功的Promise
  return Promise.resolve({
    code: 0,
    data: {
      id: userMessage.id,
      content: content,
      createTime: userMessage.createTime,
      serviceReply: {
        id: serviceMessage.id,
        content: serviceMessage.content,
        createTime: serviceMessage.createTime
      }
    },
    msg: '发送成功'
  }).then(res => {
    // 延迟1秒后添加客服回复
    setTimeout(() => {
      messages.push(serviceMessage);
      wx.setStorageSync(key, messages);
    }, 1000);
    
    return res;
  });
}

// 确保兼容性：同时使用module.exports和export
module.exports = {
  getFaqList,
  createQuestion,
  getQuestionList,
  getQuestionDetail,
  getChatSession,
  getChatHistory,
  sendChatMessage
};

export {
  getFaqList,
  createQuestion,
  getQuestionList,
  getQuestionDetail,
  getChatSession,
  getChatHistory,
  sendChatMessage
}; 