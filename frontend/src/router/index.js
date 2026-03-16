import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import MonitorDetailView from '../views/MonitorDetailView.vue'
import MonitorFormView from '../views/MonitorFormView.vue'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { public: true },
  },
  {
    path: '/',
    name: 'dashboard',
    component: DashboardView,
  },
  {
    path: '/monitors/new',
    name: 'monitor-create',
    component: MonitorFormView,
  },
  {
    path: '/monitors/:id',
    name: 'monitor-detail',
    component: MonitorDetailView,
  },
  {
    path: '/monitors/:id/edit',
    name: 'monitor-edit',
    component: MonitorFormView,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) {
    return { name: 'login' }
  }
})

export default router
