import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getAdminInfo } from '../api/admin'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  // 设置token
  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 清除token
  const clearToken = () => {
    token.value = ''
    localStorage.removeItem('token')
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      const res = await getAdminInfo()
      userInfo.value = res
      return res
    } catch (error) {
      console.error('获取用户信息失败:', error)
      return null
    }
  }

  // 退出登录
  const logout = () => {
    userInfo.value = null
    clearToken()
  }

  return {
    userInfo,
    token,
    setToken,
    clearToken,
    fetchUserInfo,
    logout
  }
}) 