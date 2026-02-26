import authRoutes from './auth.ts'
import setupRoutes from './setup.ts'
import adminRoutes from './admin.ts'
import pagesRoutes from './pages.ts'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [...authRoutes, ...setupRoutes, ...adminRoutes, ...pagesRoutes],
})
export default router
