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
export const filmService = createService('http://localhost:8007') 