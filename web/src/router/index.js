import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../layout/index.vue'

const routes = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/index.vue'),
        meta: { title: '控制台' }
      }
    ]
  },
  {
    path: '/signals',
    component: Layout,
    redirect: '/signals/list',
    children: [
      {
        path: 'list',
        name: 'SignalList',
        component: () => import('../views/signals/index.vue'),
        meta: { title: '信号列表' }
      },
      {
        path: 'detail/:id',
        name: 'SignalDetail',
        component: () => import('../views/signals/detail.vue'),
        meta: { title: '信号详情' }
      }
    ]
  },
  {
    path: '/orders',
    component: Layout,
    redirect: '/orders/list',
    children: [
      {
        path: 'list',
        name: 'OrderList',
        component: () => import('../views/orders/index.vue'),
        meta: { title: '订单列表' }
      },
      {
        path: 'detail/:id',
        name: 'OrderDetail',
        component: () => import('../views/orders/detail.vue'),
        meta: { title: '订单详情' }
      }
    ]
  },
  {
    path: '/contract-codes',
    component: Layout,
    redirect: '/contract-codes/list',
    children: [
      {
        path: 'list',
        name: 'ContractCodeList',
        component: () => import('../views/contract-codes/index.vue'),
        meta: { title: '交易对管理' }
      }
    ]
  },
  {
    path: '/strategies',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Strategies',
        component: () => import('../views/strategies/index.vue'),
        meta: { title: '策略管理' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router