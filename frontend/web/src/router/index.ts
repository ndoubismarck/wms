import authRoutes from './auth.ts'
import setupRoutes from './setup.ts'
import adminRoutes from './admin.ts'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/admin',
    },
    ...authRoutes,
    ...setupRoutes,
    ...adminRoutes,
  ],
})
export default router
