<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '240px'" class="aside-container">
      <div class="logo-container" :class="{ 'collapsed': isCollapse }">
        <img src="../assets/logo.svg" alt="Logo" class="logo">
        <span class="title" v-show="!isCollapse">胶卷商城管理</span>
      </div>
      
      <div class="menu-wrapper">
        <el-menu
          :default-active="route.path"
          class="el-menu-vertical"
          :collapse="isCollapse"
          :collapse-transition="true"
          :router="true">
          <el-menu-item index="/">
            <el-icon><Odometer /></el-icon>
            <template #title>数据概览</template>
          </el-menu-item>
          
          <el-menu-item index="/admins">
            <el-icon><User /></el-icon>
            <template #title>管理员管理</template>
          </el-menu-item>
          
          <el-menu-item index="/products">
            <el-icon><Goods /></el-icon>
            <template #title>商品管理</template>
          </el-menu-item>
          
          <el-menu-item index="/orders">
            <el-icon><List /></el-icon>
            <template #title>订单管理</template>
          </el-menu-item>
          
          <el-menu-item index="/film/list">
            <el-icon><Picture /></el-icon>
            <template #title>冲洗管理</template>
          </el-menu-item>
          
          <el-menu-item index="/payments">
            <el-icon><Money /></el-icon>
            <template #title>支付管理</template>
          </el-menu-item>
        </el-menu>
      </div>
      
      <div class="sidebar-footer" v-if="!isCollapse">
        <div class="user-info">
          <el-avatar size="small" icon="el-icon-user"></el-avatar>
          <span class="username">{{ userStore.username }}</span>
        </div>
      </div>
    </el-aside>

    <el-container class="main-container">
      <el-header>
        <div class="header-left">
          <div 
            class="collapse-btn"
            @click="toggleCollapse">
            <el-icon>
              <Fold v-if="!isCollapse"/>
              <Expand v-else/>
            </el-icon>
          </div>
          <breadcrumb />
        </div>
        
        <div class="header-right">
          <div class="action-item">
            <el-tooltip content="全屏" placement="bottom">
              <el-icon @click="toggleFullScreen"><FullScreen /></el-icon>
            </el-tooltip>
          </div>
          
          <div class="action-item">
            <el-tooltip content="刷新" placement="bottom">
              <el-icon @click="refreshPage"><Refresh /></el-icon>
            </el-tooltip>
          </div>
          
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-avatar size="small" icon="el-icon-user"></el-avatar>
              {{ userStore.username }}
              <el-icon class="el-icon--right"><CaretBottom /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="settings">系统设置</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main>
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
      
      <el-footer height="40px">
        <div class="footer-content">
          © 2023 胶卷商城后台管理系统 版权所有
        </div>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { 
  Odometer, 
  User, 
  Goods, 
  List, 
  Money, 
  Fold, 
  Expand,
  CaretBottom,
  Picture,
  FullScreen,
  Refresh
} from '@element-plus/icons-vue'
import Breadcrumb from './components/Breadcrumb.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isCollapse = ref(false)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else if (command === 'profile') {
    // 实现跳转到个人信息页面的逻辑
  } else if (command === 'settings') {
    // 实现跳转到系统设置页面的逻辑
  }
}

const toggleFullScreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen()
    }
  }
}

const refreshPage = () => {
  window.location.reload()
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
  --sidebar-bg: #001529;
  --sidebar-text: rgba(255, 255, 255, 0.65);
  --sidebar-active-text: #fff;
  --sidebar-active-bg: #1890ff;
  --header-height: 64px;
  --sidebar-logo-height: 64px;
}

.aside-container {
  transition: width 0.3s;
  overflow: hidden;
  background-color: var(--sidebar-bg);
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
}

.logo-container {
  height: var(--sidebar-logo-height);
  display: flex;
  align-items: center;
  padding: 0 24px;
  color: white;
  background-color: var(--sidebar-bg);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-container.collapsed {
  padding: 0 20px;
  justify-content: center;
}

.logo {
  height: 32px;
  min-width: 32px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  margin-left: 12px;
  white-space: nowrap;
  color: white;
}

.menu-wrapper {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.el-menu-vertical {
  border-right: none;
  background-color: var(--sidebar-bg);
}

.el-menu-vertical:not(.el-menu--collapse) {
  width: 240px;
}

.el-menu-item {
  height: 50px;
  line-height: 50px;
}

:deep(.el-menu-item) {
  color: var(--sidebar-text);
}

:deep(.el-menu-item.is-active) {
  color: var(--sidebar-active-text);
  background-color: var(--sidebar-active-bg);
}

:deep(.el-menu-item:hover) {
  background-color: rgba(255, 255, 255, 0.05);
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  color: var(--sidebar-text);
}

.username {
  margin-left: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.main-container {
  display: flex;
  flex-direction: column;
}

.el-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--header-height);
  padding: 0 24px;
  background-color: white;
  border-bottom: 1px solid #f0f0f0;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  cursor: pointer;
  margin-right: 20px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  transition: background-color 0.3s;
}

.collapse-btn:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.header-right {
  display: flex;
  align-items: center;
}

.action-item {
  padding: 0 12px;
  cursor: pointer;
  font-size: 20px;
  color: #606266;
  transition: color 0.3s;
}

.action-item:hover {
  color: #1890ff;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
  height: 50px;
  transition: background-color 0.3s;
}

.user-dropdown .el-avatar {
  margin-right: 8px;
}

.user-dropdown:hover {
  background-color: rgba(0, 0, 0, 0.025);
}

.el-main {
  padding: 24px;
  background-color: #f0f2f5;
  overflow-y: auto;
  min-height: 0;
  flex: 1;
}

.el-footer {
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: white;
  border-top: 1px solid #f0f0f0;
  color: #606266;
  font-size: 14px;
}

/* 路由过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style> 