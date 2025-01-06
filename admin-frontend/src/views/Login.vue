<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2 class="login-title">胶卷商城后台管理系统</h2>
      </template>
      
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
            prefix-icon="User">
          </el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            prefix-icon="Lock"
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
    </el-card>
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
    
    if (res && res.accessToken) {
      userStore.token = res.accessToken
      localStorage.setItem('token', res.accessToken)
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
  background-color: #f0f2f5;
  background-image: linear-gradient(45deg, #1890ff 0%, #1890ff 100%);
}

.login-card {
  width: 400px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-title {
  text-align: center;
  color: #303133;
  margin: 0;
  font-size: 24px;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
  padding: 12px 0;
  font-size: 16px;
}

:deep(.el-input__wrapper) {
  padding: 8px 11px;
}

:deep(.el-input__inner) {
  height: 40px;
}
</style> 