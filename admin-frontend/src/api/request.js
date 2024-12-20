import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const createService = (baseURL) => {
  const service = axios.create({
    baseURL,
    timeout: 5000
  })

  // 请求拦截器
  service.interceptors.request.use(
    config => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers['Authorization'] = `Bearer ${token}`
      }
      return config
    },
    error => {
      return Promise.reject(error)
    }
  )

  // 响应拦截器
  service.interceptors.response.use(
    response => {
      return response.data
    },
    error => {
      if (error.response) {
        switch (error.response.status) {
          case 401:
            ElMessage.error('未授权，请重新登录')
            localStorage.removeItem('token')
            router.push('/login')
            break
          case 403:
            ElMessage.error('拒绝访问')
            break
          case 404:
            ElMessage.error('请求错误，未找到该资源')
            break
          case 500:
            ElMessage.error('服务器错误')
            break
          default:
            ElMessage.error(error.response.data.message || '未知错误')
        }
      }
      return Promise.reject(error)
    }
  )

  return service
}

export const userService = createService('http://localhost:8000')
export const productService = createService('http://localhost:8001')
export const orderService = createService('http://localhost:8002')
export const payService = createService('http://localhost:8003') 