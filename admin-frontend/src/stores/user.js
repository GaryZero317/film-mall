import { defineStore } from 'pinia'
import { login } from '../api/admin'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    username: localStorage.getItem('username') || ''
  }),

  getters: {
    isLoggedIn: (state) => !!state.token
  },

  actions: {
    async login(username, password) {
      try {
        const res = await login({ username, password })
        this.token = res.token
        this.username = username
        localStorage.setItem('token', res.token)
        localStorage.setItem('username', username)
        return true
      } catch (error) {
        console.error('登录失败:', error)
        return false
      }
    },

    logout() {
      this.token = ''
      this.username = ''
      localStorage.removeItem('token')
      localStorage.removeItem('username')
    }
  }
}) 