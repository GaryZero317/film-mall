<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-header">
        <h2 class="login-title">胶卷商城后台管理系统</h2>
        <p class="login-subtitle">欢迎回来，请登录您的账号</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="0"
        class="login-form">
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            prefix-icon="User"
            class="custom-input">
          </el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            prefix-icon="Lock"
            class="custom-input"
            @keyup.enter="handleLogin">
          </el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="login-button"
            @click="handleLogin">
            登录
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <p>© 2023 胶卷商城 版权所有</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { login } from '../api/admin'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const loading = ref(false)
const loginFormRef = ref(null)

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    const res = await login(loginForm)
    console.log('登录响应:', res)
    
    if (res && res.code === 0) {
      userStore.token = res.data.token
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('adminInfo', JSON.stringify(res.data))
      localStorage.setItem('username', loginForm.username)
      ElMessage.success('登录成功')
      await router.push('/')
    } else {
      console.error('登录失败，响应数据:', res)
      const errorMsg = res.msg || ''
      if (errorMsg.includes('密码错误')) {
        ElMessage.error('密码错误，请重新输入')
      } else if (errorMsg.includes('管理员不存在')) {
        ElMessage.error('用户名不存在，请检查输入')
      } else {
        ElMessage.error(errorMsg || '登录失败，请检查用户名和密码')
      }
    }
  } catch (error) {
    console.error('登录错误:', error)
    console.error('错误详情:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    })
    
    if (error.message === '请输入密码') {
      ElMessage.error('请输入密码')
    } else if (error.response?.data) {
      const errorData = typeof error.response.data === 'string' 
        ? error.response.data 
        : error.response.data.msg || error.response.data.error || ''
      
      if (errorData.includes('管理员不存在')) {
        ElMessage.error('用户名不存在，请检查输入')
      } else if (errorData.includes('密码错误')) {
        ElMessage.error('密码错误，请重新输入')
      } else {
        ElMessage.error(errorData || '登录失败，请稍后重试')
      }
    } else {
      ElMessage.error('登录失败，请检查网络连接')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #2979FF 0%, #0D47A1 100%);
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  width: 200%;
  height: 200%;
  top: -50%;
  left: -50%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 60%);
  transform: rotate(30deg);
}

.login-content {
  width: 420px;
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.2);
  padding: 40px;
  backdrop-filter: blur(5px);
  animation: fadeIn 0.6s ease-out;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  color: #1a1a1a;
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.login-subtitle {
  margin-top: 12px;
  color: #666;
  font-size: 16px;
}

.login-form {
  margin-bottom: 24px;
}

.login-button {
  width: 100%;
  padding: 12px 0;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  background: linear-gradient(90deg, #2979FF 0%, #1E88E5 100%);
  border: none;
  transition: transform 0.2s, box-shadow 0.2s;
  margin-top: 12px;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(33, 150, 243, 0.25);
}

.login-footer {
  text-align: center;
  color: #999;
  font-size: 14px;
  margin-top: 24px;
}

.custom-input :deep(.el-input__wrapper) {
  padding: 12px 16px;
  border-radius: 8px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
}

.custom-input :deep(.el-input__inner) {
  height: 42px;
  font-size: 15px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style> 