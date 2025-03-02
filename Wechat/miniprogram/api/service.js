import request from '../utils/request'

// 获取常见问题列表
export const getFaqList = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/service/faq/list',
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
export const createQuestion = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/service/create',
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
export const getQuestionList = (data) => {
  return request({
    url: 'http://localhost:8000/api/user/service/list',
    method: 'POST',
    data
  });
}

// 获取问题详情
export const getQuestionDetail = (id) => {
  return request({
    url: 'http://localhost:8000/api/user/service/detail',
    method: 'POST',
    data: { id }
  });
}

// 获取或创建聊天会话
export const getChatSession = () => {
  return request({
    url: 'http://localhost:8000/api/user/chat/session',
    method: 'GET'
  }).then(res => {
    console.log('API返回的聊天会话数据:', res);
    // 转换API返回格式，确保统一
    if (res.code === 200) {
      return {
        code: 0,
        data: res.data,
        msg: res.msg
      };
    }
    return res;
  });
}

// 获取聊天历史记录
export const getChatHistory = (sessionId, page = 1, pageSize = 20) => {
  return request({
    url: 'http://localhost:8000/api/user/chat/history',
    method: 'GET',
    data: { sessionId, page, pageSize }
  }).then(res => {
    console.log('API返回的聊天历史数据:', res);
    // 转换API返回格式，确保统一
    if (res.code === 200) {
      return {
        code: 0,
        data: res.data,
        msg: res.msg
      };
    }
    return res;
  });
}

// 发送聊天消息
export const sendChatMessage = (sessionId, content, type = 1) => {
  return request({
    url: 'http://localhost:8000/api/user/chat/send',
    method: 'POST',
    data: { sessionId, content, type }
  }).then(res => {
    console.log('API返回的发送消息结果:', res);
    // 转换API返回格式，确保统一
    if (res.code === 200) {
      return {
        code: 0,
        data: res.data,
        msg: res.msg
      };
    }
    return res;
  });
} 