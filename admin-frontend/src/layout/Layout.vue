<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '200px'" class="aside-container">
      <div class="logo-container" :class="{ 'collapsed': isCollapse }">
        <img src="../assets/logo.svg" alt="Logo" class="logo">
        <span class="title" v-show="!isCollapse">胶卷商城管理</span>
      </div>
      
      <el-menu
        :default-active="route.path"
        class="el-menu-vertical"
        :collapse="isCollapse"
        :collapse-transition="true"
        :router="true">
        <el-menu-item index="/">
          <el-icon><Odometer /></el-icon>
          <template #title>首页</template>
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
    </el-aside>

    <el-container>
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
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              {{ userStore.username }}
              <el-icon class="el-icon--right"><CaretBottom /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
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
  Picture
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
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside-container {
  transition: width 0.3s;
  overflow: hidden;
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  color: var(--el-menu-text-color);
  background-color: var(--el-menu-bg-color);
  transition: all 0.3s;
  overflow: hidden;
}

.logo-container.collapsed {
  padding: 0 16px;
}

.logo {
  height: 32px;
  margin-right: 12px;
  min-width: 32px;
}

.title {
  font-size: 16px;
  font-weight: bold;
  white-space: nowrap;
  transition: opacity 0.3s;
}

.el-aside {
  background-color: var(--el-menu-bg-color);
  border-right: 1px solid var(--el-border-color-light);
}

.el-menu-vertical {
  border-right: none;
}

.el-menu-vertical:not(.el-menu--collapse) {
  width: 200px;
}

.el-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 20px;
  border-bottom: 1px solid var(--el-border-color-light);
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
  border-radius: 4px;
  transition: background-color 0.3s;
}

.collapse-btn:hover {
  background-color: var(--el-fill-color-light);
}

.header-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
  height: 40px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-dropdown:hover {
  background-color: var(--el-fill-color-light);
}

.el-main {
  padding: 20px;
  background-color: var(--el-bg-color-page);
}

/* 路由过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style> 