import axios from 'axios'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'
import router from '../router'

// 创建 axios 实例
const createService = (baseURL) => {
  const service = axios.create({
    baseURL,
    timeout: 10000
  })

  // 请求拦截器
  service.interceptors.request.use(
    (config) => {
      const userStore = useUserStore()
      if (userStore.token) {
        config.headers['Authorization'] = 'Bearer ' + userStore.token
      }
      return config
    },
    (error) => {
      console.error('请求错误:', error)
      return Promise.reject(error)
    }
  )

  // 响应拦截器
  service.interceptors.response.use(
    (response) => {
      const res = response.data
      console.log('API响应数据:', res)
      
      // 处理聊天相关API
      if (response.config.url.includes('/api/admin/chat/')) {
        console.log('检测到聊天API调用:', response.config.url)
        
        // 如果响应体中包含list字段，这是典型的列表响应
        if (res.list !== undefined || (res.data && res.data.list !== undefined)) {
          return res
        }
        
        // 对于其他聊天响应，检查是否有错误码
        if (res.code !== undefined && res.code !== 0) {
          console.error('聊天API错误:', res.msg || '请求失败')
          return Promise.reject(new Error(res.msg || '请求失败'))
        }
        
        // 没有明确的错误，返回
        return res
      }
      
      // 针对FAQ列表API的特殊处理
      if (response.config.url.includes('/api/user/service/faq/list')) {
        console.log('检测到FAQ列表API，返回完整响应:', response)
        // 如果响应没有code字段或code为0，则视为成功
        if (!res.code || res.code === 0) {
          return res
        }
        // 有code且不为0，表示有错误
        return Promise.reject(new Error(res.msg || '请求失败'))
      }
      
      // 针对问题详情API的特殊处理
      if (response.config.url.includes('/api/admin/service/detail')) {
        console.log('检测到问题详情API，返回完整响应:', response)
        // 如果响应没有code字段或code为0，则视为成功
        if (!res.code || res.code === 0) {
          return res
        }
        // 有code且不为0，表示有错误
        return Promise.reject(new Error(res.msg || '请求失败'))
      }
      
      // 如果响应中包含 accessToken，说明是登录接口，直接返回数据
      if (res.accessToken !== undefined) {
        return res
      }
      
      // 处理其他接口的响应
      if (res.code === undefined) {
        return res
      }
      
      if (res.code !== 0) {
        const errorMsg = res.msg || '请求失败'
        console.error('API错误:', errorMsg)
        return Promise.reject(new Error(errorMsg))
      }
      return res
    },
    (error) => {
      console.error('响应错误:', error)
      if (error.response) {
        const status = error.response.status
        const data = error.response.data
        console.error('错误状态码:', status)
        console.error('错误响应数据:', data)
        
        switch (status) {
          case 401:
            ElMessage.error('未授权，请重新登录')
            const userStore = useUserStore()
            userStore.logout()
            router.push('/login')
            break
          case 403:
            ElMessage.error('拒绝访问')
            break
          case 404:
            ElMessage.error('请求的资源不存在')
            break
          case 500:
            if (typeof data === 'string' && data.includes('管理员不存在')) {
              ElMessage.error('用户名不存在')
            } else if (typeof data === 'string' && data.includes('密码错误')) {
              ElMessage.error('密码错误')
            } else {
              ElMessage.error(data || '服务器内部错误')
            }
            break
          default:
            ElMessage.error(data || `请求失败(${status})`)
        }
      } else if (error.request) {
        console.error('请求未收到响应:', error.request)
        ElMessage.error('网络错误，请检查网络连接')
      } else {
        console.error('请求配置错误:', error.message)
        ElMessage.error(error.message || '请求发送失败')
      }
      return Promise.reject(error)
    }
  )

  return service
}

// 创建各个服务的实例
export const adminService = createService('http://localhost:8000')
export const productService = createService('http://localhost:8001')
export const orderService = createService('http://localhost:8002')
export const paymentService = createService('http://localhost:8003')
export const statisticsService = createService('http://localhost:8006')
export const filmService = createService('http://localhost:8007') 