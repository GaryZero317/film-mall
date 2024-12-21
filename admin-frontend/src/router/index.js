import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('../layout/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue')
      },
      {
        path: 'admins',
        name: 'AdminList',
        component: () => import('../views/admin/AdminList.vue')
      },
      {
        path: 'products',
        name: 'ProductList',
        component: () => import('../views/product/ProductList.vue')
      },
      {
        path: 'orders',
        name: 'OrderList',
        component: () => import('../views/order/OrderList.vue')
      },
      {
        path: 'payments',
        name: 'PaymentList',
        component: () => import('../views/payment/PaymentList.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth) {
    if (!token) {
      next('/login')
    } else {
      next()
    }
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router 