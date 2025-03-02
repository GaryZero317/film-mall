import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建axios实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8000', // api的base_url
  timeout: 15000 // 请求超时时间
})

// request拦截器
service.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token')
    if (token) {
      // 确保添加Bearer前缀
      config.headers['Authorization'] = token.startsWith('Bearer ') ? token : 'Bearer ' + token
      // 添加调试日志
      console.log('发送请求头Authorization:', config.headers['Authorization'])
    } else {
      console.warn('请求未包含认证令牌，可能导致授权错误')
    }
    return config
  },
  error => {
    console.log(error)
    Promise.reject(error)
  }
)

// response拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 如果是文件下载，直接返回
    if (response.config.responseType === 'blob') {
      return response
    }

    // 根据后端约定的状态码判断请求是否成功
    if (res.code !== 0) {
      ElMessage({
        message: res.msg || '请求失败',
        type: 'error',
        duration: 5 * 1000
      })

      // 401: 未登录或token过期
      if (res.code === 401) {
        // 清除token
        localStorage.removeItem('token')
        // 跳转登录页
        router.push('/login')
      }
      return Promise.reject(new Error(res.msg || '请求失败'))
    } else {
      return res
    }
  },
  error => {
    console.log('请求错误：', error)
    
    // 增强错误信息处理
    let errorMessage = '请求失败'
    
    if (error.response) {
      const status = error.response.status
      
      console.log('错误状态码：', status)
      console.log('错误响应：', error.response.data)
      
      // 根据状态码提供更具体的错误信息
      switch (status) {
        case 400:
          errorMessage = '请求参数错误 (400)'
          break
        case 401:
          errorMessage = '未授权，请重新登录 (401)'
          // 清除token
          localStorage.removeItem('token')
          // 跳转到登录页
          router.push('/login')
          break
        case 403:
          errorMessage = '拒绝访问 (403)'
          break
        case 404:
          errorMessage = '请求的资源不存在 (404)'
          break
        case 500:
          errorMessage = '服务器内部错误 (500)'
          break
        default:
          errorMessage = `请求错误 (${status})`
      }
    } else if (error.request) {
      // 请求已发出但未收到响应
      errorMessage = '服务器无响应'
      console.log('请求未收到响应：', error.request)
    } else {
      // 请求配置有误
      errorMessage = '请求配置错误'
      console.log('请求错误：', error.message)
    }
    
    ElMessage({
      message: errorMessage,
      type: 'error',
      duration: 5 * 1000
    })
    
    return Promise.reject(error)
  }
)

export default service 