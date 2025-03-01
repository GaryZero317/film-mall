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
        path: 'processing',
        redirect: 'film/list'
      },
      {
        path: 'payments',
        name: 'PaymentList',
        component: () => import('../views/payment/PaymentList.vue')
      },
      {
        path: 'film/list',
        name: 'FilmList',
        component: () => import('../views/film/FilmList.vue')
      },
      {
        path: 'film/detail/:id',
        name: 'FilmDetail',
        component: () => import('../views/film/FilmDetail.vue')
      },
      {
        path: 'customer-service/questions',
        name: 'CustomerServiceQuestions',
        component: () => import('../views/customer-service/QuestionList.vue')
      },
      {
        path: 'customer-service/chat',
        name: 'CustomerServiceChat',
        component: () => import('../views/customer-service/ChatConsole.vue')
      },
      {
        path: 'customer-service/faq',
        name: 'CustomerServiceFaq',
        component: () => import('../views/customer-service/FaqList.vue')
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