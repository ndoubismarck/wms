import AuthLogin from '@/views/pages/Auth/AuthLogin.vue'
import AuthLogout from '@/views/pages/Auth/AuthLogout.vue'

export default [
  {
    path: '/login',
    name: 'auth:login',
    meta: {
      title: 'Login',
      breadcrumb: [],
    },
    component: AuthLogin,
  },
  {
    path: '/logout',
    name: 'auth:logout',
    meta: {
      title: 'Logout',
      breadcrumb: [],
    },
    component: AuthLogout,
  },
]
